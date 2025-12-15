package usecase_impl

import (
	"context"
	"errors"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
	"hackathon-backend/settings_loader"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCaseImpl struct {
	authRepo repositories.AuthRepository
	settings *settings_loader.SettingsLoader
}

func NewAuthUseCaseImpl(authRepo repositories.AuthRepository, settings *settings_loader.SettingsLoader) *AuthUseCaseImpl {
	return &AuthUseCaseImpl{
		authRepo: authRepo,
		settings: settings,
	}
}

func (uc *AuthUseCaseImpl) Login(ctx context.Context, email, password string) (*entities.User, string, error) {
	// Validar credenciais
	user, err := uc.ValidateCredentials(ctx, email, password)
	if err != nil {
		return nil, "", err
	}

	// Gerar token (neste caso, retornamos o ID como string)
	token := fmt.Sprintf("%d", user.ID)

	return user, token, nil
}

func (uc *AuthUseCaseImpl) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	return uc.authRepo.GetUserByID(ctx, userID)
}

func (uc *AuthUseCaseImpl) ValidateCredentials(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := uc.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	// Comparar senha com hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	return user, nil
}
