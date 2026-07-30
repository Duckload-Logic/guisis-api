package logs

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	// Public/Token-authenticated backup route
	rg.POST("/logs/system/backup", h.PostBackupLog)

	// Base group for all activity logs
	activityGroup := rg.Group("/logs")
	activityGroup.Use(middleware.AuthMiddleware(redis))

	// User-specific activity route (No role check, just auth)
	activityGroup.GET("/me", h.GetLogsMe)

	// Admin-only routes (Requires SuperAdmin role)
	adminOnly := activityGroup.Group("")
	adminOnly.Use(middleware.RoleMiddleware(constants.SuperAdminRoleID))
	{
		adminOnly.GET("", h.GetLogs)
		adminOnly.GET("/audit", h.GetLogsAudit)
		adminOnly.GET("/system", h.GetLogsSystem)
		adminOnly.GET("/security", h.GetLogsSecurity)
		adminOnly.GET("/m2m", h.GetLogsM2M)
		adminOnly.GET("/stats", h.GetLogsStats)
		adminOnly.GET("/activity", h.GetLogsActivity)
		adminOnly.GET("/trace/:traceId", h.GetTraceTracks)
		adminOnly.GET("/:id", h.GetLog)
		// adminOnly.POST("/cleanup", h.PostLogsCleanup)
	}
}
