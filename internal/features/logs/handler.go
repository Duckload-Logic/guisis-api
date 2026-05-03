package logs

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetService() *Service {
	return h.service
}

// HandleListSystemLogs godoc
// @Summary      List system logs
// @Description  Retrieves a paginated list of audit, system, and security logs.
// Super Admin only.
// @Tags         SystemLogs
// @Accept       json
// @Produce      json
// @Param        page        query     int    false "Page number"
// @Param        page_size   query     int    false "Number of entries per page"
// @Param        category    query     string false "Filter by category (AUDIT, SYSTEM, SECURITY)"
// @Param        action      query     string false "Filter by action"
// @Param        user_email  query     string false "Filter by user email"
// @Param        start_date  query     string false "Filter from date (YYYY-MM-DD)"
// @Param        end_date    query     string false "Filter to date (YYYY-MM-DD)"
// @Param        search      query     string false "Search in message, action, or user email"
// @Success      200         {object}  ListSystemLogsDTO
// @Failure      400         {object}  map[string]string "Bad request"
// @Failure      500         {object}  map[string]string "Internal server error"
// @Router       /system-logs [get]
func (h *Handler) GetLogs(c *gin.Context) {
	var req audit.ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf("[GetLogs] {Bind Query}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetLogs] {Fetch Logs}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve system logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

func (h *Handler) GetLogsAudit(c *gin.Context) {
	var req audit.ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf("[GetLogsAudit] {Bind Query}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	req.Category = audit.CategoryAudit

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetLogsAudit] {Fetch Logs}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve audit logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

func (h *Handler) GetLogsSystem(c *gin.Context) {
	var req audit.ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf("[GetLogsSystem] {Bind Query}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	req.Category = audit.CategorySystem

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetLogsSystem] {Fetch Logs}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve system logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

// GetSecurityLogs returns only SECURITY category logs
func (h *Handler) GetLogsSecurity(c *gin.Context) {
	var req audit.ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf("[GetLogsSecurity] {Bind Query}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	req.Category = audit.CategorySecurity

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetLogsSecurity] {Fetch Logs}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve security logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

// GetLogStats returns log counts by category
func (h *Handler) GetLogsStats(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	stats, err := h.service.GetStats(
		c.Request.Context(),
		startDate,
		endDate,
	)
	if err != nil {
		fmt.Printf("[GetLogsStats] {GetStats}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve log stats",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, stats)
}

// GetActivityStats returns log counts grouped by hour for the last 24 hours
func (h *Handler) GetLogsActivity(c *gin.Context) {
	stats, err := h.service.GetActivityStats(c.Request.Context())
	if err != nil {
		fmt.Printf("[GetLogsActivity] {GetActivityStats}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve log activity stats",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, stats)
}

// GetMyLogs retrieves activity logs for the currently authenticated user.
func (h *Handler) GetLogsMe(c *gin.Context) {
	userEmail := c.MustGet("userEmail").(string)

	var req audit.ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf("[GetLogsMe] {Bind Query}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	req.UserEmail = userEmail

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetLogsMe] {Fetch Logs}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve your activity logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

// PostLogsCleanup godoc
// @Summary      Clean up old system logs
// @Description  Deletes logs older than N days. Super Admin only.
// @Tags         SystemLogs
// @Param        days query     int    false "Days of logs to keep (default 30)"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]string
// @Router       /activity-meta/cleanup [post]
func (h *Handler) PostLogsCleanup(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	var days int
	if _, err := fmt.Sscanf(daysStr, "%d", &days); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid days format"})
		return
	}

	rows, err := h.service.DeleteLogsOlderThan(c.Request.Context(), days)
	if err != nil {
		fmt.Printf("[PostLogsCleanup] {Delete Logs}: %v\n", err)
		response.SendError(
			c,
			"Failed to clean up logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"message":       fmt.Sprintf("Logs older than %d days cleaned up", days),
		"rows_affected": rows,
	})
}
