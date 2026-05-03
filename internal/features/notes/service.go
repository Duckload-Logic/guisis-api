package notes

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Service struct {
	repo         *Repository
	logService   audit.Logger
	notifService audit.Notifier
}

func NewService(
	repo *Repository,
	logService audit.Logger,
	notifService audit.Notifier,
) *Service {
	return &Service{
		repo:         repo,
		logService:   logService,
		notifService: notifService,
	}
}

func (s *Service) GetStudentSignificantNotes(
	ctx context.Context,
	iirID string,
) ([]SignificantNoteDTO, error) {
	notes, err := s.repo.GetStudentSignificantNotes(ctx, iirID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get student significant notes: %w",
			err,
		)
	}

	noteDTOs := make([]SignificantNoteDTO, 0, len(notes))
	for _, n := range notes {
		noteDTOs = append(noteDTOs, SignificantNoteDTO{
			ID:            n.ID,
			AppointmentID: n.AppointmentID.String,
			Note:          n.Note,
			Remarks:       n.Remarks,
			CreatedAt:     n.CreatedAt,
			UpdatedAt:     n.UpdatedAt,
		})
	}

	return noteDTOs, nil
}

func (s *Service) CreateSignificantNote(
	ctx context.Context,
	iirID string,
	noteReq SignificantNoteDTO,
) error {
	note := &SignificantNote{
		ID:    uuid.New().String(),
		IIRID: structs.StringToNullableString(iirID),
		AppointmentID: structs.StringToNullableString(
			noteReq.AppointmentID,
		),
		Note:    noteReq.Note,
		Remarks: noteReq.Remarks,
	}

	_, err := s.repo.CreateSignificantNote(ctx, note)
	if err != nil {
		audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelError,
				Category: audit.CategoryAudit,
				Action:   audit.ActionNoteCreateFailed,
				Message: fmt.Sprintf(
					"Failed to create significant note for IIR #%s",
					iirID,
				),
				Metadata: &audit.LogMetadata{
					EntityType: "Note",
					NewValues:  note,
					Error:      err.Error(),
				},
			},
		})
		return fmt.Errorf(
			"failed to create significant note: %w",
			err,
		)
	}

	audit.Dispatch(ctx, s.logService, s.notifService, audit.DispatchParams{
		Log: &audit.LogParams{
			Level:    audit.LevelInfo,
			Category: audit.CategoryAudit,
			Action:   audit.ActionNoteCreated,
			Message: fmt.Sprintf(
				"Significant note created for IIR #%s",
				iirID,
			),
			Metadata: &audit.LogMetadata{
				EntityType: "Note",
				EntityID:   note.ID,
				NewValues:  note,
			},
		},
	})

	return nil
}

func (s *Service) HasNoteForAppointment(
	ctx context.Context,
	appointmentID string,
) (bool, error) {
	return s.repo.HasNoteForAppointment(ctx, appointmentID)
}
