package notes

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
	routes := rg.Group("/notes")
	routes.Use(middleware.AuthMiddleware(redis))
	routes.Use(middleware.HydrateStudentIIRContext(db))
	routes.Use(middleware.RoleMiddleware(constants.AdminRoleID))
	{
		routes.GET("/user/id/:iirID", h.GetSignificantNotes)
		routes.POST("/user/id/:iirID", h.PostSignificantNote)
	}
}
