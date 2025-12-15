package module_impl

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"hackathon-backend/domain/entities"
	"hackathon-backend/domain/usecases"
	contextutil "hackathon-backend/utils/context"
	"hackathon-backend/utils/http_error"
	"net/http"
	"strconv"
)

type UserCrudModule struct {
	userCrudUseCase usecases.UserCrudUseCase
}

func NewUserCrudModule(userCrudUseCase usecases.UserCrudUseCase) *UserCrudModule {
	return &UserCrudModule{
		userCrudUseCase: userCrudUseCase,
	}
}

func (m *UserCrudModule) RegisterRoutes(router *mux.Router) {
	// Rotas protegidas para admin
	router.HandleFunc("/users", m.CreateUser).Methods("POST")
	router.HandleFunc("/users", m.ListUsers).Methods("GET")
	router.HandleFunc("/users/{id}", m.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", m.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", m.DeleteUser).Methods("DELETE")

	// Rota para qualquer usuário autenticado mudar sua própria senha
	router.HandleFunc("/change-password", m.ChangePassword).Methods("POST")
}

func (m *UserCrudModule) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req entities.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	user, err := m.userCrudUseCase.CreateUser(r.Context(), &req)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Usuário criado com sucesso",
		"user":    user,
	})
}

func (m *UserCrudModule) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := m.userCrudUseCase.ListUsers(r.Context())
	if err != nil {
		http_error.InternalServerError(w, "Erro ao listar usuários")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    users,
		"count":   len(users),
	})
}

func (m *UserCrudModule) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	user, err := m.userCrudUseCase.GetUserByID(r.Context(), userID)
	if err != nil {
		http_error.NotFound(w, "Usuário não encontrado")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

func (m *UserCrudModule) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	var req entities.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	user, err := m.userCrudUseCase.UpdateUser(r.Context(), userID, &req)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Usuário atualizado com sucesso",
		"user":    user,
	})
}

func (m *UserCrudModule) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	// Impedir que o usuário delete a si mesmo
	currentUser, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	if currentUser.ID == userID {
		http_error.BadRequest(w, "Você não pode deletar sua própria conta")
		return
	}

	if err := m.userCrudUseCase.DeleteUser(r.Context(), userID); err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Usuário deletado com sucesso",
	})
}

func (m *UserCrudModule) ChangePassword(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	var req entities.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	if err := m.userCrudUseCase.ChangePassword(r.Context(), currentUser.ID, &req); err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Senha alterada com sucesso",
	})
}
