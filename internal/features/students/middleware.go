package students

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// IIRLoggerMiddleware logs IIR responses (500, 200, 201, 404) to the database.
func IIRLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()
		if status != http.StatusInternalServerError &&
			status != http.StatusOK &&
			status != http.StatusCreated &&
			status != http.StatusNotFound {
			return
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		// Only inspect IIR process endpoints
		isIIR := strings.Contains(path, "/records/iir") ||
			strings.Contains(path, "/records/user")
		if !isIIR {
			return
		}

		logVal, exists := c.Get(middleware.SecurityLoggerContextKey)
		if !exists {
			return
		}
		logSvc, ok := logVal.(middleware.SecurityLogger)
		if !ok {
			return
		}

		var action string
		isDraft := strings.Contains(path, "/draft")
		isCOR := strings.Contains(path, "/cor") ||
			strings.Contains(path, "/cors")
		isPDF := strings.Contains(path, "/download")
		isCheckNum := strings.Contains(path, "/check-number")
		isProfile := strings.Contains(path, "/profile")

		var level string
		if status == http.StatusInternalServerError {
			level = audit.LevelError
		} else {
			level = audit.LevelInfo
		}

		// Map to appropriate Action constants
		switch method {
		case http.MethodPost:
			if isDraft {
				if status == http.StatusInternalServerError {
					action = "IIR_DRAFT_SAVE_FAILED"
				} else {
					action = audit.ActionIIRDraftSaved
				}
			} else if isCOR {
				if status == http.StatusInternalServerError {
					action = "IIR_COR_SUBMIT_FAILED"
				} else {
					action = "IIR_COR_SUBMITTED"
				}
			} else {
				if status == http.StatusInternalServerError {
					action = audit.ActionIIRCreateFailed
				} else {
					action = audit.ActionIIRSubmitted
				}
			}
		case http.MethodPatch:
			if status == http.StatusInternalServerError {
				action = audit.ActionIIRUpdateFailed
			} else {
				action = audit.ActionIIRUpdated
			}
		case http.MethodGet:
			if isDraft {
				if status == http.StatusInternalServerError {
					action = "IIR_DRAFT_FETCH_FAILED"
				} else if status == http.StatusNotFound {
					action = "IIR_DRAFT_NOT_FOUND"
				} else {
					action = "IIR_DRAFT_RETRIEVED"
				}
			} else if isCOR {
				if status == http.StatusInternalServerError {
					action = "IIR_COR_FETCH_FAILED"
				} else if status == http.StatusNotFound {
					action = "IIR_COR_NOT_FOUND"
				} else {
					action = "IIR_COR_RETRIEVED"
				}
			} else if isPDF {
				if status == http.StatusInternalServerError {
					action = "IIR_PDF_DOWNLOAD_FAILED"
				} else if status == http.StatusNotFound {
					action = "IIR_PDF_NOT_FOUND"
				} else {
					action = "IIR_PDF_DOWNLOADED"
				}
			} else if isCheckNum {
				if status == http.StatusInternalServerError {
					action = "IIR_STUDENT_NUMBER_CHECK_FAILED"
				} else {
					action = "IIR_STUDENT_NUMBER_CHECKED"
				}
			} else if isProfile {
				if status == http.StatusInternalServerError {
					action = "IIR_PROFILE_FETCH_FAILED"
				} else if status == http.StatusNotFound {
					action = "IIR_PROFILE_NOT_FOUND"
				} else {
					action = "IIR_PROFILE_RETRIEVED"
				}
			} else {
				if status == http.StatusInternalServerError {
					action = "IIR_FETCH_FAILED"
				} else if status == http.StatusNotFound {
					action = "IIR_NOT_FOUND"
				} else {
					action = "IIR_RETRIEVED"
				}
			}
		default:
			action = "IIR_ACTION"
		}

		var message string
		if status == http.StatusInternalServerError {
			if len(c.Errors) > 0 {
				message = c.Errors.Last().Error()
			} else {
				message = fmt.Sprintf(
					"IIR process error: Internal error on %s %s",
					method,
					path,
				)
			}
		} else {
			message = fmt.Sprintf(
				"IIR process: %s response status %d for path %s",
				action,
				status,
				path,
			)
		}

		userIDVal, _ := c.Get("userID")
		userIDStr, _ := userIDVal.(string)
		userEmailVal, _ := c.Get("userEmail")
		userEmailStr, _ := userEmailVal.(string)

		logSvc.RecordEntry(c.Request.Context(), audit.LogEntry{
			Level:     level,
			Category:  audit.CategoryAudit,
			Action:    action,
			Message:   message,
			UserID:    structs.StringToNullableString(userIDStr),
			UserEmail: structs.StringToNullableString(userEmailStr),
		})
	}
}
