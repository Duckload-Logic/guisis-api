package middleware

import (
	"log"
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
	log.Printf("IS_VERIFIED: %v", isVerified)
	if vBool, isBool := isVerified.(bool); !ok || !isBool || !vBool {
		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"error": "Client not verified for this resource"},
		)
		return
	}

	c.Next()
}

func RequireM2MPersonalInfoAccess(c *gin.Context) {
	if isM2M, ok := c.Get("isM2M"); ok && isM2M.(bool) {
		access, ok := c.Get("hasPersonalInfoAccess")
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "DPA violation: access forbidden"},
			)
			return
		}
		hasAccess, ok := access.(bool)
		if !ok || !hasAccess {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"error": "DPA restriction: Access to " +
						"student personal info is denied",
				},
			)
			return
		}
	}
	c.Next()
}
