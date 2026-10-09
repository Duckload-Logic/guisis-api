package notes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Handler struct {
	service *Service
	logger  audit.Logger
}

func NewHandler(
	service *Service,
	logger audit.Logger,
) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) GetSignificantNotes(c *gin.Context) {
	iirID := c.Param("iirID")

	significantNotes, err := h.service.GetStudentSignificantNotes(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		fmt.Printf("[GetSignificantNotes] {Fetch Notes}: %v\n", err)
		response.SendError(
			c,
			"Failed to get student significant notes",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if h.logger != nil {
		id, ip, ua, email, _, trace := audit.ExtractMeta(c.Request.Context())
		h.logger.Record(c.Request.Context(), nil, audit.LogEntry{
			Level:    audit.LevelInfo,
			Category: audit.CategoryAudit,
			Action:   audit.ActionNoteViewed,
			Message: fmt.Sprintf(
				"Confidential counseling notes viewed for IIR #%s",
				iirID,
			),
			UserID:    structs.StringToNullableString(id),
			UserEmail: structs.StringToNullableString(email),
			IPAddress: structs.StringToNullableString(ip),
			UserAgent: structs.StringToNullableString(ua),
			TraceID:   structs.StringToNullableString(trace),
			Metadata: &audit.LogMetadata{
				EntityType: "SignificantNotes",
				EntityID:   iirID,
			},
		})
	}

	response.SendSuccess(c, significantNotes)
}

func (h *Handler) PostSignificantNote(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(
			c,
			gin.H{"error": "IIR ID not found"},
			http.StatusUnauthorized,
		)
		return
	}

	var noteReq SignificantNoteDTO
	if err := c.ShouldBindJSON(&noteReq); err != nil {
		fmt.Printf("[PostSignificantNote] {JSON Bind}: %v\n", err)
		response.SendFail(c, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.service.CreateSignificantNote(
		c.Request.Context(),
		iirID,
		noteReq,
	)
	if err != nil {
		fmt.Printf("[PostSignificantNote] {Save Note}: %v\n", err)
		response.SendError(
			c,
			"Failed to save significant note",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "Significant note saved successfully"},
	)
}

func (h *Handler) DeleteSignificantNote(
	c *gin.Context,
) {

}
