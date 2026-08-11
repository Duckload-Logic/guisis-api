package students

import (
	"context"
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/testutil"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
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

func TestIIRLifecycle(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer db.Close()

	// Seed geographic data, test user and roles
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
	`)
	if err != nil {
		t.Fatalf("failed to seed test database: %v", err)
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

	var programID int
	err = db.Get(&programID, "SELECT id FROM programs LIMIT 1")
	if err != nil {
		t.Fatalf("failed to get program ID: %v", err)
	}

	var relationshipID int
	err = db.Get(
		&relationshipID,
		"SELECT id FROM student_relationship_types LIMIT 1",
	)
	if err != nil {
		t.Fatalf("failed to get relationship ID: %v", err)
	}

	var parentalStatusID int
	err = db.Get(
		&parentalStatusID,
		"SELECT id FROM parental_status_types LIMIT 1",
	)
	if err != nil {
		t.Fatalf("failed to get parental status ID: %v", err)
	}

	var natureOfResidenceID int
	err = db.Get(
		&natureOfResidenceID,
		"SELECT id FROM nature_of_residence_types LIMIT 1",
	)
	if err != nil {
		t.Fatalf("failed to get nature of residence ID: %v", err)
	}

	var educationalLevelID int
	err = db.Get(
		&educationalLevelID,
		"SELECT id FROM educational_levels LIMIT 1",
	)
	if err != nil {
		t.Fatalf("failed to get educational level ID: %v", err)
	}

	var incomeRangeID int
	err = db.Get(&incomeRangeID, "SELECT id FROM income_ranges LIMIT 1")
	if err != nil {
		t.Fatalf("failed to get income range ID: %v", err)
	}

	var eduAttainmentID int
	err = db.Get(
		&eduAttainmentID,
		"SELECT id FROM educational_attainments LIMIT 1",
	)
	if err != nil {
		t.Fatalf("failed to get educational attainment ID: %v", err)
	}

	// Setup service dependencies
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo, nil, nil)

	locationsRepo := locations.NewRepository(db)
	locationsService := locations.NewService(locationsRepo)

	studentRepo := NewRepository(db)
	svc := NewService(
		studentRepo,
		locationsService,
		userService,
		nil,
		&mockLogger{},
		&mockNotifier{},
		&config.Config{},
		nil,
	)

	// Create valid ComprehensiveProfileDTO payload
	zero := 0
	addr := locations.AddressDTO{
		StreetDetail: "123 St",
		Region:       locations.Region{Code: "130000000"},
		Province:     &locations.Province{Code: "133900000"},
		City:         locations.City{Code: "133901000"},
		Barangay:     locations.Barangay{Code: "133901001"},
	}

	req := ComprehensiveProfileDTO{}
	req.Student.Gender = Gender{
		ID:   1,
		Name: "Male",
	}
	req.Student.CivilStatus.ID = civilStatusID
	req.Student.Religion.ID = religionID
	req.Student.Program.ID = programID
	req.Student.StudentNumber = "2025-00001-MN-0"
	req.Student.DateOfBirth = "2000-01-01"
	req.Student.HeightM = 1.75
	req.Student.WeightKg = 70.0
	req.Student.Complexion = "Fair"
	req.Student.HighSchoolGWA = 85.0
	req.Student.YearLevel = 1
	req.Student.Section = 1
	req.Student.PlaceOfBirth = "Manila"
	req.Student.MobileNumber = "09123456789"
	req.Student.EmergencyContact.FirstName = "John"
	req.Student.EmergencyContact.LastName = "Doe"
	req.Student.EmergencyContact.ContactNumber = "09123456788"
	req.Student.EmergencyContact.Relationship.ID = relationshipID
	req.Student.EmergencyContact.Address = addr

	req.Family.ParentalStatus.ID = parentalStatusID
	req.Family.NatureOfResidence.ID = natureOfResidenceID
	req.Family.Brothers = &zero
	req.Family.Sisters = &zero
	req.Family.EmployedSiblings = &zero
	req.Family.OrdinalPosition = 1

	req.Family.Finance.IncomeRange.ID = incomeRangeID
	req.Family.Finance.WeeklyAllowance = 500.0

	req.Family.RelatedPersons = []RelatedPersonDTO{
		{
			FirstName:   "Jane",
			LastName:    "Doe",
			DateOfBirth: "1975-05-05",
			IsGuardian:  true,
			Relationship: StudentRelationshipType{
				ID: relationshipID,
			},
			EducationalAttainment: EducationalAttainment{
				ID: eduAttainmentID,
			},
		},
	}

	req.Education.NatureOfSchooling = "Continuous"
	req.Education.School = []SchoolDetailsDTO{
		{
			EducationalLevel: EducationalLevel{ID: educationalLevelID},
			SchoolName:       "Elementary School",
			SchoolType:       "Public",
			YearCompleted:    2012,
		},
	}

	// 1. Submit IIR
	iirID, err := svc.SubmitStudentIIR(context.Background(), "student-uuid", req)
	if err != nil {
		t.Fatalf("SubmitStudentIIR failed: %v", err)
	}
	if iirID == "" {
		t.Fatal("expected non-empty iirID")
	}

	// 2. Fetch/Verify IIR
	profile, err := svc.GetStudentProfile(context.Background(), iirID)
	if err != nil {
		t.Fatalf("GetStudentProfile failed: %v", err)
	}
	if profile == nil {
		t.Fatal("expected to find submitted IIR profile")
	}
	if profile.Student.StudentNumber != "2025-00001-MN-0" {
		t.Errorf(
			"unexpected student number: %s",
			profile.Student.StudentNumber,
		)
	}

	// 3. Update IIR
	req.Student.MobileNumber = "09999999999"
	_, err = svc.UpdateStudentIIR(context.Background(), iirID, req)
	if err != nil {
		t.Fatalf("UpdateStudentIIR failed: %v", err)
	}

	// 4. Verify Update
	updatedProfile, err := svc.GetStudentProfile(context.Background(), iirID)
	if err != nil {
		t.Fatalf("GetStudentProfile failed after update: %v", err)
	}
	if updatedProfile.Student.MobileNumber != "09999999999" {
		t.Errorf(
			"expected updated mobile number, got %s",
			updatedProfile.Student.MobileNumber,
		)
	}
}
