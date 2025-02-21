package test

import (
	"github.com/TDiblik/gofiber-swagger/gofiberswagger"

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

	with_swagger := gofiberswagger.NewRouter(app)
	with_swagger.Get("/", &gofiberswagger.RouteInfo{
		Parameters: gofiberswagger.Parameters{
			gofiberswagger.NewQueryParameter("a"),
			gofiberswagger.NewHeaderParameter("a"),
		},
		RequestBody: gofiberswagger.NewRequestBodyJSON[RequestBody](),
		Responses: gofiberswagger.NewResponses(
			map[string]*gofiberswagger.ResponseRef{
				"200": gofiberswagger.NewResponseJSON[ResponseBody]("ok"),
				"400": gofiberswagger.NewResponseJSON[Embed]("fail"),
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
			map[string]*gofiberswagger.ResponseRef{
				"200": gofiberswagger.NewResponseJSON[ResponseBody]("ok"),
				"400": gofiberswagger.NewResponseJSON[Embed]("fail"),
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

	gofiberswagger.Register(app, gofiberswagger.DefaultConfig)

	app.Listen(":3000")
}
