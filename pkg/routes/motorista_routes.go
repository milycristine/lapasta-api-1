package routes

import (
	"net/http"

	dbsql "lapasta/database"
	motorista "lapasta/internal/Motorista"
)

/// RegisterMotoristaRoutes registra endpoints do módulo Motorista.
func RegisterMotoristaRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := motorista.NovoMotoristaRepository(db)
	svc := motorista.NovoMotoristaService(repo)
	handler := motorista.NovoMotoristaHandler(svc)

	mux.HandleFunc("/motorista", handler.CriarMotorista)
	mux.HandleFunc("/editarMotorista", handler.EditarMotorista)
	mux.HandleFunc("/statusmotorista", handler.AtualizarStatusMotorista)
	mux.HandleFunc("/listarMotoristas", handler.ListarMotoristas)
	mux.HandleFunc("/buscasMotoristaPorId", handler.BuscarMotoristaPorID)
	mux.HandleFunc("/emissaoNota", handler.CriarEmissaoNota)
	mux.HandleFunc("/listarEmissaoNotas", handler.ListarEmissaoNotas)
	mux.HandleFunc("/listarEmissaoNotasPorMotorista", handler.ListarEmissaoNotasPorMotorista)
	mux.HandleFunc("/buscarEmissao", handler.BuscarEmissaoNotas)
	mux.HandleFunc("/filtrarDataEmissao", handler.FiltrarDataEmissaoNota)

	mux.HandleFunc("/notaMotorista", handler.CriarNotaMotorista)
	mux.HandleFunc("/listarNotasPorMotorista", handler.ListarNotasPorMotorista)
	mux.HandleFunc("/atualizarStatusLancamentoNota", handler.AtualizarStatusLancamentoNotaMotorista)
	mux.HandleFunc("/filtrarDataNotaMotorista", handler.FiltrarNotasMotoristaPorData)
	mux.HandleFunc("/verificarSeLancouTodas", handler.MotoristaLancouTodasAsNotas)
	mux.HandleFunc("/pagamentoMotorista", handler.CriarPagamentoMotorista)
	mux.HandleFunc("/listarPagamentoMotorista", handler.ListarPagamentosMotorista)
	mux.HandleFunc("/calculoPagamento", handler.CalcularPagamentoMotorista)
	mux.HandleFunc("/atualizarStatusPagamentoMotorista", handler.AtualizarStatusPagamentoMotorista)
	mux.HandleFunc("/buscarMotoristaPorCPFouNome", handler.BuscarMotoristaPorCPFouNome)
}
