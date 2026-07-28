package users

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
	userRoutes := rg.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware(redis))
	userRoutes.Use(middleware.RoleMiddleware(
		constants.SuperAdminRoleID,
		constants.AdminRoleID,
		constants.StudentRoleID,
		constants.DeveloperRoleID,
	))

	userRoutes.GET("/me", h.GetMe)
	userRoutes.POST("/profile-picture/:id", h.PostProfilePicture)
	userRoutes.POST("/profile-picture/upload", h.UploadProfilePicture)

	adminOnly := userRoutes.Group("")
	adminOnly.Use(middleware.RoleMiddleware(
		constants.SuperAdminRoleID,
		constants.AdminRoleID,
	))
	{
		adminOnly.GET("", h.GetUsers)
		adminOnly.GET("/search", h.GetUserByEmail)
	}

	userRoutes.GET("/distribution",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.GetRoleDistribution,
	)
	userRoutes.POST("/:id/block",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.PostUserBlock,
	)
	userRoutes.POST("/:id/unblock",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.PostUserUnblock,
	)

	// Session & Activity Audit (Super Admin only)
	userRoutes.GET("/:id/sessions",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.GetUserSessions,
	)
	userRoutes.DELETE("/:id/sessions/:session_id",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.DeleteUserSession,
	)
	userRoutes.GET("/:id/activity",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.GetUserActivity,
	)
	userRoutes.POST("/update-roles",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.PostUpdateRoles,
	)
	userRoutes.POST("/whitelist",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.PostUserToWhitelist,
	)
	userRoutes.POST("/whitelist/remove",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.PostRemoveUserFromWhitelist,
	)
	userRoutes.GET("/whitelist",
		middleware.RoleMiddleware(constants.SuperAdminRoleID),
		h.GetWhitelist,
	)
}
