package integrations

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PostLinkStudent(c *gin.Context) {}

func (h *Handler) PostEmailVerificationCode(c *gin.Context) {}

// HandleListStudents godoc
// @Summary List students
// @Description Get a paginated list of students with optional filters
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param search query string false "Search term for student name or number"
// @Param course_id query int false "Filter by course ID"
// @Param gender_id query int false "Filter by gender ID"
// @Param year_level query int false "Filter by year level"
// @Param order_by query string false "Order by field (first_name, last_name, student_number, created_at, updated_at, year_level, course_id)"
// @Param page query int false "Page number for pagination"
// @Param page_size query int false "Number of items per page for pagination"
// @Success 200 {object} ListStudentsSuccessResponse
// @Failure 400 {object} response.CommonErrorResponse "Bad Request"
// @Failure 401 {object} response.CommonErrorResponse "Unauthorized"
// @Failure 500 {object} response.CommonErrorResponse "Internal Server Error"
// @Router /integrations/students/profiles [get]
func (h *Handler) GetStudents(c *gin.Context) {
	var req OGOSListStudentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[GetStudents] {Bind Query}: %v", err)
		response.SendFail(c, gin.H{"error": "invalid query parameters"})
		return
	}

	resp, err := h.service.ListStudents(c.Request.Context(), req)
	if err != nil {
		log.Printf("[GetStudents] {Service List}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetStudentByStudentNumber godoc
// @Summary Get student by student number
// @Description Get student information by student number
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentNumber path string true "Student number of the student"
// @Success 200 {object} StudentSuccessResponse
// @Failure 400 {object} response.CommonErrorResponse "Bad Request"
// @Failure 401 {object} response.CommonErrorResponse "Unauthorized"
// @Failure 404 {object} response.CommonErrorResponse "Not Found"
// @Failure 500 {object} response.CommonErrorResponse "Internal Server Error"
// @Router /integrations/students/{studentNumber} [get]
func (h *Handler) GetStudentByStudentNumber(c *gin.Context) {
	studentNumber := c.Param("studentNumber")
	if studentNumber == "" {
		response.SendFail(
			c,
			gin.H{"error": "studentNumber parameter is required"},
		)
		return
	}

	student, err := h.service.GetStudentByStudentNumber(
		c.Request.Context(),
		studentNumber,
	)
	if err != nil {
		log.Printf("[GetStudentByStudentNumber] {Service Get}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if student == nil {
		response.SendFail(
			c,
			gin.H{"error": constants.ErrNotFound},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, student)
}

// GetStudentByEmail godoc
// @Summary Get student by Email
// @Description Get student information by their email
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param email query string true "Student Email"
// @Success 200 {object} StudentSuccessResponse
// @Failure 401 {object} response.CommonErrorResponse "Unauthorized"
// @Failure 403 {object} response.CommonErrorResponse "Forbidden"
// @Failure 404 {object} response.CommonErrorResponse "Not Found"
// @Failure 500 {object} response.CommonErrorResponse "Internal Server Error"
// @Router /integrations/students/profile [get]
func (h *Handler) GetStudentByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		response.SendFail(
			c,
			gin.H{"error": "email query parameter is required"},
		)
		return
	}

	// For M2M calls, ensure the client is formally verified (Partner status)
	isM2M, _ := c.Get("isM2M")
	isVerified, _ := c.Get("isVerified")
	if m2m, _ := isM2M.(bool); m2m {
		if verified, _ := isVerified.(bool); !verified {
			response.SendFail(
				c,
				gin.H{
					"error": "Only formally verified Partner Systems can query by email",
				},
				http.StatusForbidden,
			)
			return
		}
	}

	student, err := h.service.GetStudentByEmail(
		c.Request.Context(),
		email,
	)
	if err != nil {
		log.Printf("[GetStudentByEmail] {Service Get}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if student == nil {
		response.SendFail(
			c,
			gin.H{"error": constants.ErrNotFound},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, student)
}

// GetPersonalInfoByStudentNumber godoc
// @Summary Get personal information by student number
// @Description Get personal information of a student by their student number
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentNumber path string true "Student number of the student"
// @Success 200 {object} StudentPersonalInfoSuccessResponse
// @Failure 400 {object} response.CommonErrorResponse "Bad Request"
// @Failure 401 {object} response.CommonErrorResponse "Unauthorized"
// @Failure 404 {object} response.CommonErrorResponse "Not Found"
// @Failure 500 {object} response.CommonErrorResponse "Internal Server Error"
// @Router /integrations/students/{studentNumber}/personal-info [get]
func (h *Handler) GetPersonalInfoByStudentNumber(c *gin.Context) {
	studentNumber := c.Param("studentNumber")
	if studentNumber == "" {
		response.SendFail(
			c,
			gin.H{"error": "student_number parameter is required"},
		)
		return
	}

	student, err := h.service.GetPersonalInfoByStudentNumber(
		c.Request.Context(),
		studentNumber,
	)
	if err != nil {
		log.Printf("[GetPersonalInfoByStudentNumber] {Service Get}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if student == nil {
		response.SendFail(
			c,
			gin.H{"error": constants.ErrNotFound},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, student)
}

// GetAddressByStudentNumber godoc
// @Summary Get student addresses by student number
// @Description Get all addresses of a student by their student number
// @Tags External Students
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param studentNumber path string true "Student number of the student"
// @Success 200 {object} StudentAddressSuccessResponse
// @Failure 400 {object} response.CommonErrorResponse "Bad Request"
// @Failure 401 {object} response.CommonErrorResponse "Unauthorized"
// @Failure 404 {object} response.CommonErrorResponse "Not Found"
// @Failure 500 {object} response.CommonErrorResponse "Internal Server Error"
// @Router /integrations/students/{studentNumber}/addresses [get]
func (h *Handler) GetAddressByStudentNumber(c *gin.Context) {
	studentNumber := c.Param("studentNumber")
	if studentNumber == "" {
		response.SendFail(
			c,
			gin.H{"error": "student_number parameter is required"},
		)
		return
	}

	studentAddresses, err := h.service.GetAddressByStudentNumber(
		c.Request.Context(),
		studentNumber,
	)
	if err != nil {
		log.Printf("[GetAddressByStudentNumber] {Service Get}: %v", err)
		response.SendError(
			c,
			string(constants.ErrInternalServerError),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if len(studentAddresses) == 0 {
		response.SendFail(
			c,
			gin.H{"error": constants.ErrNotFound},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, studentAddresses)
}
