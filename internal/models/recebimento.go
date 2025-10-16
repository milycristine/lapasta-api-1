package models

type Recebimento struct {
	Id                 int     `json:"id"`
	Dia                string  `json:"dia"`
	Produto            *string `json:"produto"`
	UrlImagem          *string `json:"urlImagem"`
	ImagemBase64       *string `json:"imagem"`
	IdResponsavel      uint    `json:"idResponsavel"`
	NomeResponsavel    string  `json:"nomeResponsavel"`
	Quantidade         int     `json:"quantidade"`
	Peso               float64 `json:"peso"`
	Valor              float64 `json:"valor"`
	NumeroNota         string  `json:"numeroNota"`
	Vencimento         string  `json:"vencimento"`
	IdNota             int     `json:"idNota"`
	IdPedidoFornecedor *int    `json:"idPedidoFornecedor,omitempty"`
	NomeFornecedor     *string `json:"fornecedor"`
	FormaPagamento     string  `json:"formaPagamento"` //ex: 7 dias, 14 dias, 21 dias, outro
	OutroPrazoDias     *int    `json:"outroPrazoDias,omitempty"`
}

type DadosRecebimentoNota struct {
	IdNota             int
	NumeroNota         string
	IdPedidoFornecedor *int
	Produto            string
	Fornecedor         string
	FormaPagamento     string
	OutroPrazoDias     *int
}
