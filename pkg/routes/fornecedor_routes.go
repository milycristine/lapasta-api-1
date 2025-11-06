package routes

import (
	"net/http"

	dbsql "lapasta/database"
	fornecedor "lapasta/internal/Fornecedores"
)

/// RegisterFornecedorRoutes registra endpoints relacionados a fornecedores.
func RegisterFornecedorRoutes(mux *http.ServeMux, db  *dbsql.SQLStr) {
	repo := fornecedor.NovoFornecedorRepository(db)
	svc := fornecedor.NovoFornecedorService(repo)
	handler := fornecedor.NovoFornecedorHandler(svc)

	mux.HandleFunc("/fornecedores", handler.CriarFornecedor)
	mux.HandleFunc("/editarFornecedores", handler.EditarFornecedor)
	mux.HandleFunc("/listarFornecedores", handler.ListarFornecedores)
	mux.HandleFunc("/fornecedores/buscar", handler.BuscarFornecedorPorCNPJouNome)
	mux.HandleFunc("/buscarPedidoId", handler.BuscarPedidosFornecedorPorDescricaoOuId)
}
