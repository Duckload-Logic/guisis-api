package slips

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

func TestGenerateSlipsCSV(t *testing.T) {
	data, err := generateSlipsCSV([]SlipWithDetailsView{{
		DateNeeded: "2026-08-01", DateOfAbsence: "2026-07-31", StudentNumber: "2026-001",
		UserFirstName: "Ana", UserMiddleName: structs.StringToNullableString("D."), UserLastName: "Santos",
		UserEmail: "ana@example.com", CategoryName: "Medical", StatusName: "Pending", Reason: "@formula",
		TicketCode: structs.StringToNullableString("T-1"), CreatedAt: time.Date(2026, 8, 1, 9, 30, 0, 0, time.UTC),
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
	if records[1][3] != "Ana D. Santos" || records[1][7] != "'@formula" {
		t.Fatalf("unexpected CSV row: %#v", records[1])
	}
}
