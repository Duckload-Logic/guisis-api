package appointments

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Handler struct {
	service *Service
	cfg     *config.Config
}

// NewHandler creates a new appointments handler.
func NewHandler(service *Service, cfg *config.Config) *Handler {
	return &Handler{service: service, cfg: cfg}
}

// getIIRIDFromContext extracts iirID from context or aborts
// with Forbidden status if not found.
func getIIRIDFromContext(c *gin.Context) (string, bool) {
	iirIDVal, exists := c.Get("iirID")
	if !exists {
		response.SendFail(c, gin.H{
			"error": "Please complete your IIR profile",
		}, http.StatusForbidden)
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

// GetAppointmentCategoryList godoc
// @Summary      Get appointment categories
// @Description  Retrieves available appointment categories.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} []AppointmentCategory
// @Failure      500  {object} map[string]string
// @Router       /appointments/lookups/categories [get]
// GetAppointmentCategoryList retrieves all appointment concern categories.
func (h *Handler) GetAppointmentCategories(c *gin.Context) {
	categories, err := h.service.GetConcernCategories(c.Request.Context())
	if err != nil {
		fmt.Printf("[GetAppointmentCategories] {Fetch Categories}: %v\n", err)
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

// GetDailyStatusCountList godoc
// @Summary      Get daily status count
// @Description  Retrieves appointment status counts by date.
// @Tags         Appointments
// @Produce      json
// @Param        start_date query string true "Start date (YYYY-MM-DD)"
// @Success      200  {object} []DailyStatusCount
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /appointments/calendar/stats [get]
func (h *Handler) GetAppointmentDailyStats(c *gin.Context) {
	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	if req.StartDate == "" {
		response.SendFail(c, gin.H{"error": "start_date parameter required"})
		return
	}

	dsc, err := h.service.GetDailyStatusCount(c.Request.Context(), req.StartDate)
	if err != nil {
		fmt.Printf("[GetAppointmentDailyStats] {Fetch Daily Stats}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve statistics",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, dsc)
}

// PostAppointment godoc
// @Summary      Book a new appointment
// @Description  Schedules a new appointment for a student.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        request body      AppointmentDTO true "Appointment Details"
// @Success      201     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /appointments [post]
func (h *Handler) PostAppointment(c *gin.Context) {
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req AppointmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	appt, err := h.service.CreateAppointment(
		c.Request.Context(), iirID, req, h.cfg,
	)
	if err != nil {
		fmt.Printf("[PostAppointment] {Create Appointment}: %v\n", err)
		response.SendError(
			c,
			"Failed to create appointment",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"message": "Appointment created successfully",
		"id":      appt.ID,
	}, http.StatusCreated)
}

// GetAppointmentByID godoc
// @Summary      Get appointment details
// @Description  Retrieves the details of a specific appointment.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Appointment ID"
// @Success      200  {object}  Appointment
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /appointments/id/{id} [get]
func (h *Handler) GetAppointmentByID(c *gin.Context) {
	id := c.Param("appointmentID")
	if id == "" {
		response.SendFail(c, gin.H{"error": "Invalid ID format"})
		return
	}

	appt, err := h.service.GetAppointmentByID(c.Request.Context(), id)
	if err != nil {
		fmt.Printf("[GetAppointmentByID] {Fetch Appointment}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve appointment",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	if appt == nil {
		response.SendFail(
			c,
			gin.H{"error": "Appointment not found"},
			http.StatusNotFound,
		)
		return
	}

	response.SendSuccess(c, appt)
}

// GetAppointmentList godoc
// @Summary      Get all appointments
// @Description  Retrieves a list of all appointments.
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        status     query string false "Filter by status"
// @Param        start_date query string false "Filter by start date"
// @Param        end_date   query string false "Filter by end date"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /appointments [get]
func (h *Handler) GetAppointments(c *gin.Context) {
	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	appts, err := h.service.ListAppointments(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetAppointments] {Fetch Appointments}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve appointments",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, appts)
}

// GetAvailableTimeSlotList godoc
// @Summary      Get available time slots
// @Description  Retrieves available time slots for a date.
// @Tags         Appointments
// @Produce      json
// @Param        date query string true "Date (YYYY-MM-DD)"
// @Success      200  {object} []AvailableTimeSlotView
// @Failure      400  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /appointments/lookups/slots [get]
func (h *Handler) GetAppointmentSlots(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		response.SendFail(c, gin.H{"error": "date parameter required"})
		return
	}

	slots, err := h.service.GetAvailableTimeSlots(c.Request.Context(), date)
	if err != nil {
		fmt.Printf("[GetAppointmentSlots] {Fetch Slots}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve time slots",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, slots)
}

// GetAppointmentStatusList godoc
// @Summary      Get appointment statuses
// @Description  Retrieves available appointment statuses.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} []AppointmentStatus
// @Failure      500  {object} map[string]string
// @Router       /appointments/lookups/statuses [get]
func (h *Handler) GetAppointmentStatuses(c *gin.Context) {
	statuses, err := h.service.GetAppointmentStatuses(c.Request.Context())
	if err != nil {
		fmt.Printf("[GetAppointmentStatuses] {Fetch Statuses}: %v\n", err)
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

// GetAppointmentListByIIR godoc
// @Summary      Get student's appointments
// @Description  Retrieves appointments for the authenticated student.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} map[string]interface{}
// @Failure      403  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /appointments/me [get]
func (h *Handler) GetAppointmentMe(c *gin.Context) {
	iirID, ok := getIIRIDFromContext(c)
	if !ok {
		return
	}

	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	appointments, err := h.service.GetAppointmentsByIIRID(
		c.Request.Context(), iirID, req,
	)
	if err != nil {
		fmt.Printf("[GetAppointmentMe] {Fetch Appointments}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve appointments",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, appointments)
}

// GetAppointmentStatsList godoc
// @Summary      Get appointment statistics
// @Description  Retrieves appointment status counts.
// @Tags         Appointments
// @Produce      json
// @Success      200  {object} []StatusCount
// @Failure      500  {object} map[string]string
// @Router       /appointments/stats [get]
func (h *Handler) GetAppointmentStats(c *gin.Context) {
	iirIDVal, exists := c.Get("iirID")
	roleIDs := c.MustGet("roleIDs").([]int)

	isStudent := false
	for _, rid := range roleIDs {
		if rid == int(constants.StudentRoleID) {
			isStudent = true
			break
		}
	}

	var iirIDPtr *string
	if isStudent {
		if !exists {
			response.SendFail(c, gin.H{
				"error": "Please complete your IIR profile",
			}, http.StatusForbidden)
			return
		}
		iirID, ok := iirIDVal.(string)
		if !ok {
			response.SendError(
				c,
				"Internal server error",
				http.StatusInternalServerError,
				nil,
			)
			return
		}
		iirIDPtr = &iirID
	}

	var req ListAppointmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}

	stats, err := h.service.GetAppointmentStats(
		c.Request.Context(), req, iirIDPtr,
	)
	if err != nil {
		fmt.Printf("[GetAppointmentStats] {Fetch Stats}: %v\n", err)
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

// PatchAppointment godoc
// @Summary      Update appointment
// @Description  Updates appointment details (reschedule).
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Param        id   path      int            true  "Appointment ID"
// @Param        body body      AppointmentDTO true  "Updated appointment"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Failure      404  {object} map[string]string
// @Failure      500  {object} map[string]string
// @Router       /appointments/id/{id} [patch]
type CancelAppointmentRequest struct {
	Reason string `json:"reason"`
}

func (h *Handler) PostAppointmentCancel(c *gin.Context) {
	id := c.Param("appointmentID")
	userID := audit.ExtractUserID(c.Request.Context())

	ownerID, err := h.service.GetUserIDByAppointmentID(c.Request.Context(), id)
	if err != nil {
		response.SendError(
			c,
			"Failed to verify ownership",
			http.StatusInternalServerError,
			nil,
		)
		return
	}
	if ownerID != userID {
		response.SendFail(c, gin.H{"error": "Access denied"}, http.StatusForbidden)
		return
	}

	appt, err := h.service.GetAppointmentByID(c.Request.Context(), id)
	if err != nil {
		response.SendError(
			c,
			"Failed to fetch appointment",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	statusName := strings.ToLower(appt.Status.Name)
	if statusName != "pending" && statusName != "scheduled" {
		response.SendFail(
			c,
			gin.H{
				"error": "Only pending or scheduled appointments can be cancelled",
			},
		)
		return
	}

	statuses, err := h.service.GetAppointmentStatuses(c.Request.Context())
	if err != nil {
		response.SendError(
			c,
			"Failed to fetch statuses",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	var cancelStatusID int
	for _, s := range statuses {
		if strings.ToLower(s.Name) == "cancelled" {
			cancelStatusID = s.ID
			break
		}
	}

	if cancelStatusID == 0 {
		response.SendError(
			c,
			"Cancelled status not found",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	var req CancelAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if err.Error() != "EOF" {
			fmt.Printf("[PostAppointmentCancel] {Bind Request}: %v\n", err)
		}
	}

	updateReq := *appt
	updateReq.Status.ID = cancelStatusID

	if req.Reason != "" {
		msg := "Student Cancellation: " + req.Reason
		if updateReq.AdminNotes.Valid && updateReq.AdminNotes.String != "" {
			updateReq.AdminNotes.String += "\n" + msg
		} else {
			updateReq.AdminNotes = structs.StringToNullableString(msg)
		}
	}

	if err := h.service.UpdateAppointment(
		c.Request.Context(), id, updateReq,
	); err != nil {
		fmt.Printf("[PostAppointmentCancel] {Update}: %v\n", err)
		response.SendError(
			c,
			"Failed to cancel appointment",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"message": "Appointment cancelled successfully",
	})
}

// PatchAppointment godoc
func (h *Handler) PatchAppointment(c *gin.Context) {
	id := c.Param("appointmentID")
	if id == "" {
		response.SendFail(c, gin.H{"error": "Invalid ID format"})
		return
	}

	var req AppointmentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": "Invalid request format"})
		return
	}

	if err := h.service.UpdateAppointment(
		c.Request.Context(), id, req,
	); err != nil {
		if err == sql.ErrNoRows {
			response.SendFail(
				c,
				gin.H{"error": "Appointment not found"},
				http.StatusNotFound,
			)
			return
		}
		fmt.Printf("[PatchAppointment] {Update Appointment}: %v\n", err)
		response.SendError(
			c,
			"Failed to update appointment",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{
		"message": "Appointment updated successfully",
	})
}
