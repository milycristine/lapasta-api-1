package models

type Motorista struct {
	Id                 int     `json:"id"`
	Nome               string  `json:"nome"`
	Telefone           string  `json:"telefone"`
	CPF                string  `json:"cpf"`
	ChavePix           string  `json:"chave_pix"`
	ValorDiaria        float64 `json:"valor_diaria"`
	Ativo              bool    `json:"ativo"`
	Email              string  `json:"email,omitempty"`
	Senha              string  `json:"senha,omitempty"`
	CnhFrenteUrl       string  `json:"cnh_frente_url,omitempty"`
	CnhVersoUrl        string  `json:"cnh_verso_url,omitempty"`
	KmDiaria           float64 `json:"km_diaria,omitempty"`
	ImagemFrenteBase64 string  `json:"imagem_frente,omitempty"`
	ImagemVersoBase64  string  `json:"imagem_verso,omitempty"`
}

type EmissaoNota struct {
	Id                 int     `json:"id"`
	NumeroNota         string  `json:"numero_nota"`
	Valor              float64 `json:"valor"`
	DataEmissao        string  `json:"data_emissao"`
	DataAtribuicao     string  `json:"data_atribuicao"`
	Descricao          string  `json:"descricao"`
	MotoristaId        int     `json:"motorista_id"`
	MotoristaNome      string  `json:"motorista_nome"`
	IdStatusLancamento int     `json:"id_status_lancamento"`
}

type StatusLancamento struct {
	Id        int    `json:"id"`
	Descricao string `json:"descricao"`
}

type NotasMotorista struct {
	Id                 int     `json:"id"`
	IdMotorista        int     `json:"id_motorista"`
	IdNota             int     `json:"id_nota"`
	DataAtribuicao     string  `json:"data_atribuicao"`
	DataLancamento     *string `json:"data_lancamento,omitempty"`
	IdStatusLancamento int     `json:"id_status_lancamento"`
	Url                string  `json:"url"`
	ImagemBase64       string  `json:"imagem"`
	NumeroNota         string  `json:"numero_nota"`
	Valor              float64 `json:"valor"`
	DataEmissao        string  `json:"data_emissao"`
	Descricao          string  `json:"descricao"`
	NomeMotorista      string  `json:"nome_motorista"`
}
