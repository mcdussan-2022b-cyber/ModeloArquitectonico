package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"shipments-service/database"
	"shipments-service/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	database.Connect()

	http.HandleFunc("/health", handlers.Health)
	http.HandleFunc("/shipments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetShipments(w, r)
		} else if r.Method == http.MethodPost {
			handlers.CreateShipment(w, r)
		} else {
			http.Error(w, "Método no permitido", 405)
		}
	})

	fmt.Println("🚚 Servidor ejecutándose en puerto:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
