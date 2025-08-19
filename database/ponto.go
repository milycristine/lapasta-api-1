package sql

import (
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"time"
)

func (s *SQLStr) RegistrarEntrada(idFuncionario int) error {
	horaAtual := time.Now().Format("15:04")
	diaAtual := time.Now().Format("2006-01-02")

	var hManha sql.NullString
	query := "SELECT HManha FROM Ponto WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia"
	err := s.db.QueryRow(query, sql.Named("IdFuncionario", idFuncionario), sql.Named("Dia", diaAtual)).Scan(&hManha)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("erro ao verificar ponto de entrada existente: %w", err)
	}

	if hManha.Valid {
		_, err = s.db.Exec("UPDATE Ponto SET HManha = @HManha, Situacao = 'Presente', EntradaRegistrada = 1 WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia",
			sql.Named("HManha", horaAtual),
			sql.Named("IdFuncionario", idFuncionario),
			sql.Named("Dia", diaAtual),
		)
		if err != nil {
			return fmt.Errorf("erro ao atualizar ponto de entrada: %w", err)
		}
	} else {
		_, err = s.db.Exec("INSERT INTO Ponto (IdFuncionario, Dia, HManha, Situacao, EntradaRegistrada) VALUES (@IdFuncionario, @Dia, @HManha, 'Presente', 1)",
			sql.Named("IdFuncionario", idFuncionario),
			sql.Named("Dia", diaAtual),
			sql.Named("HManha", horaAtual),
		)
		if err != nil {
			return fmt.Errorf("erro ao inserir ponto de entrada: %w", err)
		}
	}

	return nil
}

func (s *SQLStr) RegistrarSaidaAlmoco(idFuncionario int) error {
	horaAtual := time.Now().Format("15:04")
	diaAtual := time.Now().Format("2006-01-02")

	var hManha, hAlmocoSaida sql.NullString
	query := "SELECT HManha, HAlmocoSaida FROM Ponto WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia"
	err := s.db.QueryRow(query, sql.Named("IdFuncionario", idFuncionario), sql.Named("Dia", diaAtual)).Scan(&hManha, &hAlmocoSaida)
	if err != nil {
		return fmt.Errorf("erro ao verificar ponto de saída para almoço: %w", err)
	}
	if !hManha.Valid {
		return fmt.Errorf("ponto de entrada não registrado para hoje")
	}
	if hAlmocoSaida.Valid {
		return fmt.Errorf("saída para almoço já registrada para hoje")
	}

	_, err = s.db.Exec("UPDATE Ponto SET HAlmocoSaida = @HAlmocoSaida, PausaRegistrada = 1 WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia",
		sql.Named("HAlmocoSaida", horaAtual),
		sql.Named("IdFuncionario", idFuncionario),
		sql.Named("Dia", diaAtual),
	)
	if err != nil {
		return fmt.Errorf("erro ao registrar saída para o almoço: %w", err)
	}
	return nil
}

func (s *SQLStr) RegistrarRetornoAlmoco(idFuncionario int) error {
	horaAtual := time.Now().Format("15:04")
	diaAtual := time.Now().Format("2006-01-02")

	var hAlmocoSaida, HAlmocoRetorno sql.NullString
	query := "SELECT HAlmocoSaida, HAlmocoRetorno FROM Ponto WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia"
	err := s.db.QueryRow(query, sql.Named("IdFuncionario", idFuncionario), sql.Named("Dia", diaAtual)).Scan(&hAlmocoSaida, &HAlmocoRetorno)
	if err != nil {
		return fmt.Errorf("erro ao verificar retorno do almoço: %w", err)
	}
	if !hAlmocoSaida.Valid {
		return fmt.Errorf("saída para almoço não registrada para hoje")
	}
	if HAlmocoRetorno.Valid {
		return fmt.Errorf("retorno do almoço já registrado para hoje")
	}

	_, err = s.db.Exec("UPDATE Ponto SET HAlmocoRetorno = @HAlmocoRetorno, RetornoRegistrado = 1 WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia",
		sql.Named("HAlmocoRetorno", horaAtual),
		sql.Named("IdFuncionario", idFuncionario),
		sql.Named("Dia", diaAtual),
	)
	if err != nil {
		return fmt.Errorf("erro ao registrar retorno do almoço: %w", err)
	}
	return nil
}

func (s *SQLStr) RegistrarSaida(idFuncionario int) error {
	horaAtual := time.Now().Format("15:04")
	diaAtual := time.Now().Format("2006-01-02")

	var hAlmocoRetorno, hNoite sql.NullString
	query := "SELECT HAlmocoRetorno, HNoite FROM Ponto WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia"
	err := s.db.QueryRow(query, sql.Named("IdFuncionario", idFuncionario), sql.Named("Dia", diaAtual)).Scan(&hAlmocoRetorno, &hNoite)
	if err != nil {
		return fmt.Errorf("erro ao verificar saída final: %w", err)
	}
	if !hAlmocoRetorno.Valid {
		return fmt.Errorf("retorno do almoço não registrado para hoje")
	}
	if hNoite.Valid {
		return fmt.Errorf("saída já registrada para hoje")
	}

	_, err = s.db.Exec("UPDATE Ponto SET HNoite = @HNoite, SaidaRegistrada = 1 WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia",
		sql.Named("HNoite", horaAtual),
		sql.Named("IdFuncionario", idFuncionario),
		sql.Named("Dia", diaAtual),
	)
	if err != nil {
		return fmt.Errorf("erro ao registrar saída final: %w", err)
	}
	return nil
}

func (s *SQLStr) MarcarAusente(idFuncionario int) error {
	diaAtual := time.Now().Format("2006-01-02")

	_, err := s.db.Exec(`
		UPDATE Ponto
		SET Situacao = 'Ausente'
		WHERE IdFuncionario = @IdFuncionario AND Dia = @Dia AND Situacao IS NULL
	`,
		sql.Named("IdFuncionario", idFuncionario),
		sql.Named("Dia", diaAtual),
	)
	if err != nil {
		return fmt.Errorf("erro ao marcar como ausente: %w", err)
	}
	return nil
}
func (s *SQLStr) ListarPontos(page int) ([]models.Ponto, error) {
	var pontos []models.Ponto

	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT 
			p.Id, p.HManha, p.HAlmocoRetorno, p.HAlmocoSaida, p.HNoite, 
			p.Dia, p.Situacao, p.IdFuncionario, f.Nome AS NomeFuncionario
		FROM Ponto p WITH (NOLOCK)
		JOIN Funcionarios f WITH (NOLOCK) ON p.IdFuncionario = f.Id
		ORDER BY p.Dia DESC
		OFFSET @Offset ROWS
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar pontos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Ponto
		if err := rows.Scan(
			&p.Id, &p.HManha, &p.HAlmocoRetorno, &p.HAlmocoSaida, &p.HNoite,
			&p.Dia, &p.Situacao, &p.IdFuncionario, &p.NomeFuncionario,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear ponto: %w", err)
		}
		pontos = append(pontos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return pontos, nil
}
func (s *SQLStr) ListarPontosPorId(idFuncionario, page int) ([]models.Ponto, error) {
	var pontos []models.Ponto

	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT 
			p.Id, p.HManha, p.HAlmocoRetorno, p.HAlmocoSaida, p.HNoite, 
			p.Dia, p.Situacao, p.IdFuncionario, f.Nome AS NomeFuncionario
		FROM Ponto p WITH (NOLOCK)
		JOIN Funcionarios f WITH (NOLOCK) ON p.IdFuncionario = f.Id
		WHERE p.IdFuncionario = @IdFuncionario
		ORDER BY p.Dia DESC
		OFFSET @Offset ROWS
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(
		query,
		sql.Named("IdFuncionario", idFuncionario),
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar pontos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Ponto
		if err := rows.Scan(
			&p.Id, &p.HManha, &p.HAlmocoRetorno, &p.HAlmocoSaida, &p.HNoite,
			&p.Dia, &p.Situacao, &p.IdFuncionario, &p.NomeFuncionario,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear ponto: %w", err)
		}
		pontos = append(pontos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return pontos, nil
}
func (s *SQLStr) ListarPontosPorIdEDia(idFuncionario int, dia string) (models.Ponto, error) {
	var ponto models.Ponto

	query := `
		SELECT p.Id, p.HManha, p.HAlmocoRetorno, p.HAlmocoSaida, p.HNoite, p.Dia, p.Situacao, p.IdFuncionario, 
		       f.Nome AS NomeFuncionario, p.EntradaRegistrada, p.PausaRegistrada, p.RetornoRegistrado, p.SaidaRegistrada
		FROM Ponto p WITH (NOLOCK)
		JOIN Funcionarios f WITH (NOLOCK) ON p.IdFuncionario = f.Id
		WHERE p.IdFuncionario = @IdFuncionario
		AND p.Dia = @Dia
		ORDER BY p.Dia DESC
	`

	rows, err := s.db.Query(query, sql.Named("IdFuncionario", idFuncionario), sql.Named("Dia", dia))
	if err != nil {
		return ponto, fmt.Errorf("erro ao consultar pontos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Ponto
		if err := rows.Scan(&p.Id, &p.HManha, &p.HAlmocoRetorno, &p.HAlmocoSaida, &p.HNoite, &p.Dia, &p.Situacao, &p.IdFuncionario,
			&p.NomeFuncionario, &p.EntradaRegistrada, &p.PausaRegistrada, &p.RetornoRegistrado, &p.SaidaRegistrada); err != nil {
			return ponto, fmt.Errorf("erro ao escanear ponto: %w", err)
		}
		ponto = p
	}

	if err := rows.Err(); err != nil {
		return ponto, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return ponto, nil
}
func (s *SQLStr) ListarPontosPorData(startDate time.Time, endDate time.Time, page int) ([]models.Ponto, error) {
	var pontos []models.Ponto

	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT p.Id, p.HManha, p.HAlmocoRetorno, p.HAlmocoSaida, p.HNoite, p.Dia, p.Situacao, p.IdFuncionario, 
		       f.Nome AS NomeFuncionario
		FROM Ponto p WITH (NOLOCK)
		JOIN Funcionarios f WITH (NOLOCK) ON p.IdFuncionario = f.Id
		WHERE p.Dia BETWEEN @StartDate AND @EndDate
		ORDER BY p.Dia DESC
		OFFSET @Offset ROWS
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(query,
		sql.Named("StartDate", startDate),
		sql.Named("EndDate", endDate),
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar pontos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Ponto
		if err := rows.Scan(&p.Id, &p.HManha, &p.HAlmocoRetorno, &p.HAlmocoSaida, &p.HNoite, &p.Dia, &p.Situacao,
			&p.IdFuncionario, &p.NomeFuncionario); err != nil {
			return nil, fmt.Errorf("erro ao escanear ponto: %w", err)
		}
		pontos = append(pontos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return pontos, nil
}

func (s *SQLStr) ListarPontosPorDataId(idFuncionario int, startDate time.Time, endDate time.Time, page int) ([]models.Ponto, error) {
	var pontos []models.Ponto

	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT p.Id, p.HManha, p.HAlmocoRetorno, p.HAlmocoSaida, p.HNoite, p.Dia, p.Situacao, p.IdFuncionario, 
		       f.Nome AS NomeFuncionario
		FROM Ponto p WITH (NOLOCK)
		JOIN Funcionarios f WITH (NOLOCK) ON p.IdFuncionario = f.Id
		WHERE p.IdFuncionario = @IdFuncionario
		AND p.Dia BETWEEN @StartDate AND @EndDate
		ORDER BY p.Dia DESC
		OFFSET @Offset ROWS
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(query,
		sql.Named("IdFuncionario", idFuncionario),
		sql.Named("StartDate", startDate),
		sql.Named("EndDate", endDate),
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar pontos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Ponto
		if err := rows.Scan(&p.Id, &p.HManha, &p.HAlmocoRetorno, &p.HAlmocoSaida, &p.HNoite, &p.Dia, &p.Situacao,
			&p.IdFuncionario, &p.NomeFuncionario); err != nil {
			return nil, fmt.Errorf("erro ao escanear ponto: %w", err)
		}
		pontos = append(pontos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return pontos, nil
}
