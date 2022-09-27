package awesomeapi

import (
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
)

type Address struct {
	CEP         string `json:"cep"`
	AddressType string `json:"address_type"`
	AddressName string `json:"address_name"`
	Address     string `json:"address"`
	State       string `json:"state"`
	District    string `json:"district"`
	Lat         string `json:"lat"`
	Lng         string `json:"lng"`
	City        string `json:"city"`
	CityIBGE    string `json:"city_ibge"`
	DDD         string `json:"ddd"`
}

func Find(ch chan<- Address, zipcode string) error {
	var address Address
	url := fmt.Sprintf("https://cep.awesomeapi.com.br/json/%s", zipcode)
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
