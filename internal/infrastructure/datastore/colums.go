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
	hasID := false
	for _, col := range cols {
		if col == "id" {
			hasID = true
			continue
		}
		if !contains(exclude, col) &&
			!contains(excludeOnUpsert, col) {
			updates = append(
				updates,
				fmt.Sprintf("%s = VALUES(%s)", col, col),
			)
		}
	}

	if hasID {
		t := reflect.TypeOf(s)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		isIntegerID := false
		if t.Kind() == reflect.Struct {
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				tag := field.Tag.Get("db")
				if tag == "id" ||
					(tag == "" && strings.ToLower(field.Name) == "id") {
					k := field.Type.Kind()
					switch k {
					case reflect.Int, reflect.Int8, reflect.Int16,
						reflect.Int32, reflect.Int64, reflect.Uint,
						reflect.Uint8, reflect.Uint16, reflect.Uint32,
						reflect.Uint64:
						isIntegerID = true
					}
					break
				}
			}
		}
		if isIntegerID {
			updates = append(updates, "id = LAST_INSERT_ID(id)")
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
