package integrations

import (
	"database/sql"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type OGOSStudentView struct {
	IDPUUID       sql.NullString `db:"idp_uuid"`
	StudentNumber string         `db:"student_number"`

	FirstName  string         `db:"first_name"`
	MiddleName sql.NullString `db:"middle_name,omitempty"`
	LastName   string         `db:"last_name"`
	SuffixName sql.NullString `db:"suffix_name,omitempty"`

	Email        string `db:"email"`
	MobileNumber string `db:"mobile_number"`

	ProgramID   int    `db:"program_id"`
	ProgramCode string `db:"program_code"`
	ProgramName string `db:"program_name"`
	YearLevel   int    `db:"year_level"`
	Section     string `db:"section"`
}

type OGOSStudentPersonalInfoView struct {
	IDPUUID       sql.NullString `db:"idp_uuid"`
	StudentNumber string         `db:"student_number"`

	GenderID                     int     `db:"gender_id"`
	GenderName                   string  `db:"gender_name"`
	DateOfBirth                  string  `db:"date_of_birth"`
	PlaceOfBirth                 string  `db:"place_of_birth"`
	HeightM                      float32 `db:"height_m"`
	WeightKg                     float32 `db:"weight_kg"`
	EmergencyContactName         string  `db:"emergency_contact_name"`
	EmergencyContactNumber       string  `db:"emergency_contact_number"`
	EmergencyContactRelationship string  `db:"emergency_contact_relationship"`
}

type OGOSStudentAddressView struct {
	IDPUUID       sql.NullString `db:"idp_uuid"`
	StudentNumber string         `db:"student_number"`

	AddressType  string                 `db:"address_type,omitempty"`
	StreetDetail structs.NullableString `db:"street_detail"`
	BarangayCode string                 `db:"barangay_code"`
	BarangayName string                 `db:"barangay_name"`
	CityCode     string                 `db:"city_code"`
	CityName     string                 `db:"city_name"`
	ProvinceCode sql.NullString         `db:"province_code,omitempty"`
	ProvinceName sql.NullString         `db:"province_name,omitempty"`
	RegionCode   string                 `db:"region_code,omitempty"`
	RegionName   string                 `db:"region_name,omitempty"`
}
