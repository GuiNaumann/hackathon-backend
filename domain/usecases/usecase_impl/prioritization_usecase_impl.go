package usecase_impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
	"time"
)

type PrioritizationUseCaseImpl struct {
	prioritizationRepo repositories.PrioritizationRepository
	initiativeRepo     repositories.InitiativeRepository
	authRepo           repositories.AuthRepository
	permRepo           repositories.PermissionRepository
	sectorRepo         repositories.SectorRepository // NOVO
}

func NewPrioritizationUseCaseImpl(
	prioritizationRepo repositories.PrioritizationRepository,
	initiativeRepo repositories.InitiativeRepository,
	authRepo repositories.AuthRepository,
	permRepo repositories.PermissionRepository,
	sectorRepo repositories.SectorRepository, // NOVO
) *PrioritizationUseCaseImpl {
	return &PrioritizationUseCaseImpl{
		prioritizationRepo: prioritizationRepo,
		initiativeRepo:     initiativeRepo,
		authRepo:           authRepo,
		permRepo:           permRepo,
		sectorRepo:         sectorRepo, // NOVO
	}
}

// SavePrioritization salva a priorização do setor do usuário
func (uc *PrioritizationUseCaseImpl) SavePrioritization(ctx context.Context, req *entities.SavePrioritizationRequest, userID int64) (*entities.PrioritizationWithInitiatives, error) {
	// Validações
	if req.Year < 2020 || req.Year > 2100 {
		return nil, errors.New("ano inválido")
	}

	if len(req.PriorityOrder) == 0 {
		return nil, errors.New("ordem de prioridade não pode estar vazia")
	}

	// Buscar usuário para pegar o setor
	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	if user.SectorID == nil {
		return nil, errors.New("usuário não está vinculado a um setor")
	}

	// Verificar se já existe priorização para este setor/ano
	existing, err := uc.prioritizationRepo.GetBySectorAndYear(ctx, *user.SectorID, req.Year)

	isAdminOrManager, _ := uc.isAdminOrManager(ctx, userID)

	if err == nil {
		// Já existe priorização
		if existing.IsLocked && !isAdminOrManager {
			return nil, errors.New("priorização já está bloqueada.  Solicite aprovação para alterá-la")
		}

		// Atualizar priorização existente
		existing.PriorityOrder = req.PriorityOrder
		existing.IsLocked = true // Bloquear ao salvar

		if err := uc.prioritizationRepo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("erro ao atualizar priorização: %w", err)
		}

		return uc.buildPrioritizationWithInitiatives(ctx, existing)
	}

	// Criar nova priorização
	prioritization := &entities.InitiativePrioritization{
		SectorID:        *user.SectorID,
		Year:            req.Year,
		PriorityOrder:   req.PriorityOrder,
		IsLocked:        true, // Bloquear ao salvar
		CreatedByUserID: userID,
	}

	if err := uc.prioritizationRepo.Create(ctx, prioritization); err != nil {
		return nil, fmt.Errorf("erro ao criar priorização: %w", err)
	}

	// Buscar novamente para pegar os dados completos
	created, err := uc.prioritizationRepo.GetBySectorAndYear(ctx, *user.SectorID, req.Year)
	if err != nil {
		return nil, err
	}

	return uc.buildPrioritizationWithInitiatives(ctx, created)
}

// GetPrioritization busca a priorização do setor do usuário
func (uc *PrioritizationUseCaseImpl) GetPrioritization(ctx context.Context, year int, userID int64) (*entities.PrioritizationWithInitiatives, error) {
	// Buscar usuário para pegar o setor
	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	if user.SectorID == nil {
		return nil, errors.New("usuário não está vinculado a um setor")
	}

	// Buscar priorização
	prioritization, err := uc.prioritizationRepo.GetBySectorAndYear(ctx, *user.SectorID, year)
	if err != nil {
		if err == sql.ErrNoRows {
			// Não existe priorização ainda, retornar iniciativas do setor sem ordem
			return uc.buildEmptyPrioritization(ctx, *user.SectorID, year, userID)
		}
		return nil, fmt.Errorf("erro ao buscar priorização:  %w", err)
	}

	return uc.buildPrioritizationWithInitiatives(ctx, prioritization)
}

// GetAllSectorsPrioritization busca priorização de todos os setores (Admin/Manager)
func (uc *PrioritizationUseCaseImpl) GetAllSectorsPrioritization(ctx context.Context, year int, userID int64) (*entities.AllSectorsPrioritization, error) {
	// Verificar se é admin ou manager
	isAdminOrManager, err := uc.isAdminOrManager(ctx, userID)
	if err != nil || !isAdminOrManager {
		return nil, errors.New("apenas administradores e gerentes podem visualizar todas as priorizações")
	}

	// Buscar todas as priorizações do ano
	prioritizations, err := uc.prioritizationRepo.GetAllByYear(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar priorizações:  %w", err)
	}

	// NOVO: Se não existir nenhuma priorização, buscar todos os setores ativos e criar estrutura vazia
	if len(prioritizations) == 0 {
		return uc.buildEmptyPrioritizationForAllSectors(ctx, year)
	}

	// Construir response com iniciativas de cada setor
	var sectors []*entities.PrioritizationWithInitiatives
	for _, p := range prioritizations {
		withInitiatives, err := uc.buildPrioritizationWithInitiatives(ctx, p)
		if err != nil {
			continue // Skip em caso de erro
		}
		sectors = append(sectors, withInitiatives)
	}

	return &entities.AllSectorsPrioritization{
		Year:    year,
		Sectors: sectors,
	}, nil
}

func (uc *PrioritizationUseCaseImpl) buildEmptyPrioritizationForAllSectors(ctx context.Context, year int) (*entities.AllSectorsPrioritization, error) {
	// Buscar todos os setores ativos
	sectors, err := uc.sectorRepo.ListAll(ctx, true) // true = apenas ativos
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar setores:   %w", err)
	}

	var sectorsWithInitiatives []*entities.PrioritizationWithInitiatives

	for _, sector := range sectors {
		// Buscar TODAS as iniciativas do setor (sem filtro de status)
		filter := &entities.InitiativeFilter{
			Sector: sector.Name, // Filtrar pelo nome do setor
		}

		initiatives, err := uc.initiativeRepo.ListAllWithCancellation(ctx, filter)
		if err != nil {
			continue // Skip em caso de erro
		}

		// Filtrar apenas os status desejados:  Aprovada, Em Execução, Em Análise
		var initiativesList []*entities.InitiativeListResponse
		for _, initiative := range initiatives {
			// FILTRO: Apenas status permitidos
			if initiative.Status == entities.StatusApproved ||
				initiative.Status == entities.StatusInExecution ||
				initiative.Status == entities.StatusInAnalysis {

				initiativesList = append(initiativesList, &entities.InitiativeListResponse{
					ID:          initiative.ID,
					Title:       initiative.Title,
					Description: initiative.Description,
					Status:      initiative.Status,
					Type:        initiative.Type,
					Priority:    initiative.Priority,
					Sector:      initiative.Sector,
					OwnerName:   initiative.OwnerName,
					Date:        formatDate(initiative.CreatedAt),
				})
			}
		}

		sectorsWithInitiatives = append(sectorsWithInitiatives, &entities.PrioritizationWithInitiatives{
			ID:          0, // Não existe ainda
			SectorID:    sector.ID,
			SectorName:  sector.Name,
			Year:        year,
			IsLocked:    false,
			Initiatives: initiativesList,
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	return &entities.AllSectorsPrioritization{
		Year:    year,
		Sectors: sectorsWithInitiatives,
	}, nil
}

// Atualizar a função buildEmptyPrioritization para usar sectorRepo
func (uc *PrioritizationUseCaseImpl) buildEmptyPrioritization(ctx context.Context, sectorID int64, year int, userID int64) (*entities.PrioritizationWithInitiatives, error) {
	// Buscar nome do setor
	sector, err := uc.sectorRepo.GetByID(ctx, sectorID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar setor:  %w", err)
	}

	// Buscar TODAS as iniciativas do setor
	filter := &entities.InitiativeFilter{
		Sector: sector.Name,
	}

	initiatives, err := uc.initiativeRepo.ListAllWithCancellation(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Filtrar apenas os status desejados
	var initiativesList []*entities.InitiativeListResponse
	for _, initiative := range initiatives {
		// FILTRO: Apenas status permitidos
		if initiative.Status == entities.StatusApproved ||
			initiative.Status == entities.StatusInExecution ||
			initiative.Status == entities.StatusInAnalysis {

			initiativesList = append(initiativesList, &entities.InitiativeListResponse{
				ID:          initiative.ID,
				Title:       initiative.Title,
				Description: initiative.Description,
				Status:      initiative.Status,
				Type:        initiative.Type,
				Priority:    initiative.Priority,
				Sector:      initiative.Sector,
				OwnerName:   initiative.OwnerName,
				Date:        formatDate(initiative.CreatedAt),
			})
		}
	}

	return &entities.PrioritizationWithInitiatives{
		ID:          0, // Não existe ainda
		SectorID:    sectorID,
		SectorName:  sector.Name,
		Year:        year,
		IsLocked:    false,
		Initiatives: initiativesList,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// RequestPrioritizationChange solicita mudança na priorização
func (uc *PrioritizationUseCaseImpl) RequestPrioritizationChange(ctx context.Context, req *entities.RequestPrioritizationChangeRequest, userID int64, year int) (*entities.PrioritizationChangeRequest, error) {
	// Validações
	if len(req.Reason) < 10 {
		return nil, errors.New("motivo deve ter no mínimo 10 caracteres")
	}

	if len(req.NewPriorityOrder) == 0 {
		return nil, errors.New("nova ordem de prioridade não pode estar vazia")
	}

	// Buscar usuário para pegar o setor
	user, err := uc.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("usuário não encontrado")
	}

	if user.SectorID == nil {
		return nil, errors.New("usuário não está vinculado a um setor")
	}

	// Buscar priorização existente
	prioritization, err := uc.prioritizationRepo.GetBySectorAndYear(ctx, *user.SectorID, year)
	if err != nil {
		return nil, errors.New("priorização não encontrada")
	}

	// Verificar se já tem solicitação pendente
	hasPending, err := uc.prioritizationRepo.HasPendingChangeRequest(ctx, prioritization.ID)
	if err == nil && hasPending {
		return nil, errors.New("já existe uma solicitação de mudança pendente para esta priorização")
	}

	// Criar solicitação
	changeRequest := &entities.PrioritizationChangeRequest{
		PrioritizationID:  prioritization.ID,
		RequestedByUserID: userID,
		NewPriorityOrder:  req.NewPriorityOrder,
		Reason:            req.Reason,
		Status:            entities.PrioritizationChangeStatusPending,
	}

	if err := uc.prioritizationRepo.CreateChangeRequest(ctx, changeRequest); err != nil {
		return nil, fmt.Errorf("erro ao criar solicitação: %w", err)
	}

	return uc.prioritizationRepo.GetChangeRequestByID(ctx, changeRequest.ID)
}

// ReviewPrioritizationChange aprova ou recusa solicitação de mudança
func (uc *PrioritizationUseCaseImpl) ReviewPrioritizationChange(ctx context.Context, requestID int64, req *entities.ReviewPrioritizationChangeRequest, userID int64) error {
	// Verificar se é admin ou manager
	isAdminOrManager, err := uc.isAdminOrManager(ctx, userID)
	if err != nil || !isAdminOrManager {
		return errors.New("apenas administradores e gerentes podem revisar solicitações")
	}

	// Validar motivo
	if len(req.Reason) < 5 {
		return errors.New("justificativa deve ter no mínimo 5 caracteres")
	}

	// Buscar solicitação
	changeRequest, err := uc.prioritizationRepo.GetChangeRequestByID(ctx, requestID)
	if err != nil {
		return errors.New("solicitação não encontrada")
	}

	// Verificar se já foi revisada
	if changeRequest.Status != entities.PrioritizationChangeStatusPending {
		return errors.New("esta solicitação já foi revisada")
	}

	// Definir novo status
	newStatus := entities.PrioritizationChangeStatusRejected
	if req.Approved {
		newStatus = entities.PrioritizationChangeStatusApproved
	}

	// Atualizar status da solicitação
	if err := uc.prioritizationRepo.UpdateChangeRequestStatus(ctx, requestID, newStatus, userID, req.Reason); err != nil {
		return fmt.Errorf("erro ao atualizar solicitação: %w", err)
	}

	// Se aprovado, desbloquear a priorização para permitir alteração
	if req.Approved {
		if err := uc.prioritizationRepo.UnlockPrioritization(ctx, changeRequest.PrioritizationID); err != nil {
			return fmt.Errorf("erro ao desbloquear priorização:  %w", err)
		}
	}

	return nil
}

// ListPendingChangeRequests lista solicitações pendentes
func (uc *PrioritizationUseCaseImpl) ListPendingChangeRequests(ctx context.Context, userID int64) ([]*entities.PrioritizationChangeRequest, error) {
	// Verificar se é admin ou manager
	isAdminOrManager, err := uc.isAdminOrManager(ctx, userID)
	if err != nil || !isAdminOrManager {
		return nil, errors.New("apenas administradores e gerentes podem visualizar solicitações")
	}

	return uc.prioritizationRepo.ListPendingChangeRequests(ctx)
}

// Helper:  Construir priorização com iniciativas completas
func (uc *PrioritizationUseCaseImpl) buildPrioritizationWithInitiatives(ctx context.Context, p *entities.InitiativePrioritization) (*entities.PrioritizationWithInitiatives, error) {
	// Buscar iniciativas do setor na ordem da priorização
	var initiatives []*entities.InitiativeListResponse

	for _, initiativeID := range p.PriorityOrder {
		initiative, err := uc.initiativeRepo.GetByIDWithCancellation(ctx, initiativeID)
		if err != nil {
			continue // Skip se não encontrar
		}

		initiatives = append(initiatives, &entities.InitiativeListResponse{
			ID:          initiative.ID,
			Title:       initiative.Title,
			Description: initiative.Description,
			Status:      initiative.Status,
			Type:        initiative.Type,
			Priority:    initiative.Priority,
			Sector:      initiative.Sector,
			OwnerName:   initiative.OwnerName,
			Date:        formatDate(initiative.CreatedAt),
		})
	}

	return &entities.PrioritizationWithInitiatives{
		ID:              p.ID,
		SectorID:        p.SectorID,
		SectorName:      p.SectorName,
		Year:            p.Year,
		IsLocked:        p.IsLocked,
		Initiatives:     initiatives,
		CreatedByUserID: p.CreatedByUserID,
		CreatedByName:   p.CreatedByName,
		CreatedAt:       p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// Helper: Verificar se é admin ou manager
func (uc *PrioritizationUseCaseImpl) isAdminOrManager(ctx context.Context, userID int64) (bool, error) {
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
