package service

import (
	"context"
	"strings"

	"citywalk/internal/repository"
)

type RouteService struct{ c *Container }

func (s *RouteService) List(ctx context.Context, city string, themeID int64, sortBy string, page, pageSize int, minRate float64) ([]map[string]any, int64, error) {
	limit, offset := Page{Page: page, PageSize: pageSize}.limitOffset()
	countMinRate := minRate
	if strings.TrimSpace(city) != "" && themeID > 0 {
		countMinRate = 0
	}
	return s.c.Deps.Store.ListRoutes(ctx, repository.RouteFilter{City: city, ThemeID: themeID, MinRate: minRate, CountMinRate: countMinRate, SortBy: sortBy, Limit: limit, Offset: offset})
}

func (s *RouteService) Detail(ctx context.Context, id int64) (map[string]any, error) {
	route, err := s.c.Deps.Store.FindRoute(ctx, id)
	if err != nil {
		return nil, err
	}
	_, _ = s.c.Deps.Store.Exec(ctx, `UPDATE routes SET view_count = view_count + 1 WHERE id = ?`, id)
	return route, nil
}

func (s *RouteService) Create(ctx context.Context, title, description string, themeID *int64, city string, startLat, startLng, endLat, endLng *float64, totalDistance, duration *int, difficulty int8, coverImage string, createdBy int64, waypoints []map[string]any) (map[string]any, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrInvalidInput("title required")
	}
	status := int8(1)
	id, err := s.c.Deps.Store.CreateRoute(ctx, title, description, themeID, city, startLat, startLng, endLat, endLng, totalDistance, duration, difficulty, coverImage, createdBy, status, waypoints)
	if err != nil {
		return nil, err
	}
	_ = s.c.Badge.UnlockByCondition(ctx, createdBy, "create_route", 1)
	return s.c.Deps.Store.FindRoute(ctx, id)
}

func (s *RouteService) Update(ctx context.Context, id int64, title, description string, themeID *int64, city string, startLat, startLng, endLat, endLng *float64, totalDistance, duration *int, difficulty int8, coverImage string, status int8, waypoints []map[string]any) error {
	if strings.TrimSpace(title) == "" {
		return ErrInvalidInput("title required")
	}
	return s.c.Deps.Store.UpdateRoute(ctx, id, title, description, themeID, city, startLat, startLng, endLat, endLng, totalDistance, duration, difficulty, coverImage, status, waypoints)
}

func (s *RouteService) Delete(ctx context.Context, id int64) error {
	if _, err := s.c.Deps.Store.FindRoute(ctx, id); err != nil {
		return err
	}
	return s.c.Deps.Store.SoftDeleteRoute(ctx, id)
}

func (s *RouteService) UpdateStatus(ctx context.Context, id int64, status int8) error {
	route, err := s.c.Deps.Store.FindRoute(ctx, id)
	if err != nil {
		return err
	}
	current, _ := route["status"].(int8)
	if current == 0 && status == 1 {
		status = 2
	}
	return s.c.Deps.Store.UpdateRouteStatus(ctx, id, status)
}

func (s *RouteService) Favorite(ctx context.Context, userID, routeID int64, add bool) error {
	active, err := s.c.Deps.Store.HasUserTarget(ctx, "favorites", userID, "route", routeID)
	if err != nil {
		return err
	}
	if active == add {
		return nil
	}
	return s.c.Deps.Store.FavoriteRoute(ctx, userID, routeID, add)
}
func (s *RouteService) Rate(ctx context.Context, routeID int64, score float64) error {
	return s.c.Deps.Store.RateRoute(ctx, routeID, score)
}
