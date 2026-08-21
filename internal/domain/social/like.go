package social

import "time"

type Like struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"userId"`
	TargetType string    `json:"targetType"`
	TargetID   int64     `json:"targetId"`
	CreatedAt  time.Time `json:"createdAt"`
}
