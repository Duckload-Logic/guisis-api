package tokens

import (
	"os"
	"testing"
)

func TestService_TokenWorkflow(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	svc := &Service{secret: []byte("test-secret")}

	userEmail := "test@example.com"
	userID := "user-123"
	roleID := 1
	tokenType := "access"
	expireSeconds := 60

	t.Run("generate and validate token", func(t *testing.T) {
		tokenString, claims, err := svc.GenerateToken(
			userEmail,
			userID,
			[]int{roleID},
			tokenType,
			expireSeconds,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if tokenString == "" || claims == nil {
			t.Fatal("generated token or claims are empty")
		}

		if claims.UserEmail != userEmail || claims.UserID != userID {
			t.Errorf("claims mismatch: %+v", claims)
		}

		// Validate
		validatedClaims, err := svc.ValidateToken(tokenString)
		if err != nil {
			t.Fatalf("validation failed: %v", err)
		}
		if validatedClaims.UserID != userID {
			t.Errorf(
				"validated userID %s, want %s",
				validatedClaims.UserID,
				userID,
			)
		}
	})

	t.Run("validate invalid token", func(t *testing.T) {
		_, err := svc.ValidateToken("invalid-token-string")
		if err == nil {
			t.Error("expected error for invalid token")
		}
	})

	t.Run("unverified parse", func(t *testing.T) {
		tokenString, _, _ := svc.GenerateToken(
			userEmail,
			userID,
			[]int{roleID},
			tokenType,
			expireSeconds,
		)

		// Parse without verification
		claims, err := svc.ParseTokenUnverified(tokenString)
		if err != nil {
			t.Fatalf("parse unverified failed: %v", err)
		}
		if claims.UserID != userID {
			t.Errorf("got %s, want %s", claims.UserID, userID)
		}
	})

	t.Run("validate expired token", func(t *testing.T) {
		// Generate with negative expiry
		tokenString, _, _ := svc.GenerateToken(
			userEmail,
			userID,
			[]int{roleID},
			tokenType,
			-60,
		)

		_, err := svc.ValidateToken(tokenString)
		if err == nil {
			t.Error("expected error for expired token")
		}
	})
}

func TestNewService(t *testing.T) {
	os.Setenv("JWT_SECRET", "another-secret")
	svc := NewService()
	if string(svc.secret) != "another-secret" {
		t.Errorf("got secret %s, want another-secret", string(svc.secret))
	}
}
