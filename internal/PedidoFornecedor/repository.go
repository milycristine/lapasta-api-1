package pedidofornecedor

import (
	"database/sql"
	"fmt"
	database "lapasta/database"
	"lapasta/internal/models"
)

type PedidoFornecedorRepository interface {
	CriarPedidoFornecedor(pedido *models.PedidoFornecedor) error
	ListarPedidosPorFornecedor(fornecedorId int) ([]models.PedidoFornecedor, error)
}

type pedidoFornecedorRepository struct {
	db *sql.DB
}

func NovoPedidoFornecedorRepository(conn *database.SQLStr) PedidoFornecedorRepository {
	return &pedidoFornecedorRepository{
		db: conn.DB(),
	}
}


func (r *pedidoFornecedorRepository) CriarPedidoFornecedor(p *models.PedidoFornecedor) error {
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
func (r *pedidoFornecedorRepository) ListarPedidosPorFornecedor(fornecedorId int) ([]models.PedidoFornecedor, error) {
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