package analytics

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
	analyticsRoutes := rg.Group("/analytics")
	analyticsRoutes.Use(middleware.AuditContextMiddleware())
	analyticsRoutes.Use(middleware.AuthMiddleware(redis))

	analyticsRoutes.GET("",
		middleware.RoleMiddleware(
			constants.AdminRoleID,
		),
		h.GetAdminDashboard,
	)

	analyticsRoutes.GET("/admin-dashboard",
		middleware.RoleMiddleware(
			constants.AdminRoleID,
		),
		h.GetAdminDashboard,
	)

	analyticsRoutes.GET("/reports/iir",
		middleware.RoleMiddleware(
			constants.AdminRoleID,
		),
		h.GetIIRAnalyticsReport,
	)

	analyticsRoutes.GET("/reports/iir/export",
		middleware.RoleMiddleware(
			constants.AdminRoleID,
		),
		h.ExportIIRAnalyticsReport,
	)
}
