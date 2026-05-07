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
	path := c.Param("path")
	if path == "" {
		response.SendFail(c, gin.H{"error": "File path is required"})
		return
	}

	// Ensure path starts with a slash for consistent reconstruction
	if path[0] != '/' {
		path = "/" + path
	}

	// Reconstruct the full file URL as stored in the DB
	fileURL := "/uploads" + path
	log.Printf("[GetFile] Attempting to serve file: %s", fileURL)

	mimeType, err := h.service.DownloadFile(c.Request.Context(), fileURL, c.Writer)
	if err != nil {
		log.Printf("[GetFile] {Service Error}: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found or inaccessible"})
		return
	}

	c.Header("Content-Type", mimeType)
}
