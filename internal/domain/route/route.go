package route

import "time"

type Status int8

const (
	StatusDraft Status = iota
	StatusPending
	StatusPublished
	StatusOffline
)

type Route struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ThemeID        int64     `json:"themeId"`
	City           string    `json:"city"`
	StartLat       float64   `json:"startLat"`
	StartLng       float64   `json:"startLng"`
	EndLat         float64   `json:"endLat"`
	EndLng         float64   `json:"endLng"`
	TotalDistance  int       `json:"totalDistance"`
	Duration       int       `json:"duration"`
	Difficulty     int8      `json:"difficulty"`
	CoverImage     string    `json:"coverImage"`
	CreatedBy      int64     `json:"createdBy"`
	Status         Status    `json:"status"`
	FavoriteCount  int       `json:"favoriteCount"`
	Rating         float64   `json:"rating"`
	RatingCount    int       `json:"ratingCount"`
	ViewCount      int       `json:"viewCount"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty"`
	Waypoints      []Waypoint `json:"waypoints,omitempty"`
}
