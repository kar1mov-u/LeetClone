package models

import "time"

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type User struct {
	Username   string
	Email      string
	Role       string
	Created_at time.Time
}

// username, email, password, created_at,
