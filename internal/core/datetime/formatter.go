package datetime

import (
	"strings"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

func FormatTime(t string) string {
	layouts := []string{"15:04:05", "15:04"}
	var parsedTime time.Time
	var err error

	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, t)
		if err == nil {
			return parsedTime.Format("3:04 PM")
		}
	}

	return ""
}

func FormatDate(t string) string {
	layouts := []string{time.RFC3339, constants.LayoutDateOnly}
	var parsedTime time.Time
	var err error

	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, t)
		if err == nil {
			return parsedTime.Format("January 2, 2006")
		}
	}

	return ""
}

// ExtractDateOnly extracts the YYYY-MM-DD date portion from an ISO string.
func ExtractDateOnly(t string) string {
	if idx := strings.Index(t, "T"); idx != -1 {
		return t[:idx]
	}
	return t
}

// FormatDateTime formats the given time into YYYY-MM-DD HH:MM:SS format.
func FormatDateTime(t time.Time) string {
	return t.Format(constants.LayoutDateTime)
}

// GetTodayInPHT returns the start of the current day in Philippine Time,
// converted to UTC representation.
func GetTodayInPHT() time.Time {
	loc := time.FixedZone(constants.PHTZoneName, constants.PHTZoneOffset)
	now := time.Now().In(loc)
	return time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC,
	)
}
