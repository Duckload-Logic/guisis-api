package appointments

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type ListAppointmentsRequest struct {
	structs.PaginationRequest
	StatusID   string `form:"status_id"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	CategoryID string `form:"category_id"`
	Urgency    string `form:"urgency"`
	Scope      string `form:"scope"`
}

type ListAppointmentsResponse struct {
	Appointments []AppointmentDTO           `json:"appointments"`
	Meta         structs.PaginationMetadata `json:"meta"`
}

type AppointmentDTO struct {
	ID                  string                 `json:"id,omitempty"`
	UserID              string                 `json:"userId,omitempty"`
	User                users.UserResponse     `json:"user,omitempty"`
	IIRID               string                 `json:"iirId,omitempty"`
	StudentNumber       string                 `json:"studentNumber,omitempty"`
	ContactNumber       string                 `json:"contactNumber,omitempty"`
	Reason              structs.NullableString `json:"reason,omitempty"`
	WhenDate            string                 `json:"whenDate,omitempty"`
	TimeSlot            TimeSlot               `json:"timeSlot,omitempty"`
	AppointmentCategory AppointmentCategory    `json:"appointmentCategory,omitempty"`
	AdminNotes          structs.NullableString `json:"adminNotes,omitempty"`
	Status              AppointmentStatus      `json:"status,omitempty"`
	UrgencyLevel        string                 `json:"urgencyLevel,omitempty"`
	UrgencyScore        float64                `json:"urgencyScore,omitempty"`
	UrgencyRationale    string                 `json:"urgencyRationale,omitempty"`
	PreferredCounselorID *string               `json:"preferredCounselorId,omitempty"`
	CounselorPreference  string                `json:"counselorPreference,omitempty"`
	PreferredDate1      string                 `json:"preferredDate1,omitempty"`
	PreferredTimeSlot1  *TimeSlot              `json:"preferredTimeSlot1,omitempty"`
	PreferredDate2      string                 `json:"preferredDate2,omitempty"`
	PreferredTimeSlot2  *TimeSlot              `json:"preferredTimeSlot2,omitempty"`
	PreferredDate3      string                 `json:"preferredDate3,omitempty"`
	PreferredTimeSlot3  *TimeSlot              `json:"preferredTimeSlot3,omitempty"`
	HasSignificantNote  bool                   `json:"hasSignificantNote"`
	StudentCORURL       string                 `json:"studentCorUrl,omitempty"`
	StartedAt           structs.NullableTime   `json:"startedAt,omitempty"`
	CompletedAt         structs.NullableTime   `json:"completedAt,omitempty"`
	CreatedAt           time.Time              `json:"createdAt,omitempty"`
	UpdatedAt           time.Time              `json:"updatedAt,omitempty"`
}

type UpdateCounselorProfileRequest struct {
	Title    string `json:"title" binding:"required"`
	Gender   string `json:"gender" binding:"required,oneof=FEMALE MALE OTHER"`
	Position string `json:"position" binding:"required"`
}

type StartAppointmentRequest struct {
	OffsetMinutes int `json:"offsetMinutes"`
}
