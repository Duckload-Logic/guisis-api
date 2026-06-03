package students

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Gender represents the domain entity for student gender.
type Gender struct {
	ID   int    `db:"id"          json:"id"`
	Name string `db:"gender_name" json:"name"`
}

// ParentalStatusType represents the domain entity for parental status types.
type ParentalStatusType struct {
	ID   int    `db:"id"          json:"id"`
	Name string `db:"status_name" json:"name"`
}

// StudentSupportType represents the domain entity for student support types.
type StudentSupportType struct {
	ID   int    `db:"id"                json:"id"`
	Name string `db:"support_type_name" json:"name"`
}

// EnrollmentReason represents the domain entity for enrollment reasons.
type EnrollmentReason struct {
	ID   int    `db:"id"          json:"id"`
	Text string `db:"reason_text" json:"text"`
}

// IncomeRange represents the domain entity for family income ranges.
type IncomeRange struct {
	ID   int    `db:"id"         json:"id"`
	Text string `db:"range_text" json:"text"`
}

// EducationalLevel represents the domain entity for educational levels.
type EducationalLevel struct {
	ID   int    `db:"id"         json:"id"`
	Name string `db:"level_name" json:"name"`
}

// EducationalAttainment represents the standardized educational attainment levels.
type EducationalAttainment struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

// Course represents the domain entity for academic courses.
type Course struct {
	ID   int    `db:"id"          json:"id"`
	Code string `db:"code"        json:"code"`
	Name string `db:"course_name" json:"name"`
}

// CivilStatusType represents the domain entity for civil status types.
type CivilStatusType struct {
	ID   int    `db:"id"          json:"id"`
	Name string `db:"status_name" json:"name"`
}

// Religion represents the domain entity for religious affiliations.
type Religion struct {
	ID   int    `db:"id"            json:"id"`
	Name string `db:"religion_name" json:"name"`
}

// StudentRelationshipType represents the domain entity for relationships.
type StudentRelationshipType struct {
	ID   int    `db:"id"                json:"id"`
	Name string `db:"relationship_name" json:"name"`
}

// NatureOfResidenceType represents the domain entity for residence types.
type NatureOfResidenceType struct {
	ID   int    `db:"id"                  json:"id"`
	Name string `db:"residence_type_name" json:"name"`
}

// SibilingSupportType represents the domain entity for sibling support.
type SibilingSupportType struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

// StudentStatus represents the domain entity for academic status.
type StudentStatus struct {
	ID   int    `db:"id"          json:"id"`
	Name string `db:"status_name" json:"name"`
}

// ActivityOption represents the domain entity for student activities.
type ActivityOption struct {
	ID       int    `db:"id"        json:"id"`
	Name     string `db:"name"      json:"name"`
	Category string `db:"category"  json:"category"`
	IsActive bool   `db:"is_active" json:"isActive"`
}

// StudentBasicInfoView represents a simplified view of student identification info.
type StudentBasicInfoView struct {
	UserID     string                 `db:"user_id"`
	Email      string                 `db:"email"`
	FirstName  string                 `db:"first_name"`
	MiddleName structs.NullableString `db:"middle_name"`
	LastName   string                 `db:"last_name"`
	SuffixName structs.NullableString `db:"suffix_name"`
}

type StudentProfileView struct {
	IIRID         string                 `db:"iir_id"`
	UserID        string                 `db:"user_id"`
	FirstName     string                 `db:"first_name"`
	MiddleName    structs.NullableString `db:"middle_name"`
	LastName      string                 `db:"last_name"`
	SuffixName    structs.NullableString `db:"suffix_name"`
	Email         string                 `db:"email"`
	StudentNumber string                 `db:"student_number"`
	GenderID      int                    `db:"gender_id"`
	CourseID      int                    `db:"course_id"`
	Section       int                    `db:"section"`
	YearLevel     int                    `db:"year_level"`
	StatusID      int                    `db:"status_id"`
	StatusName    string                 `db:"status_name"`
	GenderName    string                 `db:"gender_name"`
	CourseCode    string                 `db:"course_code"`
	CourseName    string                 `db:"course_name"`
}

// IIRDraft represents an unsubmitted draft of the student's IIR form.
type IIRDraft struct {
	ID        int       `db:"id"         json:"id"`
	UserID    string    `db:"user_id"    json:"userId"`
	Data      string    `db:"data"       json:"data"` // JSON string
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// IIRRecord represents a submitted student record header.
type IIRRecord struct {
	ID          string    `db:"id"           json:"id"`
	UserID      string    `db:"user_id"      json:"userId"`
	IsSubmitted bool      `db:"is_submitted" json:"isSubmitted"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updatedAt"`
}

// StudentSelectedReason represents reasons why a student chose their course.
type StudentSelectedReason struct {
	IIRID           string                 `db:"iir_id"            json:"iirId"`
	ReasonID        int                    `db:"reason_id"         json:"reasonId"`
	OtherReasonText structs.NullableString `db:"other_reason_text" json:"otherReasonText,omitempty"`
}

// StudentPersonalInfo represents comprehensive personal and academic data for a student.
type StudentPersonalInfo struct {
	ID    int    `db:"id"     json:"id"`
	IIRID string `db:"iir_id" json:"iirId"`

	StudentNumber string `db:"student_number" json:"studentNumber"`

	// Personal Details
	GenderID      int `db:"gender_id"       json:"genderId"`
	CivilStatusID int `db:"civil_status_id" json:"civilStatusId"`
	ReligionID    int `db:"religion_id"     json:"religionId"`

	// Physical Attributes
	HeightM    float64 `db:"height_m"   json:"heightM"`
	WeightKg   float64 `db:"weight_kg"  json:"weightKg"`
	Complexion string  `db:"complexion" json:"complexion"`

	// Academic Information
	HighSchoolGWA float64 `db:"high_school_gwa" json:"highSchoolGWA"`
	CourseID      int     `db:"course_id"       json:"courseId"`
	YearLevel     int     `db:"year_level"      json:"yearLevel"`
	Section       int     `db:"section"         json:"section"`

	// Additional Details
	PlaceOfBirth          string                 `db:"place_of_birth"          json:"placeOfBirth"`
	DateOfBirth           string                 `db:"date_of_birth"           json:"dateOfBirth"`
	IsEmployed            bool                   `db:"is_employed"             json:"isEmployed"`
	EmployerName          structs.NullableString `db:"employer_name"           json:"employerName,omitempty"`
	EmployerAddress       structs.NullableString `db:"employer_address"        json:"employerAddress,omitempty"`
	MobileNumber          string                 `db:"mobile_number"           json:"mobileNumber"`
	TelephoneNumber       structs.NullableString `db:"telephone_number"        json:"telephoneNumber,omitempty"`
	EmployerContactNumber structs.NullableString `db:"employer_contact_number" json:"employerContactNumber,omitempty"`
	OtherReligionText     structs.NullableString `db:"other_religion_text"`

	StatusID       int                   `db:"status_id"       json:"statusId"`
	GraduationYear structs.NullableInt64 `db:"graduation_year" json:"graduationYear,omitempty"`
	CreatedAt      time.Time             `db:"created_at"      json:"createdAt"`
	UpdatedAt      time.Time             `db:"updated_at"      json:"updatedAt"`
}

// EmergencyContact represents the contact person in case of emergency.
type EmergencyContact struct {
	ID             int                    `db:"id"              json:"id"`
	IIRID          string                 `db:"iir_id"          json:"iirId"`
	FirstName      string                 `db:"first_name"      json:"firstName"`
	MiddleName     structs.NullableString `db:"middle_name"     json:"middleName,omitempty"`
	LastName       string                 `db:"last_name"       json:"lastName"`
	SuffixName     structs.NullableString `db:"suffix_name"     json:"suffixName,omitempty"`
	ContactNumber  string                 `db:"contact_number"  json:"contactNumber"`
	RelationshipID int                    `db:"relationship_id" json:"relationshipId"`
	AddressID      int                    `db:"address_id"      json:"addressId"`
	CreatedAt      time.Time              `db:"created_at"      json:"createdAt"`
	UpdatedAt      time.Time              `db:"updated_at"      json:"updatedAt"`
}

// StudentAddress links an IIR record to an address.
type StudentAddress struct {
	ID          int       `db:"id"           json:"id"`
	IIRID       string    `db:"iir_id"       json:"iirId"`
	AddressID   int       `db:"address_id"   json:"addressId"`
	AddressType string    `db:"address_type" json:"addressType"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updatedAt"`
}

// RelatedPerson represents family members or guardians.
type RelatedPerson struct {
	ID                      int                    `db:"id"                        json:"id"`
	EducationalAttainmentID structs.NullableInt64  `db:"educational_attainment_id" json:"educationalAttainmentId,omitempty"`
	DateOfBirth             string                 `db:"date_of_birth"             json:"dateOfBirth"`
	LastName                string                 `db:"last_name"                 json:"lastName"`
	FirstName               string                 `db:"first_name"                json:"firstName"`
	MiddleName              structs.NullableString `db:"middle_name"               json:"middleName,omitempty"`
	SuffixName              structs.NullableString `db:"suffix_name"               json:"suffixName,omitempty"`
	Occupation              structs.NullableString `db:"occupation"                json:"occupation,omitempty"`
	EmployerName            structs.NullableString `db:"employer_name"             json:"employerName,omitempty"`
	EmployerAddress         structs.NullableString `db:"employer_address"          json:"employerAddress,omitempty"`
	CreatedAt               time.Time              `db:"created_at"                json:"createdAt"`
	UpdatedAt               time.Time              `db:"updated_at"                json:"updatedAt"`
}

// StudentRelatedPerson represents the link between a student and a related person.
type StudentRelatedPerson struct {
	IIRID           string    `db:"iir_id"            json:"iirId"`
	RelatedPersonID int       `db:"related_person_id" json:"relatedPersonId"`
	RelationshipID  int       `db:"relationship_id"   json:"relationshipId"`
	IsParent        bool      `db:"is_parent"         json:"isParent"`
	IsGuardian      bool      `db:"is_guardian"       json:"isGuardian"`
	IsLiving        bool      `db:"is_living"         json:"isLiving"`
	CreatedAt       time.Time `db:"created_at"        json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at"        json:"updatedAt"`
}

// FamilyBackground represents family dynamics and housing information.
type FamilyBackground struct {
	ID                    int                    `db:"id"                        json:"id"`
	IIRID                 string                 `db:"iir_id"                    json:"iirId"`
	ParentalStatusID      int                    `db:"parental_status_id"        json:"parentalStatusId"`
	ParentalStatusDetails structs.NullableString `db:"parental_status_details"   json:"parentalStatusDetails,omitempty"`
	Brothers              int                    `db:"brothers"                  json:"brothers"`
	Sisters               int                    `db:"sisters"                   json:"sisters"`
	EmployedSiblings      int                    `db:"employed_siblings"         json:"employedSiblings"`
	OrdinalPosition       int                    `db:"ordinal_position"          json:"ordinalPosition"`
	HaveQuietPlaceToStudy bool                   `db:"have_quiet_place_to_study" json:"haveQuietPlaceToStudy"`
	IsSharingRoom         bool                   `db:"is_sharing_room"           json:"isSharingRoom"`
	RoomSharingDetails    structs.NullableString `db:"room_sharing_details"      json:"roomSharingDetails,omitempty"`
	NatureOfResidenceId   int                    `db:"nature_of_residence_id"    json:"natureOfResidenceId"`
	CreatedAt             time.Time              `db:"created_at"                json:"createdAt"`
	UpdatedAt             time.Time              `db:"updated_at"                json:"updatedAt"`
}

// StudentSiblingSupport represents financial support for siblings.
type StudentSiblingSupport struct {
	FamilyBackgroundID int `db:"family_background_id" json:"familyBackgroundId"`
	SupportTypeID      int `db:"support_type_id"      json:"supportTypeId"`
}

// EducationalBackground represents overall educational history.
type EducationalBackground struct {
	ID                 int                    `db:"id"                  json:"id"`
	IIRID              string                 `db:"iir_id"              json:"iirId"`
	NatureOfSchooling  string                 `db:"nature_of_schooling" json:"natureOfSchooling"`
	InterruptedDetails structs.NullableString `db:"interrupted_details" json:"interruptedDetails,omitempty"`
	CreatedAt          time.Time              `db:"created_at"          json:"createdAt"`
	UpdatedAt          time.Time              `db:"updated_at"          json:"updatedAt"`
}

// SchoolDetails represents specific school levels and details.
type SchoolDetails struct {
	ID                 int                    `db:"id"                   json:"id"`
	EBID               int                    `db:"eb_id"                json:"ebId"`
	EducationalLevelID int                    `db:"educational_level_id" json:"educationalLevelId"`
	SchoolName         string                 `db:"school_name"          json:"schoolName"`
	SchoolAddress      string                 `db:"school_address"       json:"schoolAddress"`
	SchoolType         string                 `db:"school_type"          json:"schoolType"`
	YearStarted        int                    `db:"year_started"         json:"yearStarted"`
	YearCompleted      int                    `db:"year_completed"       json:"yearCompleted"`
	Awards             structs.NullableString `db:"awards"               json:"awards,omitempty"`
	CreatedAt          time.Time              `db:"created_at"           json:"createdAt"`
	UpdatedAt          time.Time              `db:"updated_at"           json:"updatedAt"`
}

// StudentHealthRecord represents physical health status.
type StudentHealthRecord struct {
	ID                        int                    `db:"id"                           json:"id"`
	IIRID                     string                 `db:"iir_id"                       json:"iirId"`
	VisionHasProblem          bool                   `db:"vision_has_problem"           json:"visionHasProblem"`
	VisionDetails             structs.NullableString `db:"vision_details"               json:"visionDetails,omitempty"`
	HearingHasProblem         bool                   `db:"hearing_has_problem"          json:"hearingHasProblem"`
	HearingDetails            structs.NullableString `db:"hearing_details"              json:"hearingDetails,omitempty"`
	SpeechHasProblem          bool                   `db:"speech_has_problem"           json:"speechHasProblem"`
	SpeechDetails             structs.NullableString `db:"speech_details"               json:"speechDetails,omitempty"`
	GeneralHealthHasProblem   bool                   `db:"general_health_has_problem"   json:"generalHealthHasProblem"`
	GeneralHealthDetails      structs.NullableString `db:"general_health_details"       json:"generalHealthDetails,omitempty"`
	MentalEmotionalHasProblem bool                   `db:"mental_emotional_has_problem" json:"mentalEmotionalHasProblem"`
	MentalEmotionalDetails    structs.NullableString `db:"mental_emotional_details"     json:"mentalEmotionalDetails,omitempty"`
	CreatedAt                 time.Time              `db:"created_at"                   json:"createdAt"`
	UpdatedAt                 time.Time              `db:"updated_at"                   json:"updatedAt"`
}

// StudentConsultation represents history with professional counselors.
type StudentConsultation struct {
	ID               int                    `db:"id"                json:"id"`
	IIRID            string                 `db:"iir_id"            json:"iirId"`
	ProfessionalType string                 `db:"professional_type" json:"professionalType"`
	HasConsulted     bool                   `db:"has_consulted"     json:"hasConsulted"`
	WhenDate         structs.NullableString `db:"when_date"         json:"whenDate,omitempty"`
	ForWhat          structs.NullableString `db:"for_what"          json:"forWhat,omitempty"`
	CreatedAt        time.Time              `db:"created_at"        json:"createdAt"`
	UpdatedAt        time.Time              `db:"updated_at"        json:"updatedAt"`
}

// StudentFinance represents family financial status.
type StudentFinance struct {
	ID              int                    `db:"id"                             json:"id"`
	IIRID           string                 `db:"iir_id"                         json:"iirId"`
	IncomeRangeID   int                    `db:"monthly_family_income_range_id" json:"incomeRangeId"`
	OtherIncome     structs.NullableString `db:"other_income_details"           json:"otherIncome,omitempty"`
	WeeklyAllowance float64                `db:"weekly_allowance"               json:"weeklyAllowance"`
	CreatedAt       time.Time              `db:"created_at"                     json:"createdAt"`
	UpdatedAt       time.Time              `db:"updated_at"                     json:"updatedAt"`
}

// StudentFinancialSupport represents sources of financial aid.
type StudentFinancialSupport struct {
	StudentFinanceID int       `db:"sf_id"           json:"sfId"`
	SupportTypeID    int       `db:"support_type_id" json:"supportTypeId"`
	CreatedAt        time.Time `db:"created_at"      json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at"      json:"updatedAt"`
}

// StudentActivity represents participation in clubs or sports.
type StudentActivity struct {
	ID                 int                    `db:"id"                  json:"id"`
	IIRID              string                 `db:"iir_id"              json:"iirId"`
	OptionID           int                    `db:"option_id"           json:"optionId"`
	OtherSpecification structs.NullableString `db:"other_specification" json:"otherSpecification,omitempty"`
	Role               string                 `db:"role"                json:"role"`
	RoleSpecification  structs.NullableString `db:"role_specification"  json:"roleSpecification,omitempty"`
	CreatedAt          time.Time              `db:"created_at"          json:"createdAt"`
	UpdatedAt          time.Time              `db:"updated_at"          json:"updatedAt"`
}

// StudentSubjectPreference represents favorite or difficult subjects.
type StudentSubjectPreference struct {
	ID          int       `db:"id"           json:"id"`
	IIRID       string    `db:"iir_id"       json:"iirId"`
	SubjectName string    `db:"subject_name" json:"subjectName"`
	IsFavorite  bool      `db:"is_favorite"  json:"isFavorite"`
	CreatedAt   time.Time `db:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updatedAt"`
}

// StudentHobby represents student hobbies and priority.
type StudentHobby struct {
	ID           int       `db:"id"            json:"id"`
	IIRID        string    `db:"iir_id"        json:"iirId"`
	HobbyName    string    `db:"hobby_name"    json:"hobbyName"`
	PriorityRank int       `db:"priority_rank" json:"priorityRank"`
	CreatedAt    time.Time `db:"created_at"    json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updatedAt"`
}

// TestResult represents psychological or academic test results.
type TestResult struct {
	ID          int       `db:"id"          json:"id"`
	IIRID       string    `db:"iir_id"      json:"iirId"`
	TestDate    string    `db:"test_date"   json:"testDate"`
	TestName    string    `db:"test_name"   json:"testName"`
	RawScore    string    `db:"raw_score"   json:"rawScore"`
	Percentile  string    `db:"percentile"  json:"percentile"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at"  json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updatedAt"`
}

// StudentCOR represents the Certificate of Registration file link.
type StudentCOR struct {
	FileID        string               `db:"file_id"        json:"fileId"`
	StudentID     string               `db:"student_id"     json:"studentId"`
	StudentNumber string               `db:"student_number" json:"studentNumber"`
	CourseCode    string               `db:"course_code"    json:"courseCode"`
	CourseDesc    string               `db:"course_desc"    json:"courseDesc"`
	YearLevel     int                  `db:"year_level"     json:"yearLevel"`
	Section       int                  `db:"section"        json:"section"`
	Campus        string               `db:"campus"         json:"campus"`
	YearStart     int                  `db:"year_start"     json:"yearStart"`
	YearEnd       int                  `db:"year_end"       json:"yearEnd"`
	Term          int                  `db:"term"           json:"term"`
	ValidFrom     structs.NullableTime `db:"valid_from"     json:"validFrom,omitempty"`
	ValidUntil    structs.NullableTime `db:"valid_until"    json:"validUntil,omitempty"`
}

// AcademicSetting holds the single-row global academic year + term setting
// managed by the SuperAdmin. Used to validate uploaded student CORs.
type AcademicSetting struct {
	ID               int       `db:"id"                 json:"id"`
	CurrentYearStart int       `db:"current_year_start" json:"currentYearStart"`
	CurrentYearEnd   int       `db:"current_year_end"   json:"currentYearEnd"`
	CurrentTerm      int       `db:"current_term"       json:"currentTerm"`
	UpdatedAt        time.Time `db:"updated_at"         json:"updatedAt"`
}

// StudentPersonalInfoView holds joined personal info and emergency details.
type StudentPersonalInfoView struct {
	ID                    int                    `db:"id"`
	IIRID                 string                 `db:"iir_id"`
	StudentNumber         string                 `db:"student_number"`
	GenderID              int                    `db:"gender_id"`
	GenderName            string                 `db:"gender_name"`
	CivilStatusID         int                    `db:"civil_status_id"`
	CivilStatusName       string                 `db:"civil_status_name"`
	ReligionID            int                    `db:"religion_id"`
	ReligionName          string                 `db:"religion_name"`
	HeightM               float64                `db:"height_m"`
	WeightKg              float64                `db:"weight_kg"`
	Complexion            string                 `db:"complexion"`
	HighSchoolGWA         float64                `db:"high_school_gwa"`
	CourseID              int                    `db:"course_id"`
	CourseCode            string                 `db:"course_code"`
	CourseName            string                 `db:"course_name"`
	YearLevel             int                    `db:"year_level"`
	Section               int                    `db:"section"`
	PlaceOfBirth          string                 `db:"place_of_birth"`
	DateOfBirth           string                 `db:"date_of_birth"`
	IsEmployed            bool                   `db:"is_employed"`
	EmployerName          structs.NullableString `db:"employer_name"`
	EmployerAddress       structs.NullableString `db:"employer_address"`
	MobileNumber          string                 `db:"mobile_number"`
	TelephoneNumber       structs.NullableString `db:"telephone_number"`
	EmployerContactNumber structs.NullableString `db:"employer_contact_number"`
	OtherReligionText     structs.NullableString `db:"other_religion_text"`
	StatusID              int                    `db:"status_id"`
	StatusName            string                 `db:"status_name"`
	GraduationYear        structs.NullableInt64  `db:"graduation_year"`

	// Emergency Contact Details
	EmergencyID               int                    `db:"emergency_id"`
	EmergencyFirstName        string                 `db:"emergency_first_name"`
	EmergencyMiddleName       structs.NullableString `db:"emergency_middle_name"`
	EmergencyLastName         string                 `db:"emergency_last_name"`
	EmergencyContactNumber    string                 `db:"emergency_contact_number"`
	EmergencyRelationshipID   int                    `db:"emergency_relationship_id"`
	EmergencyRelationshipName string                 `db:"emergency_relationship_name"`
	EmergencyAddressID        int                    `db:"emergency_address_id"`
}

// RelatedPersonView holds joined related person and relationship details.
type RelatedPersonView struct {
	ID                        int                    `db:"id"`
	LastName                  string                 `db:"last_name"`
	FirstName                 string                 `db:"first_name"`
	MiddleName                structs.NullableString `db:"middle_name"`
	SuffixName                structs.NullableString `db:"suffix_name"`
	DateOfBirth               string                 `db:"date_of_birth"`
	EducationalAttainmentID   int                    `db:"educational_attainment_id"`
	EducationalAttainmentName string                 `db:"educational_attainment_name"`
	Occupation                structs.NullableString `db:"occupation"`
	EmployerName              structs.NullableString `db:"employer_name"`
	EmployerAddress           structs.NullableString `db:"employer_address"`
	RelationshipID            int                    `db:"relationship_id"`
	RelationshipName          string                 `db:"relationship_name"`
	IsParent                  bool                   `db:"is_parent"`
	IsGuardian                bool                   `db:"is_guardian"`
	IsLiving                  bool                   `db:"is_living"`
}

// StudentFinanceView holds joined student finance and income range details.
type StudentFinanceView struct {
	ID              int                    `db:"id"`
	IIRID           string                 `db:"iir_id"`
	IncomeRangeID   int                    `db:"income_range_id"`
	IncomeRangeText string                 `db:"income_range_text"`
	OtherIncome     structs.NullableString `db:"other_income"`
	WeeklyAllowance float64                `db:"weekly_allowance"`
}
// SignificantNote represents a note record from the database.
type SignificantNote struct {
	ID              string                 `db:"id"`
	IIRID           structs.NullableString `db:"iir_id"`
	AppointmentID   structs.NullableString `db:"appointment_id"`
	AdmissionSlipID structs.NullableString `db:"admission_slip_id"`
	Note            string                 `db:"note"`
	Remarks         string                 `db:"remarks"`
	CreatedAt       time.Time              `db:"created_at"`
	UpdatedAt       time.Time              `db:"updated_at"`
}
