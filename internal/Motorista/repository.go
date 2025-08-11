package motorista

import (
	database "lapasta/database"
	"lapasta/internal/models"
	"time"
)

type MotoristaRepository interface {
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

type motoristaRepository struct {
	db *database.SQLStr
}

func NovoMotoristaRepository(db *database.SQLStr) MotoristaRepository {
	return &motoristaRepository{
		db: db,
	}
}

func (r *motoristaRepository) CriarMotorista(m *models.Motorista) error {
	return r.db.CriarMotorista(m)
}
func (r *motoristaRepository) EditarMotorista(m *models.Motorista) error {
	return r.db.EditarMotorista(m)
}
func (r *motoristaRepository) AtualizarStatusMotorista(id int, ativo bool) error {
	return r.db.AtualizarStatusMotorista(id, ativo)
}

func (r *motoristaRepository) ListarMotoristas() ([]models.Motorista, error) {
	return r.db.ListarMotoristas()
}
func (r *motoristaRepository) BuscarMotoristaPorCPFouNome(valor string) ([]models.Motorista, error) {
	return r.db.BuscarMotoristaPorCPFouNome(valor)
}
func (r *motoristaRepository) BuscarMotoristaPorID(id int) (*models.Motorista, error)  {
	return r.db.BuscarMotoristaPorID(id)
}

func (r *motoristaRepository) CriarEmissaoNota(n *models.EmissaoNota) error {
	return r.db.CriarEmissaoNota(n)
}

func (r *motoristaRepository) ListarEmissaoNotas(page int) ([]models.EmissaoNota, error) {
	return r.db.ListarEmissaoNotas(page)
}
func (r *motoristaRepository) ListarEmissaoNotasPorMotorista(idMotorista int) ([]models.EmissaoNota, error) {
	return r.db.ListarEmissaoNotasPorMotorista(idMotorista)
}

func (r *motoristaRepository) BuscarEmissaoNotas(valor string) ([]models.EmissaoNota, error) {
	return r.db.BuscarEmissaoNotas(valor)
}
func (r *motoristaRepository) FiltrarDataEmissaoNota(inicioData time.Time, fimData time.Time) ([]models.EmissaoNota, error) {
	return r.db.FiltrarDataEmissaoNota(inicioData, fimData)
}

func (r *motoristaRepository) CriarNotasMotorista(nm *models.NotasMotorista) error {
	return r.db.CriarNotasMotorista(nm)
}

func (r *motoristaRepository) ListarNotasMotoristaPorMotorista(idMotorista int) ([]models.NotasMotorista, error) {
	return r.db.ListarNotasMotoristaPorMotorista(idMotorista)
}

func (r *motoristaRepository) AtualizarStatusLancamentoNotaMotorista(id int, statusId int, dataLancamento *time.Time) error {
	return r.db.AtualizarStatusLancamentoNotaMotorista(id, statusId, dataLancamento)
}

func (r *motoristaRepository) MotoristaLancouTodasAsNotas(idMotorista int, inicio, fim time.Time) (bool, error) {
	return r.db.MotoristaLancouTodasAsNotas(idMotorista, inicio, fim)
}

func (r *motoristaRepository) CriarPagamentoMotorista(p *models.PagamentosMotorista) error {
	return r.db.CriarPagamentoMotorista(p)
}

func (r *motoristaRepository) ListarPagamentosMotorista(idMotorista int) ([]models.PagamentosMotorista, error) {
	return r.db.ListarPagamentosMotorista(idMotorista)
}

func (r *motoristaRepository) AtualizarStatusPagamentoMotorista(id int, statusId int) error {
	return r.db.AtualizarStatusPagamentoMotorista(id, statusId)
}
func (r *motoristaRepository) FiltrarNotasMotoristaPorData(
	idMotorista int,
	inicioData *time.Time,
	fimData *time.Time,
) ([]models.NotasMotorista, error) {
	return r.db.FiltrarNotasMotoristaPorData(idMotorista, inicioData, fimData)
}
func (r *motoristaRepository) CalcularPagamentoMotorista(idMotorista int, inicio, fim time.Time) (*models.PagamentosMotorista, error) {
	return r.db.CalcularPagamentoMotorista(idMotorista, inicio, fim)
}
