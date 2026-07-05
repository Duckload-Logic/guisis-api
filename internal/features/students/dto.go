package students

import (
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
)

// List Students
type ListStudentsRequest struct {
	structs.PaginationRequest
	CourseID  int `form:"course_id,omitempty"`
	GenderID  int `form:"gender_id,omitempty"`
	YearLevel int `form:"year_level,omitempty"`
	StatusID  int `form:"status_id,omitempty"`
}

type ListStudentsResponse struct {
	Students []StudentProfileDTO        `json:"students"`
	Meta     structs.PaginationMetadata `json:"meta"`
}

type StudentProfileDTO struct {
	IIRID             string                 `json:"iirId"`
	UserID            string                 `json:"userId"`
	FirstName         string                 `json:"firstName"`
	MiddleName        structs.NullableString `json:"middleName,omitempty"`
	LastName          string                 `json:"lastName"`
	SuffixName        structs.NullableString `json:"suffixName,omitempty"`
	Gender            Gender                 `json:"gender"`
	Email             string                 `json:"email"`
	StudentNumber     string                 `json:"studentNumber"`
	Course            Course                 `json:"course"`
	Section           int                    `json:"section"`
	YearLevel         int                    `json:"yearLevel"`
	Status            StudentStatus          `json:"status"`
	StudentCORURL     string                 `json:"studentCorUrl,omitempty"`
	IsStudentCORValid bool                   `json:"isStudentCorValid"`
	ProfilePicture    string                 `json:"profilePicture,omitempty"`
}

type ComprehensiveProfileDTO struct {
	IIRID             string `json:"iirId,omitempty"`
	StudentCORURL     string `json:"studentCorUrl,omitempty"`
	IsStudentCORValid bool   `json:"isStudentCorValid"`
	Student           struct {
		BasicInfo              StudentBasicInfoViewDTO `json:"basicInfo"`
		StudentPersonalInfoDTO `json:"personalInfo"`
		Addresses              []StudentAddressDTO `json:"addresses"`
	} `json:"student"`

	Education EducationalBackgroundDTO `json:"education"`

	Family struct {
		FamilyBackgroundDTO `json:"background"`
		RelatedPersons      []RelatedPersonDTO `json:"relatedPersons"`
		Finance             StudentFinanceDTO  `json:"finance"`
	} `json:"family"`

	Health struct {
		StudentHealthRecordDTO `json:"healthRecord"`
		Consultations          []StudentConsultationDTO `json:"consultations"`
	} `json:"health"`

	Interests struct {
		Activities         []StudentActivityDTO          `json:"activities"`
		SubjectPreferences []StudentSubjectPreferenceDTO `json:"subjectPreferences"`
		Hobbies            []StudentHobbyDTO             `json:"hobbies"`
	} `json:"interests"`

	// TestResults []TestResultDTO `json:"testResults"`
}

type StudentSelectedReasonDTO struct {
	Reason          EnrollmentReason `json:"reason"`
	OtherReasonText *string          `json:"otherReasonText,omitempty"`
}

type StudentBasicInfoViewDTO struct {
	Email      string                 `json:"email"`
	FirstName  string                 `json:"firstName"`
	MiddleName structs.NullableString `json:"middleName,omitempty"`
	LastName   string                 `json:"lastName"`
	SuffixName structs.NullableString `json:"suffixName,omitempty"`
}

type StudentPersonalInfoDTO struct {
	ID                    int                    `json:"id,omitempty"`
	IIRID                 string                 `json:"iirId,omitempty"`
	StudentNumber         string                 `json:"studentNumber"                   binding:"required"`
	Gender                Gender                 `json:"gender"                          binding:"required"`
	CivilStatus           CivilStatusType        `json:"civilStatus"                     binding:"required"`
	Religion              Religion               `json:"religion"                        binding:"required"`
	HeightM               float64                `json:"heightM"                         binding:"required"`
	WeightKg              float64                `json:"weightKg"                        binding:"required"`
	Complexion            string                 `json:"complexion"                      binding:"required"`
	HighSchoolGWA         float64                `json:"highSchoolGWA"                   binding:"required"`
	Course                Course                 `json:"course"                          binding:"required"`
	YearLevel             int                    `json:"yearLevel"                       binding:"required"`
	Section               int                    `json:"section"                         binding:"required"`
	PlaceOfBirth          string                 `json:"placeOfBirth"                    binding:"required"`
	DateOfBirth           string                 `json:"dateOfBirth"                     binding:"required"`
	IsEmployed            bool                   `json:"isEmployed"`
	EmployerName          structs.NullableString `json:"employerName,omitempty"`
	EmployerAddress       structs.NullableString `json:"employerAddress,omitempty"`
	MobileNumber          string                 `json:"mobileNumber"                    binding:"required"`
	TelephoneNumber       structs.NullableString `json:"telephoneNumber,omitempty"`
	EmployerContactNumber structs.NullableString `json:"employerContactNumber,omitempty"`
	OtherReligionText     structs.NullableString `json:"otherReligionText"`
	TwoByTwoPhotoDataUrl  string                 `json:"twoByTwoPhotoDataUrl,omitempty"`
	Status                StudentStatus          `json:"status"                          binding:"required"`
	GraduationYear        *int                   `json:"graduationYear,omitempty"`
	EmergencyContact      EmergencyContactDTO    `json:"emergencyContact,omitempty"`
}

type BulkUpdateStatusRequest struct {
	IIRIDs            []string `json:"iirIds"`
	ExcludedIIRIDs    []string `json:"excludedIirIds"`
	SelectAllMatching bool     `json:"selectAllMatching"`
	StatusID          int      `json:"statusId"                 binding:"required"`
	GraduationYear    *int     `json:"graduationYear,omitempty"`
	Filters           struct {
		Search     string `json:"search"`
		CourseID   int    `json:"courseId"`
		YearLevel  int    `form:"yearLevel"`
		EnrollYear int    `json:"enrollYear"`
	} `json:"filters"`
}

type EmergencyContactDTO struct {
	ID            int                     `json:"id,omitempty"`
	FirstName     string                  `json:"firstName"            binding:"required"`
	MiddleName    structs.NullableString  `json:"middleName,omitempty"`
	LastName      string                  `json:"lastName"             binding:"required"`
	SuffixName    structs.NullableString  `json:"suffixName,omitempty"`
	ContactNumber string                  `json:"contactNumber"        binding:"required"`
	Relationship  StudentRelationshipType `json:"relationship"         binding:"required"`
	Address       locations.AddressDTO    `json:"address"              binding:"required"`
}

type StudentAddressDTO struct {
	ID          int                  `json:"id,omitempty"`
	AddressType string               `json:"addressType"         binding:"required"`
	Address     locations.AddressDTO `json:"address"             binding:"required"`
	CreatedAt   time.Time            `json:"createdAt,omitempty"`
	UpdatedAt   time.Time            `json:"updatedAt,omitempty"`
}

type EducationalBackgroundDTO struct {
	ID                 int                    `json:"id,omitempty"`
	NatureOfSchooling  string                 `json:"natureOfSchooling"            binding:"required"`
	InterruptedDetails structs.NullableString `json:"interruptedDetails,omitempty"`
	School             []SchoolDetailsDTO     `json:"schools"                      binding:"required"`
	CreatedAt          time.Time              `json:"createdAt,omitempty"`
	UpdatedAt          time.Time              `json:"updatedAt,omitempty"`
}

type SchoolDetailsDTO struct {
	ID               int                    `json:"id,omitempty"`
	EducationalLevel EducationalLevel       `json:"educationalLevel"        binding:"required"`
	SchoolName       string                 `json:"schoolName"              binding:"required"`
	SchoolAddress    string                 `json:"schoolAddress,omitempty"`
	SchoolType       string                 `json:"schoolType"              binding:"required"`
	YearStarted      int                    `json:"yearStarted,omitempty"`
	YearCompleted    int                    `json:"yearCompleted"           binding:"required"`
	Awards           structs.NullableString `json:"awards,omitempty"`
}

type RelType = StudentRelationshipType
type EduAttain = EducationalAttainment

type RelatedPersonDTO struct {
	ID int `json:"id,omitempty"`

	LastName string `json:"lastName" binding:"omitempty"`

	FirstName string `json:"firstName" binding:"omitempty"`

	MiddleName structs.NullableString `json:"middleName,omitempty"`

	SuffixName structs.NullableString `json:"suffixName,omitempty"`

	DateOfBirth string `json:"dateOfBirth,omitempty" binding:"omitempty"`

	EducationalAttainment EduAttain `json:"educationalAttainment"`

	Occupation structs.NullableString `json:"occupation,omitempty"`

	EmployerName structs.NullableString `json:"employerName,omitempty"`

	EmployerAddress structs.NullableString `json:"employerAddress,omitempty"`

	Relationship RelType `json:"relationship" binding:"omitempty"`

	IsParent bool `json:"isParent"`

	IsGuardian bool `json:"isGuardian"`

	IsLiving bool `json:"isLiving"`
}

type FamilyBackgroundDTO struct {
	ID                    int                    `json:"id,omitempty"`
	ParentalStatus        ParentalStatusType     `json:"parentalStatus"                  binding:"required"`
	ParentalStatusDetails structs.NullableString `json:"parentalStatusDetails,omitempty"`
	Brothers              *int                   `json:"brothers"                        binding:"required"`
	Sisters               *int                   `json:"sisters"                         binding:"required"`
	EmployedSiblings      *int                   `json:"employedSiblings"                binding:"required"`
	OrdinalPosition       int                    `json:"ordinalPosition"                 binding:"required"`
	HaveQuietPlaceToStudy bool                   `json:"haveQuietPlaceToStudy"`
	SiblingSupportTypes   []SibilingSupportType  `json:"siblingSupportTypes"`
	IsSharingRoom         bool                   `json:"isSharingRoom"`
	RoomSharingDetails    structs.NullableString `json:"roomSharingDetails,omitempty"`
	NatureOfResidence     NatureOfResidenceType  `json:"natureOfResidence"               binding:"required"`
}

type EducationalBGDTO struct {
	ID               int    `json:"id"`
	EducationalLevel string `json:"educationalLevel"   binding:"required"`
	SchoolName       string `json:"schoolName"         binding:"required"`
	Location         string `json:"location,omitempty"`
	SchoolType       string `json:"schoolType"         binding:"required,oneof=Public Private"`
	YearCompleted    string `json:"yearCompleted"      binding:"required"`
	Awards           string `json:"awards,omitempty"`
}

type StudentFinanceDTO struct {
	ID                       int                    `json:"id,omitempty"`
	MonthlyFamilyIncomeRange IncomeRange            `json:"monthlyFamilyIncomeRange"        binding:"required"`
	OtherIncomeDetails       structs.NullableString `json:"otherIncomeDetails,omitempty"`
	FinancialSupportTypes    []StudentSupportType   `json:"financialSupportTypes,omitempty"`
	WeeklyAllowance          float64                `json:"weeklyAllowance"                 binding:"required"`
}

type StudentHealthRecordDTO struct {
	ID                        int                    `json:"id,omitempty"`
	VisionHasProblem          bool                   `json:"visionHasProblem"`
	VisionDetails             structs.NullableString `json:"visionDetails,omitempty"`
	HearingHasProblem         bool                   `json:"hearingHasProblem"`
	HearingDetails            structs.NullableString `json:"hearingDetails,omitempty"`
	SpeechHasProblem          bool                   `json:"speechHasProblem"`
	SpeechDetails             structs.NullableString `json:"speechDetails,omitempty"`
	GeneralHealthHasProblem   bool                   `json:"generalHealthHasProblem"`
	GeneralHealthDetails      structs.NullableString `json:"generalHealthDetails,omitempty"`
	MentalEmotionalHasProblem bool                   `json:"mentalEmotionalHasProblem"`
	MentalEmotionalDetails    structs.NullableString `json:"mentalEmotionalDetails,omitempty"`
}

type StudentConsultationDTO struct {
	ID               int                    `json:"id,omitempty"`
	ProfessionalType string                 `json:"professionalType"   binding:"required"`
	HasConsulted     bool                   `json:"hasConsulted"`
	WhenDate         structs.NullableString `json:"whenDate,omitempty"`
	ForWhat          structs.NullableString `json:"forWhat,omitempty"`
}

type StudentActivityDTO struct {
	ID                 int                    `json:"id,omitempty"`
	ActivityOption     ActivityOption         `json:"activityOption"`
	OtherSpecification structs.NullableString `json:"otherSpecification,omitempty"`
	Role               string                 `json:"role"` // "Officer", "Member", "Other"
	RoleSpecification  structs.NullableString `json:"roleSpecification,omitempty"`
}

type StudentSubjectPreferenceDTO struct {
	ID          int    `json:"id,omitempty"`
	SubjectName string `json:"subjectName"`
	IsFavorite  bool   `json:"isFavorite"`
}

type StudentHobbyDTO struct {
	ID           int    `json:"id,omitempty"`
	HobbyName    string `json:"hobbyName"`
	PriorityRank int    `json:"priorityRank"`
}

type TestResultDTO struct {
	ID          int    `json:"id,omitempty"`
	TestDate    string `json:"testDate"              binding:"required"`
	TestName    string `json:"testName"              binding:"required"`
	RawScore    string `json:"rawScore"              binding:"required"`
	Percentile  string `json:"percentile"            binding:"required"`
	Description string `json:"description,omitempty"`
}

// UpdateAcademicSettingDTO is the request body for the SuperAdmin-only
// PUT /students/settings/academic endpoint.
type UpdateAcademicSettingDTO struct {
	CurrentYearStart int `json:"currentYearStart" binding:"required,min=1900,max=2100"`
	CurrentYearEnd   int `json:"currentYearEnd"   binding:"required,min=1900,max=2100"`
	CurrentTerm      int `json:"currentTerm"      binding:"required,min=1,max=3"`
}

// ============================================================================
// DTO Mappings / ToDTO Receivers
// ============================================================================

// ToDTO converts a StudentProfileView database model to a StudentProfileDTO.
func (st *StudentProfileView) ToDTO() StudentProfileDTO {
	return StudentProfileDTO{
		IIRID:      st.IIRID,
		UserID:     st.UserID,
		FirstName:  st.FirstName,
		MiddleName: st.MiddleName,
		LastName:   st.LastName,
		SuffixName: st.SuffixName,
		Gender: Gender{
			ID:   st.GenderID,
			Name: st.GenderName,
		},
		Email:          st.Email,
		StudentNumber:  st.StudentNumber,
		ProfilePicture: st.ProfilePicture.String,
		Course: Course{
			ID:   st.CourseID,
			Code: st.CourseCode,
			Name: st.CourseName,
		},
		Section:   st.Section,
		YearLevel: st.YearLevel,
		Status: StudentStatus{
			ID:   st.StatusID,
			Name: st.StatusName,
		},
	}
}

// ToDTO converts a StudentPersonalInfoView to StudentPersonalInfoDTO.
func (view *StudentPersonalInfoView) ToDTO(
	emergencyAddr locations.AddressDTO,
) *StudentPersonalInfoDTO {
	statusDTO := StudentStatus{
		ID:   view.StatusID,
		Name: view.StatusName,
	}

	var gradYear *int
	if view.GraduationYear.Valid {
		gy := int(view.GraduationYear.Int64)
		gradYear = &gy
	}

	return &StudentPersonalInfoDTO{
		ID:            view.ID,
		IIRID:         view.IIRID,
		StudentNumber: view.StudentNumber,
		Gender: Gender{
			ID:   view.GenderID,
			Name: view.GenderName,
		},
		CivilStatus: CivilStatusType{
			ID:   view.CivilStatusID,
			Name: view.CivilStatusName,
		},
		Religion: Religion{
			ID:   view.ReligionID,
			Name: view.ReligionName,
		},
		HeightM:       view.HeightM,
		WeightKg:      view.WeightKg,
		Complexion:    view.Complexion,
		HighSchoolGWA: view.HighSchoolGWA,
		Course: Course{
			ID:   view.CourseID,
			Code: view.CourseCode,
			Name: view.CourseName,
		},
		YearLevel:             view.YearLevel,
		Section:               view.Section,
		PlaceOfBirth:          view.PlaceOfBirth,
		DateOfBirth:           view.DateOfBirth,
		TelephoneNumber:       view.TelephoneNumber,
		MobileNumber:          view.MobileNumber,
		IsEmployed:            view.IsEmployed,
		EmployerName:          view.EmployerName,
		EmployerAddress:       view.EmployerAddress,
		EmployerContactNumber: view.EmployerContactNumber,
		OtherReligionText:     view.OtherReligionText,
		TwoByTwoPhotoDataUrl:  view.TwoByTwoPhotoDataURL.String,
		Status:                statusDTO,
		GraduationYear:        gradYear,
		EmergencyContact: EmergencyContactDTO{
			ID:            view.EmergencyID,
			FirstName:     view.EmergencyFirstName,
			MiddleName:    view.EmergencyMiddleName,
			LastName:      view.EmergencyLastName,
			ContactNumber: view.EmergencyContactNumber,
			Relationship: StudentRelationshipType{
				ID:   view.EmergencyRelationshipID,
				Name: view.EmergencyRelationshipName,
			},
			Address: emergencyAddr,
		},
	}
}

// ToDTO converts a RelatedPersonView to RelatedPersonDTO.
func (view *RelatedPersonView) ToDTO() RelatedPersonDTO {
	return RelatedPersonDTO{
		ID: view.ID,
		EducationalAttainment: EducationalAttainment{
			ID:   view.EducationalAttainmentID,
			Name: view.EducationalAttainmentName,
		},
		DateOfBirth:     view.DateOfBirth,
		LastName:        view.LastName,
		FirstName:       view.FirstName,
		MiddleName:      view.MiddleName,
		SuffixName:      view.SuffixName,
		Occupation:      view.Occupation,
		EmployerName:    view.EmployerName,
		EmployerAddress: view.EmployerAddress,
		Relationship: StudentRelationshipType{
			ID:   view.RelationshipID,
			Name: view.RelationshipName,
		},
		IsParent:   view.IsParent,
		IsGuardian: view.IsGuardian,
		IsLiving:   view.IsLiving,
	}
}

// ToDTO converts a StudentFinanceView to StudentFinanceDTO.
func (view *StudentFinanceView) ToDTO(
	supports []StudentSupportType,
) StudentFinanceDTO {
	return StudentFinanceDTO{
		ID: view.ID,
		MonthlyFamilyIncomeRange: IncomeRange{
			ID:   view.IncomeRangeID,
			Text: view.IncomeRangeText,
		},
		OtherIncomeDetails:    view.OtherIncome,
		WeeklyAllowance:       view.WeeklyAllowance,
		FinancialSupportTypes: supports,
	}
}

// ToDTO converts a StudentHealthRecord to StudentHealthRecordDTO.
func (hr *StudentHealthRecord) ToDTO() StudentHealthRecordDTO {
	return StudentHealthRecordDTO{
		ID:                        hr.ID,
		VisionHasProblem:          hr.VisionHasProblem,
		VisionDetails:             hr.VisionDetails,
		HearingHasProblem:         hr.HearingHasProblem,
		HearingDetails:            hr.HearingDetails,
		SpeechHasProblem:          hr.SpeechHasProblem,
		SpeechDetails:             hr.SpeechDetails,
		GeneralHealthHasProblem:   hr.GeneralHealthHasProblem,
		GeneralHealthDetails:      hr.GeneralHealthDetails,
		MentalEmotionalHasProblem: hr.MentalEmotionalHasProblem,
		MentalEmotionalDetails:    hr.MentalEmotionalDetails,
	}
}

// ToDTO converts a StudentAddress model and address DTO to StudentAddressDTO.
func (addr *StudentAddress) ToDTO(
	addrDTO locations.AddressDTO,
) StudentAddressDTO {
	return StudentAddressDTO{
		ID:          addr.ID,
		Address:     addrDTO,
		AddressType: addr.AddressType,
		CreatedAt:   addr.CreatedAt,
		UpdatedAt:   addr.UpdatedAt,
	}
}

// ToDTO converts FamilyBackground to FamilyBackgroundDTO.
func (fb *FamilyBackground) ToDTO(
	parentalStatus ParentalStatusType,
	natureOfResidence NatureOfResidenceType,
	siblingSupports []SibilingSupportType,
) FamilyBackgroundDTO {
	return FamilyBackgroundDTO{
		ID:                    fb.ID,
		ParentalStatus:        parentalStatus,
		ParentalStatusDetails: fb.ParentalStatusDetails,
		Brothers:              &fb.Brothers,
		Sisters:               &fb.Sisters,
		EmployedSiblings:      &fb.EmployedSiblings,
		OrdinalPosition:       fb.OrdinalPosition,
		HaveQuietPlaceToStudy: fb.HaveQuietPlaceToStudy,
		IsSharingRoom:         fb.IsSharingRoom,
		SiblingSupportTypes:   siblingSupports,
		RoomSharingDetails:    fb.RoomSharingDetails,
		NatureOfResidence:     natureOfResidence,
	}
}

// ToDTO converts StudentConsultation to StudentConsultationDTO.
func (c *StudentConsultation) ToDTO() StudentConsultationDTO {
	return StudentConsultationDTO{
		ID:               c.ID,
		ProfessionalType: c.ProfessionalType,
		HasConsulted:     c.HasConsulted,
		WhenDate:         c.WhenDate,
		ForWhat:          c.ForWhat,
	}
}

// ToDTO converts StudentActivity to StudentActivityDTO.
func (a *StudentActivity) ToDTO(
	option ActivityOption,
) StudentActivityDTO {
	return StudentActivityDTO{
		ID:                 a.ID,
		ActivityOption:     option,
		OtherSpecification: a.OtherSpecification,
		Role:               a.Role,
		RoleSpecification:  a.RoleSpecification,
	}
}

// ToDTO converts StudentSubjectPreference to StudentSubjectPreferenceDTO.
func (p *StudentSubjectPreference) ToDTO() StudentSubjectPreferenceDTO {
	return StudentSubjectPreferenceDTO{
		ID:          p.ID,
		SubjectName: p.SubjectName,
		IsFavorite:  p.IsFavorite,
	}
}

// ToDTO converts StudentHobby to StudentHobbyDTO.
func (h *StudentHobby) ToDTO() StudentHobbyDTO {
	return StudentHobbyDTO{
		ID:           h.ID,
		HobbyName:    h.HobbyName,
		PriorityRank: h.PriorityRank,
	}
}

// ToDTO converts TestResult to TestResultDTO.
func (r *TestResult) ToDTO() TestResultDTO {
	return TestResultDTO{
		ID:          r.ID,
		TestDate:    r.TestDate,
		TestName:    r.TestName,
		RawScore:    r.RawScore,
		Percentile:  r.Percentile,
		Description: r.Description,
	}
}
