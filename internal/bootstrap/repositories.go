package bootstrap

import (
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/internal/features/analytics"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
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
)

type Repositories struct {
	UserRepo               *users.Repository
	StudentRepo            *students.Repository
	NoteRepo               *notes.Repository
	IntegrationStudentRepo integrations.RepositoryInterface
	AppointmentRepo        *appointments.Repository
	SlipRepo               *slips.Repository
	LocationsRepo          *locations.Repository
	AnalyticsRepo          *analytics.Repository
	M2MClientRepo          *m2mclients.Repository
	NotificationRepo       *notifications.Repository
	SystemLogRepo          *logs.Repository
	FileRepo               *files.Repository
}

func getRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		UserRepo:               users.NewRepository(db),
		StudentRepo:            students.NewRepository(db),
		NoteRepo:               notes.NewRepository(db),
		IntegrationStudentRepo: integrations.NewRepository(db),
		AppointmentRepo:        appointments.NewRepository(db),
		SlipRepo:               slips.NewRepository(db),
		LocationsRepo:          locations.NewRepository(db),
		AnalyticsRepo:          analytics.NewRepository(db),
		M2MClientRepo:          m2mclients.NewRepository(db),
		NotificationRepo:       notifications.NewRepository(db),
		SystemLogRepo:          logs.NewRepository(db),
		FileRepo:               files.NewRepository(db),
	}
}
