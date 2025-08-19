package sql

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"lapasta/internal/models"
	"mime/multipart"
	"mime/quotedprintable"
	"net/smtp"
	"net/textproto"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func (r *SQLStr) ListarBoletosComFiltro(ano, mes int, fornecedorID *int, statusIDs []int) ([]models.BoletoRelatorio, error) {
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

func (s *SQLStr) enviarEmailComAnexoPDF(emailAdmin string, fileContent []byte, mes, ano int) error {
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
func (s *SQLStr) GerarEEnviarRelatorioBoletos(emailAdmin string, ano, mes int, fornecedorID *int, statusIDs []int) error {
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
