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

func (uc *AIUseCaseImpl) RefineText(ctx context.Context, req *entities.RefineTextRequest, userID int64) (*entities.RefineTextResponse, error) {
	// Validações
	if req.Text == "" {
		return nil, errors.New("texto não pode estar vazio")
	}

	if len(req.Text) < 10 {
		return nil, errors.New("texto muito curto (mínimo 10 caracteres)")
	}

	// Validar action
	validActions := []string{entities.ActionSummarize, entities.ActionRefine, entities.ActionExpand}
	isValidAction := false
	for _, validAction := range validActions {
		if req.Action == validAction {
			isValidAction = true
			break
		}
	}

	if !isValidAction {
		req.Action = entities.ActionRefine // Default
	}

	// Chamar IA
	return uc.aiRepo.RefineText(ctx, req)
}
