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
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
}

type LoginFailureData struct {
	RemainingAttempts int   `json:"remaining_attempts"`
	RetryAfterSeconds int64 `json:"retry_after_seconds"`
	Locked            bool  `json:"locked"`
}

type UserProfile struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== Conversion ====================

// Conversion 完整记录（内部使用）
type Conversion struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Markdown  string    `json:"markdown,omitempty"`
	HTML      string    `json:"html,omitempty"`
	CharCount int       `json:"char_count"`
	Theme     string    `json:"theme"`
	IsDelete  int       `json:"is_delete"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HistoryListItem 历史列表项（不含 markdown/html 大字段）
type HistoryListItem struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CharCount int       `json:"char_count"`
	Theme     string    `json:"theme"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HistoryDetail 历史详情（含完整 markdown）
type HistoryDetail struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Markdown  string    `json:"markdown"`
	HTML      string    `json:"html,omitempty"`
	CharCount int       `json:"char_count"`
	Theme     string    `json:"theme"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ConvertRequest struct {
	Markdown string `json:"markdown" binding:"required"`
	Theme    string `json:"theme"`
}

type ConvertResponse struct {
	HTML      string `json:"html"`
	Formatted string `json:"formatted"`
}

// SaveHistoryRequest 保存历史请求（html 可选）
type SaveHistoryRequest struct {
	Markdown string `json:"markdown" binding:"required"`
	HTML     string `json:"html"`
	Theme    string `json:"theme"`
}

type RenameHistoryRequest struct {
	Title string `json:"title" binding:"required"`
}

// HistoryPageResponse 历史列表分页响应（轻量）
type HistoryPageResponse struct {
	Data       []HistoryListItem `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
	TotalPages int               `json:"total_pages"`
}

// ==================== Theme ====================

type Theme struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ==================== JSON Formatter ====================

type JsonFormatRequest struct {
	Input    string `json:"input" binding:"required"`
	Indent   int    `json:"indent"`
	Minify   bool   `json:"minify"`
}

type JsonFormatResponse struct {
	Output   string `json:"output"`
	Size     int    `json:"size"`
	Minified bool   `json:"minified"`
}

type JsonValidateRequest struct {
	Input string `json:"input" binding:"required"`
}

type JsonValidateResponse struct {
	Valid  bool   `json:"valid"`
	Error  string `json:"error,omitempty"`
	Keys   int    `json:"keys,omitempty"`
	Depth  int    `json:"depth,omitempty"`
}

type JsonHistoryListItem struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CharCount int       `json:"char_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JsonHistoryDetail struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	CharCount int       `json:"char_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type JsonSaveHistoryRequest struct {
	Input  string `json:"input" binding:"required"`
	Output string `json:"output"`
}

type JsonHistoryPageResponse struct {
	Data       []JsonHistoryListItem `json:"data"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Size       int                   `json:"size"`
	TotalPages int                   `json:"total_pages"`
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
