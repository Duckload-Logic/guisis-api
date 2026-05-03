package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo         *Repository
	notifService audit.Notifier
	userSvc      audit.UserGetter
}

func NewService(
	repo *Repository,
	notifService audit.Notifier,
	userSvc audit.UserGetter,
) *Service {
	return &Service{
		repo:         repo,
		notifService: notifService,
		userSvc:      userSvc,
	}
}

func (s *Service) GetDB() datastore.DB {
	return s.repo.GetDB()
}

func (s *Service) Record(
	ctx context.Context,
	tx datastore.DB,
	entry audit.LogEntry,
) {
	level := entry.Level
	if level == "" {
		level = audit.LevelInfo
	}

	if !entry.TraceID.Valid || entry.TraceID.String == "" {
		_, _, _, _, _, trace := audit.ExtractMeta(ctx)
		if trace != "" {
			entry.TraceID = structs.StringToNullableString(trace)
		}
	}

	var metaStr string
	if entry.Metadata != nil {
		s.sanitizeMetadata(entry.Metadata)
		b, _ := json.Marshal(entry.Metadata)
		metaStr = string(b)
	}

	sysLog := &SystemLog{
		Level:       level,
		Category:    entry.Category,
		Action:      entry.Action,
		Message:     entry.Message,
		UserID:      entry.UserID,
		TargetID:    entry.TargetID,
		TargetType:  entry.TargetType,
		UserEmail:   entry.UserEmail,
		TargetEmail: entry.TargetEmail,
		IPAddress:   entry.IPAddress,
		UserAgent:   entry.UserAgent,
		TraceID:     entry.TraceID,
		Metadata:    structs.StringToNullableString(metaStr),
	}

	if err := s.repo.Record(ctx, tx, sysLog); err != nil {
		fmt.Printf("[Record] {Database Insertion}: %v\n", err)
		return
	}

	if level == audit.LevelError && func(action string) bool {
		excluded := map[string]bool{
			audit.ActionLoginFailed:           true,
			audit.ActionInvalidToken:          true,
			audit.ActionAccessDenied:          true,
			audit.ActionRateLimitExceeded:     true,
			audit.ActionM2MAuthFailed:         true,
			audit.ActionM2MClientVerifyFailed: true,
		}
		return !excluded[action]
	}(entry.Action) {
		s.notifySuperadmins(ctx, entry)
	}
}

func (s *Service) notifySuperadmins(ctx context.Context, entry audit.LogEntry) {
	if s.userSvc == nil || s.notifService == nil {
		return
	}

	adminIDs, err := s.userSvc.GetUserIDsByRole(ctx, 3)
	if err != nil {
		fmt.Printf("[notifySuperadmins] {Fetch Admins}: %v\n", err)
		return
	}

	title := "Critical System Alert"
	if entry.Category == audit.CategorySecurity {
		title = "Security Breach/Alert"
	}

	for _, adminID := range adminIDs {
		if err := s.notifService.Send(ctx, audit.NotificationEntry{
			ReceiverID: structs.StringToNullableString(adminID),
			ActorID:    entry.UserID,
			Title:      title,
			Message: fmt.Sprintf(
				"[%s] %s: %s",
				entry.Action,
				entry.Level,
				entry.Message,
			),
			Type: constants.SystemEntityType,
		}); err != nil {
			fmt.Printf("[notifySuperadmins] {Send Notification}: %v\n", err)
		}
	}
}

func (s *Service) RecordSecurity(
	ctx context.Context,
	tx datastore.DB,
	action, message string,
	userEmail, userID, ipAddress, userAgent structs.NullableString,
) {
	s.Record(ctx, tx, audit.LogEntry{
		Category:  audit.CategorySecurity,
		Action:    action,
		Message:   message,
		UserID:    userID,
		UserEmail: userEmail,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
}

func (s *Service) ListLogs(
	ctx context.Context,
	req audit.ListSystemLogsRequest,
) (*audit.ListSystemLogsDTO, error) {
	req.SetDefaults("created_at")

	results, err := s.repo.List(
		ctx,
		req.GetOffset(), req.PageSize,
		req.Category, req.Action, req.UserEmail,
		req.TargetType, req.TargetEmail,
		req.Search, req.StartDate, req.EndDate, req.OrderBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list system logs: %w", err)
	}

	dtos := s.mapLogsToDTOs(results)

	total, err := s.repo.GetTotalCount(
		ctx,
		req.Category, req.Action, req.UserEmail,
		req.TargetType, req.TargetEmail,
		req.Search, req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to count system logs: %w", err)
	}

	return &audit.ListSystemLogsDTO{
		Logs: dtos,
		Meta: structs.CalculateMetadata(total, req.Page, req.PageSize),
	}, nil
}

func (s *Service) GetStats(
	ctx context.Context,
	startDate, endDate string,
) ([]audit.LogStatsDTO, error) {
	return s.repo.GetStats(ctx, startDate, endDate)
}

func (s *Service) GetActivityStats(
	ctx context.Context,
) ([]audit.LogActivityDTO, error) {
	return s.repo.GetActivityStats(ctx)
}

func (s *Service) mapLogsToDTOs(logs []SystemLog) []audit.SystemLogDTO {
	dtos := make([]audit.SystemLogDTO, 0, len(logs))

	for _, l := range logs {
		dto := audit.SystemLogDTO{
			ID:        l.ID,
			Category:  l.Category,
			Action:    l.Action,
			Message:   l.Message,
			UserID:    l.UserID,
			UserEmail: l.UserEmail,
			IPAddress: l.IPAddress,
			UserAgent: l.UserAgent,
			TraceID:   l.TraceID,
			CreatedAt: l.CreatedAt,
		}

		if l.Metadata.Valid {
			dto.Metadata = json.RawMessage(l.Metadata.String)
		}

		dtos = append(dtos, dto)
	}

	return dtos
}

func (s *Service) DeleteLogsOlderThan(
	ctx context.Context,
	days int,
) (int64, error) {
	return s.repo.DeleteLogsOlderThan(ctx, days)
}

func (s *Service) sanitizeMetadata(meta *audit.LogMetadata) {
	if meta == nil {
		return
	}
	s.sanitizeValue(meta.OldValues)
	s.sanitizeValue(meta.NewValues)
}

func (s *Service) sanitizeValue(val interface{}) {
	m, ok := val.(map[string]interface{})
	if !ok {
		return
	}

	sensitiveKeys := []string{
		"password",
		"secret",
		"token",
		"key",
		"authorization",
		"credential",
	}
	for k := range m {
		lowerK := strings.ToLower(k)
		for _, sk := range sensitiveKeys {
			if strings.Contains(lowerK, sk) {
				delete(m, k)
				break
			}
		}
	}
}
