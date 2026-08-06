package csvutil

import "testing"

func TestEscapeCell(t *testing.T) {
	tests := map[string]string{
		"plain text":      "plain text",
		"=SUM(A1:A2)":     "'=SUM(A1:A2)",
		" +cmd|'/C calc'": "' +cmd|'/C calc'",
		"\t@danger":       "'\t@danger",
		"-1":              "'-1",
	}

	for input, want := range tests {
		if got := EscapeCell(input); got != want {
			t.Errorf("EscapeCell(%q) = %q, want %q", input, got, want)
		}
	}
}
