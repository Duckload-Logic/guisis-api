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
	ctx := context.Background()
	cfg := &config.Config{}

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}
				if r.URL.Path != "/classify" {
					t.Errorf("expected /classify, got %s", r.URL.Path)
				}

				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(ClassifyResponse{
					Level:      "HIGH",
					Confidence: 0.95,
				})
			}),
		)
		defer server.Close()

		client := NewClient(http.DefaultClient, server.URL)
		resp, err := client.Classify(ctx, "test text", cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Level != "HIGH" {
			t.Errorf("got %s, want HIGH", resp.Level)
		}
	})

	t.Run("error_status", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
		)
		defer server.Close()

		client := NewClient(http.DefaultClient, server.URL)
		resp, err := client.Classify(ctx, "test text", cfg)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if resp != nil {
			t.Error("expected nil response, got one")
		}
	})
}
