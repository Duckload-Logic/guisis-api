package structs

import (
	"testing"

	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

func TestPaginationRequest_GetOffset(t *testing.T) {
	req := PaginationRequest{Page: 2, PageSize: 10}
	if got := req.GetOffset(); got != 10 {
		t.Errorf("GetOffset() = %d, want 10", got)
	}
}

func TestPaginationRequest_SetDefaults(t *testing.T) {
	tests := []struct {
		name           string
		input          PaginationRequest
		defaultOrderBy string
		wantPage       int
		wantSize       int
		wantOrder      string
	}{
		{
			name:           "all defaults",
			input:          PaginationRequest{},
			defaultOrderBy: "id",
			wantPage:       1,
			wantSize:       constants.DefaultPageSize,
			wantOrder:      "id",
		},
		{
			name:           "clamped size",
			input:          PaginationRequest{PageSize: 9999},
			defaultOrderBy: "id",
			wantPage:       1,
			wantSize:       constants.MaxPageSize,
			wantOrder:      "id",
		},
		{
			name: "preserve valid",
			input: PaginationRequest{
				Page:     5,
				PageSize: 50,
				OrderBy:  "name",
			},
			defaultOrderBy: "id",
			wantPage:       5,
			wantSize:       50,
			wantOrder:      "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.SetDefaults(tt.defaultOrderBy)
			if tt.input.Page != tt.wantPage {
				t.Errorf("Page = %d, want %d", tt.input.Page, tt.wantPage)
			}
			if tt.input.PageSize != tt.wantSize {
				t.Errorf(
					"PageSize = %d, want %d",
					tt.input.PageSize,
					tt.wantSize,
				)
			}
			if tt.input.OrderBy != tt.wantOrder {
				t.Errorf(
					"OrderBy = %s, want %s",
					tt.input.OrderBy,
					tt.wantOrder,
				)
			}
		})
	}
}

func TestCalculateMetadata(t *testing.T) {
	tests := []struct {
		name      string
		total     int
		page      int
		pageSize  int
		wantPages int
	}{
		{"exact division", 100, 1, 10, 10},
		{"with remainder", 105, 1, 10, 11},
		{"zero total", 0, 1, 10, 0},
		{"zero pageSize", 100, 1, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateMetadata(tt.total, tt.page, tt.pageSize)
			if got.TotalPages != tt.wantPages {
				t.Errorf(
					"TotalPages = %d, want %d",
					got.TotalPages,
					tt.wantPages,
				)
			}
			if got.Total != tt.total || got.Page != tt.page ||
				got.PageSize != tt.pageSize {
				t.Errorf("Metadata mismatch: %+v", got)
			}
		})
	}
}
