package models

import "time"

type Note struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CategoryName string    `json:"category_name"`
	CategoryID   int       `json:"category_id"`
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateNoteRequest struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	CategoryID int    `json:"category_id" validate:"required"`
}

type UpdateNoteRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID int    `json:"category_id"`
}
