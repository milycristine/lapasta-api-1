package routes

import (
	"net/http"

	dbsql "lapasta/database"
	documento "lapasta/internal/Documento"
)

/// RegisterDocumentoRoutes registra endpoints relacionados a documentos.
func RegisterDocumentoRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := documento.NovoDocumentoRepository(db)
	svc := documento.NovoDocumentoService(repo)
	handler := documento.NovoDocumentoHandler(svc)

	mux.HandleFunc("/documento", handler.CriarDocumento)
	mux.HandleFunc("/listarDocumento", handler.ListarDocumentos)
	mux.HandleFunc("/filtrarDataDocumento", handler.FiltrarDataDocumento)
}
