package user

import "time"

type Status int8

const (
	Enabled Status = 1
	Disabled Status = 0
)

type Explorer struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar"`
	Bio          string    `json:"bio"`
	City         string    `json:"city"`
	Status       Status    `json:"status"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Roles        []string  `json:"roles,omitempty"`
	Permissions  []string  `json:"permissions,omitempty"`
}
