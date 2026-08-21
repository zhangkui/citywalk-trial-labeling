package dto

type CommentRequest struct {
	TargetType string `json:"targetType"`
	TargetID   int64  `json:"targetId"`
	Content    string `json:"content"`
	ParentID   int64  `json:"parentId"`
}

type CommentUpdateRequest struct {
	Content string `json:"content"`
}

