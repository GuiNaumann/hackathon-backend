package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type UserCrudUseCase interface {
	CreateUser(ctx context.Context, req *entities.CreateUserRequest) (*entities.User, error)
	UpdateUser(ctx context.Context, userID int64, req *entities.UpdateUserRequest) (*entities.User, error)
	DeleteUser(ctx context.Context, userID int64) error
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
	ListUsers(ctx context.Context) ([]*entities.UserListResponse, error)
	ChangePassword(ctx context.Context, userID int64, req *entities.ChangePasswordRequest) error
}
