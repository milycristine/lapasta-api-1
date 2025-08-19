package sql

import (
	"database/sql"
	"fmt"
	"lapasta/internal/models"
	"time"
)

func (s *SQLStr) CriarDocumento(doc *models.Documento) error {
	var existe bool
	checkQuery := `
		SELECT 1 FROM Documentos 
		WHERE Titulo = @Titulo AND IdFuncionario = @IdFuncionario
	`
	err := s.db.QueryRow(checkQuery,
		sql.Named("Titulo", doc.Titulo),
		sql.Named("IdFuncionario", doc.IdFuncionario),
	).Scan(&existe)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("erro ao verificar duplicidade de documento: %w", err)
	}
	if existe {
		return fmt.Errorf("já existe um documento com este título para o funcionário")
	}
	query := `
		INSERT INTO Documentos (Titulo, Url, Data_Criacao, IdFuncionario, Descricao) 
		VALUES (@Titulo, @Url, @Data_Criacao, @IdFuncionario, @Descricao)
	`
	_, err = s.db.Exec(query,
		sql.Named("Titulo", doc.Titulo),
		sql.Named("Url", doc.Url),
		sql.Named("Data_Criacao", doc.DataCriacao),
		sql.Named("IdFuncionario", doc.IdFuncionario),
		sql.Named("Descricao", doc.Descricao),
	)

	if err != nil {
		return fmt.Errorf("erro ao inserir documento: %w", err)
	}
	return nil
}

func (s *SQLStr) ListarDocumentos(page int) ([]models.Documento, error) {
	var docs []models.Documento

	limit := 10
	offset := (page - 1) * limit

	query := `
        SELECT d.Id, d.Titulo, d.Url, d.Data_Criacao, d.IdFuncionario, 
               f.Nome AS NomeResponsavel, d.Descricao
        FROM Documentos d WITH (NOLOCK)
        JOIN Funcionarios f WITH (NOLOCK) ON d.IdFuncionario = f.Id
        ORDER BY d.Data_Criacao DESC
        OFFSET @Offset ROWS
        FETCH NEXT @Limit ROWS ONLY
    `

	rows, err := s.db.Query(query, sql.Named("Offset", offset), sql.Named("Limit", limit))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar documentos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d models.Documento
		if err := rows.Scan(&d.Id, &d.Titulo, &d.Url, &d.DataCriacao, &d.IdFuncionario, &d.NomeResponsavel, &d.Descricao); err != nil {
			return nil, fmt.Errorf("erro ao escanear documento: %w", err)
		}
		docs = append(docs, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return docs, nil
}

func (s *SQLStr) FiltrarDataDocumento(inicioData time.Time, fimData time.Time) ([]models.Documento, error) {
	var documentos []models.Documento

	query := `
		SELECT d.Id, d.Titulo, d.Url, d.Data_Criacao, d.IdFuncionario, 
		       f.Nome AS NomeResponsavel, d.Descricao
		FROM Documentos d WITH (NOLOCK)
		JOIN Funcionarios f WITH (NOLOCK) ON d.IdFuncionario = f.Id
		WHERE 
		d.Data_Criacao BETWEEN @InicioData AND @FimData
		ORDER BY d.Data_Criacao DESC
	`

	rows, err := s.db.Query(query, sql.Named("InicioData", inicioData), sql.Named("FimData", fimData))
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar documentos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Documento
		if err := rows.Scan(&p.Id, &p.Titulo, &p.Url, &p.DataCriacao, &p.IdFuncionario, &p.NomeResponsavel, &p.Descricao); err != nil {
			return nil, fmt.Errorf("erro ao escanear documentos: %w", err)
		}
		documentos = append(documentos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração das linhas: %w", err)
	}

	return documentos, nil
}
