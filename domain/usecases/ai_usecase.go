package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type AIUseCase interface {
	ProcessText(ctx context.Context, req *entities.AIRequest, userID int64) (*entities.AIResponse, error)
}
