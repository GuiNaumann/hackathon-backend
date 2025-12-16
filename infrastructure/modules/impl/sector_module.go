package module_impl

import (
	"encoding/json"
	"hackathon-backend/domain/entities"
	"hackathon-backend/domain/usecases"
	"hackathon-backend/utils/http_error"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SectorModule struct {
	sectorUseCase usecases.SectorUseCase
}

func NewSectorModule(sectorUseCase usecases.SectorUseCase) *SectorModule {
	return &SectorModule{
		sectorUseCase: sectorUseCase,
	}
}

func (m *SectorModule) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/sectors", m.CreateSector).Methods("POST")
	router.HandleFunc("/sectors", m.ListSectors).Methods("GET")
	router.HandleFunc("/sectors/{id}", m.GetSector).Methods("GET")
	router.HandleFunc("/sectors/{id}", m.UpdateSector).Methods("PUT")
	router.HandleFunc("/sectors/{id}", m.DeleteSector).Methods("DELETE")
}

func (m *SectorModule) CreateSector(w http.ResponseWriter, r *http.Request) {
	var req entities.CreateSectorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	sector, err := m.sectorUseCase.CreateSector(r.Context(), &req)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Setor criado com sucesso",
		"data":    sector,
	})
}

func (m *SectorModule) ListSectors(w http.ResponseWriter, r *http.Request) {
	// Query param para filtrar apenas ativos
	activeOnly := r.URL.Query().Get("active_only") == "true"

	sectors, err := m.sectorUseCase.ListSectors(r.Context(), activeOnly)
	if err != nil {
		http_error.InternalServerError(w, "Erro ao listar setores")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    sectors,
		"count":   len(sectors),
	})
}

func (m *SectorModule) GetSector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sectorID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	sector, err := m.sectorUseCase.GetSectorByID(r.Context(), sectorID)
	if err != nil {
		http_error.NotFound(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    sector,
	})
}

func (m *SectorModule) UpdateSector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sectorID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	var req entities.UpdateSectorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http_error.BadRequest(w, "Payload inválido")
		return
	}

	sector, err := m.sectorUseCase.UpdateSector(r.Context(), sectorID, &req)
	if err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Setor atualizado com sucesso",
		"data":    sector,
	})
}

func (m *SectorModule) DeleteSector(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sectorID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http_error.BadRequest(w, "ID inválido")
		return
	}

	if err := m.sectorUseCase.DeleteSector(r.Context(), sectorID); err != nil {
		http_error.BadRequest(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Setor deletado com sucesso",
	})
}
