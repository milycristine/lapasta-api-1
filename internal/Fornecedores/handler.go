package fornecedor

import (
	"encoding/json"
	"lapasta/internal/models"
	"log"
	"net/http"
	"strconv"
)

type FornecedorHandler interface {
	CriarFornecedor(w http.ResponseWriter, r *http.Request)
	EditarFornecedor(w http.ResponseWriter, r *http.Request)
	ListarFornecedores(w http.ResponseWriter, r *http.Request)
	BuscarFornecedorPorCNPJouNome(w http.ResponseWriter, r *http.Request)
	CriarPedido(w http.ResponseWriter, r *http.Request)
	ListarPedidosPorFornecedor(w http.ResponseWriter, r *http.Request)
	BuscarPedidosFornecedorPorDescricaoOuId(w http.ResponseWriter, r *http.Request)
}

type fornecedorHandler struct {
	service FornecedorService
}

func NovoFornecedorHandler(service FornecedorService) FornecedorHandler {
	return &fornecedorHandler{
		service: service,
	}
}

func (h *fornecedorHandler) CriarFornecedor(w http.ResponseWriter, r *http.Request) {
	var fornecedor models.Fornecedor
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      fornecedor,
	}

	if err := json.NewDecoder(r.Body).Decode(&fornecedor); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao decodificar o fornecedor"
		w.WriteHeader(http.StatusBadRequest)
	} else if err := h.service.CriarFornecedor(&fornecedor); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *fornecedorHandler) EditarFornecedor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var fornecedor models.Fornecedor
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      fornecedor,
	}

	if err := json.NewDecoder(r.Body).Decode(&fornecedor); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao decodificar os dados do funcionário"
		w.WriteHeader(http.StatusBadRequest)
	} else if err := h.service.EditarFornecedor(&fornecedor); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *fornecedorHandler) ListarFornecedores(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page := 1

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	fornecedores, err := h.service.ListarFornecedores(page)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      fornecedores,
	}

	if err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar fornecedores"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *fornecedorHandler) BuscarFornecedorPorCNPJouNome(w http.ResponseWriter, r *http.Request) {
	valor := r.URL.Query().Get("valor")
	if valor == "" {
		http.Error(w, "CNPJ ou nome é obrigatório", http.StatusBadRequest)
		return
	}

	fornecedor, err := h.service.BuscarFornecedorPorCNPJouNome(valor)
	if err != nil {
		http.Error(w, "Erro ao buscar fornecedor: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      fornecedor,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *fornecedorHandler) CriarPedido(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var pedido models.PedidoFornecedor
	response := models.ResponseDefaultModel{IsSuccess: true}

	if err := json.NewDecoder(r.Body).Decode(&pedido); err != nil {
		log.Printf("Erro ao decodificar pedido: %v", err)
		http.Error(w, `{"isSuccess":false,"errorMessage":"Erro ao decodificar o pedido"}`, http.StatusBadRequest)
		return
	}

	log.Printf("Pedido recebido: %+v", pedido)

	if err := h.service.CriarPedidoFornecedor(&pedido); err != nil {
		log.Printf("Erro ao criar pedido: %v", err)
		http.Error(w, `{"isSuccess":false,"errorMessage":"Erro ao criar pedido"}`, http.StatusInternalServerError)
		return
	}

	response.Data = pedido
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *fornecedorHandler) ListarPedidosPorFornecedor(w http.ResponseWriter, r *http.Request) {
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

func (h *fornecedorHandler) BuscarPedidosFornecedorPorDescricaoOuId(w http.ResponseWriter, r *http.Request) {
	fornecedorIdStr := r.URL.Query().Get("fornecedorId")
	valor := r.URL.Query().Get("valor")

	response := models.ResponseDefaultModel{IsSuccess: true}

	if fornecedorIdStr == "" || valor == "" {
		response.IsSuccess = false
		response.ErrorMessage = "fornecedorId e valor são obrigatórios"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		fornecedorId, err := strconv.Atoi(fornecedorIdStr)
		if err != nil {
			response.IsSuccess = false
			response.ErrorMessage = "fornecedorId inválido"
			w.WriteHeader(http.StatusBadRequest)
		} else {
			pedidos, err := h.service.BuscarPedidosFornecedorPorDescricaoOuId(fornecedorId, valor)
			if err != nil {
				response.IsSuccess = false
				response.ErrorMessage = "Erro ao buscar pedidos"
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
