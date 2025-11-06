package documento

type Documento struct {
	Id              int    `json:"id"`
	Titulo          string `json:"titulo"`
	Url             string `json:"url"`
	ImagemBase64    string `json:"imagem"`
	DataCriacao     string `json:"dataCriacao"`
	IdFuncionario   int    `json:"idFuncionario"`
	NomeResponsavel string `json:"nomeResponsavel"`
	Descricao       string `json:"descricao"`
}
