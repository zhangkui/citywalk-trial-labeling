package landmark

import "time"

type Status int8

const (
	Pending Status = iota
	Published
)

type Landmark struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Lat         float64    `json:"lat"`
	Lng         float64    `json:"lng"`
	CategoryID  int64      `json:"categoryId"`
	Description string     `json:"description"`
	CoverImage  string     `json:"coverImage"`
	CreatedBy   int64      `json:"createdBy"`
	Status      Status     `json:"status"`
	Rating      float64    `json:"rating"`
	RatingCount int        `json:"ratingCount"`
	ViewCount   int        `json:"viewCount"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
