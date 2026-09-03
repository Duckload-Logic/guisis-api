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
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
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
	ID              string                 `db:"id" json:"id"`
	IIRID           string                 `db:"iir_id" json:"iirId"`
	Reason          structs.NullableString `db:"reason" json:"reason"`
	AdminNotes      structs.NullableString `db:"admin_notes" json:"adminNotes"`
	WhenDate        string                 `db:"when_date" json:"whenDate"`
	TimeSlotID      int                    `db:"time_slot_id" json:"timeSlotId"`
	CategoryID      int                    `db:"appointment_category_id" json:"appointmentCategoryId"`
	StatusID        int                    `db:"status_id" json:"statusId"`
	UrgencyLevel    string                 `db:"urgency_level" json:"urgencyLevel"`
	UrgencyScore    float64                `db:"urgency_score" json:"urgencyScore"`
	PrefDate1       structs.NullableString `db:"preferred_date_1"`
	PrefTimeSlotID1 structs.NullableInt64  `db:"preferred_time_slot_id_1"`
	PrefDate2       structs.NullableString `db:"preferred_date_2"`
	PrefTimeSlotID2 structs.NullableInt64  `db:"preferred_time_slot_id_2"`
	PrefDate3       structs.NullableString `db:"preferred_date_3"`
	PrefTimeSlotID3 structs.NullableInt64  `db:"preferred_time_slot_id_3"`
	StartedAt       structs.NullableTime   `db:"started_at" json:"startedAt"`
	CompletedAt     structs.NullableTime   `db:"completed_at" json:"completedAt"`
	CreatedAt       time.Time              `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time              `db:"updated_at" json:"updatedAt"`
}

// AppointmentWithDetailsView represents a denormalized view for listings.
type AppointmentWithDetailsView struct {
	ID                 string                 `db:"id"`
	UserID             string                 `db:"user_id"`
	IIRID              string                 `db:"iir_id"`
	StudentNumber      string                 `db:"student_number"`
	UserFirstName      string                 `db:"user_first_name"`
	UserMiddleName     structs.NullableString `db:"user_middle_name"`
	UserLastName       string                 `db:"user_last_name"`
	UserEmail          string                 `db:"user_email"`
	UserProfilePicture structs.NullableString `db:"user_profile_picture"`
	Reason             structs.NullableString `db:"reason"`
	AdminNotes         structs.NullableString `db:"admin_notes"`
	WhenDate           string                 `db:"when_date"`
	StartedAt          structs.NullableTime   `db:"started_at"`
	CompletedAt        structs.NullableTime   `db:"completed_at"`
	CreatedAt          time.Time              `db:"created_at"`
	UpdatedAt          time.Time              `db:"updated_at"`
	TimeSlotID         int                    `db:"time_slot_id"`
	TimeSlotTime       string                 `db:"time_slot_time"`
	CategoryID         int                    `db:"category_id"`
	CategoryName       string                 `db:"category_name"`
	StatusID           int                    `db:"status_id"`
	StatusName         string                 `db:"status_name"`
	UrgencyLevel       string                 `db:"urgency_level"`
	UrgencyScore       float64                `db:"urgency_score"`
	PrefDate1          structs.NullableString `db:"preferred_date_1"`
	PrefTimeSlotID1    structs.NullableInt64  `db:"preferred_time_slot_id_1"`
	PrefTimeSlotTime1  structs.NullableString `db:"preferred_time_slot_time_1"`
	PrefDate2          structs.NullableString `db:"preferred_date_2"`
	PrefTimeSlotID2    structs.NullableInt64  `db:"preferred_time_slot_id_2"`
	PrefTimeSlotTime2  structs.NullableString `db:"preferred_time_slot_time_2"`
	PrefDate3          structs.NullableString `db:"preferred_date_3"`
	PrefTimeSlotID3    structs.NullableInt64  `db:"preferred_time_slot_id_3"`
	PrefTimeSlotTime3  structs.NullableString `db:"preferred_time_slot_time_3"`
}

// DailyStatusCount represents aggregated daily metrics.
type DailyStatusCount struct {
	Date             string `db:"date"              json:"date"`
	PendingCount     int    `db:"pending_count"     json:"pendingCount"`
	ScheduledCount   int    `db:"scheduled_count"   json:"scheduledCount"`
	RescheduledCount int    `db:"rescheduled_count" json:"rescheduledCount"`
}

// FullName returns the user's full name from the denormalized view.
func (a *AppointmentWithDetailsView) FullName() string {
	if a == nil {
		return "A student"
	}
	return a.UserFirstName + " " + a.UserLastName
}
