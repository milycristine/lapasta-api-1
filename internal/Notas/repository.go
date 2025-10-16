package nota

import (
	database "lapasta/database"
	"lapasta/internal/models"
	"time"
)

type NotaRepository interface {
	CriarNota(nota *models.Nota) (int, error)
	ListarNotas(page int) ([]models.Nota, error)
	FiltrarDataNota(inicioData, fimData time.Time) ([]models.Nota, error)
	BuscarNotasPorNumero(numero string) ([]models.Nota, error)
}

type notaRepository struct {
	db *database.SQLStr
}

func NovoNotaRepository(db *database.SQLStr) NotaRepository {
	return &notaRepository{
		db: db,
	}
}

func (r *notaRepository) CriarNota(nota *models.Nota) (int, error) {
	return r.db.CriarNota(nota)
}

func (r *notaRepository) ListarNotas(page int) ([]models.Nota, error) {
	return r.db.ListarNotas(page)
}

func (r *notaRepository) FiltrarDataNota(inicioData, fimData time.Time) ([]models.Nota, error) {
	return r.db.FiltrarDataNota(inicioData, fimData)
}
func (r *notaRepository) BuscarNotasPorNumero(numero string) ([]models.Nota, error) {
	return r.db.BuscarNotasPorNumero(numero)
}
