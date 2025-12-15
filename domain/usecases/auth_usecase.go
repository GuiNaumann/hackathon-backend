package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type AuthUseCase interface {
	Login(ctx context.Context, email, password string) (*entities.User, string, error)
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
	ValidateCredentials(ctx context.Context, email, password string) (*entities.User, error)
}
