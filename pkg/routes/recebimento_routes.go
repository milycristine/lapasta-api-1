package routes

import (
	"net/http"

	dbsql "lapasta/database"
	recebimento "lapasta/internal/Recebimento"
)

// / RegisterRecebimentoRoutes registra endpoints do módulo Recebimento.
func RegisterRecebimentoRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := recebimento.NovoRecebimentoRepository(db)
	svc := recebimento.NovoRecebimentoService(repo)
	handler := recebimento.NovoRecebimentoHandler(svc)

	mux.HandleFunc("/recebimento", handler.CriarRecebimento)
	mux.HandleFunc("/listarRecebimento", handler.ListarRecebimentos)
	mux.HandleFunc("/filtrarDataRecebimento", handler.FiltrarDataRecebimentos)
	mux.HandleFunc("/buscarPorNota", handler.BuscarDadosRecebimentoPorNumeroNota)
	mux.HandleFunc("/totaisAvista", handler.TotalRecebimentosAvistaMesAtual)
	mux.HandleFunc("/listarAvista", handler.ListarRecebimentosAvistaMesAtual)
}
