package sql

import (
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"time"
)

func (s *SQLStr) CriarRecebimento(recebimento *models.Recebimento) error {
	var existe bool
	checkQuery := `SELECT 1 FROM Recebimento WHERE IdNota = @IdNota`
	err := s.db.QueryRow(checkQuery, sql.Named("IdNota", recebimento.IdNota)).Scan(&existe)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("erro ao verificar recebimentos existentes: %w", err)
	}
	if existe {
		return fmt.Errorf("já existe um recebimento para a esse número de nota com ID %d", recebimento.IdNota)
	}

	query := `
        INSERT INTO Recebimento (
            Dia, UrlImagem, IdResponsavel, Quantidade, Peso, Valor, 
            Vencimento, IdNota, IdPedidoFornecedor
        ) 
        VALUES (
            @Dia, @UrlImagem, @IdResponsavel, @Quantidade, @Peso, @Valor, 
            @Vencimento, @IdNota, @IdPedidoFornecedor
        );
        SELECT SCOPE_IDENTITY();
    `
	var id int
	err = s.db.QueryRow(query,
		sql.Named("Dia", recebimento.Dia),
		sql.Named("UrlImagem", recebimento.UrlImagem),
		sql.Named("IdResponsavel", recebimento.IdResponsavel),
		sql.Named("Quantidade", recebimento.Quantidade),
		sql.Named("Peso", recebimento.Peso),
		sql.Named("Valor", recebimento.Valor),
		sql.Named("Vencimento", recebimento.Vencimento),
		sql.Named("IdNota", recebimento.IdNota),
		sql.Named("IdPedidoFornecedor", recebimento.IdPedidoFornecedor),
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("erro ao inserir recebimento: %w", err)
	}

	recebimento.Id = id
	return nil
}

func (s *SQLStr) ListarRecebimentos(page int) ([]models.Recebimento, error) {
	var recebimentos []models.Recebimento
	limit := 10
	offset := (page - 1) * limit

	query := `
	SELECT 
		R.Id, R.Dia, R.UrlImagem, R.IdResponsavel, 
		F.Nome AS NomeResponsavel, R.Quantidade, R.Peso, R.Valor, 
		R.Vencimento, R.IdNota, R.IdPedidoFornecedor,
		N.NumeroNota,
		PF.Descricao AS Produto,
		FORNE.Nome AS NomeFornecedor,
		PF.PrazoAcordadoDias
	FROM Recebimento R WITH (NOLOCK)
	JOIN Funcionarios F WITH (NOLOCK) ON F.Id = R.IdResponsavel
	JOIN Notas N WITH (NOLOCK) ON N.Id = R.IdNota
	JOIN PedidoFornecedor PF WITH (NOLOCK) ON PF.Id = R.IdPedidoFornecedor
	JOIN Fornecedores FORNE WITH (NOLOCK) ON FORNE.Id = PF.FornecedorId
	ORDER BY R.Dia DESC
	OFFSET @Offset ROWS
	FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar recebimentos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Recebimento
		if err := rows.Scan(
			&r.Id, &r.Dia, &r.UrlImagem, &r.IdResponsavel,
			&r.NomeResponsavel, &r.Quantidade, &r.Peso, &r.Valor,
			&r.Vencimento, &r.IdNota, &r.IdPedidoFornecedor,
			&r.NumeroNota, &r.Produto,
			&r.NomeFornecedor, &r.PrazoAcordadoDias,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear recebimento: %w", err)
		}
		recebimentos = append(recebimentos, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return recebimentos, nil
}
func (s *SQLStr) FiltrarDataRecebimentos(inicioData time.Time, fimData time.Time) ([]models.Recebimento, error) {
	var recebimentos []models.Recebimento

	query := `
	SELECT 
		R.Id, R.Dia, R.UrlImagem, R.IdResponsavel, 
		F.Nome AS NomeResponsavel, R.Quantidade, R.Peso, R.Valor, 
		R.Vencimento, R.IdNota, R.IdPedidoFornecedor,
		N.NumeroNota,
		PF.Descricao AS Produto,
		FORNE.Nome AS NomeFornecedor,
		PF.PrazoAcordadoDias
	FROM Recebimento R WITH (NOLOCK)
	JOIN Funcionarios F WITH (NOLOCK) ON F.Id = R.IdResponsavel
	JOIN Notas N WITH (NOLOCK) ON N.Id = R.IdNota
	JOIN PedidoFornecedor PF WITH (NOLOCK) ON PF.Id = R.IdPedidoFornecedor
	JOIN Fornecedores FORNE WITH (NOLOCK) ON FORNE.Id = PF.FornecedorId
	WHERE 
		R.Dia BETWEEN @InicioData AND @FimData
	ORDER BY R.Dia DESC
	`

	rows, err := s.db.Query(query,
		sql.Named("InicioData", inicioData),
		sql.Named("FimData", fimData),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar recebimentos por data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Recebimento
		if err := rows.Scan(
			&r.Id, &r.Dia, &r.UrlImagem, &r.IdResponsavel,
			&r.NomeResponsavel, &r.Quantidade, &r.Peso, &r.Valor,
			&r.Vencimento, &r.IdNota, &r.IdPedidoFornecedor,
			&r.NumeroNota, &r.Produto,
			&r.NomeFornecedor, &r.PrazoAcordadoDias,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear recebimento: %w", err)
		}
		recebimentos = append(recebimentos, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas: %w", err)
	}

	return recebimentos, nil
}

func (s *SQLStr) ValidarRecebimento(recebimento *models.Recebimento) (bool, string, error) {
	var prazoDias int
	err := s.db.QueryRow(`
        SELECT PrazoAcordadoDias 
        FROM PedidoFornecedor 
        WHERE Id = @IdPedidoFornecedor`,
		sql.Named("IdPedidoFornecedor", recebimento.IdPedidoFornecedor),
	).Scan(&prazoDias)
	if err != nil {
		return false, "", fmt.Errorf("erro ao consultar prazo acordado: %w", err)
	}

	var notaExiste int
	err = s.db.QueryRow(`
        SELECT COUNT(*) 
        FROM Notas 
        WHERE Id = @IdNota`,
		sql.Named("IdNota", recebimento.IdNota),
	).Scan(&notaExiste)
	if err != nil {
		return false, "", fmt.Errorf("erro ao consultar nota: %w", err)
	}
	if notaExiste == 0 {
		return false, "Nota não encontrada", nil
	}

	extractDate := func(dateTimeStr string) (string, error) {
		if len(dateTimeStr) == 10 && dateTimeStr[4] == '-' && dateTimeStr[7] == '-' {
			return dateTimeStr, nil
		}
		if t, err := time.Parse(time.RFC3339, dateTimeStr); err == nil {
			return t.Format("2006-01-02"), nil
		}
		if t, err := time.Parse("2006-01-02T15:04:05.999999", dateTimeStr); err == nil {
			return t.Format("2006-01-02"), nil
		}
		if t, err := time.Parse("2006-01-02", dateTimeStr); err == nil {
			return t.Format("2006-01-02"), nil
		}

		return "", fmt.Errorf("formato de data inválido: %s", dateTimeStr)
	}

	diaStr, err := extractDate(recebimento.Dia)
	if err != nil {
		return false, "", fmt.Errorf("formato de data inválido em Dia: %w", err)
	}

	vencimentoStr, err := extractDate(recebimento.Vencimento)
	if err != nil {
		return false, "", fmt.Errorf("formato de data inválido em Vencimento: %w", err)
	}

	diaRecebimento, err := time.Parse("2006-01-02", diaStr)
	if err != nil {
		return false, "", fmt.Errorf("erro ao converter data do recebimento: %w", err)
	}

	vencimento, err := time.Parse("2006-01-02", vencimentoStr)
	if err != nil {
		return false, "", fmt.Errorf("erro ao converter data de vencimento: %w", err)
	}

	dataMinima := diaRecebimento.AddDate(0, 0, prazoDias)
	if vencimento.Before(dataMinima) {
		return false, fmt.Sprintf(
			"Data de vencimento está antes do prazo acordado de %d dias (mínimo permitido: %s)",
			prazoDias, dataMinima.Format("02/01/2006")), nil
	}

	return true, "", nil
}

func (r *SQLStr) BuscarDadosRecebimentoPorNumeroNota(numeroNota string) (*models.DadosRecebimentoNota, error) {
	var dados models.DadosRecebimentoNota

	query := `
		SELECT 
			n.Id AS IdNota,
			n.NumeroNota,
			n.IdPedidoFornecedor,
			pf.Descricao AS Produto,
			pf.PrazoAcordadoDias,	
			f.Nome AS Fornecedor
		FROM Notas n
		INNER JOIN PedidoFornecedor pf ON pf.Id = n.IdPedidoFornecedor
		INNER JOIN Fornecedores f ON f.Id = pf.FornecedorId
		WHERE n.NumeroNota = @NumeroNota
	`

	err := r.db.QueryRow(query, sql.Named("NumeroNota", numeroNota)).Scan(
		&dados.IdNota,
		&dados.NumeroNota,
		&dados.IdPedidoFornecedor,
		&dados.Produto,
		&dados.PrazoAcordadoDias,
		&dados.Fornecedor,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados da nota: %w", err)
	}

	return &dados, nil
}
