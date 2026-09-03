package m2mclients

import (
	"time"
)

type M2MClient struct {
	ID                    int        `db:"id"                 json:"id"`
	UserID                string     `db:"user_id"            json:"userId"`
	ClientName            string     `db:"client_name"        json:"clientName"`
	ClientID              string     `db:"client_id"          json:"clientId"`
	ClientSecret          string     `db:"client_secret_hash" json:"-"`
	ClientDescription     string     `db:"client_description" json:"clientDescription"`
	IsActive              bool       `db:"is_active"          json:"isActive"`
	IsVerified            bool       `db:"is_verified"        json:"isVerified"`
	HasPersonalInfoAccess bool       `db:"has_personal_info_access" json:"hasPersonalInfoAccess"`
	LastUsedAt            *time.Time `db:"last_used_at"       json:"lastUsedAt"`
	CreatedAt             time.Time  `db:"created_at"         json:"createdAt"`
	UpdatedAt             time.Time  `db:"updated_at"         json:"updatedAt"`
}
