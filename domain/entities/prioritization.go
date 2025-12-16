package entities

import "time"

// Priorização de iniciativas por setor
type InitiativePrioritization struct {
	ID              int64     `json:"id"`
	SectorID        int64     `json:"sector_id"`
	SectorName      string    `json:"sector_name"`
	Year            int       `json:"year"`
	PriorityOrder   []int64   `json:"priority_order"` // Array com IDs das iniciativas ordenadas
	IsLocked        bool      `json:"is_locked"`      // Se está bloqueada para edição
	CreatedByUserID int64     `json:"created_by_user_id"`
	CreatedByName   string    `json:"created_by_name"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Request para salvar priorização
type SavePrioritizationRequest struct {
	Year          int     `json:"year"`
	PriorityOrder []int64 `json:"priority_order"` // Array com IDs das iniciativas na ordem
}

// Request para solicitar mudança de priorização
type RequestPrioritizationChangeRequest struct {
	NewPriorityOrder []int64 `json:"new_priority_order"`
	Reason           string  `json:"reason"`
}

// Solicitação de mudança de priorização
type PrioritizationChangeRequest struct {
	ID                int64      `json:"id"`
	PrioritizationID  int64      `json:"prioritization_id"`
	RequestedByUserID int64      `json:"requested_by_user_id"`
	RequestedByName   string     `json:"requested_by_name"`
	NewPriorityOrder  []int64    `json:"new_priority_order"`
	Reason            string     `json:"reason"`
	Status            string     `json:"status"` // Pendente, Aprovada, Reprovada
	ReviewedByUserID  *int64     `json:"reviewed_by_user_id,omitempty"`
	ReviewedByName    string     `json:"reviewed_by_name,omitempty"`
	ReviewReason      string     `json:"review_reason,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
}

// Request para revisar solicitação de mudança
type ReviewPrioritizationChangeRequest struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

// Response com iniciativas priorizadas (com dados completos)
type PrioritizationWithInitiatives struct {
	ID              int64                     `json:"id"`
	SectorID        int64                     `json:"sector_id"`
	SectorName      string                    `json:"sector_name"`
	Year            int                       `json:"year"`
	IsLocked        bool                      `json:"is_locked"`
	Initiatives     []*InitiativeListResponse `json:"initiatives"` // Iniciativas na ordem de prioridade
	CreatedByUserID int64                     `json:"created_by_user_id"`
	CreatedByName   string                    `json:"created_by_name"`
	CreatedAt       string                    `json:"created_at"`
	UpdatedAt       string                    `json:"updated_at"`
}

// Response para admin/manager com todos os setores
type AllSectorsPrioritization struct {
	Year    int                              `json:"year"`
	Sectors []*PrioritizationWithInitiatives `json:"sectors"`
}

// Status da solicitação de mudança
const (
	PrioritizationChangeStatusPending  = "Pendente"
	PrioritizationChangeStatusApproved = "Aprovada"
	PrioritizationChangeStatusRejected = "Reprovada"
)
