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
		SELECT u.id, u.email, u.name, u.password, u.sector_id, u.created_at, u.updated_at
		FROM users u
		WHERE u.email = $1
	`

	user := &entities.User{}
	var sectorID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&sectorID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if sectorID.Valid {
		user.SectorID = &sectorID.Int64
	}

	return user, nil
}

func (r *AuthRepositoryImpl) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.password, u.sector_id, u.created_at, u.updated_at
		FROM users u
		WHERE u.id = $1
	`

	user := &entities.User{}
	var sectorID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&sectorID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if sectorID.Valid {
		user.SectorID = &sectorID.Int64
	}

	return user, nil
}

func (r *AuthRepositoryImpl) CreateUser(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (email, name, password, sector_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query, user.Email, user.Name, user.Password, user.SectorID).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func (r *AuthRepositoryImpl) UpdateUser(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users
		SET email = $1, name = $2, sector_id = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	return r.db.QueryRowContext(ctx, query, user.Email, user.Name, user.SectorID, user.ID).Scan(&user.UpdatedAt)
}

func (r *AuthRepositoryImpl) DeleteUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// ATUALIZADO: Buscar usuários com informações do setor
func (r *AuthRepositoryImpl) ListAllUsers(ctx context.Context) ([]*entities.User, error) {
	query := `
		SELECT u.id, u.email, u.name, u.password, u.sector_id, s.name as sector_name, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN sectors s ON s.id = u.sector_id
		ORDER BY u.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		var sectorID sql.NullInt64
		var sectorName sql.NullString

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Password,
			&sectorID,
			&sectorName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if sectorID.Valid {
			user.SectorID = &sectorID.Int64
		}

		if sectorName.Valid {
			user.SectorName = sectorName.String
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *AuthRepositoryImpl) UpdatePassword(ctx context.Context, userID int64, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, hashedPassword, userID)
	return err
}

func (r *AuthRepositoryImpl) RemoveAllUserTypes(ctx context.Context, userID int64) error {
	query := `DELETE FROM type_user WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
