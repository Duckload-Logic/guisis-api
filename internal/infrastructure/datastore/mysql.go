package datastore

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/olazo-johnalbert/duckload-api/scripts/migrations"
)

const (
	maxOpenConnections = 50
	maxIdleConnections = 10
	connectionLifetime = 5 * time.Minute
)

func GetDBConnection(dbUrl string) (*sqlx.DB, error) {
	// Open Database instance
	db, err := sqlx.Connect("mysql", dbUrl)
	if err != nil {
		log.Fatal("Failed to open Database connection:", err)
	}

	// Optimize connection pool to stay within MySQL's limit
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(connectionLifetime)

	// Check ping connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	// Run migrations automatically
	if err := runDatabaseMigrations(db); err != nil {
		log.Fatalf("Migration runner failed: %v", err)
	}

	// Return database instance
	return db, nil
}

func runDatabaseMigrations(db *sqlx.DB) error {
	sourceDriver, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}

	dbDriver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return err
	}

	migrationInstance, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"mysql",
		dbDriver,
	)
	if err != nil {
		return err
	}

	// Apply all up migrations
	migrationError := migrationInstance.Up()
	if migrationError != nil && migrationError != migrate.ErrNoChange {
		return migrationError
	}

	log.Println("Database migrations applied successfully.")
	return nil
}
