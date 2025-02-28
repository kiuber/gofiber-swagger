package main

import (
	"github.com/TDiblik/gofiber-swagger/gofiberswagger"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	router := gofiberswagger.NewRouter(app)
	router.Get("/", nil, HelloHandler)
	router.Post("/enum", &gofiberswagger.RouteInfo{
		RequestBody: gofiberswagger.NewRequestBody[EnumHandlerRequest](),
		Responses: gofiberswagger.NewResponses(
			gofiberswagger.NewResponseInfo[EnumHandlerRequest]("200", "OK"),
		),
	}, EnumHandler)

	// You can now see your:
	// - UI at /swagger/
	// - json at /swagger/swagger.json
	// - yaml at /swagger/swagger.yaml
	gofiberswagger.Register(app, gofiberswagger.DefaultConfig)

	app.Listen(":3000")
}

// ----- Hello Handler and it's types ----- //
func HelloHandler(c fiber.Ctx) error {
	return c.SendStatus(200)
}

// ----- Enum Handler and it's types ----- //

// Define int enum
type Status int

const (
	Pending Status = iota
	Active
	Inactive
	Deleted
)

// Implement ISwaggerEnum
func (e Status) EnumValues() []any {
	return []any{Pending, Active, Inactive, Deleted}
}

// Define string enum
type StringStatus string

const (
	IGNORED      StringStatus = "IGNORED"
	NOT_FINISHED StringStatus = "NOT_FINISHED"
	IN_PROGRESS  StringStatus = "IN_PROGRESS"
	DONE         StringStatus = "DONE"
)

// Implement ISwaggerEnum
func (e StringStatus) EnumValues() []any {
	return []any{IGNORED, NOT_FINISHED, IN_PROGRESS, DONE}
}

type EnumHandlerRequest struct {
	Status             Status
	StringStatus       StringStatus
	OneOfExampleNumber Status       `validate:"oneof=1 2 3 4"`
	OneOfExampleString StringStatus `validate:"oneof=IGNORED NOT_FINISHED IN_PROGRESS DONE"`
}

func EnumHandler(c fiber.Ctx) error {
	request_body := new(EnumHandlerRequest)
	if err := c.Bind().Body(request_body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "msg": "Invalid request body"})
	}

	return c.JSON(request_body)
}
