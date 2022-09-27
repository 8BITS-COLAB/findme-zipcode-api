package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/8bits/findme/core"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New()
	v1 := app.Group("/api/v1")

	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())

	app.Get("/_health_check", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(http.StatusOK)
	})
	v1.Get("/:zipcode", core.Handle)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))

	app.Listen(addr)
}
