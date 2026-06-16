package datastore

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetDBConnection(dbUrl string) (*sqlx.DB, error) {
	// Open Database instance
	db, err := sqlx.Connect("mysql", dbUrl)
	if err != nil {
		log.Fatal("Failed to open Database connection:", err)
	}

	// Optimize connection pool to stay within MySQL's limit
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Check ping connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	// Return database instance
	return db, nil
}
