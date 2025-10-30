package pedidofornecedor

import (
	"encoding/json"
	"lapasta/internal/models"
	"net/http"
	"strconv"
)

type PedidoFornecedorHandler interface {
	CriarPedido(w http.ResponseWriter, r *http.Request)
	ListarPedidosPorFornecedor(w http.ResponseWriter, r *http.Request)
}

type pedidoFornecedorHandler struct {
	service PedidoFornecedorService
}

func NovoPedidoFornecedorHandler(service PedidoFornecedorService) PedidoFornecedorHandler {
	return &pedidoFornecedorHandler{
		service: service,
	}
}

func (h *pedidoFornecedorHandler) CriarPedido(w http.ResponseWriter, r *http.Request) {
	var pedido models.PedidoFornecedor
	response := models.ResponseDefaultModel{IsSuccess: true}

	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao decodificar o pedido"
		w.WriteHeader(http.StatusBadRequest)
	} else if err := h.service.CriarPedidoFornecedor(&pedido); err != nil {
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Data = pedido
		w.WriteHeader(http.StatusCreated)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *pedidoFornecedorHandler) ListarPedidosPorFornecedor(w http.ResponseWriter, r *http.Request) {
	valor := r.URL.Query().Get("fornecedorId")
	response := models.ResponseDefaultModel{IsSuccess: true}

	if valor == "" {
		response.IsSuccess = false
		response.ErrorMessage = "ID do fornecedor é obrigatório"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		id, err := strconv.Atoi(valor)
		if err != nil {
			response.IsSuccess = false
			response.ErrorMessage = "ID inválido"
			w.WriteHeader(http.StatusBadRequest)
		} else {
			pedidos, err := h.service.ListarPedidosPorFornecedor(id)
			if err != nil {
				response.IsSuccess = false
				response.ErrorMessage = "Erro ao listar pedidos"
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				response.Data = pedidos
				w.WriteHeader(http.StatusOK)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}
