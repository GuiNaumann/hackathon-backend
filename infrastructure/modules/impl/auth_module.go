package module_impl

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"hackathon-backend/domain/entities"
	"hackathon-backend/domain/usecases"
	"hackathon-backend/settings_loader"
	contextutil "hackathon-backend/utils/context"
	"hackathon-backend/utils/http_error"
	"net/http"
	"time"
)

type AuthModule struct {
	authUseCase usecases.AuthUseCase
	settings    *settings_loader.SettingsLoader
}

func NewAuthModule(authUseCase usecases.AuthUseCase, settings *settings_loader.SettingsLoader) *AuthModule {
	return &AuthModule{
		authUseCase: authUseCase,
		settings:    settings,
	}
}

func (m *AuthModule) RegisterPublicRoutes(router *mux.Router) {
	router.HandleFunc("/login", m.Login).Methods("POST")
	router.HandleFunc("/logout", m.Logout).Methods("POST")
}

func (m *AuthModule) RegisterPrivateRoutes(router *mux.Router) {
	router.HandleFunc("/me", m.GetCurrentUser).Methods("GET")
}

func (m *AuthModule) Login(w http.ResponseWriter, r *http.Request) {
	var req entities.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	// Validar campos obrigatórios
	if req.Email == "" || req.Password == "" {
		http_error.BadRequest(w, "Email e senha são obrigatórios")
		return
	}

	// Executar login
	user, _, err := m.authUseCase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http_error.Unauthorized(w, err.Error())
		return
	}

	// Criar cookie seguro
	sc := securecookie.New([]byte(m.settings.Security.CookieEncryptionKey), nil)
	encoded, err := sc.Encode("auth_token", user.ID)
	if err != nil {
		http_error.InternalServerError(w, "Erro ao criar sessão")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    encoded,
		Path:     "/",
		Domain:   m.settings.Security.CookieDomain,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   m.settings.Security.CookieSecure,
		HttpOnly: m.settings.Security.CookieHTTPOnly,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

func (m *AuthModule) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Domain:   m.settings.Security.CookieDomain,
		Expires:  time.Now().Add(-1 * time.Hour),
		MaxAge:   -1,
		Secure:   m.settings.Security.CookieSecure,
		HttpOnly: m.settings.Security.CookieHTTPOnly,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Logout realizado com sucesso",
	})
}

func (m *AuthModule) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := contextutil.GetUserFromContext(r.Context())
	if !ok {
		http_error.Unauthorized(w, "Usuário não autenticado")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}
