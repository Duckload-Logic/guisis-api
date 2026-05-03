// owner.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

// OwnershipMiddleware - Direct database access version
func OwnershipMiddleware(db *sqlx.DB, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInUserID := c.MustGet("userID").(string)
		roleIDs := c.MustGet("roleIDs").([]int)

		isAdmin := false
		isStudent := false
		for _, rid := range roleIDs {
			if rid == int(constants.AdminRoleID) ||
				rid == int(constants.SuperAdminRoleID) {
				isAdmin = true
				break
			}
			if rid == int(constants.StudentRoleID) {
				isStudent = true
			}
		}

		// Allow counselors and super admins to bypass
		if isAdmin {
			c.Next()
			return
		}

		// For students, check ownership
		if isStudent {
			paramValue := c.Param(paramName)

			// For email-based params, compare directly
			if paramName == "userID" {
				if paramValue != loggedInUserID {
					c.AbortWithStatusJSON(
						http.StatusForbidden,
						gin.H{"error": "Access denied"},
					)
					return
				}
				c.Next()
				return
			}

			// For int-based params, parse and check ownership
			resourceID := paramValue

			owns, err := checkStudentOwnership(
				db, loggedInUserID, paramName, resourceID,
			)
			if err != nil || !owns {
				c.AbortWithStatusJSON(
					http.StatusForbidden,
					gin.H{"error": "Access denied"},
				)
				return
			}
		}

		c.Next()
	}
}

// Direct database query - ONE function to rule them all
func checkStudentOwnership(
	db *sqlx.DB, userID, paramName, resourceID string,
) (bool, error) {
	switch paramName {
	case "iirID":
		// Check if student_record belongs to user
		query := `
			SELECT EXISTS(
				SELECT 1 FROM iir_records
				WHERE id = ? AND user_id = ?
			)`
		var exists bool
		err := db.QueryRow(query, resourceID, userID).Scan(&exists)
		return exists, err
	case "appointmentID":
		query := `
			SELECT EXISTS(
				SELECT 1 FROM appointments a
				INNER JOIN iir_records ir ON a.iir_id = ir.id
				WHERE a.id = ? AND ir.user_id = ?
			)
		`
		var exists bool
		err := db.QueryRow(query, resourceID, userID).Scan(&exists)
		return exists, err
	case "slipID":
		query := `
			SELECT EXISTS(
				SELECT 1 FROM admission_slips s
				INNER JOIN iir_records ir ON s.iir_id = ir.id
				WHERE s.id = ? AND ir.user_id = ?
			)
		`
		var exists bool
		err := db.QueryRow(query, resourceID, userID).Scan(&exists)
		return exists, err
	default:
		return false, nil
	}
}
