package notes

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/testutil"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type mockLogger struct{}

func (m *mockLogger) Record(
	ctx context.Context,
	tx datastore.DB,
	entry audit.LogEntry,
) {
}

type mockNotifier struct{}

func (m *mockNotifier) Send(
	ctx context.Context,
	notif audit.NotificationEntry,
) error {
	return nil
}

type mockEmailer struct{}

func (m *mockEmailer) Send(
	ctx context.Context,
	email audit.EmailEntry,
) error {
	return nil
}

func TestSignificantNotesLifecycle(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	// Seed required data for the student
	_, err := db.Exec(`
		INSERT INTO regions (id, name, code) VALUES (1, 'NCR', '130000000');
		INSERT INTO provinces (id, code, name, region_code) 
		VALUES (1, '133900000', 'NCR', '130000000');
		INSERT INTO cities (id, code, name, province_code, region_code) 
		VALUES (1, '133901000', 'Manila', '133900000', '130000000');
		INSERT INTO barangays (id, code, name, city_code) 
		VALUES (1, '133901001', 'Barangay 1', '133901000');

		INSERT INTO users (id, email, first_name, last_name, auth_type, is_active)
		VALUES ('student-uuid', 'student@test.com', 'Test', 'Student', 'native', 1);

		INSERT INTO user_roles (user_id, role_id) VALUES ('student-uuid', 1);

		INSERT INTO iir_records (id, user_id, is_submitted)
		VALUES ('iir-uuid', 'student-uuid', 1);

		INSERT INTO appointments (
			id, iir_id, time_slot_id, when_date,
			appointment_category_id, urgency_score, status_id
		) VALUES (
			'appt-uuid-placeholder', 'iir-uuid', 1, '2025-06-30',
			1, 0.5, 1
		);

		INSERT INTO admission_slips (
			id, iir_id, category_id, reason,
			date_of_absence, date_needed, status_id
		) VALUES (
			'slip-uuid-placeholder', 'iir-uuid', 1, 'Fever',
			'2025-06-30', '2025-06-30', 1
		);
	`)
	if err != nil {
		t.Fatalf("failed to seed test database: %v", err)
	}

	repo := NewRepository(db)
	svc := NewService(
		repo,
		&mockLogger{},
		&mockNotifier{},
		&mockEmailer{},
	)

	ctx := context.Background()

	// 1. Create a Significant Note
	noteReq := SignificantNoteDTO{
		AppointmentID:   "appt-uuid-placeholder",
		AdmissionSlipID: "slip-uuid-placeholder",
		Note:            "Student was cooperative during the session.",
		Remarks:         "Follow-up scheduled.",
	}

	err = svc.CreateSignificantNote(ctx, "iir-uuid", noteReq)
	if err != nil {
		t.Fatalf("failed to create significant note: %v", err)
	}

	// 2. Verify HasNoteForAppointment
	hasNote, err := svc.HasNoteForAppointment(ctx, "appt-uuid-placeholder")
	if err != nil {
		t.Fatalf("failed to check note for appointment: %v", err)
	}
	if !hasNote {
		t.Error("expected note to exist for appointment ID")
	}

	// 3. Get Student Significant Notes
	notes, err := svc.GetStudentSignificantNotes(ctx, "iir-uuid")
	if err != nil {
		t.Fatalf("failed to get student significant notes: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}

	targetNote := notes[0]
	if targetNote.Note != noteReq.Note {
		t.Errorf("expected note %q, got %q", noteReq.Note, targetNote.Note)
	}
	if targetNote.Remarks != noteReq.Remarks {
		t.Errorf(
			"expected remarks %q, got %q",
			noteReq.Remarks,
			targetNote.Remarks,
		)
	}

	// 4. Delete Significant Note
	deleted, err := svc.DeleteSignificantNote(ctx, targetNote.ID)
	if err != nil {
		t.Fatalf("failed to delete significant note: %v", err)
	}
	if !deleted {
		t.Error("expected deleted status to be true")
	}

	// 5. Verify notes are empty now
	notesPostDelete, err := svc.GetStudentSignificantNotes(ctx, "iir-uuid")
	if err != nil {
		t.Fatalf("failed to get notes post-delete: %v", err)
	}
	if len(notesPostDelete) != 0 {
		t.Errorf("expected 0 notes, got %d", len(notesPostDelete))
	}
}
