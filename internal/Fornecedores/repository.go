package fornecedor

import (
	database "lapasta/database"
	"lapasta/internal/models"
)

type FornecedorRepository interface {
	CriarFornecedor(fornecedor *models.Fornecedor) error
	EditarFornecedor(f *models.Fornecedor) error
	ListarFornecedores(page int) ([]models.Fornecedor, error)
	BuscarFornecedorPorCNPJouNome(valor string) ([]models.Fornecedor, error)
	CriarPedidoFornecedor(pedido *models.PedidoFornecedor) error
	ListarPedidosPorFornecedor(fornecedorId int) ([]models.PedidoFornecedor, error)
	BuscarPedidosFornecedorPorDescricaoOuId(fornecedorId int, valor string) ([]models.PedidoFornecedor, error) // <-- NOVO

}

type fornecedorRepository struct {
	db *database.SQLStr
}

func NovoFornecedorRepository(db *database.SQLStr) FornecedorRepository {
	return &fornecedorRepository{
		db: db,
	}
}

func (r *fornecedorRepository) CriarFornecedor(fornecedor *models.Fornecedor) error {
	return r.db.CriarFornecedor(fornecedor)
}
func (r *fornecedorRepository) EditarFornecedor(f *models.Fornecedor) error  {
	return r.db.EditarFornecedor(f)
}

func (r *fornecedorRepository) ListarFornecedores(page int) ([]models.Fornecedor, error) {
	return r.db.ListarFornecedores(page)
}

func (r *fornecedorRepository) BuscarFornecedorPorCNPJouNome(valor string) ([]models.Fornecedor, error) {
	return r.db.BuscarFornecedorPorCNPJouNome(valor)
}
func (r *fornecedorRepository) CriarPedidoFornecedor(pedido *models.PedidoFornecedor) error {
	return r.db.CriarPedidoFornecedor(pedido)
}

func (r *fornecedorRepository) ListarPedidosPorFornecedor(fornecedorId int) ([]models.PedidoFornecedor, error) {
	return r.db.ListarPedidosFornecedor(fornecedorId)
}
func (r *fornecedorRepository) BuscarPedidosFornecedorPorDescricaoOuId(fornecedorId int, valor string) ([]models.PedidoFornecedor, error) {
	return r.db.BuscarPedidosFornecedorPorDescricaoOuId(fornecedorId, valor)
}
