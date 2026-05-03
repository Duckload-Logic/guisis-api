package notifications

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetNotifications(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	notifications, err := h.service.GetUserNotifications(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		fmt.Printf("[GetNotifications] {Fetch Notifications}: %v\n", err)
		response.SendError(
			c,
			"Failed to fetch notifications",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	data := ListNotificationsResponse{
		Notifications: notifications,
		Total:         len(notifications),
		Page:          1,
		PageSize:      len(notifications),
		TotalPages:    1,
	}

	response.SendSuccess(c, data)
}

func (h *Handler) GetNotificationsStream(c *gin.Context) {
	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// Force headers to be sent immediately
	c.Writer.Flush()

	ch, unsubscribe := h.service.Subscribe(
		c.Request.Context(),
		c.MustGet("userID").(string),
	)

	defer unsubscribe()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-ticker.C:
			// Heartbeat to keep connection alive and detect broken pipes
			_, err := c.Writer.Write([]byte(": heartbeat\n\n"))
			if err != nil {
				return
			}
			c.Writer.Flush()
		case notif, ok := <-ch:
			if !ok {
				return
			}
			b, err := json.Marshal(notif)
			if err != nil {
				fmt.Printf(
					"[GetNotificationsStream] {Marshal}: %v\n",
					err,
				)
				continue
			}

			_, err = c.Writer.Write([]byte("data: " + string(b) + "\n\n"))
			if err != nil {
				return
			}

			c.Writer.Flush()
		}
	}
}

func (h *Handler) PatchNotificationRead(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(string)

	if err := h.service.MarkAsRead(c.Request.Context(), id, userID); err != nil {
		fmt.Printf("[PatchNotificationRead] {Mark Read}: %v\n", err)
		response.SendError(
			c,
			"Failed to mark notification as read",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "Notification marked as read"})
}
