package constants

type RoleID int

const (
	StudentRoleID RoleID = iota + 1
	AdminRoleID
	SuperAdminRoleID
	DeveloperRoleID
	StudentAssistantRoleID
)
