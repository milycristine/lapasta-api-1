package funcionario

import (
	database "lapasta/database"
	"lapasta/internal/models"
)

type FuncionarioRepository interface {
	CriarFuncionario(funcionario *models.Funcionario) error
	EditarFuncionario(funcionario *models.Funcionario) error
	ListarFuncionarios(page int) ([]models.Funcionario, error)
	BuscarFuncionarioPorCPF(cpf string) (*models.FuncionarioComPontos, error)
	BuscarFuncionarioPorID(id int) (*models.Funcionario, error)
}

type funcionarioRepository struct {
	db *database.SQLStr
}

func NovoFuncionarioRepository(db *database.SQLStr) FuncionarioRepository {
	return &funcionarioRepository{
		db: db,
	}
}

func (r *funcionarioRepository) CriarFuncionario(funcionario *models.Funcionario) error {
	return r.db.CriarFuncionario(funcionario)
}

func (r *funcionarioRepository) EditarFuncionario(funcionario *models.Funcionario) error {
	return r.db.EditarFuncionario(funcionario)
}

func (r *funcionarioRepository) ListarFuncionarios(page int) ([]models.Funcionario, error) {
	return r.db.ListarFuncionarios(page)
}
func (r *funcionarioRepository) BuscarFuncionarioPorID(id int) (*models.Funcionario, error) {
	return r.db.BuscarFuncionarioPorID(id)
}

func (r *funcionarioRepository) BuscarFuncionarioPorCPF(cpf string) (*models.FuncionarioComPontos, error) {
	funcionarioComPontos, err := r.db.BuscarFuncionarioPorCPF(cpf)
	if err != nil {
		return nil, err
	}
	return &funcionarioComPontos, nil
}
