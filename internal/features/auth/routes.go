package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/middleware"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	h *Handler,
	redis *datastore.RedisClient,
) {
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/login", h.PostLogin)
		authRoutes.POST("/refresh", h.PostRefreshToken)
		authRoutes.GET(
			"/me",
			middleware.AuthMiddleware(redis),
			h.GetMe,
		)
		authRoutes.GET(
			"/logout",
			middleware.AuthMiddleware(redis),
			h.GetLogout,
		)

		// IDP OAuth 2.0 routes
		authRoutes.GET("/idp/authorize", h.GetAuthorizeURL)
		authRoutes.POST("/idp/token", h.PostIDPToken)
	}

	// Internal Debugging (Redis Dashboard)
	RegisterDebugRoutes(rg, redis)
}
