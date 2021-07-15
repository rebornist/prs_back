package models

import "time"

type Logger struct {
	Body       string    `json:"body"`
	ConnectIp  string    `json:"connect_ip"`
	RequestId  string    `json:"request_id"`
	RequestUrl string    `json:"request_url"`
	Message    string    `json:"message"`
	Status     int       `json:"status"`
	UserAgent  string    `json:"user_agent"`
	Backoff    int64     `json:"backoff"`
	CreatedAt  time.Time `json:"created_at"`
}
