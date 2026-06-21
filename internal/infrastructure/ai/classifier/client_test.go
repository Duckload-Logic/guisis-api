package classifier

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
)

func TestClassifierClient_Classify(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		// Assert method and headers
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("X-API-KEY") != "test-key" {
			t.Errorf("expected X-API-KEY to be test-key")
		}

		var req ClassifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Text != "severe anxiety and stress" {
			t.Errorf("unexpected text: %s", req.Text)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := ClassifyResponse{
			Level:      "HIGH",
			Confidence: 0.92,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(http.DefaultClient, server.URL, "test-key")
	cfg := &config.Config{}

	res, err := client.Classify(
		context.Background(),
		"severe anxiety and stress",
		cfg,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Level != "HIGH" {
		t.Errorf("expected level HIGH, got %s", res.Level)
	}
	if res.Confidence != 0.92 {
		t.Errorf("expected confidence 0.92, got %f", res.Confidence)
	}
}
