package pagamento

import (
	"database/sql"
	"fmt"
	database "lapasta/database"
	"lapasta/internal/models"
	"time"
)

type PagamentoRepository interface {
	CriarPagamento(pagamento *models.Pagamento) error
	ListarPagamentos(page int) ([]models.Pagamento, error)
	ListarPagamentosPorMes(mes int, page int) ([]models.Pagamento, error)
	ListarPagamentosPorDia(date string, page int) ([]models.Pagamento, error)
	AtualizarStatusPagamento(id int, statusId int) error
}

type pagamentoRepository struct {
	db *sql.DB
}

func NovoPagamentoRepository(conn *database.SQLStr) PagamentoRepository {
	return &pagamentoRepository{
		db: conn.DB(),
	}
}

func ProximoDiaUtil(data time.Time) time.Time {
	weekday := data.Weekday()
	if weekday == time.Saturday {
		return data.AddDate(0, 0, 2)
	} else if weekday == time.Sunday {
		return data.AddDate(0, 0, 1)
	}
	return data
}

func (s *pagamentoRepository) CriarPagamento(pagamento *models.Pagamento) error {
	hoje := time.Now()
	var dataPagamento time.Time

	if hoje.Day() < 5 {
		dataPagamento = time.Date(hoje.Year(), hoje.Month(), 5, 0, 0, 0, 0, hoje.Location())
	} else if hoje.Day() < 20 {
		dataPagamento = time.Date(hoje.Year(), hoje.Month(), 20, 0, 0, 0, 0, hoje.Location())
	} else {
		dataPagamento = time.Date(hoje.Year(), hoje.Month()+1, 5, 0, 0, 0, 0, hoje.Location())
	}

	dataPagamento = ProximoDiaUtil(dataPagamento)

	query := `
		INSERT INTO Pagamentos (FuncionarioId, DataPagamento, Valor, StatusId) 
		VALUES (@FuncionarioId, @DataPagamento, @Valor, 2)
	`
	_, err := s.db.Exec(query,
		sql.Named("FuncionarioId", pagamento.FuncionarioId),
		sql.Named("DataPagamento", dataPagamento),
		sql.Named("Valor", pagamento.Valor))

	if err != nil {
		return fmt.Errorf("erro ao inserir pagamento: %w", err)
	}
	return nil
}

func (s *pagamentoRepository) ListarPagamentos(page int) ([]models.Pagamento, error) {
	var pagamentos []models.Pagamento
	hoje := time.Now()

	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT 
			p.Id, 
			p.FuncionarioId, 
			p.DataPagamento, 
			p.Valor, 
			p.StatusId, 
			f.Cpf, 
			f.Nome,
			f.Sobrenome,
			f.ChavePix
		FROM 
			Pagamentos p WITH (NOLOCK)
		JOIN 
			Funcionarios f WITH (NOLOCK) ON p.FuncionarioId = f.Id
		ORDER BY 
			p.DataPagamento DESC
		OFFSET @Offset ROWS
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar pagamentos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Pagamento
		if err := rows.Scan(&p.Id, &p.FuncionarioId, &p.DataPagamento, &p.Valor, &p.StatusId, &p.CpfFuncionario, &p.NomeFuncionario, &p.SobrenomeFuncionario, &p.ChavePixFuncionario); err != nil {
			return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
		}
		if p.StatusId == 2 && p.DataPagamento.Before(hoje.Truncate(24*time.Hour)) {
			p.StatusId = 3
			s.AtualizarStatusPagamento(p.Id, 3)
		}
		pagamentos = append(pagamentos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return pagamentos, nil
}
func (s *pagamentoRepository) ListarPagamentosPorDia(date string, page int) ([]models.Pagamento, error) {
	var pagamentos []models.Pagamento

	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT 
			p.Id, 
			p.FuncionarioId, 
			p.DataPagamento, 
			p.Valor, 
			p.StatusId, 
			f.Cpf, 
			f.Nome,
			f.Sobrenome,
			f.ChavePix
		FROM 
			Pagamentos p WITH (NOLOCK)
		JOIN 
			Funcionarios f WITH (NOLOCK) ON p.FuncionarioId = f.Id
		WHERE 
			p.DataPagamento = @DiaPagamento
		ORDER BY 
			p.DataPagamento DESC
		OFFSET @Offset ROWS
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(query, sql.Named("DiaPagamento", date), sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar pagamentos para o dia %s: %w", date, err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Pagamento
		if err := rows.Scan(&p.Id, &p.FuncionarioId, &p.DataPagamento, &p.Valor, &p.StatusId, &p.CpfFuncionario, &p.NomeFuncionario, &p.SobrenomeFuncionario, &p.ChavePixFuncionario); err != nil {
			return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
		}
		pagamentos = append(pagamentos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return pagamentos, nil
}

func (s *pagamentoRepository) AtualizarStatusPagamento(id int, statusId int) error {
	query := `
		UPDATE Pagamentos 
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

func (s *pagamentoRepository) ListarPagamentosPorMes(mes int, page int) ([]models.Pagamento, error) {
	var pagamentos []models.Pagamento

	limit := 10
	offset := (page - 1) * limit
	anoAtual := time.Now().Year()

	query := `
        SELECT 
            p.Id, 
            p.FuncionarioId, 
            p.DataPagamento, 
            p.Valor, 
            p.StatusId, 
            f.Cpf, 
            f.Nome,
			f.Sobrenome,
			f.ChavePix
        FROM 
            Pagamentos p WITH (NOLOCK)
        JOIN 
            Funcionarios f WITH (NOLOCK) ON p.FuncionarioId = f.Id
        WHERE 
            MONTH(p.DataPagamento) = @Mes AND YEAR(p.DataPagamento) = @Ano
        ORDER BY 
            p.DataPagamento DESC
        OFFSET @Offset ROWS
        FETCH NEXT @Limit ROWS ONLY
    `

	rows, err := s.db.Query(query,
		sql.Named("Mes", mes),
		sql.Named("Ano", anoAtual),
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar pagamentos por mês: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Pagamento
		if err := rows.Scan(&p.Id, &p.FuncionarioId, &p.DataPagamento, &p.Valor, &p.StatusId, &p.CpfFuncionario, &p.NomeFuncionario, &p.SobrenomeFuncionario, &p.ChavePixFuncionario); err != nil {
			return nil, fmt.Errorf("erro ao escanear pagamento: %w", err)
		}
		pagamentos = append(pagamentos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return pagamentos, nil
}
