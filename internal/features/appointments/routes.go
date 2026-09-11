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
		adminOnly.POST("/id/:appointmentID/start", h.PostAppointmentStart)
	}

	appointmentRoutes := routes.Group("")
	appointmentRoutes.Use(middleware.RoleMiddleware(
		constants.StudentRoleID,
		constants.AdminRoleID,
	))
	{
		appointmentRoutes.GET("/me", h.GetAppointmentsMe)
		appointmentRoutes.POST("", middleware.RequireCOR(), h.PostAppointment)
		appointmentRoutes.POST(
			"/id/:appointmentID/cancel",
			appointmentLookup,
			middleware.RequireCOR(),
			h.PostAppointmentCancellation,
		)
		appointmentRoutes.GET(
			"/id/:appointmentID",
			appointmentLookup,
			h.GetAppointmentByID,
		)
		appointmentRoutes.PATCH(
			"/id/:appointmentID",
			appointmentLookup,
			middleware.RequireCOR(),
			h.PatchAppointment,
		)
	}

	sharedRoutes := routes.Group("")
	sharedRoutes.Use(middleware.RoleMiddleware(
		constants.StudentRoleID,
		constants.AdminRoleID,
	))
	{
		sharedRoutes.GET("/stats", h.GetAppointmentStats)
	}

	sharedLookups := routes.Group("/lookups")
	sharedLookups.Use(middleware.RoleMiddleware(
		constants.StudentRoleID,
		constants.AdminRoleID,
		constants.SuperAdminRoleID,
		constants.DeveloperRoleID,
	))
	{
		sharedLookups.GET("/categories", h.GetAppointmentCategories)
		sharedLookups.GET("/slots", h.GetAppointmentSlots)
		sharedLookups.GET("/statuses", h.GetAppointmentStatuses)
	}
}
