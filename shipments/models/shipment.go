package models

import "time"

type Shipment struct {
	ID          string    `json:"id"`
	Origin      string    `json:"origin"`
	Destination string    `json:"destination"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
