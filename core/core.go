package core

import (
	"fmt"
	"time"

	"github.com/8bits/findme/apicep"
	"github.com/8bits/findme/awesomeapi"
	"github.com/8bits/findme/findcep"
	"github.com/8bits/findme/viacep"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

type Address struct {
	ZipCode  string `json:"zip_code"`
	Address  string `json:"address"`
	State    string `json:"state"`
	District string `json:"district"`
	City     string `json:"city"`
}

var c *cache.Cache = cache.New(5*time.Minute, 10*time.Minute)

func Handle(ctx *fiber.Ctx) error {
	zipcode := ctx.Params("zipcode")

	if cached, ok := c.Get(zipcode); ok {
		fmt.Println("cached")
		return ctx.JSON(cached)
	}

	apicepCh := make(chan apicep.Address)
	awesomeApiCh := make(chan awesomeapi.Address)
	viacepCh := make(chan viacep.Address)
	findcepCh := make(chan findcep.Address)

	go viacep.Find(viacepCh, zipcode)
	go apicep.Find(apicepCh, zipcode)
	go awesomeapi.Find(awesomeApiCh, zipcode)
	go findcep.Find(findcepCh, zipcode)

	select {
	case msg := <-apicepCh:
		address := Address{
			ZipCode:  msg.Code,
			District: msg.District,
			State:    msg.State,
			Address:  msg.Address,
			City:     msg.City,
		}

		c.Set(zipcode, address, cache.DefaultExpiration)

		return ctx.JSON(address)
	case msg := <-awesomeApiCh:
		address := Address{
			ZipCode:  msg.CEP,
			District: msg.District,
			State:    msg.State,
			Address:  msg.Address,
			City:     msg.City,
		}

		c.Set(zipcode, address, cache.DefaultExpiration)

		return ctx.JSON(address)
	case msg := <-viacepCh:
		address := Address{
			ZipCode:  msg.CEP,
			District: msg.Bairro,
			State:    msg.UF,
			Address:  msg.Logradouro,
			City:     msg.Localidade,
		}

		c.Set(zipcode, address, cache.DefaultExpiration)

		return ctx.JSON(address)
	case msg := <-findcepCh:
		address := Address{
			ZipCode:  msg.CEP,
			District: msg.Bairro,
			State:    msg.UF,
			Address:  msg.Logradouro,
			City:     msg.Cidade,
		}

		c.Set(zipcode, address, cache.DefaultExpiration)

		return ctx.JSON(address)
	}
}
