package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/identity/idp"
)

func TestPingIDP_Up(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	client := idp.NewIDPClient()
	cfg := &config.Config{
		IDPBaseUrl: server.URL + "/api/v1",
	}

	err := client.PingIDP(context.Background(), cfg)
	if err != nil {
		t.Fatalf("expected IDP to be up, got error: %v", err)
	}
}

func TestPingIDP_Down(t *testing.T) {
	// Use an invalid URL that will trigger a connection error
	client := idp.NewIDPClient()
	cfg := &config.Config{
		IDPBaseUrl: "http://localhost:9999/api/v1",
	}

	err := client.PingIDP(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected IDP ping to fail, but it succeeded")
	}
}

func TestIsIDPUp(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	defer server.Close()

	client := idp.NewIDPClient()
	svc := &Service{
		idpClient: *client,
		redis:     nil, // safe due to check
	}

	cfg := &config.Config{
		IDPBaseUrl: server.URL + "/api/v1",
	}

	up := svc.IsIDPUp(context.Background(), cfg)
	if !up {
		t.Fatal("expected IsIDPUp to return true")
	}
}
