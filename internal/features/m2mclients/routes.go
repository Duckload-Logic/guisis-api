package m2mclients

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// RegisterRoutes sets up management and auth endpoints for M2M clients.
func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	// Public M2M Auth Routes with brute-force rate limiting (1 rps, 5 burst)
	m2mAuthLimiter := middleware.NewIPRateLimiter(1, 5)
	m2mAuth := rg.Group("/auth/m2m")
	m2mAuth.Use(middleware.RateLimitMiddleware(m2mAuthLimiter))
	{
		m2mAuth.POST("/token", h.PostM2MToken)
		m2mAuth.POST("/refresh", h.PostM2MTokenRefresh)
	}

	// Protected Management Routes
	m2mMgmt := rg.Group("/m2m-clients")
	m2mMgmt.Use(
		middleware.AuthMiddleware(redis),
		middleware.AuditContextMiddleware(),
	)
	{
		// Common routes for both Developer and Superadmin
		common := m2mMgmt.Group("")
		common.Use(middleware.RoleMiddleware(
			constants.SuperAdminRoleID,
			constants.DeveloperRoleID,
		))
		{
			common.GET("/me", h.GetMyM2MClient)
			common.POST("/:id/secret", h.PostM2MClientSecret)
			common.DELETE("/:id", h.DeleteM2MClient)
		}

		// Admin-only routes
		adminOnly := m2mMgmt.Group("")
		adminOnly.Use(
			middleware.RoleMiddleware(constants.SuperAdminRoleID),
		)
		{
			adminOnly.GET("", h.GetM2MClients)
			adminOnly.PATCH("/:id/verify", h.PatchM2MClientVerify)
			adminOnly.PATCH("/:id/reject", h.PatchM2MClientReject)
		}

		devOnly := m2mMgmt.Group("")
		devOnly.Use(
			middleware.RoleMiddleware(constants.DeveloperRoleID),
		)
		{
			devOnly.POST("", h.PostM2MClient)
		}
	}
}
