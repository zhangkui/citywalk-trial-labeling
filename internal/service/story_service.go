package service

import (
	"context"
	"strings"

	"citywalk/internal/repository"
)

type StoryService struct{ c *Container }

func (s *StoryService) List(ctx context.Context, sortBy string, page, pageSize int) ([]map[string]any, int64, error) {
	limit, offset := Page{Page: page, PageSize: pageSize}.limitOffset()
	return s.c.Deps.Store.ListStories(ctx, repository.StoryFilter{SortBy: sortBy, Limit: limit, Offset: offset})
}

func (s *StoryService) Detail(ctx context.Context, id int64) (map[string]any, error) {
	story, err := s.c.Deps.Store.FindStory(ctx, id)
	if err != nil {
		return nil, err
	}
	_, _ = s.c.Deps.Store.Exec(ctx, `UPDATE stories SET view_count = view_count + 1 WHERE id = ?`, id)
	return story, nil
}

func (s *StoryService) Create(ctx context.Context, title, content, coverImage string, landmarkIDs []int64, routeID *int64, createdBy int64, media []map[string]any) (map[string]any, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		return nil, ErrInvalidInput("title/content required")
	}
	id, err := s.c.Deps.Store.CreateStory(ctx, title, content, coverImage, landmarkIDs, routeID, createdBy, 1)
	if err != nil {
		return nil, err
	}
	for idx, item := range media {
		if url, ok := item["url"].(string); ok {
			var typ int8 = 1
			if t, ok := item["type"].(float64); ok {
				typ = int8(t)
			}
			_ = s.c.Deps.Store.AddStoryMedia(ctx, id, typ, url, idx)
		}
	}
	_ = s.c.Badge.UnlockByCondition(ctx, createdBy, "create_story", 1)
	return s.c.Deps.Store.FindStory(ctx, id)
}

func (s *StoryService) Update(ctx context.Context, id int64, title, content, coverImage string, landmarkIDs []int64, routeID *int64, status int8, media []map[string]any) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		return ErrInvalidInput("title/content required")
	}
	if err := s.c.Deps.Store.UpdateStory(ctx, id, title, content, coverImage, landmarkIDs, routeID, status); err != nil {
		return err
	}
	return s.c.Deps.Store.ReplaceStoryMedia(ctx, id, normalizeStoryMedia(media))
}

func normalizeStoryMedia(media []map[string]any) []map[string]any {
	result := make([]map[string]any, 0, len(media))
	for _, item := range media {
		url, _ := item["url"].(string)
		var typ int8 = 1
		if value, ok := item["type"].(float64); ok {
			typ = int8(value)
		}
		result = append(result, map[string]any{"type": typ, "url": strings.TrimSpace(url)})
	}
	return result
}

func (s *StoryService) Delete(ctx context.Context, id int64) error {
	return s.c.Deps.Store.SoftDeleteStory(ctx, id)
}
func (s *StoryService) Review(ctx context.Context, id int64, status int8) error {
	return s.c.Deps.Store.UpdateStoryStatus(ctx, id, status)
}
func (s *StoryService) Like(ctx context.Context, userID, storyID int64, add bool) error {
	return s.c.Deps.Store.LikeStory(ctx, userID, storyID, add)
}
