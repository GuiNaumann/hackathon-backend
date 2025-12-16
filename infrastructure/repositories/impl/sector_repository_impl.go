package repository_impl

import (
	"context"
	"database/sql"
	"hackathon-backend/domain/entities"
)

type SectorRepositoryImpl struct {
	db *sql.DB
}

func NewSectorRepositoryImpl(db *sql.DB) *SectorRepositoryImpl {
	return &SectorRepositoryImpl{db: db}
}

func (r *SectorRepositoryImpl) Create(ctx context.Context, sector *entities.Sector) error {
	query := `
		INSERT INTO sectors (name, description, active, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		sector.Name,
		sector.Description,
		sector.Active,
	).Scan(&sector.ID, &sector.CreatedAt, &sector.UpdatedAt)
}

func (r *SectorRepositoryImpl) Update(ctx context.Context, sector *entities.Sector) error {
	query := `
		UPDATE sectors
		SET name = $1, description = $2, active = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		sector.Name,
		sector.Description,
		sector.Active,
		sector.ID,
	).Scan(&sector.UpdatedAt)
}

func (r *SectorRepositoryImpl) Delete(ctx context.Context, sectorID int64) error {
	query := `DELETE FROM sectors WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, sectorID)
	return err
}

func (r *SectorRepositoryImpl) GetByID(ctx context.Context, sectorID int64) (*entities.Sector, error) {
	query := `
		SELECT id, name, description, active, created_at, updated_at
		FROM sectors
		WHERE id = $1
	`

	sector := &entities.Sector{}
	err := r.db.QueryRowContext(ctx, query, sectorID).Scan(
		&sector.ID,
		&sector.Name,
		&sector.Description,
		&sector.Active,
		&sector.CreatedAt,
		&sector.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return sector, nil
}

func (r *SectorRepositoryImpl) GetByName(ctx context.Context, name string) (*entities.Sector, error) {
	query := `
		SELECT id, name, description, active, created_at, updated_at
		FROM sectors
		WHERE name = $1
	`

	sector := &entities.Sector{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&sector.ID,
		&sector.Name,
		&sector.Description,
		&sector.Active,
		&sector.CreatedAt,
		&sector.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return sector, nil
}

func (r *SectorRepositoryImpl) ListAll(ctx context.Context, activeOnly bool) ([]*entities.Sector, error) {
	query := `
		SELECT id, name, description, active, created_at, updated_at
		FROM sectors
		WHERE 1=1
	`

	if activeOnly {
		query += " AND active = true"
	}

	query += " ORDER BY name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sectors []*entities.Sector
	for rows.Next() {
		sector := &entities.Sector{}
		err := rows.Scan(
			&sector.ID,
			&sector.Name,
			&sector.Description,
			&sector.Active,
			&sector.CreatedAt,
			&sector.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sectors = append(sectors, sector)
	}

	return sectors, nil
}

func (r *SectorRepositoryImpl) ListWithUserCount(ctx context.Context, activeOnly bool) ([]*entities.SectorListResponse, error) {
	query := `
		SELECT s.id, s.name, s.description, s.active, COUNT(u.id) as user_count
		FROM sectors s
		LEFT JOIN users u ON u.sector_id = s.id
		WHERE 1=1
	`

	if activeOnly {
		query += " AND s. active = true"
	}

	query += " GROUP BY s. id, s.name, s. description, s.active ORDER BY s.name ASC"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sectors []*entities.SectorListResponse
	for rows.Next() {
		sector := &entities.SectorListResponse{}
		err := rows.Scan(
			&sector.ID,
			&sector.Name,
			&sector.Description,
			&sector.Active,
			&sector.UserCount,
		)
		if err != nil {
			return nil, err
		}
		sectors = append(sectors, sector)
	}

	return sectors, nil
}

func (r *SectorRepositoryImpl) CountUsersBySector(ctx context.Context, sectorID int64) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE sector_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, sectorID).Scan(&count)
	return count, err
}
