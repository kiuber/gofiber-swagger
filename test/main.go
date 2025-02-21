package main

import (
	swagger "TDiblik/gofiber-swagger/src"

	"github.com/gofiber/fiber/v3"
)

type RequestBody struct {
	A string `validate:"required"`
	E Embed
}

type Embed struct {
	B string
}

type ResponseBody struct {
	C string `validate:"required"`
}

func main() {
	app := fiber.New()

	with_swagger := swagger.NewRouter(app)
	with_swagger.Get("/", &swagger.RouteInfo{
		Parameters: swagger.Parameters{
			swagger.NewQueryParameter("a"),
			swagger.NewHeaderParameter("a"),
		},
		RequestBody: swagger.NewRequestBodyJSON[RequestBody](),
		Responses: swagger.NewResponses(
			map[string]*swagger.ResponseRef{
				"200": swagger.NewResponseJSON[ResponseBody]("ok"),
				"400": swagger.NewResponseJSON[Embed]("fail"),
			},
		),
	}, func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	with_swagger.Post("/", &swagger.RouteInfo{
		Parameters: swagger.Parameters{
			swagger.NewPathParameter("a"),
		},
		RequestBody: swagger.NewRequestBodyJSON[RequestBody](),
		Responses: swagger.NewResponses(
			map[string]*swagger.ResponseRef{
				"200": swagger.NewResponseJSON[ResponseBody]("ok"),
				"400": swagger.NewResponseJSON[Embed]("fail"),
			},
		),
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
