package ponto

import (
	"lapasta/internal/models"
	"time"
)

type PontoService interface {
	ListarPontos(page int) ([]models.Ponto, error)
	ListarPontosPorId(idFuncionario int, page int) ([]models.Ponto, error)
	ListarPontosPorIdEDia(idFuncionario int, dia string) (models.Ponto, error)
	ListarPontosPorData(startDate, endDate time.Time, page int) ([]models.Ponto, error)
	ListarPontosPorDataId(idFuncionario int, startDate time.Time, endDate time.Time, page int) ([]models.Ponto, error)
	RegistrarEntrada(idFuncionario int) error
	RegistrarSaidaAlmoco(idFuncionario int) error
	RegistrarRetornoAlmoco(idFuncionario int) error
	RegistrarSaida(idFuncionario int) error
	GerarRelatorioMensal(mes int, ano int, emailAdmin string) (string, error)
}

type pontoService struct {
	repo PontoRepository
}

func NovoPontoService(repo PontoRepository) PontoService {
	return &pontoService{
		repo: repo,
	}
}

func (s *pontoService) ListarPontos(page int) ([]models.Ponto, error) {
	return s.repo.ListarPontos(page)
}
func (s *pontoService) ListarPontosPorId(idFuncionario int, page int) ([]models.Ponto, error) {
	return s.repo.ListarPontosPorId(idFuncionario, page)
}
func (s *pontoService) ListarPontosPorIdEDia(idFuncionario int, dia string) (models.Ponto, error) {
	return s.repo.ListarPontosPorIdEDia(idFuncionario, dia)
}

func (s *pontoService) ListarPontosPorData(startDate, endDate time.Time, page int) ([]models.Ponto, error) {
	return s.repo.ListarPontosPorData(startDate, endDate, page)
}

func (s *pontoService) ListarPontosPorDataId(IdFuncionario int, startDate time.Time, endDate time.Time, page int) ([]models.Ponto, error) {
	return s.repo.ListarPontosPorDataId(IdFuncionario, startDate, endDate, page)
}

func (s *pontoService) RegistrarEntrada(idFuncionario int) error {
	return s.repo.RegistrarEntrada(idFuncionario)
}

func (s *pontoService) RegistrarSaidaAlmoco(idFuncionario int) error {
	return s.repo.RegistrarSaidaAlmoco(idFuncionario)
}

func (s *pontoService) RegistrarRetornoAlmoco(idFuncionario int) error {
	return s.repo.RegistrarRetornoAlmoco(idFuncionario)
}

func (s *pontoService) RegistrarSaida(idFuncionario int) error {
	return s.repo.RegistrarSaida(idFuncionario)
}

func (s *pontoService) GerarRelatorioMensal(mes int, ano int, emailAdmin string) (string, error) {
	return s.repo.GerarRelatorioMensal(mes, ano, emailAdmin)
}
