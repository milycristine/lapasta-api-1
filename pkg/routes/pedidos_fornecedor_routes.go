package routes

import (
	"net/http"

	dbsql "lapasta/database"
	pedidofornecedor "lapasta/internal/PedidoFornecedor"
)

// / RegisterPedidosFornecedorRoutes registra endpoints relacionados a pedidos aos fornecedores.
func RegisterPedidosFornecedorRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := pedidofornecedor.NovoPedidoFornecedorRepository(db)
	svc := pedidofornecedor.NovoPedidoFornecedorService(repo)
	handler := pedidofornecedor.NovoPedidoFornecedorHandler(svc)

	mux.HandleFunc("/pedidoFornecedor", handler.CriarPedido)
	mux.HandleFunc("/listarPedido", handler.ListarPedidosPorFornecedor)
}
