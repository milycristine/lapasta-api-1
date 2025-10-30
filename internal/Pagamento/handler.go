package pagamento

import (
	"encoding/json"
	"lapasta/internal/models"
	"log"
	"net/http"
	"strconv"
)

type PagamentoHandler struct {
	service PagamentoService
}

func NovoPagamentoHandler(service PagamentoService) *PagamentoHandler {
	return &PagamentoHandler{
		service: service,
	}
}

func (h *PagamentoHandler) CriarPagamento(w http.ResponseWriter, r *http.Request) {
	var pagamento models.Pagamento
	response := models.ResponseDefaultModel{
		IsSuccess: true,
	}

	if err := json.NewDecoder(r.Body).Decode(&pagamento); err != nil {
		log.Printf("Erro ao decodificar o pagamento: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Formato de entrada inválido"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if err := h.service.CriarPagamento(&pagamento); err != nil {
			log.Printf("Erro ao criar o pagamento: %v", err)
			response.IsSuccess = false
			response.Error = err
			response.ErrorMessage = "Erro ao criar o pagamento"
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Data = pagamento
			w.WriteHeader(http.StatusCreated)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *PagamentoHandler) ListarPagamentos(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	log.Printf("Listando pagamentos na página: %d", pageInt)
	pagamentos, err := h.service.ListarPagamentos(pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pagamentos,
	}

	if err != nil {
		log.Printf("Erro ao listar pagamentos: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar pagamentos"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}
func (h *PagamentoHandler) ListarPagamentosPorDia(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	dateStr := r.URL.Query().Get("dia")
	if dateStr == "" {
		http.Error(w, "Parâmetro 'dia' é obrigatório", http.StatusBadRequest)
		return
	}

	log.Printf("Listando pagamentos no dia %s, página: %d", dateStr, pageInt)
	pagamentos, err := h.service.ListarPagamentosPorDia(dateStr, pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pagamentos,
	}

	if err != nil {
		log.Printf("Erro ao listar pagamentos: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar pagamentos"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}
func (h *PagamentoHandler) ListarPagamentosPorMes(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	mes := r.URL.Query().Get("mes")
	if mes == "" {
		http.Error(w, "Parâmetro 'mes' é obrigatório", http.StatusBadRequest)
		return
	}

	mesInt, err := strconv.Atoi(mes)
	if err != nil || mesInt < 1 || mesInt > 12 {
		http.Error(w, "Mês inválido", http.StatusBadRequest)
		return
	}

	log.Printf("Listando pagamentos no mês %d, página: %d", mesInt, pageInt)
	pagamentos, err := h.service.ListarPagamentosPorMes(mesInt, pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pagamentos,
	}

	if err != nil {
		log.Printf("Erro ao listar pagamentos por mês: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar pagamentos por mês"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}
func (h *PagamentoHandler) AtualizarPagamento(w http.ResponseWriter, r *http.Request) {
	var pagamento models.Pagamento
	response := models.ResponseDefaultModel{
		IsSuccess: true,
	}

	if err := json.NewDecoder(r.Body).Decode(&pagamento); err != nil {
		log.Printf("Erro ao decodificar o pagamento: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Formato de entrada inválido"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if err := h.service.AtualizarStatusPagamento(pagamento.Id, pagamento.StatusId); err != nil {
			log.Printf("Erro ao atualizar o pagamento: %v", err)
			response.IsSuccess = false
			response.Error = err
			response.ErrorMessage = "Erro ao atualizar o pagamento"
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}
