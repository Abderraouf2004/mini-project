package main

import (
	"context"
	"fmt"
	"log"
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
	db, err := database.ConnectDatabase()
	if err != nil {
		fmt.Println("Database connection error:", err)
		return
	}

	if err = db.AutoMigrate(models.User{}, models.Ticket{}); err != nil {
		fmt.Println("Migration error:", err)
		return
	}
	apiRouter := http.NewServeMux()

	api.RegisterRoutes(apiRouter)

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
	// fmt.Println("Server running on port 8080")

	// if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 	log.Fatal("Server error:", err)
	// }
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	fmt.Println("Server running on port 8080")

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

	fmt.Println("Shutting down server...")

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	fmt.Println("Server stopped")

}
