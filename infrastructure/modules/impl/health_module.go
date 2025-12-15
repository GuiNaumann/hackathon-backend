package module_impl

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HealthModule struct{}

func NewHealthModule() *HealthModule {
	return &HealthModule{}
}

func (m *HealthModule) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/health", m.HealthCheck).Methods("GET")
}

func (m *HealthModule) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
