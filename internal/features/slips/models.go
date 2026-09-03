package slips

import (
	"database/sql"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// SlipCategory represents the nature of the absence or request.
type SlipCategory struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

// SlipStatus represents the lifecycle state of a slip request.
type SlipStatus struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

// SlipWithDetailsView represents a denormalized view of a slip for listings.
type SlipWithDetailsView struct {
	ID             string                 `db:"id"               json:"id"`
	IIRID          string                 `db:"iir_id"           json:"iirId"`
	UserID         string                 `db:"user_id"          json:"userId"`
	UserFirstName  string                 `db:"user_first_name"  json:"userFirstName"`
	UserMiddleName structs.NullableString `db:"user_middle_name" json:"userMiddleName,omitempty"`
	UserLastName   string                 `db:"user_last_name"   json:"userLastName"`
	UserEmail      string                 `db:"user_email"       json:"userEmail"`
	StudentNumber  string                 `db:"student_number"   json:"studentNumber"`
	Reason         string                 `db:"reason"           json:"reason"`
	DateOfAbsence  string                 `db:"date_of_absence"  json:"dateOfAbsence"`
	DateNeeded     string                 `db:"date_needed"      json:"dateNeeded"`
	AdminNotes     structs.NullableString `db:"admin_notes"      json:"adminNotes"`
	CategoryID     int                    `db:"category_id"      json:"categoryId"`
	CategoryName   string                 `db:"category_name"    json:"categoryName"`
	StatusID       int                    `db:"status_id"        json:"statusId"`
	StatusName     string                 `db:"status_name"      json:"statusName"`
	UrgencyScore   int                    `db:"urgency_score"    json:"urgencyScore"`
	TicketCode     structs.NullableString `db:"ticket_code"`
	IsVerified     sql.NullBool           `db:"is_verified"`
	VerifiedAt     structs.NullableTime   `db:"verified_at"`
	StartedAt      structs.NullableTime   `db:"started_at"       json:"startedAt"`
	CompletedAt    structs.NullableTime   `db:"completed_at"     json:"completedAt"`
	CreatedAt      time.Time              `db:"created_at"`
	UpdatedAt      time.Time              `db:"updated_at"`
}

// SlipStatusCount represents aggregated status metrics.
type SlipStatusCount struct {
	ID    int    `db:"id"    json:"id"`
	Name  string `db:"name"  json:"name"`
	Count int    `db:"count" json:"count"`
}

// Slip represents the core entity for an excuse or admission slip.
type Slip struct {
	ID            string                 `db:"id"              json:"id"`
	IIRID         string                 `db:"iir_id"          json:"iirId"`
	Reason        string                 `db:"reason"          json:"reason"`
	DateOfAbsence string                 `db:"date_of_absence" json:"dateOfAbsence"`
	DateNeeded    string                 `db:"date_needed"     json:"dateNeeded"`
	AdminNotes    structs.NullableString `db:"admin_notes"     json:"adminNotes"`
	CategoryID    int                    `db:"category_id"     json:"categoryId"`
	StatusID      int                    `db:"status_id"       json:"statusId"`
	StartedAt     structs.NullableTime   `db:"started_at"      json:"startedAt"`
	CompletedAt   structs.NullableTime   `db:"completed_at"    json:"completedAt"`
	CreatedAt     time.Time              `db:"created_at"      json:"createdAt"`
	UpdatedAt     time.Time              `db:"updated_at"      json:"updatedAt"`
}

// SlipAttachment represents a file linked to a specific slip.
type SlipAttachment struct {
	FileID         string                 `db:"file_id"           json:"fileId"`
	SlipID         structs.NullableString `db:"admission_slip_id" json:"slipId"`
	AttachmentType string                 `db:"attachment_type"   json:"attachmentType"`
	FileName       string                 `db:"file_name"         json:"fileName"`
	FileURL        string                 `db:"file_url"          json:"fileUrl"`
	FileType       string                 `db:"file_type"         json:"fileType"`
	FileSize       int64                  `db:"file_size"         json:"fileSize"`
	MimeType       string                 `db:"mime_type"         json:"mimeType"`
}

// AdmissionTicket represents the physical verification ticket for on-site approval.
type AdmissionTicket struct {
	ID              string                 `db:"id"                json:"id"`
	AdmissionSlipID string                 `db:"admission_slip_id" json:"admissionSlipId"`
	TicketCode      string                 `db:"ticket_code"       json:"ticketCode"`
	IsVerified      bool                   `db:"is_verified"       json:"isVerified"`
	VerifiedAt      structs.NullableTime   `db:"verified_at"       json:"verifiedAt"`
	VerifiedBy      structs.NullableString `db:"verified_by"       json:"verifiedBy"`
	CreatedAt       time.Time              `db:"created_at"        json:"createdAt"`
	UpdatedAt       time.Time              `db:"updated_at"        json:"updatedAt"`
}
