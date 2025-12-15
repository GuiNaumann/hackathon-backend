package repository_impl

import (
	"context"
	"database/sql"
	"fmt"
	"hackathon-backend/domain/entities"
	"strings"
)

type InitiativeRepositoryImpl struct {
	db *sql.DB
}

func NewInitiativeRepositoryImpl(db *sql.DB) *InitiativeRepositoryImpl {
	return &InitiativeRepositoryImpl{db: db}
}

func (r *InitiativeRepositoryImpl) Create(ctx context.Context, initiative *entities.Initiative) error {
	query := `
		INSERT INTO initiatives (title, description, benefits, status, type, priority, sector, owner_id, deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		initiative.Title,
		initiative.Description,
		initiative.Benefits,
		initiative.Status,
		initiative.Type,
		initiative.Priority,
		initiative.Sector,
		initiative.OwnerID,
		initiative.Deadline,
	).Scan(&initiative.ID, &initiative.CreatedAt, &initiative.UpdatedAt)
}

func (r *InitiativeRepositoryImpl) Update(ctx context.Context, initiative *entities.Initiative) error {
	query := `
		UPDATE initiatives
		SET title = $1, description = $2, benefits = $3, type = $4, priority = $5, sector = $6, deadline = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		initiative.Title,
		initiative.Description,
		initiative.Benefits,
		initiative.Type,
		initiative.Priority,
		initiative.Sector,
		initiative.Deadline,
		initiative.ID,
	).Scan(&initiative.UpdatedAt)
}

func (r *InitiativeRepositoryImpl) Delete(ctx context.Context, initiativeID int64) error {
	query := `DELETE FROM initiatives WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, initiativeID)
	return err
}

func (r *InitiativeRepositoryImpl) GetByID(ctx context.Context, initiativeID int64) (*entities.Initiative, error) {
	query := `
		SELECT i.id, i.title, i.description, i.benefits, i.status, i.type, i.priority, i.sector, 
		       i.owner_id, u.name as owner_name, i.deadline, i.created_at, i.updated_at
		FROM initiatives i
		INNER JOIN users u ON u.id = i.owner_id
		WHERE i.id = $1
	`

	initiative := &entities.Initiative{}
	err := r.db.QueryRowContext(ctx, query, initiativeID).Scan(
		&initiative.ID,
		&initiative.Title,
		&initiative.Description,
		&initiative.Benefits,
		&initiative.Status,
		&initiative.Type,
		&initiative.Priority,
		&initiative.Sector,
		&initiative.OwnerID,
		&initiative.OwnerName,
		&initiative.Deadline,
		&initiative.CreatedAt,
		&initiative.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return initiative, nil
}

func (r *InitiativeRepositoryImpl) ListAll(ctx context.Context, filter *entities.InitiativeFilter) ([]*entities.Initiative, error) {
	query := `
		SELECT i.id, i.title, i.description, i.benefits, i.status, i.type, i.priority, i.sector, 
		       i.owner_id, u.name as owner_name, i.deadline, i.created_at, i.updated_at
		FROM initiatives i
		INNER JOIN users u ON u.id = i.owner_id
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	// Aplicar filtros dinamicamente
	if filter != nil {
		if filter.Search != "" {
			query += fmt.Sprintf(" AND (LOWER(i.title) LIKE $%d OR LOWER(i. description) LIKE $%d)", argCount, argCount)
			args = append(args, "%"+strings.ToLower(filter.Search)+"%")
			argCount++
		}

		if filter.Status != "" {
			query += fmt.Sprintf(" AND i.status = $%d", argCount)
			args = append(args, filter.Status)
			argCount++
		}

		if filter.Type != "" {
			query += fmt.Sprintf(" AND i.type = $%d", argCount)
			args = append(args, filter.Type)
			argCount++
		}

		if filter.Sector != "" {
			query += fmt.Sprintf(" AND i.sector = $%d", argCount)
			args = append(args, filter.Sector)
			argCount++
		}

		if filter.Priority != "" {
			query += fmt.Sprintf(" AND i.priority = $%d", argCount)
			args = append(args, filter.Priority)
			argCount++
		}
	}

	query += " ORDER BY i.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var initiatives []*entities.Initiative
	for rows.Next() {
		initiative := &entities.Initiative{}
		err := rows.Scan(
			&initiative.ID,
			&initiative.Title,
			&initiative.Description,
			&initiative.Benefits,
			&initiative.Status,
			&initiative.Type,
			&initiative.Priority,
			&initiative.Sector,
			&initiative.OwnerID,
			&initiative.OwnerName,
			&initiative.Deadline,
			&initiative.CreatedAt,
			&initiative.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		initiatives = append(initiatives, initiative)
	}

	return initiatives, nil
}

func (r *InitiativeRepositoryImpl) ChangeStatus(ctx context.Context, initiativeID int64, status, reason string) error {
	query := `
		UPDATE initiatives
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, initiativeID)
	return err
}

func (r *InitiativeRepositoryImpl) GetByOwner(ctx context.Context, ownerID int64) ([]*entities.Initiative, error) {
	query := `
		SELECT i.id, i.title, i.description, i.benefits, i.status, i.type, i.priority, i.sector, 
		       i.owner_id, u.name as owner_name, i.deadline, i.created_at, i.updated_at
		FROM initiatives i
		INNER JOIN users u ON u.id = i.owner_id
		WHERE i.owner_id = $1
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var initiatives []*entities.Initiative
	for rows.Next() {
		initiative := &entities.Initiative{}
		err := rows.Scan(
			&initiative.ID,
			&initiative.Title,
			&initiative.Description,
			&initiative.Benefits,
			&initiative.Status,
			&initiative.Type,
			&initiative.Priority,
			&initiative.Sector,
			&initiative.OwnerID,
			&initiative.OwnerName,
			&initiative.Deadline,
			&initiative.CreatedAt,
			&initiative.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		initiatives = append(initiatives, initiative)
	}

	return initiatives, nil
}
