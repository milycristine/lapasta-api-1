package documento

import (
	"encoding/base64"
	"encoding/json"

	//"lapasta/config"
	"lapasta/internal/models"
	//"lapasta/internal/s3client"
	Utils "lapasta/internal/Utils"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DocumentoHandler interface {
	CriarDocumento(w http.ResponseWriter, r *http.Request)
	ListarDocumentos(w http.ResponseWriter, r *http.Request)
	FiltrarDataDocumento(w http.ResponseWriter, r *http.Request)
}

type documentoHandler struct {
	service DocumentoService
}

func NovoDocumentoHandler(service DocumentoService) DocumentoHandler {
	return &documentoHandler{
		service: service,
	}
}

func (h *documentoHandler) CriarDocumento(w http.ResponseWriter, r *http.Request) {
	var documento models.Documento
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      documento,
	}

	if err := json.NewDecoder(r.Body).Decode(&documento); err != nil {
		log.Println("Erro ao decodificar o documento:", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao decodificar o documento"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	base64String := documento.ImagemBase64
	if idx := strings.Index(base64String, ","); idx != -1 {
		base64String = base64String[idx+1:]
	}

	fileBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		log.Println("Erro ao decodificar base64 do documento:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao decodificar o documento base64"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	//s3Client := s3client.NovoS3Client(config.Yml.AWS.Region, config.Yml.AWS.BucketName)
	//url, err := s3Client.UploadBase64File(fileBytes)
	extensao := ".png"
	nomeImagem := Utils.GerarStringAleatoria(12) + extensao
	url, err := Utils.UploadImagemFirebase(fileBytes, nomeImagem, "documentos")
	if err != nil {
		log.Println("Erro ao salvar imagem local:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao salvar imagem local"
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	documento.Url = url

	if err := h.service.CriarDocumento(&documento); err != nil {
		log.Println("Erro ao criar o documento no banco de dados:", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Data = documento
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}
func (h *documentoHandler) ListarDocumentos(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	log.Printf("Listando documentos na página: %d", pageInt)
	documentos, err := h.service.ListarDocumentos(pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      documentos,
	}

	if err != nil {
		log.Printf("Erro ao listar documentos: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar documentos"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}

func (h *documentoHandler) FiltrarDataDocumento(w http.ResponseWriter, r *http.Request) {
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

	pontos, err := h.service.FiltrarDataDocumento(inicioData, fimData)
	if err != nil {
		log.Printf("Erro ao listar documentos por data: %v", err)
		response := models.ResponseDefaultModel{
			IsSuccess:    false,
			Error:        err,
			ErrorMessage: "Erro ao listar documentos por data",
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
