package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// SetupTestDB initializes a clean test database using local_schema.sql
// and seeds initial lookup data.
func SetupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DB_URL")
	if dsn == "" {
		dsn = loadEnvDSN()
		if dsn == "" {
			dsn = "root:root@tcp(127.0.0.1:3306)/test_db?parseTime=true"
		}
	}

	// Ensure multiStatements=true is set
	if !strings.Contains(dsn, "multiStatements=true") {
		if strings.Contains(dsn, "?") {
			dsn += "&multiStatements=true"
		} else {
			dsn += "?multiStatements=true"
		}
	}

	// Extract database name from DSN
	parts := strings.Split(dsn, "/")
	if len(parts) < 2 {
		t.Fatalf("invalid TEST_DB_URL DSN format: %s", dsn)
	}
	dbAndParams := parts[len(parts)-1]
	dbName := strings.Split(dbAndParams, "?")[0]

	// Sanitize test name for use as MySQL database name
	sanitizedTestName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, t.Name())
	uniqueDBName := fmt.Sprintf("%s_%s", dbName, sanitizedTestName)
	if len(uniqueDBName) > 60 {
		uniqueDBName = uniqueDBName[:60]
	}

	// Reconstruct admin DSN without dbName
	adminParts := make([]string, len(parts)-1)
	copy(adminParts, parts[:len(parts)-1])
	adminDSN := strings.Join(adminParts, "/") + "/?multiStatements=true"

	// Connect to root to drop & recreate
	adminDB, err := sqlx.Connect("mysql", adminDSN)
	if err != nil {
		t.Fatalf("failed to connect to root database: %v", err)
	}
	defer adminDB.Close()

	_, err = adminDB.Exec(
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", uniqueDBName),
	)
	if err != nil {
		t.Fatalf("failed to drop test database: %v", err)
	}

	_, err = adminDB.Exec(fmt.Sprintf("CREATE DATABASE %s", uniqueDBName))
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	// Register cleanup to drop the database after test finishes
	t.Cleanup(func() {
		cleanupDB, err := sqlx.Connect("mysql", adminDSN)
		if err == nil {
			_, _ = cleanupDB.Exec(
				fmt.Sprintf("DROP DATABASE IF EXISTS %s", uniqueDBName),
			)
			cleanupDB.Close()
		}
	})

	// Reconstruct DSN with unique DB name
	dbParams := ""
	if strings.Contains(dbAndParams, "?") {
		dbParams = "?" + strings.Split(dbAndParams, "?")[1]
	}
	uniqueDSN := strings.Join(parts[:len(parts)-1], "/") +
		"/" + uniqueDBName + dbParams

	// Connect to the clean database
	db, err := sqlx.Connect("mysql", uniqueDSN)
	if err != nil {
		t.Fatalf("failed to connect to clean test database: %v", err)
	}

	// Find project root to load files
	root, err := findProjectRoot()
	if err != nil {
		db.Close()
		t.Fatalf("failed to find project root: %v", err)
	}

	// Load and run local_schema.sql
	schemaPath := filepath.Join(root, "local_schema.sql")
	schemaContent, err := os.ReadFile(schemaPath)
	if err != nil {
		db.Close()
		t.Fatalf("failed to read schema file: %v", err)
	}

	// Strip DEFINER clauses from the schema to prevent privilege errors
	definerRegex := regexp.MustCompile(
		`(?i)\bDEFINER\s*=\s*(?:'[^']*'|"[^"]*"|` + "`[^`]*`" + `|[^\s]+)@` +
			`(?:'[^']*'|"[^"]*"|` + "`[^`]*`" + `|[^\s]+)`,
	)
	cleanSchema := definerRegex.ReplaceAll(schemaContent, []byte(""))

	_, err = db.Exec(string(cleanSchema))
	if err != nil {
		db.Close()
		t.Fatalf("failed to execute schema SQL: %v", err)
	}

	// Load and run seed SQL
	seedPath := filepath.Join(
		root,
		"guisis-api",
		"scripts",
		"seeds",
		"000001_seed_init.up.sql",
	)
	seedContent, err := os.ReadFile(seedPath)
	if err != nil {
		db.Close()
		t.Fatalf("failed to read seed file: %v", err)
	}

	_, err = db.Exec(string(seedContent))
	if err != nil {
		db.Close()
		t.Fatalf("failed to execute seed SQL: %v", err)
	}

	// Seed tables that were created/seeded in later migration files
	extraSeeds := `
		INSERT IGNORE INTO student_statuses (id, status_name) VALUES
		(1, 'Active'),
		(2, 'Graduated'),
		(3, 'On Leave'),
		(4, 'Archived'),
		(5, 'Withdrawn');

		INSERT IGNORE INTO educational_attainments (id, name) VALUES
		(1, 'Doctorate Degree'),
		(2, 'Master\'s Degree'),
		(3, 'College Undergraduate'),
		(4, 'College Graduate'),
		(5, 'Vocational'),
		(6, 'High School Undergraduate'),
		(7, 'High School Graduate'),
		(8, 'Elementary Undergraduate'),
		(9, 'Elementary Graduate'),
		(10, 'Not Indicated');

		INSERT IGNORE INTO academic_settings
			(id, current_year_start, current_year_end, current_term)
		VALUES (1, 2025, 2026, 1);
	`
	_, err = db.Exec(extraSeeds)
	if err != nil {
		db.Close()
		t.Fatalf("failed to execute extra seed SQL: %v", err)
	}

	return db
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for i := 0; i < 10; i++ {
		path := filepath.Join(dir, "local_schema.sql")
		if _, err := os.Stat(path); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("project root (local_schema.sql) not found")
}

func loadEnvDSN() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for i := 0; i < 4; i++ {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			content, err := os.ReadFile(envPath)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				var user, pass, host, port string
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "#") {
						continue
					}
					parts := strings.SplitN(line, "=", 2)
					if len(parts) != 2 {
						continue
					}
					key := strings.TrimSpace(parts[0])
					val := strings.TrimSpace(parts[1])
					switch key {
					case "DB_USER":
						user = val
					case "DB_PASS":
						pass = val
					case "DB_HOST":
						host = val
					case "DB_PORT":
						port = val
					}
				}
				_ = user // Keep compiler happy
				if pass != "" {
					if host == "" {
						host = "127.0.0.1"
					}
					if port == "" {
						port = "3306"
					}
					return fmt.Sprintf(
						"root:%s@tcp(%s:%s)/test_db?parseTime=true",
						pass,
						host,
						port,
					)
				}
			}
		}
		dir = filepath.Dir(dir)
	}
	return ""
}
