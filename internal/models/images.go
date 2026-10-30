package models

import "time"

const (
	StatusPending = "PENDING"
	StatusReady   = "READY"
)

type Image struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Status    string    `json:"status"`
	CDNUrl    string    `json:"cdn_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
