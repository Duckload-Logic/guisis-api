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
			constants.SuperAdminRoleID,
			constants.DeveloperRoleID,
		),
		h.GetAdminDashboard,
	)

	analyticsRoutes.GET("/dashboard",
		middleware.RoleMiddleware(
			constants.AdminRoleID,
			constants.SuperAdminRoleID,
			constants.DeveloperRoleID,
		),
		h.GetDashboard,
	)
}
