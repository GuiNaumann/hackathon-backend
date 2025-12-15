package infrastructure

import (
	_ "hackathon-backend/domain/entities"
	"hackathon-backend/domain/usecases"
	"hackathon-backend/utils/http_error"
	"net/http"
)

func NewPermissionMiddleware(permUseCase usecases.PermissionUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Pegar usuário do contexto (já injetado pelo AuthMiddleware)
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http_error.Unauthorized(w, "Usuário não autenticado")
				return
			}

			// 2. Verificar se o usuário tem permissão para este endpoint
			endpoint := r.URL.Path
			method := r.Method

			hasPermission, err := permUseCase.HasPermission(r.Context(), user.ID, endpoint, method)
			if err != nil {
				http_error.InternalServerError(w, "Erro ao verificar permissões")
				return
			}

			if !hasPermission {
				http_error.Forbidden(w, "Você não tem permissão para acessar este recurso")
				return
			}

			// 3. Usuário tem permissão, continuar
			next.ServeHTTP(w, r)
		})
	}
}

// Middleware específico para rotas que exigem um tipo específico
func RequireUserType(permUseCase usecases.PermissionUseCase, requiredTypes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http_error.Unauthorized(w, "Usuário não autenticado")
				return
			}

			// Buscar tipos do usuário
			userTypes, err := permUseCase.GetUserTypes(r.Context(), user.ID)
			if err != nil {
				http_error.InternalServerError(w, "Erro ao verificar tipo de usuário")
				return
			}

			// Verificar se o usuário tem algum dos tipos requeridos
			hasRequiredType := false
			for _, userType := range userTypes {
				for _, requiredType := range requiredTypes {
					if userType.Name == requiredType {
						hasRequiredType = true
						break
					}
				}
				if hasRequiredType {
					break
				}
			}

			if !hasRequiredType {
				http_error.Forbidden(w, "Você não tem o tipo de usuário necessário para acessar este recurso")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
