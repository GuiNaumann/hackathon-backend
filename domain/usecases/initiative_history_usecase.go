package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type InitiativeHistoryUseCase interface {
	GetHistory(ctx context.Context, initiativeID int64) ([]*entities.InitiativeHistoryResponse, error)
}
