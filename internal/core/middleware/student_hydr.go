package middleware

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

// HydrateStudentContext extracts student IIR ID from database and
// sets it in the Gin context. Only applies to Student role users.
// Day One students (no IIR record) can still proceed.
func HydrateStudentIIRContext(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, exists := c.Get("iirID"); exists {
			if _, ok := val.(string); ok && val != "" {
				c.Next()
				return
			}
		}

		userID, isStudent := extractStudentContext(c)
		if !isStudent {
			c.Next()
			return
		}

		// Query iir_records table to find IIR ID by user_id
		var iirID string
		err := db.QueryRow(`
			SELECT id FROM iir_records WHERE user_id = ?
		`, userID).Scan(&iirID)

		if err == sql.ErrNoRows {
			// Day One student - no IIR record yet
			c.Next()
			return
		}

		if err != nil {
			log.Printf(
				"[HydrateStudentContext] {Database Query IIR "+
					"Lookup}: %v",
				err)
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "Internal server error"},
			)
			return
		}

		// Set IIR ID in context for downstream handlers
		c.Set("iirID", iirID)
		c.Next()
	}
}

func HydrateStudentCORContext(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, exists := c.Get("corID"); exists {
			if _, ok := val.(string); ok && val != "" {
				c.Next()
				return
			}
		}

		userID, isStudent := extractStudentContext(c)
		if !isStudent {
			c.Next()
			return
		}

		// Query student_cors table to find COR ID by user_id
		var corID string
		err := db.QueryRow(`
			SELECT file_id FROM student_cors WHERE student_id = ?
		`, userID).Scan(&corID)

		if err == sql.ErrNoRows {
			// Day One student - no COR record yet
			c.Next()
			return
		}

		if err != nil {
			log.Printf(
				"[HydrateStudentContext] {Database Query COR "+
					"Lookup}: %v",
				err)
			c.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{"error": "Internal server error"},
			)
			return
		}

		// Set COR ID in context for downstream handlers
		c.Set("corID", corID)
		c.Next()
	}
}

// RequireCOR ensures that a student has uploaded a COR before proceeding.
// Non-student roles bypass this check automatically.
func RequireCOR() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, isStudent := extractStudentContext(c)
		if !isStudent {
			// Admins and developers bypass this requirement
			c.Next()
			return
		}

		_, exists := c.Get("corID")
		if !exists {
			// Instead of aborting with generic 403, we provide a
			// structured fail response
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{
					"status": "fail",
					"data": gin.H{
						"error": "A valid Certificate of Registration (COR) " +
							"is required to proceed. Please upload your " +
							"COR in your profile.",
					},
				},
			)
			return
		}

		c.Next()
	}
}

func extractStudentContext(c *gin.Context) (string, bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return "", false
	}

	roleIDsVal, exists := c.Get("roleIDs")
	if !exists {
		return "", false
	}

	userID, ok := userIDVal.(string)
	if !ok {
		return "", false
	}

	var roleIDs []int
	switch v := roleIDsVal.(type) {
	case []int:
		roleIDs = v
	case []interface{}:
		for _, item := range v {
			if f, ok := item.(float64); ok {
				roleIDs = append(roleIDs, int(f))
			} else if i, ok := item.(int); ok {
				roleIDs = append(roleIDs, i)
			}
		}
	default:
		return "", false
	}

	for _, rid := range roleIDs {
		if rid == int(constants.StudentRoleID) {
			return userID, true
		}
	}

	return userID, false
}
