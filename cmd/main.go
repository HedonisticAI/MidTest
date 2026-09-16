package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"midtest/config"
	"midtest/internal/infrastructure/cache"
	postgres_infra "midtest/internal/infrastructure/postgres"
	cache_repo "midtest/internal/repo/cache"
	postgres_repo "midtest/internal/repo/postgres"
	http_transport "midtest/internal/transport/http"
	"midtest/internal/usecase"
)

func main() {
	cfg := config.NewConfig()
	if cfg == nil {
		log.Fatal("failed to load application config")
	}

	ctx := context.Background()

	// PostgreSQL
	db, err := postgres_infra.Open(ctx, cfg.DB_DSN)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	log.Println("connected to postgres")

	// Cache
	appCache := cache.NewCache(
		10*time.Minute,
		5*time.Minute,
	)

	cacheRepo := &cache_repo.RepoCache{
		Cache: appCache,
	}

	// PostgreSQL repository
	postgresRepo := postgres_repo.NewRepo(db)

	// Usecase
	uc := usecase.NewUsecase(
		cfg.ADMToken,
		"./storage",
		postgresRepo,
		cacheRepo,
	)

	// HTTP router
	router := http_transport.NewRouter(uc)

	server := &http.Server{
		Addr:              cfg.HttpPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		log.Printf("HTTP server started on %s", cfg.HttpPort)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	<-stop

	log.Println("shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("HTTP server stopped")
}
