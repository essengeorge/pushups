package domain

import "time"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type Pushup struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Count     int       `json:"count"`
	CreatedAt time.Time `json:"created_at"`
}

type ChartPoint struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

type UserStats struct {
	Username string       `json:"username"`
	Color    string       `json:"color"`
	Points   []ChartPoint `json:"points"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
