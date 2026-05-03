package hash

import "testing"

func TestGetSHA256Hash(t *testing.T) {
	tests := []struct {
		name  string
		input string
		size  int
		want  string
	}{
		{
			name:  "standard hash full",
			input: "hello",
			size:  32,
			want:  "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name:  "standard hash truncated",
			input: "hello",
			size:  8,
			want:  "2cf24dba5fb0a30e",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSHA256Hash(tt.input, tt.size)
			if got != tt.want {
				t.Errorf("GetSHA256Hash() = %s, want %s", got, tt.want)
			}
		})
	}
}
