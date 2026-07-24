// Command api is the Zourleb HTTP API entrypoint.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/zourleb/zourleb-api/config"
	"github.com/zourleb/zourleb-api/internal/app"
	"github.com/zourleb/zourleb-api/internal/database"
	"github.com/zourleb/zourleb-api/internal/router"
	"github.com/zourleb/zourleb-api/internal/validator"
	"github.com/zourleb/zourleb-api/pkg/logger"
	"github.com/zourleb/zourleb-api/seeds"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.App.Env)

	// database
	db, err := database.Connect(cfg, log)
	if err != nil {
		log.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	if err := database.Migrate(db, log); err != nil {
		log.Error("migration failed", "err", err)
		os.Exit(1)
	}
	if err := seeds.Run(db, cfg, log); err != nil {
		log.Error("seeding failed", "err", err)
		os.Exit(1)
	}

	// HTTP
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	validator.Install(e)

	container := app.New(cfg, db, log)
	router.Register(e, container)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		log.Info("server starting", "port", cfg.HTTP.Port, "env", cfg.App.Env)
		if err := e.StartServer(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
	}
	_ = time.Now
}
