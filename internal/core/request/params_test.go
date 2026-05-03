package request

import "testing"

func TestPaginationParams_GetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		want     int
	}{
		{"page 1 offset 0", 1, 10, 0},
		{"page 2 offset 10", 2, 10, 10},
		{"page 3 offset 20", 3, 10, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PaginationParams{Page: tt.page, PageSize: tt.pageSize}
			if got := p.GetOffset(); got != tt.want {
				t.Errorf("GetOffset() = %d, want %d", got, tt.want)
			}
		})
	}
}
