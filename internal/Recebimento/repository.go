package recebimento

import (
	database "lapasta/database"
	"lapasta/internal/models"
	"time"
)

type RecebimentoRepository interface {
	CriarRecebimento(recebimento *models.Recebimento) error
	ListarRecebimentos(page int) ([]models.Recebimento, error)
	FiltrarDataRecebimentos(inicioData, fimData time.Time) ([]models.Recebimento, error)
	ValidarRecebimento(recebimento *models.Recebimento) (bool, string, error)
	BuscarDadosRecebimentoPorNumeroNota(numeroNota string) (*models.DadosRecebimentoNota, error)
	TotalRecebimentosAvistaMesAtual() (float64, error)
	ListarRecebimentosAvistaMesAtual() ([]models.Recebimento, error)
}

type recebimentoRepository struct {
	db *database.SQLStr
}

func NovoRecebimentoRepository(db *database.SQLStr) RecebimentoRepository {
	return &recebimentoRepository{
		db: db,
	}
}

func (r *recebimentoRepository) CriarRecebimento(recebimento *models.Recebimento) error {
	return r.db.CriarRecebimento(recebimento)
}

func (r *recebimentoRepository) ListarRecebimentos(page int) ([]models.Recebimento, error) {
	return r.db.ListarRecebimentos(page)
}
func (r *recebimentoRepository) FiltrarDataRecebimentos(inicioData, fimData time.Time) ([]models.Recebimento, error) {
	return r.db.FiltrarDataRecebimentos(inicioData, fimData)
}
func (r *recebimentoRepository) ValidarRecebimento(recebimento *models.Recebimento) (bool, string, error) {
	return r.db.ValidarRecebimento(recebimento)
}
func (r *recebimentoRepository) BuscarDadosRecebimentoPorNumeroNota(numeroNota string) (*models.DadosRecebimentoNota, error) {
	return r.db.BuscarDadosRecebimentoPorNumeroNota(numeroNota)
}
func (r *recebimentoRepository) TotalRecebimentosAvistaMesAtual() (float64, error) {
	return r.db.TotalRecebimentosAvistaMesAtual()
}
func (r *recebimentoRepository) ListarRecebimentosAvistaMesAtual() ([]models.Recebimento, error) {
	return r.db.ListarRecebimentosAvistaMesAtual()
}
