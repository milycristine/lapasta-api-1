package models

import (
	"database/sql"
)

type Ponto struct {
	Id                   int            `json:"id"`
	HManha               sql.NullString `json:"hManha"`
	HAlmocoSaida         sql.NullString `json:"hAlmocoSaida"`
	HAlmocoRetorno       sql.NullString `json:"hAlmocoRetorno"`
	HNoite               sql.NullString `json:"hNoite"`
	Dia                  string         `json:"dia"`
	Situacao             string         `json:"situacao"`
	IdFuncionario        int            `json:"idFuncionario"`
	EntradaRegistrada    bool           `json:"entradaRegistrada"`
	PausaRegistrada      bool           `json:"pausaRegistrada"`
	RetornoRegistrado    bool           `json:"retornoRegistrado"`
	SaidaRegistrada      bool           `json:"saidaRegistrada"`
	NomeFuncionario      string         `json:"nomeFuncionario"`
}
