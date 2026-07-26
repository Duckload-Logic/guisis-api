package m2mclients

import (
	"log"
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
	roleIDsVal, ok := c.Get("roleIDs")
	if !ok {
		return false
	}

	var roleIDs []int
	switch v := roleIDsVal.(type) {
	case []int:
		roleIDs = v
	case []interface{}:
		for _, item := range v {
			if f, ok := item.(float64); ok {
				roleIDs = append(roleIDs, int(f))
			} else if i, ok := item.(int); ok {
				roleIDs = append(roleIDs, i)
			}
		}
	}

	for _, id := range roleIDs {
		if id == int(constants.SuperAdminRoleID) ||
			id == int(constants.DeveloperRoleID) {
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

// @Summary     Get M2M access token using client credentials
// @Description Authenticate with client ID and client secret to get access and refresh tokens
// @Tags         M2M Clients
// @Accept       json
// @Produce      json
// @Param      credentials body M2MTokenRequest true "Client credentials"
// @Success 200  {object} M2MTokenResponse
// @Failure 401  {object} response.CommonErrorResponse "Invalid client credentials"
// @Failure 500  {object} response.CommonErrorResponse "Internal server error"
// @Router /auth/m2m/token [post]
func (h *Handler) PostM2MToken(c *gin.Context) {
	var req M2MTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Authenticate(
		c.Request.Context(),
		req.ClientID,
		req.ClientSecret,
	)
	if err != nil {
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	response.SendSuccess(c, resp)
}

// @Summary     Get new M2M access token using refresh token
// @Description Use refresh token to get a new access token
// @Tags         M2M Clients
// @Accept       json
// @Produce      json
// @Param      refreshToken body RefreshTokenRequest true "Refresh token"
// @Success 200  {object} M2MTokenResponse
// @Failure 401  {object} response.CommonErrorResponse "Invalid refresh token"
// @Failure 500  {object} response.CommonErrorResponse "Internal server error"
// @Router /auth/m2m/refresh [post]
func (h *Handler) PostM2MTokenRefresh(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.RefreshM2MToken(
		c.Request.Context(),
		req.RefreshToken,
	)
	if err != nil {
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	response.SendSuccess(c, resp)
}

func (h *Handler) GetM2MClients(c *gin.Context) {
	var req ListM2MClientsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendError(c, "Invalid query parameters", http.StatusBadRequest, nil)
		return
	}

	clients, err := h.service.ListClients(c.Request.Context(), req)
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

	var req struct {
		HasPersonalInfoAccess bool `json:"hasPersonalInfoAccess"`
	}
	_ = c.ShouldBindJSON(&req)

	err := h.service.Verify(
		c.Request.Context(),
		id,
		req.HasPersonalInfoAccess,
	)
	if err != nil {
		log.Printf(
			"[PatchM2MClientVerify] {Verify M2M Client}: %v",
			err,
		)
		response.SendError(
			c,
			err.Error(),
			http.StatusInternalServerError,
			nil,
		)
		return
	}
	response.SendSuccess(c, gin.H{"message": "Verified"})
}

func (h *Handler) PatchM2MClientReject(c *gin.Context) {
	id := c.Param("id")
	err := h.service.Reject(c.Request.Context(), id)
	if err != nil {
		log.Printf(
			"[PatchM2MClientReject] {Reject M2M Client}: %v",
			err,
		)
		response.SendError(
			c,
			err.Error(),
			http.StatusInternalServerError,
			nil,
		)
		return
	}
	response.SendSuccess(c, gin.H{"message": "Rejected"})
}
