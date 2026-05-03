package m2mclients

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) isAdmin(c *gin.Context) bool {
	roles, ok := c.Get("roleIDs")
	if !ok {
		return false
	}
	roleIDs := roles.([]int)
	for _, id := range roleIDs {
		if id == int(constants.SuperAdminRoleID) {
			return true
		}
	}
	return false
}

func (h *Handler) PostM2MClient(c *gin.Context) {
	var req CreateM2MClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	resp, err := h.service.CreateClient(c.Request.Context(), userID, req)
	if err != nil {
		response.SendError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	response.SendSuccess(c, resp)
}

func (h *Handler) PostM2MToken(c *gin.Context) {
	var req M2MTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Authenticate(c.Request.Context(), req.ClientID, req.ClientSecret)
	if err != nil {
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	response.SendSuccess(c, resp)
}

func (h *Handler) PostM2MTokenRefresh(c *gin.Context) {
	response.SendError(c, "Not implemented", http.StatusNotImplemented, nil)
}

func (h *Handler) GetM2MClients(c *gin.Context) {
	includeRevoked := c.Query("include_revoked")

	clients, err := h.service.ListClients(c.Request.Context(), includeRevoked == "true")
	if err != nil {
		response.SendError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	if clients == nil {
		clients = make([]M2MClient, 0)
	}

	response.SendSuccess(c, clients)
}

func (h *Handler) GetMyM2MClient(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	client, err := h.service.GetClientByUserID(c.Request.Context(), userID)
	if err != nil {
		response.SendError(c, err.Error(), http.StatusNotFound, nil)
		return
	}
	response.SendSuccess(c, client)
}

func (h *Handler) PostM2MClientSecret(c *gin.Context) {
	clientID := c.Param("id")
	userID := c.MustGet("userID").(string)

	secret, err := h.service.ResetSecret(
		c.Request.Context(),
		clientID,
		userID,
		h.isAdmin(c),
	)
	if err != nil {
		response.SendError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	response.SendSuccess(c, gin.H{"clientSecret": secret})
}

func (h *Handler) DeleteM2MClient(c *gin.Context) {
	id := c.Param("id")
	userID := c.MustGet("userID").(string)

	err := h.service.Deactivate(
		c.Request.Context(),
		id,
		userID,
		h.isAdmin(c),
	)
	if err != nil {
		response.SendError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	response.SendSuccess(c, gin.H{"message": "Deactivated"})
}

func (h *Handler) PatchM2MClientVerify(c *gin.Context) {
	id := c.Param("id")
	err := h.service.Verify(c.Request.Context(), id)
	if err != nil {
		response.SendError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	response.SendSuccess(c, gin.H{"message": "Verified"})
}
