package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/mayswind/ezbookkeeping/pkg/log"
)

// PaddleBillOCRResponse represents the response body of the PaddleOCR HTTP service.
// This is intentionally simple: the service should at least return { "success": true, "text": "..." }.
type PaddleBillOCRResponse struct {
	Success bool   `json:"success"`
	Text    string `json:"text"`
	Error   string `json:"error,omitempty"`
}

// RunPaddleBillOCR sends the image data to an external PaddleOCR HTTP service and returns the recognized text.
// The external service is responsible for doing OCR and returning a JSON body like:
//
//   {
//     "success": true,
//     "text": "2月7日 21:49 京东超市 -100.00\n2月7日 22:10 余额宝收益 +5.23\n"
//   }
//
// Only the "text" field is used by ezBookkeeping; it will be parsed by ParseBillListText to generate transactions.
func RunPaddleBillOCR(imageData []byte, endpoint string) (string, error) {
	if len(imageData) == 0 {
		return "", fmt.Errorf("image data is empty")
	}
	if endpoint == "" {
		return "", fmt.Errorf("paddle ocr endpoint is not configured")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("image", "bill.jpg")
	if err != nil {
		return "", fmt.Errorf("create form file: %w", err)
	}

	if _, err := part.Write(imageData); err != nil {
		return "", fmt.Errorf("write image data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, &body)
	if err != nil {
		return "", fmt.Errorf("create http request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Warnf(nil, "[ocr.RunPaddleBillOCR] request paddle ocr endpoint failed: %v", err)
		return "", fmt.Errorf("request paddle ocr endpoint failed: %w", err)
	}
	defer resp.Body.Close()

	var ocrResp PaddleBillOCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&ocrResp); err != nil {
		log.Warnf(nil, "[ocr.RunPaddleBillOCR] decode paddle ocr response failed: %v", err)
		return "", fmt.Errorf("decode paddle ocr response failed: %w", err)
	}

	if !ocrResp.Success {
		if ocrResp.Error != "" {
			return "", fmt.Errorf("paddle ocr failed: %s", ocrResp.Error)
		}
		return "", fmt.Errorf("paddle ocr failed without error message")
	}

	return strings.TrimSpace(ocrResp.Text), nil
}
