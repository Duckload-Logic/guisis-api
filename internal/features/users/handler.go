package users

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/response"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/features/files"
)

type Handler struct {
	service        *Service
	sessionService *sessions.Service
	logger         audit.Logger
	logReader      audit.LogReader
}

// NewHandler creates a new users handler.
func NewHandler(
	service *Service,
	sessionService *sessions.Service,
	logger audit.Logger,
	logReader audit.LogReader,
) *Handler {
	return &Handler{
		service:        service,
		sessionService: sessionService,
		logger:         logger,
		logReader:      logReader,
	}
}

// GetMe retrieves the currently authenticated user's information.
func (h *Handler) GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	resp, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		fmt.Printf("[GetMe] {GetUserByID}: %v\n", err)
		response.SendError(
			c,
			"Failed to get current user",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetUserByEmail retrieves user information by their email address.
func (h *Handler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		fmt.Printf(
			"[GetUserByEmail] {Check Query}: Email parameter is required\n",
		)
		response.SendFail(
			c,
			gin.H{"error": "Email query parameter is required"},
		)
		return
	}

	authType := c.DefaultQuery("auth_type", "native")

	resp, err := h.service.GetUserByEmail(c.Request.Context(), email, authType)
	if err != nil {
		fmt.Printf("[GetUserByEmail] {GetUserByEmail}: %v\n", err)
		response.SendError(
			c,
			"Failed to get user by email",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetUserByID retrieves user information by their ID.
func (h *Handler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		fmt.Printf(
			"[GetUserByID] {Check Param}: User ID is required\n",
		)
		response.SendFail(
			c,
			gin.H{"error": "User ID parameter is required"},
		)
		return
	}

	resp, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		fmt.Printf("[GetUserByID] {GetUserByID}: %v\n", err)
		response.SendError(
			c,
			"Failed to get user by ID",
			http.StatusNotFound,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetUsers retrieves a paginated list of all users.
func (h *Handler) GetUsers(c *gin.Context) {
	var params ListUsersRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.ListUsers(c.Request.Context(), params)
	if err != nil {
		fmt.Printf("[GetUsers] {ListUsers}: %v\n", err)
		response.SendError(
			c,
			"Failed to list users",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

// GetRoleDistribution returns the count of users for each role.
func (h *Handler) GetRoleDistribution(c *gin.Context) {
	resp, err := h.service.GetRoleDistribution(c.Request.Context())
	if err != nil {
		fmt.Printf("[GetRoleDistribution] {GetRoleDistribution}: %v\n", err)
		response.SendError(
			c,
			"Failed to get role distribution",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, resp)
}

func (h *Handler) PostUserBlock(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.SendFail(c, gin.H{"error": "User ID parameter is required"})
		return
	}

	err := h.service.BlockUser(c.Request.Context(), userID)
	if err != nil {
		fmt.Printf("[PostUserBlock] {BlockUser}: %v\n", err)
		response.SendError(
			c,
			"Failed to block user",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "User blocked successfully"})
}

func (h *Handler) PostUserUnblock(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.SendFail(c, gin.H{"error": "User ID parameter is required"})
		return
	}

	err := h.service.UnblockUser(c.Request.Context(), userID)
	if err != nil {
		fmt.Printf("[PostUserUnblock] {UnblockUser}: %v\n", err)
		response.SendError(
			c,
			"Failed to unblock user",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, gin.H{"message": "User unblocked successfully"})
}

func (h *Handler) GetUserSessions(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		response.SendFail(c, gin.H{"error": "User ID is required"})
		return
	}

	sessions, err := h.sessionService.ListUserSessions(
		c.Request.Context(),
		targetUserID,
	)
	if err != nil {
		fmt.Printf("[GetUserSessions] {ListSessions}: %v\n", err)
		response.SendError(
			c,
			"Failed to list user sessions",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, sessions)
}

func (h *Handler) DeleteUserSession(c *gin.Context) {
	targetUserID := c.Param("id")
	jti := c.Param("session_id")
	if targetUserID == "" || jti == "" {
		response.SendFail(
			c,
			gin.H{"error": "User ID and Session ID are required"},
		)
		return
	}

	err := h.sessionService.DeleteUserToken(
		c.Request.Context(),
		targetUserID,
		sessions.NewJTI(jti),
	)
	if err != nil {
		fmt.Printf("[DeleteUserSession] {DeleteToken}: %v\n", err)
		response.SendError(
			c,
			"Failed to revoke session",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	// Audit log for session revocation
	adminEmail := c.MustGet("userEmail").(string)
	adminID := c.MustGet("userID").(string)

	user, err := h.service.GetUserByID(c.Request.Context(), targetUserID)
	if err != nil {
		fmt.Printf("[DeleteUserSession] {GetUserByID}: %v\n", err)
		response.SendError(c, "Failed to find user", http.StatusNotFound, nil)
		return
	}

	h.logger.Record(
		c.Request.Context(),
		nil,
		audit.LogEntry{
			Level:    audit.LevelWarning,
			Category: audit.CategorySecurity,
			Action:   audit.ActionLogout,
			Message: fmt.Sprintf(
				"Superadmin %s revoked session %s for user %s",
				adminEmail,
				jti,
				targetUserID,
			),
			UserID:      structs.StringToNullableString(adminID),
			UserEmail:   structs.StringToNullableString(adminEmail),
			TargetID:    structs.StringToNullableString(targetUserID),
			TargetEmail: structs.StringToNullableString(user.Email),
		},
	)

	response.SendSuccess(c, gin.H{"message": "Session revoked successfully"})
}

func (h *Handler) GetUserActivity(c *gin.Context) {
	targetUserID := c.Param("id")
	if targetUserID == "" {
		response.SendFail(c, gin.H{"error": "User ID is required"})
		return
	}

	// Fetch user email first
	user, err := h.service.GetUserByID(c.Request.Context(), targetUserID)
	if err != nil {
		fmt.Printf("[GetUserActivity] {GetUserByID}: %v\n", err)
		response.SendError(c, "Failed to find user", http.StatusNotFound, nil)
		return
	}

	var req audit.ListSystemLogsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	// Filter by target user's email
	req.UserEmail = user.Email

	result, err := h.logReader.ListLogs(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetUserActivity] {ListLogs}: %v\n", err)
		response.SendError(
			c,
			"Failed to retrieve user activity",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, result)
}

func (h *Handler) PostUpdateRoles(c *gin.Context) {
	var req UpdateRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	adminID := c.MustGet("userID").(string)

	err := h.service.UpdateUserRoles(c.Request.Context(), req, adminID)
	if err != nil {
		fmt.Printf("[PostUpdateRoles] {UpdateUserRoles}: %v\n", err)
		if err.Error() == "cannot modify your own roles" ||
			err.Error() == "user cannot be both Counselor and SuperAdmin" ||
			err.Error() == "user cannot be both Student and Counselor" ||
			err.Error() == "user cannot be both Counselor and Developer" {
			response.SendFail(c, gin.H{"error": err.Error()})
			return
		}
		response.SendError(
			c,
			"Failed to update user roles",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	// Audit log
	adminEmail := c.MustGet("userEmail").(string)
	h.logger.Record(
		c.Request.Context(),
		nil,
		audit.LogEntry{
			Level:    audit.LevelWarning,
			Category: audit.CategorySecurity,
			Action:   audit.ActionElevateRoles,
			Message: fmt.Sprintf(
				"Superadmin %s elevated/updated roles for user %s. Reason: %s, Ref: %s",
				adminEmail,
				req.UserID,
				req.Reason,
				req.ReferenceID,
			),
			UserID:   structs.StringToNullableString(adminID),
			TargetID: structs.StringToNullableString(req.UserID),
		},
	)

	response.SendSuccess(c, gin.H{"message": "User roles updated successfully"})
}

func (h *Handler) PostUserToWhitelist(c *gin.Context) {
	var req AddUserToWhitelistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}
	err := h.service.AddUserToWhitelist(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[PostUserToWhitelist] {AddUserToWhitelist}: %v\n", err)
		if err.Error() == "user cannot be whitelisted as both Counselor and SuperAdmin" ||
			err.Error() == "user cannot be whitelisted as both Student and Counselor" ||
			err.Error() == "user cannot be whitelisted as both Counselor and Developer" {
			response.SendFail(c, gin.H{"error": err.Error()})
			return
		}
		response.SendError(
			c,
			"Failed to add user to whitelist",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "User added to whitelist successfully"},
	)
}

func (h *Handler) PostRemoveUserFromWhitelist(c *gin.Context) {
	var req RemoveUserFromWhitelistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	err := h.service.RemoveUserFromWhitelist(c.Request.Context(), req)
	if err != nil {
		fmt.Printf(
			"[PostRemoveUserFromWhitelist] {RemoveUserFromWhitelist}: %v\n",
			err,
		)
		response.SendError(
			c,
			"Failed to remove user from whitelist",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "User removed from whitelist successfully"},
	)
}

func (h *Handler) GetWhitelist(c *gin.Context) {
	var req ListWhitelistRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.SendFail(c, gin.H{"error": err.Error()})
		return
	}

	entries, err := h.service.ListWhitelist(c.Request.Context(), req)
	if err != nil {
		fmt.Printf("[GetWhitelist] {ListWhitelist}: %v\n", err)
		response.SendError(
			c,
			"Failed to fetch whitelist",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(c, entries)
}

func (h *Handler) PostProfilePicture(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	fileID := c.Param("id")
	if fileID == "" {
		response.SendFail(c, gin.H{"error": "File ID is required"})
		return
	}

	err := h.service.PostProfilePicture(c.Request.Context(), userID, fileID)
	if err != nil {
		fmt.Printf("[PostProfilePicture] {PostProfilePicture}: %v\n", err)
		response.SendError(
			c,
			"Failed to associate profile picture",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{"message": "Profile picture updated successfully"},
	)
}

func (h *Handler) UploadProfilePicture(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Printf(
			"[UploadProfilePicture] {FormFile}: %v\n",
			err,
		)
		response.SendFail(c, gin.H{"error": "File is required"})
		return
	}

	fileURL, err := h.service.UploadProfilePicture(
		c.Request.Context(),
		userID,
		fileHeader,
	)
	if err != nil {
		fmt.Printf(
			"[UploadProfilePicture] {UploadProfilePicture}: %v\n",
			err,
		)
		if errors.Is(err, files.ErrFileTooLarge) {
			response.SendFail(
				c,
				gin.H{"error": "Profile picture must not exceed 5MB"},
			)
			return
		}
		response.SendError(
			c,
			"Failed to upload profile picture",
			http.StatusInternalServerError,
			nil,
		)
		return
	}

	response.SendSuccess(
		c,
		gin.H{
			"message": "Profile picture uploaded successfully",
			"url":     fileURL,
		},
	)
}
