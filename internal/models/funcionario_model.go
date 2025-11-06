package models

type Funcionario struct {
	Id                    int      `json:"id,omitempty"`
	Nome                  string   `json:"nome,omitempty"`
	Sobrenome             string   `json:"sobrenome,omitempty"`
	Cpf                   string   `json:"cpf,omitempty"`
	Rg                    string   `json:"rg,omitempty"`
	DataNasc              string   `json:"data_nasc,omitempty"`
	Email                 string   `json:"email,omitempty"`
	Senha                 string   `json:"senha,omitempty"`
	Cargo                 string   `json:"cargo,omitempty"`
	DateAdmissao          string   `json:"data_admissao,omitempty"`
	DateFinal             *string  `json:"data_final,omitempty"`
	HoraEntrada           *string  `json:"hora_entrada,omitempty"`
	HoraSaida             *string  `json:"hora_saida,omitempty"`
	Salario               *float64 `json:"salario,omitempty"`
	Admin                 int      `json:"admin,omitempty"`
	Status                int      `json:"status"`
	ValeTransporteSemanal *float64 `json:"vtSemanal,omitempty"`
	ChavePix              *string  `json:"chavepix,omitempty"`
}
type FuncionarioComPontos struct {
	Funcionario Funcionario
	Pontos      []Ponto
	BotaoStatus string
}
