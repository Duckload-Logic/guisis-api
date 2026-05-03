package structs

func TruncateString(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}
