package entities

import "time"

type InitiativeHistory struct {
	ID           int64     `json:"id"`
	InitiativeID int64     `json:"initiative_id"`
	UserID       int64     `json:"user_id"`
	UserName     string    `json:"user_name"`
	OldStatus    string    `json:"old_status"`
	NewStatus    string    `json:"new_status"`
	Reason       string    `json:"reason,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type InitiativeHistoryResponse struct {
	ID           int64  `json:"id"`
	InitiativeID int64  `json:"initiative_id"`
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	OldStatus    string `json:"old_status"`
	NewStatus    string `json:"new_status"`
	Reason       string `json:"reason,omitempty"`
	CreatedAt    string `json:"created_at"`
	TimeAgo      string `json:"time_ago"`
}
