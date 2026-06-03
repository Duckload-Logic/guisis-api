package tokens

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID         string `json:"userId"`
	IDPUserID      string `json:"idpUserId"` // Only for IDP sessions
	UserEmail      string `json:"userEmail"`
	RoleIDs        []int  `json:"roleIds"`
	TokenType      string `json:"tokenType"`   // "native", "idp", or "m2m"
	M2MClientID    string `json:"m2mClientId"` // Only for M2M sessions
	IDPAccessToken string `json:"idpAccessToken,omitempty"`
	IsVerified     bool   `json:"isVerified"`
	ClientName     string `json:"clientName,omitempty"`
	IIRID          string `json:"iirId,omitempty"`
	CORID          string `json:"corId,omitempty"`
	HasPersonalInfoAccess bool `json:"hasPersonalInfoAccess,omitempty"`
	jwt.RegisteredClaims
}
