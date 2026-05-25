package bootstrap

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/pdf"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/auth"
	"github.com/olazo-johnalbert/duckload-api/internal/features/files"
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
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/gotenberg"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/ocr"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/storage"
)

type Services struct {
	AuthService               *auth.Service
	UserService               *users.Service
	LocationsService          *locations.Service
	StudentService            *students.Service
	FileService               *files.Service
	NoteService               *notes.Service
	IntegrationStudentService integrations.ServiceInterface
	AppointmentService        *appointments.Service
	SlipService               *slips.Service
	AnalyticsService          *analytics.Service
	M2MClientService          *m2mclients.Service
	NotificationsService      *notifications.Service
	SystemLogService          *logs.Service
	SessionService            *sessions.Service
}

func getServices(
	repos *Repositories,
	fileStorage storage.FileStorage,
	cfg *config.Config,
	redis *datastore.RedisClient,
	emailer audit.Emailer,
) *Services {
	notificationsService := notifications.NewService(repos.NotificationRepo)
	sessionService := sessions.NewService(redis)
	gotenbergClient := gotenberg.NewClient(cfg.GotenbergURL)
	pdfService := pdf.NewService(gotenbergClient)

	ocrClient := ocr.NewClient(cfg.AIBaseUrl, cfg.AiAPIKey)

	fileService := files.NewService(
		repos.FileRepo,
		fileStorage,
		ocrClient,
	)

	userService := users.NewService(repos.UserRepo, sessionService, fileService)
	systemLogService := logs.NewService(
		repos.SystemLogRepo,
		notificationsService,
		userService,
	)
	fileService.SetLogger(systemLogService)

	tokenService := tokens.NewService()

	m2mClientService := m2mclients.NewService(
		repos.M2MClientRepo,
		systemLogService,
		notificationsService,
		emailer,
		userService,
		tokenService,
		sessionService,
		cfg,
	)
	authService := auth.NewService(
		repos.UserRepo,
		redis,
		sessionService,
		emailer,
		systemLogService,
	)

	locationsService := locations.NewService(repos.LocationsRepo)

	studentService := students.NewService(
		repos.StudentRepo,
		locationsService,
		userService,
		fileService,
		systemLogService,
		notificationsService,
		cfg,
		pdfService,
	)
	noteService := notes.NewService(
		repos.NoteRepo,
		systemLogService,
		notificationsService,
		emailer,
	)
	integrationStudentService := integrations.NewService(
		repos.IntegrationStudentRepo,
	)
	appointmentService := appointments.NewService(
		repos.AppointmentRepo,
		notificationsService,
		systemLogService,
		emailer,
		userService,
		noteService,
		studentService,
		cfg,
	)
	slipService := slips.NewService(
		repos.SlipRepo,
		systemLogService,
		notificationsService,
		emailer,
		fileStorage,
		userService,
		studentService,
		fileService,
		ocrClient,
		cfg,
	)
	analyticsService := analytics.NewService(
		repos.AnalyticsRepo,
		redis,
		pdfService,
	)

	return &Services{
		AuthService:               authService,
		UserService:               userService,
		LocationsService:          locationsService,
		StudentService:            studentService,
		NoteService:               noteService,
		IntegrationStudentService: integrationStudentService,
		AppointmentService:        appointmentService,
		FileService:               fileService,
		SlipService:               slipService,
		AnalyticsService:          analyticsService,
		M2MClientService:          m2mClientService,
		NotificationsService:      notificationsService,
		SystemLogService:          systemLogService,
		SessionService:            sessionService,
	}
}
