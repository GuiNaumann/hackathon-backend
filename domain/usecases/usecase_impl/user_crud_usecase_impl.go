package usecase_impl

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
	"regexp"
)

type UserCrudUseCaseImpl struct {
	authRepo repositories.AuthRepository
	permRepo repositories.PermissionRepository
}

func NewUserCrudUseCaseImpl(authRepo repositories.AuthRepository, permRepo repositories.PermissionRepository) *UserCrudUseCaseImpl {
	return &UserCrudUseCaseImpl{
		authRepo: authRepo,
		permRepo: permRepo,
	}
}

func (uc *UserCrudUseCaseImpl) CreateUser(ctx context.Context, req *entities.CreateUserRequest) (*entities.User, error) {
	// Validar email
	if !isValidEmail(req.Email) {
		return nil, errors.New("email inválido")
	}

	// Validar nome
	if len(req.Name) < 3 {
		return nil, errors.New("nome deve ter no mínimo 3 caracteres")
	}

	// Validar senha
	if len(req.Password) < 6 {
		return nil, errors.New("senha deve ter no mínimo 6 caracteres")
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao criptografar senha: %w", err)
	}

	// Criar usuário
	user := &entities.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: string(hashedPassword),
	}

	if err := uc.authRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	// Atribuir tipos ao usuário
	for _, typeID := range req.TypeIDs {
		if err := uc.permRepo.AssignUserType(ctx, user.ID, typeID); err != nil {
			return nil, fmt.Errorf("erro ao atribuir tipo:  %w", err)
		}
	}

	return user, nil
}

func (uc *UserCrudUseCaseImpl) UpdateUser(ctx context.Context, userID int64, req *entities.UpdateUserRequest) (*entities.User, error) {
	// Buscar usuário existente
	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	// Atualizar campos se fornecidos
	if req.Email != nil {
		if !isValidEmail(*req.Email) {
			return nil, errors.New("email inválido")
		}
		user.Email = *req.Email
	}

	if req.Name != nil {
		if len(*req.Name) < 3 {
			return nil, errors.New("nome deve ter no mínimo 3 caracteres")
		}
		user.Name = *req.Name
	}

	// Atualizar no banco
	if err := uc.authRepo.UpdateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	// Atualizar tipos se fornecidos
	if req.TypeIDs != nil {
		// Remover todos os tipos atuais
		if err := uc.authRepo.RemoveAllUserTypes(ctx, userID); err != nil {
			return nil, fmt.Errorf("erro ao remover tipos: %w", err)
		}

		// Adicionar novos tipos
		for _, typeID := range *req.TypeIDs {
			if err := uc.permRepo.AssignUserType(ctx, userID, typeID); err != nil {
				return nil, fmt.Errorf("erro ao atribuir tipo: %w", err)
			}
		}
	}

	return user, nil
}

func (uc *UserCrudUseCaseImpl) DeleteUser(ctx context.Context, userID int64) error {
	// Verificar se usuário existe
	_, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Deletar usuário (cascade vai remover type_user automaticamente)
	if err := uc.authRepo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	return nil
}

func (uc *UserCrudUseCaseImpl) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	return uc.authRepo.GetUserByID(ctx, userID)
}

func (uc *UserCrudUseCaseImpl) ListUsers(ctx context.Context) ([]*entities.UserListResponse, error) {
	users, err := uc.authRepo.ListAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar usuários: %w", err)
	}

	// Buscar tipos de cada usuário
	var result []*entities.UserListResponse
	for _, user := range users {
		userTypes, err := uc.permRepo.GetUserTypes(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar tipos do usuário: %w", err)
		}

		result = append(result, &entities.UserListResponse{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			UserTypes: userTypes,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

func (uc *UserCrudUseCaseImpl) ChangePassword(ctx context.Context, userID int64, req *entities.ChangePasswordRequest) error {
	// Buscar usuário
	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	// Validar senha antiga
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("senha atual incorreta")
	}

	// Validar nova senha
	if len(req.NewPassword) < 6 {
		return errors.New("nova senha deve ter no mínimo 6 caracteres")
	}

	// Hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("erro ao criptografar senha: %w", err)
	}

	// Atualizar senha
	if err := uc.authRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return fmt.Errorf("erro ao atualizar senha: %w", err)
	}

	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
