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
		authRoutes.GET("/logout", h.GetLogout)

		// IDP OAuth 2.0 routes
		authRoutes.GET("/idp/authorize", h.GetAuthorizeURL)
		authRoutes.POST("/idp/token", h.PostIDPToken)

		// OTP Fallback routes
		authRoutes.POST("/otp/request", h.PostOTPRequest)
		authRoutes.POST("/otp/login", h.PostOTPLogin)
	}

	// Internal Debugging (Redis Dashboard)
	RegisterDebugRoutes(rg, redis)
}
