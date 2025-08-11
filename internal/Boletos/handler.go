package boleto

import (
	"encoding/json"
	"lapasta/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BoletoHandler interface {
	CriarBoleto(w http.ResponseWriter, r *http.Request)
	ListarBoletosPorFornecedor(w http.ResponseWriter, r *http.Request)
	ListarBoletosPorPedido(w http.ResponseWriter, r *http.Request)
	ListarBoletosDoDia(w http.ResponseWriter, r *http.Request)
	PagarBoleto(w http.ResponseWriter, r *http.Request)
	ListarBoletosPagos(w http.ResponseWriter, r *http.Request)
	ListarBoletosVencidos(w http.ResponseWriter, r *http.Request)
	ListarBoletosPendentes(w http.ResponseWriter, r *http.Request)
	AtualizarBoleto(w http.ResponseWriter, r *http.Request)
	GerarEEnviarRelatorioBoletos(w http.ResponseWriter, r *http.Request)
	TotaisBoletos(w http.ResponseWriter, r *http.Request)
}

type boletoHandler struct {
	service BoletoService
}

func NovoBoletoHandler(service BoletoService) BoletoHandler {
	return &boletoHandler{
		service: service,
	}
}

func (h *boletoHandler) CriarBoleto(w http.ResponseWriter, r *http.Request) {
	var boleto models.Boleto
	response := models.ResponseDefaultModel{IsSuccess: true}

	if err := json.NewDecoder(r.Body).Decode(&boleto); err != nil {
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao decodificar o boleto"
		w.WriteHeader(http.StatusBadRequest)
	} else if err := h.service.CriarBoleto(&boleto); err != nil {
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao criar boleto"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Data = boleto
		w.WriteHeader(http.StatusCreated)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *boletoHandler) ListarBoletosPorFornecedor(w http.ResponseWriter, r *http.Request) {
	valor := r.URL.Query().Get("fornecedorId")
	response := models.ResponseDefaultModel{IsSuccess: true}

	id, err := strconv.Atoi(valor)
	if err != nil {
		response.IsSuccess = false
		response.ErrorMessage = "ID inválido"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		boletos, err := h.service.ListarBoletosPorFornecedor(id)
		if err != nil {
			response.IsSuccess = false
			response.ErrorMessage = "Erro ao listar boletos por fornecedor"
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Data = boletos
			w.WriteHeader(http.StatusOK)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *boletoHandler) ListarBoletosPorPedido(w http.ResponseWriter, r *http.Request) {
	valor := r.URL.Query().Get("pedidoId")
	response := models.ResponseDefaultModel{IsSuccess: true}

	id, err := strconv.Atoi(valor)
	if err != nil {
		response.IsSuccess = false
		response.ErrorMessage = "ID inválido"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		boletos, err := h.service.ListarBoletosPorPedido(id)
		if err != nil {
			response.IsSuccess = false
			response.ErrorMessage = "Erro ao listar boletos por pedido"
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Data = boletos
			w.WriteHeader(http.StatusOK)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *boletoHandler) ListarBoletosDoDia(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query().Get("data")
	response := models.ResponseDefaultModel{IsSuccess: true}

	if data == "" {
		response.IsSuccess = false
		response.ErrorMessage = "Data é obrigatória"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		boletos, err := h.service.ListarBoletosDoDia(data)
		if err != nil {
			response.IsSuccess = false
			response.ErrorMessage = "Erro ao listar boletos do dia"
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			response.Data = boletos
			w.WriteHeader(http.StatusOK)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *boletoHandler) PagarBoleto(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CodigoBarras string `json:"codigo_barras"`
	}
	response := models.ResponseDefaultModel{IsSuccess: true}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.CodigoBarras == "" {
		response.IsSuccess = false
		response.ErrorMessage = "Código de barras é obrigatório"
		w.WriteHeader(http.StatusBadRequest)
	} else if err := h.service.PagarBoleto(body.CodigoBarras); err != nil {
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao pagar boleto"
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Data = "Boleto pago com sucesso"
		w.WriteHeader(http.StatusOK)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *boletoHandler) ListarBoletosPagos(w http.ResponseWriter, r *http.Request) {
	boletos, err := h.service.ListarBoletosPagos()
	if err != nil {
		http.Error(w, "Erro ao listar boletos pagos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boletos)
}

func (h *boletoHandler) ListarBoletosVencidos(w http.ResponseWriter, r *http.Request) {
	boletos, err := h.service.ListarBoletosVencidos()
	if err != nil {
		http.Error(w, "Erro ao listar boletos vencidos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boletos)
}

func (h *boletoHandler) ListarBoletosPendentes(w http.ResponseWriter, r *http.Request) {
	boletos, err := h.service.ListarBoletosPendentes()
	if err != nil {
		http.Error(w, "Erro ao listar boletos pendentes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(boletos)
}

func (h *boletoHandler) AtualizarBoleto(w http.ResponseWriter, r *http.Request) {
	var boleto models.Boleto
	response := models.ResponseDefaultModel{
		IsSuccess: true,
	}

	if err := json.NewDecoder(r.Body).Decode(&boleto); err != nil {
		log.Printf("Erro ao decodificar o boleto: %v", err)
		response.IsSuccess = false
		response.Error = err
		response.ErrorMessage = "Formato de entrada inválido"
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if err := h.service.AtualizarStatusBoleto(boleto.Id, boleto.StatusId); err != nil {
			log.Printf("Erro ao atualizar O boleto: %v", err)
			response.IsSuccess = false
			response.Error = err
			response.ErrorMessage = "Erro ao atualizar o boleto"
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusAccepted)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func (h *boletoHandler) GerarEEnviarRelatorioBoletos(w http.ResponseWriter, r *http.Request) {
	response := models.ResponseDefaultModel{IsSuccess: true}

	emailAdmin := r.URL.Query().Get("emailAdmin")
	anoStr := r.URL.Query().Get("ano")
	mesStr := r.URL.Query().Get("mes")
	fornecedorIDStr := r.URL.Query().Get("fornecedorId")
	statusIDsStr := r.URL.Query().Get("statusIds")

	if emailAdmin == "" {
		response.IsSuccess = false
		response.ErrorMessage = "Parâmetro 'emailAdmin' é obrigatório"
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, response)
		return
	}

	ano, _ := strconv.Atoi(anoStr)
	mes, _ := strconv.Atoi(mesStr)

	var fornecedorID *int
	if fornecedorIDStr != "" {
		id, err := strconv.Atoi(fornecedorIDStr)
		if err != nil {
			response.IsSuccess = false
			response.ErrorMessage = "Parâmetro 'fornecedorId' inválido"
			w.WriteHeader(http.StatusBadRequest)
			writeJSON(w, response)
			return
		}
		fornecedorID = &id
	}

	var statusIDs []int
	if statusIDsStr != "" {
		for _, s := range strings.Split(statusIDsStr, ",") {
			id, err := strconv.Atoi(strings.TrimSpace(s))
			if err == nil {
				statusIDs = append(statusIDs, id)
			}
		}
	}

	err := h.service.GerarEEnviarRelatorioBoletos(emailAdmin, ano, mes, fornecedorID, statusIDs)
	if err != nil {
		response.IsSuccess = false
		response.ErrorMessage = "Erro ao gerar e enviar relatório: " + err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Data = "Relatório gerado e enviado com sucesso"
		w.WriteHeader(http.StatusOK)
	}

	writeJSON(w, response)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *boletoHandler) TotaisBoletos(w http.ResponseWriter, r *http.Request) {
	hoje := time.Now().Format("2006-01-02")

	totalDia, err := h.service.TotalBoletosDia(hoje)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalAtrasados, err := h.service.TotalBoletosAtrasados()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalMesPendentes, err := h.service.TotalBoletosPendentesMesAtual()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalMesPagos, err := h.service.TotalBoletosPagosMesAtual()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]float64{
		"pendentesHoje": totalDia,
		"atrasados":     totalAtrasados,
		"pendentesMes":  totalMesPendentes,
		"pagosMes":      totalMesPagos,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
