package recebimento

import (
	"lapasta/internal/models"
	"time"
)

type RecebimentoService interface {
	CriarRecebimento(recebimento *models.Recebimento) error
	ListarRecebimentos(page int) ([]models.Recebimento, error)
	FiltrarDataRecebimentos(inicioData, fimData time.Time) ([]models.Recebimento, error)
	ValidarRecebimento(recebimento *models.Recebimento) (bool, string, error)
	BuscarDadosRecebimentoPorNumeroNota(numeroNota string) (*models.DadosRecebimentoNota, error)
	TotalRecebimentosAvistaMesAtual() (float64, error) 
	ListarRecebimentosAvistaMesAtual() ([]models.Recebimento, error)
}

type recebimentoService struct {
	repo RecebimentoRepository
}

func NovoRecebimentoService(repo RecebimentoRepository) RecebimentoService {
	return &recebimentoService{repo: repo}
}

func (s *recebimentoService) CriarRecebimento(recebimento *models.Recebimento) error {
	return s.repo.CriarRecebimento(recebimento)
}

func (s *recebimentoService) ListarRecebimentos(page int) ([]models.Recebimento, error) {
	return s.repo.ListarRecebimentos(page)
}
func (s *recebimentoService) FiltrarDataRecebimentos(inicioData, fimData time.Time) ([]models.Recebimento, error) {
	return s.repo.FiltrarDataRecebimentos(inicioData, fimData)
}
func (s *recebimentoService) ValidarRecebimento(recebimento *models.Recebimento) (bool, string, error) {
	return s.repo.ValidarRecebimento(recebimento)
}
func (s *recebimentoService) BuscarDadosRecebimentoPorNumeroNota(numeroNota string) (*models.DadosRecebimentoNota, error) {
	return s.repo.BuscarDadosRecebimentoPorNumeroNota(numeroNota)
}
func (s *recebimentoService) TotalRecebimentosAvistaMesAtual() (float64, error) {
	return s.repo.TotalRecebimentosAvistaMesAtual()
}
func (s *recebimentoService) ListarRecebimentosAvistaMesAtual() ([]models.Recebimento, error)  {
	return s.repo.ListarRecebimentosAvistaMesAtual()
}
