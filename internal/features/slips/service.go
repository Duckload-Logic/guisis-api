package slips

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/datetime"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/files"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/ocr"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

const MaxFileSize = 5 * 1024 * 1024 // 5MB limit
type Service struct {
	repo           *Repository
	logService     audit.Logger
	notifService   audit.Notifier
	emailService   audit.Emailer
	fileStorage    storage.FileStorage
	userService    *users.Service
	studentService *students.Service
	filesService   *files.Service
	ocrClient      *ocr.OCRClient
	cfg            *config.Config
}

func NewService(
	repo *Repository,
	logService audit.Logger,
	notifService audit.Notifier,
	emailService audit.Emailer,
	fileStorage storage.FileStorage,
	userService *users.Service,
	studentService *students.Service,
	filesService *files.Service,
	ocrClient *ocr.OCRClient,
	cfg *config.Config,
) *Service {
	return &Service{
		repo:           repo,
		logService:     logService,
		notifService:   notifService,
		emailService:   emailService,
		fileStorage:    fileStorage,
		userService:    userService,
		studentService: studentService,
		filesService:   filesService,
		ocrClient:      ocrClient,
		cfg:            cfg,
	}
}

func (s *Service) GetSlipStatuses(
	ctx context.Context,
) ([]SlipStatus, error) {
	statuses, err := s.repo.GetSlipStatuses(ctx)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (s *Service) GetSlipCategories(
	ctx context.Context,
) ([]SlipCategory, error) {
	categories, err := s.repo.GetSlipCategories(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) GetSlipByID(
	ctx context.Context,
	id string,
) (*SlipDTO, error) {
	slip, err := s.repo.GetSlipByIDWithDetails(ctx, s.repo.GetDB(), id)
	if err != nil {
		return nil, err
	}
	if slip == nil {
		return nil, fmt.Errorf("slip not found")
	}

	dto := s.mapToDTO(slip)

	hasNote, _ := s.repo.HasNoteForAdmissionSlip(ctx, id)
	dto.HasSignificantNote = hasNote

	corMap, _ := s.studentService.GetLatestCORsByUserIDs(
		ctx,
		[]string{slip.UserID},
	)
	dto.StudentCORURL = corMap[slip.UserID]

	return dto, nil
}

func (s *Service) GetUrgentSlips(
	ctx context.Context,
	req *ListSlipsRequest,
) (*ListSlipsResponse, error) {
	req.SetDefaults("urgency_score")

	slips, err := s.repo.GetUrgentSlips(ctx, req)
	if err != nil {
		return nil, err
	}

	// Batch fetch COR URLs
	userIDs := make([]string, 0, len(slips))
	for i := range slips {
		if slips[i].UserID != "" {
			userIDs = append(userIDs, slips[i].UserID)
		}
	}
	corMap, _ := s.studentService.GetLatestCORsByUserIDs(ctx, userIDs)

	var slipDTOs []SlipDTO
	for i := range slips {
		dto := s.mapToDTO(&slips[i])
		dto.StudentCORURL = corMap[slips[i].UserID]
		slipDTOs = append(slipDTOs, *dto)
	}

	req.StatusID = 1
	total, err := s.repo.GetTotalUrgentSlipsCount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get slips count: %w", err)
	}

	return &ListSlipsResponse{
		Slips: slipDTOs,
		Meta:  structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetSlipStats(
	ctx context.Context,
	iirID *string,
	req *ListSlipsRequest,
) ([]SlipStatusCount, error) {
	stats, err := s.repo.GetSlipStats(ctx, iirID, req)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *Service) GetAllExcuseSlips(
	ctx context.Context,
	req ListSlipsRequest,
) (*ListSlipsResponse, error) {
	req.SetDefaults("date_needed")
	if req.SortBy != "" {
		req.OrderBy = req.SortBy
	}

	slips, err := s.repo.GetAll(ctx, &req)
	if err != nil {
		return nil, err
	}

	// Batch fetch COR URLs
	userIDs := make([]string, 0, len(slips))
	for i := range slips {
		if slips[i].UserID != "" {
			userIDs = append(userIDs, slips[i].UserID)
		}
	}
	corMap, _ := s.studentService.GetLatestCORsByUserIDs(ctx, userIDs)

	var slipDTOs []SlipDTO
	for i := range slips {
		dto := s.mapToDTO(&slips[i])
		dto.StudentCORURL = corMap[slips[i].UserID]
		slipDTOs = append(slipDTOs, *dto)
	}

	total, err := s.repo.GetTotalSlipsCount(ctx, &req, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get slips count: %w", err)
	}

	return &ListSlipsResponse{
		Slips: slipDTOs,
		Meta:  structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetExcuseSlipsByIIRID(
	ctx context.Context,
	iirID string,
	req ListSlipsRequest,
) (*ListSlipsResponse, error) {
	req.SetDefaults("date_needed")
	if req.SortBy != "" {
		req.OrderBy = req.SortBy
	}

	slips, err := s.repo.GetByIIRID(ctx, iirID, &req)
	if err != nil {
		return nil, err
	}

	var slipDTOs []SlipDTO
	for i := range slips {
		dto := s.mapToDTO(&slips[i])
		slipDTOs = append(slipDTOs, *dto)
	}

	total, err := s.repo.GetTotalSlipsCount(ctx, &req, &iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get slips count: %w",
			err,
		)
	}

	return &ListSlipsResponse{
		Slips: slipDTOs,
		Meta:  structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetSlipAttachments(
	ctx context.Context,
	slipID string,
) ([]AttachmentDTO, error) {
	attachments, err := s.repo.GetSlipAttachments(ctx, slipID)
	if err != nil {
		return nil, err
	}

	var attachmentDTOs []AttachmentDTO
	for a := range attachments {
		slipID := ""
		if attachments[a].SlipID.Valid {
			slipID = attachments[a].SlipID.String
		}

		attachmentDTOs = append(attachmentDTOs, AttachmentDTO{
			ID:             attachments[a].FileID,
			SlipID:         slipID,
			FileName:       attachments[a].FileName,
			FileURL:        attachments[a].FileURL,
			FileType:       attachments[a].FileType,
			FileSize:       attachments[a].FileSize,
			MimeType:       attachments[a].MimeType,
			AttachmentType: attachments[a].AttachmentType,
		})
	}

	return attachmentDTOs, nil
}

func (s *Service) validateFiles(files []*multipart.FileHeader) error {
	allowedTypes := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	for _, file := range files {
		// Check File Size
		if file.Size > MaxFileSize {
			return fmt.Errorf(
				"file '%s' is too large: maximum 5MB allowed",
				file.Filename,
			)
		}

		// Check File Type (Content-Aware)
		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s", file.Filename)
		}
		defer f.Close()

		// Read first 512 bytes to detect content type
		buffer := make([]byte, 512)
		_, _ = f.Read(buffer)
		contentType := http.DetectContentType(buffer)

		allowedMime := map[string]bool{
			"application/pdf": true,
			"image/jpeg":      true,
			"image/png":       true,
		}

		if !allowedMime[contentType] {
			return fmt.Errorf(
				"invalid content type for '%s': "+
					"expected PDF or Image, got %s",
				file.Filename,
				contentType,
			)
		}

		// Double check extension just in case
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if !allowedTypes[ext] {
			return fmt.Errorf(
				"invalid file extension for '%s'",
				file.Filename,
			)
		}
	}
	return nil
}

func (s *Service) validateFilesOCR(
	ctx context.Context,
	files []*multipart.FileHeader,
) error {
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return fmt.Errorf(
				"failed to open file '%s' for validation: %w",
				file.Filename,
				err,
			)
		}

		res, err := s.ocrClient.ValidateDocument(ctx, file.Filename, f)
		f.Close()
		if err != nil {
			return fmt.Errorf("OCR validation service is currently offline")
		}
		if !res.IsValid {
			return fmt.Errorf(
				"file '%s' is invalid: %s",
				file.Filename,
				res.Message,
			)
		}
	}
	return nil
}

// SubmitExcuseSlip creates a new slip with attachments.
func (s *Service) SubmitExcuseSlip(
	ctx context.Context,
	iirID string,
	req CreateSlipRequest,
	files []*multipart.FileHeader,
	parentIdFiles []*multipart.FileHeader,
) (*SlipDTO, error) {
	allFiles := append([]*multipart.FileHeader{}, files...)
	allFiles = append(allFiles, parentIdFiles...)

	// Validate all files
	if err := s.validateFiles(allFiles); err != nil {
		return nil, err
	}

	// OCR check on non-ID files
	if err := s.validateFilesOCR(ctx, files); err != nil {
		return nil, err
	}

	dateOfAbsence := strings.Split(req.DateOfAbsence, "T")[0]
	parsedDate, err := time.Parse("2006-01-02", dateOfAbsence)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid date format: YYYY-MM-DD",
		)
	}

	// Evaluate current day in Philippine Time (UTC+8)
	loc := time.FixedZone("PHT", 8*3600)
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if parsedDate.After(today) {
		return nil, fmt.Errorf("absence date cannot be in future")
	}

	dateNeeded := strings.Split(req.DateNeeded, "T")[0]
	parsedDateNeeded, err := time.Parse("2006-01-02", dateNeeded)
	if err != nil {
		return nil, fmt.Errorf("invalid date needed format: YYYY-MM-DD")
	}
	if parsedDateNeeded.Before(today) {
		return nil, fmt.Errorf("date needed cannot be in the past")
	}

	// Check for duplicate active slip
	exists, err := s.repo.HasActiveSlipForDate(ctx, iirID, dateOfAbsence, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf(
			"an active excuse slip already exists for this absence date",
		)
	}


	// Unified File Implementation: Use files features
	uploadedFiles, err := s.filesService.UploadFiles(ctx, allFiles, "slips")
	if err != nil {
		return nil, fmt.Errorf("failed to upload files: %w", err)
	}

	slip := &Slip{
		ID:            uuid.New().String(),
		IIRID:         iirID,
		Reason:        req.Reason,
		DateOfAbsence: req.DateOfAbsence,
		DateNeeded:    req.DateNeeded,
		CategoryID:    req.CategoryID,
		StatusID:      1,
	}

	var createdSlip *SlipWithDetailsView
	err = s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			var err error
			createdSlip, err = s.repo.CreateSlip(ctx, tx, slip)
			if err != nil {
				return err
			}

			// Loop to create attachment records linked to files table
			for _, f := range uploadedFiles {
				attachment := &SlipAttachment{
					FileID:         f.ID,
					SlipID:         structs.StringToNullableString(slip.ID),
					AttachmentType: "OTHER",
				}
				if err := s.repo.SaveSlipAttachment(
					ctx,
					tx,
					attachment,
				); err != nil {
					return err
				}
			}
			return nil
		},
	)

	if err != nil {
		audit.Dispatch(
			ctx,
			s.logService,
			s.notifService,
			s.emailService,
			audit.DispatchParams{
				Log: &audit.LogParams{
					Level:    audit.LevelError,
					Category: audit.CategoryAudit,
					Action:   audit.ActionSlipFailed,
					Message: fmt.Sprintf(
						"Failed to create slip for IIR #%s",
						iirID,
					),
					Metadata: &audit.LogMetadata{
						EntityType: constants.SlipEntityType,
						EntityID:   slip.ID,
						Error:      err.Error(),
					},
				},
				Notifications: []audit.NotificationParams{
					{
						Title: fmt.Sprintf(
							"Slip Creation Failed for IIR #%s",
							iirID,
						),
						Message: fmt.Sprintf(
							"An error occurred while creating the slip: %s",
							err.Error(),
						),
						Type: constants.SlipEntityType,
					},
				},
			})
		return nil, err
	}

	// Fetch personalized notification targets
	userID := audit.ExtractUserID(ctx)
	student, _ := s.userService.GetUserByID(ctx, userID)
	studentName := "A student"
	if student != nil {
		studentName = fmt.Sprintf("%s %s", student.FirstName, student.LastName)
	}

	counselorIDs, _ := s.userService.GetUserIDsByRole(
		ctx,
		int(constants.AdminRoleID),
	)

	notifications := []audit.NotificationParams{
		{
			ReceiverID: structs.StringToNullableString(userID),
			TargetID:   structs.StringToNullableString(slip.ID),
			TargetType: structs.StringToNullableString(
				constants.SlipEntityType,
			),
			Title:   "Admission Slip Submitted Successfully",
			Message: "Your admission slip request has been submitted.",
			Type:    constants.SlipEntityType,
		},
	}

	for _, cid := range counselorIDs {
		notifications = append(notifications, audit.NotificationParams{
			ReceiverID: structs.StringToNullableString(cid),
			TargetID:   structs.StringToNullableString(slip.ID),
			TargetType: structs.StringToNullableString(
				constants.SlipEntityType,
			),
			Title: "New Admission Slip Request",
			Message: fmt.Sprintf(
				"New admission slip request received from %s for %s.",
				studentName,
				slip.DateOfAbsence,
			),
			Type: constants.SlipEntityType,
		})
	}

	counselorEmails, _ := s.userService.GetEmailsByRole(
		ctx,
		int(constants.AdminRoleID),
	)

	newSlipDTO := s.mapToDTO(createdSlip)

	audit.Dispatch(
		ctx,
		s.logService,
		s.notifService,
		s.emailService,
		audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelInfo,
				Category: audit.CategoryAudit,
				Action:   audit.ActionSlipCreated,
				Message:  fmt.Sprintf("Excuse slip #%s created", slip.ID),
				Metadata: &audit.LogMetadata{
					EntityType: constants.SlipEntityType,
					EntityID:   slip.ID,
				},
			},
			Notifications: notifications,
			Email: []audit.EmailParams{
				{
					To:           counselorEmails,
					Subject:      "New Admission Slip Request",
					TemplatePath: "request.html",
					TemplateData: map[string]interface{}{
						"EntityType": constants.SlipEntityType,
						"StudentName": fmt.Sprintf(
							"%s %s",
							newSlipDTO.User.FirstName,
							newSlipDTO.User.LastName,
						),
						"Category": newSlipDTO.Category.Name,
						"DateOfAbsence": datetime.FormatDate(
							newSlipDTO.DateOfAbsence,
						),
						"DateNeeded": datetime.FormatDate(
							newSlipDTO.DateNeeded,
						),
						"Status":     newSlipDTO.Status.Name,
						"AdminNotes": nil,
					},
				},
			},
		})

	return newSlipDTO, nil
}

func (s *Service) UpdateExcuseSlip(
	ctx context.Context,
	iirID string,
	slipID string,
	req CreateSlipRequest,
	files []*multipart.FileHeader,
	parentIdFiles []*multipart.FileHeader,
) (*SlipDTO, error) {
	// Fetch existing slip and validate ownership/status
	existingSlip, err := s.repo.GetSlipByID(ctx, slipID)
	if err != nil {
		return nil, err
	}
	if existingSlip == nil {
		return nil, fmt.Errorf("slip not found")
	}
	if existingSlip.IIRID != iirID {
		return nil, fmt.Errorf("access denied")
	}

	// Only allow editing if status is Pending (1) or For Revision (9)
	if existingSlip.StatusID != 1 && existingSlip.StatusID != 9 {
		return nil, fmt.Errorf("cannot edit slip in current status")
	}

	allFiles := append([]*multipart.FileHeader{}, files...)
	allFiles = append(allFiles, parentIdFiles...)

	// Validate all files
	if err := s.validateFiles(allFiles); err != nil {
		return nil, err
	}

	// OCR check on non-ID files
	if err := s.validateFilesOCR(ctx, files); err != nil {
		return nil, err
	}

	dateOfAbsence := strings.Split(req.DateOfAbsence, "T")[0]
	parsedDate, err := time.Parse("2006-01-02", dateOfAbsence)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: YYYY-MM-DD")
	}

	// Evaluate current day in Philippine Time (UTC+8)
	loc := time.FixedZone("PHT", 8*3600)
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if parsedDate.After(today) {
		return nil, fmt.Errorf("absence date cannot be in future")
	}

	dateNeeded := strings.Split(req.DateNeeded, "T")[0]
	parsedDateNeeded, err := time.Parse("2006-01-02", dateNeeded)
	if err != nil {
		return nil, fmt.Errorf("invalid date needed format: YYYY-MM-DD")
	}
	if parsedDateNeeded.Before(today) {
		return nil, fmt.Errorf("date needed cannot be in the past")
	}

	// Check for duplicate active slip (excluding this current slip)
	exists, err := s.repo.HasActiveSlipForDate(
		ctx, iirID, dateOfAbsence, slipID,
	)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf(
			"an active excuse slip already exists for this absence date",
		)
	}


	// Delete old attachments NOT in KeepFileIDs
	oldAttachments, err := s.repo.GetSlipAttachments(ctx, slipID)
	if err == nil {
		for _, att := range oldAttachments {
			keep := false
			for _, keepID := range req.KeepFileIDs {
				if att.FileID == keepID {
					keep = true
					break
				}
			}
			if !keep {
				_ = s.filesService.DeleteFile(ctx, att.FileID)
			}
		}
	}

	// Upload new files using centralized service
	uploadedFiles, err := s.filesService.UploadFiles(ctx, allFiles, "slips")
	if err != nil {
		return nil, fmt.Errorf("failed to upload files: %w", err)
	}

	existingNotes := ""
	if existingSlip.AdminNotes.Valid {
		existingNotes = existingSlip.AdminNotes.String
	}

	updatedNotes := existingNotes
	if existingSlip.StatusID != 1 {
		formattedTime := time.Now().Format("2006-01-02 15:04:05")
		newLogEntry := fmt.Sprintf(
			"[%s] STATUS: PENDING\n"+
				"Remarks: Student updated/resubmitted the slip.",
			formattedTime,
		)
		if existingNotes != "" {
			updatedNotes = newLogEntry +
				"\n\n------------------------------\n\n" +
				existingNotes
		} else {
			updatedNotes = newLogEntry
		}
	}

	// Update database in transaction
	updatedSlip := &Slip{
		ID:            slipID,
		IIRID:         iirID,
		Reason:        req.Reason,
		DateOfAbsence: req.DateOfAbsence,
		DateNeeded:    req.DateNeeded,
		CategoryID:    req.CategoryID,
		StatusID:      1, // Reset to Pending
		AdminNotes: structs.StringToNullableString(
			updatedNotes,
		),
	}

	err = s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			// Delete old attachment records
			if err := s.repo.DeleteSlipAttachments(ctx, tx, slipID); err != nil {
				return err
			}
			// Update slip
			if err := s.repo.UpdateSlip(ctx, tx, updatedSlip); err != nil {
				return err
			}
			// Save new attachments
			for _, f := range uploadedFiles {
				attachment := &SlipAttachment{
					FileID:         f.ID,
					SlipID:         structs.StringToNullableString(slipID),
					AttachmentType: "OTHER",
				}
				if err := s.repo.SaveSlipAttachment(
					ctx, tx, attachment,
				); err != nil {
					return err
				}
			}
			// Save kept attachments
			for _, keepID := range req.KeepFileIDs {
				attachment := &SlipAttachment{
					FileID:         keepID,
					SlipID:         structs.StringToNullableString(slipID),
					AttachmentType: "OTHER",
				}
				if err := s.repo.SaveSlipAttachment(
					ctx, tx, attachment,
				); err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	studentUserID, _ := s.repo.GetUserIDBySlipID(ctx, slipID)
	counselorIDs, _ := s.userService.GetUserIDsByRole(
		ctx,
		int(constants.AdminRoleID),
	)

	notifications := []audit.NotificationParams{
		{
			ReceiverID: structs.StringToNullableString(studentUserID),
			Title:      "Slip Updated",
			Message: fmt.Sprintf(
				"Your slip #%s has been updated",
				slipID,
			),
			Type: constants.SlipEntityType,
		},
	}

	userID := audit.ExtractUserID(ctx)
	student, _ := s.userService.GetUserByID(ctx, userID)
	studentName := "A student"
	if student != nil {
		studentName = fmt.Sprintf(
			"%s %s",
			student.FirstName,
			student.LastName,
		)
	}

	for _, cid := range counselorIDs {
		notifications = append(
			notifications,
			audit.NotificationParams{
				ReceiverID: structs.StringToNullableString(cid),
				TargetID: structs.StringToNullableString(
					slipID,
				),
				TargetType: structs.StringToNullableString(
					constants.SlipEntityType,
				),
				Title: "Admission Slip Resubmitted",
				Message: fmt.Sprintf(
					"%s resubmitted slip #%s for review.",
					studentName,
					slipID,
				),
				Type: constants.SlipEntityType,
			},
		)
	}

	audit.Dispatch(
		ctx,
		s.logService,
		s.notifService,
		s.emailService,
		audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelInfo,
				Category: audit.CategoryAudit,
				Action:   audit.ActionSlipUpdated,
				Message: fmt.Sprintf(
					"Excuse slip #%s updated",
					slipID,
				),
				Metadata: &audit.LogMetadata{
					EntityType: constants.SlipEntityType,
					EntityID:   slipID,
				},
			},
			Notifications: notifications,
		},
	)

	// Fetch fully populated updated slip from DB
	fullUpdatedSlip, err := s.repo.GetSlipByIDWithDetails(
		ctx,
		s.repo.GetDB(),
		slipID,
	)
	if err != nil {
		return nil, err
	}

	return s.mapToDTO(fullUpdatedSlip), nil
}

// DownloadAttachment streams an attachment after validating that it belongs to the slip.
func (s *Service) DownloadAttachment(
	ctx context.Context,
	slipID string,
	attachmentID string,
	writer http.ResponseWriter,
) (*SlipAttachment, error) {
	attachment, err := s.repo.GetAttachmentByIDAndSlipID(
		ctx,
		slipID,
		attachmentID,
	)
	if err != nil {
		return nil, err
	}
	if attachment == nil {
		return nil, fmt.Errorf("attachment not found")
	}

	blobPath, err := normalizeAttachmentBlobPath(attachment.FileURL)
	if err != nil {
		return nil, err
	}

	var fileBuffer bytes.Buffer
	if err := s.fileStorage.Download(ctx, blobPath, &fileBuffer); err != nil {
		return nil, fmt.Errorf("attachment file not found in storage: %w", err)
	}

	if fileBuffer.Len() == 0 {
		return nil, fmt.Errorf("attachment file is empty in storage")
	}

	fileName := sanitizeDownloadFileName(attachment.FileName)
	contentType := attachment.MimeType
	if contentType == "" {
		contentType = mime.TypeByExtension(path.Ext(fileName))
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Expires", "0")
	writer.Header().Set(
		"Content-Disposition",
		mime.FormatMediaType("attachment", map[string]string{
			"filename": fileName,
		}),
	)
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileBuffer.Len()))

	if _, err := writer.Write(fileBuffer.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to write attachment response: %w", err)
	}

	return attachment, nil
}

func normalizeAttachmentBlobPath(fileURL string) (string, error) {
	raw := strings.TrimSpace(strings.ReplaceAll(fileURL, "\\", "/"))
	if raw == "" {
		return "", fmt.Errorf("security: invalid file path detected")
	}

	if parsed, err := url.Parse(raw); err == nil && parsed.Path != "" {
		raw = parsed.Path
	}

	raw = strings.TrimPrefix(raw, "/")
	raw = strings.TrimPrefix(raw, "uploads/")

	for _, env := range []string{"development/", "staging/", "production/"} {
		raw = strings.TrimPrefix(raw, env)
		raw = strings.TrimPrefix(raw, "uploads/")
	}

	cleanPath := path.Clean(raw)
	if cleanPath == "." ||
		strings.HasPrefix(cleanPath, "../") ||
		strings.Contains(cleanPath, "/../") ||
		!(strings.HasPrefix(cleanPath, "slips/") ||
			strings.HasPrefix(cleanPath, "cors/")) {
		return "", fmt.Errorf("security: invalid file path detected")
	}

	return cleanPath, nil
}

func sanitizeDownloadFileName(fileName string) string {
	cleanName := path.Base(strings.ReplaceAll(fileName, "\\", "/"))
	if cleanName == "." || cleanName == "/" || cleanName == "" {
		return "attachment"
	}

	return cleanName
}

func (s *Service) UpdateExcuseSlipStatus(
	ctx context.Context,
	id string,
	newStatus string,
	adminNotes string,
) error {
	// ... (validation code was here, I'll keep it)
	validStatuses := map[string]bool{
		"Pending":      true,
		"Approved":     true,
		"Rejected":     true,
		"For Revision": true,
	}

	if !validStatuses[newStatus] {
		return fmt.Errorf(
			"invalid status: must be 'Pending', 'Approved', " +
				"'Rejected', or 'For Revision'",
		)
	}

	// Fetch old state for audit trail
	oldSlip, _ := s.repo.GetSlipByIDWithDetails(
		ctx,
		s.repo.GetDB(),
		id,
	)

	existingNotes := ""
	if oldSlip != nil && oldSlip.AdminNotes.Valid {
		existingNotes = oldSlip.AdminNotes.String
	}

	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	newLogEntry := fmt.Sprintf(
		"[%s] STATUS: %s",
		formattedTime,
		newStatus,
	)
	trimmedNotes := strings.TrimSpace(adminNotes)
	if trimmedNotes != "" {
		newLogEntry = fmt.Sprintf(
			"[%s] STATUS: %s\nRemarks: %s",
			formattedTime,
			newStatus,
			trimmedNotes,
		)
	}

	updatedNotes := newLogEntry
	if existingNotes != "" {
		updatedNotes = newLogEntry +
			"\n\n------------------------------\n\n" +
			existingNotes
	}

	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			err := s.repo.UpdateStatus(
				ctx,
				tx,
				id,
				newStatus,
				updatedNotes,
			)
			if err != nil {
				audit.Dispatch(
					ctx,
					s.logService,
					s.notifService,
					s.emailService,
					audit.DispatchParams{
						Tx: tx,
						Log: &audit.LogParams{
							Level:    audit.LevelError,
							Category: audit.CategoryAudit,
							Action:   audit.ActionSlipFailed,
							Message: fmt.Sprintf(
								"Failed to update status for admission slip #%s: %s",
								id,
								err.Error(),
							),
							Metadata: &audit.LogMetadata{
								EntityType: constants.SlipEntityType,
								EntityID:   id,
								Error:      err.Error(),
							},
						},
						Notifications: []audit.NotificationParams{
							{
								Title: "Admission Slip Status Update Failed",
								Message: fmt.Sprintf(
									"Failed to update status for admission slip #%s: %s",
									id,
									err.Error(),
								),
								Type: constants.SlipEntityType,
							},
						},
					},
				)
				return err
			}

			emails := []audit.EmailParams{}
			ticketCode := ""

			// Handle ticket and email for Approved status
			if newStatus == "Approved" {
				ticket, err := s.repo.GetTicketBySlipID(ctx, id)
				if err != nil {
					fmt.Printf(
						"[UpdateExcuseSlipStatus] {Get Ticket}: %v\n",
						err,
					)
				}

				if ticket == nil {
					// Generate new ticket
					ticketCode = fmt.Sprintf(
						"SLIP-%d-%s",
						time.Now().Year(),
						strings.ToUpper(uuid.New().String()[:8]),
					)
					newTicket := &AdmissionTicket{
						ID:              uuid.New().String(),
						AdmissionSlipID: id,
						TicketCode:      ticketCode,
					}
					if err := s.repo.CreateTicket(ctx, tx, newTicket); err != nil {
						return err
					}
				} else {
					ticketCode = ticket.TicketCode
				}

				// Always send approval email
				emails = append(emails, audit.EmailParams{
					To:           []string{oldSlip.UserEmail},
					Subject:      "Your Admission Slip Has Been Approved!",
					TemplatePath: "slip.html",
					TemplateData: map[string]interface{}{
						"IsApproved": true,
						"TicketCode": ticketCode,
					},
				})
			} else {
				emails = append(emails, audit.EmailParams{
					To:           []string{oldSlip.UserEmail},
					Subject:      "Admission Slip Status Updated",
					TemplatePath: "slip.html",
					TemplateData: map[string]interface{}{
						"IsApproved": false,
						"StudentName": fmt.Sprintf(
							"%s %s",
							oldSlip.UserFirstName,
							oldSlip.UserLastName,
						),
						"Category": oldSlip.CategoryName,
						"DateOfAbsence": datetime.FormatDate(
							oldSlip.DateOfAbsence,
						),
						"DateNeeded": datetime.FormatDate(
							oldSlip.DateNeeded,
						),
						"Status":     newStatus,
						"AdminNotes": adminNotes,
					},
				})
			}

			// Fetch student UserID for notification
			studentUserID, _ := s.repo.GetUserIDBySlipID(ctx, id)

			notifications := []audit.NotificationParams{
				{
					ReceiverID: structs.StringToNullableString(studentUserID),
					TargetID:   structs.StringToNullableString(id),
					TargetType: structs.StringToNullableString(
						constants.SlipEntityType,
					),
					Title: "Admission Slip Updated",
					Message: fmt.Sprintf(
						"Status for your admission slip has been updated to '%s'",
						newStatus,
					),
					Type: constants.SlipEntityType,
				},
				{
					TargetID: structs.StringToNullableString(id),
					TargetType: structs.StringToNullableString(
						constants.SlipEntityType,
					),
					Title: "Admission Slip Updated Successfully",
					Message: fmt.Sprintf(
						"You have successfully updated the status "+
							"of admission slip %s to '%s'.",
						structs.TruncateString(id, 7),
						newStatus,
					),
					Type: constants.SlipEntityType,
				},
			}

			audit.Dispatch(
				ctx,
				s.logService,
				s.notifService,
				s.emailService,
				audit.DispatchParams{
					Tx: tx,
					Log: &audit.LogParams{
						Level:    audit.LevelInfo,
						Category: audit.CategoryAudit,
						Action:   audit.ActionSlipStatusUpdated,
						Message: fmt.Sprintf(
							"Admission slip #%s status changed to '%s'",
							id,
							newStatus,
						),
						Metadata: &audit.LogMetadata{
							EntityType: constants.SlipEntityType,
							EntityID:   id,
						},
					},
					Notifications: notifications,
					Email:         emails,
				},
			)
			return nil
		},
	)
}

func (s *Service) ClaimTicket(
	ctx context.Context,
	code string,
	counselorID string,
) error {
	ticket, err := s.repo.GetTicketByCode(ctx, code)
	if err != nil {
		return err
	}
	if ticket == nil {
		return fmt.Errorf("ticket not found")
	}
	if ticket.IsVerified {
		return fmt.Errorf("ticket already verified")
	}

	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			if err := s.repo.UpdateTicketVerification(
				ctx,
				tx,
				ticket.ID,
				counselorID,
			); err != nil {
				return err
			}

			// Audit the verification
			audit.Dispatch(
				ctx,
				s.logService,
				s.notifService,
				s.emailService,
				audit.DispatchParams{
					Tx: tx,
					Log: &audit.LogParams{
						Level:    audit.LevelInfo,
						Category: audit.CategoryAudit,
						Action:   audit.ActionSlipStatusUpdated,
						Message:  fmt.Sprintf("Ticket #%s verified", code),
						Metadata: &audit.LogMetadata{
							EntityType: "ticket",
							EntityID:   ticket.ID,
						},
					},
				})
			return nil
		},
	)
}

func (s *Service) GetSlipByTicketCode(
	ctx context.Context,
	code string,
) (*SlipDTO, error) {
	slip, err := s.repo.GetSlipByTicketCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if slip == nil {
		return nil, nil
	}
	return s.mapToDTO(slip), nil
}

func (s *Service) mapToDTO(slip *SlipWithDetailsView) *SlipDTO {
	if slip == nil {
		return nil
	}

	dto := &SlipDTO{
		ID:     slip.ID,
		UserID: slip.UserID,
		IIRID:  slip.IIRID,
		User: users.UserResponse{
			FirstName:  slip.UserFirstName,
			MiddleName: slip.UserMiddleName,
			LastName:   slip.UserLastName,
			Email:      slip.UserEmail,
		},
		StudentNumber: slip.StudentNumber,
		Reason:        slip.Reason,
		DateOfAbsence: slip.DateOfAbsence,
		DateNeeded:    slip.DateNeeded,
		AdminNotes:    slip.AdminNotes,
		Category: SlipCategory{
			ID:   slip.CategoryID,
			Name: slip.CategoryName,
		},
		Status: SlipStatus{
			ID:   slip.StatusID,
			Name: slip.StatusName,
		},
		IsVerified: slip.IsVerified.Bool,
		CreatedAt:  slip.CreatedAt,
		UpdatedAt:  slip.UpdatedAt,
	}

	if slip.TicketCode.Valid {
		dto.Ticket = &TicketDTO{
			TicketCode: slip.TicketCode.String,
			IsVerified: slip.IsVerified.Bool,
		}
		if slip.VerifiedAt.Valid {
			dto.Ticket.VerifiedAt = slip.VerifiedAt.Time
		}
	}

	return dto
}
