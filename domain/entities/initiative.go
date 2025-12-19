package entities

import "time"

type Initiative struct {
	ID                  int64                       `json:"id"`
	Title               string                      `json:"title"`
	Description         string                      `json:"description"`
	Benefits            string                      `json:"benefits"`
	Status              string                      `json:"status"`
	Type                string                      `json:"type"`
	Priority            string                      `json:"priority"`
	Sector              string                      `json:"sector"`
	OwnerID             int64                       `json:"owner_id"`
	OwnerName           string                      `json:"owner_name"`
	Deadline            *time.Time                  `json:"deadline,omitempty"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
	CancellationRequest *InitiativeCancellationInfo `json:"cancellation_request,omitempty"` // NOVO
}

// Informações resumidas sobre solicitação de cancelamento
type InitiativeCancellationInfo struct {
	ID                int64      `json:"id"`
	Status            string     `json:"status"` // Pendente, Aprovada, Reprovada
	RequestedByUserID int64      `json:"requested_by_user_id"`
	RequestedByName   string     `json:"requested_by_name"`
	Reason            string     `json:"reason"`
	ReviewedByUserID  *int64     `json:"reviewed_by_user_id,omitempty"`
	ReviewedByName    string     `json:"reviewed_by_name,omitempty"`
	ReviewReason      string     `json:"review_reason,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	ReviewedAt        *time.Time `json:"reviewed_at,omitempty"`
}

type InitiativeListResponse struct {
	ID                  int64                       `json:"id"`
	Title               string                      `json:"title"`
	Description         string                      `json:"description"`
	Status              string                      `json:"status"`
	Type                string                      `json:"type"`
	Priority            string                      `json:"priority"`
	Sector              string                      `json:"sector"`
	OwnerName           string                      `json:"owner_name"`
	Date                string                      `json:"date"`
	CancellationRequest *InitiativeCancellationInfo `json:"cancellation_request,omitempty"` // NOVO
}

type CreateInitiativeRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Benefits    string  `json:"benefits"`
	Type        string  `json:"type"`
	Priority    string  `json:"priority"`
	Sector      string  `json:"sector"`
	Deadline    *string `json:"deadline,omitempty"`
}

type UpdateInitiativeRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Benefits    *string `json:"benefits,omitempty"`
	Type        *string `json:"type,omitempty"`
	Priority    *string `json:"priority,omitempty"`
	Sector      *string `json:"sector,omitempty"`
	Deadline    *string `json:"deadline,omitempty"`
}

type ChangeInitiativeStatusRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type InitiativeFilter struct {
	Search   string
	Status   string
	Type     string
	Sector   string
	Priority string
	SectorID *int64 // NOVO: Filtro por ID do setor
}

// Status da iniciativa
const (
	StatusSubmitted   = "Submetida"
	StatusInAnalysis  = "Em Análise"
	StatusApproved    = "Aprovada"
	StatusInExecution = "Em Execução"
	StatusReturned    = "Devolvida"
	StatusRejected    = "Reprovada"
	StatusInHomolog   = "Em Homologação"
	StatusCompleted   = "Concluída"
	StatusCancelled   = "Cancelada"
)

// Tipos de iniciativa
const (
	TypeAutomation  = "Automação"
	TypeIntegration = "Integração"
	TypeImprovement = "Melhoria"
	TypeNewProject  = "Novo Projeto"
)

// Prioridades
const (
	PriorityHigh   = "Alta"
	PriorityMedium = "Média"
	PriorityLow    = "Baixa"
)

type ReviewInitiativeRequest struct {
	Approved bool   `json:"approved"` // true = aprovar, false = reprovar
	Reason   string `json:"reason"`   // Justificativa obrigatória
}
