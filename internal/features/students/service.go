package students

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/pdf"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/files"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"golang.org/x/sync/errgroup"
)

//go:embed assets/iirf.html
var iirTemplate string

// Service provides student-related business logic and data access.
type Service struct {
	repo         *Repository
	locationsSvc *locations.Service
	userService  *users.Service
	filesSvc     *files.Service
	logService   audit.Logger
	notifService audit.Notifier
	cfg          *config.Config
	pdfService   *pdf.Service
}

// NewService creates a new student service instance.
func NewService(
	repo *Repository,
	locationsSvc *locations.Service,
	userService *users.Service,
	filesSvc *files.Service,
	logService audit.Logger,
	notifService audit.Notifier,
	cfg *config.Config,
	pdfService *pdf.Service,
) *Service {
	return &Service{
		repo:         repo,
		locationsSvc: locationsSvc,
		userService:  userService,
		filesSvc:     filesSvc,
		logService:   logService,
		notifService: notifService,
		cfg:          cfg,
		pdfService:   pdfService,
	}
}

func (s *Service) GetLatestCORsByUserIDs(
	ctx context.Context,
	userIDs []string,
) (map[string]string, error) {
	return s.repo.GetLatestCORsByUserIDs(ctx, userIDs)
}

// SubmitCOR uploads a Certificate of Registration and links it to the student.
func (s *Service) SubmitCOR(
	ctx context.Context,
	userID string,
	fileHeader *multipart.FileHeader,
) (string, error) {
	file, err := s.filesSvc.UploadFile(ctx, fileHeader, "cors")
	if err != nil {
		return "", fmt.Errorf("[StudentService] {SubmitCOR Upload}: %w", err)
	}

	cor := StudentCOR{
		FileID:    file.ID,
		StudentID: userID,
	}

	// Fetch OCR result to set validity
	ocrResult, err := s.filesSvc.GetOCRResult(ctx, file.ID)
	if err != nil {
		// Non-fatal, but we'll use fallback dates
		fmt.Printf("[StudentService] {SubmitCOR GetOCRResult}: %v\n", err)
	}

	if ocrResult != nil && ocrResult.StructuredData != "" {
		var corData struct {
			StartAcademicYear string `json:"start_academic_year"`
			EndAcademicYear   string `json:"end_academic_year"`
			Term              int    `json:"term"`
		}
		if err := json.Unmarshal([]byte(ocrResult.StructuredData), &corData); err == nil {
			// Roughly calculate validity based on term
			// Term 1: Aug to Dec
			// Term 2: Jan to May
			// Term 3/Summer: June to July
			startYear := time.Now().Year()
			if corData.StartAcademicYear != "" {
				fmt.Sscanf(corData.StartAcademicYear, "%d", &startYear)
			}

			if corData.Term == 1 {
				cor.ValidFrom = structs.TimeToNullableTime(time.Date(startYear, time.August, 1, 0, 0, 0, 0, time.Local))
				cor.ValidUntil = structs.TimeToNullableTime(time.Date(startYear, time.December, 31, 23, 59, 59, 0, time.Local))
			} else if corData.Term == 2 {
				cor.ValidFrom = structs.TimeToNullableTime(time.Date(startYear+1, time.January, 1, 0, 0, 0, 0, time.Local))
				cor.ValidUntil = structs.TimeToNullableTime(time.Date(startYear+1, time.May, 31, 23, 59, 59, 0, time.Local))
			} else {
				cor.ValidFrom = structs.TimeToNullableTime(time.Now())
				cor.ValidUntil = structs.TimeToNullableTime(time.Now().AddDate(0, 5, 0)) // Default 5 months
			}
		}
	}

	if !cor.ValidFrom.Valid {
		cor.ValidFrom = structs.TimeToNullableTime(time.Now())
		cor.ValidUntil = structs.TimeToNullableTime(time.Now().AddDate(0, 5, 0)) // Default 5 months
	}

	// Delete existing CORs to enforce 1-to-1 relationship and free up storage
	oldCORs, _ := s.repo.GetStudentCORsByUserID(ctx, userID)
	for _, oldCOR := range oldCORs {
		_ = s.filesSvc.DeleteFile(ctx, oldCOR.FileID)
	}

	err = s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		return s.repo.SaveStudentCOR(ctx, tx, cor)
	})
	if err != nil {
		return "", fmt.Errorf("[StudentService] {SubmitCOR Save}: %w", err)
	}

	return file.ID, nil
}

func (s *Service) GetStudentCOR(
	ctx context.Context,
	userID string,
) (StudentCOR, error) {
	return s.repo.GetStudentCORByUserID(ctx, userID)
}

func (s *Service) GetStudentCORs(
	ctx context.Context,
	userID string,
) ([]StudentCOR, error) {
	return s.repo.GetStudentCORsByUserID(ctx, userID)
}

func (s *Service) GetGenders(ctx context.Context) ([]Gender, error) {
	return s.repo.GetGenders(ctx)
}

func (s *Service) GetParentalStatusTypes(ctx context.Context) ([]ParentalStatusType, error) {
	return s.repo.GetParentalStatusTypes(ctx)
}

func (s *Service) GetIncomeRanges(ctx context.Context) ([]IncomeRange, error) {
	return s.repo.GetIncomeRanges(ctx)
}

func (s *Service) GetStudentSupportTypes(ctx context.Context) ([]StudentSupportType, error) {
	return s.repo.GetStudentSupportTypes(ctx)
}

func (s *Service) GetSiblingSupportTypes(ctx context.Context) ([]SibilingSupportType, error) {
	return s.repo.GetSiblingSupportTypes(ctx)
}

func (s *Service) GetEducationalLevels(ctx context.Context) ([]EducationalLevel, error) {
	return s.repo.GetEducationalLevels(ctx)
}

func (s *Service) GetCourses(ctx context.Context) ([]Course, error) {
	return s.repo.GetCourses(ctx)
}

func (s *Service) GetCivilStatusTypes(ctx context.Context) ([]CivilStatusType, error) {
	return s.repo.GetCivilStatusTypes(ctx)
}

func (s *Service) GetReligions(ctx context.Context) ([]Religion, error) {
	return s.repo.GetReligions(ctx)
}

func (s *Service) GetNatureOfResidenceTypes(ctx context.Context) ([]NatureOfResidenceType, error) {
	return s.repo.GetNatureOfResidenceTypes(ctx)
}

func (s *Service) GetActivityOptions(ctx context.Context) ([]ActivityOption, error) {
	return s.repo.GetActivityOptions(ctx)
}

func (s *Service) GetStudentRelationshipTypes(ctx context.Context) ([]StudentRelationshipType, error) {
	return s.repo.GetStudentRelationshipTypes(ctx)
}

func (s *Service) GetStudentStatuses(ctx context.Context) ([]StudentStatus, error) {
	return s.repo.GetStudentStatuses(ctx)
}

func (s *Service) ListStudents(
	ctx context.Context, req ListStudentsRequest,
) (*ListStudentsResponse, error) {
	req.SetDefaults("last_name")

	students, err := s.repo.ListStudents(
		ctx,
		req.Search,
		req.GetOffset(),
		req.PageSize,
		req.OrderBy,
		req.CourseID,
		req.GenderID,
		req.YearLevel,
		req.StatusID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	studentDTOs := make([]StudentProfileDTO, len(students))
	for i, st := range students {
		course, _ := s.repo.GetCourseByID(ctx, st.CourseID)
		gender, _ := s.repo.GetGenderByID(ctx, st.GenderID)

		studentDTOs[i] = StudentProfileDTO{
			IIRID:         st.IIRID,
			UserID:        st.UserID,
			FirstName:     st.FirstName,
			MiddleName:    st.MiddleName,
			LastName:      st.LastName,
			SuffixName:    st.SuffixName,
			Gender:        *gender,
			Email:         st.Email,
			StudentNumber: st.StudentNumber,
			Course:        *course,
			Section:       st.Section,
			YearLevel:     st.YearLevel,
			Status: StudentStatus{
				ID:   st.StatusID,
				Name: st.StatusName,
			},
		}
	}

	total, err := s.repo.GetTotalStudentsCount(
		ctx,
		req.Search,
		req.CourseID,
		req.GenderID,
		req.YearLevel,
		req.StatusID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get total students count: %w", err)
	}

	return &ListStudentsResponse{
		Students: studentDTOs,
		Meta:     structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetStudentProfile(
	ctx context.Context,
	iirID string,
) (*ComprehensiveProfileDTO, error) {
	// Fetch user ID for this IIR
	iir, err := s.repo.GetStudentIIR(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch IIR record: %w", err)
	}

	profile := &ComprehensiveProfileDTO{IIRID: iirID}

	// Fetch COR URL
	corMap, _ := s.repo.GetLatestCORsByUserIDs(ctx, []string{iir.UserID})
	profile.StudentCORURL = corMap[iir.UserID]

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		basicInfo, err := s.GetStudentBasicInfo(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Student.BasicInfo = *basicInfo
		return nil
	})

	g.Go(func() error {
		personalInfo, err := s.GetStudentPersonalInfo(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Student.StudentPersonalInfoDTO = *personalInfo
		return nil
	})

	g.Go(func() error {
		addresses, err := s.GetStudentAddresses(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Student.Addresses = addresses
		return nil
	})

	g.Go(func() error {
		education, err := s.GetEducationalBackground(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Education = *education
		return nil
	})

	g.Go(func() error {
		familyBackground, err := s.GetStudentFamilyBackground(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Family.FamilyBackgroundDTO = *familyBackground
		return nil
	})

	g.Go(func() error {
		relatedPersons, err := s.GetStudentRelatedPersons(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Family.RelatedPersons = relatedPersons
		return nil
	})

	g.Go(func() error {
		finance, err := s.GetStudentFinancialInfo(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Family.Finance = *finance
		return nil
	})

	g.Go(func() error {
		healthRecord, err := s.GetStudentHealthRecord(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Health.StudentHealthRecordDTO = *healthRecord
		return nil
	})

	g.Go(func() error {
		consultations, err := s.GetStudentConsultations(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Health.Consultations = consultations
		return nil
	})

	g.Go(func() error {
		activities, err := s.GetStudentActivities(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Interests.Activities = activities
		return nil
	})

	g.Go(func() error {
		subjectPreferences, err := s.GetStudentSubjectPreferences(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Interests.SubjectPreferences = subjectPreferences
		return nil
	})

	g.Go(func() error {
		hobbies, err := s.GetStudentHobbies(ctx, iirID)
		if err != nil {
			return err
		}
		profile.Interests.Hobbies = hobbies
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *Service) GetStudentBasicInfo(
	ctx context.Context,
	iirID string,
) (*StudentBasicInfoViewDTO, error) {
	info, err := s.repo.GetStudentBasicInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student basic info: %w", err)
	}

	return &StudentBasicInfoViewDTO{
		Email:      info.Email,
		FirstName:  info.FirstName,
		MiddleName: info.MiddleName,
		LastName:   info.LastName,
	}, nil
}

func (s *Service) GetIIRDraft(
	ctx context.Context,
	userID string,
) (*ComprehensiveProfileDTO, error) {
	draft, err := s.repo.GetIIRDraftByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get IIR draft: %w", err)
	}

	if draft == nil || len(draft.Data) == 0 {
		return nil, nil
	}

	var draftData ComprehensiveProfileDTO
	if err := json.Unmarshal([]byte(draft.Data), &draftData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IIR draft data: %w", err)
	}

	return &draftData, nil
}

func (s *Service) GetStudentIIRByUserID(
	ctx context.Context,
	userID string,
) (*IIRRecord, error) {
	return s.repo.GetStudentIIRByUserID(ctx, userID)
}

func (s *Service) GetStudentIIR(
	ctx context.Context,
	iirID string,
) (*IIRRecord, error) {
	return s.repo.GetStudentIIR(ctx, iirID)
}

func (s *Service) GetStudentPersonalInfo(
	ctx context.Context,
	iirID string,
) (*StudentPersonalInfoDTO, error) {
	personalInfo, err := s.repo.GetStudentPersonalInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student personal info: %w", err)
	}

	gender, _ := s.repo.GetGenderByID(ctx, personalInfo.GenderID)
	civilStatus, _ := s.repo.GetCivilStatusByID(ctx, personalInfo.CivilStatusID)
	religion, _ := s.repo.GetReligionByID(ctx, personalInfo.ReligionID)
	course, _ := s.repo.GetCourseByID(ctx, personalInfo.CourseID)
	emergencyContact, _ := s.repo.GetEmergencyContactByIIRID(ctx, personalInfo.IIRID)
	emergencyContactRelationship, _ := s.repo.GetStudentRelationshipByID(ctx, emergencyContact.RelationshipID)
	emergencyAddressDTO, _ := s.locationsSvc.GetAddressByID(ctx, emergencyContact.AddressID)

	emergencyContactDTO := EmergencyContactDTO{
		ID:            emergencyContact.ID,
		FirstName:     emergencyContact.FirstName,
		MiddleName:    emergencyContact.MiddleName,
		LastName:      emergencyContact.LastName,
		ContactNumber: emergencyContact.ContactNumber,
		Relationship:  *emergencyContactRelationship,
		Address:       emergencyAddressDTO,
	}

	return &StudentPersonalInfoDTO{
		ID:               personalInfo.ID,
		StudentNumber:    personalInfo.StudentNumber,
		Gender:           *gender,
		CivilStatus:      *civilStatus,
		Religion:         *religion,
		HeightFt:         personalInfo.HeightFt,
		WeightKg:         personalInfo.WeightKg,
		Complexion:       personalInfo.Complexion,
		HighSchoolGWA:    personalInfo.HighSchoolGWA,
		Course:           *course,
		YearLevel:        personalInfo.YearLevel,
		Section:          personalInfo.Section,
		PlaceOfBirth:     personalInfo.PlaceOfBirth,
		DateOfBirth:      personalInfo.DateOfBirth,
		TelephoneNumber:  personalInfo.TelephoneNumber,
		MobileNumber:     personalInfo.MobileNumber,
		IsEmployed:       personalInfo.IsEmployed,
		EmployerName:     personalInfo.EmployerName,
		EmployerAddress:  personalInfo.EmployerAddress,
		EmergencyContact: emergencyContactDTO,
	}, nil
}

func (s *Service) GetStudentAddresses(
	ctx context.Context,
	iirID string,
) ([]StudentAddressDTO, error) {
	studentAddresses, err := s.repo.GetStudentAddresses(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student addresses: %w", err)
	}

	var addresses []StudentAddressDTO
	for _, addr := range studentAddresses {
		addrDTO, _ := s.locationsSvc.GetAddressByID(ctx, addr.AddressID)
		addresses = append(addresses, StudentAddressDTO{
			ID:          addr.ID,
			Address:     addrDTO,
			AddressType: addr.AddressType,
			CreatedAt:   addr.CreatedAt,
			UpdatedAt:   addr.UpdatedAt,
		})
	}
	return addresses, nil
}

func (s *Service) GetStudentFamilyBackground(
	ctx context.Context,
	iirID string,
) (*FamilyBackgroundDTO, error) {
	studentFamily, err := s.repo.GetStudentFamilyBackground(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get family background: %w", err)
	}
	if studentFamily == nil {
		return nil, nil
	}

	parentalStatus, _ := s.repo.GetParentalStatusByID(ctx, studentFamily.ParentalStatusID)
	natureOfResidence, _ := s.repo.GetNatureOfResidenceByID(ctx, studentFamily.NatureOfResidenceId)
	siblingSupportTypes, _ := s.repo.GetStudentSiblingSupport(ctx, studentFamily.ID)

	var supportTypes []SibilingSupportType
	for _, sst := range siblingSupportTypes {
		st, _ := s.repo.GetSiblingSupportTypeByID(ctx, sst.SupportTypeID)
		supportTypes = append(supportTypes, *st)
	}

	return &FamilyBackgroundDTO{
		ID:                    studentFamily.ID,
		ParentalStatus:        *parentalStatus,
		ParentalStatusDetails: studentFamily.ParentalStatusDetails,
		Brothers:              &studentFamily.Brothers,
		Sisters:               &studentFamily.Sisters,
		EmployedSiblings:      &studentFamily.EmployedSiblings,
		OrdinalPosition:       studentFamily.OrdinalPosition,
		HaveQuietPlaceToStudy: studentFamily.HaveQuietPlaceToStudy,
		IsSharingRoom:         studentFamily.IsSharingRoom,
		SiblingSupportTypes:   supportTypes,
		RoomSharingDetails:    studentFamily.RoomSharingDetails,
		NatureOfResidence:     *natureOfResidence,
	}, nil
}

func (s *Service) GetStudentRelatedPersons(
	ctx context.Context,
	iirID string,
) ([]RelatedPersonDTO, error) {
	studentRelatedPersons, err := s.repo.GetStudentRelatedPersons(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get student related persons: %w", err)
	}

	var related []RelatedPersonDTO
	for _, srp := range studentRelatedPersons {
		rp, _ := s.repo.GetRelatedPersonByID(ctx, srp.RelatedPersonID)
		rel, _ := s.repo.GetStudentRelationshipByID(ctx, srp.RelationshipID)

		related = append(related, RelatedPersonDTO{
			ID:               rp.ID,
			EducationalLevel: rp.EducationalLevel,
			DateOfBirth:      rp.DateOfBirth,
			LastName:         rp.LastName,
			FirstName:        rp.FirstName,
			MiddleName:       rp.MiddleName,
			SuffixName:       rp.SuffixName,
			Occupation:       rp.Occupation,
			EmployerName:     rp.EmployerName,
			EmployerAddress:  rp.EmployerAddress,
			Relationship:     *rel,
			IsParent:         srp.IsParent,
			IsGuardian:       srp.IsGuardian,
			IsLiving:         srp.IsLiving,
		})
	}
	return related, nil
}

func (s *Service) GetEducationalBackground(
	ctx context.Context,
	iirID string,
) (*EducationalBackgroundDTO, error) {
	eb, err := s.repo.GetStudentEducationalBackground(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get educational background: %w", err)
	}
	if eb == nil {
		return nil, nil
	}

	schoolDetails, _ := s.repo.GetSchoolDetailsByEBID(ctx, eb.ID)
	var details []SchoolDetailsDTO
	for _, sd := range schoolDetails {
		level, _ := s.repo.GetEducationalLevelByID(ctx, sd.EducationalLevelID)
		details = append(details, SchoolDetailsDTO{
			ID:               sd.ID,
			EducationalLevel: *level,
			SchoolName:       sd.SchoolName,
			SchoolAddress:    sd.SchoolAddress,
			SchoolType:       sd.SchoolType,
			YearStarted:      sd.YearStarted,
			YearCompleted:    sd.YearCompleted,
			Awards:           sd.Awards,
		})
	}

	return &EducationalBackgroundDTO{
		ID:                 eb.ID,
		NatureOfSchooling:  eb.NatureOfSchooling,
		InterruptedDetails: eb.InterruptedDetails,
		School:             details,
	}, nil
}

func (s *Service) GetStudentFinancialInfo(
	ctx context.Context,
	iirID string,
) (*StudentFinanceDTO, error) {
	finance, err := s.repo.GetStudentFinancialInfo(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get financial info: %w", err)
	}
	if finance == nil {
		return nil, nil
	}

	supportTypes, _ := s.repo.GetFinancialSupportTypeByFinanceID(ctx, finance.ID)
	var supportDTOs []StudentSupportType
	for _, st := range supportTypes {
		support, _ := s.repo.GetStudentSupportByID(ctx, st.SupportTypeID)
		supportDTOs = append(supportDTOs, *support)
	}

	incomeRange, _ := s.repo.GetIncomeRangeByID(ctx, finance.IncomeRangeID)

	return &StudentFinanceDTO{
		ID:                       finance.ID,
		MonthlyFamilyIncomeRange: *incomeRange,
		OtherIncomeDetails:       finance.OtherIncome,
		WeeklyAllowance:          finance.WeeklyAllowance,
		FinancialSupportTypes:    supportDTOs,
	}, nil
}

func (s *Service) GetStudentHealthRecord(
	ctx context.Context,
	iirID string,
) (*StudentHealthRecordDTO, error) {
	hr, err := s.repo.GetStudentHealthRecord(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get health record: %w", err)
	}
	if hr == nil {
		return nil, nil
	}

	return &StudentHealthRecordDTO{
		ID:                      hr.ID,
		VisionHasProblem:        hr.VisionHasProblem,
		VisionDetails:           hr.VisionDetails,
		HearingHasProblem:       hr.HearingHasProblem,
		HearingDetails:          hr.HearingDetails,
		SpeechHasProblem:        hr.SpeechHasProblem,
		SpeechDetails:           hr.SpeechDetails,
		GeneralHealthHasProblem: hr.GeneralHealthHasProblem,
		GeneralHealthDetails:    hr.GeneralHealthDetails,
	}, nil
}

func (s *Service) GetStudentConsultations(
	ctx context.Context,
	iirID string,
) ([]StudentConsultationDTO, error) {
	consultations, err := s.repo.GetStudentConsultations(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consultations: %w", err)
	}

	var dtos []StudentConsultationDTO
	for _, c := range consultations {
		dtos = append(dtos, StudentConsultationDTO{
			ID:               c.ID,
			ProfessionalType: c.ProfessionalType,
			HasConsulted:     c.HasConsulted,
			WhenDate:         c.WhenDate,
			ForWhat:          c.ForWhat,
		})
	}
	return dtos, nil
}

func (s *Service) GetStudentActivities(
	ctx context.Context,
	iirID string,
) ([]StudentActivityDTO, error) {
	activities, err := s.repo.GetStudentActivities(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}

	var dtos []StudentActivityDTO
	for _, a := range activities {
		option, _ := s.repo.GetActivityOptionByID(ctx, a.OptionID)
		dtos = append(dtos, StudentActivityDTO{
			ID:                 a.ID,
			ActivityOption:     *option,
			OtherSpecification: a.OtherSpecification,
			Role:               a.Role,
			RoleSpecification:  a.RoleSpecification,
		})
	}
	return dtos, nil
}

func (s *Service) GetStudentSubjectPreferences(
	ctx context.Context,
	iirID string,
) ([]StudentSubjectPreferenceDTO, error) {
	prefs, err := s.repo.GetStudentSubjectPreferences(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subject preferences: %w", err)
	}

	var dtos []StudentSubjectPreferenceDTO
	for _, p := range prefs {
		dtos = append(dtos, StudentSubjectPreferenceDTO{
			ID:          p.ID,
			SubjectName: p.SubjectName,
			IsFavorite:  p.IsFavorite,
		})
	}
	return dtos, nil
}

func (s *Service) GetStudentHobbies(ctx context.Context, iirID string) ([]StudentHobbyDTO, error) {
	hobbies, err := s.repo.GetStudentHobbies(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hobbies: %w", err)
	}

	var dtos []StudentHobbyDTO
	for _, h := range hobbies {
		dtos = append(dtos, StudentHobbyDTO{
			ID:           h.ID,
			HobbyName:    h.HobbyName,
			PriorityRank: h.PriorityRank,
		})
	}
	return dtos, nil
}

func (s *Service) GetStudentTestResults(
	ctx context.Context,
	iirID string,
) ([]TestResultDTO, error) {
	results, err := s.repo.GetStudentTestResults(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test results: %w", err)
	}

	var dtos []TestResultDTO
	for _, r := range results {
		dtos = append(dtos, TestResultDTO{
			ID:          r.ID,
			TestDate:    r.TestDate,
			TestName:    r.TestName,
			RawScore:    r.RawScore,
			Percentile:  r.Percentile,
			Description: r.Description,
		})
	}
	return dtos, nil
}

func (s *Service) SaveIIRDraft(
	ctx context.Context,
	userID string,
	req ComprehensiveProfileDTO,
) (int, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal IIR draft: %w", err)
	}

	draft := IIRDraft{
		UserID: userID,
		Data:   string(data),
	}

	id, err := s.repo.UpsertIIRDraft(ctx, draft)
	if err != nil {
		return 0, fmt.Errorf("failed to save IIR draft: %w", err)
	}

	return id, nil
}

func (s *Service) SubmitStudentIIR(
	ctx context.Context,
	userID string,
	req ComprehensiveProfileDTO,
) (string, error) {
	var iirID string
	err := s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		id, err := s.saveComprehensiveProfile(ctx, tx, userID, req)
		if err != nil {
			return err
		}
		iirID = id
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("[StudentService] {SubmitStudentIIR}: %w", err)
	}

	return iirID, nil
}

func (s *Service) saveComprehensiveProfile(
	ctx context.Context,
	tx datastore.DB,
	userID string,
	req ComprehensiveProfileDTO,
) (string, error) {
	iirID := req.IIRID
	if iirID == "" {
		iirID = uuid.New().String()
	}

	// 1. IIR Record Header
	_, err := s.repo.UpsertIIRRecord(ctx, tx, &IIRRecord{
		ID:          iirID,
		UserID:      userID,
		IsSubmitted: true,
	})
	if err != nil {
		return "", err
	}

	// 2. Personal Info
	err = s.repo.UpsertStudentPersonalInfo(ctx, tx, &StudentPersonalInfo{
		IIRID:           iirID,
		StudentNumber:   req.Student.StudentNumber,
		GenderID:        req.Student.Gender.ID,
		CivilStatusID:   req.Student.CivilStatus.ID,
		ReligionID:      req.Student.Religion.ID,
		HeightFt:        req.Student.HeightFt,
		WeightKg:        req.Student.WeightKg,
		Complexion:      req.Student.Complexion,
		HighSchoolGWA:   req.Student.HighSchoolGWA,
		CourseID:        req.Student.Course.ID,
		YearLevel:       req.Student.YearLevel,
		Section:         req.Student.Section,
		PlaceOfBirth:    req.Student.PlaceOfBirth,
		DateOfBirth:     req.Student.DateOfBirth,
		MobileNumber:    req.Student.MobileNumber,
		TelephoneNumber: req.Student.TelephoneNumber,
		IsEmployed:      req.Student.IsEmployed,
		EmployerName:    req.Student.EmployerName,
		EmployerAddress: req.Student.EmployerAddress,
		StatusID:        1, // Default to Enrolled/Active?
	})
	if err != nil {
		return "", err
	}

	// 3. Emergency Contact
	ecAddrID, err := s.locationsSvc.UpsertAddress(ctx, tx, req.Student.EmergencyContact.Address)
	if err != nil {
		return "", err
	}

	_, err = s.repo.UpsertEmergencyContact(ctx, tx, &EmergencyContact{
		IIRID:          iirID,
		FirstName:      req.Student.EmergencyContact.FirstName,
		MiddleName:     req.Student.EmergencyContact.MiddleName,
		LastName:       req.Student.EmergencyContact.LastName,
		ContactNumber:  req.Student.EmergencyContact.ContactNumber,
		RelationshipID: req.Student.EmergencyContact.Relationship.ID,
		AddressID:      ecAddrID,
	})
	if err != nil {
		return "", err
	}

	// 4. Addresses
	for _, addrDTO := range req.Student.Addresses {
		addrID, err := s.locationsSvc.UpsertAddress(ctx, tx, addrDTO.Address)
		if err != nil {
			return "", err
		}
		_, err = s.repo.UpsertStudentAddress(ctx, tx, &StudentAddress{
			IIRID:       iirID,
			AddressID:   addrID,
			AddressType: addrDTO.AddressType,
		})
		if err != nil {
			return "", err
		}
	}

	// 5. Family Background
	fbID, err := s.repo.UpsertFamilyBackground(ctx, tx, &FamilyBackground{
		IIRID:                 iirID,
		ParentalStatusID:      req.Family.ParentalStatus.ID,
		ParentalStatusDetails: req.Family.ParentalStatusDetails,
		Brothers:              *req.Family.Brothers,
		Sisters:               *req.Family.Sisters,
		EmployedSiblings:      *req.Family.EmployedSiblings,
		OrdinalPosition:       req.Family.OrdinalPosition,
		HaveQuietPlaceToStudy: req.Family.HaveQuietPlaceToStudy,
		IsSharingRoom:         req.Family.IsSharingRoom,
		RoomSharingDetails:    req.Family.RoomSharingDetails,
		NatureOfResidenceId:   req.Family.NatureOfResidence.ID,
	})
	if err != nil {
		return "", err
	}

	_ = s.repo.DeleteStudentSiblingSupportsByFamilyID(ctx, tx, fbID)
	for _, sst := range req.Family.SiblingSupportTypes {
		_ = s.repo.CreateStudentSiblingSupport(ctx, tx, &StudentSiblingSupport{
			FamilyBackgroundID: fbID,
			SupportTypeID:      sst.ID,
		})
	}

	// 6. Educational Background
	ebID, err := s.repo.UpsertEducationalBackground(ctx, tx, &EducationalBackground{
		IIRID:              iirID,
		NatureOfSchooling:  req.Education.NatureOfSchooling,
		InterruptedDetails: req.Education.InterruptedDetails,
	})
	if err != nil {
		return "", err
	}

	_ = s.repo.DeleteSchoolDetailsByEBID(ctx, tx, ebID)
	for _, sd := range req.Education.School {
		_, _ = s.repo.UpsertSchoolDetails(ctx, tx, &SchoolDetails{
			EducationalLevelID: sd.EducationalLevel.ID,
			SchoolName:         sd.SchoolName,
			SchoolAddress:      sd.SchoolAddress,
			SchoolType:         sd.SchoolType,
			YearStarted:        sd.YearStarted,
			YearCompleted:      sd.YearCompleted,
			Awards:             sd.Awards,
			EBID:               ebID,
		})
	}

	// 7. Health Record
	_, err = s.repo.UpsertStudentHealthRecord(ctx, tx, &StudentHealthRecord{
		IIRID:                   iirID,
		VisionHasProblem:        req.Health.VisionHasProblem,
		VisionDetails:           req.Health.VisionDetails,
		HearingHasProblem:       req.Health.HearingHasProblem,
		HearingDetails:          req.Health.HearingDetails,
		SpeechHasProblem:        req.Health.SpeechHasProblem,
		SpeechDetails:           req.Health.SpeechDetails,
		GeneralHealthHasProblem: req.Health.GeneralHealthHasProblem,
		GeneralHealthDetails:    req.Health.GeneralHealthDetails,
	})
	if err != nil {
		return "", err
	}

	// 8. Finance
	sfID, err := s.repo.UpsertStudentFinance(ctx, tx, &StudentFinance{
		IIRID:           iirID,
		IncomeRangeID:   req.Family.Finance.MonthlyFamilyIncomeRange.ID,
		OtherIncome:     req.Family.Finance.OtherIncomeDetails,
		WeeklyAllowance: req.Family.Finance.WeeklyAllowance,
	})
	if err != nil {
		return "", err
	}

	_ = s.repo.DeleteStudentFinancialSupportsByFinanceID(ctx, tx, sfID)
	for _, st := range req.Family.Finance.FinancialSupportTypes {
		_ = s.repo.CreateStudentFinancialSupport(ctx, tx, &StudentFinancialSupport{
			StudentFinanceID: sfID,
			SupportTypeID:    st.ID,
		})
	}

	return iirID, nil
}

func (s *Service) GenerateIIR(
	ctx context.Context,
	iirID string,
) ([]byte, string, error) {
	profile, err := s.GetStudentProfile(ctx, iirID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get student profile: %w", err)
	}

	data := struct {
		Profile   *ComprehensiveProfileDTO
		DateToday string
	}{
		Profile:   profile,
		DateToday: s.GetFormattedDate(time.Now().Format("2006-01-02")),
	}

	pdfBytes, err := s.pdfService.GenerateFromContent(
		ctx,
		"iirf.html",
		iirTemplate,
		data,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate pdf from template: %w", err)
	}

	fileName := fmt.Sprintf(
		"%s_IIR_%s_%s.pdf",
		profile.Student.BasicInfo.LastName,
		profile.Student.StudentNumber,
		time.Now().Format("20060102"),
	)

	return pdfBytes, fileName, nil
}

func (s *Service) GetFormattedDate(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return ""
	}
	return t.Format("January 02, 2006")
}

func (s *Service) IsStudentLocked(
	ctx context.Context,
	iirID string,
) (bool, error) {
	return s.repo.IsStudentLocked(ctx, iirID)
}

const (
	diplomaFinalYear  = 3
	bachelorFinalYear = 4
	statusIDGraduated = 2
)

func (s *Service) BulkUpdateStudentStatus(
	ctx context.Context,
	req BulkUpdateStatusRequest,
) error {
	if req.StatusID != statusIDGraduated {
		return s.repo.BulkUpdateStudentStatus(ctx, req)
	}

	// Eligibility logic for graduation
	rows, err := s.repo.ListStudents(
		ctx,
		req.Filters.Search,
		0, 100000,
		"created_at DESC",
		req.Filters.CourseID,
		0,
		req.Filters.YearLevel,
		0,
	)
	if err != nil {
		return err
	}

	var eligibleIDs []string
	for _, r := range rows {
		course, _ := s.repo.GetCourseByID(ctx, r.CourseID)
		code := strings.ToUpper(course.Code)
		isDiploma := strings.HasPrefix(code, "D") && !strings.HasPrefix(code, "DS")
		isBachelor := strings.HasPrefix(code, "BS")

		if (isDiploma && r.YearLevel == diplomaFinalYear) ||
			(isBachelor && r.YearLevel == bachelorFinalYear) {
			eligibleIDs = append(eligibleIDs, r.IIRID)
		}
	}

	if len(eligibleIDs) == 0 {
		return nil
	}

	req.SelectAllMatching = false
	req.ExcludedIIRIDs = nil
	req.IIRIDs = eligibleIDs

	return s.repo.BulkUpdateStudentStatus(ctx, req)
}
