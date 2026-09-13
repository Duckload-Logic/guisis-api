package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

func RoleMiddleware(allowedRoles ...constants.RoleID) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIDsVal, exists := c.Get("roleIDs")
		if !exists {
			log.Printf(
				"[RoleMiddleware] {Auth Context}: roleIDs not found in context",
			)
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Roles not found"},
			)
			return
		}

		var userRoles []int
		switch v := roleIDsVal.(type) {
		case []int:
			userRoles = v
		case []interface{}:
			for _, item := range v {
				if f, ok := item.(float64); ok {
					userRoles = append(userRoles, int(f))
				} else if i, ok := item.(int); ok {
					userRoles = append(userRoles, i)
				}
			}
		default:
			log.Printf(
				"[RoleMiddleware] {Type Assertion}: invalid roleIDs type: %T",
				roleIDsVal,
			)
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Invalid roles type"},
			)
			return
		}

		isAuthorized := false
		for _, urid := range userRoles {
			for _, allowed := range allowedRoles {
				if int(allowed) == urid {
					isAuthorized = true
					break
				}
			}
			if isAuthorized {
				break
			}
		}

		if !isAuthorized {
			if logSvc, ok := c.Get(SecurityLoggerContextKey); ok {
				if svc, ok := logSvc.(SecurityLogger); ok {
					userEmailVal, _ := c.Get("userEmail")
					userEmail, _ := userEmailVal.(string)
					svc.RecordSecurity(
						c.Request.Context(),
						"ACCESS_DENIED",
						fmt.Sprintf(
							"Access denied for %s (roles %v) on %s %s",
							userEmail,
							userRoles,
							c.Request.Method,
							c.Request.URL.Path,
						),
						userEmail,
						c.ClientIP(),
						c.Request.UserAgent(),
					)
				}
			}
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "Access denied"},
			)
			return
		}

		c.Next()
	}
}
