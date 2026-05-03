package analytics

import "time"

// DemographicStat represents demographic statistics for analytics.
type DemographicStat struct {
	Category    string `db:"category"     json:"category"`
	MaleCount   int    `db:"male_count"   json:"maleCount"`
	FemaleCount int    `db:"female_count" json:"femaleCount"`
	Total       int    `db:"total"        json:"total"`
	RankPos     int    `db:"rank_pos"     json:"rankPos"`
}

// StudentPersonalInfo represents personal information data for analytics.
type StudentPersonalInfo struct {
	ID              int       `db:"id"               json:"id"`
	IIRID           string    `db:"iir_id"           json:"iirId"`
	StudentNumber   string    `db:"student_number"   json:"studentNumber"`
	GenderID        int       `db:"gender_id"        json:"genderId"`
	CivilStatusID   int       `db:"civil_status_id"  json:"civilStatusId"`
	ReligionID      int       `db:"religion_id"      json:"religionId"`
	HeightFt        float64   `db:"height_ft"        json:"heightFt"`
	WeightKg        float64   `db:"weight_kg"        json:"weightKg"`
	Complexion      string    `db:"complexion"       json:"complexion"`
	HighSchoolGWA   float64   `db:"high_school_gwa"  json:"highSchoolGwa"`
	CourseID        int       `db:"course_id"        json:"courseId"`
	YearLevel       int       `db:"year_level"       json:"yearLevel"`
	Section         int       `db:"section"          json:"section"`
	PlaceOfBirth    string    `db:"place_of_birth"   json:"placeOfBirth"`
	DateOfBirth     time.Time `db:"date_of_birth"    json:"dateOfBirth"`
	IsEmployed      bool      `db:"is_employed"      json:"isEmployed"`
	EmployerName    *string   `db:"employer_name"    json:"employerName"`
	EmployerAddress *string   `db:"employer_address" json:"employerAddress"`
	MobileNumber    string    `db:"mobile_number"    json:"mobileNumber"`
	CreatedAt       time.Time `db:"created_at"       json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at"       json:"updatedAt"`
}

// StudentFinances represents financial data for student analytics.
type StudentFinances struct {
	ID                         int     `db:"id"                             json:"id"`
	IIRID                      string  `db:"iir_id"                         json:"iirId"`
	MonthlyFamilyIncomeRangeID *int    `db:"monthly_family_income_range_id" json:"monthlyFamilyIncomeRangeId"`
	WeeklyAllowance            float64 `db:"weekly_allowance"               json:"weeklyAllowance"`
}

// FamilyBackground represents family background data for student analytics.
type FamilyBackground struct {
	IIRID                 string `db:"iir_id"                    json:"iirId"`
	ParentalStatusID      int    `db:"parental_status_id"        json:"parentalStatusId"`
	OrdinalPosition       int    `db:"ordinal_position"          json:"ordinalPosition"`
	HaveQuietPlaceToStudy bool   `db:"have_quiet_place_to_study" json:"haveQuietPlaceToStudy"`
}

// EducationalBackground represents educational background data for student analytics.
type EducationalBackground struct {
	ID               int     `db:"id"                 json:"id"`
	IIRID            string  `db:"iir_id"             json:"iirId"`
	EducationLevelID int     `db:"education_level_id" json:"educationLevelId"`
	SchoolDetailID   int     `db:"school_detail_id"   json:"schoolDetailId"`
	GWA              float64 `db:"gwa"                json:"gwa"`
}
