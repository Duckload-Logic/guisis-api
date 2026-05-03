package appointments

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

type ListAppointmentsRequest struct {
	structs.PaginationRequest
	StatusID  string `form:"status_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
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
	Reason              structs.NullableString `json:"reason,omitempty"`
	WhenDate            string                 `json:"whenDate,omitempty"`
	TimeSlot            TimeSlot               `json:"timeSlot,omitempty"`
	AppointmentCategory AppointmentCategory    `json:"appointmentCategory,omitempty"`
	AdminNotes          structs.NullableString `json:"adminNotes,omitempty"`
	Status              AppointmentStatus      `json:"status,omitempty"`
	UrgencyLevel        string                 `json:"urgencyLevel,omitempty"`
	UrgencyScore        float64                `json:"urgencyScore,omitempty"`
	HasSignificantNote  bool                   `json:"hasSignificantNote"`
	StudentCORURL       string                 `json:"studentCorUrl,omitempty"`
	CreatedAt           time.Time              `json:"createdAt,omitempty"`
	UpdatedAt           time.Time              `json:"updatedAt,omitempty"`
}
