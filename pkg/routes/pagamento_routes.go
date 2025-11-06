package routes

import (
	"net/http"

	dbsql "lapasta/database"
	pagamento "lapasta/internal/Pagamento"
)

/// RegisterPagamentoRoutes registra endpoints relacionados ao pagamentos dos funcionarios.
func RegisterPagamentoRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := pagamento.NovoPagamentoRepository(db)
	svc := pagamento.NovoPagamentoService(repo)
	handler := pagamento.NovoPagamentoHandler(svc)

	mux.HandleFunc("/pagamento", handler.CriarPagamento)
	mux.HandleFunc("/listarPagamento", handler.ListarPagamentos)
	mux.HandleFunc("/listarPagamentoPorDia", handler.ListarPagamentosPorDia)
	mux.HandleFunc("/atualizarStatusPagamento", handler.AtualizarPagamento)
	mux.HandleFunc("/listarPagamentoPorMes", handler.ListarPagamentosPorMes)
}
