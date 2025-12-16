package module_impl

import (
	"encoding/json"
	"fmt"
	"hackathon-backend/domain/entities"
	"hackathon-backend/domain/usecases"
	contextutil "hackathon-backend/utils/context"
	"hackathon-backend/utils/http_error"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type PrioritizationModule struct {
	prioritizationUseCase usecases.PrioritizationUseCase
}

func NewPrioritizationModule(prioritizationUseCase usecases.PrioritizationUseCase) *PrioritizationModule {
	return &PrioritizationModule{
		prioritizationUseCase: prioritizationUseCase,
	}
}

func (m *PrioritizationModule) RegisterRoutes(router *mux.Router) {
	// Rotas para usuários (seu setor)
	router.HandleFunc("/prioritization", m.GetMyPrioritization).Methods("GET")
	router.HandleFunc("/prioritization", m.SavePrioritization).Methods("POST")
	router.HandleFunc("/prioritization/request-change", m.RequestChange).Methods("POST")

	// Rotas para admin/manager (todos os setores)
	router.HandleFunc("/prioritization/all", m.GetAllSectorsPrioritization).Methods("GET")
	router.HandleFunc("/prioritization/change-requests", m.ListPendingChangeRequests).Methods("GET")
	router.HandleFunc("/prioritization/change-requests/{id}/review", m.ReviewChangeRequest).Methods("POST")
}

// GetMyPrioritization busca a priorização do setor do usuário
func (m *PrioritizationModule) GetMyPrioritization(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	// Pegar ano da query (default:  ano atual)
	yearStr := r.URL.Query().Get("year")
	year := time.Now().Year()
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	prioritization, err := m.prioritizationUseCase.GetPrioritization(r.Context(), year, user.ID)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    prioritization,
	})
}

// SavePrioritization salva a priorização do setor do usuário
func (m *PrioritizationModule) SavePrioritization(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	var req entities.SavePrioritizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	prioritization, err := m.prioritizationUseCase.SavePrioritization(r.Context(), &req, user.ID)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Priorização salva com sucesso",
		"data":    prioritization,
	})
}

// RequestChange solicita mudança na priorização (usuário normal)
func (m *PrioritizationModule) RequestChange(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	// Pegar ano da query (default: ano atual)
	yearStr := r.URL.Query().Get("year")
	year := time.Now().Year()
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	var req entities.RequestPrioritizationChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	changeRequest, err := m.prioritizationUseCase.RequestPrioritizationChange(r.Context(), &req, user.ID, year)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Solicitação de mudança criada com sucesso",
		"data":    changeRequest,
	})
}

// GetAllSectorsPrioritization busca priorização de todos os setores (Admin/Manager)
func (m *PrioritizationModule) GetAllSectorsPrioritization(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	// Pegar ano da query (default: ano atual)
	yearStr := r.URL.Query().Get("year")
	year := time.Now().Year()
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	allSectors, err := m.prioritizationUseCase.GetAllSectorsPrioritization(r.Context(), year, user.ID)
	if err != nil {
		http_error.Forbidden(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    allSectors,
	})
}

// ListPendingChangeRequests lista solicitações pendentes (Admin/Manager)
func (m *PrioritizationModule) ListPendingChangeRequests(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	requests, err := m.prioritizationUseCase.ListPendingChangeRequests(r.Context(), user.ID)
	if err != nil {
		http_error.Forbidden(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    requests,
		"count":   len(requests),
	})
}

// ReviewChangeRequest aprova ou recusa solicitação de mudança (Admin/Manager)
func (m *PrioritizationModule) ReviewChangeRequest(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	vars := mux.Vars(r)
	requestID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	var req entities.ReviewPrioritizationChangeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	if err := m.prioritizationUseCase.ReviewPrioritizationChange(r.Context(), requestID, &req, user.ID); err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	action := "reprovada"
	if req.Approved {
		action = "aprovada e a priorização foi desbloqueada para alteração"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Solicitação de mudança %s", action),
	})
}
