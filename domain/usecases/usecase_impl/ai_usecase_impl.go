package usecase_impl

import (
	"context"
	"errors"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
)

type AIUseCaseImpl struct {
	aiRepo repositories.AIRepository
}

func NewAIUseCaseImpl(aiRepo repositories.AIRepository) *AIUseCaseImpl {
	return &AIUseCaseImpl{
		aiRepo: aiRepo,
	}
}

func (uc *AIUseCaseImpl) ProcessText(ctx context.Context, req *entities.AIRequest, userID int64) (*entities.AIResponse, error) {
	// Validações
	if req.Text == "" {
		return nil, errors.New("texto não pode estar vazio")
	}

	if req.Prompt == "" {
		return nil, errors.New("prompt não pode estar vazio")
	}

	if len(req.Text) < 1 {
		return nil, errors.New("texto muito curto")
	}

	if len(req.Prompt) < 5 {
		return nil, errors.New("prompt muito curto (mínimo 5 caracteres)")
	}

	// Chamar IA
	return uc.aiRepo.ProcessText(ctx, req)
}
