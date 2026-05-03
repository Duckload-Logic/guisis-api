package slips

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service *Service
}

// NewHandler creates a new slips handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// getIIRIDFromContext extracts iirID from context or aborts with Forbidden
// status if not found.
func getIIRIDFromContext(c *gin.Context) (string, bool) {
	iirIDVal, exists := c.Get("iirID")
	if !exists {
		response.SendFail(
			c,
			gin.H{"error": "Please complete your IIR profile"},
			http.StatusForbidden,
		)
		return "", false
	}

	iirID, ok := iirIDVal.(string)
	if !ok {
		response.SendError(
			c,
			"Internal server error",
			http.StatusInternalServerError,
			nil,
		)
		return "", false
	}

	return iirID, true
}

// PostSlip godoc
// @Summary      Submit an excuse slip
// @Description  Allows a student to submit an excuse slip with
// @Description  supporting document (file upload).
// @Tags         ExcuseSlips
// @Accept       multipart/form-data
// @Produce      json
// @Param        reason      formData string true "Reason for absence"
// @Param        absenceDate formData string true "Date (YYYY-MM-DD)"
// @Param        files       formData file   true "Supporting Document"
// @Success      201         {object} map[string]interface{}
// @Failure      400         {object} map[string]string
// @Failure      403         {object} map[string]string
// @Failure      500         {object} map[string]string
// @Router       /slips [post]
// PostSlip handles the submission of student excuse slips.
func (h *Handler) PostSlip(c *gin.Context) {
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req CreateSlipRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("[PostSlip] {Bind Request}: %v\n", err)
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		fmt.Printf("[PostSlip] {Parse Form}: %v\n", err)
		response.SendFail(c, gin.H{"error": "Failed to parse form"})
		return
	}

	var files []*multipart.FileHeader
	fieldNames := []string{
		"files", "cor", "excuseLetter", "parentId", "medicalCert",
	}

	for _, field := range fieldNames {
		if f := form.File[field]; len(f) > 0 {
			files = append(files, f...)
		}
	}

	if len(files) == 0 {
		response.SendFail(c, gin.H{"error": "At least one file required"})
		return
	}

	slip, err := h.service.SubmitExcuseSlip(
		c.Request.Context(), iirID, req, files,
	)
	if err != nil {
		fmt.Printf("[PostSlip] {Submit Slip}: %v\n", err)
		response.SendError(
			c,
			"Failed to submit slip",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"message": "Excuse slip submitted successfully",
		"slipId":  slip.ID,
	}, http.StatusCreated)
}

// GetUrgentSlipList godoc
// @Summary      Get urgent excuse slips
// @Description  Retrieves urgent slips for counselor review.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]string
// @Router       /slips/urgent [get]
func (h *Handler) GetSlipUrgent(c *gin.Context) {
	var req ListSlipsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	slips, err := h.service.GetUrgentSlips(c.Request.Context(), &req)
	if err != nil {
		fmt.Printf("[GetSlipUrgent] {Fetch Slips}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve slips",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, slips)
}

// GetSlipStatsList godoc
// @Summary      Get slip statistics
// @Description  Retrieves slip status counts.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} []SlipStatusCount
// @Failure      500  {object} map[string]string
// @Router       /slips/stats [get]
func (h *Handler) GetSlipStats(c *gin.Context) {
	var req ListSlipsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	roleIDs := c.MustGet("roleIDs").([]int)
	var iirIDPtr *string

	isStudent := false
	for _, rid := range roleIDs {
		if rid == int(constants.StudentRoleID) {
			isStudent = true
			break
		}
	}

	if isStudent {
		iirID, ok := getIIRIDFromContext(c)
		if !ok {
			return
		}
		iirIDPtr = &iirID
	}

	stats, err := h.service.GetSlipStats(c.Request.Context(), iirIDPtr, &req)
	if err != nil {
		fmt.Printf("[GetSlipStats] {Fetch Stats}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve statistics",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, stats)
}

// GetSlipStatusList godoc
// @Summary      Get slip statuses
// @Description  Retrieves available slip statuses.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} []SlipStatus
// @Failure      500  {object} map[string]string
// @Router       /slips/lookups/statuses [get]
func (h *Handler) GetSlipStatuses(c *gin.Context) {
	statuses, err := h.service.GetSlipStatuses(c.Request.Context())
	if err != nil {
		fmt.Printf("[GetSlipStatuses] {Fetch Statuses}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve statuses",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, statuses)
}

// GetSlipCategoryList godoc
// @Summary      Get slip categories
// @Description  Retrieves available slip categories.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} []SlipCategory
// @Failure      500  {object} map[string]string
// @Router       /slips/lookups/categories [get]
func (h *Handler) GetSlipCategories(c *gin.Context) {
	categories, err := h.service.GetSlipCategories(c.Request.Context())
	if err != nil {
		fmt.Printf("[GetSlipCategories] {Fetch Categories}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve categories",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, categories)
}

// GetSlipList godoc
// @Summary      Get all excuse slips
// @Description  Retrieves list of all submitted slips.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      500  {object} map[string]string
// @Router       /slips [get]
func (h *Handler) GetSlips(c *gin.Context) {
	var req ListSlipsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	slips, err := h.service.GetAllExcuseSlips(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetSlips] {Fetch Slips}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve slips",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, slips)
}

// GetSlipListByIIR godoc
// @Summary      Get student's excuse slips
// @Description  Retrieves slips for the authenticated student.
// @Tags         ExcuseSlips
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/me [get]
func (h *Handler) GetSlipMe(c *gin.Context) {
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req ListSlipsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	slips, err := h.service.GetExcuseSlipsByIIRID(c.Request.Context(), iirID, req)
	if err != nil {
		fmt.Printf("[GetSlipMe] {Fetch Slips}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve slips",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, slips)
}

// GetSlipByID godoc
// @Summary      Get excuse slip by ID
// @Description  Retrieves details for a specific excuse slip.
// @Tags         ExcuseSlips
// @Produce      json
// @Param        id   path      string  true  "Slip ID"
// @Success      200  {object} SlipDTO
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/id/{id} [get]
func (h *Handler) GetSlipByID(c *gin.Context) {
	idParam := c.Param("slipID")
	slip, err := h.service.GetSlipByID(c.Request.Context(), idParam)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			response.SendFail(
				c,
				gin.H{"error": "Excuse slip not found"},
				http.StatusNotFound,
			)
			return
		}
		fmt.Printf("[GetSlipByID] {Fetch Slip}: %v\n", err)
		response.SendError(
			c,
			"Internal server error",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, slip)
}

// GetSlipAttachmentList godoc
// @Summary      Get slip attachments
// @Description  Retrieves attachments for a specific slip.
// @Tags         ExcuseSlips
// @Produce      json
// @Param        id   path      int  true  "Slip ID"
// @Success      200  {object} []AttachmentDTO
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/id/{id}/attachments [get]
func (h *Handler) GetSlipAttachments(c *gin.Context) {
	idParam := c.Param("slipID")
	attachments, err := h.service.GetSlipAttachments(c.Request.Context(), idParam)
	if err != nil {
		fmt.Printf("[GetSlipAttachments] {Fetch Attachments}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve attachments",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, attachments)
}

// GetAttachmentFile godoc
// @Summary      Download attachment
// @Description  Downloads a specific attachment file.
// @Tags         ExcuseSlips
// @Param        attachmentId path int true "Attachment ID"
// @Success      200  {file} binary
// @Failure      400  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/id/{id}/attachments/{attachmentId} [get]
func (h *Handler) GetSlipAttachmentContent(c *gin.Context) {
	attachmentIDParam := c.Param("attachmentId")
	attachment, err := h.service.DownloadAttachment(
		c.Request.Context(), attachmentIDParam, c.Writer,
	)
	if err != nil {
		if strings.Contains(err.Error(), "attachment not found") {
			response.SendFail(
				c,
				gin.H{"error": "Attachment not found"},
				http.StatusNotFound,
			)
			return
		}
		fmt.Printf("[GetSlipAttachmentContent] {Download}: %v\n", err)
		response.SendError(
			c,
			"Failed to download file",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.Header(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=%q", attachment.FileName),
	)
}

// PatchSlipStatus godoc
// @Summary      Update slip status
// @Description  Approve, reject, or request revision of slip.
// @Tags         ExcuseSlips
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "Slip ID"
// @Param        body body      UpdateStatusRequest  true  "Status and notes"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /slips/id/{id}/status [patch]
func (h *Handler) PatchSlip(c *gin.Context) {
	idParam := c.Param("slipID")
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req CreateSlipRequest
	if err := c.ShouldBind(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		response.SendFail(c, gin.H{"error": "Failed to parse form"})
		return
	}

	var files []*multipart.FileHeader
	fieldNames := []string{
		"files",
		"cor",
		"excuseLetter",
		"parentId",
		"medicalCert",
	}
	for _, field := range fieldNames {
		if f := form.File[field]; len(f) > 0 {
			files = append(files, f...)
		}
	}

	if len(files) == 0 {
		response.SendFail(c, gin.H{"error": "At least one file required"})
		return
	}

	slip, err := h.service.UpdateExcuseSlip(
		c.Request.Context(), iirID, idParam, req, files,
	)
	if err != nil {
		fmt.Printf("[PatchSlip] {Update Slip}: %v\n", err)
		response.SendError(
			c,
			err.Error(),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"message": "Excuse slip updated successfully",
		"slipId":  slip.ID,
	}, http.StatusOK)
}

// PatchSlipStatus godoc
func (h *Handler) PatchSlipStatus(c *gin.Context) {
	idParam := c.Param("slipID")
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	err := h.service.UpdateExcuseSlipStatus(
		c.Request.Context(),
		idParam,
		req.Status,
		req.AdminNotes,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			response.SendFail(
				c,
				gin.H{"error": "Slip not found"},
				http.StatusNotFound,
			)
			return
		}
		if strings.Contains(err.Error(), "invalid status") {
			response.SendFail(c, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("[PatchSlipStatus] {Update Status}: %v\n", err)
		response.SendError(
			c,
			"Failed to update status",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "Status updated successfully"})
}
