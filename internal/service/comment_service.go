package service

import (
	"context"
	"strings"
)

type CommentService struct{ c *Container }

func (s *CommentService) List(ctx context.Context, targetType string, targetID int64) ([]map[string]any, error) {
	return s.c.Deps.Store.ListComments(ctx, targetType, targetID)
}

func (s *CommentService) Create(ctx context.Context, targetType string, targetID int64, content string, parentID, createdBy int64) (int64, error) {
	return s.c.Deps.Store.CreateComment(ctx, targetType, targetID, content, parentID, createdBy)
}

func (s *CommentService) Update(ctx context.Context, id, actorID int64, content string) error {
	if strings.TrimSpace(content) == "" {
		return ErrInvalidInput("content required")
	}
	return s.c.Deps.Store.UpdateComment(ctx, id, actorID, content)
}

func (s *CommentService) Delete(ctx context.Context, id, actorID int64) error {
	return s.c.Deps.Store.SoftDeleteComment(ctx, id, actorID)
}
func (s *CommentService) Like(ctx context.Context, id int64, delta int64) error {
	return s.c.Deps.Store.TouchCommentLike(ctx, id, delta)
}
