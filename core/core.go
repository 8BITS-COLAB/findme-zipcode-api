package core

import (
	"time"

	"github.com/8bits/findme/apicep"
	"github.com/8bits/findme/awesomeapi"
	"github.com/8bits/findme/findcep"
	"github.com/8bits/findme/viacep"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v3"
)

type Address struct {
	ZipCode  string `json:"zip_code" yaml:"zip_code"`
	Address  string `json:"address"`
	State    string `json:"state"`
	District string `json:"district"`
	City     string `json:"city"`
}

var c *cache.Cache = cache.New(time.Minute, 5*time.Minute)

func reply(ctx *fiber.Ctx, address Address) error {
	fns := map[string]func(data interface{}) error{
		"json": ctx.JSON,
		"xml":  ctx.XML,
		"csv": func(data interface{}) error {
			addresses := []Address{
				data.(Address),
			}
			b, _ := gocsv.MarshalBytes(addresses)

			ctx.Response().Header.Set("Content-Type", "application/vnd.ms-excel")

			return ctx.Send(b)
		},
		"yaml": func(data interface{}) error {
			b, _ := yaml.Marshal(data)

			ctx.Response().Header.Set("Content-Type", "text/vnd.yaml")

			return ctx.Send(b)
		},
	}

	fn := fns[ctx.Query("format")]

	if fn == nil {
		fn = ctx.JSON
	}

	return fn(address)
}

func Handle(ctx *fiber.Ctx) error {
	zipcode := ctx.Params("zipcode")

	if cached, ok := c.Get(zipcode); ok {
		return reply(ctx, cached.(Address))
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

		return reply(ctx, address)
	case msg := <-awesomeApiCh:
		address := Address{
			ZipCode:  msg.CEP,
			District: msg.District,
			State:    msg.State,
			Address:  msg.Address,
			City:     msg.City,
		}

		c.Set(zipcode, address, cache.DefaultExpiration)

		return reply(ctx, address)
	case msg := <-viacepCh:
		address := Address{
			ZipCode:  msg.CEP,
			District: msg.Bairro,
			State:    msg.UF,
			Address:  msg.Logradouro,
			City:     msg.Localidade,
		}

		c.Set(zipcode, address, cache.DefaultExpiration)

		return reply(ctx, address)
	case msg := <-findcepCh:
		address := Address{
			ZipCode:  msg.CEP,
			District: msg.Bairro,
			State:    msg.UF,
			Address:  msg.Logradouro,
			City:     msg.Cidade,
		}

		c.Set(zipcode, address, cache.DefaultExpiration)

		return reply(ctx, address)
	}
}
