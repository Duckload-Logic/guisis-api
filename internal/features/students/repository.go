package students

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetDB() *sqlx.DB {
	return r.db
}

func (r *Repository) WithTransaction(
	ctx context.Context,
	fn func(datastore.DB) error,
) error {
	return datastore.RunInTransaction(ctx, r.db, fn)
}

// Lookup
func (r *Repository) GetEnrollmentYears(ctx context.Context) ([]int, error) {
	query := `
		SELECT DISTINCT
			CAST(SUBSTRING(student_number, 1, 4) AS UNSIGNED) AS year
		FROM student_personal_info
		WHERE student_number LIKE '____-%'
		ORDER BY year DESC
	`

	var years []int
	err := r.db.SelectContext(ctx, &years, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment years: %w", err)
	}

	return years, nil
}

func (r *Repository) GetGenders(ctx context.Context) ([]Gender, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM genders ORDER BY id
	`, datastore.GetColumns(Gender{}))

	var genders []Gender
	err := r.db.SelectContext(ctx, &genders, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get genders: %w", err)
	}
	return genders, nil
}

func (r *Repository) GetParentalStatusTypes(
	ctx context.Context,
) ([]ParentalStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM parental_status_types ORDER BY id
	`, datastore.GetColumns(ParentalStatusType{}))

	var statuses []ParentalStatusType
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status types: %w", err)
	}
	return statuses, nil
}

func (r *Repository) GetIncomeRanges(
	ctx context.Context,
) ([]IncomeRange, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM income_ranges ORDER BY id
	`, datastore.GetColumns(IncomeRange{}))

	var ranges []IncomeRange
	err := r.db.SelectContext(ctx, &ranges, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get income ranges: %w", err)
	}
	return ranges, nil
}

func (r *Repository) GetStudentSupportTypes(
	ctx context.Context,
) ([]StudentSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_support_types ORDER BY id
	`, datastore.GetColumns(StudentSupportType{}))

	var supportTypes []StudentSupportType
	err := r.db.SelectContext(ctx, &supportTypes, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student support types: %w", err)
	}
	return supportTypes, nil
}

func (r *Repository) GetSiblingSupportTypes(
	ctx context.Context,
) ([]SibilingSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM sibling_support_types ORDER BY id
	`, datastore.GetColumns(SibilingSupportType{}))

	var supportTypes []SibilingSupportType
	err := r.db.SelectContext(ctx, &supportTypes, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sibling support types: %w", err)
	}
	return supportTypes, nil
}

func (r *Repository) GetEducationalLevels(
	ctx context.Context,
) ([]EducationalLevel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_levels ORDER BY id
	`, datastore.GetColumns(EducationalLevel{}))

	var levels []EducationalLevel
	err := r.db.SelectContext(ctx, &levels, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational levels: %w", err)
	}
	return levels, nil
}

func (r *Repository) GetEducationalAttainments(
	ctx context.Context,
) ([]EducationalAttainment, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_attainments ORDER BY id
	`, datastore.GetColumns(EducationalAttainment{}))

	var attainments []EducationalAttainment
	err := r.db.SelectContext(ctx, &attainments, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational attainments: %w", err)
	}
	return attainments, nil
}

func (r *Repository) GetStudentStatuses(
	ctx context.Context,
) ([]StudentStatus, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_statuses ORDER BY id
	`, datastore.GetColumns(StudentStatus{}))

	var statuses []StudentStatus
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get student statuses: %w", err)
	}
	return statuses, nil
}

func (r *Repository) GetCourses(ctx context.Context) ([]Course, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM courses ORDER BY id
	`, datastore.GetColumns(Course{}))

	var courses []Course
	err := r.db.SelectContext(ctx, &courses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	return courses, nil
}

func (r *Repository) GetCivilStatusTypes(
	ctx context.Context,
) ([]CivilStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM civil_status_types ORDER BY id
	`, datastore.GetColumns(CivilStatusType{}))

	var statuses []CivilStatusType
	err := r.db.SelectContext(ctx, &statuses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get civil status types: %w", err)
	}
	return statuses, nil
}

func (r *Repository) GetReligions(ctx context.Context) ([]Religion, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM religions ORDER BY id
	`, datastore.GetColumns(Religion{}))

	var religions []Religion
	err := r.db.SelectContext(ctx, &religions, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get religions: %w", err)
	}
	return religions, nil
}

func (r *Repository) GetStudentRelationshipTypes(
	ctx context.Context,
) ([]StudentRelationshipType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_relationship_types ORDER BY id
	`, datastore.GetColumns(StudentRelationshipType{}))

	var relationships []StudentRelationshipType
	err := r.db.SelectContext(ctx, &relationships, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student relationship types: %w",
			err,
		)
	}
	return relationships, nil
}

func (r *Repository) GetNatureOfResidenceTypes(
	ctx context.Context,
) ([]NatureOfResidenceType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM nature_of_residence_types ORDER BY id
	`, datastore.GetColumns(NatureOfResidenceType{}))

	var residences []NatureOfResidenceType
	err := r.db.SelectContext(ctx, &residences, query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get nature of residence types: %w",
			err,
		)
	}
	return residences, nil
}

// Retrieve - Count
func (r *Repository) GetTotalStudentsCount(
	ctx context.Context,
	req ListStudentsRequest,
) (int, error) {
	query, args := r.applyStudentFilters(
		"SELECT COUNT(iir_id) FROM v_student_profiles WHERE 1=1",
		nil,
		req,
	)

	var total int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *Repository) applyStudentFilters(
	query string,
	args []interface{},
	req ListStudentsRequest,
) (string, []interface{}) {
	if args == nil {
		args = []interface{}{}
	}

	if req.CourseID > 0 {
		query += " AND course_id = ?"
		args = append(args, req.CourseID)
	}

	if req.GenderID > 0 {
		query += " AND gender_id = ?"
		args = append(args, req.GenderID)
	}

	if req.YearLevel > 0 {
		query += " AND year_level = ?"
		args = append(args, req.YearLevel)
	}

	if req.StatusID > 0 {
		query += " AND status_id = ?"
		args = append(args, req.StatusID)
	}

	if req.Search != "" {
		query += ` AND (first_name LIKE ?
                 OR last_name LIKE ?
                 OR email LIKE ?
                 OR student_number LIKE ?)`

		pattern := "%" + req.Search + "%"
		args = append(args, pattern, pattern, pattern, pattern)
	}

	return query, args
}

// Retrieve - List
func (r *Repository) ListStudents(
	ctx context.Context,
	req ListStudentsRequest,
) ([]StudentProfileView, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM v_student_profiles WHERE 1 = 1
	`, datastore.GetColumns(StudentProfileView{}))

	query, args := r.applyStudentFilters(
		query,
		nil,
		req,
	)

	allowedSortColumns := map[string]string{
		"last_name":      "last_name",
		"first_name":     "first_name",
		"year_level":     "year_level",
		"course_id":      "course_id",
		"created_at":     "created_at",
		"updated_at":     "updated_at",
		"student_number": "student_number",
		"iir_id":         "iir_id",
	}

	sortColumn, ok := allowedSortColumns[req.OrderBy]
	if !ok {
		sortColumn = allowedSortColumns["last_name"]
	}

	query += fmt.Sprintf(" ORDER BY %s ASC, iir_id ASC", sortColumn)
	query += " LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, req.GetOffset())

	var views []StudentProfileView
	err := r.db.SelectContext(ctx, &views, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	return views, nil
}

func (r *Repository) GetStudentBasicInfo(
	ctx context.Context,
	iirID string,
) (*StudentBasicInfoView, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM v_student_basic_info WHERE iir_id = ?
	`, datastore.GetColumns(StudentBasicInfoView{}))

	var view StudentBasicInfoView
	err := r.db.GetContext(ctx, &view, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student basic info: %w", err)
	}

	return &view, nil
}

func (r *Repository) GetIIRDraftByUserID(
	ctx context.Context,
	userID string,
) (*IIRDraft, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM iir_drafts WHERE user_id = ? LIMIT 1
	`, datastore.GetColumns(IIRDraft{}))

	var model IIRDraft
	err := r.db.GetContext(ctx, &model, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (r *Repository) GetStudentIIRByUserID(
	ctx context.Context,
	userID string,
) (*IIRRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM iir_records WHERE user_id = ? LIMIT 1
	`, datastore.GetColumns(IIRRecord{}))

	var model IIRRecord
	err := r.db.GetContext(ctx, &model, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (r *Repository) GetStudentIIR(
	ctx context.Context,
	iirID string,
) (*IIRRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM iir_records WHERE id = ? LIMIT 1
	`, datastore.GetColumns(IIRRecord{}))

	var model IIRRecord
	err := r.db.GetContext(ctx, &model, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &model, nil
}

func (r *Repository) GetStudentPersonalInfoView(
	ctx context.Context,
	iirID string,
) (*StudentPersonalInfoView, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM v_student_personal_info WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(StudentPersonalInfoView{}))

	var view StudentPersonalInfoView
	err := r.db.GetContext(ctx, &view, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf(
			"failed to get personal info view: %w",
			err,
		)
	}

	return &view, nil
}

func (r *Repository) GetCourseByID(
	ctx context.Context,
	courseID int,
) (*Course, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM courses WHERE id = ?
	`, datastore.GetColumns(Course{}))

	var model Course
	err := r.db.GetContext(ctx, &model, query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course by ID: %w", err)
	}

	return &model, nil
}

func (r *Repository) GetStudentFamilyBackground(
	ctx context.Context,
	iirID string,
) (*FamilyBackground, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM family_backgrounds WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(FamilyBackground{}))

	var model FamilyBackground
	err := r.db.GetContext(ctx, &model, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get family background: %w", err)
	}
	return &model, nil
}

func (r *Repository) GetParentalStatusByID(
	ctx context.Context,
	statusID int,
) (*ParentalStatusType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM parental_status_types WHERE id = ?
	`, datastore.GetColumns(ParentalStatusType{}))

	var model ParentalStatusType
	err := r.db.GetContext(ctx, &model, query, statusID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parental status by ID: %w", err)
	}
	return &model, nil
}

func (r *Repository) GetNatureOfResidenceByID(
	ctx context.Context,
	residenceID int,
) (*NatureOfResidenceType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM nature_of_residence_types WHERE id = ?
	`, datastore.GetColumns(NatureOfResidenceType{}))

	var model NatureOfResidenceType
	err := r.db.GetContext(ctx, &model, query, residenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get residence type by ID: %w", err)
	}
	return &model, nil
}

func (r *Repository) GetStudentSiblingSupport(
	ctx context.Context,
	fbID int,
) ([]StudentSiblingSupport, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_sibling_supports WHERE family_background_id = ?
	`, datastore.GetColumns(StudentSiblingSupport{}))

	var supports []StudentSiblingSupport
	err := r.db.SelectContext(ctx, &supports, query, fbID)
	return supports, err
}

func (r *Repository) GetSiblingSupportTypeByID(
	ctx context.Context,
	supportID int,
) (*SibilingSupportType, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM sibling_support_types WHERE id = ?
	`, datastore.GetColumns(SibilingSupportType{}))

	var model SibilingSupportType
	err := r.db.GetContext(ctx, &model, query, supportID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get sibling support type by ID: %w",
			err,
		)
	}
	return &model, nil
}

func (r *Repository) GetStudentFinancialInfoView(
	ctx context.Context,
	iirID string,
) (*StudentFinanceView, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM v_student_finances WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(StudentFinanceView{}))

	var view StudentFinanceView
	err := r.db.GetContext(ctx, &view, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf(
			"failed to get financial info view: %w",
			err,
		)
	}
	return &view, nil
}

func (r *Repository) GetFinancialSupportTypes(
	ctx context.Context,
	sfID int,
) ([]StudentSupportType, error) {
	query := `
		SELECT id, support_type_name
		FROM v_student_financial_supports
		WHERE sf_id = ?
	`

	var supports []StudentSupportType
	err := r.db.SelectContext(ctx, &supports, query, sfID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get financial support types: %w",
			err,
		)
	}
	return supports, nil
}

func (r *Repository) GetStudentHealthRecord(
	ctx context.Context,
	iirID string,
) (*StudentHealthRecord, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_health_records WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(StudentHealthRecord{}))

	var model StudentHealthRecord
	err := r.db.GetContext(ctx, &model, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get health record: %w", err)
	}
	return &model, nil
}

func (r *Repository) GetActivityOptions(
	ctx context.Context,
) ([]ActivityOption, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM activity_options WHERE is_active = 1 ORDER BY id
	`, datastore.GetColumns(ActivityOption{}))

	var options []ActivityOption
	err := r.db.SelectContext(ctx, &options, query)
	return options, err
}

func (r *Repository) GetStudentConsultations(
	ctx context.Context,
	iirID string,
) ([]StudentConsultation, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_consultations WHERE iir_id = ?
	`, datastore.GetColumns(StudentConsultation{}))

	var consultations []StudentConsultation
	err := r.db.SelectContext(ctx, &consultations, query, iirID)
	return consultations, err
}

func (r *Repository) GetStudentActivities(
	ctx context.Context,
	iirID string,
) ([]StudentActivity, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_activities WHERE iir_id = ?
	`, datastore.GetColumns(StudentActivity{}))

	var activities []StudentActivity
	err := r.db.SelectContext(ctx, &activities, query, iirID)
	return activities, err
}

func (r *Repository) GetActivityOptionByID(
	ctx context.Context,
	optionID int,
) (*ActivityOption, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM activity_options WHERE id = ?
	`, datastore.GetColumns(ActivityOption{}))

	var model ActivityOption
	err := r.db.GetContext(ctx, &model, query, optionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity option by ID: %w", err)
	}
	return &model, nil
}

func (r *Repository) GetStudentSubjectPreferences(
	ctx context.Context,
	iirID string,
) ([]StudentSubjectPreference, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_subject_preferences WHERE iir_id = ?
	`, datastore.GetColumns(StudentSubjectPreference{}))

	var prefs []StudentSubjectPreference
	err := r.db.SelectContext(ctx, &prefs, query, iirID)
	return prefs, err
}

func (r *Repository) GetStudentHobbies(
	ctx context.Context,
	iirID string,
) ([]StudentHobby, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_hobbies WHERE iir_id = ? ORDER BY priority_rank
	`, datastore.GetColumns(StudentHobby{}))

	var hobbies []StudentHobby
	err := r.db.SelectContext(ctx, &hobbies, query, iirID)
	return hobbies, err
}

func (r *Repository) GetStudentTestResults(
	ctx context.Context,
	iirID string,
) ([]TestResult, error) {
	query := fmt.Sprintf(`
		SELECT
			%s
		FROM
			student_test_results
		WHERE
			iir_id = ?
		ORDER BY
			test_date DESC
	`, datastore.GetColumns(TestResult{}))

	var results []TestResult
	err := r.db.SelectContext(ctx, &results, query, iirID)
	return results, err
}

func (r *Repository) GetStudentAddresses(
	ctx context.Context,
	iirID string,
) ([]StudentAddress, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_addresses WHERE iir_id = ?
	`, datastore.GetColumns(StudentAddress{}))

	var addresses []StudentAddress
	err := r.db.SelectContext(ctx, &addresses, query, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student addresses: %w", err)
	}

	return addresses, nil
}

func (r *Repository) GetStudentEducationalBackground(
	ctx context.Context,
	iirID string,
) (*EducationalBackground, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_backgrounds WHERE iir_id = ? LIMIT 1
	`, datastore.GetColumns(EducationalBackground{}))

	var model EducationalBackground
	err := r.db.GetContext(ctx, &model, query, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf(
			"failed to get student educational background: %w",
			err,
		)
	}

	return &model, nil
}

func (r *Repository) GetSchoolDetailsByEBID(
	ctx context.Context,
	ebID int,
) ([]SchoolDetails, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM school_details WHERE eb_id = ?
		ORDER BY educational_level_id ASC
	`, datastore.GetColumns(SchoolDetails{}))

	var details []SchoolDetails
	err := r.db.SelectContext(ctx, &details, query, ebID)
	if err != nil {
		return nil, fmt.Errorf("failed to get school details: %w", err)
	}

	return details, nil
}

func (r *Repository) GetEducationalLevelByID(
	ctx context.Context,
	levelID int,
) (*EducationalLevel, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM educational_levels WHERE id = ?
	`, datastore.GetColumns(EducationalLevel{}))
	var model EducationalLevel
	err := r.db.GetContext(ctx, &model, query, levelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational level by ID: %w", err)
	}

	return &model, nil
}

func (r *Repository) GetStudentRelatedPersonsView(
	ctx context.Context,
	iirID string,
) ([]RelatedPersonView, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM v_related_persons WHERE iir_id = ?
	`, datastore.GetColumns(RelatedPersonView{}))

	var views []RelatedPersonView
	err := r.db.SelectContext(ctx, &views, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get related persons view: %w",
			err,
		)
	}

	return views, nil
}

func (r *Repository) UpsertIIRDraft(
	ctx context.Context,
	draft IIRDraft,
) (int, error) {
	exclude := []string{"id", "created_at"}
	cols, vals := datastore.GetInsertStatement(IIRDraft{}, exclude)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		IIRDraft{},
		exclude,
	)

	query := fmt.Sprintf(`
		INSERT INTO iir_drafts (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, onDuplicate)

	result, err := r.db.NamedExecContext(ctx, query, &draft)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	return int(lastID), err
}

func (r *Repository) UpsertIIRRecord(
	ctx context.Context,
	tx datastore.DB,
	iir *IIRRecord,
) (string, error) {
	cols, vals := datastore.GetInsertStatement(IIRRecord{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		IIRRecord{},
		nil,
	)
	query := fmt.Sprintf(`
		INSERT INTO
			iir_records (id, %s)
		VALUES
			(:id, %s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	_, err := tx.NamedExecContext(ctx, query, iir)
	return iir.ID, err
}

func (r *Repository) UpsertStudentPersonalInfo(
	ctx context.Context,
	tx datastore.DB,
	info *StudentPersonalInfo,
) error {
	cols, vals := datastore.GetInsertStatement(StudentPersonalInfo{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentPersonalInfo{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			student_personal_info (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	_, err := tx.NamedExecContext(ctx, query, info)
	return err
}

func (r *Repository) UpsertEmergencyContact(
	ctx context.Context,
	tx datastore.DB,
	ec *EmergencyContact,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(EmergencyContact{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		EmergencyContact{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			emergency_contacts (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, ec)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) UpsertStudentAddress(
	ctx context.Context,
	tx datastore.DB,
	sa *StudentAddress,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(StudentAddress{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentAddress{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			student_addresses (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, sa)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) DeleteStudentAddressesByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_addresses WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	return err
}

func (r *Repository) UpsertRelatedPerson(
	ctx context.Context,
	tx datastore.DB,
	rp *RelatedPerson,
) (int, error) {
	exclude := []string{"id", "created_at"}
	cols, vals := datastore.GetInsertStatement(RelatedPerson{}, exclude)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		RelatedPerson{},
		exclude,
	)

	query := fmt.Sprintf(`
		INSERT INTO related_persons (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, rp)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) UpsertStudentRelatedPerson(
	ctx context.Context,
	tx datastore.DB,
	srp *StudentRelatedPerson,
) error {
	cols, vals := datastore.GetInsertStatement(StudentRelatedPerson{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentRelatedPerson{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			student_related_persons (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	_, err := tx.NamedExecContext(ctx, query, srp)
	return err
}

func (r *Repository) DeleteEmergencyContactByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	var addressID int
	queryGet := `
		SELECT address_id FROM emergency_contacts
		WHERE iir_id = ? LIMIT 1
	`
	err := tx.GetContext(ctx, &addressID, queryGet, iirID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return fmt.Errorf("failed to get EC address ID: %w", err)
	}

	queryDelEC := `DELETE FROM emergency_contacts WHERE iir_id = ?`
	_, err = tx.ExecContext(ctx, queryDelEC, iirID)
	if err != nil {
		return fmt.Errorf("failed to delete EC: %w", err)
	}

	queryDelAddr := `DELETE FROM addresses WHERE id = ?`
	_, err = tx.ExecContext(ctx, queryDelAddr, addressID)
	if err != nil {
		return fmt.Errorf("failed to delete EC address: %w", err)
	}

	return nil
}

func (r *Repository) DeleteStudentRelatedPersons(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	var ids []int
	queryGet := `
		SELECT related_person_id FROM student_related_persons
		WHERE iir_id = ?
	`
	err := tx.SelectContext(ctx, &ids, queryGet, iirID)
	if err != nil {
		return fmt.Errorf("failed to get related person IDs: %w", err)
	}

	queryDel := `DELETE FROM student_related_persons WHERE iir_id = ?`
	_, err = tx.ExecContext(ctx, queryDel, iirID)
	if err != nil {
		return fmt.Errorf("failed to delete student related links: %w", err)
	}

	for _, id := range ids {
		_, err = tx.ExecContext(
			ctx,
			"DELETE FROM related_persons WHERE id = ?",
			id,
		)
		if err != nil {
			return fmt.Errorf("failed to delete related person: %w", err)
		}
	}

	return nil
}

func (r *Repository) UpsertFamilyBackground(
	ctx context.Context,
	tx datastore.DB,
	fb *FamilyBackground,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(FamilyBackground{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		FamilyBackground{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			family_backgrounds (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, fb)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) CreateStudentSiblingSupport(
	ctx context.Context,
	tx datastore.DB,
	sss *StudentSiblingSupport,
) error {
	cols, vals := datastore.GetInsertStatement(StudentSiblingSupport{}, nil)
	query := fmt.Sprintf(
		`INSERT INTO student_sibling_supports (%s) VALUES (%s)`,
		cols,
		vals,
	)
	_, err := tx.NamedExecContext(ctx, query, sss)
	return err
}

func (r *Repository) DeleteStudentSiblingSupportsByFamilyID(
	ctx context.Context,
	tx datastore.DB,
	familyBackgroundID int,
) error {
	query := `
		DELETE FROM student_sibling_supports
		WHERE family_background_id = ?
	`
	_, err := tx.ExecContext(ctx, query, familyBackgroundID)
	return err
}

func (r *Repository) UpsertEducationalBackground(
	ctx context.Context,
	tx datastore.DB,
	eb *EducationalBackground,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(EducationalBackground{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		EducationalBackground{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			educational_backgrounds (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, eb)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) UpsertSchoolDetails(
	ctx context.Context,
	tx datastore.DB,
	sd *SchoolDetails,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(SchoolDetails{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		SchoolDetails{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			school_details (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, sd)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) DeleteSchoolDetailsByEBID(
	ctx context.Context,
	tx datastore.DB,
	ebID int,
) error {
	query := `DELETE FROM school_details WHERE eb_id = ?`
	_, err := tx.ExecContext(ctx, query, ebID)
	return err
}

func (r *Repository) UpsertStudentHealthRecord(
	ctx context.Context,
	tx datastore.DB,
	hr *StudentHealthRecord,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(StudentHealthRecord{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentHealthRecord{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			student_health_records (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE %s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, hr)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) UpsertStudentConsultation(
	ctx context.Context,
	tx datastore.DB,
	sc *StudentConsultation,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(StudentConsultation{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentConsultation{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO
			student_consultations (%s)
		VALUES
			(%s)
		ON DUPLICATE KEY UPDATE
			%s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, sc)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) DeleteStudentConsultationsByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_consultations WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	return err
}

func (r *Repository) UpsertStudentFinance(
	ctx context.Context,
	tx datastore.DB,
	sf *StudentFinance,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(StudentFinance{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentFinance{},
		nil,
	)

	query := fmt.Sprintf(`
		INSERT INTO student_finances (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s
	`, cols, vals, onDuplicate)

	result, err := tx.NamedExecContext(ctx, query, sf)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) CreateStudentFinancialSupport(
	ctx context.Context,
	tx datastore.DB,
	sfs *StudentFinancialSupport,
) error {
	cols, vals := datastore.GetInsertStatement(StudentFinancialSupport{}, nil)
	query := fmt.Sprintf(
		`INSERT INTO student_financial_supports (%s) VALUES (%s)`,
		cols,
		vals,
	)
	_, err := tx.NamedExecContext(ctx, query, sfs)
	return err
}

func (r *Repository) DeleteStudentFinancialSupportsByFinanceID(
	ctx context.Context,
	tx datastore.DB,
	financeID int,
) error {
	query := `DELETE FROM student_financial_supports WHERE sf_id = ?`
	_, err := tx.ExecContext(ctx, query, financeID)
	return err
}

func (r *Repository) CreateStudentActivity(
	ctx context.Context,
	tx datastore.DB,
	sa *StudentActivity,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(StudentActivity{}, nil)
	query := fmt.Sprintf(
		`INSERT INTO student_activities (%s) VALUES (%s)`,
		cols,
		vals,
	)
	result, err := tx.NamedExecContext(ctx, query, sa)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) DeleteStudentActivitiesByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_activities WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	return err
}

func (r *Repository) CreateStudentSubjectPreference(
	ctx context.Context,
	tx datastore.DB,
	ssp *StudentSubjectPreference,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(StudentSubjectPreference{}, nil)
	query := fmt.Sprintf(
		`INSERT INTO student_subject_preferences (%s) VALUES (%s)`,
		cols,
		vals,
	)
	result, err := tx.NamedExecContext(ctx, query, ssp)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) DeleteStudentSubjectPreferencesByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_subject_preferences WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	return err
}

func (r *Repository) CreateStudentHobby(
	ctx context.Context,
	tx datastore.DB,
	sh *StudentHobby,
) (int, error) {
	cols, vals := datastore.GetInsertStatement(StudentHobby{}, nil)
	query := fmt.Sprintf(
		`INSERT INTO student_hobbies (%s) VALUES (%s)`,
		cols,
		vals,
	)
	result, err := tx.NamedExecContext(ctx, query, sh)
	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return int(lastID), nil
}

func (r *Repository) DeleteStudentHobbiesByIIRID(
	ctx context.Context,
	tx datastore.DB,
	iirID string,
) error {
	query := `DELETE FROM student_hobbies WHERE iir_id = ?`
	_, err := tx.ExecContext(ctx, query, iirID)
	return err
}

func (r *Repository) SaveStudentCOR(
	ctx context.Context,
	tx datastore.DB,
	cor StudentCOR,
) error {
	cols, vals := datastore.GetInsertStatement(StudentCOR{}, nil)
	onDuplicate := datastore.GetOnDuplicateKeyUpdateStatement(
		StudentCOR{},
		nil,
	)

	query := fmt.Sprintf(
		`INSERT INTO student_cors (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s`,
		cols,
		vals,
		onDuplicate,
	)
	_, err := tx.NamedExecContext(ctx, query, &cor)
	return err
}

func (r *Repository) GetLatestCORsByUserIDs(
	ctx context.Context,
	userIDs []string,
) (map[string]string, error) {
	if len(userIDs) == 0 {
		return map[string]string{}, nil
	}

	placeholders := strings.Repeat("?,", len(userIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf(`
		SELECT student_id, file_url
		FROM v_student_cors_files
		WHERE student_id IN (%s)
	`, placeholders)

	args := make([]interface{}, len(userIDs))
	for i, id := range userIDs {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make(map[string]string)
	for rows.Next() {
		var studentID, fileURL string
		if err := rows.Scan(&studentID, &fileURL); err != nil {
			return nil, err
		}
		results[studentID] = fileURL
	}

	return results, nil
}

func (r *Repository) GetStudentCORByUserID(
	ctx context.Context,
	userID string,
) (StudentCOR, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM v_student_current_cors WHERE student_id = ?
	`, datastore.GetColumns(StudentCOR{}))

	var model StudentCOR
	err := r.db.GetContext(ctx, &model, query, userID)
	return model, err
}

func (r *Repository) GetStudentCORsByUserID(
	ctx context.Context,
	userID string,
) ([]StudentCOR, error) {
	query := fmt.Sprintf(`
		SELECT %s FROM student_cors WHERE student_id = ?
	`, datastore.GetColumns(StudentCOR{}))

	var models []StudentCOR
	err := r.db.SelectContext(ctx, &models, query, userID)
	return models, err
}

func (r *Repository) GetAcademicSetting(
	ctx context.Context,
) (*AcademicSetting, error) {
	cols := datastore.GetColumns(AcademicSetting{})
	query := fmt.Sprintf(`
		SELECT %s FROM academic_settings WHERE id = 1 LIMIT 1
	`, cols)
	var setting AcademicSetting
	err := r.db.GetContext(ctx, &setting, query)
	if err != nil {
		return nil, fmt.Errorf(
			"[Repository] {GetAcademicSetting}: %w", err,
		)
	}
	return &setting, nil
}

func (r *Repository) UpdateAcademicSetting(
	ctx context.Context,
	tx datastore.DB,
	yearStart, yearEnd, term int,
) error {
	query := `
		UPDATE academic_settings
		SET current_year_start = ?,
		    current_year_end   = ?,
		    current_term       = ?
		WHERE id = 1
	`
	_, err := tx.ExecContext(ctx, query, yearStart, yearEnd, term)
	if err != nil {
		return fmt.Errorf(
			"[Repository] {UpdateAcademicSetting}: %w", err,
		)
	}
	return nil
}

func (r *Repository) GetStudentSignificantNotes(
	ctx context.Context,
	iirID string,
) ([]SignificantNote, error) {
	query := `
		SELECT id, iir_id, appointment_id, admission_slip_id,
		       note, remarks, created_at, updated_at
		FROM significant_notes
		WHERE iir_id = ?
	`

	var results []SignificantNote
	err := r.db.SelectContext(ctx, &results, query, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"[Repository] {GetStudentSignificantNotes}: %w",
			err,
		)
	}

	return results, nil
}
