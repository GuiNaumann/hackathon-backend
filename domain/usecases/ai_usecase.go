package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type AIUseCase interface {
	RefineText(ctx context.Context, req *entities.RefineTextRequest, userID int64) (*entities.RefineTextResponse, error)
}
