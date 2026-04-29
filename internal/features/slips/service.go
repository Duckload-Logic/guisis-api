package slips

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/files"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

const MaxFileSize = 5 * 1024 * 1024 // 5MB limit
type Service struct {
	repo           *Repository
	logService     audit.Logger
	notifService   audit.Notifier
	fileStorage    storage.FileStorage
	userService    *users.Service
	studentService *students.Service
	filesService   *files.Service
}

func NewService(
	repo *Repository,
	logService audit.Logger,
	notifService audit.Notifier,
	fileStorage storage.FileStorage,
	userService *users.Service,
	studentService *students.Service,
	filesService *files.Service,
) *Service {
	return &Service{
		repo:           repo,
		logService:     logService,
		notifService:   notifService,
		fileStorage:    fileStorage,
		userService:    userService,
		studentService: studentService,
		filesService:   filesService,
	}
}

func (s *Service) GetSlipStatuses(ctx context.Context) ([]SlipStatus, error) {
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
	slip, err := s.repo.GetSlipByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	if slip == nil {
		return nil, fmt.Errorf("slip not found")
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
			ID:       slip.StatusID,
			Name:     slip.StatusName,
			ColorKey: slip.StatusColorKey,
		},
		CreatedAt: slip.CreatedAt,
		UpdatedAt: slip.UpdatedAt,
	}

	corMap, _ := s.studentService.GetLatestCORsByUserIDs(ctx, []string{slip.UserID})
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
	for s := range slips {
		slipDTOs = append(slipDTOs, SlipDTO{
			ID:     slips[s].ID,
			UserID: slips[s].UserID,
			IIRID:  slips[s].IIRID,
			User: users.UserResponse{
				FirstName:  slips[s].UserFirstName,
				MiddleName: slips[s].UserMiddleName,
				LastName:   slips[s].UserLastName,
				Email:      slips[s].UserEmail,
			},
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes:    slips[s].AdminNotes,
			Category: SlipCategory{
				ID:   slips[s].CategoryID,
				Name: slips[s].CategoryName,
			},
			Status: SlipStatus{
				ID:       slips[s].StatusID,
				Name:     slips[s].StatusName,
				ColorKey: slips[s].StatusColorKey,
			},
			StudentCORURL: corMap[slips[s].UserID],
			CreatedAt:     slips[s].CreatedAt,
			UpdatedAt:     slips[s].UpdatedAt,
		})
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
	req.SetDefaults("created_at")

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
	for s := range slips {
		slipDTOs = append(slipDTOs, SlipDTO{
			ID:     slips[s].ID,
			UserID: slips[s].UserID,
			IIRID:  slips[s].IIRID,
			User: users.UserResponse{
				FirstName:  slips[s].UserFirstName,
				MiddleName: slips[s].UserMiddleName,
				LastName:   slips[s].UserLastName,
				Email:      slips[s].UserEmail,
			},
			StudentNumber: slips[s].StudentNumber,
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes:    slips[s].AdminNotes,
			Category: SlipCategory{
				ID:   slips[s].CategoryID,
				Name: slips[s].CategoryName,
			},
			Status: SlipStatus{
				ID:       slips[s].StatusID,
				Name:     slips[s].StatusName,
				ColorKey: slips[s].StatusColorKey,
			},
			StudentCORURL: corMap[slips[s].UserID],
			CreatedAt:     slips[s].CreatedAt,
			UpdatedAt:     slips[s].UpdatedAt,
		})
	}

	total, err := s.repo.GetTotalSlipsCount(ctx, &req)
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
	req.SetDefaults("created_at")

	slips, err := s.repo.GetByIIRID(ctx, iirID, &req)
	if err != nil {
		return nil, err
	}

	var slipDTOs []SlipDTO
	for s := range slips {
		slipDTOs = append(slipDTOs, SlipDTO{
			ID:    slips[s].ID,
			IIRID: slips[s].IIRID,
			User: users.UserResponse{
				ID:         "",
				FirstName:  slips[s].UserFirstName,
				MiddleName: slips[s].UserMiddleName,
				LastName:   slips[s].UserLastName,
				Email:      slips[s].UserEmail,
			},
			StudentNumber: slips[s].StudentNumber,
			Reason:        slips[s].Reason,
			DateOfAbsence: slips[s].DateOfAbsence,
			DateNeeded:    slips[s].DateNeeded,
			AdminNotes:    slips[s].AdminNotes,
			Category: SlipCategory{
				ID:   slips[s].CategoryID,
				Name: slips[s].CategoryName,
			},
			Status: SlipStatus{
				ID:       slips[s].StatusID,
				Name:     slips[s].StatusName,
				ColorKey: slips[s].StatusColorKey,
			},
			CreatedAt: slips[s].CreatedAt,
			UpdatedAt: slips[s].UpdatedAt,
		})
	}

	total, err := s.repo.GetTotalSlipsCount(ctx, &req)
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
		// Keep FileURL as the URL path (e.g., /slips/{hash}/{filename})
		// Don't convert it to filesystem path - the frontend needs the URL path
		attachmentDTOs = append(attachmentDTOs, AttachmentDTO{
			ID:       attachments[a].FileID,
			FileName: attachments[a].FileName,
			FileURL:  attachments[a].FileURL,
		})
	}

	return attachmentDTOs, nil
}

func (s *Service) GetAttachmentFile(
	ctx context.Context,
	attachmentID string,
) (*SlipAttachment, error) {
	attachment, err := s.repo.GetAttachmentByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}

	if attachment == nil {
		return nil, fmt.Errorf("attachment not found")
	}

	return attachment, nil
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
				"invalid content type for '%s': expected PDF or Image, got %s",
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

// SubmitExcuseSlip creates a new slip with attachments.
func (s *Service) SubmitExcuseSlip(
	ctx context.Context,
	iirID string,
	req CreateSlipRequest,
	files []*multipart.FileHeader,
) (*Slip, error) {
	// Graduated Student Protocol: Lock records for Graduated or Archived students
	isLocked, err := s.studentService.IsStudentLocked(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to check student status: %w", err)
	}
	if isLocked {
		return nil, fmt.Errorf(
			"cannot submit slip: student record is locked (Graduated/Archived)",
		)
	}

	// Validate all files
	if err := s.validateFiles(files); err != nil {
		return nil, err
	}

	dateOfAbsence := strings.Split(req.DateOfAbsence, "T")[0]
	parsedDate, err := time.Parse("2006-01-02", dateOfAbsence)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid date format: YYYY-MM-DD",
		)
	}

	if parsedDate.After(time.Now()) {
		return nil, fmt.Errorf(
			"absence date cannot be in future",
		)
	}

	// Unified File Implementation: Use files features
	uploadedFiles, err := s.filesService.UploadFiles(ctx, files, "slips")
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

	err = s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			_, err := s.repo.CreateSlip(ctx, tx, slip)
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
		audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
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

	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategoryAudit,
			Action:   audit.ActionSlipCreated,
			Message:  fmt.Sprintf("Excuse slip #%s created", slip.ID),
			Metadata: &audit.LogMetadata{
				EntityType: constants.SlipEntityType,
				EntityID:   slip.ID,
				NewValues:  slip,
			},
		},
		Notifications: notifications,
	})

	return slip, nil
}

func (s *Service) UpdateExcuseSlip(
	ctx context.Context,
	iirID string,
	slipID string,
	req CreateSlipRequest,
	files []*multipart.FileHeader,
) (*Slip, error) {
	// Graduated Student Protocol: Lock records for Graduated or
	// Archived students
	isLocked, err := s.studentService.IsStudentLocked(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf("failed to check student status: %w", err)
	}
	if isLocked {
		return nil, fmt.Errorf(
			"cannot update slip: student record is locked (Graduated/Archived)",
		)
	}

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

	// Validate all files
	if err := s.validateFiles(files); err != nil {
		return nil, err
	}

	// Delete old attachments from both slip records and files table
	oldAttachments, err := s.repo.GetSlipAttachments(ctx, slipID)
	if err == nil {
		for _, att := range oldAttachments {
			_ = s.filesService.DeleteFile(ctx, att.FileID)
		}
	}

	// Upload new files using centralized service
	uploadedFiles, err := s.filesService.UploadFiles(ctx, files, "slips")
	if err != nil {
		return nil, fmt.Errorf("failed to upload files: %w", err)
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
				if err := s.repo.SaveSlipAttachment(ctx, tx, attachment); err != nil {
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
	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategoryAudit,
			Action:   audit.ActionSlipUpdated,
			Message:  fmt.Sprintf("Excuse slip #%s updated", slipID),
			Metadata: &audit.LogMetadata{
				EntityType: constants.SlipEntityType,
				EntityID:   slipID,
				NewValues:  updatedSlip,
			},
		},
		Notifications: []audit.NotificationParams{
			// Notification para sa Student
			{
				ReceiverID: structs.StringToNullableString(studentUserID),
				Title:      "Slip Updated",
				Message:    fmt.Sprintf("Your slip #%s has been updated", slipID),
				Type:       constants.SlipEntityType,
			},
			// Notification para sa Counselor (if needed)
			{
				ReceiverID: structs.StringToNullableString(audit.ExtractUserID(ctx)),
				Title:      "New Slip Update",
				Message:    fmt.Sprintf("Student updated slip #%s", slipID),
				Type:       constants.SlipEntityType,
			},
		},
	})

	return updatedSlip, nil
}

// DownloadAttachment streams the attachment from Azure Blob Storage.
func (s *Service) DownloadAttachment(
	ctx context.Context,
	attachmentID string,
	writer io.Writer,
) (*SlipAttachment, error) {
	attachment, err := s.repo.GetAttachmentByID(ctx, attachmentID)
	if err != nil {
		return nil, err
	}
	if attachment == nil {
		return nil, fmt.Errorf("attachment not found")
	}

	// Convert URL path "/uploads/slips/hash/file" to blob path "slips/hash/file"
	blobPath := strings.TrimPrefix(attachment.FileURL, "/uploads/")
	blobPath = strings.TrimPrefix(blobPath, "/")

	// Security: Path Traversal Protection (Jail Check)
	if strings.Contains(blobPath, "..") ||
		!(strings.HasPrefix(blobPath, "slips/") ||
			strings.HasPrefix(blobPath, "cors/")) {
		return nil, fmt.Errorf("security: invalid file path detected")
	}

	if err := s.fileStorage.Download(ctx, blobPath, writer); err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return attachment, nil
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
			"invalid status: must be 'Pending', 'Approved', 'Rejected', or 'For Revision'",
		)
	}

	// Fetch old state for audit trail
	oldSlip, _ := s.repo.GetSlipByID(ctx, id)

	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			err := s.repo.UpdateStatus(ctx, tx, id, newStatus, adminNotes)
			if err != nil {
				audit.Dispatch(
					ctx,
					s.logService,
					s.notifService,
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
						"You have successfully updated the status of admission slip %s to '%s'.",
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
							OldValues:  oldSlip,
							NewValues: map[string]interface{}{
								"status":     newStatus,
								"adminNotes": adminNotes,
							},
						},
					},
					Notifications: notifications,
				},
			)
			return nil
		},
	)
}

func (s *Service) GetUserIDBySlipID(
	ctx context.Context,
	id string,
) (string, error) {
	return s.repo.GetUserIDBySlipID(ctx, id)
}
