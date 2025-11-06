package routes

import (
	"net/http"

	dbsql "lapasta/database"
	boleto "lapasta/internal/Boletos"
)

/// RegisterBoletoRoutes registra endpoints do módulo Boletos.
func RegisterBoletoRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := boleto.NovoBoletoRepository(db)
	svc := boleto.NovoBoletoService(repo)
	handler := boleto.NovoBoletoHandler(svc)

	mux.HandleFunc("/boleto", handler.CriarBoletoRecebido)
	mux.HandleFunc("/listarPorPedido", handler.ListarBoletosPorRecebimento)
	mux.HandleFunc("/listarPorFornecedor", handler.ListarBoletosPorFornecedor)
	mux.HandleFunc("/boletodoDia", handler.ListarBoletosDoDia)
	mux.HandleFunc("/boletoAPagar", handler.PagarBoleto)
	mux.HandleFunc("/boletoPagos", handler.ListarBoletosPagos)
	mux.HandleFunc("/boletoVencidos", handler.ListarBoletosVencidos)
	mux.HandleFunc("/boletoPendentes", handler.ListarBoletosPendentes)
	mux.HandleFunc("/atualizarStatusBoleto", handler.AtualizarBoleto)
	mux.HandleFunc("/gerarRelatorio", handler.GerarEEnviarRelatorioBoletos)
	mux.HandleFunc("/boletostotais", handler.TotaisBoletos)
	mux.HandleFunc("/filtrarBoleto", handler.FiltrarBoletosPagosPorData)
}
