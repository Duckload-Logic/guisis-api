package support

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PostSupportTicket(c *gin.Context) {
	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf(
			"[PostSupportTicket] {BindJSON}: %v\n",
			err,
		)
		response.SendError(
			c, "Invalid request body", http.StatusBadRequest, nil,
		)
		return
	}

	if len(strings.Fields(req.Message)) > 100 {
		fmt.Println(
			"[PostSupportTicket] {ValidateMessage}: " +
				"Message exceeds 100 words limit",
		)
		response.SendError(
			c,
			"Message cannot exceed 100 words",
			http.StatusBadRequest,
			nil,
		)
		return
	}

	authUserID := h.getOptionalUserID(c)
	res, err := h.service.OpenTicket(
		c.Request.Context(), req, authUserID,
	)
	if err != nil {
		fmt.Printf(
			"[PostSupportTicket] {OpenTicket}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to open support ticket",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, res)
}

func (h *Handler) PostSupportMessage(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		fmt.Println(
			"[PostSupportMessage] {ValidateID}: Ticket ID is required",
		)
		response.SendError(
			c, "Ticket ID is required", http.StatusBadRequest, nil,
		)
		return
	}

	if _, ok := h.validateTicketAccess(c, ticketID); !ok {
		return
	}

	var req CreateMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf(
			"[PostSupportMessage] {BindJSON}: %v\n",
			err,
		)
		response.SendError(
			c, "Invalid request body", http.StatusBadRequest, nil,
		)
		return
	}

	if len(strings.Fields(req.Message)) > 100 {
		fmt.Println(
			"[PostSupportMessage] {ValidateMessage}: " +
				"Message exceeds 100 words limit",
		)
		response.SendError(
			c,
			"Message cannot exceed 100 words",
			http.StatusBadRequest,
			nil,
		)
		return
	}

	authUserID := h.getOptionalUserID(c)
	res, err := h.service.AddMessage(
		c.Request.Context(),
		ticketID,
		req,
		authUserID,
	)
	if err != nil {
		fmt.Printf(
			"[PostSupportMessage] {AddMessage}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to send message",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, res)
}

func (h *Handler) GetSupportTicketMessages(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		fmt.Println(
			"[GetSupportTicketMessages] {ValidateID}: Ticket ID is required",
		)
		response.SendError(
			c, "Ticket ID is required", http.StatusBadRequest, nil,
		)
		return
	}

	if _, ok := h.validateTicketAccess(c, ticketID); !ok {
		return
	}

	res, err := h.service.GetTicketMessages(
		c.Request.Context(), ticketID,
	)
	if err != nil {
		fmt.Printf(
			"[GetSupportTicketMessages] {GetTicketMessages}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to get messages",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, res)
}

func (h *Handler) GetSupportTickets(c *gin.Context) {
	staffUserID := h.getOptionalUserID(c)
	if staffUserID == "" {
		response.SendError(
			c,
			"Unauthorized: missing user session",
			http.StatusUnauthorized,
			nil,
		)
		return
	}

	var req GetTicketsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Printf(
			"[GetSupportTickets] {BindQuery}: %v\n",
			err,
		)
		response.SendFail(c, gin.H{"error": "Invalid query parameters"})
		return
	}
	req.SetDefaults("updated_at")

	var statusFilter string
	switch strings.ToLower(req.Status) {
	case "open":
		statusFilter = "OPEN"
	case "resolved", "closed":
		statusFilter = "CLOSED"
	default:
		statusFilter = ""
	}

	res, err := h.service.GetTickets(
		c.Request.Context(),
		staffUserID,
		req.PaginationRequest,
		statusFilter,
	)
	if err != nil {
		fmt.Printf(
			"[GetSupportTickets] {GetTickets}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to get support tickets",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, res)
}

func (h *Handler) PatchSupportTicketRead(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		response.SendError(
			c,
			"Ticket ID is required",
			http.StatusBadRequest,
			nil,
		)
		return
	}

	staffUserID := h.getOptionalUserID(c)
	if staffUserID == "" {
		response.SendError(
			c,
			"Unauthorized: missing user session",
			http.StatusUnauthorized,
			nil,
		)
		return
	}

	err := h.service.MarkTicketAsRead(
		c.Request.Context(),
		ticketID,
		staffUserID,
	)
	if err != nil {
		fmt.Printf(
			"[PatchSupportTicketRead] {MarkTicketAsRead}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to mark ticket as read",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "Ticket marked as read"},
	)
}

func (h *Handler) GetMySupportTickets(c *gin.Context) {
	authUserID := h.getOptionalUserID(c)
	if authUserID == "" {
		response.SendError(
			c,
			"Unauthorized: missing user session",
			http.StatusUnauthorized,
			nil,
		)
		return
	}

	res, err := h.service.GetTicketsByUserID(
		c.Request.Context(),
		authUserID,
	)
	if err != nil {
		fmt.Printf(
			"[GetMySupportTickets] {GetTicketsByUserID}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to get your support tickets",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, res)
}

func (h *Handler) PatchSupportTicketStatus(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		fmt.Println(
			"[PatchSupportTicketStatus] {ValidateID}: Ticket ID is required",
		)
		response.SendError(
			c, "Ticket ID is required", http.StatusBadRequest, nil,
		)
		return
	}

	err := h.service.CloseTicket(c.Request.Context(), ticketID)
	if err != nil {
		fmt.Printf(
			"[PatchSupportTicketStatus] {CloseTicket}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to resolve ticket",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "Ticket marked as resolved"})
}

func (h *Handler) getOptionalClaims(
	c *gin.Context,
) (*tokens.Claims, error) {
	var tokenString string
	if cookie, err := c.Cookie("access_token"); err == nil &&
		cookie != "" {
		tokenString = cookie
	} else {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if tokenString == "" {
		return nil, fmt.Errorf("no token provided")
	}

	return tokens.NewService().ValidateToken(tokenString)
}

func (h *Handler) getOptionalUserID(c *gin.Context) string {
	claims, err := h.getOptionalClaims(c)
	if err != nil {
		return ""
	}
	return claims.UserID
}

func (h *Handler) validateTicketAccess(
	c *gin.Context,
	ticketID string,
) (*SupportTicket, bool) {
	ticket, err := h.service.GetTicket(
		c.Request.Context(),
		ticketID,
	)
	if err != nil {
		fmt.Printf(
			"[validateTicketAccess] {GetTicket}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Ticket not found",
			http.StatusNotFound,
			nil,
		)
		return nil, false
	}

	claims, err := h.getOptionalClaims(c)
	if err == nil {
		isStaff := false
		for _, roleID := range claims.RoleIDs {
			if roleID == int(constants.AdminRoleID) ||
				roleID == int(constants.SuperAdminRoleID) ||
				roleID == int(constants.DeveloperRoleID) {
				isStaff = true
				break
			}
		}

		if !isStaff {
			if !ticket.UserID.Valid ||
				claims.UserID != ticket.UserID.String {
				response.SendError(
					c,
					"Forbidden: you do not own this ticket",
					http.StatusForbidden,
					nil,
				)
				return nil, false
			}
		}
	} else {
		if ticket.UserID.Valid && ticket.UserID.String != "" {
			response.SendError(
				c,
				"Unauthorized: ticket belongs to a user",
				http.StatusUnauthorized,
				nil,
			)
			return nil, false
		}
	}

	return ticket, true
}
