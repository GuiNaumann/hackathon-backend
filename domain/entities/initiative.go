package entities

import "time"

type Initiative struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Benefits    string     `json:"benefits"`
	Status      string     `json:"status"`
	Type        string     `json:"type"`
	Priority    string     `json:"priority"`
	Sector      string     `json:"sector"`
	OwnerID     int64      `json:"owner_id"`
	OwnerName   string     `json:"owner_name"`
	Deadline    *time.Time `json:"deadline,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateInitiativeRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Benefits    string  `json:"benefits"`
	Type        string  `json:"type"`
	Priority    string  `json:"priority"`
	Sector      string  `json:"sector"`
	Deadline    *string `json:"deadline,omitempty"` // ISO format string
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
	Reason string `json:"reason,omitempty"` // Para status como "Devolvida" ou "Reprovada"
}

type InitiativeFilter struct {
	Search   string
	Status   string
	Type     string
	Sector   string
	Priority string
}

type InitiativeListResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Priority    string `json:"priority"`
	Sector      string `json:"sector"`
	OwnerName   string `json:"owner"`
	Date        string `json:"date"`
}

// Constantes de Status
const (
	StatusSubmitted   = "Submetida"   // Recém criada
	StatusInAnalysis  = "Em Análise"  // Sendo analisada
	StatusApproved    = "Aprovada"    // Aprovada para execução
	StatusInExecution = "Em Execução" // Sendo executada
	StatusReturned    = "Devolvida"   // Devolvida para ajustes
	StatusRejected    = "Reprovada"   // Rejeitada
	StatusCompleted   = "Concluída"   // Finalizada
	StatusCancelled   = "Cancelada"   // Cancelada
)

// Constantes de Tipo
const (
	TypeAutomation  = "Automação"
	TypeIntegration = "Integração"
	TypeImprovement = "Melhoria"
	TypeNewProject  = "Novo Projeto"
)

// Constantes de Prioridade
const (
	PriorityHigh   = "Alta"
	PriorityMedium = "Média"
	PriorityLow    = "Baixa"
)
