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
			u.first_name AS first_name,
			u.middle_name AS middle_name,
			u.last_name AS last_name,
			u.email AS email,
			sp.mobile_number AS mobile_number,
			c.id AS course_id,
			c.code AS course_code,
			c.course_name AS course_name,
			sp.year_level AS year_level,
			sp.section AS section
		FROM users u
		JOIN iir_records i ON i.user_id = u.id
		JOIN student_personal_info sp ON sp.iir_id = i.id
		JOIN courses c ON sp.course_id = c.id
		WHERE 1=1
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

	if req.CourseID != 0 {
		query += " AND c.id = ?"
		args = append(args, req.CourseID)
	}

	if req.GenderID != 0 {
		query += " AND sp.gender_id = ?"
		args = append(args, req.GenderID)
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
			u.first_name AS first_name,
			u.middle_name AS middle_name,
			u.last_name AS last_name,
			u.email AS email,
			sp.mobile_number AS mobile_number,
			c.id AS course_id,
			c.code AS course_code,
			c.course_name AS course_name,
			sp.year_level AS year_level,
			sp.section AS section
		FROM users u
		JOIN iir_records i ON i.user_id = u.id
		JOIN student_personal_info sp ON sp.iir_id = i.id
		JOIN courses c ON sp.course_id = c.id
		WHERE sp.student_number = ?
		LIMIT 1
	`

	var student OGOSStudentView
	err := r.db.GetContext(ctx, &student, query, studentNumber)
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
			u.first_name AS first_name,
			u.middle_name AS middle_name,
			u.last_name AS last_name,
			u.email AS email,
			sp.mobile_number AS mobile_number,
			c.id AS course_id,
			c.code AS course_code,
			c.course_name AS course_name,
			sp.year_level AS year_level,
			sp.section AS section
		FROM users u
		JOIN iir_records i ON i.user_id = u.id
		JOIN student_personal_info sp ON sp.iir_id = i.id
		JOIN courses c ON sp.course_id = c.id
		WHERE u.email = ?
		LIMIT 1
	`

	var student OGOSStudentView
	err := r.db.GetContext(ctx, &student, query, email)
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
			g.id AS gender_id,
			g.gender_name AS gender_name,
			sp.date_of_birth AS date_of_birth,
			sp.place_of_birth AS place_of_birth,
			sp.height_m AS height_m,
			sp.weight_kg AS weight_kg
		FROM student_personal_info sp
		JOIN genders g ON sp.gender_id = g.id
		WHERE sp.student_number = ?
		LIMIT 1
	`

	var student OGOSStudentPersonalInfoView
	err := r.db.GetContext(ctx, &student, query, studentNumber)
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
		WHERE sp.student_number = ?
	`

	var addresses []OGOSStudentAddressView
	err := r.db.SelectContext(ctx, &addresses, query, studentNumber)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}
