package usecase_impl

import (
	"context"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories"
	"time"
)

type InitiativeHistoryUseCaseImpl struct {
	historyRepo repositories.InitiativeHistoryRepository
}

func NewInitiativeHistoryUseCaseImpl(historyRepo repositories.InitiativeHistoryRepository) *InitiativeHistoryUseCaseImpl {
	return &InitiativeHistoryUseCaseImpl{
		historyRepo: historyRepo,
	}
}

func (uc *InitiativeHistoryUseCaseImpl) GetHistory(ctx context.Context, initiativeID int64) ([]*entities.InitiativeHistoryResponse, error) {
	histories, err := uc.historyRepo.ListByInitiative(ctx, initiativeID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar histórico: %w", err)
	}

	var response []*entities.InitiativeHistoryResponse
	for _, history := range histories {
		response = append(response, &entities.InitiativeHistoryResponse{
			ID:           history.ID,
			InitiativeID: history.InitiativeID,
			UserID:       history.UserID,
			UserName:     history.UserName,
			OldStatus:    history.OldStatus,
			NewStatus:    history.NewStatus,
			Reason:       history.Reason,
			CreatedAt:    history.CreatedAt.Format("2006-01-02 15:04:05"),
			TimeAgo:      timeAgo(history.CreatedAt),
		})
	}

	return response, nil
}

func timeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration.Hours() < 1 {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "há 1 minuto"
		}
		return fmt.Sprintf("há %d minutos", minutes)
	}

	if duration.Hours() < 24 {
		hours := int(duration.Hours())
		if hours == 1 {
			return "há 1 hora"
		}
		return fmt.Sprintf("há %d horas", hours)
	}

	days := int(duration.Hours() / 24)
	if days == 1 {
		return "há 1 dia"
	}
	if days < 30 {
		return fmt.Sprintf("há %d dias", days)
	}

	months := days / 30
	if months == 1 {
		return "há 1 mês"
	}
	if months < 12 {
		return fmt.Sprintf("há %d meses", months)
	}

	years := months / 12
	if years == 1 {
		return "há 1 ano"
	}
	return fmt.Sprintf("há %d anos", years)
}
