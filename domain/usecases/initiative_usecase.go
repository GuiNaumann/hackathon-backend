package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type InitiativeUseCase interface {
	CreateInitiative(ctx context.Context, req *entities.CreateInitiativeRequest, ownerID int64) (*entities.Initiative, error)
	UpdateInitiative(ctx context.Context, initiativeID int64, req *entities.UpdateInitiativeRequest, userID int64) (*entities.Initiative, error)
	DeleteInitiative(ctx context.Context, initiativeID int64, userID int64) error
	GetInitiativeByID(ctx context.Context, initiativeID int64) (*entities.Initiative, error)
	ListInitiatives(ctx context.Context, filter *entities.InitiativeFilter) ([]*entities.InitiativeListResponse, error)
	ListSubmittedInitiatives(ctx context.Context, userID int64) ([]*entities.InitiativeListResponse, error)              // NOVO
	ReviewInitiative(ctx context.Context, initiativeID int64, req *entities.ReviewInitiativeRequest, userID int64) error // NOVO
	ChangeStatus(ctx context.Context, initiativeID int64, req *entities.ChangeInitiativeStatusRequest, userID int64) error
	GetMyInitiatives(ctx context.Context, userID int64) ([]*entities.Initiative, error)
}
