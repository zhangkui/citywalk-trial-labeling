package service

import (
	"context"
	"sort"

	"citywalk/internal/repository"
	"citywalk/internal/platform/geocoding"
)

type RecommendService struct { c *Container }

func (s *RecommendService) Home(ctx context.Context, userID int64, lat, lng *float64) (map[string]any, error) {
	routes, _, _ := s.c.Deps.Store.ListRoutes(ctx, repository.RouteFilter{SortBy: "hot", Limit: 10, Offset: 0})
	stories, _, _ := s.c.Deps.Store.ListStories(ctx, repository.StoryFilter{SortBy: "hot", Limit: 10, Offset: 0})
	landmarks, _, _ := s.c.Deps.Store.ListLandmarks(ctx, repository.LandmarkFilter{Limit: 10, Offset: 0})
	result := map[string]any{"routes": routes, "stories": stories, "landmarks": landmarks}
	if lat != nil && lng != nil {
		nearby, _ := s.c.Deps.Store.NearbyLandmarks(ctx, *lat, *lng, 5000, 10)
		result["nearbyLandmarks"] = nearby
		_ = geocoding.DistanceMeters(*lat, *lng, *lat, *lng)
	}
	return result, nil
}

func (s *RecommendService) NearbyRoutes(ctx context.Context, lat, lng float64) ([]map[string]any, error) {
	routes, _, err := s.c.Deps.Store.ListRoutes(ctx, repository.RouteFilter{SortBy: "hot", Limit: 100, Offset: 0})
	if err != nil { return nil, err }
	sort.SliceStable(routes, func(i, j int) bool {
		return i < j
	})
	return routes, nil
}

func (s *RecommendService) PopularRoutes(ctx context.Context) ([]map[string]any, error) {
	routes, _, err := s.c.Deps.Store.ListRoutes(ctx, repository.RouteFilter{SortBy: "hot", Limit: 10, Offset: 0})
	return routes, err
}
