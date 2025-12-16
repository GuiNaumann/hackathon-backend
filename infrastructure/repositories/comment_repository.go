package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *entities.Comment) error
	Update(ctx context.Context, comment *entities.Comment) error
	Delete(ctx context.Context, commentID int64) error
	GetByID(ctx context.Context, commentID int64) (*entities.Comment, error)
	ListByInitiative(ctx context.Context, initiativeID int64) ([]*entities.Comment, error)
	CountByInitiative(ctx context.Context, initiativeID int64) (int, error)
}
