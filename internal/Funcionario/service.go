package funcionario

import (
	"lapasta/internal/models"
)

type FuncionarioService interface {
	CriarFuncionario(funcionario *models.Funcionario) error
	EditarFuncionario(funcionario *models.Funcionario) error
	ListarFuncionarios(page int) ([]models.Funcionario, error)
	BuscarFuncionarioPorCPF(cpf string) (*models.FuncionarioComPontos, error)
	BuscarFuncionarioPorID(id int) (*models.Funcionario, error) 


}

type funcionarioService struct {
	repo FuncionarioRepository
}

func NovoFuncionarioService(repo FuncionarioRepository) FuncionarioService {
	return &funcionarioService{
		repo: repo,
	}
}

func (s *funcionarioService) CriarFuncionario(funcionario *models.Funcionario) error {
	return s.repo.CriarFuncionario(funcionario)
}

func (s *funcionarioService) EditarFuncionario( funcionario *models.Funcionario)	error  {
	return s.repo.EditarFuncionario(funcionario)
}

func (s *funcionarioService) ListarFuncionarios(page int) ([]models.Funcionario, error) {
	return s.repo.ListarFuncionarios(page)
}
func (s *funcionarioService) BuscarFuncionarioPorID(id int) (*models.Funcionario, error)  {
	return s.repo.BuscarFuncionarioPorID(id)
}

func (s *funcionarioService) BuscarFuncionarioPorCPF(cpf string) (*models.FuncionarioComPontos, error) {
	funcionarioComPontos, err := s.repo.BuscarFuncionarioPorCPF(cpf)
	if err != nil {
		return nil, err
	}
	return funcionarioComPontos, nil
}