package auth

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/testutil"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"golang.org/x/crypto/bcrypt"
)

type mockLogger struct{}

func (m *mockLogger) Record(
	ctx context.Context,
	tx datastore.DB,
	entry audit.LogEntry,
) {
}

type mockEmailer struct{}

func (m *mockEmailer) Send(
	ctx context.Context,
	email audit.EmailEntry,
) error {
	return nil
}

func setupTestRedis(t *testing.T) *datastore.RedisClient {
	cfg := &config.Config{
		RedisHost: "localhost",
		RedisPort: "6379",
	}
	client, err := datastore.NewRedisClient(cfg)
	if err != nil {
		t.Skip("skipping Redis-dependent tests: connection failed")
	}
	return client
}

func TestAuthLifecycle(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	redisClient := setupTestRedis(t)

	// Hash password
	pwd := "password123"
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to generate bcrypt hash: %v", err)
	}

	// Seed required tables first
	_, err = db.Exec(`
		INSERT INTO regions (id, name, code) VALUES (1, 'NCR', '130000000');
		INSERT INTO provinces (id, code, name, region_code) 
		VALUES (1, '133900000', 'NCR', '130000000');
		INSERT INTO cities (id, code, name, province_code, region_code) 
		VALUES (1, '133901000', 'Manila', '133900000', '130000000');
		INSERT INTO barangays (id, code, name, city_code) 
		VALUES (1, '133901001', 'Barangay 1', '133901000');
	`)
	if err != nil {
		t.Fatalf("failed to seed address tables: %v", err)
	}

	// Seed user with bcrypt hash
	_, err = db.Exec(`
		INSERT INTO users (
			id, email, first_name, last_name, auth_type,
			is_active, password_hash
		) VALUES (
			'student-uuid', 'student@test.com', 'Test', 'Student',
			'native', 1, ?
		);
	`, string(hash))
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	// Seed roles and iir
	_, err = db.Exec(`
		INSERT INTO user_roles (user_id, role_id) VALUES ('student-uuid', 1);
		INSERT INTO iir_records (id, user_id, is_submitted)
		VALUES ('iir-uuid', 'student-uuid', 1);
	`)
	if err != nil {
		t.Fatalf("failed to seed roles and iir: %v", err)
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

	// 1. Authenticate with valid credentials
	userID, token, refreshToken, err := svc.AuthenticateUser(
		ctx,
		"student@test.com",
		pwd,
		"127.0.0.1",
		"test-agent",
	)
	if err != nil {
		t.Fatalf("AuthenticateUser failed: %v", err)
	}
	if userID != "student-uuid" {
		t.Errorf("expected student-uuid, got %s", userID)
	}
	if token == "" || refreshToken == "" {
		t.Fatal("expected non-empty tokens")
	}

	// 2. Authenticate with invalid credentials
	_, _, _, err = svc.AuthenticateUser(
		ctx,
		"student@test.com",
		"wrong-password",
		"127.0.0.1",
		"test-agent",
	)
	if err == nil {
		t.Fatal("expected authentication to fail for invalid password")
	}

	// 3. Get profile (GetMe)
	me, err := svc.GetMe(ctx, userID, "native")
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}
	if me.Email != "student@test.com" {
		t.Errorf("expected student@test.com, got %s", me.Email)
	}

	// 4. Refresh token
	cfg := &config.Config{}
	newToken, newRefreshToken, err := svc.RefreshToken(
		ctx,
		refreshToken,
		cfg,
		"127.0.0.1",
		"test-agent",
	)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}
	if newToken == "" || newRefreshToken == "" {
		t.Fatal("expected non-empty refreshed tokens")
	}

	// 5. Logout
	_, err = svc.Logout(ctx, token, refreshToken, "native", cfg)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}
}
