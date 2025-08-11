package sql

import (
	"database/sql"
	"errors"
	"lapasta/internal/models"
	"time"
)

func (s *SQLStr) BuscarMotoristasAtivos() ([]models.Motorista, error) {
	rows, err := s.db.Query(`SELECT id, nome, cpf FROM Motoristas WHERE ativo = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var motoristas []models.Motorista
	for rows.Next() {
		var m models.Motorista
		if err := rows.Scan(&m.Id, &m.Nome, &m.CPF); err != nil {
			return nil, err
		}
		motoristas = append(motoristas, m)
	}

	if len(motoristas) == 0 {
		return nil, errors.New("nenhum motorista ativo encontrado")
	}
	return motoristas, nil
}

func (s *SQLStr) SalvarEmissaoNota(nota models.EmissaoNota) error {
	var exists bool
	err := s.db.QueryRow(`
		SELECT CASE WHEN EXISTS (
			SELECT 1 FROM EmissaoNotas WHERE NumeroNota = @NumeroNota
		) THEN 1 ELSE 0 END`,
		sql.Named("NumeroNota", nota.NumeroNota),
	).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if exists {
		return nil
	}

	_, err = s.db.Exec(`
		INSERT INTO EmissaoNotas (NumeroNota, Valor, DataEmissao, DataAtribuicao, Descricao, MotoristaId)
		VALUES (@NumeroNota, @Valor, @DataEmissao, @DataAtribuicao, @Descricao, @MotoristaId)`,
		sql.Named("NumeroNota", nota.NumeroNota),
		sql.Named("Valor", nota.Valor),
		sql.Named("DataEmissao", nota.DataEmissao),
		sql.Named("DataAtribuicao", time.Now().Format("2006-01-02")),
		sql.Named("Descricao", nota.Descricao),
		sql.Named("MotoristaId", nota.MotoristaId),
	)
	return err
}
