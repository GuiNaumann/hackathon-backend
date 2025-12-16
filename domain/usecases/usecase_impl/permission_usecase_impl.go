package usecase_impl

import (
	"context"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
)

type PermissionUseCaseImpl struct {
	permRepo repositories.PermissionRepository
	authRepo repositories.AuthRepository
}

func NewPermissionUseCaseImpl(permRepo repositories.PermissionRepository, authRepo repositories.AuthRepository) *PermissionUseCaseImpl {
	return &PermissionUseCaseImpl{
		permRepo: permRepo,
		authRepo: authRepo,
	}
}

func (uc *PermissionUseCaseImpl) GetUserTypes(ctx context.Context, userID int64) ([]*entities.UserType, error) {
	return uc.permRepo.GetUserTypes(ctx, userID)
}

func (uc *PermissionUseCaseImpl) GetPersonalInformation(ctx context.Context, userID int64) (*entities.PersonalInformation, error) {
	// Buscar usuário
	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuário:  %w", err)
	}

	// Buscar tipos do usuário
	userTypes, err := uc.permRepo.GetUserTypes(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tipos de usuário: %w", err)
	}

	return &entities.PersonalInformation{
		User:      user,
		UserTypes: userTypes,
	}, nil
}

func (uc *PermissionUseCaseImpl) HasPermission(ctx context.Context, userID int64, endpoint, method string) (bool, error) {
	return uc.permRepo.HasPermission(ctx, userID, endpoint, method)
}

func (uc *PermissionUseCaseImpl) AssignUserType(ctx context.Context, userID, userTypeID int64) error {
	return uc.permRepo.AssignUserType(ctx, userID, userTypeID)
}

func (uc *PermissionUseCaseImpl) RemoveUserType(ctx context.Context, userID, userTypeID int64) error {
	return uc.permRepo.RemoveUserType(ctx, userID, userTypeID)
}

// NOVO: Listar todos os tipos disponíveis
func (uc *PermissionUseCaseImpl) GetAllUserTypes(ctx context.Context) ([]*entities.UserType, error) {
	return uc.permRepo.GetAllUserTypes(ctx)
}
