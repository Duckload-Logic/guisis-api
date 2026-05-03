package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
)

type Handler struct {
	service *Service
	cfg     *config.Config
}

// NewHandler creates a new authentication handler.
func NewHandler(service *Service, cfg *config.Config) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
	}
}

// PostLogin handles traditional email/password login.
func (h *Handler) PostLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[PostLogin] {Binding Error}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	userID, accessToken, refreshToken, err := h.service.AuthenticateUser(
		c.Request.Context(),
		req.Email,
		req.Password,
		ipAddress,
		userAgent,
	)
	_ = refreshToken // explicitly ignore if unused for now
	if err != nil {
		fmt.Printf("[PostLogin] {Authentication Error}: %v\n", err)
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	// Set secure httpOnly cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		constants.AccessTokenCookieName,
		accessToken,
		int(constants.RefreshTokenMaxAge),
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)
	c.SetCookie(
		constants.RefreshTokenCookieName,
		refreshToken,
		int(constants.RefreshTokenMaxAge),
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)

	response.SendSuccess(c, gin.H{
		"userId": userID,
	})
}

// PostRegister initiates the user registration process.
func (h *Handler) PostRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[PostRegister] {Binding Error}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	registrationID, err := h.service.RegisterUser(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[PostRegister] {Service Error}: %v\n", err)
		response.SendError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	response.SendSuccess(c, gin.H{"registrationID": registrationID})
}

// PostResendVerification handles requests to resend the verification email.
func (h *Handler) PostResendVerification(c *gin.Context) {
	registrationID := c.Param("registrationID")
	if registrationID == "" {
		response.SendFail(c, gin.H{"error": "registrationID is required"})
		return
	}

	err := h.service.ResendVerification(c.Request.Context(), registrationID)
	if err != nil {
		fmt.Printf("[PostResendVerification] {Service Error}: %v\n", err)
		response.SendError(c, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	response.SendSuccess(c, gin.H{"message": "Verification email resent"})
}

// PostVerify handles email verification and finalizes registration.
func (h *Handler) PostVerify(c *gin.Context) {
	registrationID := c.Param("registrationID")
	var req VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	userID, email, err := h.service.VerifyUser(
		c.Request.Context(),
		registrationID,
		req.VerificationOTP,
	)
	if err != nil {
		fmt.Printf("[PostVerify] {Service Error}: %v\n", err)
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	response.SendSuccess(c, gin.H{
		"userId": userID,
		"email":  email,
	})
}

// GetLogout handles user logout.
func (h *Handler) GetLogout(c *gin.Context) {
	token, _ := c.Cookie(constants.AccessTokenCookieName)
	tokenType := c.GetString("tokenType")

	logoutURL, _ := h.service.Logout(
		c.Request.Context(), token, tokenType, h.cfg,
	)

	// Clear cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		constants.AccessTokenCookieName,
		"", -1, constants.CookiePathRoot, "",
		h.cfg.IsProduction, true,
	)
	c.SetCookie(
		constants.RefreshTokenCookieName,
		"", -1, constants.CookiePathRoot, "",
		h.cfg.IsProduction, true,
	)

	// Determine redirection target
	redirectTarget := "/"
	if logoutURL != "" {
		redirectTarget = logoutURL
	}

	// Handle dynamic redirect for native/local sessions
	if tokenType != string(constants.AuthTypeIDP) {
		candidate := c.Query("redirect_uri")
		if candidate != "" {
			parsedURL, err := url.Parse(candidate)
			if err == nil {
				origin := fmt.Sprintf(
					"%s://%s",
					parsedURL.Scheme,
					parsedURL.Host,
				)
				if h.isAllowedOrigin(origin) {
					redirectTarget = origin + "/"
				}
			}
		}
	}

	// Security: Prevent caching of logout state
	c.Header(
		"Cache-Control",
		"no-store, no-cache, must-revalidate",
	)
	c.Redirect(http.StatusFound, redirectTarget)
}

// isAllowedOrigin checks if the given origin is
// permitted for redirects.
func (h *Handler) isAllowedOrigin(origin string) bool {
	if h.cfg.IsProduction {
		target := "dllbsit2027.com"
		parsed, err := url.Parse(origin)
		if err != nil {
			return false
		}
		host := parsed.Hostname()
		return host == target ||
			strings.HasSuffix(host, "."+target)
	}

	return strings.HasPrefix(origin, "http://localhost")
}

// PostRefreshToken handles session refreshing using the refresh token.
func (h *Handler) PostRefreshToken(c *gin.Context) {
	if _, err := c.Cookie(constants.RefreshTokenCookieName); err != nil {
		response.SendError(c, "Refresh token missing", http.StatusUnauthorized, nil)
		return
	}

	accessTokenJTI := c.MustGet("accessTokenJTI").(sessions.JTIDTO)
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	newAccessToken, newRefreshToken, err := h.service.RefreshToken(
		c.Request.Context(),
		accessTokenJTI,
		h.cfg,
		ipAddress,
		userAgent,
	)
	if err != nil {
		fmt.Printf("[PostRefreshToken] {Service Error}: %v\n", err)
		response.SendError(c, "Session expired", http.StatusUnauthorized, nil)
		return
	}

	// Set new cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		constants.AccessTokenCookieName,
		newAccessToken,
		int(constants.RefreshTokenMaxAge),
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)
	c.SetCookie(
		constants.RefreshTokenCookieName,
		newRefreshToken,
		int(constants.RefreshTokenMaxAge),
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)

	response.SendSuccess(c, gin.H{"message": "Token refreshed"})
}

// GetAuthorizeURL initiates the IDP login flow.
func (h *Handler) GetAuthorizeURL(c *gin.Context) {
	authURL, err := h.service.GetAuthorizeURL(h.cfg)
	if err != nil {
		fmt.Printf("[GetAuthorizeURL] {Service Error}: %v\n", err)
		response.SendError(
			c,
			"Failed to initiate IDP login",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"url": authURL})
}

// PostIDPToken handles the callback/token exchange from the IDP.
func (h *Handler) PostIDPToken(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.SendFail(c, gin.H{"error": "Authorization code is missing"})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	accessToken, refreshToken,
		userID, email,
		role, err := h.service.PostIDPTokenExchange(
		c.Request.Context(),
		code,
		h.cfg,
		ipAddress,
		userAgent,
	)
	if err != nil {
		fmt.Printf("[PostIDPToken] {Service Error}: %v\n", err)
		response.SendError(
			c,
			"Failed to exchange code",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	// Set cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		constants.AccessTokenCookieName,
		accessToken,
		int(constants.RefreshTokenMaxAge),
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)
	c.SetCookie(
		constants.RefreshTokenCookieName,
		refreshToken,
		int(constants.RefreshTokenMaxAge),
		constants.CookiePathRoot,
		"",
		h.cfg.IsProduction,
		true,
	)

	// Redirect to frontend with basic info
	frontendURL := fmt.Sprintf(
		"%s/auth/callback?userID=%s&email=%s&role=%s",
		h.cfg.AppFrontendUrl,
		userID,
		email,
		role,
	)
	c.Redirect(http.StatusFound, frontendURL)
}

// GetMe returns current user information.
func (h *Handler) GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	tokenType := c.MustGet("tokenType").(string)

	user, err := h.service.GetMe(c.Request.Context(), userID, tokenType)
	if err != nil {
		fmt.Printf("[GetMe] {Service Error}: %v\n", err)
		response.SendError(
			c,
			"Failed to get profile",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, user)
}

func (h *Handler) PostBlockUser(c *gin.Context) {
	userID := c.Param("id")
	err := h.service.BlockUser(c.Request.Context(), userID)
	if err != nil {
		response.SendError(
			c,
			"Failed to block user",
			http.StatusInternalServerError,
			nil,
		)
		return
	}
	response.SendSuccess(c, gin.H{"message": "User blocked"})
}

func (h *Handler) PostUnblockUser(c *gin.Context) {
	userID := c.Param("id")
	err := h.service.UnblockUser(c.Request.Context(), userID)
	if err != nil {
		response.SendError(
			c,
			"Failed to unblock user",
			http.StatusInternalServerError,
			nil,
		)
		return
	}
	response.SendSuccess(c, gin.H{"message": "User unblocked"})
}
