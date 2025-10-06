package models

import "time"

type Notification struct {
	ID        string    `json:"id"`
	Recipient string    `json:"recipient"` // destino: email/telefono/token
	Channel   string    `json:"channel"`   // email | sms | push
	Message   string    `json:"message"`
	Status    string    `json:"status"` // PENDING | SENT | FAILED
	CreatedAt time.Time `json:"created_at"`
}
