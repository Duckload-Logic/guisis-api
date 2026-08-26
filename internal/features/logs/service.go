package logs

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

const (
	logQueueBufferSize = 1000
)

type Service struct {
	repo         *Repository
	notifService audit.Notifier
	userSvc      audit.UserGetter
	logChan      chan *SystemLog
}

func NewService(
	repo *Repository,
	notifService audit.Notifier,
	userSvc audit.UserGetter,
) *Service {
	s := &Service{
		repo:         repo,
		notifService: notifService,
		userSvc:      userSvc,
		logChan:      make(chan *SystemLog, logQueueBufferSize),
	}
	s.startWorker()
	return s
}

func (s *Service) startWorker() {
	go func() {
		for logEntry := range s.logChan {
			ctx := context.Background()
			if err := s.repo.Record(ctx, nil, logEntry); err != nil {
				fmt.Printf("[RecordWorker] {Async Write}: %v\n", err)
			}
		}
	}()
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
		status := "Success"
		if level == audit.LevelError ||
			level == audit.LevelCritical ||
			entry.Metadata.Error != "" ||
			strings.Contains(strings.ToUpper(entry.Action), "FAIL") ||
			strings.Contains(strings.ToUpper(entry.Action), "DENI") ||
			strings.Contains(strings.ToUpper(entry.Action), "INVALID") {
			status = "Failed"
		}

		secureMeta := map[string]interface{}{
			"status": status,
		}
		if entry.Metadata.EntityType != "" {
			secureMeta["entityType"] = entry.Metadata.EntityType
		}
		if entry.Metadata.EntityID != "" {
			secureMeta["entityId"] = entry.Metadata.EntityID
		}
		if entry.Metadata.Error != "" {
			secureMeta["error"] = entry.Metadata.Error
		}

		b, _ := json.Marshal(secureMeta)
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

	if tx != nil {
		if err := s.repo.Record(ctx, tx, sysLog); err != nil {
			fmt.Printf("[Record] {Database Insertion}: %v\n", err)
			return
		}
	} else {
		select {
		case s.logChan <- sysLog:
		default:
			fmt.Printf(
				"[Record] {Queue Insertion}: queue full, dropping log\n",
			)
		}
	}

	// Only notify superadmins for specific critical actions or system errors.
	shouldNotify := (level == audit.LevelError &&
		entry.Action == audit.ActionRateLimitExceeded) ||
		entry.Action == audit.ActionM2MClientCreated ||
		level == audit.LevelCritical

	if shouldNotify {
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
	action, message, userEmail, ipAddress, userAgent string,
) {
	s.Record(ctx, nil, audit.LogEntry{
		Category:  audit.CategorySecurity,
		Action:    action,
		Message:   message,
		UserEmail: structs.StringToNullableString(userEmail),
		IPAddress: structs.StringToNullableString(ipAddress),
		UserAgent: structs.StringToNullableString(userAgent),
	})
}

// RecordEntry records an audit.LogEntry directly. It satisfies the
// middleware.SecurityLogger interface without requiring a transaction.
func (s *Service) RecordEntry(
	ctx context.Context,
	entry audit.LogEntry,
) {
	s.Record(ctx, nil, entry)
}

func (s *Service) ListLogs(
	ctx context.Context,
	req audit.ListSystemLogsRequest,
) (*audit.ListSystemLogsDTO, error) {
	req.SetDefaults("created_at")

	results, err := s.repo.List(
		ctx,
		req.GetOffset(), req.PageSize,
		req.Level, req.Category, req.Action, req.UserEmail,
		req.TargetType, req.TargetEmail,
		req.Search, req.StartDate, req.EndDate,
		req.SortBy, req.SortOrder, // Now passing the two distinct sorting fields
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list system logs: %w", err)
	}

	dtos := s.mapLogsToDTOs(results)

	total, err := s.repo.GetTotalCount(
		ctx,
		req.Level, req.Category, req.Action, req.UserEmail,
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
			ID:          l.ID,
			Level:       l.Level,
			Category:    l.Category,
			Action:      l.Action,
			Message:     l.Message,
			UserID:      l.UserID,
			UserEmail:   l.UserEmail,
			TargetID:    l.TargetID,
			TargetEmail: l.TargetEmail,
			IPAddress:   l.IPAddress,
			UserAgent:   l.UserAgent,
			TraceID:     l.TraceID,
			CreatedAt:   l.CreatedAt,
		}

		if l.Metadata.Valid {
			dto.Metadata = json.RawMessage(l.Metadata.String)
		}

		dtos = append(dtos, dto)
	}

	return dtos
}

func (s *Service) GetLogByID(
	ctx context.Context,
	id int64,
) (*audit.SystemLogDTO, error) {
	result, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get log: %w", err)
	}

	dto := s.mapLogsToDTOs([]SystemLog{*result})
	if len(dto) == 0 {
		return nil, fmt.Errorf("failed to map log to dto")
	}

	return &dto[0], nil
}

func (s *Service) GetTraceTracks(
	ctx context.Context,
	traceID string,
) ([]audit.SystemLogDTO, error) {
	results, err := s.repo.GetByTraceID(ctx, traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trace tracks: %w", err)
	}

	return s.mapLogsToDTOs(results), nil
}

func (s *Service) DeleteLogsOlderThan(
	ctx context.Context,
	days int,
	includeCategories []string,
	excludeCategories []string,
) (int64, error) {
	return s.repo.DeleteLogsOlderThan(ctx, days, includeCategories, excludeCategories)
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

func (s *Service) ExportLogsCSV(
	ctx context.Context,
	req audit.ListSystemLogsRequest,
) ([]byte, error) {
	req.SetDefaults("created_at")

	logs, err := s.repo.ListAll(
		ctx,
		req.Level, req.Category, req.Action, req.UserEmail,
		req.TargetType, req.TargetEmail,
		req.Search, req.StartDate, req.EndDate,
		req.SortBy, req.SortOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logs for export: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{
		"Timestamp", "Level", "Category", "Action", "Message", "User Email",
		"Target Email", "Target Type", "IP Address", "Trace ID",
	}); err != nil {
		return nil, fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, log := range logs {
		if err := writer.Write([]string{
			log.CreatedAt.Format("2006-01-02 15:04:05"), log.Level, log.Category,
			log.Action, log.Message, log.UserEmail.String, log.TargetEmail.String,
			log.TargetType.String, log.IPAddress.String, log.TraceID.String,
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
