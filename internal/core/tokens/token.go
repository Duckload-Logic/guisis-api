package tokens

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

type Service struct {
	secret []byte
}

func NewService() *Service {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Printf(
			"[WARNING] {NewService}: JWT_SECRET is empty. Signing will fail.",
		)
	}
	return &Service{secret: []byte(secret)}
}

func (s *Service) GenerateToken(
	userEmail string,
	userID string,
	roleIDs []int,
	tokenType string,
	expireSeconds int,
) (string, *Claims, error) {
	claims := &Claims{
		UserEmail: userEmail,
		UserID:    userID,
		RoleIDs:   roleIDs,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(
				time.Duration(expireSeconds) * time.Second),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer:   constants.ClaimsIssuer,
			ID:       uuid.New().String(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	return signed, claims, err
}

func (s *Service) ValidateToken(tokenString string) (
	*Claims, error,
) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return s.secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

// ParseTokenUnverified extracts claims from a token string without
// verifying its signature or expiration. Use this ONLY to identify
// a session for refresh logic, never for authorization.
func (s *Service) ParseTokenUnverified(tokenString string) (
	*Claims, error,
) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
