package motorista

import (
	"encoding/base64"
	"encoding/json"
	Utils "lapasta/internal/Utils"
	"lapasta/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MotoristaHandler interface {
	CriarMotorista(w http.ResponseWriter, r *http.Request)
	EditarMotorista(w http.ResponseWriter, r *http.Request)
	AtualizarStatusMotorista(w http.ResponseWriter, r *http.Request)
	ListarMotoristas(w http.ResponseWriter, r *http.Request)
	BuscarMotoristaPorCPFouNome(w http.ResponseWriter, r *http.Request)
	BuscarMotoristaPorID(w http.ResponseWriter, r *http.Request)

	CriarEmissaoNota(w http.ResponseWriter, r *http.Request)
	ListarEmissaoNotas(w http.ResponseWriter, r *http.Request)
	ListarEmissaoNotasPorMotorista(w http.ResponseWriter, r *http.Request)
	BuscarEmissaoNotas(w http.ResponseWriter, r *http.Request)
	FiltrarDataEmissaoNota(w http.ResponseWriter, r *http.Request)

	CriarNotaMotorista(w http.ResponseWriter, r *http.Request)
	ListarNotasPorMotorista(w http.ResponseWriter, r *http.Request)
	AtualizarStatusLancamentoNotaMotorista(w http.ResponseWriter, r *http.Request)
	MotoristaLancouTodasAsNotas(w http.ResponseWriter, r *http.Request)
	FiltrarNotasMotoristaPorData(w http.ResponseWriter, r *http.Request)

	CriarPagamentoMotorista(w http.ResponseWriter, r *http.Request)
	ListarPagamentosMotorista(w http.ResponseWriter, r *http.Request)
	AtualizarStatusPagamentoMotorista(w http.ResponseWriter, r *http.Request)
	CalcularPagamentoMotorista(w http.ResponseWriter, r *http.Request)
}

type motoristaHandler struct {
	service MotoristaService
}

func NovoMotoristaHandler(service MotoristaService) MotoristaHandler {
	return &motoristaHandler{
		service: service,
	}
}
func (h *motoristaHandler) CriarMotorista(w http.ResponseWriter, r *http.Request) {
	var motorista models.Motorista
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      motorista,
	}

	if err := json.NewDecoder(r.Body).Decode(&motorista); err != nil {
		log.Println("Erro ao decodificar motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao decodificar motorista"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if motorista.ImagemFrenteBase64 != "" {
		base64Str := motorista.ImagemFrenteBase64
		if idx := strings.Index(base64Str, ","); idx != -1 {
			base64Str = base64Str[idx+1:]
		}

		fileBytes, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			log.Println("Erro ao decodificar imagem frente CNH:", err)
			response.IsSuccess = false
			response.ErrorMessage = "Erro ao decodificar imagem frente CNH"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		nomeImagem := Utils.GerarStringAleatoria(12) + ".png"
		url, err := Utils.UploadImagemFirebase(fileBytes, nomeImagem, "cnh_motoristas")
		if err != nil {
			log.Println("Erro ao salvar imagem frente CNH:", err)
			response.IsSuccess = false
			response.ErrorMessage = "Erro ao salvar imagem frente CNH"
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		motorista.CnhFrenteUrl = url
	}

	if motorista.ImagemVersoBase64 != "" {
		base64Str := motorista.ImagemVersoBase64
		if idx := strings.Index(base64Str, ","); idx != -1 {
			base64Str = base64Str[idx+1:]
		}

		fileBytes, err := base64.StdEncoding.DecodeString(base64Str)
		if err != nil {
			log.Println("Erro ao decodificar imagem verso CNH:", err)
			response.IsSuccess = false
			response.ErrorMessage = "Erro ao decodificar imagem verso CNH"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		nomeImagem := Utils.GerarStringAleatoria(12) + ".png"
		url, err := Utils.UploadImagemFirebase(fileBytes, nomeImagem, "cnh_motoristas")
		if err != nil {
			log.Println("Erro ao salvar imagem verso CNH:", err)
			response.IsSuccess = false
			response.ErrorMessage = "Erro ao salvar imagem verso CNH"
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}
		motorista.CnhVersoUrl = url
	}

	if err := h.service.CriarMotorista(&motorista); err != nil {
		log.Println("Erro ao criar motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Data = motorista
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) AtualizarStatusMotorista(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	ativoStr := r.URL.Query().Get("ativo")

	var response models.ResponseDefaultModel

	id, err := strconv.Atoi(idStr)
	if err != nil || (ativoStr != "true" && ativoStr != "false") {
		response.IsSuccess = false
		response.ErrorMessage = "Parâmetros inválidos (use ?id=1&ativo=true/false)"
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	ativo := ativoStr == "true"

	err = h.service.AtualizarStatusMotorista(id, ativo)
	if err != nil {
		log.Println("Erro ao atualizar status do motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao atualizar status do motorista"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.IsSuccess = true
		response.Data = map[string]interface{}{
			"id":    id,
			"ativo": ativo,
		}
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) EditarMotorista(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var motorista models.Motorista
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      motorista,
	}

	if err := json.NewDecoder(r.Body).Decode(&motorista); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao decodificar os dados do funcionário"
		w.WriteHeader(http.StatusBadRequest)
	} else if err := h.service.EditarMotorista(&motorista); err != nil {
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) ListarMotoristas(w http.ResponseWriter, r *http.Request) {
	motoristas, err := h.service.ListarMotoristas()
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      motoristas,
	}

	if err != nil {
		log.Println("Erro ao listar motoristas:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao listar motoristas"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) BuscarMotoristaPorCPFouNome(w http.ResponseWriter, r *http.Request) {
	valor := r.URL.Query().Get("valor")
	if valor == "" {
		http.Error(w, "CPF ou Nome é obrigatório", http.StatusBadRequest)
		return
	}

	motoristas, err := h.service.BuscarMotoristaPorCPFouNome(valor)
	if err != nil {
		http.Error(w, "Erro ao buscar motoristas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(motoristas)
}

func (h *motoristaHandler) BuscarMotoristaPorID(w http.ResponseWriter, r *http.Request) {
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

	motorista, err := h.service.BuscarMotoristaPorID(id)
	if err != nil {
		http.Error(w, "Erro ao buscar motorista: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if motorista == nil {
		http.Error(w, "Motorista não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(motorista)
}

func (h *motoristaHandler) CriarEmissaoNota(w http.ResponseWriter, r *http.Request) {
	var nota models.EmissaoNota
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      nota,
	}

	if err := json.NewDecoder(r.Body).Decode(&nota); err != nil {
		log.Println("Erro ao decodificar emissão de nota:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao decodificar emissão de nota"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := h.service.CriarEmissaoNota(&nota); err != nil {
		log.Println("Erro ao criar emissão de nota:", err)
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Data = nota
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) ListarEmissaoNotas(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt <= 0 {
		http.Error(w, "Página inválida", http.StatusBadRequest)
		return
	}

	emissoes, err := h.service.ListarEmissaoNotas(pageInt)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      emissoes,
	}

	if err != nil {
		log.Println("Erro ao listar emissão de notas:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao listar emissão de notas"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) ListarEmissaoNotasPorMotorista(w http.ResponseWriter, r *http.Request) {
	idMotoristaStr := r.URL.Query().Get("id_motorista")

	idMotorista, err := strconv.Atoi(idMotoristaStr)
	if err != nil || idMotorista <= 0 {
		http.Error(w, "ID do motorista inválido", http.StatusBadRequest)
		return
	}

	emissoes, err := h.service.ListarEmissaoNotasPorMotorista(idMotorista)

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      emissoes,
	}

	if err != nil {
		log.Println("Erro ao listar emissão de notas por motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao listar emissão de notas por motorista"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}
func (h *motoristaHandler) BuscarEmissaoNotas(w http.ResponseWriter, r *http.Request) {
	valor := r.URL.Query().Get("valor")
	if valor == "" {
		http.Error(w, "NumeroNota, Descrição ou Nome é obrigatório", http.StatusBadRequest)
		return
	}

	motoristas, err := h.service.BuscarEmissaoNotas(valor)
	if err != nil {
		http.Error(w, "Erro ao buscar motoristas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(motoristas)
}
func (h *motoristaHandler) FiltrarDataEmissaoNota(w http.ResponseWriter, r *http.Request) {
	inicioDataStr := r.URL.Query().Get("inicioData")
	fimDataStr := r.URL.Query().Get("fimData")

	if inicioDataStr == "" || fimDataStr == "" {
		http.Error(w, "Os parâmetros 'inicioData' e 'fimData' são obrigatórios", http.StatusBadRequest)
		return
	}

	inicio, err := time.Parse("2006-01-02", inicioDataStr)
	if err != nil {
		http.Error(w, "Formato de 'inicioData' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	fim, err := time.Parse("2006-01-02", fimDataStr)
	if err != nil {
		http.Error(w, "Formato de 'fimData' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	notas, err := h.service.FiltrarDataEmissaoNota(inicio, fim)
	if err != nil {
		log.Printf("Erro ao filtrar emissão de notas por data: %v", err)
		response := models.ResponseDefaultModel{
			IsSuccess:    false,
			Error:        err,
			ErrorMessage: "Erro ao filtrar emissão de notas por data",
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      notas,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) CriarNotaMotorista(w http.ResponseWriter, r *http.Request) {
	var nota models.NotasMotorista
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      nota,
	}

	if err := json.NewDecoder(r.Body).Decode(&nota); err != nil {
		log.Println("Erro ao decodificar a nota:", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao decodificar a nota"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	base64String := nota.ImagemBase64
	if idx := strings.Index(base64String, ","); idx != -1 {
		base64String = base64String[idx+1:]
	}

	fileBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		log.Println("Erro ao decodificar base64 da imagem:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao decodificar a imagem base64"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	extensao := ".png"
	nomeImagem := Utils.GerarStringAleatoria(12) + extensao
	url, err := Utils.UploadImagemFirebase(fileBytes, nomeImagem, "notas_motorista")
	if err != nil {
		log.Println("Erro ao salvar imagem:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao salvar imagem"
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	nota.Url = url

	if err := h.service.CriarNotasMotorista(&nota); err != nil {
		log.Println("Erro ao criar nota motorista no banco:", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Data = nota
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) ListarNotasPorMotorista(w http.ResponseWriter, r *http.Request) {
	idMotoristaStr := r.URL.Query().Get("id_motorista")
	if idMotoristaStr == "" {
		http.Error(w, "O parâmetro 'id_motorista' é obrigatório", http.StatusBadRequest)
		return
	}

	idMotorista, err := strconv.Atoi(idMotoristaStr)
	if err != nil || idMotorista <= 0 {
		http.Error(w, "O parâmetro 'id_motorista' é inválido", http.StatusBadRequest)
		return
	}

	notas, err := h.service.ListarNotasMotoristaPorMotorista(idMotorista)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      notas,
	}

	if err != nil {
		log.Println("Erro ao listar notas do motorista:", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Erro ao listar notas"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) AtualizarStatusLancamentoNotaMotorista(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id             int        `json:"id"`
		StatusId       int        `json:"statusId"`
		DataLancamento *time.Time `json:"datalancamento"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if req.Id <= 0 || req.StatusId <= 0 {
		http.Error(w, "Parâmetros inválidos", http.StatusBadRequest)
		return
	}

	dataLancamento := req.DataLancamento
	if dataLancamento == nil {
		now := time.Now()
		dataLancamento = &now
	}

	err = h.service.AtualizarStatusLancamentoNotaMotorista(req.Id, req.StatusId, dataLancamento)
	response := models.ResponseDefaultModel{IsSuccess: true}

	if err != nil {
		log.Println("Erro ao atualizar status da nota do motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) MotoristaLancouTodasAsNotas(w http.ResponseWriter, r *http.Request) {
	idMotoristaStr := r.URL.Query().Get("id_motorista")
	inicioStr := r.URL.Query().Get("inicio")
	fimStr := r.URL.Query().Get("fim")

	idMotorista, err := strconv.Atoi(idMotoristaStr)
	if err != nil || idMotorista <= 0 {
		http.Error(w, "Parâmetro 'id_motorista' inválido", http.StatusBadRequest)
		return
	}

	inicio, err := time.Parse("2006-01-02", inicioStr)
	if err != nil {
		http.Error(w, "Parâmetro 'inicio' inválido", http.StatusBadRequest)
		return
	}

	fim, err := time.Parse("2006-01-02", fimStr)
	if err != nil {
		http.Error(w, "Parâmetro 'fim' inválido", http.StatusBadRequest)
		return
	}

	lancouTodas, err := h.service.MotoristaLancouTodasAsNotas(idMotorista, inicio, fim)

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      lancouTodas,
	}

	if err != nil {
		log.Println("Erro ao verificar se motorista lançou todas as notas:", err)
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) CriarPagamentoMotorista(w http.ResponseWriter, r *http.Request) {
	var pagamento models.PagamentosMotorista
	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pagamento,
	}

	if err := json.NewDecoder(r.Body).Decode(&pagamento); err != nil {
		log.Println("Erro ao decodificar pagamento motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao decodificar pagamento"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := h.service.CriarPagamentoMotorista(&pagamento); err != nil {
		log.Println("Erro ao criar pagamento motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Data = pagamento
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) ListarPagamentosMotorista(w http.ResponseWriter, r *http.Request) {
	idMotoristaStr := r.URL.Query().Get("id_motorista")

	var idMotorista int
	var err error

	if idMotoristaStr != "" {
		idMotorista, err = strconv.Atoi(idMotoristaStr)
		if err != nil || idMotorista < 0 {
			http.Error(w, "Parâmetro 'id_motorista' inválido", http.StatusBadRequest)
			return
		}
	}

	pagamentos, err := h.service.ListarPagamentosMotorista(idMotorista)

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pagamentos,
	}

	if err != nil {
		log.Println("Erro ao listar pagamentos do motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) AtualizarStatusPagamentoMotorista(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	statusIdStr := r.URL.Query().Get("status_id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Parâmetro 'id' inválido", http.StatusBadRequest)
		return
	}

	statusId, err := strconv.Atoi(statusIdStr)
	if err != nil {
		http.Error(w, "Parâmetro 'status_id' inválido", http.StatusBadRequest)
		return
	}

	err = h.service.AtualizarStatusPagamentoMotorista(id, statusId)
	response := models.ResponseDefaultModel{
		IsSuccess: true,
	}

	if err != nil {
		log.Println("Erro ao atualizar status do pagamento do motorista:", err)
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(response)
}
func (h *motoristaHandler) FiltrarNotasMotoristaPorData(w http.ResponseWriter, r *http.Request) {
	idMotoristaStr := r.URL.Query().Get("idMotorista")
	inicioDataStr := r.URL.Query().Get("inicioData")
	fimDataStr := r.URL.Query().Get("fimData")

	if idMotoristaStr == "" {
		http.Error(w, "O parâmetro 'idMotorista' é obrigatório", http.StatusBadRequest)
		return
	}

	idMotorista, err := strconv.Atoi(idMotoristaStr)
	if err != nil || idMotorista <= 0 {
		http.Error(w, "O parâmetro 'idMotorista' deve ser um número inteiro válido", http.StatusBadRequest)
		return
	}

	var inicioData, fimData *time.Time

	if inicioDataStr != "" && fimDataStr != "" {
		inicio, err := time.Parse("2006-01-02", inicioDataStr)
		if err != nil {
			http.Error(w, "Formato de 'inicioData' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		fim, err := time.Parse("2006-01-02", fimDataStr)
		if err != nil {
			http.Error(w, "Formato de 'fimData' inválido. Use o formato YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		inicioData = &inicio
		fimData = &fim
	}

	notas, err := h.service.FiltrarNotasMotoristaPorData(idMotorista, inicioData, fimData)
	if err != nil {
		log.Printf("Erro ao filtrar notas motorista por data: %v", err)
		response := models.ResponseDefaultModel{
			IsSuccess:    false,
			Error:        err,
			ErrorMessage: "Erro ao filtrar notas motorista por data",
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      notas,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *motoristaHandler) CalcularPagamentoMotorista(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	idMotoristaStr := r.URL.Query().Get("idMotorista")
	inicioStr := r.URL.Query().Get("inicio")
	fimStr := r.URL.Query().Get("fim")

	if idMotoristaStr == "" || inicioStr == "" || fimStr == "" {
		http.Error(w, "Parâmetros 'idMotorista', 'inicio' e 'fim' são obrigatórios", http.StatusBadRequest)
		return
	}

	idMotorista, err := strconv.Atoi(idMotoristaStr)
	if err != nil || idMotorista <= 0 {
		http.Error(w, "IdMotorista inválido", http.StatusBadRequest)
		return
	}

	inicio, err := time.Parse("2006-01-02", inicioStr)
	if err != nil {
		http.Error(w, "Data 'inicio' inválida. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	fim, err := time.Parse("2006-01-02", fimStr)
	if err != nil {
		http.Error(w, "Data 'fim' inválida. Use o formato YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	log.Printf("Calculando pagamento do motorista %d, de %s até %s", idMotorista, inicio.Format("02/01"), fim.Format("02/01"))

	pagamento, err := h.service.CalcularPagamentoMotorista(idMotorista, inicio, fim)

	response := models.ResponseDefaultModel{
		IsSuccess: true,
		Data:      pagamento,
	}

	if err != nil {
		log.Printf("Erro ao calcular pagamento: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Erro ao codificar resposta JSON: %v", err)
	}
}
