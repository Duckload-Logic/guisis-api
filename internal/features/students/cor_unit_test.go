package students

import (
	"testing"
)

func TestKeepOnlyDigits(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2021-00123-TG-0", "2021001230"},
		{"2021 00123 TG 0", "2021001230"},
		{"2021_00123_tg_0", "2021001230"},
		{"2021-00123-tg-0", "2021001230"},
		{"", ""},
		{"abc-123", "123"},
	}

	for _, tc := range tests {
		result := keepOnlyDigits(tc.input)
		if result != tc.expected {
			t.Errorf(
				"keepOnlyDigits(%q) = %q; expected %q",
				tc.input, result, tc.expected,
			)
		}
	}
}

func TestMatchStudentNumbers(t *testing.T) {
	tests := []struct {
		db       string
		ocr      string
		expected bool
	}{
		{"2023-00122-TG-0", "2023-00122-TG-0", true},
		{"2023-00122-TG-0", "2023-00122-TG-O", true}, // OCR typo check digit
		{"2023-00122-TG-0", "2023-00122-TO-0", true}, // OCR typo campus code
		{"2023-00122-TG-0", "2023-00122-TG", true},
		{"2023-00122-TG-0", "00122-TG-0", true},
		{"2023-00122-TG-0", "2023-00122", true},
		{"2023-00122-TG-0", "2023-00123-TG-0", false},
		{"", "2023-00122-TG-0", true},
		{"2023-00122-TG-0", "", true},
		{"2023", "2023-00122-TG-0", false},
		{"2023-00122-TG-0", "tg", false},
	}

	for _, tc := range tests {
		result := matchStudentNumbers(tc.db, tc.ocr)
		if result != tc.expected {
			t.Errorf(
				"matchStudentNumbers(%q, %q) = %t; expected %t",
				tc.db, tc.ocr, result, tc.expected,
			)
		}
	}
}
