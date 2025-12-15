package entities

import "time"

type Comment struct {
	ID           int64     `json:"id"`
	InitiativeID int64     `json:"initiative_id"`
	UserID       int64     `json:"user_id"`
	UserName     string    `json:"user_name"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateCommentRequest struct {
	Content string `json:"content"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}

type CommentListResponse struct {
	ID           int64  `json:"id"`
	InitiativeID int64  `json:"initiative_id"`
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	Content      string `json:"content"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
