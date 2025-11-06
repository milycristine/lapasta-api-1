package models

type TinyExpedicaoResponse struct {
	Retorno struct {
		Agrupamentos []AgrupamentoTiny `json:"agrupamentos"`
	} `json:"retorno"`
}

type AgrupamentoTiny struct {
	IdAgrupamento string             `json:"idAgrupamento"`
	Data          string             `json:"data"`
	FormaEnvio    string             `json:"formaEnvio"`
	Expedicoes    []ExpedicaoWrapper `json:"expedicoes"`
}

type ExpedicaoWrapper struct {
	Expedicao ExpedicaoTiny `json:"expedicao"`
}

type ExpedicaoTiny struct {
	Id             string           `json:"id"`
	TipoObjeto     string           `json:"tipoObjeto"`
	IdObjeto       string           `json:"idObjeto"` 
	IdAgrupamento  string           `json:"idAgrupamento"`
	Situacao       string           `json:"situacao"`
	DataEmissao    string           `json:"dataEmissao"`
	FormaEnvio     string           `json:"formaEnvio"`
	Identificacao  string           `json:"identificacao"` 
	QtdVolumes     string           `json:"qtdVolumes"`
	ValorDeclarado string           `json:"valorDeclarado"`
	PossuiValor    string           `json:"possuiValorDeclarado"`
	PesoBruto      string           `json:"pesoBruto"`
	Destinatario   DestinatarioTiny `json:"destinatario"`
}

type DestinatarioTiny struct {
	Nome        string `json:"nome"`
	Endereco    string `json:"endereco"`
	Numero      string `json:"numero"`
	Bairro      string `json:"bairro"`
	Cidade      string `json:"cidade"`
	Uf          string `json:"uf"`
	Complemento string `json:"complemento"`
	Cep         string `json:"cep"`
}

type ConfigTiny struct {
	Id      int    `json:"id"`
	NomeApi string `json:"nome_api"`
	Token   string `json:"token"`
}

type ExpedicaoNotaMotorista struct {
	MotoristaCPF  string
	MotoristaNome string
	NumExpedicoes int
	ValorTotal    float64
	PesoTotal     float64
	Expedicoes    []ExpedicaoTiny
}
