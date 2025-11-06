package models

type Login struct {
	Username              string `json:"username,omitempty"`
	Password              string `json:"password,omitempty"`
	PasswordCriptografado []byte
	IsLogged              bool `json:"isLogged"`
}

type User struct {
	ID        int    `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Nome      string `json:"nome,omitempty"`
	Sobrenome string `json:"sobrenome,omitempty"`
	Admin     int    `json:"admin,omitempty"`
	Tipo      string `json:"tipo"`
}

// Response ...
type Response struct {
	Status string
	Error  string
	Data   interface{}
}
