package service

import "time"

type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (p Page) limitOffset() (int, int) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 || p.PageSize > 100 {
		p.PageSize = 20
	}
	return p.PageSize, (p.Page - 1) * p.PageSize
}

type AuthResult struct {
	User         map[string]any `json:"user"`
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	ExpiresIn    int64          `json:"expiresIn"`
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type AuditInfo struct {
	UserID    *int64
	Action    string
	Resource  string
	ResourceID *int64
	Detail    string
	IP        string
	UA        string
	Result    bool
}

type CurrentUser struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	Bio         string    `json:"bio"`
	City        string    `json:"city"`
	Status      int8      `json:"status"`
	Roles       []string  `json:"roles"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

