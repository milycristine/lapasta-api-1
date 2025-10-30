package valetransporte

import (
	"encoding/json"
	"lapasta/internal/models"
	"log"
	"net/http"
	"strconv"
)

type ValeHandler struct {
	service ValeService
}

func NovoValeHandler(service ValeService) *ValeHandler {
	return &ValeHandler{
		service: service,
	}
}

func (h *ValeHandler) ListarVales(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	vales, err := h.service.ListarVales(pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      vales,
	}

	if err != nil {
		log.Printf("Erro ao listar vales: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar vales"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *ValeHandler) ListarValesDaSemana(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	semana := r.URL.Query().Get("semana")
	if semana == "" {
		http.Error(w, "Parâmetro 'semana' é obrigatório", http.StatusBadRequest)
		return
	}

	semanaInt, err := strconv.Atoi(semana)
	if err != nil || semanaInt <= 0 || semanaInt > 53 {
		http.Error(w, "Semana inválida", http.StatusBadRequest)
		return
	}

	mes := r.URL.Query().Get("mes")
	mesInt := 0
	if mes != "" {
		mesInt, err = strconv.Atoi(mes)
		if err != nil || mesInt <= 0 || mesInt > 12 {
			http.Error(w, "Mês inválido", http.StatusBadRequest)
			return
		}
	}

	vales, err := h.service.ListarValesDaSemana(semanaInt, pageInt, mesInt)

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      vales,
	}

	if err != nil {
		log.Printf("Erro ao listar vales da semana: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar vales da semana"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *ValeHandler) AtualizarVale(w http.ResponseWriter, r *http.Request) {
	var vale models.ValeTransportePagamento
	response := models.ResponseDefaultModel{
		IsSuccess: true,
	}

	if err := json.NewDecoder(r.Body).Decode(&vale); err != nil {
		log.Printf("Erro ao decodificar o vale: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Formato de entrada inválido"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if err := h.service.AtualizarStatusVale(vale.Id, vale.StatusIdVale); err != nil {
			log.Printf("Erro ao atualizar o vale: %v", err)
			response.IsSuccess = false
			response.Error = err
			response.ErrorMessage = "Erro ao atualizar o vale"
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(response)
}
