package models

type Fornecedor struct {
	Id       int    `json:"id,omitempty"`
	Nome     string `json:"nome,omitempty"`
	CNPJ     string `json:"cnpj,omitempty"`
	Email    string `json:"email,omitempty"`
	Telefone string `json:"telefone,omitempty"`
}

type FornecedorComHistorico struct {
	Fornecedor Fornecedor `json:"fornecedor"`
	Boletos    []Boleto   `json:"boletos"`
}

type PedidoFornecedor struct {
	Id                int    `json:"id,omitempty"`
	FornecedorId      int    `json:"fornecedorId"`
	Descricao         string `json:"descricao"`
	DataPedido        string `json:"dataPedido"`
	PrazoAcordadoDias int    `json:"prazoAcordadoDias"`
}
