package usecase_impl

import (
	"context"
	"errors"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
	"time"
)

type InitiativeUseCaseImpl struct {
	initiativeRepo repositories.InitiativeRepository
	historyRepo    repositories.InitiativeHistoryRepository
	permRepo       repositories.PermissionRepository
}

func NewInitiativeUseCaseImpl(
	initiativeRepo repositories.InitiativeRepository,
	historyRepo repositories.InitiativeHistoryRepository,
	permRepo repositories.PermissionRepository,
) *InitiativeUseCaseImpl {
	return &InitiativeUseCaseImpl{
		initiativeRepo: initiativeRepo,
		historyRepo:    historyRepo,
		permRepo:       permRepo,
	}
}

func (uc *InitiativeUseCaseImpl) CreateInitiative(ctx context.Context, req *entities.CreateInitiativeRequest, ownerID int64) (*entities.Initiative, error) {
	// Validações
	if len(req.Title) < 5 {
		return nil, errors.New("título deve ter no mínimo 5 caracteres")
	}

	if len(req.Description) < 20 {
		return nil, errors.New("descrição deve ter no mínimo 20 caracteres")
	}

	if len(req.Benefits) < 10 {
		return nil, errors.New("benefícios devem ter no mínimo 10 caracteres")
	}

	// Validar tipo
	if !isValidType(req.Type) {
		return nil, errors.New("tipo de iniciativa inválido")
	}

	// Validar prioridade
	if !isValidPriority(req.Priority) {
		return nil, errors.New("prioridade inválida")
	}

	// Parse deadline se fornecido
	var deadline *time.Time
	if req.Deadline != nil && *req.Deadline != "" {
		t, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			return nil, errors.New("formato de data inválido, use YYYY-MM-DD")
		}
		deadline = &t
	}

	initiative := &entities.Initiative{
		Title:       req.Title,
		Description: req.Description,
		Benefits:    req.Benefits,
		Status:      entities.StatusSubmitted, // Status inicial:  Submetida
		Type:        req.Type,
		Priority:    req.Priority,
		Sector:      req.Sector,
		OwnerID:     ownerID,
		Deadline:    deadline,
	}

	if err := uc.initiativeRepo.Create(ctx, initiative); err != nil {
		return nil, fmt.Errorf("erro ao criar iniciativa: %w", err)
	}

	return initiative, nil
}

// Listar iniciativas submetidas (aguardando aprovação)
func (uc *InitiativeUseCaseImpl) ListSubmittedInitiatives(ctx context.Context, userID int64) ([]*entities.InitiativeListResponse, error) {
	// Verificar se é admin ou manager
	isAdminOrManager, err := uc.isAdminOrManager(ctx, userID)
	if err != nil || !isAdminOrManager {
		return nil, errors.New("apenas administradores e gerentes podem visualizar iniciativas submetidas")
	}

	// Filtro para buscar apenas iniciativas "Submetidas"
	filter := &entities.InitiativeFilter{
		Status: entities.StatusSubmitted,
	}

	initiatives, err := uc.initiativeRepo.ListAllWithCancellation(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar iniciativas submetidas:  %w", err)
	}

	// Converter para response
	var response []*entities.InitiativeListResponse
	for _, initiative := range initiatives {
		listItem := &entities.InitiativeListResponse{
			ID:          initiative.ID,
			Title:       initiative.Title,
			Description: truncateDescription(initiative.Description, 150),
			Status:      initiative.Status,
			Type:        initiative.Type,
			Priority:    initiative.Priority,
			Sector:      initiative.Sector,
			OwnerName:   initiative.OwnerName,
			Date:        formatDate(initiative.CreatedAt),
		}

		if initiative.CancellationRequest != nil {
			listItem.CancellationRequest = initiative.CancellationRequest
		}

		response = append(response, listItem)
	}

	return response, nil
}

// CORRIGIDO: Aprovar ou reprovar iniciativa
func (uc *InitiativeUseCaseImpl) ReviewInitiative(ctx context.Context, initiativeID int64, req *entities.ReviewInitiativeRequest, userID int64) error {
	// Verificar se é admin ou manager
	isAdminOrManager, err := uc.isAdminOrManager(ctx, userID)
	if err != nil || !isAdminOrManager {
		return errors.New("apenas administradores e gerentes podem revisar iniciativas")
	}

	// Validar justificativa
	if len(req.Reason) < 10 {
		return errors.New("justificativa deve ter no mínimo 10 caracteres")
	}

	// Buscar iniciativa
	initiative, err := uc.initiativeRepo.GetByID(ctx, initiativeID)
	if err != nil {
		return errors.New("iniciativa não encontrada")
	}

	// Verificar se está no status "Submetida"
	if initiative.Status != entities.StatusSubmitted {
		return errors.New("apenas iniciativas submetidas podem ser revisadas")
	}

	var newStatus string
	var historyReason string

	if req.Approved {
		// CORRIGIDO: Aprovado → muda para "Aprovada" (não mais "Em Análise")
		newStatus = entities.StatusApproved
		historyReason = fmt.Sprintf("✅ Iniciativa aprovada:  %s", req.Reason)
	} else {
		// Reprovado:  muda para "Reprovada"
		newStatus = entities.StatusRejected
		historyReason = fmt.Sprintf("❌ Iniciativa reprovada:  %s", req.Reason)
	}

	// Atualizar status e registrar no histórico
	if err := uc.initiativeRepo.ChangeStatusWithUser(ctx, initiativeID, newStatus, historyReason, userID); err != nil {
		return fmt.Errorf("erro ao revisar iniciativa: %w", err)
	}

	return nil
}

func (uc *InitiativeUseCaseImpl) UpdateInitiative(ctx context.Context, initiativeID int64, req *entities.UpdateInitiativeRequest, userID int64) (*entities.Initiative, error) {
	// Buscar iniciativa
	initiative, err := uc.initiativeRepo.GetByID(ctx, initiativeID)
	if err != nil {
		return nil, errors.New("iniciativa não encontrada")
	}

	// Verificar se é o dono ou admin
	isAdmin, _ := uc.isAdmin(ctx, userID)
	if initiative.OwnerID != userID && !isAdmin {
		return nil, errors.New("você não tem permissão para editar esta iniciativa")
	}

	// Atualizar campos
	if req.Title != nil {
		if len(*req.Title) < 5 {
			return nil, errors.New("título deve ter no mínimo 5 caracteres")
		}
		initiative.Title = *req.Title
	}

	if req.Description != nil {
		if len(*req.Description) < 20 {
			return nil, errors.New("descrição deve ter no mínimo 20 caracteres")
		}
		initiative.Description = *req.Description
	}

	if req.Benefits != nil {
		initiative.Benefits = *req.Benefits
	}

	if req.Type != nil {
		if !isValidType(*req.Type) {
			return nil, errors.New("tipo de iniciativa inválido")
		}
		initiative.Type = *req.Type
	}

	if req.Priority != nil {
		if !isValidPriority(*req.Priority) {
			return nil, errors.New("prioridade inválida")
		}
		initiative.Priority = *req.Priority
	}

	if req.Sector != nil {
		initiative.Sector = *req.Sector
	}

	if req.Deadline != nil {
		if *req.Deadline != "" {
			t, err := time.Parse("2006-01-02", *req.Deadline)
			if err != nil {
				return nil, errors.New("formato de data inválido")
			}
			initiative.Deadline = &t
		} else {
			initiative.Deadline = nil
		}
	}

	if err := uc.initiativeRepo.Update(ctx, initiative); err != nil {
		return nil, fmt.Errorf("erro ao atualizar iniciativa:  %w", err)
	}

	return initiative, nil
}

func (uc *InitiativeUseCaseImpl) DeleteInitiative(ctx context.Context, initiativeID int64, userID int64) error {
	// Buscar iniciativa
	initiative, err := uc.initiativeRepo.GetByID(ctx, initiativeID)
	if err != nil {
		return errors.New("iniciativa não encontrada")
	}

	// Verificar se é o dono ou admin
	isAdmin, _ := uc.isAdmin(ctx, userID)
	if initiative.OwnerID != userID && !isAdmin {
		return errors.New("você não tem permissão para deletar esta iniciativa")
	}

	return uc.initiativeRepo.Delete(ctx, initiativeID)
}

func (uc *InitiativeUseCaseImpl) GetInitiativeByID(ctx context.Context, initiativeID int64) (*entities.Initiative, error) {
	return uc.initiativeRepo.GetByIDWithCancellation(ctx, initiativeID)
}

func (uc *InitiativeUseCaseImpl) ListInitiatives(ctx context.Context, filter *entities.InitiativeFilter) ([]*entities.InitiativeListResponse, error) {
	initiatives, err := uc.initiativeRepo.ListAllWithCancellation(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar iniciativas: %w", err)
	}

	// Converter para response
	var response []*entities.InitiativeListResponse
	for _, initiative := range initiatives {
		listItem := &entities.InitiativeListResponse{
			ID:          initiative.ID,
			Title:       initiative.Title,
			Description: truncateDescription(initiative.Description, 150),
			Status:      initiative.Status,
			Type:        initiative.Type,
			Priority:    initiative.Priority,
			Sector:      initiative.Sector,
			OwnerName:   initiative.OwnerName,
			Date:        formatDate(initiative.CreatedAt),
		}

		if initiative.CancellationRequest != nil {
			listItem.CancellationRequest = initiative.CancellationRequest
		}

		response = append(response, listItem)
	}

	return response, nil
}

func (uc *InitiativeUseCaseImpl) ChangeStatus(ctx context.Context, initiativeID int64, req *entities.ChangeInitiativeStatusRequest, userID int64) error {
	// Apenas admin pode mudar status manualmente
	isAdmin, err := uc.isAdmin(ctx, userID)
	if err != nil || !isAdmin {
		return errors.New("apenas administradores podem alterar o status")
	}

	// Validar status
	if !isValidStatus(req.Status) {
		return errors.New("status inválido")
	}

	// Validar justificativa
	if len(req.Reason) < 5 {
		return errors.New("justificativa deve ter no mínimo 5 caracteres")
	}

	// Usar o método com userID para registrar corretamente no histórico
	return uc.initiativeRepo.ChangeStatusWithUser(ctx, initiativeID, req.Status, req.Reason, userID)
}

func (uc *InitiativeUseCaseImpl) GetMyInitiatives(ctx context.Context, userID int64) ([]*entities.Initiative, error) {
	return uc.initiativeRepo.GetByOwner(ctx, userID)
}

// Helper functions
func (uc *InitiativeUseCaseImpl) isAdmin(ctx context.Context, userID int64) (bool, error) {
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

func (uc *InitiativeUseCaseImpl) isAdminOrManager(ctx context.Context, userID int64) (bool, error) {
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

func isValidType(t string) bool {
	validTypes := []string{
		entities.TypeAutomation,
		entities.TypeIntegration,
		entities.TypeImprovement,
		entities.TypeNewProject,
	}

	for _, valid := range validTypes {
		if t == valid {
			return true
		}
	}
	return false
}

func isValidPriority(p string) bool {
	validPriorities := []string{
		entities.PriorityHigh,
		entities.PriorityMedium,
		entities.PriorityLow,
	}

	for _, valid := range validPriorities {
		if p == valid {
			return true
		}
	}
	return false
}

func isValidStatus(s string) bool {
	validStatuses := []string{
		entities.StatusSubmitted,
		entities.StatusInAnalysis,
		entities.StatusApproved,
		entities.StatusInExecution,
		entities.StatusReturned,
		entities.StatusRejected,
		entities.StatusCompleted,
		entities.StatusCancelled,
	}

	for _, valid := range validStatuses {
		if s == valid {
			return true
		}
	}
	return false
}

func truncateDescription(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func formatDate(t time.Time) string {
	months := map[time.Month]string{
		time.January:   "jan",
		time.February:  "fev",
		time.March:     "mar",
		time.April:     "abr",
		time.May:       "mai",
		time.June:      "jun",
		time.July:      "jul",
		time.August:    "ago",
		time.September: "set",
		time.October:   "out",
		time.November:  "nov",
		time.December:  "dez",
	}

	return fmt.Sprintf("%02d de %s, %d", t.Day(), months[t.Month()], t.Year())
}
