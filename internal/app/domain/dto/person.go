package dto

import (
	"errors"
	"strconv"
	"strings"

	"github.com/EvertonTomalok/go-template/internal/app/domain/entities"
	"github.com/klassmann/cpfcnpj"
)

type Person struct {
	Name string `json:"name"`
	CPF  string `json:"cpf"`
}

func (p *Person) DTOtoModel() (entities.Person, error) {
	cpfValidator := cpfcnpj.NewCPF(p.CPF)
	if !cpfValidator.IsValid() {
		return entities.Person{}, errors.New("CPF is invalid")
	}

	cpf, err := strconv.Atoi(sanitizeCPF(cpfValidator.String()))
	if err != nil {
		return entities.Person{}, err
	}

	person := entities.Person{
		Name: p.Name,
		CPF:  cpf,
	}
	return person, nil
}

func sanitizeCPF(cpf string) string {
	cpfStr := strings.ReplaceAll(cpf, ".", "")

	return strings.ReplaceAll(cpfStr, "-", "")
}
