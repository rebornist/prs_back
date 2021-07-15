package models

import "time"

type Session struct {
	ID         uint      `json:"id"`
	UID        string    `json:"uid"`
	TokenValue string    `json:"token_value"`
	Expires    int64     `json:"expires"`
	CreatedAt  time.Time `json:"created_at"`
}
