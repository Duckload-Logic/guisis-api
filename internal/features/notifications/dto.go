package notifications

import (
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type NotificationResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    []audit.NotificationEntry `json:"data"`
}

type ListNotificationsRequest struct {
	structs.PaginationRequest
	UnreadOnly bool `form:"unread_only" json:"unreadOnly"`
}

type ListNotificationsResponse struct {
	Notifications  []audit.NotificationEntry `json:"notifications"`
	Total          int                       `json:"total"`
	Page           int                       `json:"page"`
	PageSize       int                       `json:"pageSize"`
	TotalPages     int                       `json:"totalPages"`
	UnreadCount    int                       `json:"unreadCount"`
	UntouchedCount int                       `json:"untouchedCount"`
}
