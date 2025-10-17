package recebimento

import (
	//"encoding/base64"
	"encoding/json"
	//Utils "lapasta/internal/Utils"
	"lapasta/internal/models"
	"log"
	"net/http"
	"strconv"

	//"strings"
	"time"
)

type RecebimentoHandler struct {
	service RecebimentoService
}

func NovoRecebimentoHandler(service RecebimentoService) *RecebimentoHandler {
	return &RecebimentoHandler{
		service: service,
	}
}

func (h *RecebimentoHandler) CriarRecebimento(w http.ResponseWriter, r *http.Request) {
	var recebimento models.Recebimento
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      recebimento,
	}

	if err := json.NewDecoder(r.Body).Decode(&recebimento); err != nil {
		log.Println(err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao decodificar o recebimento"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	valido, msg, err := h.service.ValidarRecebimento(&recebimento)
	if err != nil {
		log.Println("Erro na validação do recebimento:", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro interno na validação do recebimento"
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	if !valido {
		response.IsSuccess = false
		response.ErrorMessage = msg
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	//base64String := recebimento.ImagemBase64
	//if idx := strings.Index(base64String, ","); idx != -1 {
	//	base64String = base64String[idx+1:]
	//}
	//
	//fileBytes, err := base64.StdEncoding.DecodeString(base64String)
	//if err != nil {
	//	log.Println("Erro ao decodificar base64:", err)
	//	response.IsSuccess = false
	//	response.ErrorMessage = "Erro ao decodificar a imagem base64"
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Header().Set("Content-Type", "application/json")
	//	json.NewEncoder(w).Encode(response)
	//	return
	//}
	//extensao := ".png"
	//nomeImagem := Utils.GerarStringAleatoria(12) + extensao
	//urlImagem, err := Utils.UploadImagemFirebase(fileBytes, nomeImagem, "recebimentos")
	//if err != nil {
	//	log.Println("Erro ao salvar imagem local:", err)
	//	response.IsSuccess = false
	//	response.ErrorMessage = "Erro ao salvar imagem local"
	//	w.WriteHeader(http.StatusInternalServerError)
	//	w.Header().Set("Content-Type", "application/json")
	//	json.NewEncoder(w).Encode(response)
	//	return
	//}
	//recebimento.UrlImagem = urlImagem

	if err := h.service.CriarRecebimento(&recebimento); err != nil {
		log.Println(err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Data = recebimento
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RecebimentoHandler) ListarRecebimentos(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	log.Printf("Listando recebimentos na página: %d", pageInt)
	recebimentos, err := h.service.ListarRecebimentos(pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      recebimentos,
	}

	if err != nil {
		log.Printf("Erro ao listar recebimentos: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar recebimentos"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar a resposta: %v", err)
	}
}
func (h *RecebimentoHandler) FiltrarDataRecebimentos(w http.ResponseWriter, r *http.Request) {
	inicioDataStr := r.URL.Query().Get("inicioData")
	fimDataStr := r.URL.Query().Get("fimData")

	if inicioDataStr == "" || fimDataStr == "" {
		http.Error(w, "Os parâmetros 'inicioData' e 'fimData' são obrigatórios", http.StatusBadRequest)
		return
	}

	// Parse detalhado com log
	inicioData, err := time.Parse("2006-01-02", inicioDataStr)
	if err != nil {
		log.Printf("[Erro Parse inicioData] Valor: %s, Erro: %v", inicioDataStr, err)
		http.Error(w, "Formato de 'inicioData' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	fimData, err := time.Parse("2006-01-02", fimDataStr)
	if err != nil {
		log.Printf("[Erro Parse fimData] Valor: %s, Erro: %v", fimDataStr, err)
		http.Error(w, "Formato de 'fimData' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Ajustar hora para incluir o dia todo
	inicioData = inicioData.Truncate(24 * time.Hour)
	fimData = fimData.Truncate(24*time.Hour).AddDate(0, 0, 1).Add(-time.Nanosecond)

	log.Printf("[Info] Filtrando recebimentos de %s até %s", inicioData.Format(time.RFC3339), fimData.Format(time.RFC3339))

	pontos, err := h.service.FiltrarDataRecebimentos(inicioData, fimData)
	if err != nil {
		log.Printf("[Erro SQL] ao listar recebimentos: %v", err)
		response := models.ResponseDefaultModel{
			IsSuccess:    false,
			Error:        err,
			ErrorMessage: "Erro ao listar recebimentos por data",
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Log do resultado para depuração
	log.Printf("[Info] Recebimentos encontrados: %d", len(pontos))

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pontos,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("[Erro Encode JSON] %v", err)
	}
}

func (h *RecebimentoHandler) BuscarDadosRecebimentoPorNumeroNota(w http.ResponseWriter, r *http.Request) {
	numeroNota := r.URL.Query().Get("numeroNota")
	if numeroNota == "" {
		http.Error(w, "Parâmetro 'numeroNota' é obrigatório", http.StatusBadRequest)
		return
	}

	dados, err := h.service.BuscarDadosRecebimentoPorNumeroNota(numeroNota)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      dados,
	}

	if err != nil {
		log.Printf("Erro ao buscar dados por número da nota: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Nota não encontrada ou erro ao buscar dados"
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}
func (h *RecebimentoHandler) TotalRecebimentosAvistaMesAtual(w http.ResponseWriter, r *http.Request) {
	total, err := h.service.TotalRecebimentosAvistaMesAtual()
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      total,
	}

	if err != nil {
		log.Printf("Erro ao buscar total de pagamentos à vista: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao buscar total de pagamentos à vista"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar resposta: %v", err)
	}
}

func (h *RecebimentoHandler) ListarRecebimentosAvistaMesAtual(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	log.Printf("Listando pagamentos à vista (Pix, dinheiro, débito) do mês atual, página %d", pageInt)
	pagamentos, err := h.service.ListarRecebimentosAvistaMesAtual()
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pagamentos,
	}

	if err != nil {
		log.Printf("Erro ao listar pagamentos à vista: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar pagamentos à vista"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar resposta: %v", err)
	}
}
