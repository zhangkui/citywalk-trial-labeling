package social

import "time"

type Comment struct {
	ID         int64      `json:"id"`
	TargetType string     `json:"targetType"`
	TargetID   int64      `json:"targetId"`
	Content    string     `json:"content"`
	ParentID   int64      `json:"parentId"`
	CreatedBy  int64      `json:"createdBy"`
	LikeCount  int        `json:"likeCount"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}
