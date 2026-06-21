package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func init() {
	gin.SetMode(gin.TestMode)
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

func TestAuthAndRoleMiddlewares(t *testing.T) {
	// Setup Token Service
	tokenSvc := tokens.NewService()

	// Generate valid test tokens
	studentToken, studentClaims, err := tokenSvc.GenerateSessionToken(
		"student@test.com",
		"student-id",
		[]int{int(constants.StudentRoleID)},
		string(constants.AuthTypeNative),
		3600,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to generate student token: %v", err)
	}

	adminToken, _, err := tokenSvc.GenerateSessionToken(
		"admin@test.com",
		"admin-id",
		[]int{int(constants.AdminRoleID)},
		string(constants.AuthTypeNative),
		3600,
		nil,
	)
	if err != nil {
		t.Fatalf("failed to generate admin token: %v", err)
	}

	// Setup Gin router
	r := gin.New()
	r.Use(AuthMiddleware(nil))

	r.GET(
		"/student",
		RoleMiddleware(constants.StudentRoleID),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		},
	)

	r.GET(
		"/admin",
		RoleMiddleware(constants.AdminRoleID),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		},
	)

	// Test Case 1: No Auth Token
	req1 := httptest.NewRequest("GET", "/student", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w1.Code)
	}

	// Test Case 2: Invalid Auth Token
	req2 := httptest.NewRequest("GET", "/student", nil)
	req2.Header.Set("Authorization", "Bearer invalid-token")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w2.Code)
	}

	// Test Case 3: Valid Student accessing Student route
	req3 := httptest.NewRequest("GET", "/student", nil)
	req3.Header.Set("Authorization", "Bearer "+studentToken)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	if w3.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w3.Code)
	}

	// Test Case 4: Student accessing Admin route
	req4 := httptest.NewRequest("GET", "/admin", nil)
	req4.Header.Set("Authorization", "Bearer "+studentToken)
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)
	if w4.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w4.Code)
	}

	// Test Case 5: Admin accessing Admin route
	req5 := httptest.NewRequest("GET", "/admin", nil)
	req5.Header.Set("Authorization", "Bearer "+adminToken)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)
	if w5.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w5.Code)
	}
	// Test Case 6: Session Whitelisting and Revocation (using Redis)
	t.Run("Redis Whitelist and Revocation", func(t *testing.T) {
		redisClient := setupTestRedis(t)

		// Setup route that uses Redis session check
		rRev := gin.New()
		rRev.Use(AuthMiddleware(redisClient))
		rRev.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		ctx := context.Background()
		sessionSvc := sessions.NewService(redisClient)

		// Whitelist the session
		err := sessionSvc.WhitelistSession(
			ctx,
			"student-id",
			studentClaims.ID,
			"refresh-jti-placeholder",
			60,
		)
		if err != nil {
			t.Fatalf("failed to whitelist session: %v", err)
		}
		defer sessionSvc.RevokeUserSession(ctx, "student-id")

		// 1. Verify access succeeds when whitelisted
		reqValid := httptest.NewRequest("GET", "/test", nil)
		reqValid.Header.Set("Authorization", "Bearer "+studentToken)
		wValid := httptest.NewRecorder()
		rRev.ServeHTTP(wValid, reqValid)
		if wValid.Code != http.StatusOK {
			t.Errorf(
				"expected 200 for whitelisted session, got %d",
				wValid.Code,
			)
		}

		// 2. Revoke the whitelist key
		err = sessionSvc.RevokeUserSession(ctx, "student-id")
		if err != nil {
			t.Fatalf("failed to revoke session: %v", err)
		}

		// 3. Verify access is rejected when revoked
		reqRevoked := httptest.NewRequest("GET", "/test", nil)
		reqRevoked.Header.Set("Authorization", "Bearer "+studentToken)
		wRevoked := httptest.NewRecorder()
		rRev.ServeHTTP(wRevoked, reqRevoked)
		if wRevoked.Code != http.StatusUnauthorized {
			t.Errorf(
				"expected 401 for revoked session, got %d",
				wRevoked.Code,
			)
		}

		var resp map[string]string
		_ = json.Unmarshal(wRevoked.Body.Bytes(), &resp)
		if resp["error"] != "Session has been revoked or logged out" {
			t.Errorf("unexpected error message: %s", resp["error"])
		}
	})
}
