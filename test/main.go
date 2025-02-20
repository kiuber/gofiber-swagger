package main

import (
	swagger "TDiblik/gofiber-swagger/src"

	"github.com/gofiber/fiber/v3"
)

type RequestBody struct {
	A string
	E Embed
}

type Embed struct {
	B string
}

func main() {
	app := fiber.New()

	with_swagger := swagger.NewRouter(app)
	with_swagger.Get("/", &swagger.RouteInfo{
		Parameters: swagger.Parameters{
			swagger.NewPathParameter("a"),
		},
		RequestBody: swagger.NewRequestBodyJSON[RequestBody](),
		Responses:   swagger.NewResponses(swagger.NewResponseJSON[RequestBody]()),
	}, func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	with_swagger.Get("/test", nil, func(c fiber.Ctx) error {
		return c.SendString("Test!")
	})
	abc := with_swagger.Group("/abc")
	abc.Get("/aa", nil, func(c fiber.Ctx) error {
		return c.SendString("ABC!")
	})
	abc.Get("/bb/*/*", nil, func(c fiber.Ctx) error {
		return c.SendString("ABC!")
	})
	abc.Get("/dd/:id", nil, func(c fiber.Ctx) error {
		return c.SendString("ABC!")
	})

	swagger.Register(app, swagger.DefaultConfig)

	app.Listen(":3000")
}
