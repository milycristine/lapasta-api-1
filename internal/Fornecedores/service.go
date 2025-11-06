package fornecedor

import (
	"lapasta/internal/models"
)

type FornecedorService interface {
	CriarFornecedor(fornecedor *models.Fornecedor) error
	EditarFornecedor(f *models.Fornecedor) error
	ListarFornecedores(page int) ([]models.Fornecedor, error)
	BuscarFornecedorPorCNPJouNome(valor string) ([]models.Fornecedor, error)
	BuscarPedidosFornecedorPorDescricaoOuId(fornecedorId int, valor string) ([]models.PedidoFornecedor, error) 

}

type fornecedorService struct {
	repo FornecedorRepository
}

func NovoFornecedorService(repo FornecedorRepository) FornecedorService {
	return &fornecedorService{
		repo: repo,
	}
}

func (s *fornecedorService) CriarFornecedor(fornecedor *models.Fornecedor) error {
	return s.repo.CriarFornecedor(fornecedor)
}
func (s *fornecedorService) EditarFornecedor(f *models.Fornecedor) error  {
	return s.repo.EditarFornecedor(f)
}

func (s *fornecedorService) ListarFornecedores(page int) ([]models.Fornecedor, error) {
	return s.repo.ListarFornecedores(page)
}

func (s *fornecedorService) BuscarFornecedorPorCNPJouNome(valor string) ([]models.Fornecedor, error) {
	return s.repo.BuscarFornecedorPorCNPJouNome(valor)
}
func (s *fornecedorService) BuscarPedidosFornecedorPorDescricaoOuId(fornecedorId int, valor string) ([]models.PedidoFornecedor, error) {
	return s.repo.BuscarPedidosFornecedorPorDescricaoOuId(fornecedorId, valor)
}
