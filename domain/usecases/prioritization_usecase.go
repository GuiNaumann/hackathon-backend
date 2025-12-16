package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type PrioritizationUseCase interface {
	// Priorização
	SavePrioritization(ctx context.Context, req *entities.SavePrioritizationRequest, userID int64) (*entities.PrioritizationWithInitiatives, error)
	GetPrioritization(ctx context.Context, year int, userID int64) (*entities.PrioritizationWithInitiatives, error)
	GetAllSectorsPrioritization(ctx context.Context, year int, userID int64) (*entities.AllSectorsPrioritization, error)

	// Solicitações de mudança
	RequestPrioritizationChange(ctx context.Context, req *entities.RequestPrioritizationChangeRequest, userID int64, year int) (*entities.PrioritizationChangeRequest, error)
	ReviewPrioritizationChange(ctx context.Context, requestID int64, req *entities.ReviewPrioritizationChangeRequest, userID int64) error
	ListPendingChangeRequests(ctx context.Context, userID int64) ([]*entities.PrioritizationChangeRequest, error)
}
