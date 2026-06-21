package slips

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/testutil"
	"github.com/olazo-johnalbert/duckload-api/internal/features/files"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/ocr"
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

type mockStorage struct{}

func (m *mockStorage) Upload(
	ctx context.Context,
	path string,
	reader io.ReadSeeker,
	contentType string,
) error {
	return nil
}

func (m *mockStorage) Download(
	ctx context.Context,
	path string,
	writer io.Writer,
) error {
	return nil
}

func (m *mockStorage) Delete(ctx context.Context, path string) error {
	return nil
}

func createMockPDFFileHeader(
	filename string,
) (*multipart.FileHeader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	// Write standard PDF magic header bytes
	_, err = part.Write([]byte("%PDF-1.4\n" + strings.Repeat("A", 100)))
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	err = req.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}

	files := req.MultipartForm.File["file"]
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found")
	}
	return files[0], nil
}

func TestSlipLifecycle(t *testing.T) {
	// Setup OCR mock server
	ocrServer := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]any{
			"is_valid": true,
			"message":  "Document is valid",
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ocrServer.Close()

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

	// Retrieve valid lookup IDs from seeded database
	var genderID int
	err = db.Get(&genderID, "SELECT id FROM genders LIMIT 1")
	if err != nil {
		t.Fatalf("failed to get gender ID: %v", err)
	}

	var civilStatusID int
	err = db.Get(&civilStatusID, "SELECT id FROM civil_status_types LIMIT 1")
	if err != nil {
		t.Fatalf("failed to get civil status ID: %v", err)
	}

	var religionID int
	err = db.Get(&religionID, "SELECT id FROM religions LIMIT 1")
	if err != nil {
		t.Fatalf("failed to get religion ID: %v", err)
	}

	var courseID int
	err = db.Get(&courseID, "SELECT id FROM courses LIMIT 1")
	if err != nil {
		t.Fatalf("failed to get course ID: %v", err)
	}

	var categoryID int
	err = db.Get(
		&categoryID,
		"SELECT id FROM admission_slip_categories LIMIT 1",
	)
	if err != nil {
		t.Fatalf("failed to get slip category ID: %v", err)
	}

	// Seed student_personal_info to satisfy views and queries
	_, err = db.Exec(`
		INSERT INTO student_personal_info (
			iir_id, student_number, gender_id, civil_status_id,
			religion_id, height_m, weight_kg, complexion,
			high_school_gwa, course_id, year_level, section,
			place_of_birth, date_of_birth, mobile_number, status_id
		) VALUES (
			'iir-uuid', '2025-00001-MN-0', ?, ?,
			?, 1.75, 70.0, 'Fair',
			85.0, ?, 1, 1,
			'Manila', '2000-01-01', '09123456789', 1
		)
	`, genderID, civilStatusID, religionID, courseID)
	if err != nil {
		t.Fatalf("failed to seed student_personal_info: %v", err)
	}

	// Setup service dependencies
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo, nil, nil)

	fileRepo := files.NewRepository(db)
	cfg := &config.Config{
		AIBaseUrl: ocrServer.URL,
		AIAPIKey:  "test-key",
	}

	ocrClient := ocr.NewClient(ocrServer.URL, "test-key")
	storageClient := &mockStorage{}

	filesService := files.NewService(
		fileRepo,
		storageClient,
		ocrClient,
		cfg,
	)

	slipRepo := NewRepository(db)
	svc := NewService(
		slipRepo,
		&mockLogger{},
		&mockNotifier{},
		&mockEmailer{},
		storageClient,
		userService,
		nil,
		filesService,
		ocrClient,
		cfg,
	)

	// Create mock files
	excuseFile, err := createMockPDFFileHeader("excuse.pdf")
	if err != nil {
		t.Fatalf("failed to create excuse file header: %v", err)
	}

	parentIDFile, err := createMockPDFFileHeader("parent_id.pdf")
	if err != nil {
		t.Fatalf("failed to create parent ID file header: %v", err)
	}

	// 1. Submit Excuse Slip
	// We set Dates to be safe (not in future for absence, not in past for needed)
	todayStr := time.Now().Format("2006-01-02")
	req := CreateSlipRequest{
		Reason:        "Under the weather",
		DateOfAbsence: todayStr,
		DateNeeded:    todayStr,
		CategoryID:    categoryID,
	}

	// Setup context with user metadata for audit
	ctx := audit.WithContext(
		context.Background(),
		"127.0.0.1",
		"test-agent",
		"student-uuid",
		"student@test.com",
		"Student",
		"trace-id",
	)

	slipDTO, err := svc.SubmitExcuseSlip(
		ctx,
		"iir-uuid",
		req,
		[]*multipart.FileHeader{excuseFile},
		[]*multipart.FileHeader{parentIDFile},
	)
	if err != nil {
		t.Fatalf("SubmitExcuseSlip failed: %v", err)
	}
	if slipDTO == nil {
		t.Fatal("expected non-nil slipDTO")
	}

	// Verify details
	if slipDTO.Reason != "Under the weather" {
		t.Errorf("expected reason 'Under the weather', got %s", slipDTO.Reason)
	}

	// 2. Update Excuse Slip Status (e.g. Approve)
	err = svc.UpdateExcuseSlipStatus(
		ctx,
		slipDTO.ID,
		"Approved",
		"Approved after verification",
	)
	if err != nil {
		t.Fatalf("UpdateExcuseSlipStatus failed: %v", err)
	}

	// 3. Verify status in database
	updated, err := slipRepo.GetSlipByIDWithDetails(ctx, db, slipDTO.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated slip: %v", err)
	}
	if updated.StatusName != "Approved" {
		t.Errorf("expected StatusName 'Approved', got %s", updated.StatusName)
	}
}
