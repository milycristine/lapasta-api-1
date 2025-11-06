package routes

import (
	"net/http"

	dbsql "lapasta/database"
	funcionario "lapasta/internal/Funcionario"
)

// / RegisterFuncionarioRoutes registra endpoints relacionados a funcionarios.
func RegisterFuncionarioRoutes(mux *http.ServeMux, db *dbsql.SQLStr) {
	repo := funcionario.NovoFuncionarioRepository(db)
	svc := funcionario.NovoFuncionarioService(repo)
	handler := funcionario.NovoFuncionarioHandler(svc)

	mux.HandleFunc("/funcionario", handler.CriarFuncionario)
	mux.HandleFunc("/editarFuncionario", handler.EditarFuncionario)
	mux.HandleFunc("/listarFuncionario", handler.ListarFuncionarios)
	mux.HandleFunc("/buscarFuncionario", handler.BuscarFuncionarioPorCPF)
	mux.HandleFunc("/buscarFuncionarioPorId", handler.BuscarFuncionarioPorID)
	mux.HandleFunc("/statusFuncionario", handler.AtualizarStatusFuncionario)

}
