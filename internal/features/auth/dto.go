package auth

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type TTL int

type MeResponse struct {
	ID             string       `json:"id"`
	Email          string       `json:"email"`
	FirstName      string       `json:"firstName"`
	LastName       string       `json:"lastName"`
	SuffixName     string       `json:"suffixName,omitempty"`
	MiddleName     string       `json:"middleName,omitempty"`
	ProfilePicture string       `json:"profilePicture,omitempty"`
	CreatedAt      time.Time    `json:"createdAt"`
	Roles          []users.Role `json:"roles"`
	StudentCORURL  string       `json:"studentCorUrl,omitempty"`
	IsStudentCORValid bool      `json:"isStudentCorValid"`
	Type           string       `json:"type"` // "native" or "idp"
}

type IDPRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type IDPTokenRequest struct {
	Code string `json:"code" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email      string `json:"email"      binding:"required,email"`
	Password   string `json:"password"   binding:"required,min=8"`
	FirstName  string `json:"firstName"  binding:"required"`
	LastName   string `json:"lastName"   binding:"required"`
	SuffixName string `json:"suffixName"`
	MiddleName string `json:"middleName"`
}

type VerifyRequest struct {
	VerificationOTP string `json:"otp" binding:"required"`
}
