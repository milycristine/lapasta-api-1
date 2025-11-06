package boleto

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	database "lapasta/database"
	"lapasta/internal/models"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type BoletoRepository interface {
	CriarBoletoRecebido(boleto *models.Boleto) error
	ListarBoletosPorFornecedor(fornecedorId int) ([]models.Boleto, error)
	ListarBoletosPorRecebimento(recebimentoId int) ([]models.Boleto, error)
	ListarBoletosDoDia(data string) ([]models.Boleto, error)
	PagarBoleto(codigoBarras string) error
	ListarBoletosVencidos() ([]models.Boleto, error)
	ListarBoletosPagos() ([]models.Boleto, error)
	ListarBoletosPendentes() ([]models.Boleto, error)
	AtualizarStatusBoleto(id int, statusId int) error
	GerarEEnviarRelatorioBoletos(emailAdmin string, ano, mes int, fornecedorID *int, statusIDs []int) error
	TotalBoletosDia(data string) (float64, error)
	TotalBoletosAtrasados() (float64, error)
	TotalBoletosPendentesMesAtual() (float64, error)
	TotalBoletosPagosMesAtual() (float64, error)
	FiltrarBoletosPagosPorData(inicioData, fimData string) ([]models.Boleto, error)
}

type boletoRepository struct {
	db *sql.DB
}

func NovoBoletoRepository(conn *database.SQLStr) BoletoRepository {
	return &boletoRepository{
		db: conn.DB(),
	}
}

func (r *boletoRepository) CriarBoletoRecebido(b *models.Boleto) error {
	var receb models.Recebimento

	queryReceb := `
		SELECT Id, Dia, Valor, Vencimento
		FROM Recebimento WITH (NOLOCK)
		WHERE Id = @RecebimentoId
	`
	err := r.db.QueryRow(queryReceb,
		sql.Named("RecebimentoId", b.RecebimentoId),
	).Scan(&receb.Id, &receb.Dia, &receb.Valor, &receb.Vencimento)
	if err != nil {
		log.Printf("Erro ao buscar recebimento: %v", err)
		return fmt.Errorf("erro ao buscar recebimento: %w", err)
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

func (r *boletoRepository) ListarBoletosPorFornecedor(fornecedorId int) ([]models.Boleto, error) {
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

func (r *boletoRepository) ListarBoletosPorRecebimento(recebimentoId int) ([]models.Boleto, error) {
	query := `
		SELECT b.Id, b.RecebimentoId, b.CodigoBarras, b.DataCadastro, b.DataVencimento, b.Valor, b.StatusId, b.DataPagamento
		FROM BoletosRecebidos b
		WHERE b.RecebimentoId = @RecebimentoId
	`
	rows, err := r.db.Query(query, sql.Named("RecebimentoId", recebimentoId))
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos por recebimento: %w", err)
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

func (r *boletoRepository) ListarBoletosDoDia(data string) ([]models.Boleto, error) {
	query := `
	SELECT b.Id, b.RecebimentoId, b.CodigoBarras, b.DataVencimento, b.Valor, b.StatusId
	FROM BoletosRecebidos b
	INNER JOIN Recebimento r ON r.Id = b.RecebimentoId
	WHERE CONVERT(date, b.DataVencimento) = @DataHoje
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

		err := rows.Scan(&boleto.Id, &boleto.RecebimentoId, &boleto.CodigoBarras, &dataVencimentoStr, &boleto.Valor, &statusId)
		if err != nil {
			return nil, err
		}

		dataVencimento, err := time.Parse(time.RFC3339, dataVencimentoStr)
		if err != nil {
			dataVencimento, err = time.Parse("2006-01-02", dataVencimentoStr)
			if err != nil {
				return nil, fmt.Errorf("erro ao converter data de vencimento: %w", err)
			}
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

func (r *boletoRepository) PagarBoleto(codigoBarras string) error {
	query := `	
		UPDATE BoletosRecebidos
		SET StatusId = 2, DataPagamento = CURRENT_TIMESTAMP
		WHERE CodigoBarras = @CodigoBarras 
		  AND StatusId <> 2 
		  AND DataVencimento = CONVERT(date, GETDATE())
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
		return fmt.Errorf("boleto já pago, vencimento diferente de hoje ou não encontrado")
	}

	_, err = r.db.Exec(`
    INSERT INTO PagamentosLog (CodigoBarras, DataHora) 
    VALUES (@CodigoBarras, CURRENT_TIMESTAMP)
`, sql.Named("CodigoBarras", codigoBarras))
	if err != nil {
		log.Printf("Erro ao registrar log de pagamento do boleto %s: %v", codigoBarras, err)
	}

	return nil
}

func (r *boletoRepository) ListarBoletosPagos() ([]models.Boleto, error) {
	query := `
		SELECT 
			b.Id,
			b.Valor, 
			b.DataVencimento, 
			b.DataPagamento, 
			b.CodigoBarras, 
			b.StatusId,
			r.Dia
		FROM 
			BoletosRecebidos b
		INNER JOIN 
			Recebimento r ON r.Id = b.RecebimentoId
		WHERE 
			b.StatusId = 2
		ORDER BY 
			b.DataPagamento DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos pagos: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		var dataVencimentoStr, dataPagamentoStr, diaStr sql.NullString

		err := rows.Scan(
			&b.Id,
			&b.Valor,
			&dataVencimentoStr,
			&dataPagamentoStr,
			&b.CodigoBarras,
			&b.StatusId,
			&diaStr,
		)
		if err != nil {
			return nil, err
		}

		if dataVencimentoStr.Valid {
			b.DataVencimento = dataVencimentoStr.String
		}
		if dataPagamentoStr.Valid {
			b.DataPagamento = &dataPagamentoStr.String
		}
		if diaStr.Valid {
			b.DataCadastro = diaStr.String
		}

		boletos = append(boletos, b)
	}
	return boletos, nil
}

func (r *boletoRepository) ListarBoletosVencidos() ([]models.Boleto, error) {
	query := `
		SELECT 
			b.Id,
			b.Valor, 
			b.DataVencimento, 
			b.CodigoBarras, 
			b.StatusId,
			r.Dia
		FROM 
			BoletosRecebidos b
		INNER JOIN 
			Recebimento r ON r.Id = b.RecebimentoId
		WHERE 
			b.StatusId = 3
		ORDER BY 
			b.DataVencimento
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos vencidos: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		var dataVencimentoStr, diaStr string

		err := rows.Scan(&b.Id, &b.Valor, &dataVencimentoStr, &b.CodigoBarras, &b.StatusId, &diaStr)
		if err != nil {
			return nil, err
		}

		b.DataVencimento = dataVencimentoStr
		b.DataCadastro = diaStr
		boletos = append(boletos, b)
	}
	return boletos, nil
}

func (r *boletoRepository) ListarBoletosPendentes() ([]models.Boleto, error) {
	query := `
		SELECT 
			b.Id,
			b.Valor, 
			b.DataVencimento, 
			b.CodigoBarras, 
			b.StatusId,
			r.Dia
		FROM 
			BoletosRecebidos b
		INNER JOIN 
			Recebimento r ON r.Id = b.RecebimentoId
		WHERE 
			b.StatusId = 1
		ORDER BY 
			b.DataVencimento
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos pendentes: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		var dataVencimentoStr, diaStr string

		err := rows.Scan(&b.Id, &b.Valor, &dataVencimentoStr, &b.CodigoBarras, &b.StatusId, &diaStr)
		if err != nil {
			return nil, err
		}

		b.DataVencimento = dataVencimentoStr
		b.DataCadastro = diaStr
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

func (s *boletoRepository) AtualizarStatusBoleto(id int, statusId int) error {
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

func (r *boletoRepository) TotalBoletosDia(data string) (float64, error) {
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

func (r *boletoRepository) TotalBoletosAtrasados() (float64, error) {
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

func (r *boletoRepository) TotalBoletosPendentesMesAtual() (float64, error) {
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
func (r *boletoRepository) TotalBoletosPagosMesAtual() (float64, error) {
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

func (r *boletoRepository) FiltrarBoletosPagosPorData(inicioData, fimData string) ([]models.Boleto, error) {
	layout := "2006-01-02"
	start, _ := time.Parse(layout, inicioData)
	end, _ := time.Parse(layout, fimData)

	query := `
	SELECT b.Id, b.RecebimentoId, b.CodigoBarras, b.DataCadastro, b.DataVencimento, b.Valor, b.StatusId, f.Nome AS FornecedorNome
	FROM BoletosRecebidos b
	INNER JOIN Recebimento r ON r.Id = b.RecebimentoId
	INNER JOIN PedidoFornecedor pf ON pf.Id = r.IdPedidoFornecedor
	INNER JOIN Fornecedores f ON f.Id = pf.FornecedorId
	WHERE b.StatusId = 2 AND CONVERT(date, b.DataPagamento) BETWEEN @InicioData AND @FimData
	ORDER BY b.DataPagamento DESC
	`

	rows, err := r.db.Query(query,
		sql.Named("InicioData", start.Format("2006-01-02")),
		sql.Named("FimData", end.Format("2006-01-02")),
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos pagos por data: %w", err)
	}
	defer rows.Close()

	var boletos []models.Boleto
	for rows.Next() {
		var b models.Boleto
		if err := rows.Scan(
			&b.Id, &b.RecebimentoId, &b.CodigoBarras, &b.DataCadastro,
			&b.DataVencimento, &b.Valor, &b.StatusId, &b.FornecedorNome,
		); err != nil {
			return nil, fmt.Errorf("erro ao escanear boleto pago: %w", err)
		}
		boletos = append(boletos, b)
	}

	return boletos, nil
}

func (r *boletoRepository) ListarBoletosComFiltro(ano, mes int, fornecedorID *int, statusIDs []int) ([]models.BoletoRelatorio, error) {
	query := `
		SELECT 
			f.Nome,
			b.DataVencimento,
			b.Valor,
			b.CodigoBarras,
			b.StatusId
		FROM 
			BoletosRecebidos b WITH (NOLOCK)
		INNER JOIN 
			Recebimento r WITH (NOLOCK) ON r.Id = b.RecebimentoId
		INNER JOIN 
			PedidoFornecedor p WITH (NOLOCK) ON p.Id = r.IdPedidoFornecedor
		INNER JOIN 
			Fornecedores f WITH (NOLOCK) ON f.Id = p.FornecedorId
		WHERE 
			YEAR(b.DataVencimento) = @Ano AND MONTH(b.DataVencimento) = @Mes
	`

	params := []interface{}{
		sql.Named("Ano", ano),
		sql.Named("Mes", mes),
	}

	if fornecedorID != nil {
		query += " AND f.Id = @FornecedorId"
		params = append(params, sql.Named("FornecedorId", *fornecedorID))
	}

	if len(statusIDs) > 0 {
		query += " AND b.StatusId IN ("
		for i := range statusIDs {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("@StatusId%d", i)
		}
		query += ")"
		for i, s := range statusIDs {
			params = append(params, sql.Named(fmt.Sprintf("StatusId%d", i), s))
		}
	}

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar boletos com filtro: %w", err)
	}
	defer rows.Close()

	var boletos []models.BoletoRelatorio
	for rows.Next() {
		var b models.BoletoRelatorio
		err := rows.Scan(&b.FornecedorNome, &b.DataVencimento, &b.Valor, &b.CodigoBarras, &b.StatusId)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear boleto: %w", err)
		}
		boletos = append(boletos, b)
	}

	return boletos, nil
}

func gerarRelatorioPDF(boletos []models.BoletoRelatorio, statusIDs []int) ([]byte, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")

	pdf.AddUTF8Font("OpenSans", "", "fonts/OpenSans-Regular.ttf")
	pdf.AddUTF8Font("OpenSans", "B", "fonts/OpenSans-Bold.ttf")

	pdf.SetTitle("Relatório de Boletos", false)
	pdf.AddPage()

	pdf.SetFont("OpenSans", "B", 14)
	pdf.Cell(0, 10, "Relatório de Boletos")
	pdf.Ln(12)

	pdf.SetFont("OpenSans", "B", 12)
	headers := []string{"Fornecedor", "Vencimento", "Valor", "Código de Barras", "Status"}
	widths := []float64{50, 30, 30, 100, 30}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 10, h, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("OpenSans", "", 10)

	var total float64

	for _, b := range boletos {
		pdf.CellFormat(widths[0], 8, b.FornecedorNome, "1", 0, "", false, 0, "")
		pdf.CellFormat(widths[1], 8, b.DataVencimento.Format("02/01/2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[2], 8, fmt.Sprintf("R$ %.2f", b.Valor), "1", 0, "R", false, 0, "")
		pdf.CellFormat(widths[3], 8, b.CodigoBarras, "1", 0, "", false, 0, "")

		statusTexto := "Desconhecido"
		switch b.StatusId {
		case 1:
			statusTexto = "Pendente"
		case 2:
			statusTexto = "Pago"
		case 3:
			statusTexto = "Atrasado"
		}

		pdf.CellFormat(widths[4], 8, statusTexto, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)

		total += b.Valor
	}

	statusLabel := "Geral"
	if len(statusIDs) == 1 {
		switch statusIDs[0] {
		case 1:
			statusLabel = "Pendentes"
		case 2:
			statusLabel = "Pagos"
		case 3:
			statusLabel = "Atrasados"
		}
	}

	pdf.SetFont("OpenSans", "B", 11)
	pdf.CellFormat(widths[0]+widths[1], 10, "TOTAL "+statusLabel, "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[2], 10, fmt.Sprintf("R$ %.2f", total), "1", 0, "R", false, 0, "")
	pdf.CellFormat(widths[3]+widths[4], 10, "", "1", 0, "", false, 0, "")
	pdf.Ln(-1)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("erro ao gerar PDF: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *boletoRepository) enviarEmailComAnexoPDF(emailAdmin string, fileContent []byte, mes, ano int) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	senderEmail := "email"
	senderPassword := "senha"
	subject := fmt.Sprintf("Relatório Mensal de Boletos - %02d/%d", mes, ano)
	body := fmt.Sprintf("Olá, em anexo está o relatório mensal de boletos do mês %02d do ano %d.", mes, ano)

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

	filePart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{"application/pdf"},
		"Content-Disposition":       []string{`attachment; filename="relatorio_boletos.pdf"`},
		"Content-Transfer-Encoding": []string{"base64"},
	})
	if err != nil {
		return fmt.Errorf("erro ao criar a parte do anexo: %w", err)
	}

	b64Writer := base64.NewEncoder(base64.StdEncoding, filePart)
	_, err = b64Writer.Write(fileContent)
	if err != nil {
		return fmt.Errorf("erro ao escrever o anexo: %w", err)
	}
	b64Writer.Close()

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
func (s *boletoRepository) GerarEEnviarRelatorioBoletos(emailAdmin string, ano, mes int, fornecedorID *int, statusIDs []int) error {
	if ano == 0 || mes == 0 {
		now := time.Now()
		if ano == 0 {
			ano = now.Year()
		}
		if mes == 0 {
			mes = int(now.Month())
		}
	}

	boletos, err := s.ListarBoletosComFiltro(ano, mes, fornecedorID, statusIDs)
	if err != nil {
		return fmt.Errorf("erro ao buscar boletos: %w", err)
	}

	if len(boletos) == 0 {
		return fmt.Errorf("nenhum boleto encontrado para os filtros informados")
	}

	pdfBytes, err := gerarRelatorioPDF(boletos, statusIDs)
	if err != nil {
		return fmt.Errorf("erro ao gerar PDF: %w", err)
	}

	err = s.enviarEmailComAnexoPDF(emailAdmin, pdfBytes, mes, ano)
	if err != nil {
		return fmt.Errorf("erro ao enviar e-mail: %w", err)
	}

	return nil
}
