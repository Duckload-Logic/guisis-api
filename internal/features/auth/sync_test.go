package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/testutil"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/identity/idp"
)

func TestPostIDPTokenExchange_SyncNames(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)

	// Mock IDP server responses
	var mockIDPUser idp.IDPUserInfo
	idpServer := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.URL.Path == "/auth/token" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(idp.IDPTokenResponse{
				AccessToken: "mock-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			})
			return
		}
		if r.URL.Path == "/me" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(mockIDPUser)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer idpServer.Close()

	cfg := &config.Config{
		IDPBaseUrl:      idpServer.URL,
		IDPClientID:     "client-id",
		IDPClientSecret: "client-secret",
	}

	userRepo := users.NewRepository(db)
	sessionSvc := sessions.NewService(redisClient)
	svc := NewService(
		userRepo,
		redisClient,
		sessionSvc,
		&mockEmailer{},
		&mockLogger{},
	)

	ctx := context.Background()

	// 1. Scenario: New JIT user
	mockIDPUser = idp.IDPUserInfo{
		ID:         "19486623-db22-43ff-b63d-043224b7253e",
		Email:      "dionkylo123@gmail.com",
		FirstName:  "JENNIFER",
		LastName:   "MORALES",
		MiddleName: "A",
		SuffixName: "Jr.",
	}

	_, _, err := svc.PostIDPTokenExchange(
		ctx,
		"code123",
		cfg,
		"127.0.0.1",
		"test-agent",
	)
	if err != nil {
		t.Fatalf("JIT Provision failed: %v", err)
	}

	// Verify database record
	localUser, err := userRepo.GetUserByIDPUUID(ctx, mockIDPUser.ID)
	if err != nil {
		t.Fatalf("Failed to fetch provisioned user: %v", err)
	}

	if localUser.FirstName != "JENNIFER" ||
		localUser.MiddleName.String != "A" ||
		localUser.LastName != "MORALES" ||
		localUser.SuffixName.String != "Jr." {
		t.Errorf("Unexpected user fields: %+v", localUser)
	}

	// 2. Scenario: Update name on IDP and login again
	mockIDPUser.FirstName = "JENNIFER REAL"
	mockIDPUser.SuffixName = "III"

	_, _, err = svc.PostIDPTokenExchange(
		ctx,
		"code124",
		cfg,
		"127.0.0.1",
		"test-agent",
	)
	if err != nil {
		t.Fatalf("Sync update failed: %v", err)
	}

	localUser, err = userRepo.GetUserByIDPUUID(ctx, mockIDPUser.ID)
	if err != nil {
		t.Fatalf("Failed to fetch user after sync: %v", err)
	}

	if localUser.FirstName != "JENNIFER REAL" ||
		localUser.SuffixName.String != "III" {
		t.Errorf("Sync did not update names in DB: %+v", localUser)
	}
}

func TestPostIDPTokenExchange_NativeEmailLink(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)

	var mockIDPUser idp.IDPUserInfo
	idpServer := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.URL.Path == "/auth/token" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(idp.IDPTokenResponse{
				AccessToken: "mock-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			})
			return
		}
		if r.URL.Path == "/me" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(mockIDPUser)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer idpServer.Close()

	cfg := &config.Config{
		IDPBaseUrl:      idpServer.URL,
		IDPClientID:     "client-id",
		IDPClientSecret: "client-secret",
	}

	userRepo := users.NewRepository(db)
	sessionSvc := sessions.NewService(redisClient)
	svc := NewService(
		userRepo,
		redisClient,
		sessionSvc,
		&mockEmailer{},
		&mockLogger{},
	)

	ctx := context.Background()

	// Seed pre-existing native user
	_, err := db.Exec(`
		INSERT INTO users (
			id, email, first_name, last_name, auth_type, is_active
		) VALUES (
			'native-user-id', 'native@test.com', 'OldFirst', 'OldLast',
			'native', 1
		)
	`)
	if err != nil {
		t.Fatalf("Failed to seed native user: %v", err)
	}

	// IDP login with same email
	mockIDPUser = idp.IDPUserInfo{
		ID:         "idp-uuid-999",
		Email:      "native@test.com",
		FirstName:  "NewFirst",
		LastName:   "NewLast",
		MiddleName: "NewMiddle",
		SuffixName: "NewSuffix",
	}

	_, _, err = svc.PostIDPTokenExchange(
		ctx,
		"code999",
		cfg,
		"127.0.0.1",
		"test-agent",
	)
	if err != nil {
		t.Fatalf("Failed IDP token exchange for native link: %v", err)
	}

	// Verify that the native user record was linked/updated
	linkedUser, err := userRepo.GetUserByIDPUUID(ctx, "idp-uuid-999")
	if err != nil {
		t.Fatalf("Failed to fetch linked user: %v", err)
	}

	if linkedUser.ID != "native-user-id" {
		t.Errorf(
			"Expected linked user ID to be native-user-id, got %s",
			linkedUser.ID,
		)
	}

	if linkedUser.FirstName != "NewFirst" ||
		linkedUser.LastName != "NewLast" ||
		linkedUser.MiddleName.String != "NewMiddle" ||
		linkedUser.SuffixName.String != "NewSuffix" {
		t.Errorf("Linked user fields were not updated: %+v", linkedUser)
	}
}

func TestPostIDPTokenExchange_IDPEmailLink(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)

	var mockIDPUser idp.IDPUserInfo
	idpServer := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.URL.Path == "/auth/token" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(idp.IDPTokenResponse{
				AccessToken: "mock-token",
				TokenType:   "Bearer",
				ExpiresIn:   3600,
			})
			return
		}
		if r.URL.Path == "/me" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(mockIDPUser)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer idpServer.Close()

	cfg := &config.Config{
		IDPBaseUrl:      idpServer.URL,
		IDPClientID:     "client-id",
		IDPClientSecret: "client-secret",
	}

	userRepo := users.NewRepository(db)
	sessionSvc := sessions.NewService(redisClient)
	svc := NewService(
		userRepo,
		redisClient,
		sessionSvc,
		&mockEmailer{},
		&mockLogger{},
	)

	ctx := context.Background()

	// Seed pre-existing IDP user with NULL idp_uuid
	_, err := db.Exec(`
		INSERT INTO users (
			id, email, first_name, last_name, auth_type, is_active
		) VALUES (
			'idp-user-id', 'idp_linked@test.com', 'OldFirst', 'OldLast',
			'idp', 1
		)
	`)
	if err != nil {
		t.Fatalf("Failed to seed IDP user: %v", err)
	}

	// IDP login with same email and new idp-uuid
	mockIDPUser = idp.IDPUserInfo{
		ID:         "idp-uuid-888",
		Email:      "idp_linked@test.com",
		FirstName:  "NewFirst",
		LastName:   "NewLast",
		MiddleName: "NewMiddle",
		SuffixName: "NewSuffix",
	}

	_, _, err = svc.PostIDPTokenExchange(
		ctx,
		"code888",
		cfg,
		"127.0.0.1",
		"test-agent",
	)
	if err != nil {
		t.Fatalf("Failed IDP token exchange for IDP link: %v", err)
	}

	// Verify that the IDP user record was linked/updated correctly
	linkedUser, err := userRepo.GetUserByIDPUUID(ctx, "idp-uuid-888")
	if err != nil {
		t.Fatalf("Failed to fetch linked user: %v", err)
	}

	if linkedUser.ID != "idp-user-id" {
		t.Errorf(
			"Expected linked user ID to be idp-user-id, got %s",
			linkedUser.ID,
		)
	}

	if linkedUser.FirstName != "NewFirst" ||
		linkedUser.LastName != "NewLast" ||
		linkedUser.MiddleName.String != "NewMiddle" ||
		linkedUser.SuffixName.String != "NewSuffix" {
		t.Errorf("Linked user fields were not updated: %+v", linkedUser)
	}
}
