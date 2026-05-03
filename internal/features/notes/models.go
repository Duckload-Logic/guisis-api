package notes

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// SignificantNote represents a professional note, used for both business logic and data persistence.
type SignificantNote struct {
	ID              string                 `db:"id"                json:"id"`
	IIRID           structs.NullableString `db:"iir_id"            json:"iirId"`
	AppointmentID   structs.NullableString `db:"appointment_id"    json:"appointmentId"`
	AdmissionSlipID structs.NullableString `db:"admission_slip_id" json:"admissionSlipId"`
	Note            string                 `db:"note"              json:"note"`
	Remarks         string                 `db:"remarks"           json:"remarks"`
	CreatedAt       time.Time              `db:"created_at"        json:"createdAt"`
	UpdatedAt       time.Time              `db:"updated_at"        json:"updatedAt"`
}
