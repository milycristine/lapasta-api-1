package funcionario

import (
	"encoding/json"
	"lapasta/internal/models"
	"log"
	"net/http"
	"strconv"
)

type FuncionarioHandler interface {
	CriarFuncionario(w http.ResponseWriter, r *http.Request)
	EditarFuncionario(w http.ResponseWriter, r *http.Request)
	ListarFuncionarios(w http.ResponseWriter, r *http.Request)
	BuscarFuncionarioPorCPF(w http.ResponseWriter, r *http.Request)
	BuscarFuncionarioPorID(w http.ResponseWriter, r *http.Request)
	AtualizarStatusFuncionario(w http.ResponseWriter, r *http.Request)
}

type funcionarioHandler struct {
	service FuncionarioService
}

func NovoFuncionarioHandler(service FuncionarioService) FuncionarioHandler {
	return &funcionarioHandler{
		service: service,
	}
}

func (h *funcionarioHandler) CriarFuncionario(w http.ResponseWriter, r *http.Request) {
	var funcionario models.Funcionario
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      funcionario,
	}

	if err := json.NewDecoder(r.Body).Decode(&funcionario); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao decodificar o funcionário"
		w.WriteHeader(http.StatusBadRequest)
	} else if err := h.service.CriarFuncionario(&funcionario); err != nil {
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

func (h *funcionarioHandler) EditarFuncionario(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var funcionario models.Funcionario
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      funcionario,
	}

	if err := json.NewDecoder(r.Body).Decode(&funcionario); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao decodificar os dados do funcionário"
		w.WriteHeader(http.StatusBadRequest)
	} else if err := h.service.EditarFuncionario(&funcionario); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *funcionarioHandler) ListarFuncionarios(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page := 1

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	funcionarios, err := h.service.ListarFuncionarios(page)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      funcionarios,
	}

	if err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar funcionários"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *funcionarioHandler) BuscarFuncionarioPorCPF(w http.ResponseWriter, r *http.Request) {
	cpf := r.URL.Query().Get("cpf")
	if cpf == "" {
		http.Error(w, "CPF é obrigatório", http.StatusBadRequest)
		return
	}

	funcionarioComPontos, err := h.service.BuscarFuncionarioPorCPF(cpf)
	if err != nil {
		http.Error(w, "Erro ao buscar funcionário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if funcionarioComPontos == nil {
		http.Error(w, "Funcionário não encontrado", http.StatusNotFound)
		return
	}

	response := struct {
		Id          int            `json:"id"`
		Nome        string         `json:"nome"`
		Pontos      []models.Ponto `json:"pontos"`
		BotaoStatus string         `json:"botaoStatus"`
	}{
		Id:          funcionarioComPontos.Funcionario.Id,
		Nome:        funcionarioComPontos.Funcionario.Nome,
		Pontos:      funcionarioComPontos.Pontos,
		BotaoStatus: funcionarioComPontos.BotaoStatus,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *funcionarioHandler) BuscarFuncionarioPorID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID é obrigatório", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	funcionario, err := h.service.BuscarFuncionarioPorID(id)
	if err != nil {
		http.Error(w, "Erro ao buscar funcionário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if funcionario == nil {
		http.Error(w, "Funcionário não encontrado", http.StatusNotFound)
		return
	}

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      funcionario,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *funcionarioHandler) AtualizarStatusFuncionario(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	statusStr := r.URL.Query().Get("status")

	var response models.ResponseDefaultModel

	id, err := strconv.Atoi(idStr)
	if err != nil || (statusStr != "true" && statusStr != "false") {
		response.IsSuccess = false
		response.ErrorMessage = "Parâmetros inválidos (use ?id=1&status=true/false)"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	status := statusStr == "true"

	err = h.service.AtualizarStatusFuncionario(id, status)
	if err != nil {
		log.Println("Erro ao atualizar status do funcionario:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao atualizar status do funcionario"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.IsSuccess = true
		response.Data = map[string]interface{}{
			"id":     id,
			"status": status,
		}
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}
