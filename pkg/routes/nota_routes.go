package routes

import (
	nota "lapasta/internal/Notas"

	"net/http"
		dbsql "lapasta/database"

)

/// RegisterNotaRoutes registra endpoints para notas.
func RegisterNotaRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := nota.NovoNotaRepository(db)
	svc := nota.NovaNotaService(repo)
	handler := nota.NovaNotaHandler(svc)

	mux.HandleFunc("/nota", handler.CriarNota)
	mux.HandleFunc("/listarNotas", handler.ListarNotas)
	mux.HandleFunc("/filtrarDataNota", handler.FiltrarDataNota)
	mux.HandleFunc("/buscarNota", handler.BuscarNotasPorNumero)
}
