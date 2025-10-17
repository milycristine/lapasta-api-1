package gastos

import (
	"encoding/json"
	"lapasta/internal/models"
	"log"
	"net/http"
)

type GastosHandler struct {
	service GastosService
}

func NovoGastosHandler(service GastosService) *GastosHandler {
	return &GastosHandler{service: service}
}

func (h *GastosHandler) TotalGastos(w http.ResponseWriter, r *http.Request) {
	tipo := r.URL.Query().Get("tipo")
	total, err := h.service.TotalGastos(tipo)

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      total,
	}

	if err != nil {
		log.Printf("erro ao listar gastos: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar gastos"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *GastosHandler) ListarGastos(w http.ResponseWriter, r *http.Request) {
	tipo := r.URL.Query().Get("tipo")
	inicioData := r.URL.Query().Get("inicioData")
	fimData := r.URL.Query().Get("fimData")

	gastos, err := h.service.ListarGastos(tipo, inicioData, fimData)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      gastos,
	}

	if err != nil {
		log.Printf("erro ao listar gastos: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar gastos"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}
