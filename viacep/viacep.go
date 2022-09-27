package viacep

import (
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
)

type Address struct {
	CEP         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	GIA         string `json:"gia"`
	DDD         string `json:"ddd"`
	SIAFI       string `json:"siafi"`
}

func Find(ch chan<- Address, zipcode string) error {
	var address Address
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipcode)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := sonic.ConfigDefault.NewDecoder(res.Body).Decode(&address); err != nil {
		return err
	}

	ch <- address

	return nil
}
