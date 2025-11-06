package ponto

import (
	"encoding/json"
	"fmt"
	"lapasta/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PontoHandler struct {
	service PontoService
}

func NovoPontoHandler(service PontoService) *PontoHandler {
	return &PontoHandler{
		service: service,
	}
}

func (h *PontoHandler) ListarPontos(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	log.Printf("Listando pontos na página: %d", pageInt)
	documentos, err := h.service.ListarPontos(pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      documentos,
	}

	if err != nil {
		log.Printf("Erro ao listar pontos: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar pontos"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}

func (h *PontoHandler) ListarPontosPorId(w http.ResponseWriter, r *http.Request) {
	idFuncionario := r.URL.Query().Get("idFuncionario")
	if idFuncionario == "" {
		http.Error(w, "idFuncionario é obrigatório", http.StatusBadRequest)
		return
	}

	idFuncionarioInt, err := strconv.Atoi(idFuncionario)
	if err != nil {
		http.Error(w, "idFuncionario inválido", http.StatusBadRequest)
		return
	}
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "page inválido", http.StatusBadRequest)
			return
		}
	}
	pontos, err := h.service.ListarPontosPorId(idFuncionarioInt, page)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pontos,
	}

	if err != nil {
		log.Printf("Erro ao listar pontos: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar pontos"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *PontoHandler) ListarPontosPorIdEDia(w http.ResponseWriter, r *http.Request) {
	idFuncionario := r.URL.Query().Get("idFuncionario")
	dia := r.URL.Query().Get("dia")

	if idFuncionario == "" {
		http.Error(w, "idFuncionario é obrigatório", http.StatusBadRequest)
		return
	}

	if dia == "" {
		http.Error(w, "Dia é obrigatório", http.StatusBadRequest)
		return
	}
	idFuncionarioInt, err := strconv.Atoi(idFuncionario)
	if err != nil {
		http.Error(w, "idFuncionario inválido", http.StatusBadRequest)
		return
	}

	pontos, err := h.service.ListarPontosPorIdEDia(idFuncionarioInt, dia)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pontos,
	}

	if err != nil {
		log.Printf("Erro ao listar por id e dia: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar por id e dia"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *PontoHandler) RegistrarEntrada(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IdFuncionario int `json:"id_funcionario"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.service.RegistrarEntrada(body.IdFuncionario)
	response := models.ResponseDefaultModel{
		IsSuccess: err == nil,
	}

	if err != nil {
		log.Printf("Erro ao registrar entrada: %v", err)
		response.Error = err
		response.ErrorMessage = "Erro ao registrar entrada"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *PontoHandler) RegistrarSaidaAlmoco(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IdFuncionario int `json:"id_funcionario"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.service.RegistrarSaidaAlmoco(body.IdFuncionario)
	response := models.ResponseDefaultModel{
		IsSuccess: err == nil,
	}

	if err != nil {
		log.Printf("Erro ao registrar saída para almoço: %v", err)
		response.Error = err
		response.ErrorMessage = "Erro ao registrar saída para almoço"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *PontoHandler) RegistrarRetornoAlmoco(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IdFuncionario int `json:"id_funcionario"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.service.RegistrarRetornoAlmoco(body.IdFuncionario)
	response := models.ResponseDefaultModel{
		IsSuccess: err == nil,
	}

	if err != nil {
		log.Printf("Erro ao registrar retorno do almoço: %v", err)
		response.Error = err
		response.ErrorMessage = "Erro ao registrar retorno do almoço"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *PontoHandler) RegistrarSaida(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IdFuncionario int `json:"id_funcionario"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.service.RegistrarSaida(body.IdFuncionario)
	response := models.ResponseDefaultModel{
		IsSuccess: err == nil,
	}

	if err != nil {
		log.Printf("Erro ao registrar saída final: %v", err)
		response.Error = err
		response.ErrorMessage = "Erro ao registrar saída final"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}
func (h *PontoHandler) ListarPontosPorData(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Os parâmetros 'startDate' e 'endDate' são obrigatórios", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Formato de 'startDate' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Formato de 'endDate' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	log.Printf("Listando pontos por data na página: %d", pageInt)
	pontos, err := h.service.ListarPontosPorData(startDate, endDate, pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pontos,
	}

	if err != nil {
		log.Printf("Erro ao listar pontos por data: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar pontos por data"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}

func (h *PontoHandler) ListarPontosPorDataId(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	idFuncionarioStr := r.URL.Query().Get("idFuncionario")
	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	if idFuncionarioStr == "" || startDateStr == "" || endDateStr == "" {
		http.Error(w, "Os parâmetros 'idFuncionario', 'startDate' e 'endDate' são obrigatórios", http.StatusBadRequest)
		return
	}

	idFuncionario, err := strconv.Atoi(idFuncionarioStr)
	if err != nil {
		http.Error(w, "O parâmetro 'idFuncionario' deve ser um número inteiro", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Formato de 'startDate' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Formato de 'endDate' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	log.Printf("Listando pontos por data e ID de funcionário na página: %d", pageInt)
	pontos, err := h.service.ListarPontosPorDataId(idFuncionario, startDate, endDate, pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pontos,
	}

	if err != nil {
		log.Printf("Erro ao listar pontos por data e ID de funcionário: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar pontos por data e ID de funcionário"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}
func (h *PontoHandler) GerarRelatorioMensal(w http.ResponseWriter, r *http.Request) {
	mesParam := r.URL.Query().Get("mes")
	anoParam := r.URL.Query().Get("ano")

	var mes, ano int
	var err error

	if mesParam != "" {
		mes, err = strconv.Atoi(mesParam)
		if err != nil || mes < 1 || mes > 12 {
			http.Error(w, "Mês inválido", http.StatusBadRequest)
			return
		}
	} else {
		mes = int(time.Now().Month())
	}

	if anoParam != "" {
		ano, err = strconv.Atoi(anoParam)
		if err != nil || ano < 2000 || ano > time.Now().Year() {
			http.Error(w, "Ano inválido", http.StatusBadRequest)
			return
		}
	} else {
		ano = time.Now().Year()
	}

	filePath, err := h.service.GerarRelatorioMensal(mes, ano, "email@gmail.com")
	if err != nil {
		log.Printf("Erro ao gerar relatório mensal: %v", err)
		http.Error(w, "Erro ao gerar relatório mensal", http.StatusInternalServerError)
		return
	}

	errorMessage := fmt.Sprintf("Relatório mensal gerado com sucesso para %02d/%d. O arquivo está disponível em: %s", mes, ano, filePath)

	response := models.ResponseDefaultModel{
		IsSuccess:    true,
		Data:         errorMessage,
		Error:        nil,
		ErrorMessage: "",
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}
