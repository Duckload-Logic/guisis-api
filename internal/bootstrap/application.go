package bootstrap

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/email"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/maintenance"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

type Application struct {
	Handlers *Handlers
}

func Initialize(db *sqlx.DB, cfg *config.Config) (*Application, error) {
	var fileStorage storage.FileStorage
	var emailer email.Emailer

	if cfg.IsProduction {
		uploadDir := cfg.LocalUploadDIR
		fileStorage = storage.NewDiskStorage(uploadDir)
		emailer = email.NewSMTPMailer(
			cfg.SMTPHost,
			cfg.SMTPPort,
			cfg.SMTPUser,
			cfg.SMTPPass,
		)
	} else {
		uploadDir := cfg.LocalUploadDIR
		fileStorage = storage.NewDiskStorage(uploadDir)

		mailpit, err := email.NewMailPit(cfg.MailPitHost, cfg.MailPitPort)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to initialize MailPit: %w",
				err,
			)
		}
		emailer = mailpit
	}

	repos := getRepositories(db)

	redis, err := datastore.NewRedisClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	services := getServices(repos, fileStorage, cfg, redis, emailer)
	handlers := getHandlers(services, cfg, redis)

	// Start Background Maintenance Worker
	maintenanceWorker := maintenance.NewMaintenanceWorker(
		services.SystemLogService,
		services.NotificationsService,
	)
	maintenanceWorker.Start(context.Background())

	return &Application{
		Handlers: handlers,
	}, nil
}
