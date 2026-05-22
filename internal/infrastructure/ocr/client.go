package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type OCRClient struct {
	baseURL string
	http    *http.Client
	apiKey  string
}

func NewClient(baseURL string, apiKey string) *OCRClient {
	return &OCRClient{
		baseURL: baseURL,
		http:    &http.Client{},
		apiKey:  apiKey,
	}
}

// ProcessCOR sends a file to the specialized COR endpoint.
func (o *OCRClient) ProcessCOR(
	ctx context.Context,
	filename string,
	file io.Reader,
) (*CORResponse, error) {
	respBody, err := o.upload(ctx, "/ocr/cor", filename, file)
	if err != nil {
		return nil, fmt.Errorf("[OCRClient] {ProcessCOR}: %w", err)
	}
	defer respBody.Close()

	var result CORResponse
	if err := json.NewDecoder(respBody).Decode(&result); err != nil {
		return nil, fmt.Errorf("[OCRClient] {Decode COR}: %w", err)
	}

	return &result, nil
}

// ProcessDocument sends a file to the generic OCR endpoint.
func (o *OCRClient) ProcessDocument(
	ctx context.Context,
	filename string,
	file io.Reader,
) (*OCRResponse, error) {
	respBody, err := o.upload(ctx, "/ocr", filename, file)
	if err != nil {
		return nil, fmt.Errorf("[OCRClient] {ProcessDocument}: %w", err)
	}
	defer respBody.Close()

	var result OCRResponse
	if err := json.NewDecoder(respBody).Decode(&result); err != nil {
		return nil, fmt.Errorf("[OCRClient] {Decode Document}: %w", err)
	}

	return &result, nil
}

// upload is a helper to perform multipart uploads to the AI service.
func (o *OCRClient) upload(
	ctx context.Context,
	path string,
	filename string,
	file io.Reader,
) (io.ReadCloser, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filename))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s%s", o.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", o.apiKey)

	resp, err := o.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		return nil, fmt.Errorf("AI service returned status: %s", resp.Status)
	}

	return resp.Body, nil
}
