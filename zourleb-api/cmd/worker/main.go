// Command worker runs scheduled background jobs (boost expiry, notifications,
// OTP cleanup, departure reconciliation).
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zourleb/zourleb-api/config"
	"github.com/zourleb/zourleb-api/internal/database"
	"github.com/zourleb/zourleb-api/internal/worker"
	"github.com/zourleb/zourleb-api/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.App.Env)

	db, err := database.Connect(cfg, log)
	if err != nil {
		log.Error("database connection failed", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runner := worker.NewRunner(db, log)
	go runner.Start(ctx, time.Minute)

	log.Info("worker started", "interval", "1m")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("worker shutting down")
	cancel()
}
