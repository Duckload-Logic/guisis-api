package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

// parseUsersFromCSV reads a CSV file and returns a slice of users.User objects.
// Expected Headers: id,email,role_id,first_name,middle_name,last_name,
// suffix_name,password_hash,auth_type,is_active,created_at,updated_at
func parseUsersFromCSV(filePath string) ([]users.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// read header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var usersList []users.User
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

	// map record to User struct
		// id
		// email
		// first_name
		// middle_name
		// last_name
		// suffix_name
		// password_hash
		// auth_type
		// is_active
		// created_at
		// updated_at

		isActiveInt, _ := strconv.Atoi(record[8])
		isActive := isActiveInt != 0

		user := users.User{
			ID:           record[0],
			Email:        record[1],
			Roles:        []users.Role{{ID: 1}}, // Default: Student
			FirstName:    record[2],
			MiddleName:   stringToNullString(record[3]),
			LastName:     record[4],
			SuffixName:   stringToNullString(record[5]),
			PasswordHash: stringToNullString(record[6]),
			AuthType:     record[7],
			IsActive:     isActive,
			CreatedAt:    stringToNullTime(record[9]),
			UpdatedAt:    stringToNullTime(record[10]),
		}

		usersList = append(usersList, user)
	}

	return usersList, nil
}

func parseAppointmentsDatasetFromCSV(
	filePath string,
) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// read header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var appointmentsDataset []map[string]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV record: %w", err)
		}

		data := map[string]string{
			"text":          record[0],
			"urgency_level": strings.ToUpper(record[1]),
			"category":      strings.ToLower(record[2]),
		}

		appointmentsDataset = append(appointmentsDataset, data)
	}

	return appointmentsDataset, nil
}
