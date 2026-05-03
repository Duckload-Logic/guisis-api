package datastore

import (
	"fmt"
	"reflect"
	"strings"
)

var excludeOnUpsert = []string{"id", "created_at", "updated_at"}

func GetColumns(s interface{}) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("db")

		if tag == "" || tag == "-" {
			continue
		}
		columns = append(columns, tag)
	}
	return strings.Join(columns, ", ")
}

func GetInsertStatement(s interface{}, exclude []string) (string, string) {
	cols := strings.Split(GetColumns(s), ", ")
	var filteredCols []string
	var placeholders []string

	for _, col := range cols {
		if !contains(excludeOnUpsert, col) && !contains(exclude, col) {
			filteredCols = append(filteredCols, col)
			placeholders = append(placeholders, ":"+col)
		}
	}

	return strings.Join(filteredCols, ", "), strings.Join(placeholders, ", ")
}

func GetOnDuplicateKeyUpdateStatement(
	s interface{}, exclude []string,
) string {
	cols := strings.Split(GetColumns(s), ", ")
	var updates []string
	for _, col := range cols {
		if !contains(exclude, col) &&
			!contains(excludeOnUpsert, col) {
			updates = append(
				updates,
				fmt.Sprintf("%s = VALUES(%s)", col, col),
			)
		}
	}

	return strings.Join(updates, ", ")
}

func GetPrefixColumns(s interface{}, prefix string) string {
	cols := strings.Split(GetColumns(s), ", ")
	var prefixedCols []string
	for _, col := range cols {
		prefixedCols = append(prefixedCols, fmt.Sprintf("%s.%s", prefix, col))
	}

	return strings.Join(prefixedCols, ", ")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
