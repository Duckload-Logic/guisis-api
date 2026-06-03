package files

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetFile(c *gin.Context) {
	// The path will be everything after /uploads/
	// e.g., /api/v1/uploads/cors/abc.pdf -> path is /cors/abc.pdf
	path := c.Param("path")
	if path == "" {
		response.SendFail(c, gin.H{"error": "File path is required"})
		return
	}

	// Reconstruct the full file URL as stored in the DB
	fileURL := "/uploads" + path

	mimeType, err := h.service.DownloadFile(
		c.Request.Context(),
		fileURL,
		c.Writer,
	)
	if err != nil {
		log.Printf("[GetFile] {Service Call}: %v", err)
		c.JSON(
			http.StatusNotFound,
			gin.H{"error": "File not found or inaccessible"},
		)
		return
	}

	c.Header("Content-Type", mimeType)
}
