package gastos

import (
	"lapasta/internal/models"
)

type GastosService interface {
	TotalGastos(tipo string) (float64, error)
	ListarGastos(tipo, inicioData, fimData string) ([]models.Gasto, error)
}

type gastosService struct {
	repo GastosRepository
}

func NovoGastosService(repo GastosRepository) GastosService {
	return &gastosService{repo: repo}
}

func (s *gastosService) TotalGastos(tipo string) (float64, error) {
	return s.repo.TotalGastos(tipo)
}

func (s *gastosService) ListarGastos(tipo, inicioData, fimData string) ([]models.Gasto, error) {
	return s.repo.ListarGastos(tipo, inicioData, fimData)
}
