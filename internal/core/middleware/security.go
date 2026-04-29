package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds essential security headers to every response.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Removed overly strict X-Frame-Options to allow frontend iframes
		// c.Header("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Basic XSS protection for older browsers
		c.Header("X-XSS-Protection", "1; mode=block")

		// Relaxed CSP to allow our frontend to embed files
		c.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors *;")

		// Prevent search engines from indexing the API
		c.Header("X-Robots-Tag", "noindex, nofollow")

		// Enforce HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Next()
	}
}

// BodySizeLimitMiddleware limits the maximum size of the request body.
func BodySizeLimitMiddleware(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to POST, PUT, PATCH
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch {

			// For multipart forms, Gin handles this differently via MaxMultipartMemory,
			// but for JSON payloads, we need MaxBytesReader.
			if !strings.HasPrefix(c.GetHeader("Content-Type"), "multipart/form-data") {
				c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
			}
		}
		c.Next()
	}
}
