package module_impl

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"hackathon-backend/domain/entities"
	"hackathon-backend/domain/usecases"
	//"hackathon-backend/infrastructure"
	"hackathon-backend/utils/http_error"
	"net/http"
	"strconv"
)

type PermissionModule struct {
	permUseCase usecases.PermissionUseCase
}

func NewPermissionModule(permUseCase usecases.PermissionUseCase) *PermissionModule {
	return &PermissionModule{
		permUseCase: permUseCase,
	}
}

func (m *PermissionModule) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/personal-information", m.GetPersonalInformation).Methods("GET")
	router.HandleFunc("/admin/users/{userId}/types/{typeId}", m.AssignUserType).Methods("POST")
	router.HandleFunc("/admin/users/{userId}/types/{typeId}", m.RemoveUserType).Methods("DELETE")
}

func (m *PermissionModule) GetPersonalInformation(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("auth-ctx-user-data").(*entities.User)
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	personalInfo, err := m.permUseCase.GetPersonalInformation(r.Context(), user.ID)
	if err != nil {
		http_error.InternalServerError(w, "Erro ao buscar informações pessoais")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    personalInfo,
	})
}

func (m *PermissionModule) AssignUserType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID de usuário inválido")
		return
	}

	typeID, err := strconv.ParseInt(vars["typeId"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID de tipo inválido")
		return
	}

	if err := m.permUseCase.AssignUserType(r.Context(), userID, typeID); err != nil {
		http_error.InternalServerError(w, "Erro ao atribuir tipo ao usuário")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tipo atribuído com sucesso",
	})
}

func (m *PermissionModule) RemoveUserType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID de usuário inválido")
		return
	}

	typeID, err := strconv.ParseInt(vars["typeId"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID de tipo inválido")
		return
	}

	if err := m.permUseCase.RemoveUserType(r.Context(), userID, typeID); err != nil {
		http_error.InternalServerError(w, "Erro ao remover tipo do usuário")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tipo removido com sucesso",
	})
}
