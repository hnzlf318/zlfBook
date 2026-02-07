package ocr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/models"
)

func formatBillListDateToLongDateTime(year, month, day, hour, min int) string {
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:00", year, month, day, hour, min)
}

// ParseBillListText parses OCR text from a bill/transaction list screenshot (e.g. 账单)
// and returns a list of recognized transaction results.
// Date format expected: "M月D日 HH:MM" (e.g. "2月7日 21:49"). Year is taken from refTime.
func ParseBillListText(text string, refTime time.Time) []*models.RecognizedReceiptImageResult {
	year := refTime.Year()
	lines := strings.Split(text, "\n")
	var results []*models.RecognizedReceiptImageResult

	// Amount at end of line: -123.45 or 123.45
	amountRe := regexp.MustCompile(`(-?\d+\.?\d*)\s*$`)
	// Chinese date time: 2月7日 21:49 or 12月31日 09:00
	dateRe := regexp.MustCompile(`(\d{1,2})月(\d{1,2})日\s*(\d{1,2}):(\d{2})`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		amountIdx := amountRe.FindStringSubmatchIndex(line)
		if amountIdx == nil {
			continue
		}
		amountStr := strings.TrimSpace(line[amountIdx[2]:amountIdx[3]])
		rest := strings.TrimSpace(line[:amountIdx[0]])
		dateIdx := dateRe.FindStringSubmatchIndex(rest)
		if dateIdx == nil {
			continue
		}
		month, _ := strconv.Atoi(rest[dateIdx[2]:dateIdx[3]])
		day, _ := strconv.Atoi(rest[dateIdx[4]:dateIdx[5]])
		hour, _ := strconv.Atoi(rest[dateIdx[6]:dateIdx[7]])
		min, _ := strconv.Atoi(rest[dateIdx[8]:dateIdx[9]])
		description := strings.TrimSpace(rest[:dateIdx[0]])
		if description == "" {
			description = rest[dateIdx[1]:]
			description = strings.TrimSpace(description)
		}
		if description == "" {
			description = "OCR"
		}
		// Build long date time YYYY-MM-DD HH:mm:00
		longDateTime := formatBillListDateToLongDateTime(year, month, day, hour, min)
		// Default type: expense if amount is negative
		txType := "expense"
		if len(amountStr) > 0 && amountStr[0] == '+' {
			txType = "income"
		} else if len(amountStr) > 0 && amountStr[0] != '-' {
			txType = "income"
		}
		results = append(results, &models.RecognizedReceiptImageResult{
			Type:        txType,
			Time:        longDateTime,
			Amount:      amountStr,
			Description: description,
		})
	}
	return results
}
