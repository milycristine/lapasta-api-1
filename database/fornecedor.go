package sql

import (
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"strings"
)

func (r *SQLStr) CriarFornecedor(f *models.Fornecedor) error {
	var existe bool
	checkQuery := `SELECT * FROM Fornecedores WHERE CNPJ = @CNPJ`
	err := r.db.QueryRow(checkQuery, sql.Named("CNPJ", f.CNPJ)).Scan(&existe)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("erro ao verificar CNPJ existente: %w", err)
	}
	if existe {
		return fmt.Errorf("já existe um fornecedor com o CNPJ %s", f.CNPJ)
	}

	query := `
		INSERT INTO Fornecedores (Nome, CNPJ, Email, Telefone)
		VALUES (@Nome, @CNPJ, @Email, @Telefone)
	`
	_, err = r.db.Exec(query,
		sql.Named("Nome", f.Nome),
		sql.Named("CNPJ", f.CNPJ),
		sql.Named("Email", f.Email),
		sql.Named("Telefone", f.Telefone),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir fornecedor: %w", err)
	}

	return nil
}

func (r *SQLStr) ListarFornecedores(page int) ([]models.Fornecedor, error) {
	var fornecedores []models.Fornecedor

	limit := 10
	offset := (page - 1) * limit

	query := `
			SELECT Id, Nome, CNPJ, Email, Telefone
			FROM Fornecedores WITH (NOLOCK)
			ORDER BY Nome
			OFFSET @Offset ROWS
			FETCH NEXT @Limit ROWS ONLY
		`

	rows, err := r.db.Query(query,
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar fornecedores: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Fornecedor
		if err := rows.Scan(&f.Id, &f.Nome, &f.CNPJ, &f.Email, &f.Telefone); err != nil {
			return nil, fmt.Errorf("erro ao escanear fornecedor: %w", err)
		}
		fornecedores = append(fornecedores, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro na iteração de fornecedores: %w", err)
	}

	return fornecedores, nil
}

func (r *SQLStr) BuscarFornecedorPorCNPJouNome(valor string) ([]models.Fornecedor, error) {
	var fornecedores []models.Fornecedor

	query := `
			SELECT Id, Nome, CNPJ, Email, Telefone
			FROM Fornecedores WITH (NOLOCK)
			WHERE CNPJ = @Valor OR Nome LIKE @LikeValor
		`

	rows, err := r.db.Query(query,
		sql.Named("Valor", valor),
		sql.Named("LikeValor", "%"+valor+"%"),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar fornecedores: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var f models.Fornecedor
		if err := rows.Scan(&f.Id, &f.Nome, &f.CNPJ, &f.Email, &f.Telefone); err != nil {
			return nil, fmt.Errorf("erro ao escanear fornecedor: %w", err)
		}
		fornecedores = append(fornecedores, f)
	}

	if len(fornecedores) == 0 {
		return nil, fmt.Errorf("nenhum fornecedor encontrado com '%s'", valor)
	}

	return fornecedores, nil
}

func (r *SQLStr) CriarPedidoFornecedor(p *models.PedidoFornecedor) error {
	var existe bool
	checkQuery := `
		SELECT 1 FROM PedidoFornecedor 
		WHERE FornecedorId = @FornecedorId AND DataPedido = @DataPedido AND Descricao = @Descricao
	`
	err := r.db.QueryRow(checkQuery,
		sql.Named("FornecedorId", p.FornecedorId),
		sql.Named("DataPedido", p.DataPedido),
		sql.Named("Descricao", p.Descricao),
	).Scan(&existe)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("erro ao verificar duplicidade de pedido: %w", err)
	}
	if existe {
		return fmt.Errorf("já existe um pedido com essa descrição para o fornecedor na mesma data")
	}

	query := `
		INSERT INTO PedidoFornecedor (FornecedorId, DataPedido, PrazoAcordadoDias, Descricao)
		VALUES (@FornecedorId, @DataPedido, @PrazoAcordadoDias, @Descricao)
	`
	_, err = r.db.Exec(query,
		sql.Named("FornecedorId", p.FornecedorId),
		sql.Named("DataPedido", p.DataPedido),
		sql.Named("PrazoAcordadoDias", p.PrazoAcordadoDias),
		sql.Named("Descricao", p.Descricao),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir pedido de fornecedor: %w", err)
	}

	return nil
}

func (r *SQLStr) ListarPedidosFornecedor(fornecedorId int) ([]models.PedidoFornecedor, error) {
	var pedidos []models.PedidoFornecedor

	query := `
			SELECT Id, FornecedorId, DataPedido, PrazoAcordadoDias, Descricao
			FROM PedidoFornecedor WITH (NOLOCK)
			WHERE FornecedorId = @FornecedorId
			ORDER BY DataPedido DESC
		`

	rows, err := r.db.Query(query, sql.Named("FornecedorId", fornecedorId))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pedidos do fornecedor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PedidoFornecedor
		if err := rows.Scan(&p.Id, &p.FornecedorId, &p.DataPedido, &p.PrazoAcordadoDias, &p.Descricao); err != nil {
			return nil, fmt.Errorf("erro ao escanear pedido: %w", err)
		}
		pedidos = append(pedidos, p)
	}

	return pedidos, nil
}
func (r *SQLStr) BuscarPedidosFornecedorPorDescricaoOuId(fornecedorId int, valor string) ([]models.PedidoFornecedor, error) {
	var pedidos []models.PedidoFornecedor

	query := `
			SELECT Id, FornecedorId, DataPedido, PrazoAcordadoDias, Descricao
			FROM PedidoFornecedor WITH (NOLOCK)
			WHERE FornecedorId = @FornecedorId
			AND (
				CAST(Id AS NVARCHAR) = @Valor
				OR Descricao LIKE @LikeValor
			)
			ORDER BY DataPedido DESC
		`

	rows, err := r.db.Query(query,
		sql.Named("FornecedorId", fornecedorId),
		sql.Named("Valor", valor),
		sql.Named("LikeValor", "%"+valor+"%"),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos do fornecedor: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.PedidoFornecedor
		if err := rows.Scan(&p.Id, &p.FornecedorId, &p.DataPedido, &p.PrazoAcordadoDias, &p.Descricao); err != nil {
			return nil, fmt.Errorf("erro ao escanear pedido: %w", err)
		}
		pedidos = append(pedidos, p)
	}

	if len(pedidos) == 0 {
		return nil, fmt.Errorf("nenhum pedido encontrado com '%s'", valor)
	}

	return pedidos, nil
}

func (r *SQLStr) EditarFornecedor(f *models.Fornecedor) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	var (
		setClauses []string
		args       []any
	)

	if f.Nome != "" {
		setClauses = append(setClauses, "Nome = @Nome")
		args = append(args, sql.Named("Nome", f.Nome))
	}
	if f.Email != "" {
		setClauses = append(setClauses, "Email = @Email")
		args = append(args, sql.Named("Email", f.Email))
	}
	if f.Telefone != "" {
		setClauses = append(setClauses, "Telefone = @Telefone")
		args = append(args, sql.Named("Telefone", f.Telefone))
	}

	if len(setClauses) == 0 {
		tx.Rollback()
		return fmt.Errorf("nenhum campo para atualizar")
	}

	args = append(args, sql.Named("CNPJ", f.CNPJ))
	query := fmt.Sprintf("UPDATE Fornecedores SET %s WHERE CNPJ = @CNPJ", strings.Join(setClauses, ", "))

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao atualizar fornecedor: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao finalizar transação: %w", err)
	}

	return nil
}
