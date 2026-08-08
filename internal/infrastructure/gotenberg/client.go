package gotenberg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// Client handles communication with a Gotenberg instance.
type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new Gotenberg client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ConvertHTML sends HTML content to Gotenberg's Chromium conversion endpoint.
func (c *Client) ConvertHTML(
	ctx context.Context,
	htmlContent string,
) ([]byte, error) {
	// Gotenberg expects 'index.html' in a multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.WriteString(part, htmlContent); err != nil {
		return nil, fmt.Errorf("failed to write html content: %w", err)
	}

	if err := writer.WriteField("preferCssPageSize", "true"); err != nil {
		return nil, fmt.Errorf(
			"failed to write field preferCssPageSize: %w",
			err,
		)
	}

	if err := writer.WriteField("waitDelay", "1s"); err != nil {
		return nil, fmt.Errorf(
			"failed to write field waitDelay: %w",
			err,
		)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	reqURL := fmt.Sprintf("%s/forms/chromium/convert/html", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gotenberg request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"gotenberg returned error (%d): %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	pdfBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read gotenberg response: %w", err)
	}

	return pdfBytes, nil
}
