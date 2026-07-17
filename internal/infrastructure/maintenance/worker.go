package maintenance

import (
	"context"
	"log"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
)

type MaintenanceWorker struct {
	logService   *logs.Service
	notifService *notifications.Service
}

func NewMaintenanceWorker(
	logService *logs.Service,
	notifService *notifications.Service,
) *MaintenanceWorker {
	return &MaintenanceWorker{
		logService:   logService,
		notifService: notifService,
	}
}

func (w *MaintenanceWorker) Start(ctx context.Context) {
	// Run immediately on start
	w.runCleanup()

	ticker := time.NewTicker(24 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.runCleanup()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *MaintenanceWorker) runCleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	log.Println("[MaintenanceWorker] Starting scheduled cleanup...")

	// Cleanup System Logs (90 days)
	logRows, err := w.logService.DeleteLogsOlderThan(
		ctx,
		90,
		nil,
		[]string{audit.CategorySecurity},
	)
	if err != nil {
		log.Printf("[MaintenanceWorker] Error cleaning up logs: %v", err)
	} else if logRows > 0 {
		log.Printf("[MaintenanceWorker] Purged %d old system logs", logRows)
	}

	// Cleanup Security Logs (1 year)
	secLogRows, err := w.logService.DeleteLogsOlderThan(
		ctx,
		365,
		[]string{audit.CategorySecurity},
		nil,
	)
	if err != nil {
		log.Printf(
			"[MaintenanceWorker] Error cleaning up security logs: %v",
			err,
		)
	} else if secLogRows > 0 {
		log.Printf(
			"[MaintenanceWorker] Purged %d old security logs",
			secLogRows,
		)
	}

	// Cleanup Notifications (30 days)
	notifRows, err := w.notifService.DeleteOldNotifications(ctx, 30)
	if err != nil {
		log.Printf(
			"[MaintenanceWorker] Error cleaning up notifications: %v",
			err,
		)
	} else if notifRows > 0 {
		log.Printf(
			"[MaintenanceWorker] Purged %d old notifications",
			notifRows,
		)
	}

	log.Println(
		"[MaintenanceWorker] Cleanup task completed.")
}
