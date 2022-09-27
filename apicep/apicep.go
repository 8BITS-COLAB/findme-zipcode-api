package apicep

import (
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
)

type Address struct {
	Code     string `json:"code"`
	State    string `json:"state"`
	City     string `json:"city"`
	District string `json:"district"`
	Address  string `json:"address"`
}

func Find(ch chan<- Address, zipcode string) error {
	var address Address
	url := fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", zipcode)
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
