package models

import "time"

type User struct {
	Id           string    `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	DataLimit    int64     `json:"data_limit"`
	DataUsed     int64     `json:"data_used"`
	AllowedPools []string  `json:"allowed_pools"`
	IPWhitelist  []string  `json:"ip_whitelist,omitempty"`
	Status       string    `json:"status"` // active, suspended
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username     string   `json:"username" binding:"required"`
	Password     string   `json:"password" binding:"required"`
	DataLimit    int64    `json:"data_limit"`
	AllowedPools []string `json:"allowed_pools"`
	IPWhitelist  []string `json:"ip_whitelist"`
}

type UpdateUserRequest struct {
	Password     *string   `json:"password,omitempty"`
	DataLimit    *int64    `json:"data_limit,omitempty"`
	AllowedPools *[]string `json:"allowed_pools,omitempty"`
	IPWhitelist  *[]string `json:"ip_whitelist,omitempty"`
	Status       *string   `json:"status,omitempty"`
}

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Success      bool     `json:"success"`
	UserID       string   `json:"user_id,omitempty"`
	AllowedPools []string `json:"allowed_pools,omitempty"`
	DataLimit    int64    `json:"data_limit,omitempty"`
	DataUsed     int64    `json:"data_used,omitempty"`
	Message      string   `json:"message,omitempty"`
}

type UsageReport struct {
	UserID string `json:"user_id" binding:"required"`
	Bytes  int64  `json:"bytes" binding:"required"`
}

type GenerateRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	UpStream string `json:"upstream" binding:"required"`
	Country  string `json:"country" binding:"required"`
	IsSticky bool   `json:"issticky" binding:"required"`
}
