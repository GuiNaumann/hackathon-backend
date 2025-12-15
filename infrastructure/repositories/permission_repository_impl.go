package repositories

import (
	"context"
	"database/sql"
	"hackathon-backend/domain/entities"
	"regexp"
	"strings"
)

type PermissionRepositoryImpl struct {
	db *sql.DB
}

func NewPermissionRepositoryImpl(db *sql.DB) *PermissionRepositoryImpl {
	return &PermissionRepositoryImpl{db: db}
}

func (r *PermissionRepositoryImpl) GetUserTypes(ctx context.Context, userID int64) ([]*entities.UserType, error) {
	query := `
		SELECT ut.id, ut.name, ut.description, ut.created_at, ut.updated_at
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
		SELECT utp.endpoint
		FROM user_type_permissions utp
		INNER JOIN type_user tu ON tu.user_type_id = utp. user_type_id
		WHERE tu.user_id = $1
		AND utp.method = $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, method)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Buscar todos os padrões de endpoint que o usuário tem permissão
	for rows.Next() {
		var pattern string
		if err := rows.Scan(&pattern); err != nil {
			return false, err
		}

		// Verificar se o endpoint atual bate com o padrão
		if matchesPattern(endpoint, pattern) {
			return true, nil
		}
	}

	return false, nil
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

// Helper: Verificar se uma rota real bate com um padrão
func matchesPattern(actualPath, pattern string) bool {
	// Converter o padrão para regex
	// Exemplo: /api/private/users/{id} → ^/api/private/users/[^/]+$
	regexPattern := regexp.QuoteMeta(pattern)
	regexPattern = strings.ReplaceAll(regexPattern, `\{`, `{`)
	regexPattern = strings.ReplaceAll(regexPattern, `\}`, `}`)
	regexPattern = regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(regexPattern, `[^/]+`)
	regexPattern = "^" + regexPattern + "$"

	regex := regexp.MustCompile(regexPattern)
	return regex.MatchString(actualPath)
}
