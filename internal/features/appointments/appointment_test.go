package appointments

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/core/testutil"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
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

func TestAppointmentLifecycle(t *testing.T) {
	// Setup mock AI classifier server
	aiServer := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]any{
			"level":      "CRITICAL",
			"confidence": 0.99,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer aiServer.Close()

	db := testutil.SetupTestDB(t)
	defer db.Close()

	// Seed geographic data, test user, role, and IIR record
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
	`)
	if err != nil {
		t.Fatalf("failed to seed test database: %v", err)
	}

	// Setup service dependencies
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo, nil, nil)

	noteRepo := notes.NewRepository(db)
	noteService := notes.NewService(
		noteRepo,
		&mockLogger{},
		&mockNotifier{},
		&mockEmailer{},
	)

	apptRepo := NewRepository(db)
	cfg := &config.Config{
		AIBaseUrl: aiServer.URL,
		AIAPIKey:  "test-key",
	}

	svc := NewService(
		apptRepo,
		&mockNotifier{},
		&mockLogger{},
		&mockEmailer{},
		userService,
		noteService,
		nil,
		cfg,
	)

	// Create context with metadata for auditing
	ctx := audit.WithContext(
		context.Background(),
		"127.0.0.1",
		"test-agent",
		"student-uuid",
		"student@test.com",
		"Student",
		"trace-id",
	)

	// 1. Create Appointment
	req := AppointmentDTO{
		Reason: structs.StringToNullableString(
			"I am feeling depressed and helpless",
		),
		WhenDate:            "2025-06-30",
		TimeSlot:            TimeSlot{ID: 1},
		AppointmentCategory: AppointmentCategory{ID: 1},
	}

	created, err := svc.CreateAppointment(
		ctx,
		"iir-uuid",
		req,
		cfg,
	)
	if err != nil {
		t.Fatalf("CreateAppointment failed: %v", err)
	}
	if created == nil {
		t.Fatal("expected created appointment to be non-nil")
	}

	// Verify AI urgency classification was triggered and applied
	if created.UrgencyLevel != "CRITICAL" {
		t.Errorf("expected level CRITICAL, got %s", created.UrgencyLevel)
	}
	if created.UrgencyScore != 0.99 {
		t.Errorf("expected score 0.99, got %f", created.UrgencyScore)
	}

	// 2. Update Status (e.g. Schedule it)
	req.Status = AppointmentStatus{ID: 2} // Scheduled status
	req.AdminNotes = structs.StringToNullableString("Scheduled successfully")
	req.WhenDate = "2025-06-30"
	req.TimeSlot = TimeSlot{ID: 1}
	req.AppointmentCategory = AppointmentCategory{ID: 1}

	err = svc.UpdateAppointment(ctx, created.ID, req)
	if err != nil {
		t.Fatalf("UpdateAppointment failed: %v", err)
	}

	// 3. Verify Update in database
	updated, err := apptRepo.GetAppointment(ctx, db, created.ID)
	if err != nil {
		t.Fatalf("failed to retrieve appointment: %v", err)
	}
	if updated.StatusID != 2 {
		t.Errorf("expected status ID 2, got %d", updated.StatusID)
	}
}
