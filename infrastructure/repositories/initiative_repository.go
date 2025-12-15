package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type InitiativeRepository interface {
	Create(ctx context.Context, initiative *entities.Initiative) error
	Update(ctx context.Context, initiative *entities.Initiative) error
	Delete(ctx context.Context, initiativeID int64) error
	GetByID(ctx context.Context, initiativeID int64) (*entities.Initiative, error)
	ListAll(ctx context.Context, filter *entities.InitiativeFilter) ([]*entities.Initiative, error)
	ChangeStatus(ctx context.Context, initiativeID int64, status, reason string) error
	ChangeStatusWithUser(ctx context.Context, initiativeID int64, status, reason string, userID int64) error
	GetByOwner(ctx context.Context, ownerID int64) ([]*entities.Initiative, error)
}
