package main

import (
	"fmt"
	"log"
	"mini-project/backend/api"
	"mini-project/backend/database"
	"mini-project/backend/models"
	"net/http"
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

	// http.Handle("/api/", http.StripPrefix("/api", apiRouter))
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
	fmt.Println("Server running on port 8080")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server error:", err)
	}
	// if err = http.ListenAndServe(":8080", nil); err != nil {
	// 	fmt.Println("Server error:", err)
	// }
}
