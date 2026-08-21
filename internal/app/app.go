package app

import (
	"context"
	"database/sql"
	"log/slog"

	"citywalk/internal/config"
	"citywalk/internal/platform/cache"
	"citywalk/internal/platform/crypto"
	"citywalk/internal/platform/db"
	"citywalk/internal/platform/logger"
	"citywalk/internal/platform/scheduler"
	"citywalk/internal/platform/storage"
	"citywalk/internal/repository"
	"citywalk/internal/service"
	httptransport "citywalk/internal/transport/http"
)

type App struct {
	cfg       *config.Config
	logger    *slog.Logger
	db        *sql.DB
	redis     *cache.Redis
	jwt       *crypto.JWT
	store     *repository.Store
	services  *service.Container
	router    httptransport.Router
	scheduler *scheduler.Scheduler
}

func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	if log == nil {
		log = logger.New("info")
	}
	database, err := db.Open(cfg.MySQLDSN())
	if err != nil {
		return nil, err
	}
	if err := db.RunMigrations(context.Background(), database, "migrations"); err != nil {
		return nil, err
	}
	if err := db.SeedDefaults(context.Background(), database, db.SeedConfig{AdminUsername: cfg.AdminUsername, AdminPassword: cfg.AdminPassword}); err != nil {
		return nil, err
	}
	redisClient := cache.New(cfg.RedisAddr(), cfg.RedisPassword)
	jwtSvc := crypto.NewJWT(cfg.JWTSecret)
	store := repository.NewStore(database)
	fileStore := storage.NewLocal(cfg.DataDir, "/static")
	services := service.NewContainer(service.Dependencies{
		Cfg:     cfg,
		DB:      database,
		Redis:   redisClient,
		JWT:     jwtSvc,
		Store:   store,
		Storage: fileStore,
		Logger:  log,
	})
	app := &App{
		cfg:       cfg,
		logger:    log,
		db:        database,
		redis:     redisClient,
		jwt:       jwtSvc,
		store:     store,
		services:  services,
		scheduler: scheduler.New(),
	}
	app.router = httptransport.NewRouter(httptransport.Dependencies{
		Cfg:      cfg,
		Logger:   log,
		Redis:    redisClient,
		JWT:      jwtSvc,
		Store:    store,
		Services: services,
	})
	return app, nil
}

func (a *App) Router() httptransport.Router { return a.router }

func (a *App) Close() error {
	if a.scheduler != nil {
		// no-op for now
	}
	if a.redis != nil {
		_ = a.redis.Close()
	}
	if a.db != nil {
		return a.db.Close()
	}
	return nil
}
