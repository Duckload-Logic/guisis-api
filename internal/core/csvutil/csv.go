// Package csvutil contains small safety helpers for CSV downloads.
package csvutil

import "strings"

// EscapeCell prevents spreadsheet applications from interpreting untrusted
// values as formulas when a CSV file is opened. The apostrophe is not shown
// as part of the value by common spreadsheet applications.
func EscapeCell(value string) string {
	trimmed := strings.TrimLeft(value, " \t\r\n")
	if trimmed == "" {
		return value
	}

	switch trimmed[0] {
	case '=', '+', '-', '@':
		return "'" + value
	default:
		return value
	}
}
