package structs

import "github.com/olazo-johnalbert/duckload-api/internal/core/constants"

// PaginationRequest represents a standard paginated request.
type PaginationRequest struct {
	Page      int    `json:"page"       form:"page"`
	PageSize  int    `json:"page_size"  form:"page_size"`
	OrderBy   string `json:"order_by"   form:"order_by"`
	SortOrder string `json:"sort_order" form:"sort_order"`
	Search    string `json:"search"     form:"search"`
}

// GetOffset calculates the SQL offset.
func (r *PaginationRequest) GetOffset() int {
	return (r.Page - 1) * r.PageSize
}

// SetDefaults sets default pagination values if not provided or invalid.
func (r *PaginationRequest) SetDefaults(defaultOrderBy string) {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = constants.DefaultPageSize
	}
	if r.PageSize > constants.MaxPageSize {
		r.PageSize = constants.MaxPageSize
	}
	if r.OrderBy == "" {
		r.OrderBy = defaultOrderBy
	}
	if r.SortOrder == "" {
		r.SortOrder = "desc"
	}
}

// PaginationMetadata contains calculated pagination information for responses.
type PaginationMetadata struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	PageSize   int `json:"pagesSize"`
	TotalPages int `json:"totalPages"`
}

// CalculateMetadata generates pagination metadata.
func CalculateMetadata(total, page, pageSize int) PaginationMetadata {
	totalPages := 0
	if pageSize > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}

	return PaginationMetadata{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
