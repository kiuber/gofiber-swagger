package main

import (
	"log"

	"github.com/TDiblik/gofiber-swagger/gofiberswagger"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	// equivalent to:
	// router := gofiberswagger.NewRouter(app)
	// router.Get("/", nil, HelloHandler)
	// router.Get("/abc", nil, HelloHandler)
	// router.Get("/bca", nil, HelloHandler)

	app.Get("/", HelloHandler)
	app.Get("/abc", HelloHandler)
	app.Get("/bca", HelloHandler)
	gofiberswagger.RegisterRoute("GET", "/", &gofiberswagger.RouteInfo{})
	gofiberswagger.RegisterRoute("GET", "/abc", &gofiberswagger.RouteInfo{})
	gofiberswagger.RegisterRoute("GET", "/bca", &gofiberswagger.RouteInfo{})

	// You can now see your:
	// - UI at /swagger/
	// - json at /swagger/swagger.json
	// - yaml at /swagger/swagger.yaml
	gofiberswagger.Register(app, &gofiberswagger.DefaultConfig)

	log.Fatal(app.Listen(":3000"))
}

// ----- Hello Handler and it's types ----- //
func HelloHandler(c fiber.Ctx) error {
	return c.SendStatus(200)
}
