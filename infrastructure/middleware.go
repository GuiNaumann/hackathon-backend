package infrastructure

import (
	"github.com/gorilla/securecookie"
	//"hackathon-backend/domain/entities"
	"hackathon-backend/infrastructure/repositories/impl"
	"hackathon-backend/settings_loader"
	contextutil "hackathon-backend/utils/context"
	"hackathon-backend/utils/http_error"
	"net/http"
)

func NewAuthMiddleware(authRepo *repository_impl.AuthRepositoryImpl, settings *settings_loader.SettingsLoader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Ler cookie
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				http_error.Unauthorized(w, "Token não encontrado")
				return
			}

			// 2. Decodificar token
			sc := securecookie.New([]byte(settings.Security.CookieEncryptionKey), nil)
			var userID int64
			if err := sc.Decode("auth_token", cookie.Value, &userID); err != nil {
				http_error.Unauthorized(w, "Token inválido")
				return
			}

			// 3. Buscar usuário no banco
			user, err := authRepo.GetUserByID(r.Context(), userID)
			if err != nil {
				http_error.Unauthorized(w, "Usuário não encontrado")
				return
			}

			// 4. Injetar usuário no contexto
			ctx := contextutil.SetUserInContext(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
