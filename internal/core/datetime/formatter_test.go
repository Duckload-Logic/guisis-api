package datetime

import "testing"

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"24h with seconds", "14:30:05", "2:30 PM"},
		{"24h without seconds", "14:30", "2:30 PM"},
		{"AM time", "09:15", "9:15 AM"},
		{"Invalid time", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatTime(tt.input)
			if got != tt.want {
				t.Errorf("FormatTime(%s) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"RFC3339", "2024-03-15T14:30:00Z", "March 15, 2024"},
		{"Simple ISO", "2024-03-15", "March 15, 2024"},
		{"Invalid date", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDate(tt.input)
			if got != tt.want {
				t.Errorf("FormatDate(%s) = %s, want %s", tt.input, got, tt.want)
			}
		})
	}
}
