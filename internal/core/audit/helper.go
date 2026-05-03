package audit

import (
	"context"
	"log"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

// DispatchParams holds the parameters for the Dispatch helper.
type DispatchParams struct {
	Log           *LogParams
	Notifications []NotificationParams
	Tx            datastore.DB
}

// LogParams holds the parameters for a log entry.
type LogParams struct {
	Level      string
	Category   string
	Action     string
	Message    string
	TargetID   structs.NullableString
	TargetType structs.NullableString
	Metadata   *LogMetadata
}

// NotificationParams holds the parameters for a notification.
type NotificationParams struct {
	ReceiverID structs.NullableString
	TargetID   structs.NullableString
	TargetType structs.NullableString
	Title      string
	Message    string
	Type       string
}

// Dispatch is a centralized helper to record logs and send notifications.
// It automatically extracts audit metadata from the context.
func Dispatch(
	ctx context.Context,
	logger Logger,
	notifier Notifier,
	params DispatchParams,
) {
	// Extract Meta
	id, ip, ua, email, _, trace := ExtractMeta(ctx)

	// Prepare and Record Log
	if logger != nil && params.Log != nil {
		entry := LogEntry{
			Level:      params.Log.Level,
			Category:   params.Log.Category,
			Action:     params.Log.Action,
			Message:    params.Log.Message,
			UserID:     structs.StringToNullableString(id),
			UserEmail:  structs.StringToNullableString(email),
			IPAddress:  structs.StringToNullableString(ip),
			UserAgent:  structs.StringToNullableString(ua),
			TraceID:    structs.StringToNullableString(trace),
			TargetID:   params.Log.TargetID,
			TargetType: params.Log.TargetType,
			Metadata:   params.Log.Metadata,
		}
		logger.Record(ctx, params.Tx, entry)
	}

	// Prepare and Send Notifications
	if notifier != nil && len(params.Notifications) > 0 {
		for _, n := range params.Notifications {
			receiverID := n.ReceiverID
			if !receiverID.Valid || receiverID.String == "" {
				receiverID = structs.StringToNullableString(id)
			}

			// Safety check: Don't try to send if we still don't have a receiver
			if !receiverID.Valid || receiverID.String == "" {
				log.Printf(
					"[Audit:Dispatch] {Skip Notification}: "+
						"No valid receiver found (Actor ID: '%s')",
					id,
				)
				continue
			}

			notif := NotificationEntry{
				ReceiverID: receiverID,
				ActorID:    structs.StringToNullableString(id),
				TargetID:   n.TargetID,
				TargetType: n.TargetType,
				Title:      n.Title,
				Message:    n.Message,
				Type:       n.Type,
			}
			err := notifier.Send(ctx, notif)
			if err != nil {
				log.Printf(
					`[Audit:Dispatch] {Send Notification}: %v`,
					err,
				)
			}
		}
	}
}
