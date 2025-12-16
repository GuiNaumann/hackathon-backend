package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type AIRepository interface {
	ProcessText(ctx context.Context, req *entities.AIRequest) (*entities.AIResponse, error)
}
