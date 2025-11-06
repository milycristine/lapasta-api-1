package routes

import (
	"net/http"

	dbsql "lapasta/database"
	valetransporte "lapasta/internal/ValeTransporte"
)

// / RegisterValeRoutes registra endpoints relacionados ao vale transporte.
func RegisterValeRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := valetransporte.NovoValeRepository(db)
	svc := valetransporte.NovoValeService(repo)
	handler := valetransporte.NovoValeHandler(svc)

	mux.HandleFunc("/listarVale", handler.ListarVales)
	mux.HandleFunc("/atualizarStatusVale", handler.AtualizarVale)
	mux.HandleFunc("/listarValePorSemana", handler.ListarValesDaSemana)
}
