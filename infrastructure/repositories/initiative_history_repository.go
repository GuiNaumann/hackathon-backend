package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type InitiativeHistoryRepository interface {
	Create(ctx context.Context, history *entities.InitiativeHistory) error
	ListByInitiative(ctx context.Context, initiativeID int64) ([]*entities.InitiativeHistory, error)
	GetLatestStatus(ctx context.Context, initiativeID int64) (*entities.InitiativeHistory, error)
}
