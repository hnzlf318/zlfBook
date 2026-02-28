package services

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/mail"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// BackupService represents backup related service
type BackupService struct {
	ServiceUsingConfig
	ServiceUsingMailer
}

// Initialize a backup service singleton instance
var (
	Backup = &BackupService{
		ServiceUsingConfig: ServiceUsingConfig{
			container: settings.Container,
		},
		ServiceUsingMailer: ServiceUsingMailer{
			container: mail.Container,
		},
	}
)

// SendDailyEmailBackup creates database/config backup and sends it via email
func (s *BackupService) SendDailyEmailBackup(c *core.CronContext) error {
	config := s.CurrentConfig()

	if config == nil {
		log.Warnf(c, "[backup.SendDailyEmailBackup] current config is nil, skip backup")
		return nil
	}

	if !config.EnableDailyEmailBackup {
		// feature disabled, nothing to do
		return nil
	}

	if !config.EnableSMTP || config.SMTPConfig == nil {
		log.Warnf(c, "[backup.SendDailyEmailBackup] SMTP is not enabled, skip backup email")
		return nil
	}

	if config.DailyEmailBackupToAddress == "" {
		log.Warnf(c, "[backup.SendDailyEmailBackup] backup target email is empty, skip backup email")
		return nil
	}

	if config.DatabaseConfig.DatabasePath == "" {
		log.Warnf(c, "[backup.SendDailyEmailBackup] database path is empty, skip backup email")
		return nil
	}

	if config.ConfigFilePath == "" {
		log.Warnf(c, "[backup.SendDailyEmailBackup] config file path is empty, skip backup email")
		return nil
	}

	dbPath := config.DatabaseConfig.DatabasePath
	cfgPath := config.ConfigFilePath

	dbInfo, err := os.Stat(dbPath)

	if err != nil {
		log.Errorf(c, "[backup.SendDailyEmailBackup] database file \"%s\" not found, because %s", dbPath, err.Error())
		return nil
	}

	if _, err := os.Stat(cfgPath); err != nil {
		log.Errorf(c, "[backup.SendDailyEmailBackup] config file \"%s\" not found, because %s", cfgPath, err.Error())
		return nil
	}

	backupRoot := filepath.Join(config.WorkingPath, "backup")

	if err := os.MkdirAll(backupRoot, 0o755); err != nil {
		log.Errorf(c, "[backup.SendDailyEmailBackup] cannot create backup directory \"%s\", because %s", backupRoot, err.Error())
		return err
	}

	// Check whether database has changed since last successful backup
	metaPath := filepath.Join(backupRoot, "last_backup_meta.json")

	changed, err := isDatabaseChangedSinceLastBackup(metaPath, dbPath, dbInfo)

	if err != nil {
		log.Warnf(c, "[backup.SendDailyEmailBackup] cannot read last backup meta, will still create backup, because %s", err.Error())
	} else if !changed {
		log.Infof(c, "[backup.SendDailyEmailBackup] database file \"%s\" has not been changed since last backup, skip backup and email", dbPath)
		return nil
	}

	timestamp := time.Now().Format("20060102-150405")

	dbFileName := filepath.Base(dbPath)
	dbNameWithoutExt := strings.TrimSuffix(dbFileName, filepath.Ext(dbFileName))
	zipFileName := fmt.Sprintf("%s_%s.zip", dbNameWithoutExt, timestamp)
	zipFilePath := filepath.Join(backupRoot, zipFileName)

	// create a temporary directory to hold copies before zipping
	tempDir, err := os.MkdirTemp(backupRoot, "tmp_backup_")

	if err != nil {
		log.Errorf(c, "[backup.SendDailyEmailBackup] cannot create temp directory in \"%s\", because %s", backupRoot, err.Error())
		return err
	}

	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	dbCopyPath := filepath.Join(tempDir, filepath.Base(dbPath))
	cfgCopyPath := filepath.Join(tempDir, filepath.Base(cfgPath))

	if err := copyFile(dbPath, dbCopyPath); err != nil {
		log.Errorf(c, "[backup.SendDailyEmailBackup] cannot copy database file from \"%s\" to \"%s\", because %s", dbPath, dbCopyPath, err.Error())
		return err
	}

	if err := copyFile(cfgPath, cfgCopyPath); err != nil {
		log.Errorf(c, "[backup.SendDailyEmailBackup] cannot copy config file from \"%s\" to \"%s\", because %s", cfgPath, cfgCopyPath, err.Error())
		return err
	}

	if err := createZipFromFiles(zipFilePath, []string{dbCopyPath, cfgCopyPath}); err != nil {
		log.Errorf(c, "[backup.SendDailyEmailBackup] cannot create zip file \"%s\", because %s", zipFilePath, err.Error())
		return err
	}

	log.Infof(c, "[backup.SendDailyEmailBackup] backup zip file \"%s\" has been created", zipFilePath)

	subject := fmt.Sprintf("ezBookkeeping 数据备份 %s", timestamp)
	body := fmt.Sprintf("这是 ezBookkeeping 在 %s 生成的自动备份，请妥善保存。", time.Now().Format("2006-01-02 15:04:05"))

	message := &mail.MailMessage{
		To:          config.DailyEmailBackupToAddress,
		Subject:     subject,
		Body:        body,
		Attachments: []string{zipFilePath},
	}

	if err := s.SendMail(message); err != nil {
		log.Errorf(c, "[backup.SendDailyEmailBackup] failed to send backup email to \"%s\", because %s", config.DailyEmailBackupToAddress, err.Error())
		return err
	}

	log.Infof(c, "[backup.SendDailyEmailBackup] backup email has been sent to \"%s\"", config.DailyEmailBackupToAddress)

	// Save latest backup meta for change detection
	if err := saveLastBackupMeta(metaPath, dbPath, dbInfo); err != nil {
		log.Warnf(c, "[backup.SendDailyEmailBackup] cannot write last backup meta file \"%s\", because %s", metaPath, err.Error())
	}

	// Auto clean old backup files
	if config.DailyEmailBackupRetentionDays > 0 {
		if err := cleanupOldBackups(backupRoot, int(config.DailyEmailBackupRetentionDays)); err != nil {
			log.Warnf(c, "[backup.SendDailyEmailBackup] cannot cleanup old backups in \"%s\", because %s", backupRoot, err.Error())
		}
	}

	return nil
}

type backupMeta struct {
	DatabasePath string `json:"database_path"`
	DatabaseSize int64  `json:"database_size"`
	DatabaseMTime int64 `json:"database_mtime"`
}

func isDatabaseChangedSinceLastBackup(metaPath string, dbPath string, dbInfo os.FileInfo) (bool, error) {
	currentMeta := &backupMeta{
		DatabasePath: dbPath,
		DatabaseSize: dbInfo.Size(),
		DatabaseMTime: dbInfo.ModTime().Unix(),
	}

	file, err := os.Open(metaPath)

	if err != nil {
		// meta file not exists or cannot open means we treat database as changed
		return true, nil
	}

	defer file.Close()

	var lastMeta backupMeta

	if err := json.NewDecoder(file).Decode(&lastMeta); err != nil {
		return true, err
	}

	if lastMeta.DatabasePath != currentMeta.DatabasePath {
		return true, nil
	}

	if lastMeta.DatabaseSize != currentMeta.DatabaseSize {
		return true, nil
	}

	if lastMeta.DatabaseMTime != currentMeta.DatabaseMTime {
		return true, nil
	}

	return false, nil
}

func saveLastBackupMeta(metaPath string, dbPath string, dbInfo os.FileInfo) error {
	meta := &backupMeta{
		DatabasePath: dbPath,
		DatabaseSize: dbInfo.Size(),
		DatabaseMTime: dbInfo.ModTime().Unix(),
	}

	tempPath := metaPath + ".tmp"

	file, err := os.Create(tempPath)

	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(meta)

	closeErr := file.Close()

	if err != nil {
		_ = os.Remove(tempPath)
		return err
	}

	if closeErr != nil {
		_ = os.Remove(tempPath)
		return closeErr
	}

	return os.Rename(tempPath, metaPath)
}

func cleanupOldBackups(backupRoot string, retentionDays int) error {
	entries, err := os.ReadDir(backupRoot)

	if err != nil {
		return err
	}

	if retentionDays <= 0 {
		return nil
	}

	now := time.Now()
	cutoff := now.AddDate(0, 0, -retentionDays)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".zip" {
			continue
		}

		fullPath := filepath.Join(backupRoot, entry.Name())
		info, err := entry.Info()

		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			_ = os.Remove(fullPath)
		}
	}

	return nil
}

func copyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)

	if err != nil {
		return err
	}

	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)

	if err != nil {
		return err
	}

	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return dstFile.Sync()
}

func createZipFromFiles(zipPath string, filePaths []string) error {
	zipFile, err := os.Create(zipPath)

	if err != nil {
		return err
	}

	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	defer zipWriter.Close()

	for _, filePath := range filePaths {
		if filePath == "" {
			continue
		}

		if err := addFileToZip(zipWriter, filePath); err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return err
	}

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	header, err := zip.FileInfoHeader(fileInfo)

	if err != nil {
		return err
	}

	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)

	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}


