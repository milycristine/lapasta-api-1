package sql

import (
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"log"
	"strings"
	"time"
)

func (r *SQLStr) CriarBoletoRecebido(b *models.Boleto) error {
	var receb models.Recebimento
	queryReceb := `
		SELECT Id, IdPedidoFornecedor, Dia	
		FROM Recebimento WITH (NOLOCK)
		WHERE Id = @RecebimentoId
	`
	err := r.db.QueryRow(queryReceb, sql.Named("RecebimentoId", b.RecebimentoId)).
		Scan(&receb.Id, &receb.IdPedidoFornecedor, &receb.Dia)
	if err != nil {
		log.Printf("Erro ao buscar recebimento: %v", err)
		return fmt.Errorf("erro ao buscar recebimento: %w", err)
	}

	var prazoAcordado int
	queryPedido := `
		SELECT PrazoAcordadoDias
		FROM PedidoFornecedor 
		WHERE Id = @PedidoId
	`
	err = r.db.QueryRow(queryPedido, sql.Named("PedidoId", receb.IdPedidoFornecedor)).Scan(&prazoAcordado)
	if err != nil {
		log.Printf("Erro ao buscar pedido fornecedor: %v", err)
		return fmt.Errorf("erro ao buscar pedido fornecedor: %w", err)
	}

	codigo := b.CodigoBarras
	if len(codigo) == 47 {
		convertido, err := ConverterLinhaDigitavelParaCodigoBarras(codigo)
		if err != nil {
			log.Printf("Erro ao converter linha digitável: %v", err)
			return fmt.Errorf("linha digitável inválida: %w", err)
		}
		codigo = convertido
	}

	var existe bool
	queryCheck := `
		SELECT 1 
		FROM BoletosRecebidos 
		WHERE CodigoBarras = @CodigoBarras AND RecebimentoId = @RecebimentoId
	`
	err = r.db.QueryRow(queryCheck,
		sql.Named("CodigoBarras", codigo),
		sql.Named("RecebimentoId", b.RecebimentoId),
	).Scan(&existe)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Erro ao verificar boleto duplicado: %v", err)
		return fmt.Errorf("erro ao verificar duplicidade: %w", err)
	}
	if existe {
		return fmt.Errorf("boleto com esse código de barras já foi registrado para este recebimento")
	}

	queryInsert := `
		INSERT INTO BoletosRecebidos 
		(RecebimentoId, CodigoBarras, DataCadastro, DataVencimento, Valor, StatusId, DataPagamento)
		VALUES 
		(@RecebimentoId, @CodigoBarras, @DataCadastro, @DataVencimento, @Valor, @StatusId, @DataPagamento)
	`
	_, err = r.db.Exec(queryInsert,
		sql.Named("RecebimentoId", b.RecebimentoId),
		sql.Named("CodigoBarras", codigo),
		sql.Named("DataCadastro", b.DataCadastro),
		sql.Named("DataVencimento", b.DataVencimento),
		sql.Named("Valor", b.Valor),
		sql.Named("StatusId", 1),
		sql.Named("DataPagamento", b.DataPagamento),
	)

	if err != nil {
		log.Printf("Erro ao inserir boleto recebido: %v", err)
		return fmt.Errorf("erro ao inserir boleto recebido: %w", err)
	}

	return nil
}

func (r *SQLStr) ListarBoletosPorFornecedor(fornecedorId int) ([]models.Boleto, error) {
	query := `
		SELECT b.Id, b.RecebimentoId, b.CodigoBarras, b.DataCadastro, b.DataVencimento, b.Valor, b.StatusId, b.DataPagamento
		FROM BoletosRecebidos b
		INNER JOIN Recebimento r ON r.Id = b.RecebimentoId
		INNER JOIN PedidoFornecedor p ON p.Id = r.IdPedidoFornecedor
		WHERE p.FornecedorId = @FornecedorId
		`
	rows, err := r.db.Query(query, sql.Named("FornecedorId", fornecedorId))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos por fornecedor: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		err := rows.Scan(&b.Id, &b.RecebimentoId, &b.CodigoBarras, &b.DataCadastro, &b.DataVencimento, &b.Valor, &b.StatusId, &b.DataPagamento)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear boleto: %w", err)
		}
		boletos = append(boletos, b)
	}
	return boletos, nil
}
func (r *SQLStr) ListarBoletosPorPedido(pedidoId int) ([]models.Boleto, error) {
	query := `
		SELECT b.Id, b.RecebimentoId, b.CodigoBarras, b.DataCadastro, b.DataVencimento, b.Valor, b.StatusId, b.DataPagamento
		FROM BoletosRecebidos b
		INNER JOIN Recebimento r ON r.Id = b.RecebimentoId
		WHERE r.IdPedidoFornecedor = @PedidoId
		`
	rows, err := r.db.Query(query, sql.Named("PedidoId", pedidoId))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos por pedido: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		err := rows.Scan(&b.Id, &b.RecebimentoId, &b.CodigoBarras, &b.DataCadastro, &b.DataVencimento, &b.Valor, &b.StatusId, &b.DataPagamento)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear boleto: %w", err)
		}
		boletos = append(boletos, b)
	}
	return boletos, nil
}
func (r *SQLStr) ListarBoletosDoDia(data string) ([]models.Boleto, error) {
	query := `
		SELECT 
			b.Id,
			f.Nome, 
			b.Valor, 
			b.DataVencimento, 
			b.CodigoBarras, 
			b.StatusId
		FROM 
			BoletosRecebidos b
		INNER JOIN 
			Recebimento r ON r.Id = b.RecebimentoId
		INNER JOIN 
			PedidoFornecedor p ON p.Id = r.IdPedidoFornecedor
		INNER JOIN 
			Fornecedores f ON f.Id = p.FornecedorId
		WHERE 
			b.DataVencimento <= @DataHoje
	`

	dataHoje, err := time.Parse("2006-01-02", data)
	if err != nil {
		return nil, fmt.Errorf("data inválida: %w", err)
	}
	dataHoje = dataHoje.Truncate(24 * time.Hour)

	rows, err := r.db.Query(query, sql.Named("DataHoje", dataHoje.Format("2006-01-02")))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos do dia: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto

	for rows.Next() {
		var boleto models.Boleto
		var statusId int
		var dataVencimentoStr string

		err := rows.Scan(&boleto.Id, &boleto.FornecedorNome, &boleto.Valor, &dataVencimentoStr, &boleto.CodigoBarras, &statusId)
		if err != nil {
			return nil, err
		}

		dataVencimento, err := time.Parse("2006-01-02", dataVencimentoStr)
		if err != nil {
			return nil, fmt.Errorf("erro ao converter data de vencimento: %w", err)
		}

		boleto.DataVencimento = dataVencimento.Format("2006-01-02")

		if statusId == 1 && dataVencimento.Before(dataHoje) {
			err := r.AtualizarStatusBoleto(boleto.Id, 3)
			if err != nil {
				log.Printf("Erro ao atualizar status do boleto %d: %v", boleto.Id, err)
			}
			boleto.StatusId = 3
		} else {
			boleto.StatusId = statusId
		}

		boletos = append(boletos, boleto)
	}

	return boletos, nil
}

func (r *SQLStr) PagarBoleto(codigoBarras string) error {
	query := `
		UPDATE BoletosRecebidos
		SET StatusId = 2, DataPagamento = CURRENT_TIMESTAMP
		WHERE CodigoBarras = @CodigoBarras AND StatusId <> 2 AND DataVencimento <= CONVERT(date, GETDATE())
	`
	result, err := r.db.Exec(query, sql.Named("CodigoBarras", codigoBarras))
	if err != nil {
		return fmt.Errorf("erro ao pagar boleto: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar atualização: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("boleto já pago, vencimento futuro ou não encontrado")
	}

	_, err = r.db.Exec(`
		INSERT INTO PagamentosLog (CodigoBarras, DataHora) 
		VALUES (@CodigoBarras, CURRENT_TIMESTAMP)
	`, sql.Named("CodigoBarras", codigoBarras))
	if err != nil {
		return fmt.Errorf("erro ao registrar log de pagamento: %w", err)
	}

	return nil
}

func (r *SQLStr) ListarBoletosPagos() ([]models.Boleto, error) {
	query := `
		SELECT f.Nome, b.Valor, b.DataVencimento, b.DataPagamento, b.CodigoBarras, b.StatusId
		FROM BoletosRecebidos b
		INNER JOIN Recebimento r ON r.Id = b.RecebimentoId
		INNER JOIN PedidoFornecedor p ON p.Id = r.IdPedidoFornecedor
		INNER JOIN Fornecedores f ON f.Id = p.FornecedorId
		WHERE b.StatusId = 2
		ORDER BY b.DataPagamento DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos pagos: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		err := rows.Scan(&b.FornecedorNome, &b.Valor, &b.DataVencimento, &b.DataPagamento, &b.CodigoBarras, &b.StatusId)
		if err != nil {
			return nil, err
		}
		boletos = append(boletos, b)
	}
	return boletos, nil
}

func (r *SQLStr) ListarBoletosVencidos() ([]models.Boleto, error) {
	query := `
		SELECT f.Nome, b.Valor, b.DataVencimento, b.CodigoBarras, b.StatusId
		FROM BoletosRecebidos b
		INNER JOIN Recebimento r ON r.Id = b.RecebimentoId
		INNER JOIN PedidoFornecedor p ON p.Id = r.IdPedidoFornecedor
		INNER JOIN Fornecedores f ON f.Id = p.FornecedorId
		WHERE b.StatusId = 3
		ORDER BY b.DataVencimento
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos vencidos: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		err := rows.Scan(&b.FornecedorNome, &b.Valor, &b.DataVencimento, &b.CodigoBarras, &b.StatusId)
		if err != nil {
			return nil, err
		}
		boletos = append(boletos, b)
	}
	return boletos, nil
}

func (r *SQLStr) ListarBoletosPendentes() ([]models.Boleto, error) {
	query := `
		SELECT f.Nome, b.Valor, b.DataVencimento, b.CodigoBarras, b.StatusId
		FROM BoletosRecebidos b
		INNER JOIN Recebimento r ON r.Id = b.RecebimentoId
		INNER JOIN PedidoFornecedor p ON p.Id = r.IdPedidoFornecedor
		INNER JOIN Fornecedores f ON f.Id = p.FornecedorId
		WHERE b.StatusId = 1
		ORDER BY b.DataVencimento
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos pendentes: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		err := rows.Scan(&b.FornecedorNome, &b.Valor, &b.DataVencimento, &b.CodigoBarras, &b.StatusId)
		if err != nil {
			return nil, err
		}
		boletos = append(boletos, b)
	}
	return boletos, nil
}

func ConverterLinhaDigitavelParaCodigoBarras(linha string) (string, error) {
	linha = strings.ReplaceAll(linha, " ", "")
	linha = strings.ReplaceAll(linha, ".", "")

	if len(linha) != 47 {
		return "", fmt.Errorf("linha digitável inválida: deve conter 47 dígitos")
	}

	codigoBarras := linha[0:4] +
		linha[32:33] +
		linha[33:47] +
		linha[4:9] + linha[10:20] +
		linha[21:31]

	return codigoBarras, nil
}

func (s *SQLStr) AtualizarStatusBoleto(id int, statusId int) error {
	query := `
		UPDATE BoletosRecebidos  
		SET StatusId = @StatusId 
		WHERE Id = @Id
	`
	_, err := s.db.Exec(query,
		sql.Named("Id", id),
		sql.Named("StatusId", statusId))

	if err != nil {
		return fmt.Errorf("erro ao atualizar pagamento: %w", err)
	}
	return nil
}

func (r *SQLStr) TotalBoletosDia(data string) (float64, error) {
	query := `
        SELECT COALESCE(SUM(Valor), 0)
        FROM BoletosRecebidos
        WHERE StatusId = 1
        AND DataVencimento = CONVERT(date, @Data)
    `
	var total float64
	err := r.db.QueryRow(query, sql.Named("Data", data)).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("erro ao calcular total do dia: %w", err)
	}
	return total, nil
}

func (r *SQLStr) TotalBoletosAtrasados() (float64, error) {
	query := `
        SELECT COALESCE(SUM(Valor), 0)
        FROM BoletosRecebidos
        WHERE StatusId = 3
    `
	var total float64
	err := r.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("erro ao calcular total atrasado: %w", err)
	}
	return total, nil
}

func (r *SQLStr) TotalBoletosPendentesMesAtual() (float64, error) {
	query := `
        SELECT COALESCE(SUM(Valor), 0)
        FROM BoletosRecebidos
        WHERE StatusId = 1
        AND YEAR(DataVencimento) = YEAR(GETDATE())
        AND MONTH(DataVencimento) = MONTH(GETDATE())
    `
	var total float64
	err := r.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("erro ao calcular total pendente do mês: %w", err)
	}
	return total, nil 
}
func (r *SQLStr) TotalBoletosPagosMesAtual() (float64, error) {
	query := `
        SELECT COALESCE(SUM(Valor), 0)
        FROM BoletosRecebidos
        WHERE StatusId = 2 
        AND YEAR(DataPagamento) = YEAR(GETDATE())
        AND MONTH(DataPagamento) = MONTH(GETDATE())
    `
	var total float64
	err := r.db.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("erro ao calcular total pago do mês: %w", err)
	}
	return total, nil
}
