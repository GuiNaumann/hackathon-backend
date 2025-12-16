package repository_impl

import (
	"context"
	"database/sql"
	"hackathon-backend/domain/entities"
)

type CancellationRepositoryImpl struct {
	db *sql.DB
}

func NewCancellationRepositoryImpl(db *sql.DB) *CancellationRepositoryImpl {
	return &CancellationRepositoryImpl{db: db}
}

func (r *CancellationRepositoryImpl) Create(ctx context.Context, request *entities.InitiativeCancellationRequest) error {
	query := `
		INSERT INTO initiative_cancellation_requests (initiative_id, requested_by_user_id, reason, status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(ctx, query,
		request.InitiativeID,
		request.RequestedByUserID,
		request.Reason,
		entities.CancellationStatusPending,
	).Scan(&request.ID, &request.CreatedAt)
}

func (r *CancellationRepositoryImpl) GetByID(ctx context.Context, requestID int64) (*entities.InitiativeCancellationRequest, error) {
	query := `
		SELECT cr.id, cr.initiative_id, cr.requested_by_user_id, u1.name as requested_by_name,
		       cr.reason, cr.status, cr.reviewed_by_user_id, u2.name as reviewed_by_name,
		       cr.review_reason, cr.created_at, cr.reviewed_at
		FROM initiative_cancellation_requests cr
		INNER JOIN users u1 ON u1.id = cr.requested_by_user_id
		LEFT JOIN users u2 ON u2.id = cr.reviewed_by_user_id
		WHERE cr.id = $1
	`

	req := &entities.InitiativeCancellationRequest{}
	var reviewedByUserID sql.NullInt64
	var reviewedByName, reviewReason sql.NullString
	var reviewedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, requestID).Scan(
		&req.ID,
		&req.InitiativeID,
		&req.RequestedByUserID,
		&req.RequestedByName,
		&req.Reason,
		&req.Status,
		&reviewedByUserID,
		&reviewedByName,
		&reviewReason,
		&req.CreatedAt,
		&reviewedAt,
	)

	if err != nil {
		return nil, err
	}

	if reviewedByUserID.Valid {
		req.ReviewedByUserID = &reviewedByUserID.Int64
		req.ReviewedByName = reviewedByName.String
		req.ReviewReason = reviewReason.String
	}

	if reviewedAt.Valid {
		req.ReviewedAt = &reviewedAt.Time
	}

	return req, nil
}

func (r *CancellationRepositoryImpl) GetPendingByInitiative(ctx context.Context, initiativeID int64) (*entities.InitiativeCancellationRequest, error) {
	query := `
		SELECT cr.id, cr.initiative_id, cr.requested_by_user_id, u. name as requested_by_name,
		       cr.reason, cr.status, cr.created_at
		FROM initiative_cancellation_requests cr
		INNER JOIN users u ON u.id = cr.requested_by_user_id
		WHERE cr.initiative_id = $1 AND cr.status = $2
	`

	req := &entities.InitiativeCancellationRequest{}
	err := r.db.QueryRowContext(ctx, query, initiativeID, entities.CancellationStatusPending).Scan(
		&req.ID,
		&req.InitiativeID,
		&req.RequestedByUserID,
		&req.RequestedByName,
		&req.Reason,
		&req.Status,
		&req.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return req, nil
}

func (r *CancellationRepositoryImpl) ListPending(ctx context.Context) ([]*entities.InitiativeCancellationRequest, error) {
	query := `
		SELECT cr.id, cr.initiative_id, cr.requested_by_user_id, u.name as requested_by_name,
		       cr. reason, cr.status, cr. created_at
		FROM initiative_cancellation_requests cr
		INNER JOIN users u ON u.id = cr.requested_by_user_id
		WHERE cr.status = $1
		ORDER BY cr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, entities.CancellationStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*entities.InitiativeCancellationRequest
	for rows.Next() {
		req := &entities.InitiativeCancellationRequest{}
		err := rows.Scan(
			&req.ID,
			&req.InitiativeID,
			&req.RequestedByUserID,
			&req.RequestedByName,
			&req.Reason,
			&req.Status,
			&req.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}

	return requests, nil
}

func (r *CancellationRepositoryImpl) UpdateStatus(ctx context.Context, requestID int64, status string, reviewedByUserID int64, reviewReason string) error {
	query := `
		UPDATE initiative_cancellation_requests
		SET status = $1, reviewed_by_user_id = $2, review_reason = $3, reviewed_at = NOW()
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, status, reviewedByUserID, reviewReason, requestID)
	return err
}

func (r *CancellationRepositoryImpl) HasPendingRequest(ctx context.Context, initiativeID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM initiative_cancellation_requests WHERE initiative_id = $1 AND status = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, initiativeID, entities.CancellationStatusPending).Scan(&exists)
	return exists, err
}
