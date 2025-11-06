package ponto

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	database "lapasta/database"
	"lapasta/internal/models"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"time"
)

type PontoRepository interface {
	ListarPontos(page int) ([]models.Ponto, error)
	ListarPontosPorId(idFuncionario int, page int) ([]models.Ponto, error)
	ListarPontosPorIdEDia(idFuncionario int, dia string) (models.Ponto, error)
	ListarPontosPorData(startDate, endDate time.Time, page int) ([]models.Ponto, error)
	ListarPontosPorDataId(IdFuncionario int, startDate time.Time, endDate time.Time, page int) ([]models.Ponto, error)
	RegistrarEntrada(idFuncionario int) error
	RegistrarSaidaAlmoco(idFuncionario int) error
	RegistrarRetornoAlmoco(idFuncionario int) error
	RegistrarSaida(idFuncionario int) error
	GerarRelatorioMensal(mes int, ano int, emailAdmin string) (string, error)
}

type pontoRepository struct {
	db *sql.DB
}

func NovoPontoRepository(conn *database.SQLStr) PontoRepository {
	return &pontoRepository{
		db: conn.DB(),
	}
}

func (s *pontoRepository) RegistrarEntrada(idFuncionario int) error {
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

func (s *pontoRepository) RegistrarSaidaAlmoco(idFuncionario int) error {
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

func (s *pontoRepository) RegistrarRetornoAlmoco(idFuncionario int) error {
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

func (s *pontoRepository) RegistrarSaida(idFuncionario int) error {
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

func (s *pontoRepository) MarcarAusente(idFuncionario int) error {
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
func (s *pontoRepository) ListarPontos(page int) ([]models.Ponto, error) {
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
func (s *pontoRepository) ListarPontosPorId(idFuncionario, page int) ([]models.Ponto, error) {
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
func (s *pontoRepository) ListarPontosPorIdEDia(idFuncionario int, dia string) (models.Ponto, error) {
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
func (s *pontoRepository) ListarPontosPorData(startDate time.Time, endDate time.Time, page int) ([]models.Ponto, error) {
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

func (s *pontoRepository) ListarPontosPorDataId(idFuncionario int, startDate time.Time, endDate time.Time, page int) ([]models.Ponto, error) {
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

func (s *pontoRepository) GerarRelatorioMensal(mes int, ano int, emailAdmin string) (string, error) {
	if mes == 0 {
		mes = int(time.Now().Month())
		ano = time.Now().Year()
	}

	query := `
		SELECT 
			p.IdFuncionario, f.Nome, p.Dia, p.HManha, p.HAlmocoSaida, p.HAlmocoRetorno, p.HNoite
		FROM Ponto p
		JOIN Funcionarios f ON p.IdFuncionario = f.Id
		WHERE MONTH(p.Dia) = @Mes AND YEAR(p.Dia) = @Ano
		ORDER BY p.IdFuncionario, p.Dia
	`

	rows, err := s.db.Query(query, sql.Named("Mes", mes), sql.Named("Ano", ano))
	if err != nil {
		return "", fmt.Errorf("erro ao consultar pontos do mês: %w", err)
	}
	defer rows.Close()

	relatorio := make(map[int][]models.Ponto)
	for rows.Next() {
		var ponto models.Ponto
		if err := rows.Scan(&ponto.IdFuncionario, &ponto.NomeFuncionario, &ponto.Dia, &ponto.HManha, &ponto.HAlmocoSaida, &ponto.HAlmocoRetorno, &ponto.HNoite); err != nil {
			return "", fmt.Errorf("erro ao escanear ponto: %w", err)
		}
		relatorio[ponto.IdFuncionario] = append(relatorio[ponto.IdFuncionario], ponto)
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Comma = ';'

	header := []string{"IdFuncionario", "Nome", "Dia", "HManha", "HAlmocoSaida", "HAlmocoRetorno", "HNoite"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("erro ao escrever cabeçalho: %w", err)
	}

	formatarHorario := func(hora sql.NullString) string {
		if !hora.Valid || hora.String == "0001-01-01T00:00:00Z" || hora.String == "" {
			return " - "
		}
		t, err := time.Parse("0001-01-01T15:04:05Z", hora.String)
		if err != nil {
			return " - "
		}
		return t.Format("15:04:05")
	}

	var ultimoIdFuncionario int
	for _, pontos := range relatorio {
		for i, ponto := range pontos {
			if i == 0 && ultimoIdFuncionario != 0 {
				if err := writer.Write([]string{}); err != nil {
					return "", fmt.Errorf("erro ao escrever linha em branco no buffer: %w", err)
				}
			}

			var dataFormatada string
			if t, err := time.Parse("2006-01-02T15:04:05Z", ponto.Dia); err == nil {
				dataFormatada = t.Format("02/01/2006")
			} else {
				dataFormatada = "Null"
			}

			linha := []string{
				fmt.Sprintf("%d", ponto.IdFuncionario),
				ponto.NomeFuncionario,
				dataFormatada,
				formatarHorario(ponto.HManha),
				formatarHorario(ponto.HAlmocoSaida),
				formatarHorario(ponto.HAlmocoRetorno),
				formatarHorario(ponto.HNoite),
			}

			if err := writer.Write(linha); err != nil {
				return "", fmt.Errorf("erro ao escrever linha no buffer: %w", err)
			}

			ultimoIdFuncionario = ponto.IdFuncionario
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("erro ao finalizar escrita no buffer: %w", err)
	}

	if err := s.enviarEmailComAnexo(emailAdmin, buf.Bytes(), mes, ano); err != nil {
		return "", fmt.Errorf("erro ao enviar e-mail: %w", err)
	}

	return "Relatório enviado com sucesso!", nil
}

func (s *pontoRepository) enviarEmailComAnexo(emailAdmin string, fileContent []byte, mes, ano int) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := "email@gmail.com"
	senderPassword := "senha"
	subject := fmt.Sprintf("Relatório Mensal de Pontos - %02d/%d", mes, ano)
	body := fmt.Sprintf("Olá, em anexo está o relatório mensal de pontos do mês %02d do ano %d.", mes, ano)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	textPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{"text/plain; charset=UTF-8"},
		"Content-Transfer-Encoding": []string{"quoted-printable"},
	})
	if err != nil {
		return fmt.Errorf("erro ao criar a parte do corpo do e-mail: %w", err)
	}
	encoder := quotedprintable.NewWriter(textPart)
	_, err = encoder.Write([]byte(body))
	if err != nil {
		return fmt.Errorf("erro ao escrever o corpo do e-mail: %w", err)
	}
	encoder.Close()

	filePart, err := writer.CreateFormFile("attachment", "relatorio_pontos.csv")
	if err != nil {
		return fmt.Errorf("erro ao criar a parte do anexo: %w", err)
	}
	_, err = filePart.Write(fileContent)
	if err != nil {
		return fmt.Errorf("erro ao escrever o anexo: %w", err)
	}
	writer.Close()

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	to := []string{emailAdmin}
	msg := []byte("From: " + senderEmail + "\r\n" +
		"To: " + emailAdmin + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: multipart/mixed; boundary=\"" + writer.Boundary() + "\"\r\n" +
		"\r\n" + buf.String())

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, to, msg)
	if err != nil {
		return fmt.Errorf("erro ao enviar e-mail: %w", err)
	}

	return nil
}
