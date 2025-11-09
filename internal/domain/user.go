package domain

import "time"

type User struct {
	MaxID       string     `json:"max_id" db:"max_id"`
	Name        string     `json:"name" db:"name"`
	Geolocation string     `json:"geolocation" db:"geolocation"`
	Age         int        `json:"age" db:"age"`
	Sex         Sex        `json:"sex" db:"sex"`
	About       string     `json:"about" db:"about"`
	Role        UserRole   `json:"role" db:"role"`
	Status      UserStatus `json:"status" db:"status"`

	ReputationGroupID int              `json:"reputation_group_id" db:"reputation_group_id" default:"1"`
	ReputationGroup   *ReputationGroup `json:"reputation_group" db:"-"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)
