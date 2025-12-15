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

type CommentModule struct {
	commentUseCase usecases.CommentUseCase
}

func NewCommentModule(commentUseCase usecases.CommentUseCase) *CommentModule {
	return &CommentModule{
		commentUseCase: commentUseCase,
	}
}

func (m *CommentModule) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/initiatives/{initiativeId}/comments", m.CreateComment).Methods("POST")
	router.HandleFunc("/initiatives/{initiativeId}/comments", m.ListComments).Methods("GET")
	router.HandleFunc("/comments/{id}", m.UpdateComment).Methods("PUT")
	router.HandleFunc("/comments/{id}", m.DeleteComment).Methods("DELETE")
}

func (m *CommentModule) CreateComment(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	vars := mux.Vars(r)
	initiativeID, err := strconv.ParseInt(vars["initiativeId"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID de iniciativa inválido")
		return
	}

	var req entities.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	comment, err := m.commentUseCase.CreateComment(r.Context(), initiativeID, &req, user.ID)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Comentário criado com sucesso",
		"data":    comment,
	})
}

func (m *CommentModule) ListComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	initiativeID, err := strconv.ParseInt(vars["initiativeId"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID de iniciativa inválido")
		return
	}

	comments, err := m.commentUseCase.ListComments(r.Context(), initiativeID)
	if err != nil {
		http_error.InternalServerError(w, "Erro ao listar comentários")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    comments,
		"count":   len(comments),
	})
}

func (m *CommentModule) UpdateComment(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	vars := mux.Vars(r)
	commentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID de comentário inválido")
		return
	}

	var req entities.UpdateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	comment, err := m.commentUseCase.UpdateComment(r.Context(), commentID, &req, user.ID)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Comentário atualizado com sucesso",
		"data":    comment,
	})
}

func (m *CommentModule) DeleteComment(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	vars := mux.Vars(r)
	commentID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID de comentário inválido")
		return
	}

	if err := m.commentUseCase.DeleteComment(r.Context(), commentID, user.ID); err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Comentário deletado com sucesso",
	})
}
