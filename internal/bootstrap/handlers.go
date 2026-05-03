package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/logs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/m2mclients"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notes"
	"github.com/olazo-johnalbert/duckload-api/internal/features/notifications"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students/integrations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Handlers struct {
	DB                        *sqlx.DB
	AuthHandler               *auth.Handler
	UserHandler               *users.Handler
	LocationsHandler          *locations.Handler
	StudentHandler            *students.Handler
	NoteHandler               *notes.Handler
	IntegrationStudentHandler *integrations.Handler
	AppointmentHandler        *appointments.Handler
	SlipHandler               *slips.Handler
	AnalyticsHandler          *analytics.Handler
	M2MClientHandler          *m2mclients.Handler
	NotificationsHandler      *notifications.Handler
	SystemLogHandler          *logs.Handler
	Redis                     *datastore.RedisClient
}

func getHandlers(
	services *Services,
	cfg *config.Config,
	redis *datastore.RedisClient,
) *Handlers {
	systemLogHandler := logs.NewHandler(services.SystemLogService)
	analyticsHandler := analytics.NewHandler(services.AnalyticsService)
	notificationsHandler := notifications.NewHandler(
		services.NotificationsService,
	)

	return &Handlers{
		AuthHandler: auth.NewHandler(
			services.AuthService,
			cfg,
		),
		UserHandler: users.NewHandler(
			services.UserService,
			services.SessionService,
			services.SystemLogService, // implements audit.Logger
			services.SystemLogService, // implements audit.LogReader
		),
		LocationsHandler: locations.NewHandler(services.LocationsService),
		StudentHandler:   students.NewHandler(services.StudentService),
		NoteHandler: notes.NewHandler(
			services.NoteService,
			services.SystemLogService,
		),
		IntegrationStudentHandler: integrations.NewHandler(
			services.IntegrationStudentService,
		),
		AppointmentHandler: appointments.NewHandler(
			services.AppointmentService, cfg,
		),
		SlipHandler:          slips.NewHandler(services.SlipService),
		AnalyticsHandler:     analyticsHandler,
		M2MClientHandler:     m2mclients.NewHandler(services.M2MClientService),
		NotificationsHandler: notificationsHandler,
		SystemLogHandler:     systemLogHandler,
		Redis:                redis,
	}
}
