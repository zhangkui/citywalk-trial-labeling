package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"citywalk/internal/repository"
)

type ExportService struct { c *Container }

func (s *ExportService) ExportRoutes(ctx context.Context, path string) (string, error) {
	routes, _, err := s.c.Deps.Store.ListRoutes(ctx, repository.RouteFilter{SortBy: "latest", Limit: 1000, Offset: 0})
	if err != nil { return "", err }
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return "", err }
	file, err := os.Create(path)
	if err != nil { return "", err }
	defer file.Close()
	if err := json.NewEncoder(file).Encode(map[string]any{"generatedAt": time.Now(), "routes": routes}); err != nil { return "", err }
	return path, nil
}

func (s *ExportService) ExportStories(ctx context.Context, path string) (string, error) {
	stories, _, err := s.c.Deps.Store.ListStories(ctx, repository.StoryFilter{SortBy: "latest", Limit: 1000, Offset: 0})
	if err != nil { return "", err }
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return "", err }
	file, err := os.Create(path)
	if err != nil { return "", err }
	defer file.Close()
	if err := json.NewEncoder(file).Encode(map[string]any{"generatedAt": time.Now(), "stories": stories}); err != nil { return "", err }
	return path, nil
}
