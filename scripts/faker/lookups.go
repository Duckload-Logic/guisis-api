package main

import (
	"log"
	"strings"
)

func loadLookups() {
	// genders
	rows, err := db.Query("SELECT id FROM genders")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		genderIDs = append(genderIDs, id)
	}

	// civil status
	civilStatusByName = make(map[string]int)
	rows, err = db.Query("SELECT id, status_name FROM civil_status_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var statusName string
		if err := rows.Scan(&id, &statusName); err != nil {
			log.Fatal(err)
		}
		civilStatusIDs = append(civilStatusIDs, id)
		civilStatusByName[strings.ToLower(statusName)] = id
	}

	// religions
	rows, err = db.Query("SELECT id FROM religions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		religionIDs = append(religionIDs, id)
	}

	// courses
	rows, err = db.Query("SELECT id FROM courses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		courseIDs = append(courseIDs, id)
	}

	// enrollment reasons
	rows, err = db.Query("SELECT id FROM enrollment_reasons")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		enrollmentReasonIDs = append(enrollmentReasonIDs, id)
	}

	// student support types
	rows, err = db.Query("SELECT id FROM student_support_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		supportTypeIDs = append(supportTypeIDs, id)
	}

	// income ranges
	rows, err = db.Query("SELECT id FROM income_ranges")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		incomeRangeIDs = append(incomeRangeIDs, id)
	}

	// parental status
	rows, err = db.Query("SELECT id FROM parental_status_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		parentalStatusIDs = append(parentalStatusIDs, id)
	}

	// educational levels
	educationalLevelByName = make(map[string]int)
	rows, err = db.Query("SELECT id, level_name FROM educational_levels")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var levelName string
		if err := rows.Scan(&id, &levelName); err != nil {
			log.Fatal(err)
		}
		educationalLevelIDs = append(educationalLevelIDs, id)
		educationalLevelByName[strings.ToLower(levelName)] = id
	}

	// relationship types (map by name for easy lookup)
	relationshipTypeIDs = make(map[string]int)
	rows, err = db.Query(
		"SELECT id, relationship_name FROM student_relationship_types",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		relationshipTypeIDs[name] = id
	}

	// nature of residence
	rows, err = db.Query("SELECT id FROM nature_of_residence_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		natureOfResidenceIDs = append(natureOfResidenceIDs, id)
	}

	// sibling support types
	rows, err = db.Query("SELECT id FROM sibling_support_types")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		siblingSupportTypeIDs = append(siblingSupportTypeIDs, id)
	}

	// activity options
	rows, err = db.Query(
		"SELECT id FROM activity_options WHERE is_active = TRUE",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		activityOptionIDs = append(activityOptionIDs, id)
	}

	// appointment time slots
	rows, err = db.Query("SELECT id FROM time_slots")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		timeSlotIDs = append(timeSlotIDs, id)
	}

	// appointment statuses
	appointmentStatusByName = make(map[string]int)
	rows, err = db.Query(
		`SELECT id, name
		FROM statuses
		WHERE status_type IN ('appointment', 'both')`,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		appointmentStatusIDs = append(appointmentStatusIDs, id)
		appointmentStatusByName[strings.ToLower(name)] = id
	}

	admissionSlipStatusesByName = make(map[string]int)
	rows, err = db.Query(
		"SELECT id, name FROM statuses WHERE status_type IN ('slip', 'both')",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}

		admissionSlipStatusIDs = append(admissionSlipStatusIDs, id)
		admissionSlipStatusesByName[strings.ToLower(name)] = id
	}

	// appointment categories
	rows, err = db.Query("SELECT id, name FROM appointment_categories")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		rows.Scan(&id, &name)
		appointmentCategories = append(appointmentCategories, map[string]string{
			"id":   id,
			"name": name,
		})
	}

	// admission_slip_categories
	row, err := db.Query("SELECT id FROM admission_slip_categories")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var id int
		row.Scan(&id)
		admissionSlipCategoryIDs = append(admissionSlipCategoryIDs, id)
	}

	// student_statuses
	studentStatusByName = make(map[string]int)
	rows, err = db.Query("SELECT id, status_name FROM student_statuses")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		studentStatusIDs = append(studentStatusIDs, id)
		studentStatusByName[strings.ToLower(name)] = id
	}
}
