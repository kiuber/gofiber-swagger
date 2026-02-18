package main

import (
	"log"

	"github.com/TDiblik/gofiber-swagger/gofiberswagger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	router := gofiberswagger.NewRouter(app)
	router.Get("/", &gofiberswagger.RouteInfo{
		Responses: gofiberswagger.NewResponses(
			gofiberswagger.NewResponseInfo[HelloHandlerResponse]("200", "example response 👀"),
		),
	}, HelloHandler)

	gofiberswagger.Register(app, gofiberswagger.DefaultConfig)

	log.Fatal(app.Listen(":3000"))
}

// ----- Hello Handler and it's types ----- //
// ----- Handler with custom body & embedded struct and it's types ----- //
type HelloHandlerResponse struct {
	Status string `json:"status"`
	Test
}
type Test struct {
	EmbeddedField
}
type EmbeddedField struct {
	A int32
	B string
	C []string
}

func HelloHandler(c fiber.Ctx) error {
	response := HelloHandlerResponse{
		Status: "ok",
		Test: Test{EmbeddedField: EmbeddedField{
			A: 0,
			B: "hey there :D",
			C: []string{"i", "am", "an", "array"},
		}},
	}
	return c.Status(200).JSON(response)
}
