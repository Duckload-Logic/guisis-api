package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func AuthMiddleware(redis *datastore.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getTokenString(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "No authentication token provided"},
			)
			return
		}

		claims, err := tokens.NewService().ValidateToken(tokenString)
		if err != nil {
			log.Printf("[AuthMiddleware] {Token}: Invalid or expired: %v", err)
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "Invalid or expired token"},
			)
			return
		}

		if redis != nil {
			if !validateSession(c, redis, claims) {
				return
			}
		}

		if claims.M2MClientID != "" && !validateM2MPath(c, claims.M2MClientID) {
			return
		}

		setContextInfo(c, claims)
		c.Next()
	}
}

func getTokenString(c *gin.Context) string {
	if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
		return cookie
	}
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return ""
}

func validateSession(
	c *gin.Context,
	redis *datastore.RedisClient,
	claims *tokens.Claims,
) bool {
	jti := sessions.NewJTI(claims.ID)
	val, err := redis.Get(c.Request.Context(), jti.ToSessionKey())
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "Session has been revoked or expired"},
		)
		return false
	}

	var sessionData map[string]string
	if err := json.Unmarshal([]byte(val), &sessionData); err == nil {
		if idpToken, ok := sessionData["idpAccessToken"]; ok {
			c.Set("idpAccessToken", idpToken)
		}
		if isVerified, ok := sessionData["isVerified"]; ok {
			c.Set("isVerified", isVerified == "true")
		}
		if clientName, ok := sessionData["clientName"]; ok {
			c.Set("clientName", clientName)
		}
	}
	return true
}

func validateM2MPath(c *gin.Context, clientID string) bool {
	fullPath := c.Request.URL.Path

	// Prevent simple path traversal attempts
	if strings.Contains(fullPath, "..") {
		log.Printf(
			"[Security] {M2M Path Traversal Attempt}: client %s tried %s",
			clientID,
			fullPath,
		)
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"error": "Unauthorized path access"},
		)
		return false
	}

	allowedPrefix := "/api/v1/integrations"
	if !strings.HasPrefix(fullPath, allowedPrefix) {
		log.Printf(
			"[Security] {M2M Out-of-Scope Access}: client %s tried %s",
			clientID,
			fullPath,
		)
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"error": "M2M clients are restricted to integration routes"},
		)
		return false
	}
	return true
}

func setContextInfo(c *gin.Context, claims *tokens.Claims) {
	if claims.M2MClientID != "" {
		c.Set("m2mClientID", claims.M2MClientID)
		c.Set("isM2M", true)
	} else {
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.UserEmail)
		c.Set("roleIDs", claims.RoleIDs)
	}
	c.Set("tokenType", claims.TokenType)
	if claims.IDPUserID != "" {
		c.Set("idpUserID", claims.IDPUserID)
	}
}
