package domain

import "time"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	IsApproved   int    `json:"is_approved"`
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

type Friendship struct {
	UserID   int    `json:"user_id"`
	FriendID int    `json:"friend_id"`
	Username string `json:"username,omitempty"`
	Status   string `json:"status"`
}

type FriendRequest struct {
	TargetUsername string `json:"target_username"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
