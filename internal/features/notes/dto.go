package notes

import "time"

type SignificantNoteDTO struct {
	ID              string    `json:"id,omitempty"`
	AppointmentID   string    `json:"appointmentId,omitempty"`
	AdmissionSlipID string    `json:"admissionSlipId,omitempty"`
	Note            string    `json:"note"                    binding:"required"`
	Remarks         string    `json:"remarks"                 binding:"required"`
	CreatedAt       time.Time `json:"createdAt,omitempty"                        db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"                        db:"updated_at"`
}
