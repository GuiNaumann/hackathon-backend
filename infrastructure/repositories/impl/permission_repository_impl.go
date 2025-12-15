package repository_impl

import (
	"context"
	"database/sql"
	"hackathon-backend/domain/entities"
)

type PermissionRepositoryImpl struct {
	db *sql.DB
}

func NewPermissionRepositoryImpl(db *sql.DB) *PermissionRepositoryImpl {
	return &PermissionRepositoryImpl{db: db}
}

func (r *PermissionRepositoryImpl) GetUserTypes(ctx context.Context, userID int64) ([]*entities.UserType, error) {
	query := `
		SELECT ut.id, ut.name, ut.description, ut.created_at, ut. updated_at
		FROM user_type ut
		INNER JOIN type_user tu ON tu.user_type_id = ut.id
		WHERE tu.user_id = $1
		ORDER BY ut.name
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTypes []*entities.UserType
	for rows.Next() {
		ut := &entities.UserType{}
		err := rows.Scan(&ut.ID, &ut.Name, &ut.Description, &ut.CreatedAt, &ut.UpdatedAt)
		if err != nil {
			return nil, err
		}
		userTypes = append(userTypes, ut)
	}

	return userTypes, nil
}

func (r *PermissionRepositoryImpl) HasPermission(ctx context.Context, userID int64, endpoint, method string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_type_permissions utp
			INNER JOIN type_user tu ON tu.user_type_id = utp. user_type_id
			WHERE tu.user_id = $1
			AND utp.endpoint = $2
			AND utp.method = $3
		)
	`

	var hasPermission bool
	err := r.db.QueryRowContext(ctx, query, userID, endpoint, method).Scan(&hasPermission)
	if err != nil {
		return false, err
	}

	return hasPermission, nil
}

func (r *PermissionRepositoryImpl) AssignUserType(ctx context.Context, userID, userTypeID int64) error {
	query := `
		INSERT INTO type_user (user_id, user_type_id, created_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (user_id, user_type_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, userID, userTypeID)
	return err
}

func (r *PermissionRepositoryImpl) RemoveUserType(ctx context.Context, userID, userTypeID int64) error {
	query := `
		DELETE FROM type_user
		WHERE user_id = $1 AND user_type_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, userID, userTypeID)
	return err
}

func (r *PermissionRepositoryImpl) GetAllUserTypes(ctx context.Context) ([]*entities.UserType, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM user_type
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTypes []*entities.UserType
	for rows.Next() {
		ut := &entities.UserType{}
		err := rows.Scan(&ut.ID, &ut.Name, &ut.Description, &ut.CreatedAt, &ut.UpdatedAt)
		if err != nil {
			return nil, err
		}
		userTypes = append(userTypes, ut)
	}

	return userTypes, nil
}
