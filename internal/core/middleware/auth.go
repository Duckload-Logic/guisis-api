package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
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
	if claims.M2MClientID != "" {
		return true
	}

	key := fmt.Sprintf(
		"%s%s",
		constants.RedisUserSessionKeyPrefix,
		claims.UserID,
	)

	fields, err := redis.HGetAll(c.Request.Context(), key)
	if err != nil {
		log.Printf(
			"[AuthMiddleware] {Session Validate}: Redis error: %v",
			err,
		)
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "Session has been revoked or logged out"},
		)
		return false
	}

	if len(fields) == 0 {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "Session has been revoked or logged out"},
		)
		return false
	}

	whitelistedJTI := fields[constants.RedisSessionAccessJTIField]
	if whitelistedJTI != claims.ID {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{"error": "Session has been revoked or logged out"},
		)
		return false
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
		c.Set("hasPersonalInfoAccess", claims.HasPersonalInfoAccess)
		c.Set("isVerified", claims.IsVerified)
		c.Set("clientName", claims.ClientName)
	} else {
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.UserEmail)
		c.Set("roleIDs", claims.RoleIDs)
		if claims.IIRID != "" {
			c.Set("iirID", claims.IIRID)
		}
		if claims.CORID != "" {
			c.Set("corID", claims.CORID)
		}
	}
	c.Set("tokenType", claims.TokenType)
	if claims.IDPUserID != "" {
		c.Set("idpUserID", claims.IDPUserID)
	}
	if claims.IDPAccessToken != "" {
		c.Set("idpAccessToken", claims.IDPAccessToken)
	}
}
