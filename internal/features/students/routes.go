package students

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
	// Root group: /api/v1/students
	routes := rg.Group("/students")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.HydrateStudentIIRContext(db))

	// Define lookups
	userResourceLookup := middleware.OwnershipMiddleware(db, "userID")
	iirResourceLookup := middleware.OwnershipMiddleware(db, "iirID")
	// inventoryRecordLookup := middleware.OwnershipMiddleware(db,
	// "inventoryRecordID")

	lookupRoutes := routes.Group("/lookups")
	{
		lookupRoutes.GET("/genders", h.GetGenders)
		lookupRoutes.GET("/enrollment-years", h.GetEnrollmentYears)
		lookupRoutes.GET("/religions", h.GetReligions)
		lookupRoutes.GET("/parental-status-types", h.GetParentalStatusTypes)
		lookupRoutes.GET("/income-ranges", h.GetIncomeRanges)
		lookupRoutes.GET("/activity-options", h.GetActivityOptions)
		lookupRoutes.GET("/support-types", h.GetStudentSupportTypes)
		lookupRoutes.GET("/support-types/siblings", h.GetSiblingSupportTypes)
		lookupRoutes.GET("/courses", h.GetCourses)
		lookupRoutes.GET("/civil-statuses", h.GetCivilStatusTypes)
		lookupRoutes.GET(
			"/nature-of-residence-types",
			h.GetNatureOfResidenceTypes,
		)
		lookupRoutes.GET(
			"/student-relationship-types",
			h.GetStudentRelationshipTypes,
		)
		lookupRoutes.GET("/student-statuses", h.GetStudentStatuses)
		lookupRoutes.GET("/educational-levels", h.GetEducationalLevels)
		lookupRoutes.GET("/educational-attainments", h.GetEducationalAttainments)
	}

	inventoryRoutes := routes.Group("/inventory")

	counselorRoutes := inventoryRoutes.Group("/")
	counselorRoutes.Use(middleware.RoleMiddleware(
		constants.AdminRoleID,
	))
	{
		counselorRoutes.GET("/records", h.GetStudentList)
		counselorRoutes.PATCH("/records/bulk-status", h.PatchStudentStatusBulk)
	}

	userRoutes := inventoryRoutes.Group("/")
	userRoutes.Use(middleware.RoleMiddleware(
		constants.StudentRoleID,
		constants.AdminRoleID,
	))
	{
		userRoutes.GET(
			"/records/user/:userID",
			userResourceLookup,
			h.GetStudentIIRByUserID,
		)
		userRoutes.GET(
			"/records/iir/:iirID",
			iirResourceLookup,
			h.GetStudentIIRByIIRID,
		)
		userRoutes.PATCH(
			"/records/iir/:iirID",
			iirResourceLookup,
			h.PatchStudentIIR,
		)
		userRoutes.GET(
			"/records/iir/:iirID/profile",
			iirResourceLookup,
			h.GetStudentProfile,
		)
		userRoutes.GET(
			"/records/iir/:iirID/basic-info",
			iirResourceLookup,
			h.GetStudentBasicInfo,
		)
		userRoutes.GET(
			"/records/iir/:iirID/personal-info",
			iirResourceLookup,
			h.GetStudentPersonalInfo,
		)
		userRoutes.GET(
			"/records/iir/:iirID/addresses",
			iirResourceLookup,
			h.GetStudentAddresses,
		)
		userRoutes.GET(
			"/records/iir/:iirID/family-background",
			iirResourceLookup,
			h.GetStudentFamilyBackground,
		)
		userRoutes.GET(
			"/records/iir/:iirID/related-persons",
			iirResourceLookup,
			h.GetStudentRelatedPersons,
		)
		userRoutes.GET(
			"/records/iir/:iirID/education",
			iirResourceLookup,
			h.GetStudentEducationalBackground,
		)
		userRoutes.GET(
			"/records/iir/:iirID/finance",
			iirResourceLookup,
			h.GetStudentFinancialInfo,
		)
		userRoutes.GET(
			"/records/iir/:iirID/health",
			iirResourceLookup,
			h.GetStudentHealthRecord,
		)
		userRoutes.GET(
			"/records/iir/:iirID/consultations",
			iirResourceLookup,
			h.GetStudentConsultations,
		)
		userRoutes.GET(
			"/records/iir/:iirID/activities",
			iirResourceLookup,
			h.GetStudentActivities,
		)
		userRoutes.GET(
			"/records/iir/:iirID/subject-preferences",
			iirResourceLookup,
			h.GetStudentSubjectPreferences,
		)
		userRoutes.GET(
			"/records/iir/:iirID/hobbies",
			iirResourceLookup,
			h.GetStudentHobbies,
		)
		userRoutes.GET(
			"/records/iir/:iirID/test-results",
			iirResourceLookup,
			h.GetStudentTestResults,
		)
		userRoutes.GET(
			"/records/iir/:iirID/download",
			iirResourceLookup,
			h.GetStudentIIRPDF,
		)
	}

	studentRoutes := inventoryRoutes.Group("/")
	studentRoutes.Use(middleware.RoleMiddleware(
		constants.StudentRoleID,
	))
	{
		studentRoutes.GET("/records/iir/draft", h.GetStudentIIRDraft)
		studentRoutes.POST("/records/iir/draft", h.PostStudentIIRDraft)

		studentRoutes.POST("/records/iir", h.PostStudentIIR)

		// COR management
		studentRoutes.POST("/records/iir/:iirID/cor", h.PostStudentCORByIIRID)
		studentRoutes.POST("/cors", h.PostStudentCOR)
		studentRoutes.GET(
			"/cors/user/:userID",
			userResourceLookup,
			h.GetStudentCORByUserID,
		)
		studentRoutes.GET("/cors", h.GetStudentCORs)
	}

	// Counselor view of CORs
	counselorRoutes.GET("/cors/user/:userID/current", h.GetStudentCORByUserID)

	// Academic settings — GET is open to all authenticated roles;
	// PUT is restricted to SuperAdmin only.
	routes.GET("/settings/academic", h.GetAcademicSetting)

	counselorSettingsRoutes := routes.Group("/settings")
	counselorSettingsRoutes.Use(
		middleware.RoleMiddleware(constants.AdminRoleID),
	)
	{
		counselorSettingsRoutes.PUT("/academic", h.PutAcademicSetting)
	}
}

