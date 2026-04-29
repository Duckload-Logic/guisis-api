package users

import "github.com/olazo-johnalbert/duckload-api/internal/core/structs"

type UserResponse struct {
	ID             string                 `json:"id"`
	Roles          []Role                 `json:"roles"`
	FirstName      string                 `json:"firstName"`
	MiddleName     structs.NullableString `json:"middleName,omitempty"`
	LastName       string                 `json:"lastName"`
	SuffixName     structs.NullableString `json:"suffixName,omitempty"`
	Email          string                 `json:"email,omitempty"`
	IsActive       bool                   `json:"isActive"`
	ProfilePicture string                 `json:"profilePicture,omitempty"`
	CreatedAt      string                 `json:"createdAt,omitempty"`
	UpdatedAt      string                 `json:"updatedAt,omitempty"`
}

type CreateUserRequest struct {
	RoleID     int    `json:"roleId"     binding:"required"`
	FirstName  string `json:"firstName"  binding:"required"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"   binding:"required"`
	SuffixName string `json:"suffixName"`
	Email      string `json:"email"      binding:"required,email"`
	Password   string `json:"password"   binding:"required,min=8"`
}

type ListUsersRequest struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	RoleID   int    `form:"role_id"`
	Search   string `form:"search"`
	Active   *bool  `form:"active"`
}

type ListUsersResponse struct {
	Users []UserResponse             `json:"users"`
	Meta  structs.PaginationMetadata `json:"meta"`
}

type RoleDistributionDTO struct {
	RoleName string `json:"roleName" db:"role_name"`
	Count    int    `json:"count"    db:"count"`
}

type UpdateRolesRequest struct {
	UserID      string `json:"userId"      binding:"required"`
	RoleIDs     []int  `json:"roleIds"     binding:"required,min=1"`
	Reason      string `json:"reason"      binding:"required"`
	ReferenceID string `json:"referenceId"`
}

type AddUserToWhitelistRequest struct {
	Email   string `json:"email" binding:"required,email"`
	RoleIDs []int  `json:"roleIds" binding:"required,min=1"`
}

type RemoveUserFromWhitelistRequest struct {
	Email string `json:"email" binding:"required,email"`
}
