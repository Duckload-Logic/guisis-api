package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
)

// createSuperAdmin creates a super admin user.
func createSuperAdmin(index int, password string, userFromCSV *users.User) {
	var user users.User
	if userFromCSV != nil {
		user = *userFromCSV
	} else {
		email := fmt.Sprintf("superadmin%d@gmail.com", index+1)

		user = users.User{
			ID:           uuid.New().String(),
			Roles:        []users.Role{{ID: 3, Name: "SuperAdmin"}},
			FirstName:    "Super",
			MiddleName:   randomMiddleName(),
			LastName:     "Admin",
			SuffixName:   structs.NullableString{Valid: false},
			Email:        email,
			PasswordHash: structs.NullableString{Valid: true, String: password},
			AuthType:     "native",
			IsActive:     true,
		}
	}

	err := usersRepo.CreateUser(context.Background(), db, user)
	if err != nil {
		log.Fatalf("[Seeder] {Create SuperAdmin}: %v", err)
	}

	fmt.Printf(
		"Created super admin: %s %s (%s / SuperAdmin@123)\n",
		user.FirstName,
		user.LastName,
		user.Email,
	)
}

// createCounselor creates a counselor (admin) user and profile.
func createCounselor(index int, password string, userFromCSV *users.User) {
	var user users.User
	if userFromCSV != nil {
		user = *userFromCSV
	} else {
		firstName := gofakeit.FirstName()
		lastName := gofakeit.LastName()
		email := fmt.Sprintf("counselor%d@gmail.com", index+1)

		user = users.User{
			ID:           uuid.New().String(),
			Roles:        []users.Role{{ID: 2, Name: "Counselor"}},
			FirstName:    firstName,
			MiddleName:   randomMiddleName(),
			LastName:     lastName,
			SuffixName:   nullStringIf(rand.Float32() < 0.1, gofakeit.RandomString([]string{"Jr.", "Sr.", "III", "IV"})),
			Email:        email,
			PasswordHash: structs.NullableString{Valid: true, String: password},
			AuthType:     "native",
			IsActive:     true,
		}
	}

	err := usersRepo.CreateUser(context.Background(), db, user)
	if err != nil {
		log.Fatalf("[Seeder] {Create Counselor User}: %v", err)
	}

	fmt.Printf("Created counselor: %s %s\n", user.FirstName, user.LastName)
}
