package server

import (
	"net/http"

	dbsql "lapasta/database"
	routes "lapasta/pkg/routes"
)

/// SetupRoutes registra os handlers de todos os módulos no mux.
/// - mux: *http.ServeMux no qual as rotas serão registradas.
/// - db: *dbsql.SQLStr, injetado para criação de repos/services/handlers.
func SetupRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	routes.RegisterAuthRoutes(mux, db)
	routes.RegisterRecebimentoRoutes(mux, db)
	routes.RegisterPontoRoutes(mux, db)
	routes.RegisterNotaRoutes(mux, db)
	routes.RegisterDocumentoRoutes(mux, db)
	routes.RegisterPagamentoRoutes(mux, db)
	routes.RegisterFuncionarioRoutes(mux, db)
	routes.RegisterValeRoutes(mux, db)
	routes.RegisterFornecedorRoutes(mux, db)
	routes.RegisterBoletoRoutes(mux, db)
	routes.RegisterGastosRoutes(mux, db)
	routes.RegisterMotoristaRoutes(mux, db)
	routes.RegisterTinyRoutes(mux, db)
}
