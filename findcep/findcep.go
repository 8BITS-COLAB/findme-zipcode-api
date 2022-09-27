package findcep

import (
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
)

type Address struct {
	UF          string `json:"uf"`
	Logradouro  string `json:"logradouro"`
	Tipo        string `json:"tipo"`
	Bairro      string `json:"bairro"`
	CEP         string `json:"cep"`
	Complemento string `json:"complemento"`
	Cidade      string `json:"cidade"`
}

func Find(ch chan<- Address, zipcode string) error {
	var address Address
	url := fmt.Sprintf("https://website.api.findcep.com/v1/cep/%s.json", zipcode)
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
