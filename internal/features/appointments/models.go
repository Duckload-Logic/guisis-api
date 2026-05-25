package appointments

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// TimeSlot represents a predefined appointment window.
type TimeSlot struct {
	ID   int    `db:"id"   json:"id"`
	Time string `db:"time" json:"time"`
}

// AppointmentStatus represents the lifecycle state of an appointment.
type AppointmentStatus struct {
	ID   int    `db:"id"        json:"id"`
	Name string `db:"name"      json:"name"`
}

// StatusCount represents a summary of appointments per status.
type StatusCount struct {
	ID    int    `db:"id"    json:"id"`
	Name  string `db:"name"  json:"name"`
	Count int    `db:"count" json:"count"`
}

// AppointmentCategory represents the nature of the consultation concern.
type AppointmentCategory struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

// AvailableTimeSlotView represents a slot's availability on a specific date.
type AvailableTimeSlotView struct {
	ID          int    `db:"time_slot_id" json:"id"`
	Time        string `db:"time"         json:"time"`
	IsAvailable bool   `db:"is_available" json:"isAvailable"`
}

// Appointment represents the core entity for a scheduled consultation.
type Appointment struct {
	ID                    string                 `db:"id"                      json:"id"`
	IIRID                 string                 `db:"iir_id"                  json:"iirId"`
	Reason                structs.NullableString `db:"reason"                  json:"reason"`
	AdminNotes            structs.NullableString `db:"admin_notes"            json:"adminNotes"`
	WhenDate              string                 `db:"when_date"               json:"whenDate"`
	TimeSlotID            int                    `db:"time_slot_id"            json:"timeSlotId"`
	AppointmentCategoryID int                    `db:"appointment_category_id" json:"appointmentCategoryId"`
	StatusID              int                    `db:"status_id"               json:"statusId"`
	UrgencyLevel          string                 `db:"urgency_level"           json:"urgencyLevel"`
	UrgencyScore          float64                `db:"urgency_score"           json:"urgencyScore"`
	CreatedAt             time.Time              `db:"created_at"              json:"createdAt"`
	UpdatedAt             time.Time              `db:"updated_at"              json:"updatedAt"`
}

// AppointmentWithDetailsView represents a denormalized view for listings.
type AppointmentWithDetailsView struct {
	ID             string                 `db:"id"               json:"id"`
	UserID         string                 `db:"user_id"          json:"userId"`
	IIRID          string                 `db:"iir_id"           json:"iirId"`
	StudentNumber  string                 `db:"student_number"   json:"studentNumber"`
	UserFirstName  string                 `db:"user_first_name"  json:"userFirstName"`
	UserMiddleName structs.NullableString `db:"user_middle_name" json:"userMiddleName,omitempty"`
	UserLastName   string                 `db:"user_last_name"   json:"userLastName"`
	UserEmail      string                 `db:"user_email"       json:"userEmail"`
	Reason         structs.NullableString `db:"reason"           json:"reason"`
	AdminNotes     structs.NullableString `db:"admin_notes"       json:"adminNotes"`
	WhenDate       string                 `db:"when_date"         json:"whenDate"`
	CreatedAt      time.Time              `db:"created_at"        json:"createdAt"`
	UpdatedAt      time.Time              `db:"updated_at"        json:"updatedAt"`
	TimeSlotID     int                    `db:"time_slot_id"     json:"timeSlotId"`
	TimeSlotTime   string                 `db:"time_slot_time"   json:"timeSlotTime"`
	CategoryID     int                    `db:"category_id"      json:"categoryId"`
	CategoryName   string                 `db:"category_name"    json:"categoryName"`
	StatusID       int                    `db:"status_id"        json:"statusId"`
	StatusName     string                 `db:"status_name"      json:"statusName"`
	UrgencyLevel   string                 `db:"urgency_level"     json:"urgencyLevel"`
	UrgencyScore   float64                `db:"urgency_score"     json:"urgencyScore"`
}

// DailyStatusCount represents aggregated daily metrics.
type DailyStatusCount struct {
	Date             string `db:"date"              json:"date"`
	PendingCount     int    `db:"pending_count"     json:"pendingCount"`
	ScheduledCount   int    `db:"scheduled_count"   json:"scheduledCount"`
	RescheduledCount int    `db:"rescheduled_count" json:"rescheduledCount"`
}
