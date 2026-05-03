package m2mclients

import "time"

// Requests

type CreateM2MClientRequest struct {
	ClientName        string   `json:"clientName"          binding:"required"`
	ClientDescription string   `json:"clientDescription"   binding:"required"`
	Scopes            []string `json:"scopes,omitempty"`
	ExpiresAt         *string  `json:"expiresAt,omitempty"` // RFC3339 format, optional
}

type ListM2MClientsRequest struct {
	IncludeRevoked bool `form:"include_revoked,omitempty"`
}

type M2MTokenRequest struct {
	ClientID     string `json:"clientId"     binding:"required"`
	ClientSecret string `json:"clientSecret" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Responses

type M2MClientDTO struct {
	ID                int        `json:"id"`
	ClientName        string     `json:"clientName"`
	ClientID          string     `json:"clientId"`
	ClientDescription string     `json:"clientDescription"`
	Scopes            []string   `json:"scopes,omitempty"`
	IsActive          bool       `json:"isActive"`
	IsVerified        bool       `json:"isVerified"`
	LastUsedAt        *time.Time `json:"lastUsedAt,omitempty"`
	ExpiresAt         *time.Time `json:"expiresAt,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
}

type CreateM2MClientResponse struct {
	M2MClientDTO
	ClientSecret string `json:"clientSecret"`
}

type M2MTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"` // Seconds
}
