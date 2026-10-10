package main

import (
	"encoding/json"
	"fmt"
	"mini-project/backend/api"
	"mini-project/backend/database"
	"mini-project/backend/models"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello from my API",
	})
}
func main() {
	fmt.Println("Server running on port 8080")
	http.HandleFunc("/", helloHandler)
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

	http.Handle("/api/", http.StripPrefix("/api", apiRouter))
	if err = http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
