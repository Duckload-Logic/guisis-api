package students

import (
	"encoding/json"
	"fmt"
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

func (h *Handler) GetGenders(c *gin.Context) {
	genders, err := h.service.GetGenders(c.Request.Context())
	if err != nil {
		log.Printf("[GetGenders] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get genders",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, genders)
}

func (h *Handler) GetParentalStatusTypes(c *gin.Context) {
	statuses, err := h.service.GetParentalStatusTypes(c.Request.Context())
	if err != nil {
		log.Printf("[GetParentalStatusTypes] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get parental status types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, statuses)
}

func (h *Handler) GetIncomeRanges(c *gin.Context) {
	ranges, err := h.service.GetIncomeRanges(c.Request.Context())
	if err != nil {
		log.Printf("[GetIncomeRanges] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get income ranges",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, ranges)
}

func (h *Handler) GetStudentSupportTypes(c *gin.Context) {
	supportTypes, err := h.service.GetStudentSupportTypes(c.Request.Context())
	if err != nil {
		log.Printf("[GetStudentSupportTypes] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student support types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, supportTypes)
}

func (h *Handler) GetSiblingSupportTypes(c *gin.Context) {
	supportTypes, err := h.service.GetSiblingSupportTypes(c.Request.Context())
	if err != nil {
		log.Printf("[GetSiblingSupportTypes] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get sibling support types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, supportTypes)
}

func (h *Handler) GetEducationalLevels(c *gin.Context) {
	levels, err := h.service.GetEducationalLevels(c.Request.Context())
	if err != nil {
		log.Printf("[GetEducationalLevels] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get educational levels",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, levels)
}

func (h *Handler) GetCourses(c *gin.Context) {
	courses, err := h.service.GetCourses(c.Request.Context())
	if err != nil {
		log.Printf("[GetCourses] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get courses",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, courses)
}

func (h *Handler) GetCivilStatusTypes(c *gin.Context) {
	civilStatusTypes, err := h.service.GetCivilStatusTypes(c.Request.Context())
	if err != nil {
		log.Printf("[GetCivilStatusTypes] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get civil status types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, civilStatusTypes)
}

func (h *Handler) GetReligions(c *gin.Context) {
	religions, err := h.service.GetReligions(c.Request.Context())
	if err != nil {
		log.Printf("[GetReligions] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get religions",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, religions)
}

func (h *Handler) GetStudentRelationshipTypes(c *gin.Context) {
	relationshipTypes, err := h.service.GetStudentRelationshipTypes(
		c.Request.Context(),
	)
	if err != nil {
		log.Printf("[GetStudentRelationshipTypes] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student relationship types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, relationshipTypes)
}

func (h *Handler) GetStudentStatuses(c *gin.Context) {
	statuses, err := h.service.GetStudentStatuses(c.Request.Context())
	if err != nil {
		log.Printf("[GetStudentStatuses] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student statuses",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, statuses)
}

func (h *Handler) GetNatureOfResidenceTypes(c *gin.Context) {
	types, err := h.service.GetNatureOfResidenceTypes(c.Request.Context())
	if err != nil {
		log.Printf("[GetNatureOfResidenceTypes] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get nature of residence types",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, types)
}

func (h *Handler) GetActivityOptions(c *gin.Context) {
	options, err := h.service.GetActivityOptions(c.Request.Context())
	if err != nil {
		log.Printf("[GetActivityOptions] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get activity options",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, options)
}

// GetStudentList godoc
// @Summary      List Students
// @Description  Retrieves a paginated list of students with optional filters.
// @Tags         Students
// @Accept       json
// @Produce      json
// @Param        course     query     string  false  "Filter by course"
// @Param        year_level query     int     false  "Filter by year level"
// @Param		gender_id query     int     false  "Filter by gender ID"
// @Param        page       query     int     false  "Page number" default(1)
// @Param        page_size  query     int     false  "Number of records per page" default(20)
// @Success      200        {object}  ListStudentsResponse
// @Failure      500        {object}  map[string]string "Failed to list students"
// @Router       /students [get]
func (h *Handler) GetStudentList(c *gin.Context) {
	var req ListStudentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.ListStudents(c.Request.Context(), req)
	if err != nil {
		log.Printf("[GetStudentList] {Service Error}: %v", err)
		response.SendError(
			c,
			"Failed to list students",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

func (h *Handler) GetStudentProfile(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	profile, err := h.service.GetStudentProfile(c.Request.Context(), iirID)
	if err != nil {
		log.Printf("[GetStudentProfile] {Service Error}: %v", err)
		response.SendError(
			c,
			"Failed to get student profile",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, profile)
}

func (h *Handler) GetStudentBasicInfo(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	basicInfo, err := h.service.GetStudentBasicInfo(c.Request.Context(), iirID)
	if err != nil {
		log.Printf("[GetStudentBasicInfo] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student basic info",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, basicInfo)
}

func (h *Handler) GetStudentIIRDraft(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	draft, err := h.service.GetIIRDraft(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetStudentIIRDraft] {Fetch Draft}: %v", err)
		response.SendError(
			c,
			"Failed to get IIR draft",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, draft)
}

func (h *Handler) GetStudentIIRByUserID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		response.SendFail(c, gin.H{"error": "User ID is required"})
		return
	}

	iir, err := h.service.GetStudentIIRByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetStudentIIRByUserID] {Fetch IIR}: %v", err)
		response.SendError(
			c,
			"Failed to get student IIR by user ID",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, iir)
}

func (h *Handler) GetStudentIIRByIIRID(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	iir, err := h.service.GetStudentIIR(c.Request.Context(), iirID)
	if err != nil {
		log.Printf("[GetStudentIIRByIIRID] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student IIR by IIR ID",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, iir)
}

func (h *Handler) GetStudentPersonalInfo(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	personalInfo, err := h.service.GetStudentPersonalInfo(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentPersonalInfo] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student personal info",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, personalInfo)
}

func (h *Handler) GetStudentAddresses(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	addresses, err := h.service.GetStudentAddresses(c.Request.Context(), iirID)
	if err != nil {
		log.Printf("[GetStudentAddresses] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student addresses",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, addresses)
}

func (h *Handler) GetStudentFamilyBackground(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	familyBackground, err := h.service.GetStudentFamilyBackground(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentFamilyBackground] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student family background",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, familyBackground)
}

func (h *Handler) GetStudentRelatedPersons(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	relatedPersons, err := h.service.GetStudentRelatedPersons(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentRelatedPersons] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student related persons",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, relatedPersons)
}

func (h *Handler) GetStudentEducationalBackground(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	educationalBackground, err := h.service.GetEducationalBackground(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentEducationalBackground] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get educational background",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, educationalBackground)
}

func (h *Handler) GetStudentFinancialInfo(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	financialInfo, err := h.service.GetStudentFinancialInfo(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentFinancialInfo] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student financial info",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, financialInfo)
}

func (h *Handler) GetStudentHealthRecord(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	healthRecord, err := h.service.GetStudentHealthRecord(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		response.SendError(
			c,
			"Failed to get student health record",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, healthRecord)
}

func (h *Handler) GetStudentConsultations(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	consultations, err := h.service.GetStudentConsultations(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentConsultations] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student consultations",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, consultations)
}

func (h *Handler) GetStudentActivities(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	activities, err := h.service.GetStudentActivities(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentActivities] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student activities",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, activities)
}

func (h *Handler) GetStudentSubjectPreferences(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	subjectPreferences, err := h.service.GetStudentSubjectPreferences(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentSubjectPreferences] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student subject preferences",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, subjectPreferences)
}

func (h *Handler) GetStudentHobbies(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	hobbies, err := h.service.GetStudentHobbies(c.Request.Context(), iirID)
	if err != nil {
		log.Printf("[GetStudentHobbies] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student hobbies",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, hobbies)
}

func (h *Handler) GetStudentTestResults(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	testResults, err := h.service.GetStudentTestResults(
		c.Request.Context(),
		iirID,
	)
	if err != nil {
		log.Printf("[GetStudentTestResults] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to get student test results",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, testResults)
}

func (h *Handler) PostStudentIIRDraft(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	var req ComprehensiveProfileDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		log.Printf("[PostStudentIIRDraft] {JSON Decode}: %v", err)
		response.SendFail(c, gin.H{"error": "Invalid JSON format"})
		return
	}

	draftID, err := h.service.SaveIIRDraft(c.Request.Context(), userID, req)
	if err != nil {
		log.Printf("[PostStudentIIRDraft] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to save IIR draft",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"id": draftID, "message": "IIR draft saved successfully"},
	)
}

func (h *Handler) PostStudentIIR(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	var req ComprehensiveProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[PostStudentIIR] {JSON Bind}: %s", err.Error())
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	iirID, err := h.service.SubmitStudentIIR(c.Request.Context(), userID, req)
	if err != nil {
		log.Printf("[PostStudentIIR] {Service Call}: %s", err.Error())
		response.SendError(
			c,
			"Failed to submit student IIR",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"id": iirID, "message": "Student IIR submitted successfully"},
	)
}

// GenerateIIR godoc
// @Summary      Generate Student IIR PDF
// @Description  Generates and downloads the student's Initial Interview Record
// as a PDF.
// @Tags         Students
// @Produce      application/pdf
// @Param        iirID   path      string  true  "IIR ID"
// @Success      200     {file}    binary
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /students/{iirID}/iir/download [get]
func (h *Handler) GetStudentIIRPDF(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	pdfBytes, fileName, err := h.service.GenerateIIR(c.Request.Context(), iirID)
	if err != nil {
		log.Printf("[GetStudentIIRPDF] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to generate IIR PDF",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	c.Header(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"%s\"", fileName),
	)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// PatchStudentBulkStatus handles PATCH /students/inventory/records/bulk-status.
// It applies a lifecycle status transition to multiple student records in a
// single transaction. For the "Graduated" status the service layer enforces
// eligibility (Diploma/Year-3, Bachelor/Year-4) and quietly skips ineligible
// records — the response always reports the attempted count so the frontend
// can surface a warning if the actual count differs.
func (h *Handler) PatchStudentStatusBulk(c *gin.Context) {
	var req BulkUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[PatchStudentStatusBulk] {Bind JSON}: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.service.BulkUpdateStudentStatus(c.Request.Context(), req); err != nil {
		log.Printf("[PatchStudentStatusBulk] {Service Call}: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update student status pool",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bulk status update successful"})
}

func (h *Handler) PostStudentCOR(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	file, err := c.FormFile("cor")
	if err != nil {
		response.SendFail(c, gin.H{"error": "COR file is required (field: cor)"})
		return
	}

	fileID, err := h.service.SubmitCOR(c.Request.Context(), userID, file)
	if err != nil {
		log.Printf("[PostStudentCOR] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to submit COR",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"fileId":  fileID,
		"message": "COR submitted and processing successfully",
	})
}

func (h *Handler) PostStudentCORByIIRID(c *gin.Context) {
	iirID := c.Param("iirID")
	if iirID == "" {
		response.SendFail(c, gin.H{"error": "Invalid IIR ID format"})
		return
	}

	// Resolve UserID from IIRID
	iir, err := h.service.GetStudentIIR(c.Request.Context(), iirID)
	if err != nil || iir == nil {
		log.Printf("[PostStudentCORByIIRID] {Resolve UserID}: %v", err)
		response.SendError(
			c,
			"Student record not found",
			http.StatusNotFound,
			nil,
		)
		return
	}

	file, err := c.FormFile("cor")
	if err != nil {
		response.SendFail(c, gin.H{"error": "COR file is required (field: cor)"})
		return
	}

	fileID, err := h.service.SubmitCOR(c.Request.Context(), iir.UserID, file)
	if err != nil {
		log.Printf("[PostStudentCORByIIRID] {Service Call}: %v", err)
		response.SendError(
			c,
			"Failed to submit COR",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"fileId":  fileID,
		"message": "COR submitted and processing successfully",
	})
}

func (h *Handler) GetStudentCORByUserID(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		response.SendFail(c, gin.H{"error": "Invalid User ID format"})
		return
	}

	cor, err := h.service.GetStudentCOR(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetStudentCORByUserID] {Fetch COR}: %v", err)
		response.SendError(
			c,
			"Failed to retrieve COR",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, cor)
}

func (h *Handler) GetStudentCORs(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	cors, err := h.service.GetStudentCORs(c.Request.Context(), userID)
	if err != nil {
		log.Printf("[GetStudentCORs] {Fetch CORs}: %v", err)
		response.SendError(
			c,
			"Failed to retrieve CORs",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, cors)
}
