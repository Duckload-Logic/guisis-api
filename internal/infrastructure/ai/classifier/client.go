package classifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
)

const defaultClassifierTimeout = 30 * time.Second

type ClassifierClient struct {
	httpClient *http.Client
	baseUrl    string
	apiKey     string
}

// NewClient creates a new classifier client
func NewClient(
	httpClient *http.Client,
	baseUrl string,
	apiKey string,
) *ClassifierClient {
	return &ClassifierClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
		apiKey:     apiKey,
	}
}

func (c *ClassifierClient) Classify(
	ctx context.Context,
	text string,
	cfg *config.Config,
) (*ClassifyResponse, error) {
	if strings.TrimSpace(c.baseUrl) == "" {
		return nil, fmt.Errorf(
			"[ClassifierClient] {Validate}: baseUrl is empty",
		)
	}

	payload := ClassifyRequest{
		Text: text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("[ClassifierClient] {Marshal JSON}: %w", err)
	}

	timeout := defaultClassifierTimeout
	if c.httpClient != nil && c.httpClient.Timeout > 0 {
		timeout = c.httpClient.Timeout
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		timeoutCtx,
		http.MethodPost,
		fmt.Sprintf("%s/classify", c.baseUrl),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("[ClassifierClient] {Create Request}: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[ClassifierClient] {Do Request}: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"[ClassifierClient] {Status Code}: %d",
			resp.StatusCode,
		)
	}

	var response ClassifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("[ClassifierClient] {Decode Response}: %w", err)
	}

	return &response, nil
}
