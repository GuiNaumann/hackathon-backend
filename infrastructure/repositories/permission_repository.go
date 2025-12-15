package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type PermissionRepository interface {
	GetUserTypes(ctx context.Context, userID int64) ([]*entities.UserType, error)
	HasPermission(ctx context.Context, userID int64, endpoint, method string) (bool, error)
	AssignUserType(ctx context.Context, userID, userTypeID int64) error
	RemoveUserType(ctx context.Context, userID, userTypeID int64) error
	GetAllUserTypes(ctx context.Context) ([]*entities.UserType, error)
}
