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

	roleIDs, ok := roleIDsVal.([]int)
	if !ok {
		return "", false
	}

	for _, rid := range roleIDs {
		if rid == int(constants.StudentRoleID) {
			return userID, true
		}
	}

	return userID, false
}
