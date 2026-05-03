package audit

import (
	"encoding/json"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

// Log levels
const (
	LevelInfo     = "INFO"
	LevelWarning  = "WARNING"
	LevelError    = "ERROR"
	LevelCritical = "CRITICAL"
)

// Log categories
const (
	CategoryAudit    = "AUDIT"
	CategorySystem   = "SYSTEM"
	CategorySecurity = "SECURITY"
	CategoryConsent  = "CONSENT"
)

// Audit log actions — track data changes and admin operations
const (
	ActionUserCreated       = "USER_CREATED"
	ActionUserCreateFailed  = "USER_CREATE_FAILED"
	ActionUserUpdated       = "USER_UPDATED"
	ActionUserUpdateFailed  = "USER_UPDATE_FAILED"
	ActionUserDeleted       = "USER_DELETED"
	ActionUserDeleteFailed  = "USER_DELETE_FAILED"
	ActionUserBlocked       = "USER_BLOCKED"
	ActionUserBlockFailed   = "USER_BLOCK_FAILED"
	ActionUserUnblocked     = "USER_UNBLOCKED"
	ActionUserUnblockFailed = "USER_UNBLOCK_FAILED"

	ActionRoleChanged      = "ROLE_CHANGED"
	ActionRoleChangeFailed = "ROLE_CHANGE_FAILED"

	ActionAppointmentCreated      = "APPOINTMENT_CREATED"
	ActionAppointmentCreateFailed = "APPOINTMENT_CREATE_FAILED"
	ActionAppointmentUpdated      = "APPOINTMENT_UPDATED"
	ActionAppointmentUpdateFailed = "APPOINTMENT_UPDATE_FAILED"
	ActionAppointmentDeleted      = "APPOINTMENT_DELETED"
	ActionAppointmentDeleteFailed = "APPOINTMENT_DELETE_FAILED"
	ActionAppointmentFailed       = "APPOINTMENT_FAILED"

	ActionSlipCreated       = "SLIP_CREATED"
	ActionSlipCreateFailed  = "SLIP_CREATE_FAILED"
	ActionSlipStatusUpdated = "SLIP_STATUS_UPDATED"
	ActionSlipUpdated       = "SLIP_UPDATED"
	ActionSlipUpdateFailed  = "SLIP_UPDATE_FAILED"
	ActionSlipDeleted       = "SLIP_DELETED"
	ActionSlipDeleteFailed  = "SLIP_DELETE_FAILED"
	ActionSlipFailed        = "SLIP_FAILED"

	ActionNoteCreated      = "NOTE_CREATED"
	ActionNoteCreateFailed = "NOTE_CREATE_FAILED"
	ActionNoteUpdated      = "NOTE_UPDATED"
	ActionNoteUpdateFailed = "NOTE_UPDATE_FAILED"
	ActionNoteDeleted      = "NOTE_DELETED"
	ActionNoteDeleteFailed = "NOTE_DELETE_FAILED"

	ActionIIRCreated      = "IIR_CREATED"
	ActionIIRCreateFailed = "IIR_CREATE_FAILED"
	ActionIIRUpdated      = "IIR_UPDATED"
	ActionIIRUpdateFailed = "IIR_UPDATE_FAILED"
	ActionIIRDeleted      = "IIR_DELETED"
	ActionIIRDeleteFailed = "IIR_DELETE_FAILED"
	ActionIIRDraftSaved   = "IIR_DRAFT_SAVED"
	ActionIIRSubmitted    = "IIR_SUBMITTED"
)

// System log actions — track system-level events
const (
	ActionM2MClientCreated            = "M2M_CLIENT_CREATED"
	ActionM2MClientRevoked            = "M2M_CLIENT_REVOKED"
	ActionM2MClientVerified           = "M2M_CLIENT_VERIFIED"
	ActionM2MClientSecretRotated      = "M2M_CLIENT_SECRET_ROTATED" // nolint:gosec
	ActionSettingChanged              = "SETTING_CHANGED"
	ActionM2MClientCreateFailed       = "M2M_CLIENT_CREATE_FAILED"
	ActionM2MClientRevokeFailed       = "M2M_CLIENT_REVOKE_FAILED"
	ActionM2MClientVerifyFailed       = "M2M_CLIENT_VERIFY_FAILED"
	ActionM2MClientSecretRotateFailed = "M2M_CLIENT_SECRET_ROTATE_FAILED" // nolint:gosec
	ActionSettingChangeFailed         = "SETTING_CHANGE_FAILED"
)

// Security log actions — track authentication and access events
const (
	ActionLoginSuccess       = "LOGIN_SUCCESS"
	ActionLoginFailed        = "LOGIN_FAILED"
	ActionLogout             = "LOGOUT"
	ActionTokenRefreshed     = "TOKEN_REFRESHED"
	ActionAccessDenied       = "ACCESS_DENIED"
	ActionRateLimitExceeded  = "RATE_LIMIT_EXCEEDED"
	ActionInvalidToken       = "INVALID_TOKEN"
	ActionM2MClientUsed      = "M2M_CLIENT_USED"
	ActionM2MClientInvalid   = "M2M_CLIENT_INVALID"
	ActionM2MAuthSuccess     = "M2M_AUTH_SUCCESS"
	ActionM2MAuthFailed      = "M2M_AUTH_FAILED"
	ActionM2MTokenRefreshed  = "M2M_TOKEN_REFRESHED" // nolint:gosec
	ActionElevateRoles       = "ELEVATE_ROLES"
	ActionElevateRolesFailed = "ELEVATE_ROLES_FAILED"
)

// LogEntry is the input struct used by other services to record a log.
type LogEntry struct {
	Level    string
	Category string
	Action   string
	Message  string

	UserID   structs.NullableString // "" means no user (system event)
	TargetID structs.NullableString // Optional ID of the affected resource

	// TargetType can be used to specify the type of the target resource (e.g.,
	// "User", "Appointment")
	TargetType structs.NullableString

	UserEmail structs.NullableString // "" means no user (system event)
	// Optional email of the affected resource
	TargetEmail structs.NullableString

	IPAddress structs.NullableString
	UserAgent structs.NullableString
	Metadata  *LogMetadata
	TraceID   structs.NullableString
}

// NotificationEntry is the input struct used by other services to send a
// notification.
type NotificationEntry struct {
	ID         string                 `json:"id"`
	ReceiverID structs.NullableString `json:"receiverId,omitempty"`
	ActorID    structs.NullableString `json:"actorId,omitempty"`
	TargetID   structs.NullableString `json:"targetId,omitempty"`
	TargetType structs.NullableString `json:"targetType,omitempty"`
	Title      string                 `json:"title"`
	Message    string                 `json:"message"`
	Type       string                 `json:"type"`
	IsRead     bool                   `json:"isRead"`
	CreatedAt  time.Time              `json:"createdAt"`
}

// LogMetadata defines a structured format for audit log metadata.
type LogMetadata struct {
	EntityType string      `json:"entityType,omitempty"`
	EntityID   string      `json:"entityId,omitempty"`
	OldValues  interface{} `json:"oldValues,omitempty"`
	NewValues  interface{} `json:"newValues,omitempty"`
	Error      string      `json:"error,omitempty"`
}

// ListSystemLogsRequest holds query parameters for listing system logs
type ListSystemLogsRequest struct {
	structs.PaginationRequest
	Level       string `form:"level,omitempty"        binding:"omitempty,oneof=INFO WARNING ERROR CRITICAL"`
	Category    string `form:"category,omitempty"     binding:"omitempty,oneof=AUDIT SYSTEM SECURITY"`
	Action      string `form:"action,omitempty"`
	UserEmail   string `form:"user_email,omitempty"`
	TargetType  string `form:"target_type,omitempty"`
	TargetEmail string `form:"target_email,omitempty"`
	StartDate   string `form:"start_date,omitempty"`
	EndDate     string `form:"end_date,omitempty"`
}

// ListSystemLogsDTO is the paginated response for system logs
type ListSystemLogsDTO struct {
	Logs []SystemLogDTO             `json:"logs"`
	Meta structs.PaginationMetadata `json:"meta"`
}

// SystemLogDTO is the response DTO for a single system log entry
type SystemLogDTO struct {
	ID       int    `json:"id"`
	Level    string `json:"level"`
	Category string `json:"category"`
	Action   string `json:"action"`
	Message  string `json:"message"`

	UserID     structs.NullableString `json:"userId,omitempty"`
	TargetID   structs.NullableString `json:"targetId,omitempty"`
	TargetType structs.NullableString `json:"targetType,omitempty"`

	UserEmail   structs.NullableString `json:"userEmail,omitempty"`
	TargetEmail structs.NullableString `json:"targetEmail,omitempty"`

	IPAddress structs.NullableString `json:"ipAddress,omitempty"`
	UserAgent structs.NullableString `json:"userAgent,omitempty"`
	Metadata  json.RawMessage        `json:"metadata,omitempty"`
	TraceID   structs.NullableString `json:"traceId,omitempty"`
	CreatedAt time.Time              `json:"createdAt"`
}

// LogStatsDTO returns summary counts by category
type LogStatsDTO struct {
	Category string `json:"category" db:"category"`
	Count    int    `json:"count"    db:"count"`
}

// LogActivityDTO holds activity stats grouped by hour
type LogActivityDTO struct {
	Time     string `json:"time"     db:"time"`
	Requests int    `json:"requests" db:"requests"`
	Errors   int    `json:"errors"   db:"errors"`
}
