package sql

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"strings"
	"time"
)

func (s *SQLStr) CriarMotorista(m *models.Motorista) error {
	hashedPassword := sha256.Sum256([]byte(m.Senha))

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	insertMotorista := `
		INSERT INTO Motoristas (Nome, Telefone, CPF, ChavePix, ValorDiaria, Ativo, Email, Senha, CnhFrenteUrl, CnhVersoUrl, KmDiaria)
		VALUES (@Nome, @Telefone, @CPF, @ChavePix, @ValorDiaria, @Ativo, @Email, @Senha, @CnhFrenteUrl, @CnhVersoUrl, @KmDiaria)
	`

	_, err = tx.Exec(insertMotorista,
		sql.Named("Nome", m.Nome),
		sql.Named("Telefone", m.Telefone),
		sql.Named("CPF", m.CPF),
		sql.Named("ChavePix", m.ChavePix),
		sql.Named("ValorDiaria", m.ValorDiaria),
		sql.Named("Ativo", m.Ativo),
		sql.Named("Email", m.Email),
		sql.Named("Senha", hashedPassword[:]),
		sql.Named("CnhFrenteUrl", m.CnhFrenteUrl),
		sql.Named("CnhVersoUrl", m.CnhVersoUrl),
		sql.Named("KmDiaria", m.KmDiaria),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir motorista: %w", err)
	}

	authInsert := `INSERT INTO AUTH (username, password) VALUES (@Username, @Password)`
	_, err = tx.Exec(authInsert,
		sql.Named("Username", m.Email),
		sql.Named("Password", hashedPassword[:]),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir credenciais de motorista: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao finalizar transação: %w", err)
	}

	return nil
}

func (s *SQLStr) AtualizarStatusMotorista(id int, ativo bool) error {
	query := `UPDATE Motoristas SET Ativo = @Ativo WHERE Id = @Id`
	_, err := s.db.Exec(query,
		sql.Named("Ativo", ativo),
		sql.Named("Id", id),
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar status do motorista: %w", err)
	}
	return nil
}

func (s *SQLStr) ListarMotoristas() ([]models.Motorista, error) {
	query := `SELECT Id, Nome, Telefone, CPF, ChavePix, ValorDiaria, Ativo FROM Motoristas`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar motoristas: %w", err)
	}
	defer rows.Close()

	var motoristas []models.Motorista
	for rows.Next() {
		var m models.Motorista
		if err := rows.Scan(&m.Id, &m.Nome, &m.Telefone, &m.CPF, &m.ChavePix, &m.ValorDiaria, &m.Ativo); err != nil {
			return nil, fmt.Errorf("erro ao escanear motorista: %w", err)
		}
		motoristas = append(motoristas, m)
	}
	return motoristas, nil
}

func (s *SQLStr) EditarMotorista(m *models.Motorista) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	var (
		setClauses []string
		args       []any
	)

	
	if m.Nome != "" {
		setClauses = append(setClauses, "Nome = @Nome")
		args = append(args, sql.Named("Nome", m.Nome))
	}
	if m.Telefone != "" {
		setClauses = append(setClauses, "Telefone = @Telefone")
		args = append(args, sql.Named("Telefone", m.Telefone))
	}
	if m.ChavePix != "" {
		setClauses = append(setClauses, "ChavePix = @ChavePix")
		args = append(args, sql.Named("ChavePix", m.ChavePix))
	}
	if m.Email != "" {
		setClauses = append(setClauses, "Email = @Email")
		args = append(args, sql.Named("Email", m.Email))
	}

	setClauses = append(setClauses, "ValorDiaria = @ValorDiaria")
	args = append(args, sql.Named("ValorDiaria", m.ValorDiaria))

	setClauses = append(setClauses, "KmDiaria = @KmDiaria")
	args = append(args, sql.Named("KmDiaria", m.KmDiaria))

	setClauses = append(setClauses, "Ativo = @Ativo")
	args = append(args, sql.Named("Ativo", m.Ativo))

	if len(setClauses) == 0 {
		tx.Rollback()
		return fmt.Errorf("nenhum campo para atualizar")
	}

	args = append(args, sql.Named("Cpf", m.CPF))
	query := fmt.Sprintf("UPDATE Motoristas SET %s WHERE CPF = @Cpf", strings.Join(setClauses, ", "))

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao atualizar motorista: %w", err)
	}

	if m.Senha != "" {
		hashedPassword := sha256.Sum256([]byte(m.Senha))
		authUpdate := `
			UPDATE AUTH SET 
				password = @Password
			WHERE username = @Username
		`
		_, err = tx.Exec(authUpdate,
			sql.Named("Username", m.Email),
			sql.Named("Password", hashedPassword[:]),
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("erro ao atualizar senha: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao finalizar transação: %w", err)
	}

	return nil
}


func (s *SQLStr) BuscarMotoristaPorID(id int) (*models.Motorista, error) {
	query := `
		SELECT Id, Nome, Telefone, CPF, ChavePix, ValorDiaria, Ativo, Email, 
		       CnhFrenteUrl, CnhVersoUrl, KmDiaria
		FROM Motoristas
		WHERE Id = @Id
	`

	row := s.db.QueryRow(query, sql.Named("Id", id))

	var m models.Motorista
	err := row.Scan(
		&m.Id,
		&m.Nome,
		&m.Telefone,
		&m.CPF,
		&m.ChavePix,
		&m.ValorDiaria,
		&m.Ativo,
		&m.Email,
		&m.CnhFrenteUrl,
		&m.CnhVersoUrl,
		&m.KmDiaria,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar motorista por ID: %w", err)
	}

	return &m, nil
}

func (s *SQLStr) CriarEmissaoNota(n *models.EmissaoNota) error {
	query := `
		INSERT INTO EmissaoNotas 
		(NumeroNota, Valor, DataEmissao, DataAtribuicao,  Descricao, MotoristaId, IdStatusLancamento)
		VALUES 
		(@NumeroNota, @Valor, @DataEmissao, @DataAtribuicao,  @Descricao, @MotoristaId, @IdStatusLancamento)
	`

	_, err := s.db.Exec(query,
		sql.Named("NumeroNota", n.NumeroNota),
		sql.Named("Valor", n.Valor),
		sql.Named("DataEmissao", n.DataEmissao),
		sql.Named("DataAtribuicao", time.Now().Format("2006-01-02")),
		sql.Named("Descricao", n.Descricao),
		sql.Named("MotoristaId", n.MotoristaId),
		sql.Named("IdStatusLancamento", n.IdStatusLancamento),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir emissão de nota: %w", err)
	}
	return nil
}

func (s *SQLStr) ListarEmissaoNotas(page int) ([]models.EmissaoNota, error) {
	limit := 10
	offset := (page - 1) * limit

	query := `
		SELECT 
			e.Id, 
			e.NumeroNota, 
			e.Valor, 
			e.DataEmissao, 
			e.DataAtribuicao, 
			e.Descricao,
			e.MotoristaId,
			m.Nome AS MotoristaNome,
			e.IdStatusLancamento
		FROM 
			EmissaoNotas e
		JOIN 
			Motoristas m ON e.MotoristaId = m.Id
		ORDER BY 
			e.DataEmissao DESC
		OFFSET @Offset ROWS 
		FETCH NEXT @Limit ROWS ONLY
	`

	rows, err := s.db.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar emissão de notas: %w", err)
	}
	defer rows.Close()

	var notas []models.EmissaoNota

	for rows.Next() {
		var n models.EmissaoNota
		err := rows.Scan(
			&n.Id,
			&n.NumeroNota,
			&n.Valor,
			&n.DataEmissao,
			&n.DataAtribuicao,
			&n.Descricao,
			&n.MotoristaId,
			&n.MotoristaNome,
			&n.IdStatusLancamento,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear emissão de nota: %w", err)
		}
		notas = append(notas, n)
	}

	return notas, nil
}

func (s *SQLStr) FiltrarDataEmissaoNota(inicioData time.Time, fimData time.Time) ([]models.EmissaoNota, error) {
	var notas []models.EmissaoNota

	query := `
		SELECT 
			e.Id, 
			e.NumeroNota, 
			e.Valor, 
			e.DataEmissao, 
			e.DataAtribuicao, 
			e.Descricao,
			e.MotoristaId,
			m.Nome AS MotoristaNome,
			e.IdStatusLancamento
		FROM 
			EmissaoNotas e
		JOIN 
			Motoristas m ON e.MotoristaId = m.Id
		WHERE 
			e.DataAtribuicao BETWEEN @InicioData AND @FimData
		ORDER BY 
			e.DataAtribuicao DESC
	`

	rows, err := s.db.Query(query, sql.Named("InicioData", inicioData), sql.Named("FimData", fimData))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar emissão de notas: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n models.EmissaoNota
		err := rows.Scan(
			&n.Id,
			&n.NumeroNota,
			&n.Valor,
			&n.DataEmissao,
			&n.DataAtribuicao,
			&n.Descricao,
			&n.MotoristaId,
			&n.MotoristaNome,
			&n.IdStatusLancamento,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear emissão de nota: %w", err)
		}
		notas = append(notas, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return notas, nil
}

func (s *SQLStr) BuscarEmissaoNotas(valor string) ([]models.EmissaoNota, error) {
	query := `
		SELECT 
			e.Id, 
			e.NumeroNota, 
			e.Valor, 
			e.DataEmissao, 
			e.DataAtribuicao, 
			e.Descricao,
			e.MotoristaId,
			m.Nome AS MotoristaNome,
			e.IdStatusLancamento
		FROM 
			EmissaoNotas e
		JOIN 
			Motoristas m ON e.MotoristaId = m.Id
		WHERE 
			e.NumeroNota LIKE @Valor OR 
			e.Descricao LIKE @Valor OR 
			m.Nome LIKE @Valor
	`

	rows, err := s.db.Query(query, sql.Named("Valor", "%"+valor+"%"))
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar emissão de notas: %w", err)
	}
	defer rows.Close()

	var notas []models.EmissaoNota

	for rows.Next() {
		var n models.EmissaoNota
		err := rows.Scan(
			&n.Id,
			&n.NumeroNota,
			&n.Valor,
			&n.DataEmissao,
			&n.DataAtribuicao,
			&n.Descricao,
			&n.MotoristaId,
			&n.MotoristaNome,
			&n.IdStatusLancamento,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear emissão de nota: %w", err)
		}
		notas = append(notas, n)
	}

	if len(notas) == 0 {
		return nil, fmt.Errorf("nenhuma emissão de nota encontrada com '%s'", valor)
	}

	return notas, nil
}

func (s *SQLStr) CriarNotasMotorista(nm *models.NotasMotorista) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	query := `
		INSERT INTO NotasMotorista (IdMotorista, IdNota, IdStatusLancamento, Url)
		VALUES (@IdMotorista, @IdNota, @StatusId, @Url)
	`
	_, err = tx.Exec(query,
		sql.Named("IdMotorista", nm.IdMotorista),
		sql.Named("IdNota", nm.IdNota),
		sql.Named("StatusId", 1),
		sql.Named("Url", nm.Url),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao inserir notasMotorista: %w", err)
	}

	updateQuery := `
		UPDATE EmissaoNotas
		SET IdStatusLancamento = 1
		WHERE Id = @IdNota
	`
	_, err = tx.Exec(updateQuery, sql.Named("IdNota", nm.IdNota))
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("erro ao atualizar status da EmissaoNotas: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}
func (s *SQLStr) ListarNotasMotoristaPorMotorista(idMotorista int) ([]models.NotasMotorista, error) {
	query := `
		SELECT
			nm.Id,
			nm.IdMotorista,
			nm.IdNota,
			nm.DataAtribuicao,
			nm.DataLancamento,
			nm.IdStatusLancamento,
			nm.Url,
			en.NumeroNota,
			en.Valor,
			en.DataEmissao,
			en.Descricao,
			m.Nome as NomeMotorista
		FROM NotasMotorista nm
		INNER JOIN Motoristas m ON nm.IdMotorista = m.Id
		INNER JOIN EmissaoNotas en ON nm.IdNota = en.Id
		WHERE nm.IdMotorista = @IdMotorista
		  AND nm.IdStatusLancamento = 1
		ORDER BY nm.DataAtribuicao DESC
	`

	rows, err := s.db.Query(query, sql.Named("IdMotorista", idMotorista))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar notasMotorista: %w", err)
	}
	defer rows.Close()

	var resultado []models.NotasMotorista
	for rows.Next() {
		var nm models.NotasMotorista
		var dataAtribuicao, dataLancamento, dataEmissao sql.NullTime
		var valor sql.NullFloat64
		var numeroNota, descricao, nomeMotorista sql.NullString

		if err := rows.Scan(
			&nm.Id,
			&nm.IdMotorista,
			&nm.IdNota,
			&dataAtribuicao,
			&dataLancamento,
			&nm.IdStatusLancamento,
			&nm.Url,
			&numeroNota,
			&valor,
			&dataEmissao,
			&descricao,
			&nomeMotorista,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear notasMotorista: %w", err)
		}

		if dataAtribuicao.Valid {
			nm.DataAtribuicao = dataAtribuicao.Time.Format("2006-01-02 15:04:05")
		} else {
			nm.DataAtribuicao = ""
		}

		if numeroNota.Valid {
			nm.NumeroNota = numeroNota.String
		} else {
			nm.NumeroNota = ""
		}

		if valor.Valid {
			nm.Valor = valor.Float64
		} else {
			nm.Valor = 0.0
		}

		if dataEmissao.Valid {
			nm.DataEmissao = dataEmissao.Time.Format("2006-01-02 15:04:05")
		} else {
			nm.DataEmissao = ""
		}

		if descricao.Valid {
			nm.Descricao = descricao.String
		} else {
			nm.Descricao = ""
		}

		if nomeMotorista.Valid {
			nm.NomeMotorista = nomeMotorista.String
		} else {
			nm.NomeMotorista = ""
		}

		resultado = append(resultado, nm)
	}

	return resultado, nil
}

func (s *SQLStr) AtualizarStatusLancamentoNotaMotorista(id int, statusId int, dataLancamento *time.Time) error {
	dataLancamentoFinal := time.Now()
	if dataLancamento != nil {
		dataLancamentoFinal = *dataLancamento
	}

	_, err := s.db.Exec(`
		UPDATE NotasMotorista
		SET IdStatusLancamento = @StatusId, DataLancamento = @DataLancamento
		WHERE Id = @Id
	`,
		sql.Named("StatusId", statusId),
		sql.Named("DataLancamento", dataLancamentoFinal),
		sql.Named("Id", id),
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar NotasMotorista: %w", err)
	}

	var idNota int
	err = s.db.QueryRow(`
		SELECT IdNota FROM NotasMotorista WHERE Id = @Id
	`, sql.Named("Id", id)).Scan(&idNota)
	if err != nil {
		return fmt.Errorf("erro ao buscar IdNota relacionado: %w", err)
	}

	_, err = s.db.Exec(`
		UPDATE EmissaoNotas
		SET IdStatusLancamento = @StatusId
		WHERE Id = @IdNota
	`,
		sql.Named("StatusId", statusId),
		sql.Named("IdNota", idNota),
	)
	if err != nil {
		return fmt.Errorf("erro ao atualizar EmissaoNotas: %w", err)
	}

	return nil
}

const StatusLancado = 1

func (s *SQLStr) MotoristaLancouTodasAsNotas(idMotorista int, inicio, fim time.Time) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM (
			SELECT IdStatusLancamento
			FROM NotasMotorista
			WHERE IdMotorista = @IdMotorista
			  AND CAST(DataAtribuicao AS DATE) BETWEEN @Inicio AND @Fim

			UNION ALL

			SELECT IdStatusLancamento
			FROM EmissaoNotas
			WHERE MotoristaId = @IdMotorista
			  AND CAST(DataEmissao AS DATE) BETWEEN @Inicio AND @Fim
		) AS TodasNotas
		WHERE IdStatusLancamento IS NULL OR IdStatusLancamento != @StatusLancado
	`

	var count int
	err := s.db.QueryRow(query,
		sql.Named("IdMotorista", idMotorista),
		sql.Named("Inicio", inicio),
		sql.Named("Fim", fim),
		sql.Named("StatusLancado", StatusLancado),
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar notas do motorista: %w", err)
	}

	return count == 0, nil
}

func (r *SQLStr) BuscarMotoristaPorCPFouNome(valor string) ([]models.Motorista, error) {
	var motoristas []models.Motorista

	query := `
		SELECT Id, Nome, Telefone, CPF, ChavePix, ValorDiaria, Ativo
		FROM Motoristas WITH (NOLOCK)
		WHERE CPF = @Valor 
		   OR Nome COLLATE Latin1_General_CI_AI LIKE @LikeValor
	`

	rows, err := r.db.Query(query,
		sql.Named("Valor", valor),
		sql.Named("LikeValor", "%"+valor+"%"),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar motoristas: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Motorista
		if err := rows.Scan(&m.Id, &m.Nome, &m.Telefone, &m.CPF, &m.ChavePix, &m.ValorDiaria, &m.Ativo); err != nil {
			return nil, fmt.Errorf("erro ao escanear motorista: %w", err)
		}
		motoristas = append(motoristas, m)
	}

	if len(motoristas) == 0 {
		return nil, fmt.Errorf("nenhum motorista encontrado com '%s'", valor)
	}

	return motoristas, nil
}
func (s *SQLStr) ListarEmissaoNotasPorMotorista(idMotorista int) ([]models.EmissaoNota, error) {
	query := `
        SELECT 
            e.Id, 
            e.NumeroNota, 
            e.Valor, 
            e.DataEmissao, 
            e.DataAtribuicao, 
            e.Descricao,
            e.MotoristaId,
            m.Nome AS MotoristaNome,
            e.IdStatusLancamento
        FROM 
            EmissaoNotas e
        JOIN 
            Motoristas m ON e.MotoristaId = m.Id
        WHERE 
            e.MotoristaId = @IdMotorista
            AND (e.IdStatusLancamento IS NULL OR e.IdStatusLancamento != 1)
        ORDER BY 
            e.DataEmissao DESC
    `

	rows, err := s.db.Query(query, sql.Named("IdMotorista", idMotorista))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar emissão de notas pendentes: %w", err)
	}
	defer rows.Close()

	var notas []models.EmissaoNota

	for rows.Next() {
		var n models.EmissaoNota
		err := rows.Scan(
			&n.Id,
			&n.NumeroNota,
			&n.Valor,
			&n.DataEmissao,
			&n.DataAtribuicao,
			&n.Descricao,
			&n.MotoristaId,
			&n.MotoristaNome,
			&n.IdStatusLancamento,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear emissão de nota pendente: %w", err)
		}
		notas = append(notas, n)
	}

	return notas, nil
}

func (s *SQLStr) FiltrarNotasMotoristaPorData(
	idMotorista int,
	inicioData *time.Time,
	fimData *time.Time,
) ([]models.NotasMotorista, error) {
	var resultado []models.NotasMotorista

	query := `
		SELECT
			nm.Id,
			nm.IdMotorista,
			nm.IdNota,
			nm.DataAtribuicao,
			nm.DataLancamento,
			nm.IdStatusLancamento,
			nm.Url,
			en.NumeroNota,
			en.Valor,
			en.DataEmissao,
			en.DataAtribuicao,
			en.Descricao,
			m.Nome as NomeMotorista
		FROM NotasMotorista nm
		INNER JOIN Motoristas m ON nm.IdMotorista = m.Id
		INNER JOIN EmissaoNotas en ON nm.IdNota = en.Id
		WHERE nm.IdMotorista = @IdMotorista
		  AND nm.IdStatusLancamento = 1
	`

	if inicioData != nil && fimData != nil {
		query += " AND nm.DataAtribuicao BETWEEN @InicioData AND @FimData"
	}

	query += " ORDER BY nm.DataAtribuicao DESC"

	params := []any{
		sql.Named("IdMotorista", idMotorista),
	}
	if inicioData != nil && fimData != nil {
		params = append(params,
			sql.Named("InicioData", *inicioData),
			sql.Named("FimData", *fimData),
		)
	}

	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("erro ao filtrar notas motorista: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var nm models.NotasMotorista

		err := rows.Scan(
			&nm.Id,
			&nm.IdMotorista,
			&nm.IdNota,
			&nm.DataAtribuicao,
			&nm.DataLancamento,
			&nm.IdStatusLancamento,
			&nm.Url,
			&nm.NumeroNota,
			&nm.Valor,
			&nm.DataEmissao,
			&nm.DataAtribuicao,
			&nm.Descricao,
			&nm.NomeMotorista,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear notasMotorista: %w", err)
		}

		resultado = append(resultado, nm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return resultado, nil
}
