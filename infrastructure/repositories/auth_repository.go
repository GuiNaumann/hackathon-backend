package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
}
