package analytics

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetIIRAnalyticsReport(c *gin.Context) {
	yearStr := c.DefaultQuery("year", "0")
	courseIDStr := c.DefaultQuery("course_id", "0")

	var year, courseID int
	fmt.Sscanf(yearStr, "%d", &year)
	fmt.Sscanf(courseIDStr, "%d", &courseID)

	dashboardData, err := h.service.GetIIRAnalyticsReport(
		c.Request.Context(),
		year,
		courseID,
	)
	if err != nil {
		fmt.Printf("[GetIIRAnalyticsReport] {Fetch Data}: %v\n", err)
		response.SendError(
			c,
			"Failed to generate analytics report",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, dashboardData)
}

func (h *Handler) ExportIIRAnalyticsReport(c *gin.Context) {
	yearStr := c.DefaultQuery("year", "0")
	courseIDStr := c.DefaultQuery("course_id", "0")

	var year, courseID int
	fmt.Sscanf(yearStr, "%d", &year)
	fmt.Sscanf(courseIDStr, "%d", &courseID)

	pdfBytes, err := h.service.ExportIIRAnalyticsReport(
		c.Request.Context(),
		year,
		courseID,
	)
	if err != nil {
		fmt.Printf("[ExportIIRAnalyticsReport] {Generate PDF}: %v\n", err)
		response.SendError(
			c,
			"Failed to export analytics report",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header(
		"Content-Disposition",
		"attachment; filename=iir_analytics_report.pdf",
	)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *Handler) GetAdminDashboard(c *gin.Context) {
	timeRange := c.Query("filter")
	if timeRange == "" {
		timeRange = c.DefaultQuery("range", "monthly")
	}
	source := c.DefaultQuery("source", "appointments")

	dashboardData, err := h.service.GetAdminDashboard(
		c.Request.Context(),
		timeRange,
		source,
	)
	if err != nil {
		fmt.Printf("[GetAdminDashboard] {Fetch Statistics}: %v\n", err)
		response.SendError(
			c,
			"Failed to generate admin analytics dashboard",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, dashboardData)
}
