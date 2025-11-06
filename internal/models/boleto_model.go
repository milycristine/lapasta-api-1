package models

import "time"

type Boleto struct {
	Id             int     `json:"id,omitempty"`
	RecebimentoId  int     `json:"recebimento_id,omitempty"`
	CodigoBarras   string  `json:"codigo_barras,omitempty"`
	DataCadastro   string  `json:"data_cadastro,omitempty"`
	DataVencimento string  `json:"data_vencimento,omitempty"`
	Valor          float64 `json:"valor,omitempty"`
	StatusId       int     `json:"statusId"`
	DataPagamento  *string `json:"data_pagamento,omitempty"`
	FornecedorNome string  `json:"fornecedorNome,omitempty"`
}

type StatusBoleto struct {
	Id        int    `json:"id"`
	Descricao string `json:"descricao"`
}

type BoletoRelatorio struct {
	FornecedorNome string
	DataVencimento time.Time
	Valor          float64
	CodigoBarras   string
	StatusId       int
}

func (b *Boleto) DataPagamentoOrDefault() string {
	if b.DataPagamento != nil {
		return *b.DataPagamento
	}
	return ""
}
