package repository_impl

import (
	"context"
	"database/sql"
	"hackathon-backend/domain/entities"
)

type InitiativeHistoryRepositoryImpl struct {
	db *sql.DB
}

func NewInitiativeHistoryRepositoryImpl(db *sql.DB) *InitiativeHistoryRepositoryImpl {
	return &InitiativeHistoryRepositoryImpl{db: db}
}

func (r *InitiativeHistoryRepositoryImpl) Create(ctx context.Context, history *entities.InitiativeHistory) error {
	query := `
		INSERT INTO initiative_history (initiative_id, user_id, old_status, new_status, reason, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at
	`

	return r.db.QueryRowContext(ctx, query,
		history.InitiativeID,
		history.UserID,
		history.OldStatus,
		history.NewStatus,
		history.Reason,
	).Scan(&history.ID, &history.CreatedAt)
}

func (r *InitiativeHistoryRepositoryImpl) ListByInitiative(ctx context.Context, initiativeID int64) ([]*entities.InitiativeHistory, error) {
	query := `
		SELECT ih.id, ih.initiative_id, ih.user_id, u.name as user_name, 
		       ih.old_status, ih.new_status, ih.reason, ih.created_at
		FROM initiative_history ih
		INNER JOIN users u ON u.id = ih.user_id
		WHERE ih.initiative_id = $1
		ORDER BY ih.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, initiativeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []*entities.InitiativeHistory
	for rows.Next() {
		history := &entities.InitiativeHistory{}
		var reason sql.NullString

		err := rows.Scan(
			&history.ID,
			&history.InitiativeID,
			&history.UserID,
			&history.UserName,
			&history.OldStatus,
			&history.NewStatus,
			&reason,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if reason.Valid {
			history.Reason = reason.String
		}

		histories = append(histories, history)
	}

	return histories, nil
}

func (r *InitiativeHistoryRepositoryImpl) GetLatestStatus(ctx context.Context, initiativeID int64) (*entities.InitiativeHistory, error) {
	query := `
		SELECT ih.id, ih.initiative_id, ih.user_id, u.name as user_name, 
		       ih.old_status, ih.new_status, ih.reason, ih.created_at
		FROM initiative_history ih
		INNER JOIN users u ON u.id = ih.user_id
		WHERE ih.initiative_id = $1
		ORDER BY ih.created_at DESC
		LIMIT 1
	`

	history := &entities.InitiativeHistory{}
	var reason sql.NullString

	err := r.db.QueryRowContext(ctx, query, initiativeID).Scan(
		&history.ID,
		&history.InitiativeID,
		&history.UserID,
		&history.UserName,
		&history.OldStatus,
		&history.NewStatus,
		&reason,
		&history.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	if reason.Valid {
		history.Reason = reason.String
	}

	return history, nil
}
