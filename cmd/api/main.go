package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"citywalk/internal/app"
	"citywalk/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	application, err := app.New(cfg, logger)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:         cfg.ListenAddr,
		Handler:      application.Router(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:   cfg.IdleTimeout,
	}
	go func() {
		logger.Info("api listening", "addr", cfg.ListenAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen failed", "err", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownAfter)
	defer cancel()
	_ = server.Shutdown(ctx)
	if application != nil {
		_ = application.Close()
	}
	time.Sleep(100 * time.Millisecond)
}
