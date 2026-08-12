package students

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
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

// CheckStudentNumberExists checks if a student number is already taken.
func (s *Service) CheckStudentNumberExists(
	ctx context.Context,
	studentNumber string,
	currentUserID string,
) (bool, error) {
	return s.repo.CheckStudentNumberExists(
		ctx,
		studentNumber,
		currentUserID,
	)
}

func (s *Service) GetLatestCORsByUserIDs(
	ctx context.Context,
	userIDs []string,
) (map[string]string, error) {
	return s.repo.GetLatestCORsByUserIDs(ctx, userIDs)
}

var (
	ErrOutdatedCOR = errors.New(
		"uploaded COR is for an outdated academic year or term",
	)
	ErrCOROwnerMismatch = errors.New("uploaded COR does not belong to the user")
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// SubmitCOR uploads a Certificate of Registration and links it to the student.
func (s *Service) SubmitCOR(
	ctx context.Context,
	userID string,
	fileHeader *multipart.FileHeader,
) (string, error) {
	file, err := s.filesSvc.UploadFile(ctx, fileHeader, "cors")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	cor := StudentCOR{
		FileID:    file.ID,
		StudentID: userID,
	}

	// Fetch OCR result to set validity
	ocrResult, err := s.filesSvc.GetOCRResult(ctx, file.ID)
	if err != nil {
		// Non-fatal, but we'll use fallback dates
		fmt.Printf("%v\n", err)
	}

	if ocrResult != nil && ocrResult.StructuredData != "" {
		var corData struct {
			StudentNumber     string `json:"student_number"`
			ProgramCode       string `json:"program_code"`
			ProgramDesc       string `json:"program_desc"`
			YearLevel         int    `json:"year_level"`
			Section           int    `json:"section"`
			Campus            string `json:"campus"`
			StartAcademicYear string `json:"start_academic_year"`
			EndAcademicYear   string `json:"end_academic_year"`
			Term              int    `json:"term"`
		}
		if err := json.Unmarshal(
			[]byte(ocrResult.StructuredData), &corData,
		); err == nil {
			startYear := time.Now().Year()
			if corData.StartAcademicYear != "" {
				fmt.Sscanf(
					corData.StartAcademicYear,
					"%d",
					&startYear,
				)
			}
			endYear := startYear + 1
			if corData.EndAcademicYear != "" {
				fmt.Sscanf(
					corData.EndAcademicYear,
					"%d",
					&endYear,
				)
			}
			cor.StudentNumber = corData.StudentNumber
			cor.ProgramCode = corData.ProgramCode
			cor.ProgramDesc = corData.ProgramDesc
			cor.YearLevel = corData.YearLevel
			cor.Section = corData.Section
			cor.Campus = corData.Campus
			cor.Term = corData.Term
			cor.YearStart = startYear
			cor.YearEnd = endYear

			// Validate against the current global AcademicSetting.
			// If OCR year + term do not match, automatically reject the COR.
			setting, sErr := s.repo.GetAcademicSetting(ctx)
			if sErr != nil {
				_ = s.filesSvc.DeleteFile(ctx, file.ID)
				return "", fmt.Errorf(
					"%w",
					sErr,
				)
			}

			if startYear != setting.CurrentYearStart ||
				corData.Term != setting.CurrentTerm {
				_ = s.filesSvc.DeleteFile(ctx, file.ID)
				return "", ErrOutdatedCOR
			}

			// OPTIONAL / FUTURE IMPLEMENTATION:
			// Verify that the COR really belongs to the student.
			// Check user's first name, last name, and student number.
			/*
				userRecord, uErr := s.repo.GetUserByID(ctx, userID)
				if uErr == nil {
					ocrStudNum := strings.TrimSpace(corData.StudentNumber)
					dbStudNum := strings.TrimSpace(userRecord.StudentNumber)

					// Perform matches on student number, first name, last name
					if ocrStudNum != "" && dbStudNum != "" && ocrStudNum != dbStudNum {
						_ = s.filesSvc.DeleteFile(ctx, file.ID)
						return "", ErrCOROwnerMismatch
					}
				}
			*/

			// If valid, set ValidFrom/ValidUntil
			cor.ValidFrom = structs.TimeToNullableTime(time.Now())
			cor.ValidUntil = structs.TimeToNullableTime(
				time.Now().AddDate(0, 5, 0),
			)
		}
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

func (s *Service) GetAcademicSetting(
	ctx context.Context,
) (*AcademicSetting, error) {
	return s.repo.GetAcademicSetting(ctx)
}

// UpdateAcademicSetting updates the global academic year + term setting.
// It fetches the old values first for audit log comparison.
func (s *Service) UpdateAcademicSetting(
	ctx context.Context,
	req UpdateAcademicSettingDTO,
	updaterID string,
	updaterEmail string,
) error {
	old, err := s.repo.GetAcademicSetting(ctx)
	if err != nil {
		return fmt.Errorf(
			"[StudentService] {UpdateAcademicSetting Fetch Old}: %w",
			err,
		)
	}

	err = s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		return s.repo.UpdateAcademicSetting(
			ctx,
			tx,
			req.CurrentYearStart,
			req.CurrentYearEnd,
			req.CurrentTerm,
			*req.AllowExpeditedIIR,
		)
	})
	if err != nil {
		return fmt.Errorf(
			"[StudentService] {UpdateAcademicSetting Save}: %w",
			err,
		)
	}

	s.logService.Record(ctx, s.repo.GetDB(), audit.LogEntry{
		Level:    audit.LevelInfo,
		Category: audit.CategoryAudit,
		Action:   audit.ActionSettingChanged,
		Message:  "Global academic year/term setting updated",
		UserID:   structs.StringToNullableString(updaterID),
		UserEmail: structs.StringToNullableString(
			updaterEmail,
		),
		Metadata: &audit.LogMetadata{
			EntityType: "AcademicSetting",
			OldValues:  old,
			NewValues:  req,
		},
	})

	return nil
}

func (s *Service) GetGenders(ctx context.Context) ([]Gender, error) {
	return s.repo.GetGenders(ctx)
}

func (s *Service) GetEnrollmentYears(ctx context.Context) ([]int, error) {
	return s.repo.GetEnrollmentYears(ctx)
}

func (s *Service) GetParentalStatusTypes(
	ctx context.Context,
) ([]ParentalStatusType, error) {
	return s.repo.GetParentalStatusTypes(ctx)
}

func (s *Service) GetIncomeRanges(ctx context.Context) ([]IncomeRange, error) {
	return s.repo.GetIncomeRanges(ctx)
}

func (s *Service) GetStudentSupportTypes(
	ctx context.Context,
) ([]StudentSupportType, error) {
	return s.repo.GetStudentSupportTypes(ctx)
}

func (s *Service) GetSiblingSupportTypes(
	ctx context.Context,
) ([]SibilingSupportType, error) {
	return s.repo.GetSiblingSupportTypes(ctx)
}

func (s *Service) GetEducationalLevels(
	ctx context.Context,
) ([]EducationalLevel, error) {
	return s.repo.GetEducationalLevels(ctx)
}

func (s *Service) GetEducationalAttainments(
	ctx context.Context,
) ([]EducationalAttainment, error) {
	return s.repo.GetEducationalAttainments(ctx)
}

func (s *Service) GetPrograms(
	ctx context.Context,
) ([]Program, error) {
	return s.repo.GetPrograms(ctx)
}

func (s *Service) GetCivilStatusTypes(
	ctx context.Context,
) ([]CivilStatusType, error) {
	return s.repo.GetCivilStatusTypes(ctx)
}

func (s *Service) GetReligions(ctx context.Context) ([]Religion, error) {
	return s.repo.GetReligions(ctx)
}

func (s *Service) GetNatureOfResidenceTypes(
	ctx context.Context,
) ([]NatureOfResidenceType, error) {
	return s.repo.GetNatureOfResidenceTypes(ctx)
}

func (s *Service) GetActivityOptions(
	ctx context.Context,
) ([]ActivityOption, error) {
	return s.repo.GetActivityOptions(ctx)
}

func (s *Service) GetStudentRelationshipTypes(
	ctx context.Context,
) ([]StudentRelationshipType, error) {
	return s.repo.GetStudentRelationshipTypes(ctx)
}

func (s *Service) GetStudentStatuses(
	ctx context.Context,
) ([]StudentStatus, error) {
	return s.repo.GetStudentStatuses(ctx)
}

func (s *Service) ListStudents(
	ctx context.Context, req ListStudentsRequest,
) (*ListStudentsResponse, error) {
	req.SetDefaults("last_name")

	students, err := s.repo.ListStudents(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	studentDTOs := make([]StudentProfileDTO, len(students))
	for i, st := range students {
		studentDTOs[i] = st.ToDTO()
	}

	total, err := s.repo.GetTotalStudentsCount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get total students count: %w",
			err,
		)
	}

	filterCounts, err := s.repo.GetStudentFilterCounts(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student filter counts: %w",
			err,
		)
	}

	return &ListStudentsResponse{
		Students: studentDTOs,
		Meta: structs.CalculateMetadata(
			total,
			req.Page,
			req.PageSize,
		),
		FilterCounts: filterCounts,
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
	if profile.StudentCORURL != "" {
		_, corErr := s.repo.GetStudentCORByUserID(ctx, iir.UserID)
		profile.IsStudentCORValid = corErr == nil
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		basicInfo, err := s.GetStudentBasicInfo(ctx, iirID)
		if err != nil {
			return err
		}
		if basicInfo != nil {
			profile.Student.BasicInfo = *basicInfo
		}
		return nil
	})

	g.Go(func() error {
		personalInfo, err := s.GetStudentPersonalInfo(ctx, iirID)
		if err != nil {
			return err
		}
		if personalInfo != nil {
			profile.Student.StudentPersonalInfoDTO = *personalInfo
		}
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
		if education != nil {
			profile.Education = *education
		}
		return nil
	})

	g.Go(func() error {
		familyBackground, err := s.GetStudentFamilyBackground(ctx, iirID)
		if err != nil {
			return err
		}
		if familyBackground != nil {
			profile.Family.FamilyBackgroundDTO = *familyBackground
		}
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
		if finance != nil {
			profile.Family.Finance = *finance
		}
		return nil
	})

	g.Go(func() error {
		healthRecord, err := s.GetStudentHealthRecord(ctx, iirID)
		if err != nil {
			return err
		}
		if healthRecord != nil {
			profile.Health.StudentHealthRecordDTO = *healthRecord
		}
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
	view, err := s.repo.GetStudentPersonalInfoView(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get personal info view: %w", err)
	}
	if view == nil {
		return nil, nil
	}

	emergencyAddressDTO, _ := s.locationsSvc.GetAddressByID(
		ctx,
		view.EmergencyAddressID,
	)

	return view.ToDTO(emergencyAddressDTO), nil
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
		addresses = append(addresses, addr.ToDTO(addrDTO))
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

	parentalStatus, err := s.repo.GetParentalStatusByID(
		ctx,
		studentFamily.ParentalStatusID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get parental status: %w",
			err,
		)
	}

	natureOfResidence, err := s.repo.GetNatureOfResidenceByID(
		ctx,
		studentFamily.NatureOfResidenceId,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get nature of residence: %w",
			err,
		)
	}

	siblingSupportTypes, err := s.repo.GetStudentSiblingSupport(
		ctx,
		studentFamily.ID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get sibling support types: %w",
			err,
		)
	}

	supportTypes := make([]SibilingSupportType, 0)
	for _, sst := range siblingSupportTypes {
		st, err := s.repo.GetSiblingSupportTypeByID(ctx, sst.SupportTypeID)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get sibling support type: %w",
				err,
			)
		}
		supportTypes = append(supportTypes, *st)
	}

	fbDTO := studentFamily.ToDTO(
		*parentalStatus,
		*natureOfResidence,
		supportTypes,
	)
	return &fbDTO, nil
}

func (s *Service) GetStudentRelatedPersons(
	ctx context.Context,
	iirID string,
) ([]RelatedPersonDTO, error) {
	views, err := s.repo.GetStudentRelatedPersonsView(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get related persons: %w", err)
	}

	related := make([]RelatedPersonDTO, len(views))
	for i, view := range views {
		related[i] = view.ToDTO()
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
	view, err := s.repo.GetStudentFinancialInfoView(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get financial info view: %w", err)
	}
	if view == nil {
		return nil, nil
	}

	supportDTOs, err := s.repo.GetFinancialSupportTypes(ctx, view.ID)
	if err != nil {
		supportDTOs = []StudentSupportType{}
	}

	fbDTO := view.ToDTO(supportDTOs)
	return &fbDTO, nil
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

	hrDTO := hr.ToDTO()
	return &hrDTO, nil
}

func (s *Service) GetStudentConsultations(
	ctx context.Context,
	iirID string,
) ([]StudentConsultationDTO, error) {
	consultations, err := s.repo.GetStudentConsultations(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consultations: %w", err)
	}

	dtos := make([]StudentConsultationDTO, len(consultations))
	for i, c := range consultations {
		dtos[i] = c.ToDTO()
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

	dtos := make([]StudentActivityDTO, len(activities))
	for i, a := range activities {
		option, err := s.repo.GetActivityOptionByID(ctx, a.OptionID)
		if err != nil || option == nil {
			dtos[i] = StudentActivityDTO{
				ID: a.ID,
				ActivityOption: ActivityOption{
					ID:       a.OptionID,
					Name:     "Unknown activity",
					Category: "",
					IsActive: false,
				},
				OtherSpecification: a.OtherSpecification,
				Role:               a.Role,
				RoleSpecification:  a.RoleSpecification,
			}
			continue
		}

		dtos[i] = a.ToDTO(*option)
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

	dtos := make([]StudentSubjectPreferenceDTO, len(prefs))
	for i, p := range prefs {
		dtos[i] = p.ToDTO()
	}
	return dtos, nil
}

func (s *Service) GetStudentHobbies(
	ctx context.Context,
	iirID string,
) ([]StudentHobbyDTO, error) {
	hobbies, err := s.repo.GetStudentHobbies(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hobbies: %w", err)
	}

	dtos := make([]StudentHobbyDTO, len(hobbies))
	for i, h := range hobbies {
		dtos[i] = h.ToDTO()
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

	dtos := make([]TestResultDTO, len(results))
	for i, r := range results {
		dtos[i] = r.ToDTO()
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

	// Clean up orphaned addresses asynchronously to avoid lock contention.
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.repo.DeleteOrphanedAddresses(bgCtx, s.repo.GetDB())
	}()

	return iirID, nil
}

func (s *Service) UpdateStudentIIR(
	ctx context.Context,
	iirID string,
	req ComprehensiveProfileDTO,
) (string, error) {
	existing, err := s.repo.GetStudentIIR(ctx, iirID)
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {UpdateStudentIIR Find}: %w",
			err,
		)
	}
	if existing == nil {
		return "", fmt.Errorf("IIR record not found")
	}

	req.IIRID = iirID

	err = s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		_, err := s.saveComprehensiveProfile(ctx, tx, existing.UserID, req)
		return err
	})
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {UpdateStudentIIR Save}: %w",
			err,
		)
	}

	// Clean up orphaned addresses asynchronously to avoid lock contention.
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.repo.DeleteOrphanedAddresses(bgCtx, s.repo.GetDB())
	}()

	// Notify counselors when student updates their IIR details
	counselorIDs, err := s.userService.GetUserIDsByRole(
		ctx,
		int(constants.AdminRoleID),
	)
	if err == nil && len(counselorIDs) > 0 {
		student, _ := s.userService.GetUserByID(ctx, existing.UserID)
		studentName := "A student"
		if student != nil {
			studentName = fmt.Sprintf(
				"%s %s",
				student.FirstName,
				student.LastName,
			)
		}

		notifications := make([]audit.NotificationParams, 0, len(counselorIDs))
		for _, cid := range counselorIDs {
			notifications = append(notifications, audit.NotificationParams{
				ReceiverID: structs.StringToNullableString(cid),
				TargetID:   structs.StringToNullableString(iirID),
				TargetType: structs.StringToNullableString(
					constants.IIREntityType,
				),
				Title: "Student IIR Details Updated",
				Message: fmt.Sprintf(
					"%s updated their Individual Inventory Record details.",
					studentName,
				),
				Type: constants.IIREntityType,
			})
		}

		audit.Dispatch(
			ctx,
			s.logService,
			s.notifService,
			nil,
			audit.DispatchParams{
				Notifications: notifications,
			},
		)
	}

	return iirID, nil
}

func (s *Service) validateDate(dateStr string, fieldName string) error {
	if dateStr == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	manilaOffset := 8 * 60 * 60
	loc := time.FixedZone("PHT", manilaOffset)
	t, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		return fmt.Errorf(
			"invalid date format for %s: must be YYYY-MM-DD",
			fieldName,
		)
	}
	if t.After(time.Now().In(loc)) {
		return fmt.Errorf("%s cannot be in the future", fieldName)
	}
	return nil
}

func (s *Service) saveComprehensiveProfile(
	ctx context.Context,
	tx datastore.DB,
	userID string,
	req ComprehensiveProfileDTO,
) (string, error) {
	iirID := req.IIRID
	if iirID == "" {
		var existingID string
		err := tx.GetContext(
			ctx,
			&existingID,
			"SELECT id FROM iir_records WHERE user_id = ? LIMIT 1",
			userID,
		)
		if err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
		if existingID != "" {
			iirID = existingID
		} else {
			iirID = uuid.New().String()
		}
	}

	isCompleted := true
	if req.IsCompleted != nil {
		isCompleted = *req.IsCompleted
	}

	if ok, err := s.repo.ValidateCivilStatusExists(
		ctx, tx, req.Student.CivilStatus.ID,
	); err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w", err,
		)
	} else if !ok {
		return "", &ValidationError{
			Message: fmt.Sprintf(
				"invalid civil status ID: %d",
				req.Student.CivilStatus.ID,
			),
		}
	}

	if ok, err := s.repo.ValidateReligionExists(
		ctx, tx, req.Student.Religion.ID,
	); err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w", err,
		)
	} else if !ok {
		return "", &ValidationError{
			Message: fmt.Sprintf(
				"invalid religion ID: %d", req.Student.Religion.ID,
			),
		}
	}

	if ok, err := s.repo.ValidateProgramExists(
		ctx, tx, req.Student.Program.ID,
	); err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w", err,
		)
	} else if !ok {
		return "", &ValidationError{
			Message: fmt.Sprintf(
				"invalid program ID: %d", req.Student.Program.ID,
			),
		}
	}

	if isCompleted {
		if ok, err := s.repo.ValidateParentalStatusExists(
			ctx, tx, req.Family.ParentalStatus.ID,
		); err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w", err,
			)
		} else if !ok {
			return "", &ValidationError{
				Message: fmt.Sprintf(
					"invalid parental status ID: %d",
					req.Family.ParentalStatus.ID,
				),
			}
		}

		if ok, err := s.repo.ValidateNatureOfResidenceExists(
			ctx, tx, req.Family.NatureOfResidence.ID,
		); err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w", err,
			)
		} else if !ok {
			return "", &ValidationError{
				Message: fmt.Sprintf(
					"invalid nature of residence ID: %d",
					req.Family.NatureOfResidence.ID,
				),
			}
		}
	}

	// Validate Critical Dates
	if err := s.validateDate(
		req.Student.DateOfBirth,
		"Student Date of Birth",
	); err != nil {
		return "", &ValidationError{
			Message: err.Error(),
		}
	}

	_, err := s.repo.UpsertIIRRecord(ctx, tx, &IIRRecord{
		ID:          iirID,
		UserID:      userID,
		IsSubmitted: true,
		IsCompleted: isCompleted,
	})
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
	}

	// 2. Personal Info
	err = s.repo.UpsertStudentPersonalInfo(ctx, tx, &StudentPersonalInfo{
		IIRID:         iirID,
		StudentNumber: req.Student.StudentNumber,
		Gender: func() string {
			if req.Student.Gender.ID == 2 {
				return "Female"
			}
			return "Male"
		}(),
		CivilStatusID:         req.Student.CivilStatus.ID,
		ReligionID:            req.Student.Religion.ID,
		OtherReligionText:     req.Student.OtherReligionText,
		HeightM:               req.Student.HeightM,
		WeightKg:              req.Student.WeightKg,
		Complexion:            req.Student.Complexion,
		HighSchoolGWA:         req.Student.HighSchoolGWA,
		ProgramID:             req.Student.Program.ID,
		YearLevel:             req.Student.YearLevel,
		Section:               req.Student.Section,
		PlaceOfBirth:          req.Student.PlaceOfBirth,
		DateOfBirth:           req.Student.DateOfBirth,
		MobileNumber:          req.Student.MobileNumber,
		TelephoneNumber:       req.Student.TelephoneNumber,
		IsEmployed:            req.Student.IsEmployed,
		EmployerName:          req.Student.EmployerName,
		EmployerAddress:       req.Student.EmployerAddress,
		EmployerContactNumber: req.Student.EmployerContactNumber,

		StatusID: 1, // Default to Enrolled/Active?
	})
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
	}

	// 3. Emergency Contact
	err = s.repo.DeleteEmergencyContactByIIRID(ctx, tx, iirID)
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
	}

	ecAddrID, err := s.locationsSvc.UpsertAddress(
		ctx,
		tx,
		req.Student.EmergencyContact.Address,
	)
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
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

	// Delete all existing student addresses before reinserting from payload
	err = s.repo.DeleteStudentAddressesByIIRID(ctx, tx, iirID)
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
	}

	for _, addrDTO := range req.Student.Addresses {
		addrID, err := s.locationsSvc.UpsertAddress(ctx, tx, addrDTO.Address)
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
		_, err = s.repo.UpsertStudentAddress(ctx, tx, &StudentAddress{
			IIRID:       iirID,
			AddressID:   addrID,
			AddressType: addrDTO.AddressType,
		})
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	if !isCompleted {
		return iirID, nil
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
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
	}

	_ = s.repo.DeleteStudentSiblingSupportsByFamilyID(ctx, tx, fbID)
	for _, sst := range req.Family.SiblingSupportTypes {
		err = s.repo.CreateStudentSiblingSupport(
			ctx,
			tx,
			&StudentSiblingSupport{
				FamilyBackgroundID: fbID,
				SupportTypeID:      sst.ID,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	// 6. Educational Background
	ebID, err := s.repo.UpsertEducationalBackground(
		ctx,
		tx,
		&EducationalBackground{
			IIRID:              iirID,
			NatureOfSchooling:  req.Education.NatureOfSchooling,
			InterruptedDetails: req.Education.InterruptedDetails,
		})
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
	}

	_ = s.repo.DeleteSchoolDetailsByEBID(ctx, tx, ebID)
	for _, sd := range req.Education.School {
		if strings.TrimSpace(sd.SchoolName) == "" {
			continue
		}
		_, err = s.repo.UpsertSchoolDetails(ctx, tx, &SchoolDetails{
			EducationalLevelID: sd.EducationalLevel.ID,
			SchoolName:         sd.SchoolName,
			SchoolAddress:      sd.SchoolAddress,
			SchoolType:         sd.SchoolType,
			YearStarted:        sd.YearStarted,
			YearCompleted:      sd.YearCompleted,
			Awards:             sd.Awards,
			EBID:               ebID,
		})
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	// 7. Health Record
	_, err = s.repo.UpsertStudentHealthRecord(ctx, tx, &StudentHealthRecord{
		IIRID:                     iirID,
		VisionHasProblem:          req.Health.VisionHasProblem,
		VisionDetails:             req.Health.VisionDetails,
		HearingHasProblem:         req.Health.HearingHasProblem,
		HearingDetails:            req.Health.HearingDetails,
		SpeechHasProblem:          req.Health.SpeechHasProblem,
		SpeechDetails:             req.Health.SpeechDetails,
		GeneralHealthHasProblem:   req.Health.GeneralHealthHasProblem,
		GeneralHealthDetails:      req.Health.GeneralHealthDetails,
		MentalEmotionalHasProblem: req.Health.MentalEmotionalHasProblem,
		MentalEmotionalDetails:    req.Health.MentalEmotionalDetails,
	})
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
	}

	// 8. Finance
	sfID, err := s.repo.UpsertStudentFinance(ctx, tx, &StudentFinance{
		IIRID:           iirID,
		IncomeRangeID:   req.Family.Finance.IncomeRange.ID,
		OtherIncome:     req.Family.Finance.OtherIncomeDetails,
		WeeklyAllowance: req.Family.Finance.WeeklyAllowance,
	})
	if err != nil {
		return "", fmt.Errorf(
			"[StudentService] {saveComprehensiveProfile}: %w",
			err,
		)
	}

	_ = s.repo.DeleteStudentFinancialSupportsByFinanceID(ctx, tx, sfID)
	for _, st := range req.Family.Finance.FinancialSupportTypes {
		err = s.repo.CreateStudentFinancialSupport(
			ctx,
			tx,
			&StudentFinancialSupport{
				StudentFinanceID: sfID,
				SupportTypeID:    st.ID,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	// 9. Related Persons
	_ = s.repo.DeleteStudentRelatedPersons(ctx, tx, iirID)
	hasGuardian := false
	for _, rpDTO := range req.Family.RelatedPersons {
		if rpDTO.IsGuardian {
			hasGuardian = true
			if strings.TrimSpace(rpDTO.FirstName) == "" ||
				strings.TrimSpace(rpDTO.LastName) == "" {
				return "", &ValidationError{
					Message: "guardian first and last name are required",
				}
			}
		}
	}
	if !hasGuardian {
		return "", &ValidationError{
			Message: "guardian information is required",
		}
	}

	for _, rpDTO := range req.Family.RelatedPersons {
		if !rpDTO.IsGuardian &&
			strings.TrimSpace(rpDTO.FirstName) == "" &&
			strings.TrimSpace(rpDTO.LastName) == "" {
			continue
		}
		// Validate DOB if not empty or if required
		if rpDTO.DateOfBirth != "" {
			if err := s.validateDate(
				rpDTO.DateOfBirth,
				"Related Person Date of Birth",
			); err != nil {
				return "", &ValidationError{
					Message: err.Error(),
				}
			}
		} else {
			// If it's mandatory in DB, we MUST provide a valid date.
			return "", &ValidationError{
				Message: fmt.Sprintf(
					"date of birth is required for related person: %s %s",
					rpDTO.FirstName,
					rpDTO.LastName,
				),
			}
		}

		eduID := structs.Int64ToNullableInt64(
			int64(rpDTO.EducationalAttainment.ID),
		)

		rpID, err := s.repo.UpsertRelatedPerson(
			ctx, tx, &RelatedPerson{
				FirstName:               rpDTO.FirstName,
				MiddleName:              rpDTO.MiddleName,
				LastName:                rpDTO.LastName,
				SuffixName:              rpDTO.SuffixName,
				DateOfBirth:             rpDTO.DateOfBirth,
				EducationalAttainmentID: eduID,
				Occupation:              rpDTO.Occupation,
				EmployerName:            rpDTO.EmployerName,
				EmployerAddress:         rpDTO.EmployerAddress,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}

		err = s.repo.UpsertStudentRelatedPerson(
			ctx, tx, &StudentRelatedPerson{
				IIRID:           iirID,
				RelatedPersonID: rpID,
				RelationshipID:  rpDTO.Relationship.ID,
				IsParent:        rpDTO.IsParent,
				IsGuardian:      rpDTO.IsGuardian,
				IsLiving:        rpDTO.IsLiving,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	// 10. Consultations
	_ = s.repo.DeleteStudentConsultationsByIIRID(ctx, tx, iirID)
	for _, cDTO := range req.Health.Consultations {
		_, err = s.repo.UpsertStudentConsultation(ctx, tx, &StudentConsultation{
			IIRID:            iirID,
			ProfessionalType: cDTO.ProfessionalType,
			HasConsulted:     cDTO.HasConsulted,
			WhenDate:         cDTO.WhenDate,
			ForWhat:          cDTO.ForWhat,
		})
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	// 11. Activities
	_ = s.repo.DeleteStudentActivitiesByIIRID(ctx, tx, iirID)
	for _, aDTO := range req.Interests.Activities {
		_, err = s.repo.CreateStudentActivity(ctx, tx, &StudentActivity{
			IIRID:              iirID,
			OptionID:           aDTO.ActivityOption.ID,
			OtherSpecification: aDTO.OtherSpecification,
			Role:               aDTO.Role,
			RoleSpecification:  aDTO.RoleSpecification,
		})
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	// 12. Subject Preferences
	_ = s.repo.DeleteStudentSubjectPreferencesByIIRID(ctx, tx, iirID)
	seenSubjects := make(map[string]bool)
	for _, sspDTO := range req.Interests.SubjectPreferences {
		subjName := strings.TrimSpace(sspDTO.SubjectName)
		if subjName == "" {
			continue
		}

		// Case-insensitive deduplication
		key := strings.ToLower(subjName)
		if seenSubjects[key] {
			continue
		}
		seenSubjects[key] = true

		_, err = s.repo.CreateStudentSubjectPreference(
			ctx,
			tx,
			&StudentSubjectPreference{
				IIRID:       iirID,
				SubjectName: subjName,
				IsFavorite:  sspDTO.IsFavorite,
			})
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	// 13. Hobbies
	_ = s.repo.DeleteStudentHobbiesByIIRID(ctx, tx, iirID)
	seenHobbies := make(map[string]bool)
	for _, hDTO := range req.Interests.Hobbies {
		hName := strings.TrimSpace(hDTO.HobbyName)
		if hName == "" {
			continue
		}

		// Case-insensitive deduplication
		key := strings.ToLower(hName)
		if seenHobbies[key] {
			continue
		}
		seenHobbies[key] = true

		_, err = s.repo.CreateStudentHobby(ctx, tx, &StudentHobby{
			IIRID:        iirID,
			HobbyName:    hName,
			PriorityRank: hDTO.PriorityRank,
		})
		if err != nil {
			return "", fmt.Errorf(
				"[StudentService] {saveComprehensiveProfile}: %w",
				err,
			)
		}
	}

	return iirID, nil
}

func (s *Service) GenerateIIR(
	ctx context.Context,
	iirID string,
	isCounselor bool,
) ([]byte, string, error) {
	profile, err := s.GetStudentProfile(ctx, iirID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get student profile: %w", err)
	}

	iir, err := s.repo.GetStudentIIR(ctx, iirID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch IIR record: %w", err)
	}

	dateStr := s.GetFormattedDate(iir.UpdatedAt.Format("2006-01-02"))
	if dateStr == "" {
		dateStr = s.GetFormattedDate(time.Now().Format("2006-01-02"))
	}

	var notes []SignificantNote
	if isCounselor {
		notes, err = s.repo.GetStudentSignificantNotes(ctx, iirID)
		if err != nil {
			return nil, "", fmt.Errorf(
				"failed to fetch significant notes: %w",
				err,
			)
		}
	}

	data := struct {
		Profile              *ComprehensiveProfileDTO
		RecentDate           string
		ShowSignificantNotes bool
		SignificantNotes     []SignificantNote
	}{
		Profile:              profile,
		RecentDate:           dateStr,
		ShowSignificantNotes: isCounselor,
		SignificantNotes:     notes,
	}

	pdfBytes, err := s.pdfService.GenerateFromContent(
		ctx,
		"iirf.html",
		iirTemplate,
		data,
	)
	if err != nil {
		return nil, "", fmt.Errorf(
			"failed to generate pdf from template: %w",
			err,
		)
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

func (s *Service) ExportStudentsCSV(
	ctx context.Context,
	req ListStudentsRequest,
) ([]byte, error) {
	// Apply the same sort defaults as the paginated list, but never apply
	// pagination: an export must include the complete matching dataset.
	req.SetDefaults("last_name")

	students, err := s.repo.ListAllStudents(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch students for export: %w", err)
	}

	studentDTOs := make([]StudentProfileDTO, len(students))
	for i, student := range students {
		studentDTOs[i] = student.ToDTO()
	}

	return s.generateStudentsCSV(studentDTOs)
}

func (s *Service) generateStudentsCSV(students []StudentProfileDTO) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{
		"Student Number",
		"Last Name",
		"First Name",
		"Email",
		"Program",
		"Year Level",
		"Status",
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("error writing csv headers: %w", err)
	}

	for _, st := range students {

		row := []string{
			st.StudentNumber,
			st.LastName,
			st.FirstName,
			st.Email,
			st.Program.Code,
			fmt.Sprintf("%d", st.YearLevel),
			st.Status.Name,
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("error writing csv row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing csv writer: %w", err)
	}

	return buf.Bytes(), nil
}
