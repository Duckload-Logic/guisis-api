package users

import (
	"database/sql"

	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
)

type Role struct {
	ID   int    `db:"id"   json:"id"`
	Name string `db:"name" json:"name"`
}

type RoleAssignment struct {
	UserID      string                 `db:"user_id"`
	RoleID      int                    `db:"role_id"`
	AssignedAt  structs.NullableTime   `db:"assigned_at"`
	AssignedBy  structs.NullableString `db:"assigned_by"`
	Reason      structs.NullableString `db:"reason"`
	ReferenceID structs.NullableString `db:"reference_id"`
}

type ProfilePicture struct {
	UserID string `db:"user_id"`
	FileID string `db:"file_id"`
}

type User struct {
	ID             string                 `db:"id"            json:"id"`
	FirstName      string                 `db:"first_name"    json:"firstName"`
	MiddleName     structs.NullableString `db:"middle_name"   json:"middleName"`
	LastName       string                 `db:"last_name"     json:"lastName"`
	SuffixName     structs.NullableString `db:"suffix_name"   json:"suffixName"`
	Email          string                 `db:"email"         json:"email"`
	PasswordHash   structs.NullableString `db:"password_hash" json:"-"`
	AuthType       string                 `db:"auth_type"     json:"authType"`
	IsActive       bool                   `db:"is_active"     json:"isActive"`
	IDPUUID        structs.NullableString `db:"idp_uuid"      json:"idpUuid"`
	CreatedAt      structs.NullableTime   `db:"created_at"    json:"createdAt"`
	UpdatedAt      structs.NullableTime   `db:"updated_at"    json:"updatedAt"`
	ProfilePicture structs.NullableString `db:"-"             json:"profilePicture"`
	Roles          []Role                 `db:"-"             json:"roles"`
}

// UserRoleDB is still needed for mapping join table results if necessary,
// but we can often avoid it.
type UserRoleDB struct {
	UserID      string         `db:"user_id"`
	RoleID      int            `db:"role_id"`
	AssignedAt  sql.NullTime   `db:"assigned_at"`
	AssignedBy  sql.NullString `db:"assigned_by"`
	Reason      sql.NullString `db:"reason"`
	ReferenceID sql.NullString `db:"reference_id"`
}

// FullName returns the user's full name, fallback to a default if user is nil.
func (u *User) FullName() string {
	if u == nil {
		return "A student"
	}
	return u.FirstName + " " + u.LastName
}

