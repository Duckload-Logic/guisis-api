package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds essential security headers to every response.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Basic XSS protection for older browsers
		c.Header("X-XSS-Protection", "1; mode=block")

		// Content Security Policy
		// For docs, we need to allow styles, scripts and images
		if strings.HasPrefix(c.Request.URL.Path, "/api/v1/docs") {
			c.Header(
				"Content-Security-Policy",
				"default-src 'self'; style-src 'self' 'unsafe-inline'; "+
					"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
					"img-src 'self' data:; frame-ancestors 'none'; "+
					"form-action 'self';",
			)
		} else if strings.Contains(c.Request.URL.Path, "/uploads/") ||
			strings.Contains(c.Request.URL.Path, "/export") ||
			strings.Contains(c.Request.URL.Path, "/attachments/") {
			c.Header(
				"Content-Security-Policy",
				"default-src 'self' data: blob:; img-src 'self' data: "+
					"blob:; style-src 'self' 'unsafe-inline'; "+
					"frame-ancestors 'self' http://localhost:* "+
					"https://*.dllbsit2027.com;",
			)
		} else {
			c.Header(
				"Content-Security-Policy",
				"default-src 'none'; frame-ancestors 'none'; "+
					"form-action 'none';",
			)
		}

		// Prevent search engines from indexing the API
		c.Header("X-Robots-Tag", "noindex, nofollow")

		// Anti-clickjacking Header
		if strings.Contains(c.Request.URL.Path, "/uploads/") ||
			strings.Contains(c.Request.URL.Path, "/export") ||
			strings.Contains(c.Request.URL.Path, "/attachments/") {
			c.Header("X-Frame-Options", "SAMEORIGIN")
		} else {
			c.Header("X-Frame-Options", "DENY")
		}

		// Set a generic server name to mask actual server details
		c.Header("Server", "GuiSIS-API")

		// Restrict Flash/Silverlight cross-domain policies
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// Re-examine Cache-control Directives
		if !strings.HasPrefix(c.Request.URL.Path, "/api/v1/docs") {
			c.Header(
				"Cache-Control",
				"no-store, no-cache, must-revalidate, proxy-revalidate",
			)
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		// Enforce HTTPS
		c.Header(
			"Strict-Transport-Security",
			"max-age=31536000; includeSubDomains",
		)

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
			if !strings.HasPrefix(
				c.GetHeader("Content-Type"),
				"multipart/form-data",
			) {
				c.Request.Body = http.MaxBytesReader(
					c.Writer,
					c.Request.Body,
					limit,
				)
			}
		}
		c.Next()
	}
}
