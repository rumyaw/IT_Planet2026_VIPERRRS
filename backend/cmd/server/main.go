package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trumplin/internal/auth"
	"trumplin/internal/config"
	"trumplin/internal/db"
	apirouter "trumplin/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	database, err := db.New(ctx, cfg.DB.DSN)
	if err != nil {
		panic(err)
	}
	defer database.Close()

	if err := db.ApplyMigrationsFromDSN(cfg.DB.DSN, cfg.DB.MigrationsDir); err != nil {
		panic(err)
	}

	// Bootstrap default admin for curator/moderation flows.
	if err := auth.EnsureAdmin(ctx, database.DB, cfg.Auth.AdminEmail, cfg.Auth.AdminPassword); err != nil {
		panic(err)
	}
	if err := auth.EnsureDemoData(ctx, database.DB); err != nil {
		panic(err)
	}

	router := apirouter.NewRouter(apirouter.RouterConfig{
		Cfg:        cfg,
		Database:  database,
		JWTSecret:  cfg.Auth.JWTSecret,
	})

	srv := &http.Server{
		Addr:              cfg.HTTP.Addr(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("server started", "addr", cfg.HTTP.ListenAddr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	slog.Info("server shutting down")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}

