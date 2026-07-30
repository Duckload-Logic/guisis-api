package appointments

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/datetime"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/ai/classifier"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"

	"github.com/google/uuid"
)

const statusPending = 1

type Service struct {
	repo           *Repository
	notifService   audit.Notifier
	logService     audit.Logger
	emailService   audit.Emailer
	userService    *users.Service
	noteService    *notes.Service
	studentService *students.Service
	classifier     *classifier.ClassifierClient
}

func NewService(
	repo *Repository,
	notifService audit.Notifier,
	logService audit.Logger,
	emailService audit.Emailer,
	userService *users.Service,
	noteService *notes.Service,
	studentService *students.Service,
	cfg *config.Config,
) *Service {
	return &Service{
		repo:           repo,
		notifService:   notifService,
		logService:     logService,
		emailService:   emailService,
		userService:    userService,
		noteService:    noteService,
		studentService: studentService,
		classifier: classifier.NewClient(
			http.DefaultClient,
			cfg.AIBaseUrl,
			cfg.AIAPIKey,
		),
	}
}

func (s *Service) GetConcernCategories(
	ctx context.Context,
) ([]AppointmentCategory, error) {
	return s.repo.GetCategories(ctx)
}

func (s *Service) CreateAppointment(
	ctx context.Context,
	iirID string,
	req AppointmentDTO,
	cfg *config.Config,
) (*AppointmentDTO, error) {
	appt := &Appointment{
		ID:         uuid.New().String(),
		IIRID:      iirID,
		Reason:     req.Reason,
		WhenDate:   datetime.ExtractDateOnly(req.WhenDate),
		TimeSlotID: req.TimeSlot.ID,
		CategoryID: req.AppointmentCategory.ID,
		StatusID:   statusPending,
	}

	if req.PreferredDate1 != "" {
		appt.PrefDate1 = structs.StringToNullableString(
			datetime.ExtractDateOnly(req.PreferredDate1),
		)
	}
	if req.PreferredTimeSlot1 != nil && req.PreferredTimeSlot1.ID != 0 {
		appt.PrefTimeSlotID1 = structs.Int64ToNullableInt64(
			int64(req.PreferredTimeSlot1.ID),
		)
	}
	if req.PreferredDate2 != "" {
		appt.PrefDate2 = structs.StringToNullableString(
			datetime.ExtractDateOnly(req.PreferredDate2),
		)
	}
	if req.PreferredTimeSlot2 != nil && req.PreferredTimeSlot2.ID != 0 {
		appt.PrefTimeSlotID2 = structs.Int64ToNullableInt64(
			int64(req.PreferredTimeSlot2.ID),
		)
	}
	if req.PreferredDate3 != "" {
		appt.PrefDate3 = structs.StringToNullableString(
			datetime.ExtractDateOnly(req.PreferredDate3),
		)
	}
	if req.PreferredTimeSlot3 != nil && req.PreferredTimeSlot3.ID != 0 {
		appt.PrefTimeSlotID3 = structs.Int64ToNullableInt64(
			int64(req.PreferredTimeSlot3.ID),
		)
	}

	// Default values for urgency level and score
	appt.UrgencyLevel = "MEDIUM"
	appt.UrgencyScore = 0.0

	classification, err := s.classifier.Classify(ctx, appt.Reason.String, cfg)
	if err == nil {
		appt.UrgencyLevel = classification.Level
		appt.UrgencyScore = classification.Confidence
	} else {
		s.logService.Record(ctx, nil, audit.LogEntry{
			Level:    audit.LevelError,
			Category: audit.CategorySystem,
			Action:   "AI_CLASSIFICATION_FAILED",
			Message: fmt.Sprintf(
				"HuggingFace AI classification failed: %v",
				err,
			),
		})
	}

	var createdAppt *AppointmentWithDetailsView
	err = s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			available, err := s.repo.IsSlotAvailableForUpdate(
				ctx, tx, appt.WhenDate, appt.TimeSlotID,
			)
			if err != nil {
				return err
			}
			if !available {
				return fmt.Errorf("selected time slot is no longer available")
			}

			createdAppt, err = s.repo.CreateAppointment(ctx, tx, appt)
			return err
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
					Action:   audit.ActionAppointmentFailed,
					Message: fmt.Sprintf(
						"Failed to create appointment for IIR #%s",
						iirID,
					),
					Metadata: &audit.LogMetadata{
						EntityType: constants.AppointmentEntityType,
					},
				},
			})
		return nil, err
	}

	// Fetch personalized notification targets
	userID := audit.ExtractUserID(ctx)
	student, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	studentName := student.FullName()

	counselors, err := s.userService.GetUsersByRole(
		ctx,
		int(constants.AdminRoleID),
	)
	if err != nil {
		return nil, err
	}

	counselorEmails := make([]string, 0, len(counselors))
	for _, counselor := range counselors {
		counselorEmails = append(counselorEmails, counselor.Email)
	}

	notifications := []audit.NotificationParams{
		{
			ReceiverID: structs.StringToNullableString(userID),
			TargetID:   structs.StringToNullableString(appt.ID),
			TargetType: structs.StringToNullableString(
				constants.AppointmentEntityType,
			),
			Title: "Appointment Created Successfully",
			Message: "Your appointment has been created and " +
				"is pending approval.",
			Type: constants.AppointmentEntityType,
		},
	}

	for _, counselor := range counselors {
		notifications = append(notifications, audit.NotificationParams{
			ReceiverID: structs.StringToNullableString(counselor.ID),
			TargetID:   structs.StringToNullableString(appt.ID),
			TargetType: structs.StringToNullableString(
				constants.AppointmentEntityType,
			),
			Title: "New Appointment Request",
			Message: fmt.Sprintf(
				"New appointment request received from %s.",
				studentName,
			),
			Type: constants.AppointmentEntityType,
		})
	}

	newApptDTO := s.mapToDTO(createdAppt)

	audit.Dispatch(
		ctx,
		s.logService,
		s.notifService,
		s.emailService,
		audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelInfo,
				Category: audit.CategoryAudit,
				Action:   audit.ActionAppointmentCreated,
				Message:  fmt.Sprintf("Appointment #%s created", appt.ID),
				Metadata: &audit.LogMetadata{
					EntityType: constants.AppointmentEntityType,
					EntityID:   appt.ID,
				},
			},
			Notifications: notifications,
			Email: []audit.EmailParams{
				{
					To:           counselorEmails,
					Subject:      "New Appointment Request",
					TemplatePath: "request.html",
					TemplateData: map[string]any{
						"EntityType":   constants.AppointmentEntityType,
						"StudentName":  studentName,
						"UrgencyLevel": newApptDTO.UrgencyLevel,
						"Category":     newApptDTO.AppointmentCategory.Name,
						"Reason":       newApptDTO.Reason.String,
						"Time": datetime.FormatTime(
							newApptDTO.TimeSlot.Time,
						),
						"Date": datetime.FormatDate(
							newApptDTO.WhenDate,
						),
						"Status": newApptDTO.Status.Name,
						"RequestURL": fmt.Sprintf(
							"%s/admin/appointments/%s",
							cfg.BaseURL,
							appt.ID,
						),
					},
				},
			},
		})

	return newApptDTO, nil
}

func (s *Service) GetAppointmentByID(
	ctx context.Context,
	id string,
) (*AppointmentDTO, error) {
	appt, err := s.repo.GetAppointment(ctx, s.repo.GetDB(), id)
	if err != nil {
		return nil, err
	}
	if appt == nil {
		return nil, nil
	}

	dto := s.mapToDTO(appt)

	hasNote, _ := s.noteService.HasNoteForAppointment(ctx, id)
	dto.HasSignificantNote = hasNote

	// Fetch student COR URL
	userID, _ := s.repo.GetUserIDByAppointmentID(ctx, id)
	if userID != "" {
		corMap, _ := s.studentService.GetLatestCORsByUserIDs(
			ctx,
			[]string{userID},
		)
		dto.StudentCORURL = corMap[userID]
	}

	return dto, nil
}

func (s *Service) GetDailyStatusCount(
	ctx context.Context,
	startDate string,
) ([]DailyStatusCount, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, startDate)
	if err != nil {
		return nil, err
	}

	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	endOfMonth := time.Date(
		t.Year(),
		t.Month()+2,
		0,
		23,
		59,
		59,
		0,
		t.Location(),
	)

	startStr := startOfMonth.Format(layout)
	endStr := endOfMonth.Format(layout)

	return s.repo.GetDailyStatusCount(ctx, startStr, endStr)
}

func (s *Service) ListAppointments(
	ctx context.Context,
	req ListAppointmentsRequest,
) (*ListAppointmentsResponse, error) {
	req.SetDefaults("created_at")

	statusIDs := []string{}
	if req.StatusID != "" {
		statusIDs = strings.Split(req.StatusID, ",")
	}

	appts, err := s.repo.List(
		ctx,
		req.GetOffset(),
		req.PageSize,
		req.Search,
		req.OrderBy,
		req.SortOrder,
		strings.Join(statusIDs, ","),
		req.StartDate,
		req.EndDate,
		req.CategoryID,
		req.Urgency,
	)
	if err != nil {
		return nil, err
	}

	// Batch fetch COR URLs
	userIDs := make([]string, 0, len(appts))
	for i := range appts {
		if appts[i].UserID != "" {
			userIDs = append(userIDs, appts[i].UserID)
		}
	}
	corMap, _ := s.studentService.GetLatestCORsByUserIDs(
		ctx,
		userIDs,
	)

	dtos := make([]AppointmentDTO, 0, len(appts))
	for i := range appts {
		dto := s.mapToDTO(&appts[i])
		dto.StudentCORURL = corMap[appts[i].UserID]

		hasNote, _ := s.noteService.HasNoteForAppointment(
			ctx,
			appts[i].ID,
		)
		dto.HasSignificantNote = hasNote

		dtos = append(dtos, *dto)
	}

	total, err := s.repo.GetTotalAppointmentsCount(
		ctx,
		req.Search,
		req.StatusID,
		req.StartDate,
		req.EndDate,
		req.CategoryID,
		req.Urgency,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentsResponse{
		Appointments: dtos,
		Meta:         structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) ExportAppointmentsCSV(
	ctx context.Context,
	req ListAppointmentsRequest,
) ([]byte, error) {
	req.SetDefaults("created_at")

	statusIDs := req.StatusID
	appointments, err := s.repo.ListAll(
		ctx,
		req.Search, req.OrderBy, req.SortOrder, statusIDs,
		req.StartDate, req.EndDate, req.CategoryID, req.Urgency,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch appointments for export: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{
		"Appointment Date", "Time Slot", "Student Number", "Student Name",
		"Email", "Category", "Status", "Urgency", "Reason", "Created At",
	}); err != nil {
		return nil, fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, appointment := range appointments {
		if err := writer.Write([]string{
			appointment.WhenDate, appointment.TimeSlotTime, appointment.StudentNumber,
			appointment.FullName(), appointment.UserEmail, appointment.CategoryName,
			appointment.StatusName, appointment.UrgencyLevel, appointment.Reason.String,
			appointment.CreatedAt.Format("2006-01-02 15:04:05"),
		}); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush CSV writer: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *Service) GetAppointmentsByUserID(
	ctx context.Context,
	userID string,
	req ListAppointmentsRequest,
) (*ListAppointmentsResponse, error) {
	req.SetDefaults("created_at")

	appts, err := s.repo.ListByUserID(
		ctx,
		userID,
		req.GetOffset(),
		req.PageSize,
		req.OrderBy,
		req.SortOrder,
		req.StatusID,
		req.StartDate,
		req.EndDate,
		req.CategoryID,
		req.Urgency,
	)
	if err != nil {
		return nil, err
	}

	dtos := make([]AppointmentDTO, 0, len(appts))
	for i := range appts {
		dto := s.mapToDTO(&appts[i])

		hasNote, _ := s.noteService.HasNoteForAppointment(
			ctx,
			appts[i].ID,
		)
		dto.HasSignificantNote = hasNote

		dtos = append(dtos, *dto)
	}

	total, err := s.repo.GetTotalAppointmentsCount(
		ctx,
		"",
		req.StatusID,
		req.StartDate,
		req.EndDate,
		req.CategoryID,
		req.Urgency,
		&userID,
	)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentsResponse{
		Appointments: dtos,
		Meta: structs.CalculateMetadata(
			total,
			req.Page,
			req.PageSize,
		),
	}, nil
}
func (s *Service) GetAppointmentsByIIRID(
	ctx context.Context,
	iirID string,
	req ListAppointmentsRequest,
) (*ListAppointmentsResponse, error) {
	req.SetDefaults("created_at")

	appts, err := s.repo.ListByIIRID(
		ctx,
		iirID,
		req.GetOffset(),
		req.PageSize,
		req.OrderBy,
		req.SortOrder,
		req.StatusID,
		req.StartDate,
		req.EndDate,
		req.CategoryID,
		req.Urgency,
	)
	if err != nil {
		return nil, err
	}

	dtos := make([]AppointmentDTO, 0, len(appts))
	for i := range appts {
		dto := s.mapToDTO(&appts[i])

		hasNote, _ := s.noteService.HasNoteForAppointment(
			ctx,
			appts[i].ID,
		)
		dto.HasSignificantNote = hasNote

		dtos = append(dtos, *dto)
	}

	total, err := s.repo.GetTotalAppointmentsCount(
		ctx,
		"",
		req.StatusID,
		req.StartDate,
		req.EndDate,
		req.CategoryID,
		req.Urgency,
		&iirID,
	)
	if err != nil {
		return nil, err
	}

	return &ListAppointmentsResponse{
		Appointments: dtos,
		Meta:         structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetAppointmentStats(
	ctx context.Context,
	req ListAppointmentsRequest,
	iirID *string,
) ([]StatusCount, error) {
	return s.repo.GetAppointmentStats(
		ctx,
		req.StatusID,
		req.StartDate,
		req.EndDate,
		req.CategoryID,
		req.Urgency,
		iirID,
	)
}

func (s *Service) GetAvailableTimeSlots(
	ctx context.Context,
	date string,
) ([]AvailableTimeSlotView, error) {
	availableSlots, err := s.repo.GetAvailableTimeSlots(ctx, date)
	if err != nil {
		return nil, err
	}

	return availableSlots, nil
}

func (s *Service) GetAppointmentStatuses(
	ctx context.Context,
) ([]AppointmentStatus, error) {
	return s.repo.GetStatuses(ctx)
}

// handles Status updates AND Rescheduling
func (s *Service) UpdateAppointment(
	ctx context.Context,
	id string,
	req AppointmentDTO,
) error {
	// Fetch old state for audit trail
	oldAppt, err := s.repo.GetAppointment(
		ctx,
		s.repo.GetDB(),
		id,
	)
	if err != nil {
		return err
	}
	if oldAppt == nil {
		return fmt.Errorf("appointment not found")
	}

	reqDateOnly := datetime.ExtractDateOnly(req.WhenDate)
	statusChanged := req.Status.ID != oldAppt.StatusID
	scheduleChanged := (reqDateOnly != oldAppt.WhenDate) ||
		(req.TimeSlot.ID != oldAppt.TimeSlotID)

	// Enforce remarks for reschedule
	if scheduleChanged &&
		strings.TrimSpace(req.AdminNotes.String) == "" {
		return fmt.Errorf(
			"counselor remarks/reason for reschedule is required",
		)
	}

	existingNotes := ""
	if oldAppt.AdminNotes.Valid {
		existingNotes = oldAppt.AdminNotes.String
	}

	formattedTime := datetime.FormatDateTime(time.Now())
	newLogEntry := ""

	if scheduleChanged {
		newTimeSlot, err := s.repo.GetTimeSlotByID(ctx, req.TimeSlot.ID)
		if err != nil {
			return err
		}
		newTimeSlotTime := ""
		if newTimeSlot != nil {
			newTimeSlotTime = datetime.FormatTime(newTimeSlot.Time)
		}
		newDateFormatted := datetime.FormatDate(reqDateOnly)
		oldDateFormatted := datetime.FormatDate(oldAppt.WhenDate)
		oldTimeFormatted := datetime.FormatTime(oldAppt.TimeSlotTime)

		newLogEntry = fmt.Sprintf(
			"[%s] STATUS: RESCHEDULED\n"+
				"Remarks: %s\n"+
				"Rescheduled from %s at %s to %s at %s.",
			formattedTime,
			strings.TrimSpace(req.AdminNotes.String),
			oldDateFormatted,
			oldTimeFormatted,
			newDateFormatted,
			newTimeSlotTime,
		)
	} else if statusChanged {
		newStatus, err := s.repo.GetStatusByID(ctx, req.Status.ID)
		if err != nil {
			return err
		}
		newStatusName := ""
		if newStatus != nil {
			newStatusName = newStatus.Name
		}

		newLogEntry = fmt.Sprintf(
			"[%s] STATUS: %s",
			formattedTime,
			strings.ToUpper(newStatusName),
		)
		trimmedNotes := strings.TrimSpace(req.AdminNotes.String)
		if trimmedNotes != "" {
			newLogEntry = fmt.Sprintf(
				"[%s] STATUS: %s\nRemarks: %s",
				formattedTime,
				strings.ToUpper(newStatusName),
				trimmedNotes,
			)
		}
	}

	updatedNotes := existingNotes
	if newLogEntry != "" {
		if existingNotes != "" {
			updatedNotes = newLogEntry +
				"\n\n------------------------------\n\n" +
				existingNotes
		} else {
			updatedNotes = newLogEntry
		}
	}

	appt := Appointment{
		ID:         id,
		StatusID:   req.Status.ID,
		Reason:     req.Reason,
		AdminNotes: structs.StringToNullableString(updatedNotes),
		WhenDate:   reqDateOnly,
		TimeSlotID: req.TimeSlot.ID,
		CategoryID: req.AppointmentCategory.ID,
	}

	err = s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.UpdateAppointment(ctx, tx, appt)
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
					Action:   audit.ActionAppointmentUpdateFailed,
					Message: fmt.Sprintf(
						"Failed to update appointment #%s",
						id,
					),
					Metadata: &audit.LogMetadata{
						EntityType: "appointment",
						EntityID:   id,
						Error:      err.Error(),
					},
				},
			})

		return err
	}

	newAppt, _ := s.repo.GetAppointment(ctx, s.repo.GetDB(), id)

	// Fetch student UserID for notification
	studentUserID, _ := s.repo.GetUserIDByAppointmentID(ctx, id)
	actorUserID, _, _, adminEmail, _, _ := audit.ExtractMeta(ctx)

	var notifications []audit.NotificationParams
	if actorUserID == studentUserID &&
		strings.ToLower(newAppt.StatusName) == "cancelled" {
		counselorIDs, _ := s.userService.GetUserIDsByRole(
			ctx,
			int(constants.AdminRoleID),
		)
		studentName := newAppt.FullName()
		for _, cid := range counselorIDs {
			notifications = append(notifications, audit.NotificationParams{
				ReceiverID: structs.StringToNullableString(cid),
				TargetID:   structs.StringToNullableString(newAppt.ID),
				TargetType: structs.StringToNullableString(
					constants.AppointmentEntityType,
				),
				Title: "Appointment Cancelled by Student",
				Message: fmt.Sprintf(
					"Appointment scheduled on %s at %s has "+
						"been cancelled by %s.",
					datetime.FormatDate(newAppt.WhenDate),
					datetime.FormatTime(newAppt.TimeSlotTime),
					studentName,
				),
				Type: constants.AppointmentEntityType,
			})
		}
		// Confirm to student
		notifications = append(notifications, audit.NotificationParams{
			ReceiverID: structs.StringToNullableString(studentUserID),
			TargetID:   structs.StringToNullableString(newAppt.ID),
			TargetType: structs.StringToNullableString(
				constants.AppointmentEntityType,
			),
			Title: "Appointment Cancelled Successfully",
			Message: fmt.Sprintf(
				"You have cancelled your appointment scheduled "+
					"on %s at %s.",
				datetime.FormatDate(newAppt.WhenDate),
				datetime.FormatTime(newAppt.TimeSlotTime),
			),
			Type: constants.AppointmentEntityType,
		})
	} else {
		notifications = []audit.NotificationParams{
			{
				ReceiverID: structs.StringToNullableString(studentUserID),
				TargetID:   structs.StringToNullableString(newAppt.ID),
				TargetType: structs.StringToNullableString(
					constants.AppointmentEntityType,
				),
				Title: fmt.Sprintf(
					"Appointment Status Updated By %s",
					adminEmail,
				),
				Message: fmt.Sprintf(
					"Appointment scheduled on %s at %s has "+
						"been updated to '%s'",
					datetime.FormatDate(newAppt.WhenDate),
					datetime.FormatTime(newAppt.TimeSlotTime),
					newAppt.StatusName,
				),
				Type: constants.AppointmentEntityType,
			},
			{
				TargetID: structs.StringToNullableString(oldAppt.ID),
				TargetType: structs.StringToNullableString(
					constants.AppointmentEntityType,
				),
				Title: "Appointment Updated Successfully",
				Message: fmt.Sprintf(
					"You have successfully updated the status of "+
						"appointment #%s scheduled on %s at %s to '%s'.",
					structs.TruncateString(oldAppt.ID, 7),
					datetime.FormatDate(newAppt.WhenDate),
					datetime.FormatTime(newAppt.TimeSlotTime),
					newAppt.StatusName,
				),
				Type: constants.AppointmentEntityType,
			},
		}
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
				Action:   audit.ActionAppointmentUpdated,
				Message:  fmt.Sprintf("Appointment #%s updated", id),
				Metadata: &audit.LogMetadata{
					EntityType: constants.AppointmentEntityType,
					EntityID:   id,
				},
			},
			Notifications: notifications,
			Email: []audit.EmailParams{
				{
					To:           []string{newAppt.UserEmail},
					Subject:      "Appointment Status Updated",
					TemplatePath: "appointment.html",
					TemplateData: map[string]interface{}{
						"StudentName": fmt.Sprintf(
							"%s %s",
							newAppt.UserFirstName,
							newAppt.UserLastName,
						),
						"Date":       datetime.FormatDate(newAppt.WhenDate),
						"Time":       datetime.FormatTime(newAppt.TimeSlotTime),
						"Category":   newAppt.CategoryName,
						"Status":     newAppt.StatusName,
						"AdminNotes": newAppt.AdminNotes.String,
					},
				},
			},
		})

	// Add special prompt for counselors if appointment is completed
	// status_id 3 = Completed
	if req.Status.ID == 3 {
		hasNote, _ := s.noteService.HasNoteForAppointment(ctx, id)
		if !hasNote {
			audit.Dispatch(
				ctx,
				s.logService,
				s.notifService,
				s.emailService,
				audit.DispatchParams{
					Notifications: []audit.NotificationParams{
						{
							ReceiverID: structs.StringToNullableString(
								audit.ExtractUserID(ctx),
							),
							TargetID: structs.StringToNullableString(id),
							TargetType: structs.StringToNullableString(
								constants.AppointmentEntityType,
							),
							Title: "Action Required: Significant Note",
							Message: "Appointment completed. Please record " +
								"any significant notes or incidents for this " +
								"student.",
							Type: constants.AppointmentEntityType,
						},
					},
				},
			)
		}
	}

	return nil
}

func (s *Service) GetUserIDByAppointmentID(
	ctx context.Context,
	id string,
) (string, error) {
	return s.repo.GetUserIDByAppointmentID(ctx, id)
}

func (s *Service) mapToDTO(
	appt *AppointmentWithDetailsView,
) *AppointmentDTO {
	if appt == nil {
		return nil
	}

	dto := &AppointmentDTO{
		ID:     appt.ID,
		UserID: appt.UserID,
		User: users.UserResponse{
			ID:             appt.UserID,
			FirstName:      appt.UserFirstName,
			MiddleName:     appt.UserMiddleName,
			LastName:       appt.UserLastName,
			Email:          appt.UserEmail,
			ProfilePicture: appt.UserProfilePicture.String,
		},
		IIRID:         appt.IIRID,
		StudentNumber: appt.StudentNumber,
		Reason:        appt.Reason,
		AdminNotes:    appt.AdminNotes,
		WhenDate:      appt.WhenDate,
		TimeSlot: TimeSlot{
			ID:   appt.TimeSlotID,
			Time: appt.TimeSlotTime,
		},
		AppointmentCategory: AppointmentCategory{
			ID:   appt.CategoryID,
			Name: appt.CategoryName,
		},
		Status: AppointmentStatus{
			ID:   appt.StatusID,
			Name: appt.StatusName,
		},
		UrgencyLevel: appt.UrgencyLevel,
		UrgencyScore: appt.UrgencyScore,
		CreatedAt:    appt.CreatedAt,
		UpdatedAt:    appt.UpdatedAt,
	}

	if appt.PrefDate1.Valid {
		dto.PreferredDate1 = appt.PrefDate1.String
	}
	if appt.PrefTimeSlotID1.Valid {
		dto.PreferredTimeSlot1 = &TimeSlot{
			ID:   int(appt.PrefTimeSlotID1.Int64),
			Time: appt.PrefTimeSlotTime1.String,
		}
	}
	if appt.PrefDate2.Valid {
		dto.PreferredDate2 = appt.PrefDate2.String
	}
	if appt.PrefTimeSlotID2.Valid {
		dto.PreferredTimeSlot2 = &TimeSlot{
			ID:   int(appt.PrefTimeSlotID2.Int64),
			Time: appt.PrefTimeSlotTime2.String,
		}
	}
	if appt.PrefDate3.Valid {
		dto.PreferredDate3 = appt.PrefDate3.String
	}
	if appt.PrefTimeSlotID3.Valid {
		dto.PreferredTimeSlot3 = &TimeSlot{
			ID:   int(appt.PrefTimeSlotID3.Int64),
			Time: appt.PrefTimeSlotTime3.String,
		}
	}

	return dto
}
