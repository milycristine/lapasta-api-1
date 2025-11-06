package routes

import (
	"net/http"

	dbsql "lapasta/database"
	ponto "lapasta/internal/Ponto"
)

/// RegisterPontoRoutes registra endpoints relacionados a ponto.
func RegisterPontoRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := ponto.NovoPontoRepository(db)
	svc := ponto.NovoPontoService(repo)
	handler := ponto.NovoPontoHandler(svc)

	mux.HandleFunc("/listarPonto", handler.ListarPontos)
	mux.HandleFunc("/listarPontoId", handler.ListarPontosPorId)
	mux.HandleFunc("/listarPontoIdEDia", handler.ListarPontosPorIdEDia)
	mux.HandleFunc("/listarPontoPorData", handler.ListarPontosPorData)
	mux.HandleFunc("/listarPontoPorDataId", handler.ListarPontosPorDataId)
	mux.HandleFunc("/horaEntrada", handler.RegistrarEntrada)
	mux.HandleFunc("/horaSaidaAlmoco", handler.RegistrarSaidaAlmoco)
	mux.HandleFunc("/horaRetornoAlmoco", handler.RegistrarRetornoAlmoco)
	mux.HandleFunc("/horaSaida", handler.RegistrarSaida)
	mux.HandleFunc("/relatorioExcel", handler.GerarRelatorioMensal)
}
