package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		"Local",
	)
	if dsn == "::@tcp(:)/?parseTime=true&loc=Local" {
		dsn = "mysqladmin:secure@tcp(localhost:3307)/guidance_db?parseTime=true&loc=Local"
	}

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Verifying temporal constraints...")

	// Future Appointment Inconsistencies
	var futureApptErr int
	query1 := `
		SELECT COUNT(*) FROM appointments
		WHERE when_date > CURDATE()
		AND status_id IN (SELECT id FROM statuses WHERE LOWER(name) IN ('completed', 'no-show'))
	`
	db.Get(&futureApptErr, query1)
	fmt.Printf(
		"Future appts with 'Completed/No-show' (Error): %d\n",
		futureApptErr,
	)

	// Past Appointment Inconsistencies
	var pastApptErr int
	query2 := `
		SELECT COUNT(*) FROM appointments
		WHERE when_date < CURDATE()
		AND status_id IN (SELECT id FROM statuses WHERE LOWER(name) IN ('pending', 'scheduled'))
	`
	db.Get(&pastApptErr, query2)
	fmt.Printf("Past appts with 'Pending/Scheduled' (Error): %d\n", pastApptErr)

	// Future Slip Inconsistencies
	var futureSlipErr int
	query3 := `
		SELECT COUNT(*) FROM admission_slips
		WHERE date_of_absence > CURDATE()
		AND status_id IN (SELECT id FROM statuses WHERE LOWER(name) IN ('approved', 'rejected', 'for revision'))
	`
	db.Get(&futureSlipErr, query3)
	fmt.Printf(
		"Future slips with 'Processed' status (Error): %d\n",
		futureSlipErr,
	)

	if futureApptErr == 0 && pastApptErr == 0 && futureSlipErr == 0 {
		fmt.Println("SUCCESS: All temporal constraints verified.")
	} else {
		fmt.Println("FAILURE: Some records violate temporal constraints.")
	}
}
