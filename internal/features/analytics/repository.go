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
	programID int,
) (int, error) {
	var total int
	filter, args := r.buildFilter(year, programID)
	query := "SELECT COUNT(*) FROM student_personal_info spi WHERE 1=1" + filter
	err := r.db.GetContext(ctx, &total, query, args...)
	return total, err
}

func (r *Repository) GetGenderStats(
	ctx context.Context,
	year int,
	programID int,
) ([]DemographicStat, error) {
	var results []DemographicStat
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			spi.gender as category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		WHERE 1=1 ` + filter + `
		GROUP BY spi.gender;`

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

func (r *Repository) GetStudentsTrend(ctx context.Context) (int, error) {
	var total int
	err := r.db.GetContext(
		ctx,
		&total,
		`SELECT COUNT(*) FROM student_personal_info
		 WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)`,
	)
	return total, err
}

func (r *Repository) GetReportsTrend(ctx context.Context) (int, error) {
	var total int
	err := r.db.GetContext(
		ctx,
		&total,
		`SELECT COUNT(*) FROM significant_notes
		 WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)`,
	)
	return total, err
}

func (r *Repository) GetAppointmentsTrend(ctx context.Context) (int, error) {
	var total int
	err := r.db.GetContext(
		ctx,
		&total,
		`SELECT COUNT(*) FROM appointments
		 WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)`,
	)
	return total, err
}

func (r *Repository) GetSlipsTrend(ctx context.Context) (int, error) {
	var total int
	err := r.db.GetContext(
		ctx,
		&total,
		`SELECT COUNT(*) FROM admission_slips
		 WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)`,
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
		// Start of current week
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
			DATE_FORMAT(
				DATE_ADD(created_at, INTERVAL 8 HOUR),
				'` + format + `'
			) as period,
			DATE_FORMAT(
				DATE_ADD(created_at, INTERVAL 8 HOUR),
				'` + format + `'
			) as month,
			SUM(
				CASE WHEN action = 'LOGIN_SUCCESS' THEN 1 ELSE 0 END
			) as logins,
			COUNT(*) as activity,
			SUM(
				CASE WHEN action = 'LOGIN_SUCCESS' THEN 1 ELSE 0 END
			) as count
		FROM system_logs
		WHERE created_at >= DATE_SUB(` + baseDate + `,
			INTERVAL ` + interval + `)
		GROUP BY DATE_FORMAT(
			DATE_ADD(created_at, INTERVAL 8 HOUR),
			'` + groupBy + `'
		), period, month
		ORDER BY DATE_FORMAT(
			DATE_ADD(created_at, INTERVAL 8 HOUR),
			'` + groupBy + `'
		) ASC;
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
		WHERE when_date >= DATE_SUB(` + baseDate + `,
			INTERVAL ` + interval + `)
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
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			CAST(TIMESTAMPDIFF(
				YEAR, spi.date_of_birth, CURDATE()
			) AS CHAR) AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		WHERE spi.date_of_birth IS NOT NULL ` + filter + `
		GROUP BY category
		ORDER BY category ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetCivilStatusStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(status_name, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		LEFT JOIN civil_status_types cs ON spi.civil_status_id = cs.id
		WHERE 1=1 ` + filter + `
		GROUP BY category
		ORDER BY category ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetReligionStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			CASE
				WHEN rel.religion_name = 'Others'
					AND spi.other_religion_text IS NOT NULL
					AND spi.other_religion_text != ''
					THEN 'Others (' || spi.other_religion_text || ')'
				ELSE COALESCE(rel.religion_name, 'Not Indicated')
			END AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		LEFT JOIN religions rel ON spi.religion_id = rel.id
		WHERE 1=1 ` + filter + `
		GROUP BY category
		ORDER BY
			CASE WHEN category = 'Not Indicated' THEN 1 ELSE 0 END ASC,
			category ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetCityAddressStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			CASE
				WHEN c.type = 'SubMun' THEN 'City of Manila'
				ELSE COALESCE(c.name, 'Not Indicated')
			END AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		LEFT JOIN student_addresses sa ON spi.iir_id = sa.iir_id
		LEFT JOIN addresses a ON sa.address_id = a.id
		LEFT JOIN cities c ON a.city_code = c.code
		WHERE sa.address_type = "Residential" ` + filter + `
		GROUP BY category
		ORDER BY rank_pos ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

// --- FAMILY & FINANCIAL BACKGROUND ---

func (r *Repository) GetMonthlyIncomeStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(ir.range_text, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN student_finances sf ON spi.iir_id = sf.iir_id
		LEFT JOIN income_ranges ir ON sf.monthly_family_income_range_id = ir.id
		WHERE 1=1 ` + filter + `
		GROUP BY category
		ORDER BY category ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetOrdinalPositionStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			CASE
				WHEN fb.ordinal_position = 1 AND (fb.brothers + fb.sisters = 0) THEN 'Only Child'
				WHEN fb.ordinal_position = 1 AND (fb.brothers + fb.sisters > 0) THEN 'Eldest'
				WHEN fb.ordinal_position = (fb.brothers + fb.sisters + 1) THEN 'Youngest'
				WHEN fb.ordinal_position > 1 AND fb.ordinal_position < (fb.brothers + fb.sisters + 1) THEN 'Middle'
				WHEN fb.ordinal_position = 0 THEN 'Not Indicated'
				ELSE 'Others'
			END AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		WHERE 1=1 ` + filter + `
		GROUP BY category
		ORDER BY category ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetFatherEducationStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(ea.name, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN student_related_persons srp ON spi.iir_id = srp.iir_id
		JOIN related_persons rp ON srp.related_person_id = rp.id
		JOIN student_relationship_types srt ON srp.relationship_id = srt.id
		LEFT JOIN educational_attainments ea ON rp.educational_attainment_id = ea.id
		WHERE srt.relationship_name = 'Father' ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category,
			'Doctorate Degree',
			'Master''s Degree',
			'College Graduate',
			'College Undergraduate',
			'Vocational',
			'High School Graduate',
			'High School Undergraduate',
			'Elementary Graduate',
			'Elementary Undergraduate',
			'Not Indicated'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetMotherEducationStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(ea.name, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN student_related_persons srp ON spi.iir_id = srp.iir_id
		JOIN related_persons rp ON srp.related_person_id = rp.id
		JOIN student_relationship_types srt ON srp.relationship_id = srt.id
		LEFT JOIN educational_attainments ea ON rp.educational_attainment_id = ea.id
		WHERE srt.relationship_name = 'Mother' ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category,
			'Doctorate Degree',
			'Master''s Degree',
			'College Graduate',
			'College Undergraduate',
			'Vocational',
			'High School Graduate',
			'High School Undergraduate',
			'Elementary Graduate',
			'Elementary Undergraduate',
			'Not Indicated'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetParentsMaritalStatusStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(pst.status_name, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		LEFT JOIN parental_status_types pst ON fb.parental_status_id = pst.id
		WHERE 1=1 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetQuietStudyPlaceStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			CASE WHEN fb.have_quiet_place_to_study = 1 THEN 'Yes' ELSE 'No' END AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN family_backgrounds fb ON spi.iir_id = fb.iir_id
		WHERE 1=1 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

// --- ACADEMIC BACKGROUND ---

func (r *Repository) GetHSGWAStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			CASE
				WHEN spi.high_school_gwa >= 97.00 THEN '97.00 - 100.00'
				WHEN spi.high_school_gwa >= 94.00 THEN '94.00 - 96.99'
				WHEN spi.high_school_gwa >= 91.00 THEN '91.00 - 93.99'
				WHEN spi.high_school_gwa >= 88.00 THEN '88.00 - 90.99'
				WHEN spi.high_school_gwa >= 85.00 THEN '85.00 - 87.99'
				WHEN spi.high_school_gwa >= 82.00 THEN '82.00 - 84.99'
				WHEN spi.high_school_gwa >= 79.00 THEN '79.00 - 81.99'
				WHEN spi.high_school_gwa >= 76.00 THEN '76.00 - 78.99'
				WHEN spi.high_school_gwa >= 75.00 THEN '75.00'
				ELSE 'Below 75.00'
			END AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		WHERE spi.high_school_gwa IS NOT NULL
		  AND spi.high_school_gwa > 0 ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category,
			'97.00 - 100.00',
			'94.00 - 96.99',
			'91.00 - 93.99',
			'88.00 - 90.99',
			'85.00 - 87.99',
			'82.00 - 84.99',
			'79.00 - 81.99',
			'76.00 - 78.99',
			'75.00',
			'Below 75.00'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetElementaryStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN (
			SELECT s1.*
			FROM school_details s1
			JOIN (
				SELECT eb_id, educational_level_id,
					MAX(year_completed) AS max_yr
				FROM school_details
				GROUP BY eb_id, educational_level_id
			) s2 ON s1.eb_id = s2.eb_id
				AND s1.educational_level_id = s2.educational_level_id
				AND s1.year_completed = s2.max_yr
		) sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		WHERE el.level_name = 'Elementary' ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category, 'Public', 'Private', 'Not Indicated'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetJuniorHighStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN (
			SELECT s1.*
			FROM school_details s1
			JOIN (
				SELECT eb_id, educational_level_id,
					MAX(year_completed) AS max_yr
				FROM school_details
				GROUP BY eb_id, educational_level_id
			) s2 ON s1.eb_id = s2.eb_id
				AND s1.educational_level_id = s2.educational_level_id
				AND s1.year_completed = s2.max_yr
		) sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		WHERE el.level_name = 'Junior High School' ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category, 'Public', 'Private', 'Not Indicated'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetSeniorHighStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN (
			SELECT s1.*
			FROM school_details s1
			JOIN (
				SELECT eb_id, educational_level_id,
					MAX(year_completed) AS max_yr
				FROM school_details
				GROUP BY eb_id, educational_level_id
			) s2 ON s1.eb_id = s2.eb_id
				AND s1.educational_level_id = s2.educational_level_id
				AND s1.year_completed = s2.max_yr
		) sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		WHERE el.level_name = 'Senior High School' ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category, 'Public', 'Private', 'Not Indicated'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetHighSchoolStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN (
			SELECT s1.*
			FROM school_details s1
			JOIN (
				SELECT eb_id, educational_level_id,
					MAX(year_completed) AS max_yr
				FROM school_details
				GROUP BY eb_id, educational_level_id
			) s2 ON s1.eb_id = s2.eb_id
				AND s1.educational_level_id = s2.educational_level_id
				AND s1.year_completed = s2.max_yr
		) sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		WHERE el.level_name IN (
			'High School', 'Junior High School',
			'Senior High School'
		) ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category, 'Public', 'Private', 'Not Indicated'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetVocationalStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN (
			SELECT s1.*
			FROM school_details s1
			JOIN (
				SELECT eb_id, educational_level_id,
					MAX(year_completed) AS max_yr
				FROM school_details
				GROUP BY eb_id, educational_level_id
			) s2 ON s1.eb_id = s2.eb_id
				AND s1.educational_level_id = s2.educational_level_id
				AND s1.year_completed = s2.max_yr
		) sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		WHERE el.level_name = 'Vocational' ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category, 'Public', 'Private', 'Not Indicated'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetCollegeStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(sd.school_type, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		JOIN (
			SELECT s1.*
			FROM school_details s1
			JOIN (
				SELECT eb_id, educational_level_id,
					MAX(year_completed) AS max_yr
				FROM school_details
				GROUP BY eb_id, educational_level_id
			) s2 ON s1.eb_id = s2.eb_id
				AND s1.educational_level_id = s2.educational_level_id
				AND s1.year_completed = s2.max_yr
		) sd ON eb.id = sd.eb_id
		JOIN educational_levels el ON sd.educational_level_id = el.id
		WHERE el.level_name = 'College' ` + filter + `
		GROUP BY category
		ORDER BY FIELD(
			category, 'Public', 'Private', 'Not Indicated'
		);`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetNatureOfSchoolingStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			COALESCE(eb.nature_of_schooling, 'Not Indicated') AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END) as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END) as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN educational_backgrounds eb ON spi.iir_id = eb.iir_id
		WHERE 1=1 ` + filter + `
		GROUP BY category;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetFatherLifeStatusStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			CASE
				WHEN srp.is_living = 1 THEN 'Living'
				ELSE 'Deceased'
			END AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END)
				as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END)
				as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN student_related_persons srp ON spi.iir_id = srp.iir_id
		JOIN student_relationship_types srt ON srp.relationship_id = srt.id
		WHERE srt.relationship_name = 'Father' ` + filter + `
		GROUP BY category
		ORDER BY category ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

func (r *Repository) GetMotherLifeStatusStats(
	ctx context.Context, year int, programID int,
) ([]DemographicStat, error) {
	filter, args := r.buildFilter(year, programID)
	query := `
		SELECT
			CASE
				WHEN srp.is_living = 1 THEN 'Living'
				ELSE 'Deceased'
			END AS category,
			SUM(CASE WHEN spi.gender = 'Male' THEN 1 ELSE 0 END)
				as male_count,
			SUM(CASE WHEN spi.gender = 'Female' THEN 1 ELSE 0 END)
				as female_count,
			COUNT(*) as total,
			RANK() OVER (ORDER BY COUNT(*) DESC) +
			(COUNT(*) OVER (
				PARTITION BY COUNT(*)
			) - 1) / 2.0 as rank_pos
		FROM student_personal_info spi
		JOIN student_related_persons srp ON spi.iir_id = srp.iir_id
		JOIN student_relationship_types srt ON srp.relationship_id = srt.id
		WHERE srt.relationship_name = 'Mother' ` + filter + `
		GROUP BY category
		ORDER BY category ASC;`
	return r.executeStatQuery(ctx, query, args...)
}

// --- HELPERS ---

func (r *Repository) buildFilter(
	year int,
	programID int,
) (string, []interface{}) {
	filter := " AND spi.iir_id IN " +
		"(SELECT id FROM iir_records WHERE is_completed = 1)"
	args := []interface{}{}

	if year > 0 {
		filter += " AND spi.student_number LIKE ?"
		args = append(args, fmt.Sprintf("%d-%%", year))
	}
	if programID > 0 {
		filter += " AND spi.program_id = ?"
		args = append(args, programID)
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

func (r *Repository) GetProgram(
	ctx context.Context,
	programID int,
) (string, string, error) {
	var code, name string
	query := "SELECT code, program_name FROM programs WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, programID).Scan(&code, &name)
	return code, name, err
}
