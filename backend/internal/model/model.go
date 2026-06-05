package model

import "time"

// ==================== User ====================

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=128"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type UserProfile struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== Conversion ====================

type Conversion struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Markdown  string    `json:"markdown"`
	HTML      string    `json:"html"`
	Theme     string    `json:"theme"`
	CreatedAt time.Time `json:"created_at"`
}

type ConvertRequest struct {
	Markdown string `json:"markdown" binding:"required"`
	Theme    string `json:"theme"`
}

type ConvertResponse struct {
	HTML       string `json:"html"`
	Formatted  string `json:"formatted"`
}

type SaveHistoryRequest struct {
	Markdown string `json:"markdown" binding:"required"`
	HTML     string `json:"html" binding:"required"`
	Theme    string `json:"theme"`
}

type HistoryPageResponse struct {
	Data  []Conversion `json:"data"`
	Total int64        `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

// ==================== Theme ====================

type Theme struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ==================== Common ====================

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type PageRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"pageSize,default=20"`
}
