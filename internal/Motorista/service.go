package motorista

import (
	"lapasta/internal/models"
	"time"
)

type MotoristaService interface {
	CriarMotorista(m *models.Motorista) error
	EditarMotorista(m *models.Motorista) error
	AtualizarStatusMotorista(id int, ativo bool) error
	ListarMotoristas() ([]models.Motorista, error)
	BuscarMotoristaPorCPFouNome(valor string) ([]models.Motorista, error)
	BuscarMotoristaPorID(id int) (*models.Motorista, error)

	CriarEmissaoNota(n *models.EmissaoNota) error
	ListarEmissaoNotas(page int) ([]models.EmissaoNota, error)
	ListarEmissaoNotasPorMotorista(idMotorista int) ([]models.EmissaoNota, error)
	BuscarEmissaoNotas(valor string) ([]models.EmissaoNota, error)
	FiltrarDataEmissaoNota(inicioData time.Time, fimData time.Time) ([]models.EmissaoNota, error)

	CriarNotasMotorista(nm *models.NotasMotorista) error
	ListarNotasMotoristaPorMotorista(idMotorista int) ([]models.NotasMotorista, error)
	AtualizarStatusLancamentoNotaMotorista(id int, statusId int, dataLancamento *time.Time) error
	MotoristaLancouTodasAsNotas(idMotorista int, inicio, fim time.Time) (bool, error)
	FiltrarNotasMotoristaPorData(
		idMotorista int,
		inicioData *time.Time,
		fimData *time.Time,
	) ([]models.NotasMotorista, error)

	CriarPagamentoMotorista(p *models.PagamentosMotorista) error
	ListarPagamentosMotorista(idMotorista int) ([]models.PagamentosMotorista, error)
	AtualizarStatusPagamentoMotorista(id int, statusId int) error
	CalcularPagamentoMotorista(idMotorista int, inicio, fim time.Time) (*models.PagamentosMotorista, error)
}

type motoristaService struct {
	repo MotoristaRepository
}

func NovoMotoristaService(repo MotoristaRepository) MotoristaService {
	return &motoristaService{
		repo: repo,
	}
}

func (s *motoristaService) CriarMotorista(m *models.Motorista) error {
	return s.repo.CriarMotorista(m)
}
func (s *motoristaService) EditarMotorista(m *models.Motorista) error {
	return s.repo.EditarMotorista(m)
}
func (s *motoristaService) AtualizarStatusMotorista(id int, ativo bool) error {
	return s.repo.AtualizarStatusMotorista(id, ativo)
}

func (s *motoristaService) ListarMotoristas() ([]models.Motorista, error) {
	return s.repo.ListarMotoristas()
}
func (s *motoristaService) BuscarMotoristaPorCPFouNome(valor string) ([]models.Motorista, error) {
	return s.repo.BuscarMotoristaPorCPFouNome(valor)
}
func (s *motoristaService) BuscarMotoristaPorID(id int) (*models.Motorista, error) {
	return s.repo.BuscarMotoristaPorID(id)
}

func (s *motoristaService) CriarEmissaoNota(n *models.EmissaoNota) error {
	return s.repo.CriarEmissaoNota(n)
}

func (s *motoristaService) ListarEmissaoNotas(page int) ([]models.EmissaoNota, error) {
	return s.repo.ListarEmissaoNotas(page)
}
func (s *motoristaService) ListarEmissaoNotasPorMotorista(idMotorista int) ([]models.EmissaoNota, error) {
	return s.repo.ListarEmissaoNotasPorMotorista(idMotorista)
}
func (s *motoristaService) BuscarEmissaoNotas(valor string) ([]models.EmissaoNota, error) {
	return s.repo.BuscarEmissaoNotas(valor)
}
func (s *motoristaService) FiltrarDataEmissaoNota(inicioData time.Time, fimData time.Time) ([]models.EmissaoNota, error) {
	return s.repo.FiltrarDataEmissaoNota(inicioData, fimData)
}

func (s *motoristaService) CriarNotasMotorista(nm *models.NotasMotorista) error {
	return s.repo.CriarNotasMotorista(nm)
}

func (s *motoristaService) ListarNotasMotoristaPorMotorista(idMotorista int) ([]models.NotasMotorista, error) {
	return s.repo.ListarNotasMotoristaPorMotorista(idMotorista)
}

func (s *motoristaService) AtualizarStatusLancamentoNotaMotorista(id int, statusId int, dataLancamento *time.Time) error {
	return s.repo.AtualizarStatusLancamentoNotaMotorista(id, statusId, dataLancamento)
}

func (s *motoristaService) MotoristaLancouTodasAsNotas(idMotorista int, inicio, fim time.Time) (bool, error) {
	return s.repo.MotoristaLancouTodasAsNotas(idMotorista, inicio, fim)
}

func (s *motoristaService) CriarPagamentoMotorista(p *models.PagamentosMotorista) error {
	return s.repo.CriarPagamentoMotorista(p)
}

func (s *motoristaService) ListarPagamentosMotorista(idMotorista int) ([]models.PagamentosMotorista, error) {
	return s.repo.ListarPagamentosMotorista(idMotorista)
}

func (s *motoristaService) AtualizarStatusPagamentoMotorista(id int, statusId int) error {
	return s.repo.AtualizarStatusPagamentoMotorista(id, statusId)
}
func (s *motoristaService) FiltrarNotasMotoristaPorData(
	idMotorista int,
	inicioData *time.Time,
	fimData *time.Time,
) ([]models.NotasMotorista, error) {
	return s.repo.FiltrarNotasMotoristaPorData(idMotorista, inicioData, fimData)
}
func (s *motoristaService) CalcularPagamentoMotorista(idMotorista int, inicio, fim time.Time) (*models.PagamentosMotorista, error) {
	return s.repo.CalcularPagamentoMotorista(idMotorista, inicio, fim)
}
