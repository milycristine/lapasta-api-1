package nota

import (
	"lapasta/internal/models"
	"time"
)

type NotaService interface {
	CriarNota(nota *models.Nota) (int, error)
	ListarNotas(page int) ([]models.Nota, error)
	FiltrarDataNota(inicioData, fimData time.Time) ([]models.Nota, error)
	BuscarNotasPorNumero(numero string) ([]models.Nota, error)
}

type notaService struct {
	repo NotaRepository
}

func NovaNotaService(repo NotaRepository) NotaService {
	return &notaService{
		repo: repo,
	}
}

func (s *notaService) CriarNota(nota *models.Nota) (int, error) {
	return s.repo.CriarNota(nota)
}

func (s *notaService) ListarNotas(page int) ([]models.Nota, error) {
	return s.repo.ListarNotas(page)
}

func (s *notaService) FiltrarDataNota(inicioData, fimData time.Time) ([]models.Nota, error) {
	return s.repo.FiltrarDataNota(inicioData, fimData)
}
func (s *notaService) BuscarNotasPorNumero(numero string) ([]models.Nota, error) {
	return s.repo.BuscarNotasPorNumero(numero)
}
