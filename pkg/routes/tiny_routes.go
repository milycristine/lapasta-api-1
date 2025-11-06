package routes

import (
	"net/http"

	dbsql "lapasta/database"
	"lapasta/internal/tiny"
)

// / RegisterTinyRoutes registra endpoints utilitários do pacote tiny.
func RegisterTinyRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	mux.HandleFunc("/gerarNotasMotoristas", tiny.GerarNotasMotoristasHandler(db))
}
