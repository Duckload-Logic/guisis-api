package appointments

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

func TestGenerateAppointmentsCSV(t *testing.T) {
	data, err := generateAppointmentsCSV([]AppointmentWithDetailsView{{
		WhenDate: "2026-08-01", TimeSlotTime: "09:00", StudentNumber: "=A1",
		UserFirstName: "Ana", UserLastName: "Santos", UserEmail: "ana@example.com",
		CategoryName: "Counseling", StatusName: "Pending", UrgencyLevel: "High",
		Reason: structs.StringToNullableString("+formula"), CreatedAt: time.Date(2026, 8, 1, 9, 30, 0, 0, time.UTC),
	}})
	if err != nil {
		t.Fatal(err)
	}

	records, err := csv.NewReader(strings.NewReader(string(data))).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 || len(records[0]) != 10 {
		t.Fatalf("unexpected CSV shape: %#v", records)
	}
	if records[1][2] != "'=A1" || records[1][8] != "'+formula" {
		t.Fatalf("unsafe cells were not escaped: %#v", records[1])
	}
}
