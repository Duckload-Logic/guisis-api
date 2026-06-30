package constants

import "time"

const (
	ClaimsIssuer = "guisis-api"
)

type AuthType string

const (
	AuthTypeNative AuthType = "native"
	AuthTypeIDP    AuthType = "idp"
)

// Timeout constants for IDP requests
const (
	// IDPRequestTimeout is the timeout duration for IDP HTTP requests
	// (10 seconds)
	IDPRequestTimeout = 30 * time.Second
)

// Cookie configuration constants
const (
	// AccessTokenCookieName is the name of the access token cookie
	AccessTokenCookieName = "access_token"

	// RefreshTokenCookieName is the name of the refresh token cookie
	RefreshTokenCookieName = "refresh_token"

	// AccessTokenMaxAge is the maximum age in seconds for access token
	// cookie (30 minutes = 1800 seconds)
	AccessTokenMaxAge    = 1800
	M2MAccessTokenMaxAge = 3600

	// RefreshTokenMaxAge is the maximum age in seconds for refresh token
	// cookie (12 hours = 43200 seconds)
	RefreshTokenMaxAge    = 43200
	M2MRefreshTokenMaxAge = 86400

	// CookiePathRoot sets cookies to be accessible from root path
	CookiePathRoot = "/"
)


// Redis Key Constants
const (
	// RedisSessionKeyPrefix is the prefix for session keys (session:jti)
	RedisSessionKeyPrefix = "session:"

	// RedisIDPRefreshKeyPrefix is the prefix for IDP refresh tokens
	// (idp_refresh:jti)
	RedisIDPRefreshKeyPrefix = "idp_refresh:"

	// RedisUserSessionKeyPrefix is the prefix for whitelisted user sessions
	RedisUserSessionKeyPrefix = "user:session:"

	// RedisSessionAccessJTIField is the hash field for access token JTI
	RedisSessionAccessJTIField = "access_jti"

	// RedisSessionRefreshJTIField is the hash field for refresh token JTI
	RedisSessionRefreshJTIField = "refresh_jti"
)
