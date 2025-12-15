package usecases

import (
	"context"
	"hackathon-backend/domain/entities"
)

type CommentUseCase interface {
	CreateComment(ctx context.Context, initiativeID int64, req *entities.CreateCommentRequest, userID int64) (*entities.Comment, error)
	UpdateComment(ctx context.Context, commentID int64, req *entities.UpdateCommentRequest, userID int64) (*entities.Comment, error)
	DeleteComment(ctx context.Context, commentID int64, userID int64) error
	ListComments(ctx context.Context, initiativeID int64) ([]*entities.CommentListResponse, error)
}
