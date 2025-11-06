package models

import "time"

type ValeTransportePagamento struct {
	Id                   int       `json:"id"`
	IdFuncionario        int       `json:"id_funcionario"`
	DataPagamento        time.Time `json:"data_pagamento"`
	StatusIdVale         int       `json:"statusIdVale"`
	Valor                float64   `json:"valor"`
	CpfFuncionario       string    `json:"cpfFuncionario"`
	NomeFuncionario      string    `json:"nomeFuncionario"`
	SobrenomeFuncionario string    `json:"sobrenomeFuncionario"`
	ChavePixFuncionario  string    `json:"chavepix,omitempty"`
}

type StatusValeTransporte struct {
	Id        int    `json:"id"`
	Descricao string `json:"descricao"`
}
