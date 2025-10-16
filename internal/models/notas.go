package models

type Nota struct {
	Id                 int     `json:"id"`
	Tipo               string  `json:"tipo"`
	Valor              float64 `json:"valor"`
	IdFuncionario      int     `json:"idFuncionario"`
	UrlImagem          *string `json:"urlImagem"`
	ImagemBase64       *string `json:"imagem"`
	Dia                string  `json:"dia"`
	NomeFuncionario    string  `json:"nomeFuncionario"`
	Descricao          string  `json:"descricao"`
	NumeroNota         string  `json:"numeroNota"`
	DataEmissao        string  `json:"dataEmissao"`
	IdFornecedor       *int    `json:"idFornecedor"`
	NomeFornecedor     *string  `json:"nomeFornecedor,omitempty"`
	IdPedidoFornecedor *int    `json:"idPedidoFornecedor"`
	FormaPagamento     string  `json:"formaPagamento"`
	OutroPrazoDias     *int    `json:"outroPrazoDias,omitempty"`
}
