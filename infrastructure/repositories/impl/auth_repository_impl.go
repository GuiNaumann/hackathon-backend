package repository_impl

import (
	"context"
	"database/sql"
	"hackathon-backend/domain/entities"
)

type AuthRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepositoryImpl(db *sql.DB) *AuthRepositoryImpl {
	return &AuthRepositoryImpl{db: db}
}

func (r *AuthRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return user, nil
}

func (r *AuthRepositoryImpl) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	query := `
		SELECT id, email, name, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &entities.User{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *AuthRepositoryImpl) CreateUser(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (email, name, password, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query, user.Email, user.Name, user.Password).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}
