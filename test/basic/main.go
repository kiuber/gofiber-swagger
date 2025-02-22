package main

import (
	"github.com/TDiblik/gofiber-swagger/gofiberswagger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	// Create wrapper around the fiber router
	router := gofiberswagger.NewRouter(app)

	// Normally create your routes.
	// You don't have to provide the swagger options (second argument)
	router.Get("/", nil, HelloHandler)

	// path parameters get automatically recognized :)
	parameters_group := router.Group("/parameters/")
	parameters_group.Get("/:id", nil, HandlerWithId)
	parameters_group.Post("/:id", &gofiberswagger.RouteInfo{
		Responses: gofiberswagger.NewResponses(
			gofiberswagger.NewResponseInfo[HandlerWithIdResponse]("200", "example response 👀"),
		),
	}, HandlerWithId)

	// wildcard_group := router.Group("/wildcards/")
	// wildcard_group.Get("/1/*", Handler)
	// wildcard_group.Get("/1/*", nil)

	// Register swagger. Without this line, nothing will get generated.
	// You probably want to skip this line for production builds.
	gofiberswagger.Register(app, gofiberswagger.DefaultConfig)

	app.Listen(":3000")
}

func HelloHandler(c fiber.Ctx) error {
	return c.SendStatus(200)
}

type HandlerWithIdResponse struct {
	Id string `json:"id"`
}

func HandlerWithId(c fiber.Ctx) error {
	id := c.Params("id", "no id provided!")
	response := HandlerWithIdResponse{
		Id: id,
	}
	return c.Status(200).JSON(response)
}

type HandlerWithBodyResponse struct {
	Status        string        `json:"status"`
	EmbeddedField EmbeddedField `json:"embedded_field"`
}
type EmbeddedField struct {
	A int32
	B string
	C []string
}

func HandlerWithBody(c fiber.Ctx) error {
	response := HandlerWithBodyResponse{
		Status: "ok",
		EmbeddedField: EmbeddedField{
			A: 0,
			B: "hey there :D",
			C: []string{"i", "am", "an", "array"},
		},
	}
	return c.Status(200).JSON(response)
}
