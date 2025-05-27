package model

import "time"

type Paste struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	TTL       int       `json:"ttl_seconds"`
	ExpiresAt time.Time `json:"expires_at"`
}
