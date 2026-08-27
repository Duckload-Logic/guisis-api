package students

import (
	"testing"
)

func TestCleanStudentNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2021-00123-TG-0", "202100123tg0"},
		{"2021 00123 TG 0", "202100123tg0"},
		{"2021_00123_tg_0", "202100123tg0"},
		{"2021-00123-tg-0", "202100123tg0"},
		{"", ""},
		{"abc-123", "abc123"},
	}

	for _, tc := range tests {
		result := cleanStudentNumber(tc.input)
		if result != tc.expected {
			t.Errorf(
				"cleanStudentNumber(%q) = %q; expected %q",
				tc.input, result, tc.expected,
			)
		}
	}
}
