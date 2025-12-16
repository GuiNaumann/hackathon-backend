package usecase_impl

import (
	"context"
	"errors"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
)

type CancellationUseCaseImpl struct {
	cancellationRepo repositories.CancellationRepository
	initiativeRepo   repositories.InitiativeRepository
	permRepo         repositories.PermissionRepository
}

func NewCancellationUseCaseImpl(
	cancellationRepo repositories.CancellationRepository,
	initiativeRepo repositories.InitiativeRepository,
	permRepo repositories.PermissionRepository,
) *CancellationUseCaseImpl {
	return &CancellationUseCaseImpl{
		cancellationRepo: cancellationRepo,
		initiativeRepo:   initiativeRepo,
		permRepo:         permRepo,
	}
}

func (uc *CancellationUseCaseImpl) RequestCancellation(ctx context.Context, initiativeID int64, req *entities.RequestCancellationRequest, userID int64) (*entities.InitiativeCancellationRequest, error) {
	// Validar razão
	if len(req.Reason) < 10 {
		return nil, errors.New("motivo deve ter no mínimo 10 caracteres")
	}

	// Verificar se iniciativa existe
	initiative, err := uc.initiativeRepo.GetByID(ctx, initiativeID)
	if err != nil {
		return nil, errors.New("iniciativa não encontrada")
	}

	// Verificar se já está cancelada
	if initiative.Status == entities.StatusCancelled {
		return nil, errors.New("iniciativa já está cancelada")
	}

	// Verificar se já tem solicitação pendente
	hasPending, err := uc.cancellationRepo.HasPendingRequest(ctx, initiativeID)
	if err == nil && hasPending {
		return nil, errors.New("já existe uma solicitação de cancelamento pendente para esta iniciativa")
	}

	// Criar solicitação
	cancellationReq := &entities.InitiativeCancellationRequest{
		InitiativeID:      initiativeID,
		RequestedByUserID: userID,
		Reason:            req.Reason,
		Status:            entities.CancellationStatusPending,
	}

	if err := uc.cancellationRepo.Create(ctx, cancellationReq); err != nil {
		return nil, fmt.Errorf("erro ao criar solicitação de cancelamento: %w", err)
	}

	return uc.cancellationRepo.GetByID(ctx, cancellationReq.ID)
}

func (uc *CancellationUseCaseImpl) ReviewCancellation(ctx context.Context, requestID int64, req *entities.ReviewCancellationRequest, userID int64) error {
	// Verificar se é admin ou manager
	isAdminOrManager, err := uc.isAdminOrManager(ctx, userID)
	if err != nil || !isAdminOrManager {
		return errors.New("apenas administradores e gerentes podem revisar solicitações de cancelamento")
	}

	// Validar razão
	if len(req.Reason) < 5 {
		return errors.New("justificativa deve ter no mínimo 5 caracteres")
	}

	// Buscar solicitação
	cancellationReq, err := uc.cancellationRepo.GetByID(ctx, requestID)
	if err != nil {
		return errors.New("solicitação de cancelamento não encontrada")
	}

	// Verificar se já foi revisada
	if cancellationReq.Status != entities.CancellationStatusPending {
		return errors.New("esta solicitação já foi revisada")
	}

	// Definir novo status
	newStatus := entities.CancellationStatusRejected
	if req.Approved {
		newStatus = entities.CancellationStatusApproved
	}

	// Atualizar solicitação
	if err := uc.cancellationRepo.UpdateStatus(ctx, requestID, newStatus, userID, req.Reason); err != nil {
		return fmt.Errorf("erro ao atualizar solicitação:   %w", err)
	}

	// Se aprovado, cancelar a iniciativa
	if req.Approved {
		statusReq := &entities.ChangeInitiativeStatusRequest{
			Status: entities.StatusCancelled,
			Reason: fmt.Sprintf("Cancelamento aprovado:  %s", req.Reason),
		}

		if err := uc.initiativeRepo.ChangeStatusWithUser(ctx, cancellationReq.InitiativeID, statusReq.Status, statusReq.Reason, userID); err != nil {
			return fmt.Errorf("erro ao cancelar iniciativa:  %w", err)
		}
	}

	return nil
}

func (uc *CancellationUseCaseImpl) ListPendingCancellations(ctx context.Context, userID int64) ([]*entities.CancellationRequestResponse, error) {
	// Verificar se é admin ou manager
	isAdminOrManager, err := uc.isAdminOrManager(ctx, userID)
	if err != nil || !isAdminOrManager {
		return nil, errors.New("apenas administradores e gerentes podem visualizar solicitações pendentes")
	}

	requests, err := uc.cancellationRepo.ListPending(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar solicitações:  %w", err)
	}

	var response []*entities.CancellationRequestResponse
	for _, req := range requests {
		// Buscar título da iniciativa
		initiative, err := uc.initiativeRepo.GetByID(ctx, req.InitiativeID)
		if err != nil {
			continue
		}

		response = append(response, &entities.CancellationRequestResponse{
			ID:                req.ID,
			InitiativeID:      req.InitiativeID,
			InitiativeTitle:   initiative.Title,
			RequestedByUserID: req.RequestedByUserID,
			RequestedByName:   req.RequestedByName,
			Reason:            req.Reason,
			Status:            req.Status,
			CreatedAt:         req.CreatedAt.Format("2006-01-02 15:04:05"),
			TimeAgo:           timeAgo(req.CreatedAt),
		})
	}

	return response, nil
}

func (uc *CancellationUseCaseImpl) GetCancellationRequest(ctx context.Context, requestID int64) (*entities.InitiativeCancellationRequest, error) {
	return uc.cancellationRepo.GetByID(ctx, requestID)
}

func (uc *CancellationUseCaseImpl) isAdminOrManager(ctx context.Context, userID int64) (bool, error) {
	userTypes, err := uc.permRepo.GetUserTypes(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, userType := range userTypes {
		if userType.Name == "admin" || userType.Name == "manager" {
			return true, nil
		}
	}

	return false, nil
}
