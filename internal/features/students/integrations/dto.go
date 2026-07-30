package integrations

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
)

type OGOSLinkStudentRequest struct {
	VerificationCode string `json:"verificationCode"`
}

type OGOSListStudentsRequest struct {
	structs.PaginationRequest
	ProgramID int `form:"program_id,omitempty"`
	GenderID  int `form:"gender_id,omitempty"`
	YearLevel int `form:"year_level,omitempty"`
}

type OGOSListStudentsResponse struct {
	Students []OGOSStudentDTO           `json:"students"`
	Meta     structs.PaginationMetadata `json:"meta"`
}

type OGOSStudentDTO struct {
	StudentNumber string `json:"studentNumber"`

	FirstName  string                 `json:"firstName"`
	MiddleName structs.NullableString `json:"middleName,omitempty"`
	LastName   string                 `json:"lastName"`
	SuffixName string                 `json:"suffixName,omitempty"`

	Email        string `json:"email"`
	MobileNumber string `json:"mobileNumber"`

	Program   students.Program `json:"program"`
	YearLevel int             `json:"yearLevel"`
	Section   string          `json:"section"`
}

type OGOSStudentPersonalInfoDTO struct {
	StudentNumber string `json:"studentNumber"`

	Gender                 students.Gender `json:"gender"`
	DateOfBirth            string          `json:"dateOfBirth"`
	PlaceOfBirth           string          `json:"placeOfBirth"`
	HeightM                float32         `json:"heightM"`
	WeightKg               float32         `json:"weightKg"`
	EmergencyContactName   string          `json:"emergencyContactName"`
	EmergencyContactNumber string          `json:"emergencyContactNumber"`
}

type OGOSStudentAddressDTO struct {
	StudentNumber string `json:"studentNumber"`

	AddressType  string                 `json:"addressType,omitempty"`
	StreetDetail structs.NullableString `json:"streetDetail"`
	Barangay     locations.Barangay     `json:"barangay"`
	City         locations.City         `json:"city"`
	Province     *locations.ProvinceDTO `json:"province,omitempty"`
	Region       locations.Region       `json:"region,omitempty"`
}
