package repository_impl

import (
	"context"
	"database/sql"
	"hackathon-backend/domain/entities"
)

type CommentRepositoryImpl struct {
	db *sql.DB
}

func NewCommentRepositoryImpl(db *sql.DB) *CommentRepositoryImpl {
	return &CommentRepositoryImpl{db: db}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, comment *entities.Comment) error {
	query := `
		INSERT INTO initiative_comments (initiative_id, user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query,
		comment.InitiativeID,
		comment.UserID,
		comment.Content,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
}

func (r *CommentRepositoryImpl) Update(ctx context.Context, comment *entities.Comment) error {
	query := `
		UPDATE initiative_comments
		SET content = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at
	`

	return r.db.QueryRowContext(ctx, query, comment.Content, comment.ID).Scan(&comment.UpdatedAt)
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, commentID int64) error {
	query := `DELETE FROM initiative_comments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, commentID)
	return err
}

func (r *CommentRepositoryImpl) GetByID(ctx context.Context, commentID int64) (*entities.Comment, error) {
	query := `
		SELECT c.id, c.initiative_id, c.user_id, u.name as user_name, c.content, c.created_at, c.updated_at
		FROM initiative_comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.id = $1
	`

	comment := &entities.Comment{}
	err := r.db.QueryRowContext(ctx, query, commentID).Scan(
		&comment.ID,
		&comment.InitiativeID,
		&comment.UserID,
		&comment.UserName,
		&comment.Content,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *CommentRepositoryImpl) ListByInitiative(ctx context.Context, initiativeID int64) ([]*entities.Comment, error) {
	query := `
		SELECT c.id, c.initiative_id, c.user_id, u.name as user_name, c.content, c.created_at, c.updated_at
		FROM initiative_comments c
		INNER JOIN users u ON u.id = c.user_id
		WHERE c.initiative_id = $1
		ORDER BY c.created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, initiativeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*entities.Comment
	for rows.Next() {
		comment := &entities.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.InitiativeID,
			&comment.UserID,
			&comment.UserName,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (r *CommentRepositoryImpl) CountByInitiative(ctx context.Context, initiativeID int64) (int, error) {
	query := `SELECT COUNT(*) FROM initiative_comments WHERE initiative_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, initiativeID).Scan(&count)
	return count, err
}
