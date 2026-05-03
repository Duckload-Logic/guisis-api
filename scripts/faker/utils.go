package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

func randomMiddleName() structs.NullableString {
	return nullStringIf(gofakeit.Bool(), gofakeit.FirstName())
}

func nullStringIf(cond bool, val string) structs.NullableString {
	if cond && val != "" {
		return structs.NullableString{String: val, Valid: true}
	}
	return structs.NullableString{Valid: false}
}

func fakePasswordHash() string {
	// static hash for all dummy users
	return "$2y$10$gxeDD.IKlEkqJmqmyVxy6eU9tFvC4ZK8KL3VZc2ex3BvNLo8DL5Dq"
}

func stringToNullString(s string) structs.NullableString {
	if s == "" {
		return structs.NullableString{Valid: false}
	}
	return structs.NullableString{String: s, Valid: true}
}

func stringToNullTime(s string) structs.NullableTime {
	if s == "" {
		return structs.NullableTime{Valid: false}
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		// attempt date only
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			return structs.NullableTime{Valid: false}
		}
	}
	return structs.FromSqlNullTime(sql.NullTime{Time: t, Valid: true})
}

// helper: random choice from slice
func randomChoice(slice interface{}) interface{} {
	switch v := slice.(type) {
	case []int:
		return v[rand.Intn(len(v))]
	case []string:
		return v[rand.Intn(len(v))]
	default:
		panic("unsupported type")
	}
}

func pickUniqueStrings(pool []string, count int) []string {
	if count <= 0 || len(pool) == 0 {
		return []string{}
	}
	if count >= len(pool) {
		copyPool := append([]string{}, pool...)
		rand.Shuffle(len(copyPool), func(i, j int) {
			copyPool[i], copyPool[j] = copyPool[j], copyPool[i]
		})
		return copyPool
	}

	shuffled := append([]string{}, pool...)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled[:count]
}

func relationshipID(name string) int {
	if id, ok := relationshipTypeIDs[name]; ok {
		return id
	}
	for relName, id := range relationshipTypeIDs {
		if strings.EqualFold(relName, name) {
			return id
		}
	}
	return 0
}

func buildDSNFromEnv() string {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		user,
		pass,
		host,
		port,
		name,
	)
	if os.Getenv("DB_TLS") == "true" {
		dsn += "&tls=true"
	}
	return dsn
}
