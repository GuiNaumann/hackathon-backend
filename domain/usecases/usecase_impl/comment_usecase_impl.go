package usecase_impl

import (
	"context"
	"errors"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
	"time"
)

type CommentUseCaseImpl struct {
	commentRepo    repositories.CommentRepository
	initiativeRepo repositories.InitiativeRepository
	permRepo       repositories.PermissionRepository
}

func NewCommentUseCaseImpl(
	commentRepo repositories.CommentRepository,
	initiativeRepo repositories.InitiativeRepository,
	permRepo repositories.PermissionRepository,
) *CommentUseCaseImpl {
	return &CommentUseCaseImpl{
		commentRepo:    commentRepo,
		initiativeRepo: initiativeRepo,
		permRepo:       permRepo,
	}
}

func (uc *CommentUseCaseImpl) CreateComment(ctx context.Context, initiativeID int64, req *entities.CreateCommentRequest, userID int64) (*entities.Comment, error) {
	// Validar conteúdo
	if len(req.Content) < 3 {
		return nil, errors.New("comentário deve ter no mínimo 3 caracteres")
	}

	if len(req.Content) > 1000 {
		return nil, errors.New("comentário deve ter no máximo 1000 caracteres")
	}

	// Verificar se a iniciativa existe
	_, err := uc.initiativeRepo.GetByID(ctx, initiativeID)
	if err != nil {
		return nil, errors.New("iniciativa não encontrada")
	}

	comment := &entities.Comment{
		InitiativeID: initiativeID,
		UserID:       userID,
		Content:      req.Content,
	}

	if err := uc.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("erro ao criar comentário: %w", err)
	}

	// Buscar o comentário completo com nome do usuário
	return uc.commentRepo.GetByID(ctx, comment.ID)
}

func (uc *CommentUseCaseImpl) UpdateComment(ctx context.Context, commentID int64, req *entities.UpdateCommentRequest, userID int64) (*entities.Comment, error) {
	// Validar conteúdo
	if len(req.Content) < 3 {
		return nil, errors.New("comentário deve ter no mínimo 3 caracteres")
	}

	if len(req.Content) > 1000 {
		return nil, errors.New("comentário deve ter no máximo 1000 caracteres")
	}

	// Buscar comentário
	comment, err := uc.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return nil, errors.New("comentário não encontrado")
	}

	// Verificar se é o dono ou admin
	isAdmin, _ := uc.isAdmin(ctx, userID)
	if comment.UserID != userID && !isAdmin {
		return nil, errors.New("você não tem permissão para editar este comentário")
	}

	comment.Content = req.Content

	if err := uc.commentRepo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("erro ao atualizar comentário: %w", err)
	}

	return uc.commentRepo.GetByID(ctx, commentID)
}

func (uc *CommentUseCaseImpl) DeleteComment(ctx context.Context, commentID int64, userID int64) error {
	// Buscar comentário
	comment, err := uc.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return errors.New("comentário não encontrado")
	}

	// Verificar se é o dono ou admin
	isAdmin, _ := uc.isAdmin(ctx, userID)
	if comment.UserID != userID && !isAdmin {
		return errors.New("você não tem permissão para deletar este comentário")
	}

	return uc.commentRepo.Delete(ctx, commentID)
}

func (uc *CommentUseCaseImpl) ListComments(ctx context.Context, initiativeID int64) ([]*entities.CommentListResponse, error) {
	comments, err := uc.commentRepo.ListByInitiative(ctx, initiativeID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar comentários: %w", err)
	}

	var response []*entities.CommentListResponse
	for _, comment := range comments {
		response = append(response, &entities.CommentListResponse{
			ID:           comment.ID,
			InitiativeID: comment.InitiativeID,
			UserID:       comment.UserID,
			UserName:     comment.UserName,
			Content:      comment.Content,
			CreatedAt:    formatCommentDate(comment.CreatedAt),
			UpdatedAt:    formatCommentDate(comment.UpdatedAt),
		})
	}

	return response, nil
}

func (uc *CommentUseCaseImpl) isAdmin(ctx context.Context, userID int64) (bool, error) {
	userTypes, err := uc.permRepo.GetUserTypes(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, userType := range userTypes {
		if userType.Name == "admin" {
			return true, nil
		}
	}

	return false, nil
}

func formatCommentDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
