package analytics

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetTotalStudents(
	ctx context.Context,
	year int,
	courseID int,
) (int, error) {
	var total int
	filter, args := r.buildFilter(year, courseID)
	query := "SELECT COUNT(*) FROM student_personal_info spi WHERE 1=1" + filter
	err := r.db.GetContext(ctx, &total, query, args...)
	return total, err
}

func (r *Repository) GetGenderStats(
	ctx context.Context,
	year int,
	courseID int,
) ([]DemographicStat, error) {
	var results []DemographicStat
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			g.gender_name as category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 ` + filter + `
		GROUP BY g.gender_name;`

	err := r.db.SelectContext(ctx, &results, query, args...)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *Repository) GetTotalReports(ctx context.Context) (int, error) {
	var total int
	err := r.db.GetContext(
		ctx,
		&total,
		"SELECT COUNT(*) FROM significant_notes",
	)
	return total, err
}

func (r *Repository) GetTotalAppointments(ctx context.Context) (int, error) {
	var total int
	err := r.db.GetContext(
		ctx,
		&total,
		`SELECT COUNT(*) FROM appointments
		 WHERE status_id != (SELECT id FROM statuses WHERE name = 'Cancelled')`,
	)
	return total, err
}

func (r *Repository) GetTotalSlips(ctx context.Context) (int, error) {
	var total int
	err := r.db.GetContext(
		ctx,
		&total,
		"SELECT COUNT(*) FROM admission_slips",
	)
	return total, err
}

func (r *Repository) GetMonthlyVisitorStats(
	ctx context.Context,
	timeRange string,
) ([]MonthlyVisitorStatDTO, error) {
	var interval, format, groupBy, baseDate string

	switch timeRange {
	case "daily":
		interval = "29 DAY"
		format = "%d %b"
		groupBy = "%Y-%m-%d"
		baseDate = "CURDATE()"
	case "weekly":
		interval = "11 WEEK"
		format = "Week %u"
		groupBy = "%Y-%u"
		baseDate = "DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY)" // Start of current week
	case "yearly":
		interval = "4 YEAR"
		format = "%Y"
		groupBy = "%Y"
		baseDate = "DATE_FORMAT(NOW(), '%Y-01-01')"
	case "monthly":
		fallthrough
	default:
		interval = "11 MONTH"
		format = "%b"
		groupBy = "%Y-%m"
		baseDate = "DATE_FORMAT(NOW(), '%Y-%m-01')"
	}

	query := `
		SELECT
			DATE_FORMAT(created_at, '` + format + `') as period,
			DATE_FORMAT(created_at, '` + format + `') as month,
			SUM(CASE WHEN action = 'LOGIN_SUCCESS' THEN 1 ELSE 0 END) as logins,
			COUNT(*) as activity,
			SUM(CASE WHEN action = 'LOGIN_SUCCESS' THEN 1 ELSE 0 END) as count
		FROM system_logs
		WHERE created_at >= DATE_SUB(` + baseDate + `, INTERVAL ` + interval + `)
		GROUP BY DATE_FORMAT(created_at, '` + groupBy + `'), period, month
		ORDER BY DATE_FORMAT(created_at, '` + groupBy + `') ASC;
	`
	var stats []MonthlyVisitorStatDTO
	err := r.db.SelectContext(ctx, &stats, query)
	return stats, err
}

func (r *Repository) GetMonthlyAppointmentStats(
	ctx context.Context,
	timeRange string,
) ([]MonthlyVisitorStatDTO, error) {
	var interval, format, groupBy, baseDate string

	switch timeRange {
	case "daily":
		interval = "29 DAY"
		format = "%d %b"
		groupBy = "%Y-%m-%d"
		baseDate = "CURDATE()"
	case "weekly":
		interval = "11 WEEK"
		format = "Week %u"
		groupBy = "%Y-%u"
		baseDate = "DATE_SUB(CURDATE(), INTERVAL WEEKDAY(CURDATE()) DAY)"
	case "yearly":
		interval = "4 YEAR"
		format = "%Y"
		groupBy = "%Y"
		baseDate = "DATE_FORMAT(NOW(), '%Y-01-01')"
	case "monthly":
		fallthrough
	default:
		interval = "11 MONTH"
		format = "%b"
		groupBy = "%Y-%m"
		baseDate = "DATE_FORMAT(NOW(), '%Y-%m-01')"
	}

	query := `
		SELECT
			DATE_FORMAT(when_date, '` + format + `') as period,
			DATE_FORMAT(when_date, '` + format + `') as month,
			0 as logins,
			0 as activity,
			COUNT(*) as count
		FROM appointments a
		JOIN statuses s ON s.id = a.status_id
		WHERE when_date >= DATE_SUB(` + baseDate + `, INTERVAL ` + interval + `)
		  AND UPPER(s.name) = 'COMPLETED'
		GROUP BY DATE_FORMAT(when_date, '` + groupBy + `'), period, month
		ORDER BY DATE_FORMAT(when_date, '` + groupBy + `') ASC;
	`

	stats := []MonthlyVisitorStatDTO{}
	err := r.db.SelectContext(ctx, &stats, query)
	return stats, err
}

// --- PERSONAL INFORMATION ---

func (r *Repository) GetAgeStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			CAST(TIMESTAMPDIFF(YEAR, spi.date_of_birth, CURDATE()) AS CHAR) AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE spi.date_of_birth IS NOT NULL ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetCivilStatusStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(status_name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN civil_status_types cs ON spi.civil_status_id = cs.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 ` + filter + `
		GROUP BY category
		ORDER BY rank_pos ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetReligionStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(religion_name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN religions rel ON spi.religion_id = rel.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 ` + filter + `
		GROUP BY category
		ORDER BY rank_pos ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetCityAddressStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(c.name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN student_addresses sa ON spi.iir_id = sa.iir_id
		LEFT JOIN addresses a ON sa.address_id = a.id
		LEFT JOIN cities c ON a.city_code = c.code
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 AND sa.address_type = "Residential" ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

// --- FAMILY & FINANCIAL BACKGROUND ---

func (r *Repository) GetMonthlyIncomeStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(ir.range_text, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN student_finances sf ON spi.iir_id = sf.iir_id
		LEFT JOIN income_ranges ir ON sf.monthly_family_income_range_id = ir.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetOrdinalPositionStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			CAST(fb.ordinal_position AS CHAR) AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetFatherEducationStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(rp.educational_level, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN student_related_persons srp ON spi.iir_id = srp.iir_id
		JOIN related_persons rp ON srp.related_person_id = rp.id
		JOIN student_relationship_types srt ON srp.relationship_id = srt.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE srt.relationship_name = 'Father' ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetMotherEducationStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(rp.educational_level, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN student_related_persons srp ON spi.iir_id = srp.iir_id
		JOIN related_persons rp ON srp.related_person_id = rp.id
		JOIN student_relationship_types srt ON srp.relationship_id = srt.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE srt.relationship_name = 'Mother' ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetParentsMaritalStatusStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(pst.status_name, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		LEFT JOIN parental_status_types pst ON fb.parental_status_id = pst.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetQuietStudyPlaceStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			CASE WHEN fb.have_quiet_place_to_study = 1 THEN 'Yes' ELSE 'No' END AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

// --- ACADEMIC BACKGROUND ---

func (r *Repository) GetHSGWAStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			CASE
				WHEN spi.high_school_gwa >= 95.00 THEN '95.00 - 100.00'
				WHEN spi.high_school_gwa >= 90.00 THEN '90.00 - 94.99'
				WHEN spi.high_school_gwa >= 85.00 THEN '85.00 - 89.99'
				WHEN spi.high_school_gwa >= 80.00 THEN '80.00 - 84.99'
				WHEN spi.high_school_gwa >= 75.00 THEN '75.00 - 79.99'
				ELSE 'Below 75.00'
			END AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE spi.high_school_gwa IS NOT NULL AND spi.high_school_gwa > 0 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetElementaryStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN school_details sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE el.level_name = 'Elementary' ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetJuniorHighStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN school_details sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE el.level_name = 'Junior High School' ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetSeniorHighStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN school_details sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE el.level_name = 'Senior High School' ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetNatureOfSchoolingStats(
	ctx context.Context, year int, courseID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, courseID)
	query := `
		SELECT
			COALESCE(eb.nature_of_schooling, 'Not Indicated') AS category,
			SUM(CASE WHEN g.gender_name = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN g.gender_name = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		LEFT JOIN genders g ON spi.gender_id = g.id
		WHERE 1=1 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

// --- HELPERS ---

func (r *Repository) buildFilter(
	year int,
	courseID int,
) (string, []interface{}) {
	filter := ""
	args := []interface{}{}

	if year > 0 {
		filter += " AND spi.student_number LIKE ?"
		args = append(args, fmt.Sprintf("%d-%%", year))
	}
	if courseID > 0 {
		filter += " AND spi.course_id = ?"
		args = append(args, courseID)
	}
	return filter, args
}

func (r *Repository) executeStatQuery(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]DemographicStat, error) {
	results := make([]DemographicStat, 0)
	var err error
	if len(args) > 0 {
		err = r.db.SelectContext(ctx, &results, query, args...)
	} else {
		err = r.db.SelectContext(ctx, &results, query)
	}
	if err != nil {
		return nil, err
	}
	return results, nil
}
