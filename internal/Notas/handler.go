package nota

import (
	//"encoding/base64"
	"encoding/json"

	//"lapasta/config"
	//Utils "lapasta/internal/Utils"

	"lapasta/internal/models"
	//"lapasta/internal/s3client"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type NotaHandler struct {
	service NotaService
}

func NovaNotaHandler(service NotaService) *NotaHandler {
	return &NotaHandler{
		service: service,
	}
}

func (h *NotaHandler) CriarNota(w http.ResponseWriter, r *http.Request) {
	var nota models.Nota
	response := models.ResponseDefaultModel{
		IsSuccess: true,
	}

	if err := json.NewDecoder(r.Body).Decode(&nota); err != nil {
		log.Printf("Erro ao decodificar a nota: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Formato de entrada inválido"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	switch strings.ToUpper(nota.Tipo) {
	case "FORNECEDOR":
		if nota.IdFornecedor == nil || *nota.IdFornecedor == 0 {
			response.IsSuccess = false
			response.ErrorMessage = "Fornecedor obrigatório para nota do tipo FORNECEDOR"
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
			return
		}

	case "AVULSO":
		nota.IdFornecedor = nil
		nota.IdPedidoFornecedor = nil

	default:
		response.IsSuccess = false
		response.ErrorMessage = "Tipo de nota inválido. Permitidos: FORNECEDOR ou AVULSA"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	//base64String := nota.ImagemBase64
	//if idx := strings.Index(base64String, ","); idx != -1 {
	//	base64String = base64String[idx+1:]
	//}

	//fileBytes, err := base64.StdEncoding.DecodeString(base64String)
	//if err != nil {
	//	log.Println("Erro ao decodificar base64:", err)
	//	response.IsSuccess = false
	//	response.ErrorMessage = "Erro ao decodificar a imagem base64"
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	json.NewEncoder(w).Encode(response)
	//	return
	//}

	//extensao := ".png"
	//nomeImagem := Utils.GerarStringAleatoria(12) + extensao
	//urlImagem, err := Utils.UploadImagemFirebase(fileBytes, nomeImagem, "notas")
	//if err != nil {
	//	log.Println("Erro ao salvar imagem:", err)
	//	response.IsSuccess = false
	//	response.ErrorMessage = "Erro ao salvar imagem"
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	json.NewEncoder(w).Encode(response)
	//	return
	//}
	//
	//nota.UrlImagem = urlImagem

	idNota, err := h.service.CriarNota(&nota)
	if err != nil {
		log.Printf("Erro ao criar a nota: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	nota.Id = idNota

	response.Data = nota
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *NotaHandler) ListarNotas(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	log.Printf("Listando notas na página: %d", pageInt)
	notas, err := h.service.ListarNotas(pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      notas,
	}

	if err != nil {
		log.Printf("Erro ao listar notas: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar notas"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}

func (h *NotaHandler) FiltrarDataNota(w http.ResponseWriter, r *http.Request) {
	inicioDataStr := r.URL.Query().Get("inicioData")
	fimDataStr := r.URL.Query().Get("fimData")

	if inicioDataStr == "" || fimDataStr == "" {
		http.Error(w, "Os parâmetros 'inicioData' e 'fimData' são obrigatórios", http.StatusBadRequest)
		return
	}

	inicioData, err := time.Parse("2006-01-02", inicioDataStr)
	if err != nil {
		http.Error(w, "Formato de 'inicioData' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	fimData, err := time.Parse("2006-01-02", fimDataStr)
	if err != nil {
		http.Error(w, "Formato de 'fimData' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	pontos, err := h.service.FiltrarDataNota(inicioData, fimData)
	if err != nil {
		log.Printf("Erro ao listar notas por data: %v", err)
		response := models.ResponseDefaultModel{
			IsSuccess:    false,
			Error:        err,
			ErrorMessage: "Erro ao listar notas por data",
		}
		//w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pontos,
	}

	//w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *NotaHandler) BuscarNotasPorNumero(w http.ResponseWriter, r *http.Request) {
	numero := r.URL.Query().Get("numero")
	if numero == "" {
		http.Error(w, "Parâmetro 'numero' é obrigatório", http.StatusBadRequest)
		return
	}

	notas, err := h.service.BuscarNotasPorNumero(numero)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      notas,
	}

	if err != nil {
		log.Printf("Erro ao buscar notas por número: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao buscar notas por número"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}
