package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"shipments-service/database"
	"shipments-service/models"

	"github.com/google/uuid"
)

// 🟢 Verificación de estado
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "shipments"})
}

// 📦 Obtener todos los envíos
func GetShipments(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, origin, destination, status, created_at FROM shipments ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var shipments []models.Shipment
	for rows.Next() {
		var s models.Shipment
		rows.Scan(&s.ID, &s.Origin, &s.Destination, &s.Status, &s.CreatedAt)
		shipments = append(shipments, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shipments)
}

// 🚚 Crear un nuevo envío y notificar al otro servicio
func CreateShipment(w http.ResponseWriter, r *http.Request) {
	var s models.Shipment
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "Formato JSON inválido", 400)
		return
	}

	s.ID = uuid.New().String()
	s.Status = "CREATED"

	// Guardar en la base de datos
	query := `INSERT INTO shipments (id, origin, destination, status) VALUES ($1,$2,$3,$4)`
	_, err = database.DB.Exec(query, s.ID, s.Origin, s.Destination, s.Status)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 🔔 Enviar notificación al microservicio Notifications
	notif := map[string]string{
		"recipient": "cliente@fasttrack.com",
		"channel":   "email",
		"message":   fmt.Sprintf("Tu envío desde %s hacia %s ha sido registrado correctamente.", s.Origin, s.Destination),
	}
	notifJSON, _ := json.Marshal(notif)

	resp, err := http.Post("http://localhost:3002/notifications", "application/json", bytes.NewBuffer(notifJSON))
	if err != nil {
		fmt.Println("⚠️ No se pudo enviar notificación:", err)
	} else {
		defer resp.Body.Close()
		fmt.Println("✅ Notificación enviada al servicio Notifications")
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}
