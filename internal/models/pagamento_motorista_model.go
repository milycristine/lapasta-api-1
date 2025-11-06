package models

import "time"

type StatusPagamentoMotorista struct {
	Id        int    `json:"id"`
	Descricao string `json:"descricao"`
}

type PagamentosMotorista struct {
	Id                int       `json:"id"`
	IdMotorista       int       `json:"id_motorista"`
	SemanaInicio      time.Time `json:"semana_inicio"`
	SemanaFim         time.Time `json:"semana_fim"`
	DiasTrabalhados   int       `json:"dias_trabalhados"`
	ValorPago         float64   `json:"valor_pago"`
	DataPagamento     time.Time `json:"data_pagamento"`
	IdStatusPagamento int       `json:"id_status_pagamento"`
	NomeMotorista     string    `json:"nome_motorista,omitempty"`
	ChavePix          string    `json:"chave_pix,omitempty"`
}
