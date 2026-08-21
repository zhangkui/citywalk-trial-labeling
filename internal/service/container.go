package service

import (
	"database/sql"
	"log/slog"

	"citywalk/internal/config"
	"citywalk/internal/platform/cache"
	"citywalk/internal/platform/crypto"
	"citywalk/internal/platform/storage"
	"citywalk/internal/repository"
)

type Dependencies struct {
	Cfg     *config.Config
	DB      *sql.DB
	Redis   *cache.Redis
	JWT     *crypto.JWT
	Store   *repository.Store
	Storage *storage.Local
	Logger  *slog.Logger
}

type Container struct {
	Deps       Dependencies
	Auth       *AuthService
	Admin      *AdminService
	Route      *RouteService
	Story      *StoryService
	Landmark   *LandmarkService
	Comment    *CommentService
	Event      *EventService
	Recommend  *RecommendService
	Badge      *BadgeService
}

func NewContainer(deps Dependencies) *Container {
	c := &Container{Deps: deps}
	c.Auth = &AuthService{c}
	c.Admin = &AdminService{c}
	c.Route = &RouteService{c}
	c.Story = &StoryService{c}
	c.Landmark = &LandmarkService{c}
	c.Comment = &CommentService{c}
	c.Event = &EventService{c}
	c.Recommend = &RecommendService{c}
	c.Badge = &BadgeService{c}
	return c
}

