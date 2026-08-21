package service

import (
	"context"
	"strings"

	"citywalk/internal/repository"
)

type LandmarkService struct{ c *Container }

func (s *LandmarkService) List(ctx context.Context, categoryID int64, page, pageSize int) ([]map[string]any, int64, error) {
	limit, offset := Page{Page: page, PageSize: pageSize}.limitOffset()
	return s.c.Deps.Store.ListLandmarks(ctx, repository.LandmarkFilter{CategoryID: categoryID, Limit: limit, Offset: offset})
}

func (s *LandmarkService) Detail(ctx context.Context, id int64) (map[string]any, error) {
	landmark, err := s.c.Deps.Store.FindLandmark(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = s.c.Deps.Store.TouchLandmarkView(ctx, id)
	return landmark, nil
}

func (s *LandmarkService) Create(ctx context.Context, name, address string, lat, lng float64, categoryID *int64, description, coverImage string, createdBy int64) (map[string]any, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidInput("name required")
	}
	id, err := s.c.Deps.Store.CreateLandmark(ctx, name, address, lat, lng, categoryID, description, coverImage, createdBy, 1)
	if err != nil {
		return nil, err
	}
	if err := s.c.Badge.UnlockByCondition(ctx, createdBy, "create_landmark", 1); err != nil {
		return nil, err
	}
	return s.c.Deps.Store.FindLandmark(ctx, id)
}

func (s *LandmarkService) Update(ctx context.Context, id int64, name, address string, lat, lng float64, categoryID *int64, description, coverImage string, status int8) error {
	return s.c.Deps.Store.UpdateLandmark(ctx, id, name, address, lat, lng, categoryID, description, coverImage, status)
}

func (s *LandmarkService) Delete(ctx context.Context, id int64) error {
	return s.c.Deps.Store.SoftDeleteLandmark(ctx, id)
}
func (s *LandmarkService) Nearby(ctx context.Context, lat, lng, radius float64, limit int) ([]map[string]any, error) {
	return s.c.Deps.Store.NearbyLandmarks(ctx, lat, lng, radius, limit)
}
func (s *LandmarkService) Rate(ctx context.Context, id int64, score float64) error {
	return s.c.Deps.Store.RateLandmark(ctx, id, score)
}
