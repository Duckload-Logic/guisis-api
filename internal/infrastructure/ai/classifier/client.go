package classifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
)

type ClassifierClient struct {
	httpClient *http.Client
	baseUrl    string
}

// NewIDPClient creates a new IDP client with configured timeout
func NewClient(httpClient *http.Client, baseUrl string) *ClassifierClient {
	return &ClassifierClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (c *ClassifierClient) Classify(
	ctx context.Context,
	text string,
	cfg *config.Config,
) (*ClassifyResponse, error) {
	payload := ClassifyRequest{
		Text: text,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("[ClassifierClient] {Marshal JSON}: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/classify", c.baseUrl),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("[ClassifierClient] {Create Request}: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

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
