package repositories

import (
	"context"
	"hackathon-backend/domain/entities"
)

type SectorRepository interface {
	Create(ctx context.Context, sector *entities.Sector) error
	Update(ctx context.Context, sector *entities.Sector) error
	Delete(ctx context.Context, sectorID int64) error
	GetByID(ctx context.Context, sectorID int64) (*entities.Sector, error)
	GetByName(ctx context.Context, name string) (*entities.Sector, error)
	ListAll(ctx context.Context, activeOnly bool) ([]*entities.Sector, error)
	ListWithUserCount(ctx context.Context, activeOnly bool) ([]*entities.SectorListResponse, error)
	CountUsersBySector(ctx context.Context, sectorID int64) (int, error)
}
