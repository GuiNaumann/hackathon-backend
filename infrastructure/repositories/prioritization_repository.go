package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type PrioritizationRepository interface {
	// Priorização
	Create(ctx context.Context, prioritization *entities.InitiativePrioritization) error
	Update(ctx context.Context, prioritization *entities.InitiativePrioritization) error
	GetBySectorAndYear(ctx context.Context, sectorID int64, year int) (*entities.InitiativePrioritization, error)
	GetAllByYear(ctx context.Context, year int) ([]*entities.InitiativePrioritization, error)
	LockPrioritization(ctx context.Context, prioritizationID int64) error
	UnlockPrioritization(ctx context.Context, prioritizationID int64) error

	// Solicitações de mudança
	CreateChangeRequest(ctx context.Context, request *entities.PrioritizationChangeRequest) error
	GetChangeRequestByID(ctx context.Context, requestID int64) (*entities.PrioritizationChangeRequest, error)
	ListPendingChangeRequests(ctx context.Context) ([]*entities.PrioritizationChangeRequest, error)
	UpdateChangeRequestStatus(ctx context.Context, requestID int64, status string, reviewedByUserID int64, reviewReason string) error
	HasPendingChangeRequest(ctx context.Context, prioritizationID int64) (bool, error)
}
