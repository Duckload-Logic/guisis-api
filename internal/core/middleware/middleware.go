package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireVerifiedM2M(c *gin.Context) {
	isM2M, ok := c.Get("isM2M")
	if m2mBool, isBool := isM2M.(bool); !ok || !isBool || !m2mBool {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"error": "M2M access only"},
		)
		return
	}

	isVerified, ok := c.Get("isVerified")
	if vBool, isBool := isVerified.(bool); !ok || !isBool || !vBool {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"error": "Client not verified for this resource"},
		)
		return
	}

	c.Next()
}
