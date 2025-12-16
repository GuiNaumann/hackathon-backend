package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type AIRepository interface {
	RefineText(ctx context.Context, req *entities.RefineTextRequest) (*entities.RefineTextResponse, error)
}
