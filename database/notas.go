package sql

import (
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"time"
)

func (s *SQLStr) CriarNota(nota *models.Nota) error {
	var exists bool
	checkQuery := `SELECT 1 FROM Notas WHERE NumeroNota = @NumeroNota`
	err := s.db.QueryRow(checkQuery, sql.Named("NumeroNota", nota.NumeroNota)).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("erro ao verificar existencia da nota: %w", err)
	}
	if exists {
		return fmt.Errorf("já existe uma nota com o número %s", nota.NumeroNota)
	}

	query := `
        INSERT INTO Notas (
            Tipo, Valor, IdFuncionario, Url_Imagem, Dia, Descricao,
            NumeroNota, DataEmissao, IdFornecedor, IdPedidoFornecedor
        ) 
        VALUES (
            @Tipo, @Valor, @IdFuncionario, @UrlImagem, @Dia, @Descricao,
            @NumeroNota, @DataEmissao, @IdFornecedor, @IdPedidoFornecedor
        )
    `
	_, err = s.db.Exec(query,
		sql.Named("Tipo", nota.Tipo),
		sql.Named("Valor", nota.Valor),
		sql.Named("IdFuncionario", nota.IdFuncionario),
		sql.Named("UrlImagem", nota.UrlImagem),
		sql.Named("Dia", nota.Dia),
		sql.Named("Descricao", nota.Descricao),
		sql.Named("NumeroNota", nota.NumeroNota),
		sql.Named("DataEmissao", nota.DataEmissao),
		sql.Named("IdFornecedor", nota.IdFornecedor),
		sql.Named("IdPedidoFornecedor", nota.IdPedidoFornecedor),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir nota: %w", err)
	}
	return nil
}

func (s *SQLStr) ListarNotas(page int) ([]models.Nota, error) {
	var notas []models.Nota

	limit := 10
	offset := (page - 1) * limit

	query := `
    SELECT 
        n.Id, n.Tipo, n.Valor, n.IdFuncionario, n.Url_Imagem, n.Dia, n.Descricao,
        n.NumeroNota, n.DataEmissao, n.IdFornecedor, n.IdPedidoFornecedor,
        f.Nome AS NomeFuncionario, fo.Nome AS NomeFornecedor
    FROM Notas n WITH (NOLOCK)
    JOIN Funcionarios f WITH (NOLOCK) ON n.IdFuncionario = f.Id
    LEFT JOIN Fornecedores fo WITH (NOLOCK) ON n.IdFornecedor = fo.Id
    ORDER BY n.Dia DESC
    OFFSET @Offset ROWS
    FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar notas: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Nota
		if err := rows.Scan(
			&n.Id, &n.Tipo, &n.Valor, &n.IdFuncionario, &n.UrlImagem, &n.Dia, &n.Descricao,
			&n.NumeroNota, &n.DataEmissao, &n.IdFornecedor, &n.IdPedidoFornecedor,
			&n.NomeFuncionario, &n.NomeFornecedor,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear nota: %w", err)
		}
		notas = append(notas, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return notas, nil
}

func (s *SQLStr) FiltrarDataNota(inicioData time.Time, fimData time.Time) ([]models.Nota, error) {
	var notas []models.Nota

	query := `
		SELECT 
			n.Id, n.Tipo, n.Valor, n.IdFuncionario, n.Url_Imagem, n.Dia, n.Descricao,
			n.NumeroNota, n.DataEmissao, n.IdFornecedor, n.IdPedidoFornecedor,
			f.Nome AS NomeFuncionario, fo.Nome AS NomeFornecedor
		FROM Notas n WITH (NOLOCK)
		JOIN Funcionarios f WITH (NOLOCK) ON n.IdFuncionario = f.Id
		LEFT JOIN Fornecedores fo WITH (NOLOCK) ON n.IdFornecedor = fo.Id
		WHERE n.Dia BETWEEN @InicioData AND @FimData
		ORDER BY n.Dia DESC
	`

	rows, err := s.db.Query(query, sql.Named("InicioData", inicioData), sql.Named("FimData", fimData))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar notas: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Nota
		if err := rows.Scan(
			&n.Id, &n.Tipo, &n.Valor, &n.IdFuncionario, &n.UrlImagem, &n.Dia, &n.Descricao,
			&n.NumeroNota, &n.DataEmissao, &n.IdFornecedor, &n.IdPedidoFornecedor,
			&n.NomeFuncionario, &n.NomeFornecedor,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear notas: %w", err)
		}
		notas = append(notas, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return notas, nil
}

func (s *SQLStr) BuscarNotasPorNumero(numero string) ([]models.Nota, error) {
	var notas []models.Nota

	query := `
		SELECT 
			n.Id, n.Tipo, n.Valor, n.IdFuncionario, n.Url_Imagem, n.Dia, n.Descricao,
			n.NumeroNota, n.DataEmissao, n.IdFornecedor, n.IdPedidoFornecedor,
			f.Nome AS NomeFuncionario, fo.Nome AS NomeFornecedor
		FROM Notas n WITH (NOLOCK)
		JOIN Funcionarios f WITH (NOLOCK) ON n.IdFuncionario = f.Id
		LEFT JOIN Fornecedores fo WITH (NOLOCK) ON n.IdFornecedor = fo.Id
		WHERE n.NumeroNota LIKE @Numero
		ORDER BY n.Dia DESC
	`

	likePattern := "%" + numero + "%"
	rows, err := s.db.Query(query, sql.Named("Numero", likePattern))
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar por número da nota: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n models.Nota
		if err := rows.Scan(
			&n.Id, &n.Tipo, &n.Valor, &n.IdFuncionario, &n.UrlImagem, &n.Dia, &n.Descricao,
			&n.NumeroNota, &n.DataEmissao, &n.IdFornecedor, &n.IdPedidoFornecedor,
			&n.NomeFuncionario, &n.NomeFornecedor,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear nota: %w", err)
		}
		notas = append(notas, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return notas, nil
}
