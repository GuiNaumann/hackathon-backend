package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type CancellationUseCase interface {
	RequestCancellation(ctx context.Context, initiativeID int64, req *entities.RequestCancellationRequest, userID int64) (*entities.InitiativeCancellationRequest, error)
	ReviewCancellation(ctx context.Context, requestID int64, req *entities.ReviewCancellationRequest, userID int64) error
	ListPendingCancellations(ctx context.Context, userID int64) ([]*entities.CancellationRequestResponse, error)
	GetCancellationRequest(ctx context.Context, requestID int64) (*entities.InitiativeCancellationRequest, error)
}
