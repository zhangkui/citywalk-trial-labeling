package story

import "time"

type Status int8

const (
	Draft Status = iota
	Pending
	Published
)

type Story struct {
	ID         int64      `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	CoverImage string     `json:"coverImage"`
	LandmarkIDs string     `json:"landmarkIds"`
	RouteID    int64      `json:"routeId"`
	CreatedBy  int64      `json:"createdBy"`
	Status     Status     `json:"status"`
	LikeCount  int        `json:"likeCount"`
	ViewCount  int        `json:"viewCount"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	Media      []Media    `json:"media,omitempty"`
}
