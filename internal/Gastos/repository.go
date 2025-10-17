package gastos

import (
	database "lapasta/database"
	"lapasta/internal/models"
)

type GastosRepository interface {
	TotalGastos(tipo string) (float64, error)
	ListarGastos(tipo, inicioData, fimData string) ([]models.Gasto, error)
}

type gastosRepository struct {
	db *database.SQLStr
}

func NovoGastosRepository(db *database.SQLStr) GastosRepository {
	return &gastosRepository{db: db}
}

func (r *gastosRepository) TotalGastos(tipo string) (float64, error) {
	return r.db.TotalGastos(tipo)
}

func (r *gastosRepository) ListarGastos(tipo, inicioData, fimData string) ([]models.Gasto, error) {
	return r.db.ListarGastos(tipo, inicioData, fimData)
}
