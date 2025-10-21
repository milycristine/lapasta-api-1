package sql

import (
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"strconv"
	"time"
)

func (s *SQLStr) CriarRecebimento(r *models.Recebimento) error {
	var existe int
	checkQuery := `SELECT 1 FROM Recebimento WHERE IdNota = @IdNota`
	err := s.db.QueryRow(checkQuery, sql.Named("IdNota", r.IdNota)).Scan(&existe)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("erro ao verificar recebimentos existentes: %w", err)
	}
	if existe == 1 {
		return fmt.Errorf("já existe um recebimento para a nota com ID %d", r.IdNota)
	}
	if r.UrlImagem == nil {
		r.UrlImagem = new(string)
	}

	query := `
		INSERT INTO Recebimento (
			Dia, UrlImagem, IdResponsavel, Quantidade, Peso, Valor,
			Vencimento, IdNota, IdPedidoFornecedor, FormaPagamento, OutroPrazoDias
		)
		OUTPUT INSERTED.Id
		VALUES (
			@Dia, @UrlImagem, @IdResponsavel, @Quantidade, @Peso, @Valor,
			@Vencimento, @IdNota, @IdPedidoFornecedor, @FormaPagamento, @OutroPrazoDias
		)
	`

	var id sql.NullInt64
	err = s.db.QueryRow(
		query,
		sql.Named("Dia", r.Dia),
		sql.Named("UrlImagem", r.UrlImagem),
		sql.Named("IdResponsavel", r.IdResponsavel),
		sql.Named("Quantidade", r.Quantidade),
		sql.Named("Peso", r.Peso),
		sql.Named("Valor", r.Valor),
		sql.Named("Vencimento", r.Vencimento),
		sql.Named("IdNota", r.IdNota),
		sql.Named("IdPedidoFornecedor", r.IdPedidoFornecedor),
		sql.Named("FormaPagamento", r.FormaPagamento),
		sql.Named("OutroPrazoDias", r.OutroPrazoDias),
	).Scan(&id)

	if err != nil {
		return fmt.Errorf("erro ao inserir recebimento: %w", err)
	}

	if !id.Valid {
		return fmt.Errorf("erro: ID do recebimento não retornado após inserção")
	}

	r.Id = int(id.Int64)
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
			N.Descricao AS Produto,
			COALESCE(FORNE.Nome, '') AS NomeFornecedor,
			R.FormaPagamento,
			R.OutroPrazoDias
		FROM Recebimento R WITH (NOLOCK)
		JOIN Funcionarios F WITH (NOLOCK) ON F.Id = R.IdResponsavel
		JOIN Notas N WITH (NOLOCK) ON N.Id = R.IdNota
		LEFT JOIN Fornecedores FORNE WITH (NOLOCK) ON FORNE.Id = N.IdFornecedor
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
		var formaPagamento sql.NullString
		var outroPrazoDias sql.NullInt64

		if err := rows.Scan(
			&r.Id,
			&r.Dia,
			&r.UrlImagem,
			&r.IdResponsavel,
			&r.NomeResponsavel,
			&r.Quantidade,
			&r.Peso,
			&r.Valor,
			&r.Vencimento,
			&r.IdNota,
			&r.IdPedidoFornecedor,
			&r.NumeroNota,
			&r.Produto,
			&r.NomeFornecedor,
			&formaPagamento,
			&outroPrazoDias,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear recebimento: %w", err)
		}

		if formaPagamento.Valid {
			r.FormaPagamento = formaPagamento.String
		} else {
			r.FormaPagamento = ""
		}

		if outroPrazoDias.Valid {
			val := int(outroPrazoDias.Int64)
			r.OutroPrazoDias = &val
		} else {
			r.OutroPrazoDias = nil
		}

		recebimentos = append(recebimentos, r)
	}

	return recebimentos, nil
}
func (s *SQLStr) FiltrarDataRecebimentos(inicioData, fimData time.Time) ([]models.Recebimento, error) {
	var recebimentos []models.Recebimento

	query := `
		SELECT 
		    R.Id, R.Dia, R.UrlImagem, R.IdResponsavel, 
		    F.Nome AS NomeResponsavel, R.Quantidade, R.Peso, R.Valor, 
		    R.Vencimento, R.IdNota, R.IdPedidoFornecedor,
		    N.NumeroNota,
		    N.Descricao AS Produto,
		    COALESCE(FORNE.Nome, '') AS NomeFornecedor,
		    R.FormaPagamento,
		    R.OutroPrazoDias
		FROM Recebimento R WITH (NOLOCK)
		JOIN Funcionarios F WITH (NOLOCK) ON F.Id = R.IdResponsavel
		JOIN Notas N WITH (NOLOCK) ON N.Id = R.IdNota
		LEFT JOIN Fornecedores FORNE WITH (NOLOCK) ON FORNE.Id = N.IdFornecedor
		WHERE R.Dia BETWEEN @InicioData AND @FimData
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
		var formaPagamento sql.NullString
		var outroPrazoDias sql.NullInt64
		var produto sql.NullString
		var nomeFornecedor sql.NullString

		if err := rows.Scan(
			&r.Id, &r.Dia, &r.UrlImagem, &r.IdResponsavel,
			&r.NomeResponsavel, &r.Quantidade, &r.Peso, &r.Valor,
			&r.Vencimento, &r.IdNota, &r.IdPedidoFornecedor,
			&r.NumeroNota, &produto,
			&nomeFornecedor, &formaPagamento,
			&outroPrazoDias,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear recebimento: %w", err)
		}

		if produto.Valid {
			r.Produto = &produto.String
		} else {
			r.Produto = nil
		}

		if nomeFornecedor.Valid {
			r.NomeFornecedor = &nomeFornecedor.String
		} else {
			r.NomeFornecedor = nil
		}

		if formaPagamento.Valid {
			r.FormaPagamento = formaPagamento.String
		} else {
			r.FormaPagamento = ""
		}

		if outroPrazoDias.Valid {
			val := int(outroPrazoDias.Int64)
			r.OutroPrazoDias = &val
		} else {
			r.OutroPrazoDias = nil
		}

		recebimentos = append(recebimentos, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante iteração das linhas: %w", err)
	}

	return recebimentos, nil
}

func (r *SQLStr) BuscarDadosRecebimentoPorNumeroNota(numeroNota string) (*models.DadosRecebimentoNota, error) {
	var dados models.DadosRecebimentoNota
	var outroPrazoDias *int

	query := `
	SELECT 
		n.Id AS IdNota,
		n.NumeroNota,
		n.IdPedidoFornecedor,
		pf.Descricao AS Produto,
		pf.FormaPagamento,
		pf.OutroPrazoDias,
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
		&dados.FormaPagamento,
		&outroPrazoDias,
		&dados.Fornecedor,
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados da nota: %w", err)
	}

	dados.OutroPrazoDias = outroPrazoDias
	return &dados, nil
}

func (r *SQLStr) TotalRecebimentosAvistaMesAtual() (float64, error) {
	query := `
        SELECT COALESCE(SUM(Valor), 0)
        FROM Recebimento
        WHERE LOWER(FormaPagamento) IN ('avista', 'pix', 'dinheiro', 'debito')
        AND Id NOT IN (SELECT RecebimentoId FROM BoletosRecebidos)
        AND YEAR(Dia) = YEAR(GETDATE())
        AND MONTH(Dia) = MONTH(GETDATE())
    `
	var total float64
	err := r.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("erro ao calcular total de recebimentos à vista/pix: %w", err)
	}
	return total, nil
}

func (r *SQLStr) ListarRecebimentosAvistaMesAtual() ([]models.Recebimento, error) {
	query := `
        SELECT 
            R.Id,
            R.Dia,
            PF.Descricao AS Produto,
            R.Valor,
            R.FormaPagamento,
            F.Nome AS NomeResponsavel,
            FORNE.Nome AS NomeFornecedor
        FROM Recebimento R
        JOIN Funcionarios F ON F.Id = R.IdResponsavel
        JOIN PedidoFornecedor PF ON PF.Id = R.IdPedidoFornecedor
        JOIN Fornecedores FORNE ON FORNE.Id = PF.FornecedorId
        WHERE LOWER(R.FormaPagamento) IN ('avista', 'pix', 'dinheiro', 'debito')
        AND R.Id NOT IN (SELECT RecebimentoId FROM BoletosRecebidos)
        AND YEAR(R.Dia) = YEAR(GETDATE())
        AND MONTH(R.Dia) = MONTH(GETDATE())
        ORDER BY R.Dia DESC
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar recebimentos à vista/pix: %w", err)
	}
	defer rows.Close()

	var recebimentos []models.Recebimento
	for rows.Next() {
		var rcv models.Recebimento
		var formaPagamento sql.NullString

		if err := rows.Scan(
			&rcv.Id,
			&rcv.Dia,
			&rcv.Produto,
			&rcv.Valor,
			&formaPagamento,
			&rcv.NomeResponsavel,
			&rcv.NomeFornecedor,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear recebimento à vista: %w", err)
		}

		if formaPagamento.Valid {
			rcv.FormaPagamento = formaPagamento.String
		} else {
			rcv.FormaPagamento = ""
		}

		recebimentos = append(recebimentos, rcv)
	}

	return recebimentos, nil
}

func (s *SQLStr) ValidarRecebimento(r *models.Recebimento) (bool, string, error) {
	parseDate := func(dateStr string) (time.Time, error) {
		formats := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02T15:04:05.999999",
			"2006-01-02T15:04:05",
			"2006-01-02",
		}
		for _, f := range formats {
			if t, err := time.Parse(f, dateStr); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("formato inválido de data: %s", dateStr)
	}

	dia, err := parseDate(r.Dia)
	if err != nil {
		return false, "", fmt.Errorf("data de recebimento inválida: %w", err)
	}
	venc, err := parseDate(r.Vencimento)
	if err != nil {
		return false, "", fmt.Errorf("data de vencimento inválida: %w", err)
	}

	toDate := func(t time.Time) time.Time {
		loc := t.Location()
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	}

	var dias int
	switch r.FormaPagamento {
	case "À vista", "avista", "pix", "dinheiro", "debito":
		dias = 0
	case "7", "14", "21":
		dias, _ = strconv.Atoi(r.FormaPagamento)
	case "Outro":
		if r.OutroPrazoDias == nil || *r.OutroPrazoDias <= 0 {
			return false, "Tipo de pagamento 'Outro' informado sem número de dias", nil
		}
		dias = *r.OutroPrazoDias
	default:
		return false, fmt.Sprintf("Tipo de pagamento inválido: %s", r.FormaPagamento), nil
	}

	dataMinima := toDate(dia).AddDate(0, 0, dias)
	vencData := toDate(venc)

	if vencData.Before(dataMinima) {
		return false, fmt.Sprintf(
				"Data de vencimento antes do permitido (%d dias, mínimo: %s)",
				dias, dataMinima.Format("02/01/2006")),
			nil
	}

	return true, "", nil
}
