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
	// Base group for all activity logs
	activityGroup := rg.Group("/activity-meta")
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
		adminOnly.GET("/stats", h.GetLogsStats)
		adminOnly.GET("/activity-stats", h.GetLogsActivity)
		adminOnly.POST("/cleanup", h.PostLogsCleanup)
	}
}
