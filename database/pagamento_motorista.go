package sql

import (
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"time"
)

func (s *SQLStr) ContarDiasComNotasLancadas(idMotorista int, inicio, fim time.Time) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM (
			SELECT CAST(DataAtribuicao AS DATE) AS Dia
			FROM NotasMotorista
			WHERE IdMotorista = @IdMotorista
			  AND CAST(DataAtribuicao AS DATE) BETWEEN @DataInicio AND @DataFim
			GROUP BY CAST(DataAtribuicao AS DATE)
			HAVING COUNT(*) = SUM(CASE WHEN IdStatusLancamento = @StatusLancado THEN 1 ELSE 0 END)
		) AS DiasValidos
	`

	var dias int
	err := s.db.QueryRow(query,
		sql.Named("IdMotorista", idMotorista),
		sql.Named("DataInicio", inicio),
		sql.Named("DataFim", fim),
		sql.Named("StatusLancado", StatusLancado),
	).Scan(&dias)
	if err != nil {
		return 0, fmt.Errorf("erro ao contar dias com todas notas lançadas: %w", err)
	}

	return dias, nil
}

func (s *SQLStr) CriarPagamentoMotorista(p *models.PagamentosMotorista) error {
	var existe int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM PagamentosMotorista
		WHERE IdMotorista = @IdMotorista
		  AND SemanaInicio = @SemanaInicio
		  AND SemanaFim = @SemanaFim
	`,
		sql.Named("IdMotorista", p.IdMotorista),
		sql.Named("SemanaInicio", p.SemanaInicio),
		sql.Named("SemanaFim", p.SemanaFim),
	).Scan(&existe)

	if err != nil {
		return fmt.Errorf("erro ao verificar pagamento existente: %w", err)
	}
	if existe > 0 {
		return fmt.Errorf("pagamento já existe para esse motorista nesta semana")
	}

	diasValidos, err := s.ContarDiasComNotasLancadas(p.IdMotorista, p.SemanaInicio, p.SemanaFim)
	if err != nil {
		return fmt.Errorf("erro ao contar dias válidos: %w", err)
	}
	if diasValidos == 0 {
		return fmt.Errorf("nenhum dia com todas notas lançadas no período")
	}

	var valorDiaria float64
	err = s.db.QueryRow(`SELECT ValorDiaria FROM Motoristas WHERE Id = @IdMotorista`,
		sql.Named("IdMotorista", p.IdMotorista),
	).Scan(&valorDiaria)

	if err != nil {
		return fmt.Errorf("erro ao buscar valor diária: %w", err)
	}

	valorPago := float64(diasValidos) * valorDiaria

	query := `
		INSERT INTO PagamentosMotorista
		(IdMotorista, SemanaInicio, SemanaFim, DiasTrabalhados, ValorPago, IdStatusPagamento)
		VALUES (@IdMotorista, @SemanaInicio, @SemanaFim, @DiasTrabalhados, @ValorPago, @IdStatusPagamento)
	`
	_, err = s.db.Exec(query,
		sql.Named("IdMotorista", p.IdMotorista),
		sql.Named("SemanaInicio", p.SemanaInicio),
		sql.Named("SemanaFim", p.SemanaFim),
		sql.Named("DiasTrabalhados", diasValidos),
		sql.Named("ValorPago", valorPago),
		sql.Named("IdStatusPagamento", p.IdStatusPagamento),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir pagamento motorista: %w", err)
	}

	return nil
}

func (s *SQLStr) CalcularPagamentoMotorista(idMotorista int, inicio, fim time.Time) (*models.PagamentosMotorista, error) {
	diasValidos, err := s.ContarDiasComNotasLancadas(idMotorista, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao contar dias válidos: %w", err)
	}
	if diasValidos == 0 {
		return nil, fmt.Errorf("não há notas lançadas no período")
	}

	var nome, chavePix string
	var valorDiaria float64
	err = s.db.QueryRow(`
		SELECT Nome, ChavePix, ValorDiaria 
		FROM Motoristas 
		WHERE Id = @IdMotorista`,
		sql.Named("IdMotorista", idMotorista),
	).Scan(&nome, &chavePix, &valorDiaria)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados do motorista: %w", err)
	}

	valorPago := float64(diasValidos) * valorDiaria

	return &models.PagamentosMotorista{
		IdMotorista:     idMotorista,
		SemanaInicio:    inicio,
		SemanaFim:       fim,
		DiasTrabalhados: diasValidos,
		ValorPago:       valorPago,
		NomeMotorista:   nome,
		ChavePix:        chavePix,
	}, nil
}

func (s *SQLStr) ListarPagamentosMotorista(idMotorista int) ([]models.PagamentosMotorista, error) {
	query := `
		SELECT 
			p.Id,
			p.IdMotorista,
			p.SemanaInicio,
			p.SemanaFim,
			p.DiasTrabalhados,
			p.ValorPago,
			ISNULL(p.DataPagamento, '1900-01-01'),
			p.IdStatusPagamento,
			m.Nome,
			m.ChavePix
		FROM PagamentosMotorista p WITH (NOLOCK) 
		INNER JOIN Motoristas m WITH (NOLOCK) ON p.IdMotorista = m.Id
	`
	var args []any

	if idMotorista > 0 {
		query += " WHERE p.IdMotorista = @IdMotorista"
		args = append(args, sql.Named("IdMotorista", idMotorista))
	}

	query += " ORDER BY p.SemanaInicio DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pagamentos motorista: %w", err)
	}
	defer rows.Close()

	var resultado []models.PagamentosMotorista
	for rows.Next() {
		var p models.PagamentosMotorista

		if err := rows.Scan(
			&p.Id,
			&p.IdMotorista,
			&p.SemanaInicio,
			&p.SemanaFim,
			&p.DiasTrabalhados,
			&p.ValorPago,
			&p.DataPagamento,
			&p.IdStatusPagamento,
			&p.NomeMotorista,
			&p.ChavePix,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear pagamento motorista: %w", err)
		}

		resultado = append(resultado, p)
	}

	return resultado, nil
}

func (s *SQLStr) AtualizarStatusPagamentoMotorista(id int, statusId int) error {
	query := `
		UPDATE PagamentosMotorista
		SET IdStatusPagamento = @StatusId
		WHERE Id = @Id
	`
	_, err := s.db.Exec(query,
		sql.Named("StatusId", statusId),
		sql.Named("Id", id),
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status pagamento: %w", err)
	}
	return nil
}
