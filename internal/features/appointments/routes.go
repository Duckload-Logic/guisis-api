package appointments

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
	routes := rg.Group("/appointments")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.HydrateStudentIIRContext(db))
	routes.Use(middleware.HydrateStudentCORContext(db))
	routes.Use(middleware.AuditContextMiddleware())
	appointmentLookup := middleware.OwnershipMiddleware(db, "appointmentID")

	adminOnly := routes.Group("")
	adminOnly.Use(middleware.RoleMiddleware(
		constants.AdminRoleID,
	))
	{
		adminOnly.GET("", h.GetAppointments)
		adminOnly.GET("/calendar/stats", h.GetAppointmentDailyStats)
	}

	studentOnly := routes.Group("")
	studentOnly.Use(middleware.RoleMiddleware(
		constants.StudentRoleID,
	))
	{
		studentOnly.GET("/me", h.GetAppointmentMe)
		studentOnly.POST("", h.PostAppointment)
		studentOnly.POST(
			"/id/:appointmentID/cancel",
			appointmentLookup,
			h.PostAppointmentCancel,
		)
	}

	sharedRoutes := routes.Group("")
	sharedRoutes.Use(middleware.RoleMiddleware(
		constants.StudentRoleID,
		constants.AdminRoleID,
	))
	{
		sharedRoutes.GET("/id/:appointmentID", appointmentLookup, h.GetAppointmentByID)
		sharedRoutes.GET("/stats", h.GetAppointmentStats)
		sharedRoutes.GET("/lookups/categories", h.GetAppointmentCategories)
		sharedRoutes.GET("/lookups/slots", h.GetAppointmentSlots)
		sharedRoutes.GET("/lookups/statuses", h.GetAppointmentStatuses)
		sharedRoutes.PATCH(
			"/id/:appointmentID",
			appointmentLookup,
			h.PatchAppointment,
		)
	}
}
