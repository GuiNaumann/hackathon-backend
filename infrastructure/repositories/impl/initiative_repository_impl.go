package repository_impl

import (
	"context"
	"database/sql"
	"fmt"
	"hackathon-backend/domain/entities"
	"strings"
)

type InitiativeRepositoryImpl struct {
	db          *sql.DB
	historyRepo *InitiativeHistoryRepositoryImpl
}

func NewInitiativeRepositoryImpl(db *sql.DB) *InitiativeRepositoryImpl {
	return &InitiativeRepositoryImpl{
		db:          db,
		historyRepo: NewInitiativeHistoryRepositoryImpl(db),
	}
}

func (r *InitiativeRepositoryImpl) Create(ctx context.Context, initiative *entities.Initiative) error {
	query := `
		INSERT INTO initiatives (title, description, benefits, status, type, priority, sector, owner_id, deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
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

	if err != nil {
		return err
	}

	// Registrar no histórico:  criação da iniciativa
	history := &entities.InitiativeHistory{
		InitiativeID: initiative.ID,
		UserID:       initiative.OwnerID,
		OldStatus:    "Rascunho",
		NewStatus:    initiative.Status,
		Reason:       "Iniciativa criada",
	}

	return r.historyRepo.Create(ctx, history)
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

// NOVO: GetByID com informação de cancelamento
func (r *InitiativeRepositoryImpl) GetByIDWithCancellation(ctx context.Context, initiativeID int64) (*entities.Initiative, error) {
	query := `
		SELECT i.id, i. title, i.description, i. benefits, i.status, i. type, i.priority, i. sector, 
		       i. owner_id, u.name as owner_name, i.deadline, i.created_at, i. updated_at,
		       cr.id, cr.status, cr.requested_by_user_id, u2.name, cr.reason, 
		       cr.reviewed_by_user_id, u3.name, cr.review_reason, cr.created_at, cr.reviewed_at
		FROM initiatives i
		INNER JOIN users u ON u.id = i. owner_id
		LEFT JOIN initiative_cancellation_requests cr ON cr.initiative_id = i.id 
		    AND cr.status IN ('Pendente', 'Aprovada', 'Reprovada')
		LEFT JOIN users u2 ON u2.id = cr.requested_by_user_id
		LEFT JOIN users u3 ON u3.id = cr.reviewed_by_user_id
		WHERE i.id = $1
		ORDER BY cr.created_at DESC
		LIMIT 1
	`

	initiative := &entities.Initiative{}
	var crID sql.NullInt64
	var crStatus, crRequestedByName, crReason sql.NullString
	var crRequestedByUserID sql.NullInt64
	var crReviewedByUserID sql.NullInt64
	var crReviewedByName, crReviewReason sql.NullString
	var crCreatedAt sql.NullTime
	var crReviewedAt sql.NullTime

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
		&crID,
		&crStatus,
		&crRequestedByUserID,
		&crRequestedByName,
		&crReason,
		&crReviewedByUserID,
		&crReviewedByName,
		&crReviewReason,
		&crCreatedAt,
		&crReviewedAt,
	)

	if err != nil {
		return nil, err
	}

	// Se existe solicitação de cancelamento
	if crID.Valid {
		cancellationInfo := &entities.InitiativeCancellationInfo{
			ID:                crID.Int64,
			Status:            crStatus.String,
			RequestedByUserID: crRequestedByUserID.Int64,
			RequestedByName:   crRequestedByName.String,
			Reason:            crReason.String,
			CreatedAt:         crCreatedAt.Time,
		}

		if crReviewedByUserID.Valid {
			cancellationInfo.ReviewedByUserID = &crReviewedByUserID.Int64
			cancellationInfo.ReviewedByName = crReviewedByName.String
			cancellationInfo.ReviewReason = crReviewReason.String
		}

		if crReviewedAt.Valid {
			cancellationInfo.ReviewedAt = &crReviewedAt.Time
		}

		initiative.CancellationRequest = cancellationInfo
	}

	return initiative, nil
}

func (r *InitiativeRepositoryImpl) ListAll(ctx context.Context, filter *entities.InitiativeFilter) ([]*entities.Initiative, error) {
	query := `
		SELECT i.id, i.title, i.description, i.benefits, i.status, i.type, i.priority, i.sector, 
		       i.owner_id, u. name as owner_name, i. deadline, i.created_at, i.updated_at
		FROM initiatives i
		INNER JOIN users u ON u.id = i. owner_id
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	if filter != nil {
		if filter.Search != "" {
			query += fmt.Sprintf(" AND (LOWER(i.title) LIKE $%d OR LOWER(i. description) LIKE $%d)", argCount, argCount)
			args = append(args, "%"+strings.ToLower(filter.Search)+"%")
			argCount++
		}

		if filter.Status != "" {
			query += fmt.Sprintf(" AND i. status = $%d", argCount)
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
			query += fmt.Sprintf(" AND i. priority = $%d", argCount)
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

// NOVO: ListAll com informação de cancelamento
// ATUALIZADO: ListAllWithCancellation com filtro por sector_id
func (r *InitiativeRepositoryImpl) ListAllWithCancellation(ctx context.Context, filter *entities.InitiativeFilter) ([]*entities.Initiative, error) {
	query := `
		SELECT i.id, i.title, i.description, i.benefits, i.status, i.type, i.priority, i.sector, 
		       i.owner_id, u.name as owner_name, i.deadline, i.created_at, i.updated_at,
		       cr.id, cr.status, cr.requested_by_user_id, u2.name, cr.reason, 
		       cr.reviewed_by_user_id, u3.name, cr.review_reason, cr.created_at, cr.reviewed_at
		FROM initiatives i
		INNER JOIN users u ON u.id = i.owner_id
		LEFT JOIN LATERAL (
		    SELECT * FROM initiative_cancellation_requests 
		    WHERE initiative_id = i.id 
		    AND status IN ('Pendente', 'Aprovada', 'Reprovada')
		    ORDER BY created_at DESC
		    LIMIT 1
		) cr ON true
		LEFT JOIN users u2 ON u2.id = cr.requested_by_user_id
		LEFT JOIN users u3 ON u3.id = cr.reviewed_by_user_id
		WHERE 1=1
	`

	var args []interface{}
	argCount := 1

	if filter != nil {
		if filter.Search != "" {
			query += fmt.Sprintf(" AND (LOWER(i.title) LIKE $%d OR LOWER(i. description) LIKE $%d)", argCount, argCount)
			args = append(args, "%"+strings.ToLower(filter.Search)+"%")
			argCount++
		}

		if filter.Status != "" {
			query += fmt.Sprintf(" AND i. status = $%d", argCount)
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

		// NOVO: Filtrar por sector_id dos owners (usuários do mesmo setor)
		if filter.SectorID != nil {
			query += fmt.Sprintf(" AND u.sector_id = $%d", argCount)
			args = append(args, *filter.SectorID)
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
		var crID sql.NullInt64
		var crStatus, crRequestedByName, crReason sql.NullString
		var crRequestedByUserID sql.NullInt64
		var crReviewedByUserID sql.NullInt64
		var crReviewedByName, crReviewReason sql.NullString
		var crCreatedAt sql.NullTime
		var crReviewedAt sql.NullTime

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
			&crID,
			&crStatus,
			&crRequestedByUserID,
			&crRequestedByName,
			&crReason,
			&crReviewedByUserID,
			&crReviewedByName,
			&crReviewReason,
			&crCreatedAt,
			&crReviewedAt,
		)
		if err != nil {
			return nil, err
		}

		// Se existe solicitação de cancelamento
		if crID.Valid {
			cancellationInfo := &entities.InitiativeCancellationInfo{
				ID:                crID.Int64,
				Status:            crStatus.String,
				RequestedByUserID: crRequestedByUserID.Int64,
				RequestedByName:   crRequestedByName.String,
				Reason:            crReason.String,
				CreatedAt:         crCreatedAt.Time,
			}

			if crReviewedByUserID.Valid {
				cancellationInfo.ReviewedByUserID = &crReviewedByUserID.Int64
				cancellationInfo.ReviewedByName = crReviewedByName.String
				cancellationInfo.ReviewReason = crReviewReason.String
			}

			if crReviewedAt.Valid {
				cancellationInfo.ReviewedAt = &crReviewedAt.Time
			}

			initiative.CancellationRequest = cancellationInfo
		}

		initiatives = append(initiatives, initiative)
	}

	return initiatives, nil
}

func (r *InitiativeRepositoryImpl) ChangeStatus(ctx context.Context, initiativeID int64, status, reason string) error {
	// Primeiro, buscar o status atual
	var oldStatus string
	err := r.db.QueryRowContext(ctx, `SELECT status FROM initiatives WHERE id = $1`, initiativeID).Scan(&oldStatus)
	if err != nil {
		return err
	}

	// Atualizar o status
	query := `
		UPDATE initiatives
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = r.db.ExecContext(ctx, query, status, initiativeID)
	if err != nil {
		return err
	}

	// Registrar no histórico (precisamos do userID, vamos passar pelo contexto depois)
	// Por enquanto vamos usar o owner_id da iniciativa
	var ownerID int64
	r.db.QueryRowContext(ctx, `SELECT owner_id FROM initiatives WHERE id = $1`, initiativeID).Scan(&ownerID)

	history := &entities.InitiativeHistory{
		InitiativeID: initiativeID,
		UserID:       ownerID, // Isso será corrigido no UseCase
		OldStatus:    oldStatus,
		NewStatus:    status,
		Reason:       reason,
	}

	return r.historyRepo.Create(ctx, history)
}

func (r *InitiativeRepositoryImpl) ChangeStatusWithUser(ctx context.Context, initiativeID int64, status, reason string, userID int64) error {
	// Primeiro, buscar o status atual
	var oldStatus string
	err := r.db.QueryRowContext(ctx, `SELECT status FROM initiatives WHERE id = $1`, initiativeID).Scan(&oldStatus)
	if err != nil {
		return err
	}

	// Atualizar o status
	query := `
		UPDATE initiatives
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err = r.db.ExecContext(ctx, query, status, initiativeID)
	if err != nil {
		return err
	}

	// Registrar no histórico
	history := &entities.InitiativeHistory{
		InitiativeID: initiativeID,
		UserID:       userID,
		OldStatus:    oldStatus,
		NewStatus:    status,
		Reason:       reason,
	}

	return r.historyRepo.Create(ctx, history)
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
