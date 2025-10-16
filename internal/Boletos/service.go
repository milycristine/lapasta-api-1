package boleto

import (
	"lapasta/internal/models"
)

type BoletoService interface {
	CriarBoleto(boleto *models.Boleto) error
	ListarBoletosPorFornecedor(fornecedorId int) ([]models.Boleto, error)
	ListarBoletosPorRecebimento(recebimentoId int) ([]models.Boleto, error)
	ListarBoletosDoDia(data string) ([]models.Boleto, error)
	PagarBoleto(codigoBarras string) error
	ListarBoletosVencidos() ([]models.Boleto, error)
	ListarBoletosPagos() ([]models.Boleto, error)
	ListarBoletosPendentes() ([]models.Boleto, error)
	AtualizarStatusBoleto(id int, statusId int) error
	GerarEEnviarRelatorioBoletos(emailAdmin string, ano, mes int, fornecedorID *int, statusIDs []int) error
	TotalBoletosDia(data string) (float64, error)
	TotalBoletosAtrasados() (float64, error)
	TotalBoletosPendentesMesAtual() (float64, error)
	TotalBoletosPagosMesAtual() (float64, error)
}

type boletoService struct {
	repo BoletoRepository
}

func NovoBoletoService(repo BoletoRepository) BoletoService {
	return &boletoService{
		repo: repo,
	}
}

func (s *boletoService) CriarBoleto(boleto *models.Boleto) error {
	return s.repo.CriarBoleto(boleto)
}

func (s *boletoService) ListarBoletosPorFornecedor(fornecedorId int) ([]models.Boleto, error) {
	return s.repo.ListarBoletosPorFornecedor(fornecedorId)
}

func (s *boletoService) ListarBoletosPorRecebimento(recebimentoId int) ([]models.Boleto, error) {
	return s.repo.ListarBoletosPorRecebimento(recebimentoId)
}

func (s *boletoService) ListarBoletosDoDia(data string) ([]models.Boleto, error) {
	return s.repo.ListarBoletosDoDia(data)
}

func (s *boletoService) PagarBoleto(codigoBarras string) error {
	return s.repo.PagarBoleto(codigoBarras)
}
func (s *boletoService) ListarBoletosPagos() ([]models.Boleto, error) {
	return s.repo.ListarBoletosPagos()
}

func (s *boletoService) ListarBoletosVencidos() ([]models.Boleto, error) {
	return s.repo.ListarBoletosVencidos()
}
func (s *boletoService) ListarBoletosPendentes() ([]models.Boleto, error) {
	return s.repo.ListarBoletosPendentes()
}
func (s *boletoService) AtualizarStatusBoleto(id int, statusId int) error {
	return s.repo.AtualizarStatusBoleto(id, statusId)
}
func (s *boletoService) GerarEEnviarRelatorioBoletos(emailAdmin string, ano, mes int, fornecedorID *int, statusIDs []int) error {
	return s.repo.GerarEEnviarRelatorioBoletos(emailAdmin, ano, mes, fornecedorID, statusIDs)
}
func (s *boletoService) TotalBoletosDia(data string) (float64, error) {
	return s.repo.TotalBoletosDia(data)
}
func (s *boletoService) TotalBoletosAtrasados() (float64, error) {
	return s.repo.TotalBoletosAtrasados()
}
func (s *boletoService) TotalBoletosPendentesMesAtual() (float64, error) {
	return s.repo.TotalBoletosPendentesMesAtual()
}
func (s *boletoService) TotalBoletosPagosMesAtual() (float64, error) {
	return s.repo.TotalBoletosPagosMesAtual()
}
