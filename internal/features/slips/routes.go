package slips

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
	routes := rg.Group("/slips")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.HydrateStudentIIRContext(db))
	routes.Use(middleware.HydrateStudentCORContext(db))
	routes.Use(middleware.AuditContextMiddleware())
	slipLookup := middleware.OwnershipMiddleware(db, "slipID")

	adminOnly := routes.Group("")
	adminOnly.Use(middleware.RoleMiddleware(
		constants.AdminRoleID,
	))
	{
		adminOnly.GET("", h.GetSlips)
		adminOnly.GET("/urgent", h.GetSlipUrgent)
		adminOnly.PATCH("/id/:slipID/status", h.PatchSlipStatus)
	}

	studentOnly := routes.Group("")
	studentOnly.Use(middleware.RoleMiddleware(
		constants.StudentRoleID,
	))
	{
		studentOnly.GET("/me", h.GetSlipMe)
		studentOnly.POST("", h.PostSlip)
		studentOnly.PATCH("/id/:slipID", slipLookup, h.PatchSlip)
	}

	sharedRoutes := routes.Group("")
	sharedRoutes.Use(middleware.RoleMiddleware(
		constants.AdminRoleID,
		constants.StudentRoleID,
	))
	{
		sharedRoutes.GET("/id/:slipID", slipLookup, h.GetSlipByID)
		sharedRoutes.GET("/stats", h.GetSlipStats)
		sharedRoutes.GET(
			"/id/:slipID/attachments",
			slipLookup,
			h.GetSlipAttachments,
		)
		sharedRoutes.GET(
			"/id/:slipID/attachments/:attachmentId",
			slipLookup,
			h.GetSlipAttachmentContent,
		)
		sharedRoutes.GET(
			"/lookups/statuses",
			h.GetSlipStatuses,
		)
		sharedRoutes.GET(
			"/lookups/categories",
			h.GetSlipCategories,
		)
	}
}
