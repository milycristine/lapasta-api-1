package valetransporte

import (
	"database/sql"
	"fmt"
	"time"

	database "lapasta/database"
	"lapasta/internal/models"
)

type ValeRepository interface {
	ListarVales(page int) ([]models.ValeTransportePagamento, error)
	ListarValesDaSemana(semana int, page int, mes int) ([]models.ValeTransportePagamento, error)
	AtualizarStatusVale(id int, statusIdVale int) error
}
type valeRepository struct {
	db *sql.DB
}
func NovoValeRepository(conn *database.SQLStr) ValeRepository {
	return &valeRepository{
		db: conn.DB(),
	}
}

func (r *valeRepository) ListarVales(page int) ([]models.ValeTransportePagamento, error) {
	var vales []models.ValeTransportePagamento
	hoje := time.Now()

	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT 
			p.Id, 
			p.IdFuncionario, 
			p.DataPagamento, 
			p.Valor, 
			p.StatusIdVale, 
			f.Cpf, 
			f.Nome,
			f.Sobrenome,
			f.ChavePix
		FROM 
			ValeTransportePagamentos p WITH (NOLOCK)
		JOIN 
			Funcionarios f WITH (NOLOCK) ON p.IdFuncionario = f.Id
		ORDER BY 
			CASE 
				WHEN p.StatusIdVale = 3 THEN 1
				WHEN p.StatusIdVale = 2 THEN 2
				WHEN p.StatusIdVale = 1 THEN 3
				ELSE 4
			END,
			p.DataPagamento DESC
		OFFSET @Offset ROWS
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := r.db.Query(query,
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar Vale Transporte: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.ValeTransportePagamento
		if err := rows.Scan(
			&p.Id,
			&p.IdFuncionario,
			&p.DataPagamento,
			&p.Valor,
			&p.StatusIdVale,
			&p.CpfFuncionario,
			&p.NomeFuncionario,
			&p.SobrenomeFuncionario,
			&p.ChavePixFuncionario,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
		}

		if p.StatusIdVale == 2 && p.DataPagamento.Before(hoje.Truncate(24*time.Hour)) {
			_ = r.AtualizarStatusVale(p.Id, 3)
			p.StatusIdVale = 3
		}

		vales = append(vales, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return vales, nil
}

func (r *valeRepository) AtualizarStatusVale(id int, statusIdVale int) error {
	query := `
		UPDATE ValeTransportePagamentos 
		SET StatusIdVale = @StatusIdVale 
		WHERE Id = @Id
	`
	_, err := r.db.Exec(query,
		sql.Named("Id", id),
		sql.Named("StatusIdVale", statusIdVale),
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar o vale: %w", err)
	}
	return nil
}

func (r *valeRepository) ListarValesDaSemana(semana int, page int, mes int) ([]models.ValeTransportePagamento, error) {
	var vales []models.ValeTransportePagamento
	limit := 10
	offset := (page - 1) * limit

	if semana <= 0 {
		_, week := time.Now().ISOWeek()
		semana = week
	}

	anoAtual := time.Now().Year()
	if mes <= 0 || mes > 12 {
		mes = int(time.Now().Month())
	}

	primeiroDiaMes := time.Date(anoAtual, time.Month(mes), 1, 0, 0, 0, 0, time.Local)
	offsetDias := (semana - 1) * 7
	inicioSemana := primeiroDiaMes.AddDate(0, 0, offsetDias)

	for inicioSemana.Weekday() != time.Sunday {
		inicioSemana = inicioSemana.AddDate(0, 0, -1)
	}
	fimSemana := inicioSemana.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	query := `
		SELECT 
			p.Id, 
			p.IdFuncionario, 
			p.DataPagamento, 
			p.Valor, 
			p.StatusIdVale, 
			f.Cpf, 
			f.Nome,
			f.Sobrenome,
			f.ChavePix
		FROM 
			ValeTransportePagamentos p WITH (NOLOCK)
		JOIN 
			Funcionarios f WITH (NOLOCK) ON p.IdFuncionario = f.Id
		WHERE 
			p.DataPagamento BETWEEN @InicioSemana AND @FimSemana
		ORDER BY 
			p.DataPagamento DESC
		OFFSET @Offset ROWS
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := r.db.Query(query,
		sql.Named("InicioSemana", inicioSemana),
		sql.Named("FimSemana", fimSemana),
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar vales da semana %d: %w", semana, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.ValeTransportePagamento
		if err := rows.Scan(
			&p.Id,
			&p.IdFuncionario,
			&p.DataPagamento,
			&p.Valor,
			&p.StatusIdVale,
			&p.CpfFuncionario,
			&p.NomeFuncionario,
			&p.SobrenomeFuncionario,
			&p.ChavePixFuncionario,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
		}
		vales = append(vales, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return vales, nil
}
