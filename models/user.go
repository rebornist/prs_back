package models

import "time"

type User struct {
	ID        uint      `json:"id"`
	UID       string    `json:"uid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Age       uint8     `json:"age"`
	Birthday  string    `json:"birthday"`
	Grade     uint8     `json:"grade"`
	Mobile    string    `json:"mobile"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseUser struct {
	UID      string `json:"uid"`
	Username string `json:"username"`
	Grade    uint8  `json:"grade"`
}
