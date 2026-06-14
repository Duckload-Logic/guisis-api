package notifications

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
				fmt.Printf(
					"[GetNotificationsStream] {Write Heartbeat}: %v\n",
					err,
				)
				return
			}
			c.Writer.Flush()
		case notif, ok := <-ch:
			if !ok {
				return
			}
			b, err := json.Marshal(notif)
			if err != nil {
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

func (h *Handler) PatchNotificationsTouched(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	err := h.service.MarkAllAsTouched(c.Request.Context(), userID)
	if err != nil {
		fmt.Printf(
			"[PatchNotificationsTouched] {MarkAllAsTouched}: %s\n",
			err.Error(),
		)
		response.SendError(
			c,
			"Failed to mark notifications as touched",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "Notifications marked as touched"},
	)
}

type pushSubscribeRequest struct {
	Endpoint  string `json:"endpoint"  binding:"required"`
	P256dhKey string `json:"p256dhKey" binding:"required"`
	AuthKey   string `json:"authKey"   binding:"required"`
}

func (h *Handler) PostPushSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var req pushSubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf(
			"[PostPushSubscription] {ShouldBindJSON}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Invalid request payload",
			http.StatusBadRequest,
			nil,
		)
		return
	}

	sub := PushSubscription{
		ID:        uuid.New().String(),
		UserID:    userID,
		Endpoint:  req.Endpoint,
		P256dhKey: req.P256dhKey,
		AuthKey:   req.AuthKey,
	}

	err := h.service.repo.SavePushSubscription(
		c.Request.Context(),
		nil,
		sub,
	)
	if err != nil {
		fmt.Printf(
			"[PostPushSubscription] {SavePushSubscription}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to save push subscription",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "Push subscription saved successfully"},
	)
}

func (h *Handler) DeletePushSubscription(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	endpoint := c.Query("endpoint")

	if endpoint == "" {
		fmt.Println("[DeletePushSubscription] {CheckQuery}: missing endpoint")
		response.SendError(
			c,
			"Endpoint query parameter is required",
			http.StatusBadRequest,
			nil,
		)
		return
	}

	err := h.service.repo.DeletePushSubscription(
		c.Request.Context(),
		nil,
		endpoint,
		userID,
	)
	if err != nil {
		fmt.Printf(
			"[DeletePushSubscription] {DeletePushSubscription}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to delete push subscription",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "Push subscription deleted successfully"},
	)
}

