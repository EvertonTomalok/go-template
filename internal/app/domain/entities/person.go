package entities

type Person struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	CPF  int    `json:"cpf"`
}
