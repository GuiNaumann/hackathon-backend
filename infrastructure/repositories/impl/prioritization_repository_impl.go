package repository_impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"hackathon-backend/domain/entities"
)

type PrioritizationRepositoryImpl struct {
	db *sql.DB
}

func NewPrioritizationRepositoryImpl(db *sql.DB) *PrioritizationRepositoryImpl {
	return &PrioritizationRepositoryImpl{db: db}
}

// Create cria uma nova priorização
func (r *PrioritizationRepositoryImpl) Create(ctx context.Context, prioritization *entities.InitiativePrioritization) error {
	priorityOrderJSON, err := json.Marshal(prioritization.PriorityOrder)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO initiative_prioritization (sector_id, year, priority_order, is_locked, created_by_user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		prioritization.SectorID,
		prioritization.Year,
		priorityOrderJSON,
		prioritization.IsLocked,
		prioritization.CreatedByUserID,
	).Scan(&prioritization.ID, &prioritization.CreatedAt, &prioritization.UpdatedAt)
}

// Update atualiza uma priorização existente
func (r *PrioritizationRepositoryImpl) Update(ctx context.Context, prioritization *entities.InitiativePrioritization) error {
	priorityOrderJSON, err := json.Marshal(prioritization.PriorityOrder)
	if err != nil {
		return err
	}

	query := `
		UPDATE initiative_prioritization
		SET priority_order = $1, is_locked = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		priorityOrderJSON,
		prioritization.IsLocked,
		prioritization.ID,
	).Scan(&prioritization.UpdatedAt)
}

// GetBySectorAndYear busca priorização por setor e ano
func (r *PrioritizationRepositoryImpl) GetBySectorAndYear(ctx context.Context, sectorID int64, year int) (*entities.InitiativePrioritization, error) {
	query := `
		SELECT p.id, p.sector_id, s.name as sector_name, p.year, p.priority_order, p. is_locked, 
		       p.created_by_user_id, u.name as created_by_name, p.created_at, p.updated_at
		FROM initiative_prioritization p
		INNER JOIN sectors s ON s.id = p.sector_id
		INNER JOIN users u ON u.id = p.created_by_user_id
		WHERE p.sector_id = $1 AND p.year = $2
	`

	prioritization := &entities.InitiativePrioritization{}
	var priorityOrderJSON []byte

	err := r.db.QueryRowContext(ctx, query, sectorID, year).Scan(
		&prioritization.ID,
		&prioritization.SectorID,
		&prioritization.SectorName,
		&prioritization.Year,
		&priorityOrderJSON,
		&prioritization.IsLocked,
		&prioritization.CreatedByUserID,
		&prioritization.CreatedByName,
		&prioritization.CreatedAt,
		&prioritization.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(priorityOrderJSON, &prioritization.PriorityOrder); err != nil {
		return nil, err
	}

	return prioritization, nil
}

// GetAllByYear busca todas as priorizações de um ano
func (r *PrioritizationRepositoryImpl) GetAllByYear(ctx context.Context, year int) ([]*entities.InitiativePrioritization, error) {
	query := `
		SELECT p.id, p.sector_id, s.name as sector_name, p.year, p.priority_order, p.is_locked, 
		       p.created_by_user_id, u. name as created_by_name, p.created_at, p. updated_at
		FROM initiative_prioritization p
		INNER JOIN sectors s ON s.id = p.sector_id
		INNER JOIN users u ON u.id = p.created_by_user_id
		WHERE p.year = $1
		ORDER BY s.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prioritizations []*entities.InitiativePrioritization
	for rows.Next() {
		prioritization := &entities.InitiativePrioritization{}
		var priorityOrderJSON []byte

		err := rows.Scan(
			&prioritization.ID,
			&prioritization.SectorID,
			&prioritization.SectorName,
			&prioritization.Year,
			&priorityOrderJSON,
			&prioritization.IsLocked,
			&prioritization.CreatedByUserID,
			&prioritization.CreatedByName,
			&prioritization.CreatedAt,
			&prioritization.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(priorityOrderJSON, &prioritization.PriorityOrder); err != nil {
			return nil, err
		}

		prioritizations = append(prioritizations, prioritization)
	}

	return prioritizations, nil
}

// LockPrioritization bloqueia uma priorização
func (r *PrioritizationRepositoryImpl) LockPrioritization(ctx context.Context, prioritizationID int64) error {
	query := `UPDATE initiative_prioritization SET is_locked = true, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, prioritizationID)
	return err
}

// UnlockPrioritization desbloqueia uma priorização
func (r *PrioritizationRepositoryImpl) UnlockPrioritization(ctx context.Context, prioritizationID int64) error {
	query := `UPDATE initiative_prioritization SET is_locked = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, prioritizationID)
	return err
}

// CreateChangeRequest cria uma solicitação de mudança
func (r *PrioritizationRepositoryImpl) CreateChangeRequest(ctx context.Context, request *entities.PrioritizationChangeRequest) error {
	newOrderJSON, err := json.Marshal(request.NewPriorityOrder)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO prioritization_change_requests (prioritization_id, requested_by_user_id, new_priority_order, reason, status, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(ctx, query,
		request.PrioritizationID,
		request.RequestedByUserID,
		newOrderJSON,
		request.Reason,
		entities.PrioritizationChangeStatusPending,
	).Scan(&request.ID, &request.CreatedAt)
}

// GetChangeRequestByID busca uma solicitação de mudança por ID
func (r *PrioritizationRepositoryImpl) GetChangeRequestByID(ctx context.Context, requestID int64) (*entities.PrioritizationChangeRequest, error) {
	query := `
		SELECT cr.id, cr.prioritization_id, cr.requested_by_user_id, u1.name as requested_by_name,
		       cr.new_priority_order, cr.reason, cr.status, cr.reviewed_by_user_id, u2.name as reviewed_by_name,
		       cr.review_reason, cr.created_at, cr.reviewed_at
		FROM prioritization_change_requests cr
		INNER JOIN users u1 ON u1.id = cr.requested_by_user_id
		LEFT JOIN users u2 ON u2.id = cr.reviewed_by_user_id
		WHERE cr.id = $1
	`

	request := &entities.PrioritizationChangeRequest{}
	var newOrderJSON []byte
	var reviewedByUserID sql.NullInt64
	var reviewedByName, reviewReason sql.NullString
	var reviewedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&request.ID,
		&request.PrioritizationID,
		&request.RequestedByUserID,
		&request.RequestedByName,
		&newOrderJSON,
		&request.Reason,
		&request.Status,
		&reviewedByUserID,
		&reviewedByName,
		&reviewReason,
		&request.CreatedAt,
		&reviewedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(newOrderJSON, &request.NewPriorityOrder); err != nil {
		return nil, err
	}

	if reviewedByUserID.Valid {
		request.ReviewedByUserID = &reviewedByUserID.Int64
		request.ReviewedByName = reviewedByName.String
		request.ReviewReason = reviewReason.String
	}

	if reviewedAt.Valid {
		request.ReviewedAt = &reviewedAt.Time
	}

	return request, nil
}

// ListPendingChangeRequests lista todas as solicitações pendentes
func (r *PrioritizationRepositoryImpl) ListPendingChangeRequests(ctx context.Context) ([]*entities.PrioritizationChangeRequest, error) {
	query := `
		SELECT cr.id, cr.prioritization_id, cr. requested_by_user_id, u. name as requested_by_name,
		       cr.new_priority_order, cr.reason, cr.status, cr.created_at
		FROM prioritization_change_requests cr
		INNER JOIN users u ON u.id = cr.requested_by_user_id
		WHERE cr.status = $1
		ORDER BY cr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, entities.PrioritizationChangeStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*entities.PrioritizationChangeRequest
	for rows.Next() {
		request := &entities.PrioritizationChangeRequest{}
		var newOrderJSON []byte

		err := rows.Scan(
			&request.ID,
			&request.PrioritizationID,
			&request.RequestedByUserID,
			&request.RequestedByName,
			&newOrderJSON,
			&request.Reason,
			&request.Status,
			&request.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(newOrderJSON, &request.NewPriorityOrder); err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	return requests, nil
}

// UpdateChangeRequestStatus atualiza o status de uma solicitação
func (r *PrioritizationRepositoryImpl) UpdateChangeRequestStatus(ctx context.Context, requestID int64, status string, reviewedByUserID int64, reviewReason string) error {
	query := `
		UPDATE prioritization_change_requests
		SET status = $1, reviewed_by_user_id = $2, review_reason = $3, reviewed_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, status, reviewedByUserID, reviewReason, requestID)
	return err
}

// HasPendingChangeRequest verifica se existe solicitação pendente
func (r *PrioritizationRepositoryImpl) HasPendingChangeRequest(ctx context.Context, prioritizationID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM prioritization_change_requests WHERE prioritization_id = $1 AND status = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, prioritizationID, entities.PrioritizationChangeStatusPending).Scan(&exists)
	return exists, err
}
