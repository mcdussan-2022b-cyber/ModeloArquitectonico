package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"notifications/database"
	"notifications/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3002" // diferente a Shipments (3001)
	}

	database.Connect()

	http.HandleFunc("/health", handlers.Health)
	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.ListNotifications(w, r)
		case http.MethodPost:
			handlers.CreateNotification(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("📣 notifications-service en puerto:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
