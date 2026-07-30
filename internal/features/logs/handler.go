package logs

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
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
// @Param        sort_by     query     string false "Column to sort by (timestamp, message, actor, ipAddress)"
// @Param        sort_order  query     string false "Direction to sort (asc or desc)"
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
	if h.exportLogsCSV(c, req) {
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
	if h.exportLogsCSV(c, req) {
		return
	}

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
	if h.exportLogsCSV(c, req) {
		return
	}

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
	if h.exportLogsCSV(c, req) {
		return
	}

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

// GetLogsM2M returns only M2M category logs
func (h *Handler) GetLogsM2M(c *gin.Context) {
	var req audit.ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf("[GetLogsM2M] {Bind Query}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	req.Category = audit.CategoryM2M
	if h.exportLogsCSV(c, req) {
		return
	}

	result, err := h.service.ListLogs(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetLogsM2M] {Fetch Logs}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve M2M logs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

func (h *Handler) exportLogsCSV(
	c *gin.Context,
	req audit.ListSystemLogsRequest,
) bool {
	if c.Query("export") != "csv" {
		return false
	}

	csvData, err := h.service.ExportLogsCSV(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetLogs] {Export CSV}: %v\n", err)
		response.SendError(c, "Failed to generate CSV report", http.StatusInternalServerError, nil)
		return true
	}

	c.Header("Content-Disposition", "attachment; filename=system_logs_report.csv")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", csvData)
	return true
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

// GetLog retrieves a single system log by its ID.
func (h *Handler) GetLog(c *gin.Context) {
	idStr := c.Param("id")
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		fmt.Printf("[GetLog] {Parse ID}: invalid log id format: %s\n", idStr)
		response.SendFail(c, gin.H{"error": "Invalid log ID format"})
		return
	}

	result, err := h.service.GetLogByID(c.Request.Context(), id)
	if err != nil {
		fmt.Printf("[GetLog] {Fetch Log}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve log",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

// GetTraceTracks retrieves system logs by their trace ID.
func (h *Handler) GetTraceTracks(c *gin.Context) {
	traceID := c.Param("traceId")
	if traceID == "" {
		fmt.Printf(
			"[GetTraceTracks] {Validate Param}: trace id is required\n",
		)
		response.SendFail(c, gin.H{"error": "Trace ID is required"})
		return
	}

	result, err := h.service.GetTraceTracks(c.Request.Context(), traceID)
	if err != nil {
		fmt.Printf("[GetTraceTracks] {Fetch Trace Tracks}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve trace tracks",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

type BackupLogRequest struct {
	Status  string `json:"status" binding:"required,oneof=SUCCESS FAILED"`
	Message string `json:"message" binding:"required"`
}

func (h *Handler) PostBackupLog(c *gin.Context) {
	expectedToken := os.Getenv("BACKUP_TOKEN")
	if expectedToken == "" {
		response.SendError(
			c,
			"Backup logging is not configured (BACKUP_TOKEN is empty)",
			http.StatusForbidden,
			nil,
		)
		return
	}

	authHeader := c.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token != expectedToken {
		response.SendError(
			c,
			"Unauthorized backup token",
			http.StatusUnauthorized,
			nil,
		)
		return
	}

	var req BackupLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	level := audit.LevelInfo
	action := "BACKUP_COMPLETED"
	if req.Status == "FAILED" {
		level = audit.LevelError
		action = "BACKUP_FAILED"
	}

	h.service.Record(c.Request.Context(), nil, audit.LogEntry{
		Level:    level,
		Category: audit.CategorySystem,
		Action:   action,
		Message:  req.Message,
		UserID:   structs.StringToNullableString("System"),
	})

	response.SendSuccess(c, gin.H{"status": "recorded"})
}
