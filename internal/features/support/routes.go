package support

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	db *sqlx.DB,
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	supportLimiter := middleware.NewIPRateLimiter(1, 5)

	routes := rg.Group("/support")
	{
		// Public support routes (guests / students)
		routes.POST(
			"/tickets",
			middleware.RateLimitMiddleware(supportLimiter),
			h.PostSupportTicket,
		)
		routes.POST(
			"/tickets/:id/messages",
			middleware.RateLimitMiddleware(supportLimiter),
			h.PostSupportMessage,
		)
		routes.GET("/tickets/:id/messages", h.GetSupportTicketMessages)

		// Authenticated student routes
		studentRoutes := routes.Group("")
		studentRoutes.Use(middleware.AuthMiddleware(redis))
		{
			studentRoutes.GET("/my-tickets", h.GetMySupportTickets)
		}

		// Restricted support management routes
		adminRoutes := routes.Group("")
		adminRoutes.Use(middleware.AuthMiddleware(redis))
		adminRoutes.Use(middleware.RoleMiddleware(
			constants.AdminRoleID,
			constants.SuperAdminRoleID,
		))
		{
			adminRoutes.GET("/tickets", h.GetSupportTickets)
			adminRoutes.PATCH(
				"/tickets/:id/status",
				h.PatchSupportTicketStatus,
			)
			adminRoutes.PATCH(
				"/tickets/:id/read",
				h.PatchSupportTicketRead,
			)
		}
	}
}
