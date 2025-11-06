package models

import "time"

type Pagamento struct {
	Id                   int       `json:"id"`
	FuncionarioId        int       `json:"funcionarioId"`
	DataPagamento        time.Time `json:"dataPagamento"`
	Valor                float64   `json:"valor"`
	StatusId             int       `json:"statusId"`
	CpfFuncionario       string    `json:"cpfFuncionario"`
	NomeFuncionario      string    `json:"nomeFuncionario"`
	SobrenomeFuncionario string    `json:"sobrenomeFuncionario"`
	ChavePixFuncionario  string    `json:"chavePixFuncionario"`
}

type StatusPagamento struct {
	Id        int    `json:"id"`
	Descricao string `json:"descricao"`
}

type PagamentosPorDia struct {
	Id           int `json:"id"`
	PagamentoId  int `json:"pagamentoId"`
	DiaPagamento int `json:"diaPagamento"`
}
