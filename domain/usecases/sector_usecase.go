package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type SectorUseCase interface {
	CreateSector(ctx context.Context, req *entities.CreateSectorRequest) (*entities.Sector, error)
	UpdateSector(ctx context.Context, sectorID int64, req *entities.UpdateSectorRequest) (*entities.Sector, error)
	DeleteSector(ctx context.Context, sectorID int64) error
	GetSectorByID(ctx context.Context, sectorID int64) (*entities.Sector, error)
	ListSectors(ctx context.Context, activeOnly bool) ([]*entities.SectorListResponse, error)
}
