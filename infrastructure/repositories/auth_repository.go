package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type AuthRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, userID int64) error
	ListAllUsers(ctx context.Context) ([]*entities.User, error)
	UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error
	RemoveAllUserTypes(ctx context.Context, userID int64) error
}
