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

	analyticsRoutes.GET("/dashboard",
		middleware.RoleMiddleware(constants.AdminRoleID),
		h.GetDashboard,
	)

	analyticsRoutes.GET("/admin-dashboard",
		middleware.RoleMiddleware(
			constants.AdminRoleID,
			constants.SuperAdminRoleID,
		),
		h.GetAdminDashboard,
	)
}
