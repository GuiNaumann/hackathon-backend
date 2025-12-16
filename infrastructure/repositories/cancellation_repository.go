package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type CancellationRepository interface {
	Create(ctx context.Context, request *entities.InitiativeCancellationRequest) error
	GetByID(ctx context.Context, requestID int64) (*entities.InitiativeCancellationRequest, error)
	GetPendingByInitiative(ctx context.Context, initiativeID int64) (*entities.InitiativeCancellationRequest, error)
	ListPending(ctx context.Context) ([]*entities.InitiativeCancellationRequest, error)
	UpdateStatus(ctx context.Context, requestID int64, status string, reviewedByUserID int64, reviewReason string) error
	HasPendingRequest(ctx context.Context, initiativeID int64) (bool, error)
}
