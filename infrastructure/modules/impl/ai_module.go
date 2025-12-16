package module_impl

import (
	"encoding/json"
	"hackathon-backend/domain/entities"
	"hackathon-backend/domain/usecases"
	contextutil "hackathon-backend/utils/context"
	"hackathon-backend/utils/http_error"
	"net/http"

	"github.com/gorilla/mux"
)

type AIModule struct {
	aiUseCase usecases.AIUseCase
}

func NewAIModule(aiUseCase usecases.AIUseCase) *AIModule {
	return &AIModule{
		aiUseCase: aiUseCase,
	}
}

func (m *AIModule) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/ai/refine-text", m.RefineText).Methods("POST")
}

func (m *AIModule) RefineText(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	var req entities.RefineTextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	result, err := m.aiUseCase.RefineText(r.Context(), &req, user.ID)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    result,
	})
}
