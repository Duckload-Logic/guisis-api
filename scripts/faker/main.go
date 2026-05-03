package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/olazo-johnalbert/duckload-api/internal/features/appointments"
	"github.com/olazo-johnalbert/duckload-api/internal/features/locations"
	"github.com/olazo-johnalbert/duckload-api/internal/features/slips"
	"github.com/olazo-johnalbert/duckload-api/internal/features/students"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

func main() {
	// ---------- OPTIONAL CSV INPUTS ----------
	studentsCSV := flag.String("students-csv", "", "Path to students CSV file")
	adminsCSV := flag.String("admins-csv", "", "Path to admins CSV file")
	superadminsCSV := flag.String(
		"superadmins-csv",
		"",
		"Path to superadmins CSV file",
	)
	developersCSV := flag.String(
		"developers-csv",
		"",
		"Path to developers CSV file",
	)
	appointmentsCSV := flag.String(
		"appointments-csv",
		"",
		"Path to appointments CSV file",
	)
	// backfillIIR := flag.Bool("backfill-iir", false, "Backfill missing IIR records for existing students")
	flag.Parse()

	// ---------- CONFIGURATION ----------
	numStudents := 50  // number of students to generate
	numCounselors := 5 // number of counselors (admins)
	numSuperAdmin := 4 // number of super admins
	numWorkers := 100  // number of concurrent student workers
	numDevelopers := 2 // number of developers
	passwordHash := fakePasswordHash()
	_ = godotenv.Load()
	dsn := buildDSNFromEnv()
	fmt.Println(dsn)
	startTime := time.Now()
	// -----------------------------------

	var err error
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}
	defer db.Close()

	usersRepo = users.NewRepository(db)
	studentsRepo = students.NewRepository(db)
	appointmentsRepo = appointments.NewRepository(db)
	slipsRepo = slips.NewRepository(db)
	locationsRepo = locations.NewRepository(db)

	// load all lookup data from reference tables
	loadLookups()

	// seed random generator
	gofakeit.Seed(time.Now().UnixNano())

	// if *backfillIIR {
	// 	runBackfill(passwordHash)
	// 	fmt.Println("Backfill completed successfully.")
	// 	log.Println("Time taken:", time.Since(startTime))
	// 	return
	// }

	// clear existing student data (optional but keeps the run idempotent)
	clearStudentData()

	// Process Super Admins
	var superadminsFromCSV []users.User
	if *superadminsCSV != "" {
		superadminsFromCSV, err = parseUsersFromCSV(*superadminsCSV)
		if err != nil {
			log.Fatal(err)
		}
		numSuperAdmin = len(superadminsFromCSV)
	}

	superadminJobs := make(chan int, numSuperAdmin)
	var superadminWg sync.WaitGroup
	for sw := 0; sw < numWorkers; sw++ {
		superadminWg.Add(1)
		go func() {
			defer superadminWg.Done()
			for index := range superadminJobs {
				var u *users.User
				if index < len(superadminsFromCSV) {
					u = &superadminsFromCSV[index]
				}
				createSuperAdmin(index, passwordHash, u)
			}
		}()
	}
	for i := 0; i < numSuperAdmin; i++ {
		superadminJobs <- i
	}
	close(superadminJobs)
	superadminWg.Wait()

	// Process Counselors (Admins)
	var adminsFromCSV []users.User
	if *adminsCSV != "" {
		adminsFromCSV, err = parseUsersFromCSV(*adminsCSV)
		if err != nil {
			log.Fatal(err)
		}
		numCounselors = len(adminsFromCSV)
	}

	counselorJobs := make(chan int, numCounselors)
	var counselorWg sync.WaitGroup
	for cw := 0; cw < numWorkers; cw++ {
		counselorWg.Add(1)
		go func() {
			defer counselorWg.Done()
			for index := range counselorJobs {
				var u *users.User
				if index < len(adminsFromCSV) {
					u = &adminsFromCSV[index]
				}
				createCounselor(index, passwordHash, u)
			}
		}()
	}
	for i := 0; i < numCounselors; i++ {
		counselorJobs <- i
	}
	close(counselorJobs)
	counselorWg.Wait()

	// Process Students
	var studentsFromCSV []users.User
	if *studentsCSV != "" {
		studentsFromCSV, err = parseUsersFromCSV(*studentsCSV)
		if err != nil {
			log.Fatal(err)
		}
		numStudents = len(studentsFromCSV)
	}

	appointmentsDataset := make([]map[string]string, 0)
	if *appointmentsCSV != "" {
		appointmentsDataset, err = parseAppointmentsDatasetFromCSV(
			(*appointmentsCSV),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	rand.Shuffle(len(appointmentsDataset), func(i, j int) {
		appointmentsDataset[i], appointmentsDataset[j] = appointmentsDataset[j], appointmentsDataset[i]
	})

	jobs := make(chan int, numStudents)
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				var u *users.User
				if index < len(studentsFromCSV) {
					u = &studentsFromCSV[index]
				}
				createStudent(index, passwordHash, u, appointmentsDataset)
			}
		}()
	}
	for i := 0; i < numStudents; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	// Process Developers
	var devsFromCSV []users.User
	if *developersCSV != "" {
		devsFromCSV, err = parseUsersFromCSV(*developersCSV)
		if err != nil {
			log.Fatal(err)
		}
		numDevelopers = len(devsFromCSV)
	}

	developerJobs := make(chan int, numDevelopers)
	var developerWg sync.WaitGroup
	for dw := 0; dw < numWorkers; dw++ {
		developerWg.Add(1)
		go func() {
			defer developerWg.Done()
			for index := range developerJobs {
				var u *users.User
				if index < len(devsFromCSV) {
					u = &devsFromCSV[index]
				}
				createDeveloper(index, passwordHash, u)
			}
		}()
	}
	for i := 0; i < numDevelopers; i++ {
		developerJobs <- i
	}
	close(developerJobs)
	developerWg.Wait()

	fmt.Println("Dummy data generation completed successfully.")
	log.Println("Time taken:", time.Since(startTime))
}

// clearStudentData removes existing student data (keeps reference tables and
// seeded admin).
func clearStudentData() {
	tables := []string{
		"admission_slips",
		"slip_attachments",
		"profile_pictures",
		"student_cors",
		"files",
		"student_hobbies",
		"student_subject_preferences",
		"student_activities",
		"student_financial_supports",
		"student_finances",
		"significant_notes",
		"test_results",
		"student_consultations",
		"student_health_records",
		"school_details",
		"educational_backgrounds",
		"student_sibling_supports",
		"family_backgrounds",
		"student_related_persons",
		"related_persons",
		"student_addresses",
		"addresses",
		"student_selected_reasons",
		"emergency_contacts",
		"student_personal_info",
		"iir_records",
		"appointments",
		"notifications",
		"users",
	}

	// disable FK checks temporarily
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, tbl := range tables {
		if _, err := db.Exec(
			fmt.Sprintf("TRUNCATE TABLE %s", tbl),
		); err != nil {
			log.Fatal(err)
		}
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

// func runBackfill(passwordHash string) {
// 	ctx := context.Background()
// 	fmt.Println("Checking for students missing IIR records...")

// 	// Find students (role_id=1) who DON'T have a record in iir_records
// 	query := `
// 		SELECT u.* FROM users u
// 		LEFT JOIN iir_records iir ON u.id = iir.user_id
// 		WHERE u.role_id = 1 AND iir.id IS NULL
// 	`

// 	var missingStudents []users.User
// 	err := db.SelectContext(ctx, &missingStudents, query)
// 	if err != nil {
// 		log.Fatalf("[Backfill] Failed to query missing students: %v", err)
// 	}

// 	count := len(missingStudents)
// 	if count == 0 {
// 		fmt.Println("No students missing IIR records. Everything is in sync!")
// 		return
// 	}

// 	fmt.Printf("Found %d students missing IIRs. Starting backfill...\n", count)

// 	for i, u := range missingStudents {
// 		// generate core data needed for IIR seeder
// 		dob := gofakeit.DateRange(
// 			time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
// 			time.Now().AddDate(-18, 0, 0),
// 		)
// 		birthYear := dob.Year()

// 		tx, err := db.BeginTxx(ctx, nil)
// 		if err != nil {
// 			log.Fatalf("[Backfill] Failed to start transaction: %v", err)
// 		}

// 		// Ensure we don't double notify or something, but generateFullStudentIIR
// 		// handles the rest.
// 		generateFullStudentIIR(ctx, tx, u.ID, dob, birthYear, i)

// 		if err := tx.Commit(); err != nil {
// 			tx.Rollback()
// 			log.Fatalf("[Backfill] Failed to commit student %s: %v", u.ID, err)
// 		}
// 	}
// }
