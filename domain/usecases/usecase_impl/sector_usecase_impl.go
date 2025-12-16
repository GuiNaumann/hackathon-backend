package usecase_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
	"strings"
)

type SectorUseCaseImpl struct {
	sectorRepo repositories.SectorRepository
}

func NewSectorUseCaseImpl(sectorRepo repositories.SectorRepository) *SectorUseCaseImpl {
	return &SectorUseCaseImpl{
		sectorRepo: sectorRepo,
	}
}

func (uc *SectorUseCaseImpl) CreateSector(ctx context.Context, req *entities.CreateSectorRequest) (*entities.Sector, error) {
	// Validações
	if len(strings.TrimSpace(req.Name)) < 3 {
		return nil, errors.New("nome do setor deve ter no mínimo 3 caracteres")
	}

	// Verificar se já existe
	existing, _ := uc.sectorRepo.GetByName(ctx, req.Name)
	if existing != nil {
		return nil, errors.New("já existe um setor com este nome")
	}

	sector := &entities.Sector{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Active:      true,
	}

	if err := uc.sectorRepo.Create(ctx, sector); err != nil {
		return nil, fmt.Errorf("erro ao criar setor: %w", err)
	}

	return sector, nil
}

func (uc *SectorUseCaseImpl) UpdateSector(ctx context.Context, sectorID int64, req *entities.UpdateSectorRequest) (*entities.Sector, error) {
	// Buscar setor existente
	sector, err := uc.sectorRepo.GetByID(ctx, sectorID)
	if err != nil {
		return nil, errors.New("setor não encontrado")
	}

	// Atualizar campos se fornecidos
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 3 {
			return nil, errors.New("nome do setor deve ter no mínimo 3 caracteres")
		}

		// Verificar se já existe outro setor com este nome
		existing, _ := uc.sectorRepo.GetByName(ctx, name)
		if existing != nil && existing.ID != sectorID {
			return nil, errors.New("já existe um setor com este nome")
		}

		sector.Name = name
	}

	if req.Description != nil {
		sector.Description = strings.TrimSpace(*req.Description)
	}

	if req.Active != nil {
		sector.Active = *req.Active
	}

	if err := uc.sectorRepo.Update(ctx, sector); err != nil {
		return nil, fmt.Errorf("erro ao atualizar setor: %w", err)
	}

	return sector, nil
}

func (uc *SectorUseCaseImpl) DeleteSector(ctx context.Context, sectorID int64) error {
	// Verificar se setor existe
	_, err := uc.sectorRepo.GetByID(ctx, sectorID)
	if err != nil {
		return errors.New("setor não encontrado")
	}

	// Verificar se tem usuários vinculados
	userCount, err := uc.sectorRepo.CountUsersBySector(ctx, sectorID)
	if err != nil {
		return fmt.Errorf("erro ao verificar usuários: %w", err)
	}

	if userCount > 0 {
		return fmt.Errorf("não é possível deletar o setor pois existem %d usuário(s) vinculado(s). Remova os usuários primeiro.", userCount)
	}

	if err := uc.sectorRepo.Delete(ctx, sectorID); err != nil {
		return fmt.Errorf("erro ao deletar setor: %w", err)
	}

	return nil
}

func (uc *SectorUseCaseImpl) GetSectorByID(ctx context.Context, sectorID int64) (*entities.Sector, error) {
	sector, err := uc.sectorRepo.GetByID(ctx, sectorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("setor não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar setor: %w", err)
	}

	return sector, nil
}

func (uc *SectorUseCaseImpl) ListSectors(ctx context.Context, activeOnly bool) ([]*entities.SectorListResponse, error) {
	return uc.sectorRepo.ListWithUserCount(ctx, activeOnly)
}
