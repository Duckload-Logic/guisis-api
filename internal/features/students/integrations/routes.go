package integrations

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
	routes := rg.Group("/integrations/students")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.RequireVerifiedM2M)
	routes.Use(middleware.M2MAuditMiddleware())

	routes.GET("/profiles", h.GetStudents)
	routes.GET("/profile", h.GetStudentByEmail)

	// Student number lookups
	routes.GET("/:studentNumber", h.GetStudentByStudentNumber)
	routes.GET(
		"/:studentNumber/personal-info",
		middleware.RequireM2MPersonalInfoAccess,
		h.GetPersonalInfoByStudentNumber,
	)
	routes.GET(
		"/:studentNumber/addresses",
		middleware.RequireM2MPersonalInfoAccess,
		h.GetAddressByStudentNumber,
	)

	// IDP UUID lookups
	routes.GET("/idp/:idpUuid", h.GetStudentByIDPUUID)
	routes.GET(
		"/idp/:idpUuid/personal-info",
		middleware.RequireM2MPersonalInfoAccess,
		h.GetPersonalInfoByIDPUUID,
	)
	routes.GET(
		"/idp/:idpUuid/addresses",
		middleware.RequireM2MPersonalInfoAccess,
		h.GetAddressByIDPUUID,
	)
}
