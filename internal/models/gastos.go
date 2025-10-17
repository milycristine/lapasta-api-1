package models

type Gasto struct {
	Id         int     `json:"id"`
	Data       string  `json:"data"`
	Descricao  string  `json:"descricao"`
	Valor      float64 `json:"valor"`
	Tipo       string  `json:"tipo"`       
	Fornecedor string  `json:"fornecedor"`
}
