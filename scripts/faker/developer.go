package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

func createDeveloper(index int, passwordHash string, userFromCSV *users.User) {
	var user users.User
	if userFromCSV != nil {
		user = *userFromCSV
	} else {
		// Generate unique identifier
		userID := uuid.New().String()

		// Generate developer data
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()
		email := fmt.Sprintf(
			"developer%d@gmail.com",
			index+1,
		)

		// Create developer user
		user = users.User{
			ID:         userID,
			Email:      email,
			FirstName:  firstName,
			LastName:   lastName,
			MiddleName: randomMiddleName(),
			SuffixName: structs.NullableString{Valid: false},
			PasswordHash: structs.NullableString{
				String: passwordHash,
				Valid:  true,
			},
			Roles: []users.Role{{ID: int(constants.DeveloperRoleID), Name: "Developer"}},
			AuthType: string(constants.AuthTypeNative),
			IsActive: true,
		}
	}

	// Save to database
	if err := usersRepo.CreateUser(context.Background(), db, user); err != nil {
		log.Printf("Failed to create developer %d: %v\n", index+1, err)
		return
	}

	log.Printf(
		"Created developer: %s %s (%s / %s)\n",
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash.String,
	)
}
