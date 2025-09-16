package main

import (
	"log"
	"mime/multipart"

	"github.com/TDiblik/gofiber-swagger/gofiberswagger"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	router := gofiberswagger.NewRouter(app)
	router.Get("/", nil, HelloHandler)
	router.Post("/upload", &gofiberswagger.RouteInfo{
		RequestBody: gofiberswagger.NewRequestBodyFormData[struct {
			some_file multipart.FileHeader `validate:"required"`
		}](),
		Responses: gofiberswagger.NewResponses(
			gofiberswagger.NewResponseInfo[struct {
				status string
				file   multipart.FileHeader
			}]("200", "OK"),
		),
	}, UploadHandler)

	// You can now see your:
	// - UI at /swagger/
	// - json at /swagger/swagger.json
	// - yaml at /swagger/swagger.yaml
	gofiberswagger.Register(app, gofiberswagger.DefaultConfig)

	log.Fatal(app.Listen(":3000"))
}

// ----- Hello Handler and it's types ----- //
func HelloHandler(c fiber.Ctx) error {
	return c.SendStatus(200)
}

func UploadHandler(c fiber.Ctx) error {
	file, err := c.FormFile("some_file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "msg": "This API endpoint requires \"some_file\" submitted as a form file."})
	}

	// todo: do anything you desire with the uploaded file

	return c.Status(200).JSON(fiber.Map{"status": "ok", "file": file})

}
