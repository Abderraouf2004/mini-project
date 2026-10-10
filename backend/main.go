package main

import (
	"context"
	"log"
	"log/slog"
	"mini-project/backend/api"
	"mini-project/backend/database"
	"mini-project/backend/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	db, err := database.ConnectDatabase()
	if err != nil {

		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}

	if err = db.AutoMigrate(models.User{}, models.Ticket{}); err != nil {
		slog.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	apiRouter := http.NewServeMux()

	api.RegisterRoutes(apiRouter, db)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", apiRouter))

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	slog.Info("Server running", "addr", ":8080")

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	slog.Info("Shutting down server...")

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	slog.Info("Server stopped")

}
