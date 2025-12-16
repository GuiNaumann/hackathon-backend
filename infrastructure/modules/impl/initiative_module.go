package module_impl

import (
	"encoding/json"
	"hackathon-backend/domain/entities"
	"hackathon-backend/domain/usecases"
	contextutil "hackathon-backend/utils/context"
	"hackathon-backend/utils/http_error"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type InitiativeModule struct {
	initiativeUseCase        usecases.InitiativeUseCase
	initiativeHistoryUseCase usecases.InitiativeHistoryUseCase
}

func NewInitiativeModule(initiativeUseCase usecases.InitiativeUseCase, historyUseCase usecases.InitiativeHistoryUseCase) *InitiativeModule {
	return &InitiativeModule{
		initiativeUseCase:        initiativeUseCase,
		initiativeHistoryUseCase: historyUseCase,
	}
}

func (m *InitiativeModule) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/initiatives", m.CreateInitiative).Methods("POST")
	router.HandleFunc("/initiatives", m.ListInitiatives).Methods("GET")
	router.HandleFunc("/initiatives/{id}", m.GetInitiative).Methods("GET")
	router.HandleFunc("/initiatives/{id}", m.UpdateInitiative).Methods("PUT")
	router.HandleFunc("/initiatives/{id}", m.DeleteInitiative).Methods("DELETE")
	router.HandleFunc("/initiatives/{id}/status", m.ChangeStatus).Methods("PATCH")
	router.HandleFunc("/initiatives/{id}/history", m.GetHistory).Methods("GET") // NOVO
	router.HandleFunc("/my-initiatives", m.GetMyInitiatives).Methods("GET")
}

func (m *InitiativeModule) CreateInitiative(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	var req entities.CreateInitiativeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	initiative, err := m.initiativeUseCase.CreateInitiative(r.Context(), &req, user.ID)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Iniciativa criada com sucesso",
		"data":    initiative,
	})
}

func (m *InitiativeModule) ListInitiatives(w http.ResponseWriter, r *http.Request) {
	// Parse query params para filtros
	filter := &entities.InitiativeFilter{
		Search:   r.URL.Query().Get("search"),
		Status:   r.URL.Query().Get("status"),
		Type:     r.URL.Query().Get("type"),
		Sector:   r.URL.Query().Get("sector"),
		Priority: r.URL.Query().Get("priority"),
	}

	initiatives, err := m.initiativeUseCase.ListInitiatives(r.Context(), filter)
	if err != nil {
		http_error.InternalServerError(w, "Erro ao listar iniciativas")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    initiatives,
		"count":   len(initiatives),
	})
}

func (m *InitiativeModule) GetInitiative(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	initiative, err := m.initiativeUseCase.GetInitiativeByID(r.Context(), id)
	if err != nil {
		http_error.NotFound(w, "Iniciativa não encontrada")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    initiative,
	})
}

func (m *InitiativeModule) UpdateInitiative(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	var req entities.UpdateInitiativeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	initiative, err := m.initiativeUseCase.UpdateInitiative(r.Context(), id, &req, user.ID)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Iniciativa atualizada com sucesso",
		"data":    initiative,
	})
}

func (m *InitiativeModule) DeleteInitiative(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	if err := m.initiativeUseCase.DeleteInitiative(r.Context(), id, user.ID); err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Iniciativa deletada com sucesso",
	})
}

func (m *InitiativeModule) ChangeStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	var req entities.ChangeInitiativeStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	if err := m.initiativeUseCase.ChangeStatus(r.Context(), id, &req, user.ID); err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Status alterado com sucesso",
	})
}

func (m *InitiativeModule) GetMyInitiatives(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	initiatives, err := m.initiativeUseCase.GetMyInitiatives(r.Context(), user.ID)
	if err != nil {
		http_error.InternalServerError(w, "Erro ao buscar suas iniciativas")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    initiatives,
		"count":   len(initiatives),
	})
}

// Adicionar o handler:
func (m *InitiativeModule) GetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	history, err := m.initiativeHistoryUseCase.GetHistory(r.Context(), id)
	if err != nil {
		http_error.InternalServerError(w, "Erro ao buscar histórico")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    history,
		"count":   len(history),
	})
}
