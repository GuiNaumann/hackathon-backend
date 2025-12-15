package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type PermissionUseCase interface {
	GetUserTypes(ctx context.Context, userID int64) ([]*entities.UserType, error)
	GetPersonalInformation(ctx context.Context, userID int64) (*entities.PersonalInformation, error)
	HasPermission(ctx context.Context, userID int64, endpoint, method string) (bool, error)
	AssignUserType(ctx context.Context, userID, userTypeID int64) error
	RemoveUserType(ctx context.Context, userID, userTypeID int64) error
}
