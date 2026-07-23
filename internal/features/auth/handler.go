package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
)

type Handler struct {
	service *Service
	cfg     *config.Config
	logger  audit.Logger
}

// NewHandler creates a new authentication handler.
func NewHandler(
	service *Service,
	cfg *config.Config,
	logger audit.Logger,
) *Handler {
	return &Handler{
		service: service,
		cfg:     cfg,
		logger:  logger,
	}
}

func (h *Handler) setAuthCookies(
	c *gin.Context,
	accessToken,
	refreshToken string,
) {
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
}

// PostLogin handles traditional email/password login.
func (h *Handler) PostLogin(c *gin.Context) {
	if h.cfg.IsProduction && !h.cfg.IsStaging {
		fmt.Printf(
			"[PostLogin] {Environment Check}: " +
				"traditional login disabled in production\n",
		)
		response.SendError(
			c,
			"Traditional login is disabled in production",
			http.StatusForbidden,
			nil,
		)
		return
	}

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
	if err != nil {
		fmt.Printf("[PostLogin] {Authentication Error}: %v\n", err)
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	h.setAuthCookies(c, accessToken, refreshToken)

	response.SendSuccess(c, gin.H{
		"userId": userID,
	})
}

// GetLogout handles user logout.
func (h *Handler) GetLogout(c *gin.Context) {
	token, _ := c.Cookie(constants.AccessTokenCookieName)
	refreshToken, _ := c.Cookie(constants.RefreshTokenCookieName)
	tokenType := c.GetString("tokenType")

	var uIDStr, uEmailStr string
	if id, ok := c.Get("userID"); ok {
		uIDStr, _ = id.(string)
	}
	if email, ok := c.Get("userEmail"); ok {
		uEmailStr, _ = email.(string)
	}

	// Retrieve details from token if middleware did not run
	if token != "" {
		claims, err := tokens.NewService().ParseTokenUnverified(token)
		if err == nil {
			if uIDStr == "" {
				uIDStr = claims.UserID
			}
			if uEmailStr == "" {
				uEmailStr = claims.UserEmail
			}
			if tokenType == "" {
				tokenType = claims.TokenType
			}
		}
	}

	logoutURL, _ := h.service.Logout(
		c.Request.Context(), token, refreshToken, tokenType, h.cfg,
	)

	h.logger.Record(c.Request.Context(), nil, audit.LogEntry{
		Level:     audit.LevelInfo,
		Category:  audit.CategorySecurity,
		Action:    audit.ActionLogout,
		Message:   fmt.Sprintf("User %s logged out", uEmailStr),
		UserID:    structs.StringToNullableString(uIDStr),
		UserEmail: structs.StringToNullableString(uEmailStr),
		IPAddress: structs.StringToNullableString(c.ClientIP()),
		UserAgent: structs.StringToNullableString(c.Request.UserAgent()),
	})

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

	// Determine redirection target (fallback to frontend)
	redirectTarget := "http://localhost:5173/"
	if h.cfg.IsProduction {
		redirectTarget = "https://guisis.dllbsit2027.com/"
	} else if h.cfg.IsStaging {
		redirectTarget = "https://www.staging.guisis.dllbsit2027.com/"
	}

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
	if h.cfg.IsProduction || h.cfg.IsStaging {
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
	refreshToken, err := c.Cookie(constants.RefreshTokenCookieName)
	if err != nil {
		response.SendError(
			c,
			"Refresh token missing",
			http.StatusUnauthorized,
			nil,
		)
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	newAccessToken, newRefreshToken, err := h.service.RefreshToken(
		c.Request.Context(),
		refreshToken,
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
	h.setAuthCookies(c, newAccessToken, newRefreshToken)

	response.SendSuccess(c, gin.H{"message": "Token refreshed"})
}

// GetAuthorizeURL initiates the IDP login flow.
func (h *Handler) GetAuthorizeURL(c *gin.Context) {
	if !h.service.IsIDPUp(c.Request.Context(), h.cfg) {
		fmt.Printf(
			"[GetAuthorizeURL] {IDP Health Check}: " +
				"IDP is down, redirecting to native fallback\n",
		)
		uiBaseURL := "http://localhost:5173"

		if referer := c.GetHeader("Referer"); referer != "" {
			if u, err := url.Parse(referer); err == nil {
				if u.Scheme != "" && u.Host != "" {
					uiBaseURL = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
				}
			}
		} else {
			if h.cfg.IsStaging {
				uiBaseURL = "https://staging.guisis.dllbsit2027.com"
			}
			if h.cfg.IsProduction {
				uiBaseURL = "https://guisis.dllbsit2027.com"
			}
		}
		fallbackURL := fmt.Sprintf("%s/login?fallback=true", uiBaseURL)
		c.Redirect(http.StatusFound, fallbackURL)
		return
	}

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

	c.Redirect(http.StatusFound, authURL)
}

// PostIDPToken handles the callback/token exchange from the IDP.
func (h *Handler) PostIDPToken(c *gin.Context) {
	var req IDPTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[PostIDPToken] {Binding Error}: %v\n", err)
		response.SendFail(c, gin.H{"error": "Authorization code is required"})
		return
	}

	code := req.Code

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	accessToken, refreshToken, err := h.service.PostIDPTokenExchange(
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

	h.setAuthCookies(c, accessToken, refreshToken)

	// Return success for AJAX flow; cookies are already set.
	response.SendSuccess(c, gin.H{"message": "IDP authentication successful"})
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

// verifyFallbackAllowed checks if OTP fallback is permitted in the
// current environment based on IDP availability.
func (h *Handler) verifyFallbackAllowed(
	c *gin.Context,
	handlerName string,
) bool {
	if !h.cfg.IsProduction && !h.cfg.IsStaging {
		return true
	}

	if h.service.IsIDPUp(c.Request.Context(), h.cfg) {
		fmt.Printf(
			"[%s] {IDP Check}: "+
				"OTP fallback is only allowed when IDP is down\n",
			handlerName,
		)
		response.SendError(
			c,
			"OTP fallback is only allowed when IDP is down",
			http.StatusForbidden,
			nil,
		)
		return false
	}

	return true
}

// GetIDPStatus returns whether the IDP is up.
func (h *Handler) GetIDPStatus(c *gin.Context) {
	isUp := h.service.IsIDPUp(c.Request.Context(), h.cfg)
	response.SendSuccess(c, gin.H{
		"up": isUp,
	})
}

// PostOTPRequest triggers sending the OTP to the user's email.
func (h *Handler) PostOTPRequest(c *gin.Context) {
	if !h.verifyFallbackAllowed(c, "PostOTPRequest") {
		return
	}

	var req OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[PostOTPRequest] {Binding Error}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	err := h.service.GenerateAndSendOTP(c.Request.Context(), req.Email)
	if err != nil {
		fmt.Printf("[PostOTPRequest] {Service Error}: %v\n", err)
		response.SendError(
			c,
			err.Error(),
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "Verification code sent"})
}

// PostOTPLogin authenticates the user using email and OTP.
func (h *Handler) PostOTPLogin(c *gin.Context) {
	if !h.verifyFallbackAllowed(c, "PostOTPLogin") {
		return
	}

	var req OTPLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("[PostOTPLogin] {Binding Error}: %v\n", err)
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	userID, accessToken, refreshToken, err := h.service.AuthenticateUserOTP(
		c.Request.Context(),
		req.Email,
		req.OTP,
		ipAddress,
		userAgent,
	)
	if err != nil {
		fmt.Printf("[PostOTPLogin] {Authentication Error}: %v\n", err)
		response.SendError(c, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	h.setAuthCookies(c, accessToken, refreshToken)

	response.SendSuccess(c, gin.H{"userId": userID})
}

func (h *Handler) IsIDPUp(ctx context.Context) bool {
	return h.service.IsIDPUp(ctx, h.cfg)
}
