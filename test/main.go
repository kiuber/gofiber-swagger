package main

import (
	"github.com/TDiblik/gofiber-swagger/gofiberswagger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
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

	app.Use(cors.New())
	app.Use(logger.New())

	with_swagger_a := gofiberswagger.NewRouter(app)
	with_swagger := with_swagger_a.Group("/abc/cba/")
	with_swagger.Get("/", &gofiberswagger.RouteInfo{
		Parameters: gofiberswagger.Parameters{
			gofiberswagger.NewQueryParameter("a"),
			gofiberswagger.NewHeaderParameter("a"),
		},
		RequestBody: gofiberswagger.NewRequestBodyJSON[RequestBody](),
		Responses: gofiberswagger.NewResponsesRaw(
			map[string]*gofiberswagger.ResponseRef{
				"200": gofiberswagger.NewResponseRawJSON[ResponseBody]("ok"),
				"400": gofiberswagger.NewResponseRawJSON[Embed]("fail"),
			},
		),
	}, func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	with_swagger.Post("/", &gofiberswagger.RouteInfo{
		Parameters: gofiberswagger.Parameters{
			gofiberswagger.NewPathParameter("a"),
		},
		RequestBody: gofiberswagger.NewRequestBodyJSON[RequestBody](),
		Responses: gofiberswagger.NewResponses(
			gofiberswagger.NewResponseInfo[ResponseBody]("200", "ok"),
			gofiberswagger.NewResponseInfo[Embed]("400", "fail"),
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

	gofiberswagger.Register(app, gofiberswagger.DefaultConfig)

	app.Listen(":3000")
}
