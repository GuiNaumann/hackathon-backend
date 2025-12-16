package entities

import "time"

type InitiativeCancellationRequest struct {
	ID                int64      `json:"id"`
	InitiativeID      int64      `json:"initiative_id"`
	RequestedByUserID int64      `json:"requested_by_user_id"`
	RequestedByName   string     `json:"requested_by_name"`
	Reason            string     `json:"reason"`
	Status            string     `json:"status"` // Pendente, Aprovada, Reprovada
	ReviewedByUserID  *int64     `json:"reviewed_by_user_id,omitempty"`
	ReviewedByName    string     `json:"reviewed_by_name,omitempty"`
	ReviewReason      string     `json:"review_reason,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
}

type RequestCancellationRequest struct {
	Reason string `json:"reason"`
}

type ReviewCancellationRequest struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

type CancellationRequestResponse struct {
	ID                int64  `json:"id"`
	InitiativeID      int64  `json:"initiative_id"`
	InitiativeTitle   string `json:"initiative_title"`
	RequestedByUserID int64  `json:"requested_by_user_id"`
	RequestedByName   string `json:"requested_by_name"`
	Reason            string `json:"reason"`
	Status            string `json:"status"`
	ReviewedByName    string `json:"reviewed_by_name,omitempty"`
	ReviewReason      string `json:"review_reason,omitempty"`
	CreatedAt         string `json:"created_at"`
	ReviewedAt        string `json:"reviewed_at,omitempty"`
	TimeAgo           string `json:"time_ago"`
}

// Status de solicitação de cancelamento
const (
	CancellationStatusPending  = "Pendente"
	CancellationStatusApproved = "Aprovada"
	CancellationStatusRejected = "Reprovada"
)
