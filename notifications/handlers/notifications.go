package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"notifications/database"
	"notifications/models"

	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "notifications",
	})
}

func ListNotifications(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(
		`SELECT id, recipient, channel, message, status, created_at
		 FROM notifications
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	defer rows.Close()

	items := make([]models.Notification, 0)
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.Recipient, &n.Channel, &n.Message, &n.Status, &n.CreatedAt); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		items = append(items, n)
	}
	writeJSON(w, http.StatusOK, items)
}

func CreateNotification(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Recipient string `json:"recipient"`
		Channel   string `json:"channel"` // email | sms | push
		Message   string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON inválido"})
		return
	}
	if in.Recipient == "" || in.Channel == "" || in.Message == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "recipient, channel y message son requeridos"})
		return
	}

	id := uuid.New().String()
	_, err := database.DB.Exec(
		`INSERT INTO notifications (id, recipient, channel, message, status)
		 VALUES ($1, $2, $3, $4, 'PENDING')`,
		id, in.Recipient, in.Channel, in.Message,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var out models.Notification
	err = database.DB.QueryRow(
		`SELECT id, recipient, channel, message, status, created_at
		   FROM notifications WHERE id = $1`, id,
	).Scan(&out.ID, &out.Recipient, &out.Channel, &out.Message, &out.Status, &out.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, out)
}
