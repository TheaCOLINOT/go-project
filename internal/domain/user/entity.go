package user

import "time"

type User struct {
	ID           string
	BusinessID   string
	Email        string
	PasswordHash string
	Role         string
	IsConfirmed  bool
	CreatedAt    time.Time
}

type AuthToken struct {
	ID        string
	UserID    string
	TokenHash string
	Purpose   string
	ExpiresAt time.Time
	UsedAt    *time.Time
}

type AdminAction struct {
	ID               string
	AdminUserID      string
	ActionType       string
	TargetType       string
	TargetBusinessID *string
	Details          map[string]any
	CreatedAt        time.Time
}
