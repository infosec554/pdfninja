package gotenberg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Client interface {
	PDFToWord(ctx context.Context, pdfPath string) ([]byte, error)
	WordToPDF(ctx context.Context, wordPath string) ([]byte, error)
}

type gotenbergClient struct {
	baseURL string
}

func New(url string) Client {
	return &gotenbergClient{baseURL: url}
}

// PDF -> Word
func (g *gotenbergClient) PDFToWord(ctx context.Context, pdfPath string) ([]byte, error) {
	file, err := os.Open(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("files", filepath.Base(pdfPath))
	if err != nil {
		return nil, fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("cannot copy file: %w", err)
	}

	_ = writer.WriteField("output", "docx") // optional
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/libreoffice/convert", &requestBody)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("conversion failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("conversion failed: %s", string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

// Word -> PDF
func (g *gotenbergClient) WordToPDF(ctx context.Context, wordPath string) ([]byte, error) {
	file, err := os.Open(wordPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("files", filepath.Base(wordPath))
	if err != nil {
		return nil, fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("cannot copy file: %w", err)
	}

	_ = writer.WriteField("waitTimeout", "30s")
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/libreoffice/convert", &requestBody)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("conversion failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("conversion failed: %s", string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}
