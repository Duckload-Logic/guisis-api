package integrations

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListStudents(
	ctx context.Context,
	req OGOSListStudentsRequest,
) ([]OGOSStudentView, int, error) {
	query := `
		SELECT
			sp.student_number AS student_number,
			u.idp_uuid AS idp_uuid,
			u.first_name AS first_name,
			u.middle_name AS middle_name,
			u.last_name AS last_name,
			u.suffix_name AS suffix_name,
			u.email AS email,
			sp.mobile_number AS mobile_number,
			p.id AS program_id,
			p.code AS program_code,
			p.program_name AS program_name,
			sp.year_level AS year_level,
			sp.section AS section
		FROM users u
		JOIN iir_records i ON i.user_id = u.id
		JOIN student_personal_info sp ON sp.iir_id = i.id
		JOIN programs p ON sp.program_id = p.id
		WHERE 1=1
		AND u.is_active = true
		AND i.is_submitted = true
	`
	var args []interface{}
	if req.Search != "" {
		query += `
			AND (
				u.first_name LIKE ?
				OR u.last_name LIKE ?
				OR sp.student_number LIKE ?
			)
		`
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if req.ProgramID != 0 {
		query += " AND p.id = ?"
		args = append(args, req.ProgramID)
	}

	if req.GenderID != 0 {
		query += " AND sp.gender = ?"
		genderStr := "Male"
		if req.GenderID == 2 {
			genderStr = "Female"
		}
		args = append(args, genderStr)
	}

	if req.YearLevel != 0 {
		query += " AND sp.year_level = ?"
		args = append(args, req.YearLevel)
	}

	// Get total count for pagination
	countQuery := "SELECT COUNT(*) FROM (" + query + ") AS count_table"
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	query += " LIMIT ? OFFSET ?"
	args = append(args, req.PageSize, req.GetOffset())

	var students []OGOSStudentView
	err = r.db.SelectContext(ctx, &students, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

func (r *Repository) GetStudentByStudentNumber(
	ctx context.Context,
	studentNumber string,
) (*OGOSStudentView, error) {
	query := `
		SELECT
			sp.student_number AS student_number,
			u.idp_uuid AS idp_uuid,
			u.first_name AS first_name,
			u.middle_name AS middle_name,
			u.last_name AS last_name,
			u.suffix_name AS suffix_name,
			u.email AS email,
			sp.mobile_number AS mobile_number,
			p.id AS program_id,
			p.code AS program_code,
			p.program_name AS program_name,
			sp.year_level AS year_level,
			sp.section AS section
		FROM users u
		JOIN iir_records i ON i.user_id = u.id
		JOIN student_personal_info sp ON sp.iir_id = i.id
		JOIN programs p ON sp.program_id = p.id
		WHERE (sp.student_number = ? OR u.idp_uuid = ?)
		AND u.is_active = true
		AND i.is_submitted = true
		LIMIT 1
	`

	var student OGOSStudentView
	err := r.db.GetContext(
		ctx,
		&student,
		query,
		studentNumber,
		studentNumber,
	)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *Repository) GetStudentByEmail(
	ctx context.Context,
	email string,
) (*OGOSStudentView, error) {
	query := `
		SELECT
			sp.student_number AS student_number,
			u.idp_uuid AS idp_uuid,
			u.first_name AS first_name,
			u.middle_name AS middle_name,
			u.last_name AS last_name,
			u.suffix_name AS suffix_name,
			u.email AS email,
			sp.mobile_number AS mobile_number,
			p.id AS program_id,
			p.code AS program_code,
			p.program_name AS program_name,
			sp.year_level AS year_level,
			sp.section AS section
		FROM users u
		JOIN iir_records i ON i.user_id = u.id
		JOIN student_personal_info sp ON sp.iir_id = i.id
		JOIN programs p ON sp.program_id = p.id
		WHERE (u.email = ? OR u.idp_uuid = ?)
		AND u.is_active = true
		AND i.is_submitted = true
		LIMIT 1
	`

	var student OGOSStudentView
	err := r.db.GetContext(ctx, &student, query, email, email)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *Repository) GetPersonalInfoByStudentNumber(
	ctx context.Context,
	studentNumber string,
) (*OGOSStudentPersonalInfoView, error) {
	query := `
		SELECT
			sp.student_number AS student_number,
			u.idp_uuid AS idp_uuid,
			CASE WHEN sp.gender = 'Male' THEN 1 ELSE 2 END AS gender_id,
			sp.gender AS gender_name,
			sp.date_of_birth AS date_of_birth,
			sp.place_of_birth AS place_of_birth,
			sp.height_m AS height_m,
			sp.weight_kg AS weight_kg,
			COALESCE(TRIM(CONCAT(
				ec.first_name,
				' ',
				COALESCE(CONCAT(ec.middle_name, ' '), ''),
				ec.last_name,
				COALESCE(CONCAT(' ', ec.suffix_name), '')
			)), '') AS emergency_contact_name,
			COALESCE(ec.contact_number, '') AS emergency_contact_number
		FROM student_personal_info sp
		JOIN iir_records i ON sp.iir_id = i.id
		JOIN users u ON i.user_id = u.id
		LEFT JOIN emergency_contacts ec ON ec.iir_id = sp.iir_id
		WHERE (sp.student_number = ? OR u.idp_uuid = ?)
		AND u.is_active = true
		AND i.is_submitted = true
		LIMIT 1
	`

	var student OGOSStudentPersonalInfoView
	err := r.db.GetContext(
		ctx,
		&student,
		query,
		studentNumber,
		studentNumber,
	)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *Repository) GetAddressByStudentNumber(
	ctx context.Context,
	studentNumber string,
) ([]OGOSStudentAddressView, error) {
	query := `
		SELECT
			sp.student_number AS student_number,
			u.idp_uuid AS idp_uuid,
			sa.address_type AS address_type,
			a.street_detail AS street_detail,
			a.barangay_code AS barangay_code,
			b.name AS barangay_name,
			a.city_code AS city_code,
			ci.name AS city_name,
			a.province_code AS province_code,
			p.name AS province_name,
			a.region_code AS region_code,
			r.name AS region_name
		FROM student_addresses sa
		JOIN addresses a ON a.id = sa.address_id
		JOIN barangays b ON a.barangay_code = b.code
		JOIN cities ci ON a.city_code = ci.code
		LEFT JOIN provinces p ON a.province_code = p.code
		JOIN regions r ON a.region_code = r.code
		JOIN student_personal_info sp ON sp.iir_id = sa.iir_id
		JOIN iir_records i ON sp.iir_id = i.id
		JOIN users u ON i.user_id = u.id
		WHERE (sp.student_number = ? OR u.idp_uuid = ?)
		AND u.is_active = true
		AND i.is_submitted = true
	`

	var addresses []OGOSStudentAddressView
	err := r.db.SelectContext(
		ctx,
		&addresses,
		query,
		studentNumber,
		studentNumber,
	)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}
