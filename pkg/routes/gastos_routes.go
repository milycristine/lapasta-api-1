package routes

import (
	"net/http"

	dbsql "lapasta/database"
	//gastos "lapasta/internal/Gastos"
)

// / RegisterGastosRoutes registra endpoints relacionados a gastos.
func RegisterGastosRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	//repo := gastos.NovoGastosRepository(db)
	//svc := gastos.NovoGastosService(repo)
	//handler := gastos.NovoGastosHandler(svc)

	//mux.HandleFunc("/totalGastos", handler.TotalGastos)
	//mux.HandleFunc("/listarGastos", handler.ListarGastos)
}
