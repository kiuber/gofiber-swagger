package main

import (
	"github.com/TDiblik/gofiber-swagger/gofiberswagger"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	// Create wrapper around the fiber router (optional, recommended)
	router := gofiberswagger.NewRouter(app)

	// Normally create your routes. You don't have to provide the swagger options (second argument)
	router.Get("/", nil, GETHelloHandler)

	// You can group your routes normally
	parametersGroup := router.Group("/parameters/")

	// Path parameters get automatically recognized :)
	parametersGroup.Get("/:id", nil, GETHandlerWithId)

	// This is how you specify RequestBody / ResponseBody / different parameters / any other openapi property tied to a request
	parametersGroup.Post("/:id", &gofiberswagger.RouteInfo{
		Parameters: gofiberswagger.NewParameters(
			// We can also specify additonal parameters, for example query parameters
			gofiberswagger.NewQueryParameter("queryParam"),
		),
		RequestBody: gofiberswagger.NewRequestBodyJSON[POSTHandlerWithIdRequestBody](),
		Responses: gofiberswagger.NewResponses(
			gofiberswagger.NewResponseInfo[POSTHandlerWithIdResponse]("200", "example response 👀"),
		),
	}, POSTHandlerWithId)

	// You can also manully register routes without touching your existing code!
	app.Get("/abc", GETHelloHandler)
	gofiberswagger.RegisterRoute("GET", "/", &gofiberswagger.RouteInfo{})

	// Register swagger. Without this line, nothing will get generated.
	// For more config customizability, see /examples/custom-config/main.go
	// You can now see your:
	// - UI at /swagger/
	// - json at /swagger/swagger.json
	// - yaml at /swagger/swagger.yaml
	gofiberswagger.Register(app, gofiberswagger.DefaultConfig)

	app.Listen(":3000")
}

func GETHelloHandler(c fiber.Ctx) error {
	return c.SendStatus(200)
}
func GETHandlerWithId(c fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"id": c.Params("id", "no id provided!")})
}

// ----- Handler with custom request body it's types ----- //
type POSTHandlerWithIdRequestBody struct {
	A int32    `json:"a" validate:"required,min=1,max=10"`
	B string   `json:"b" validate:"required"`
	C []string `json:"c" validate:"require,min=1"`
}
type EmbeddedField struct {
	A int32
	B string
	C []string
}
type POSTHandlerWithIdResponse struct {
	Status        string        `json:"status"`
	Id            string        `json:"id"`
	QueryParam    string        `json:"query_param"`
	EmbeddedField EmbeddedField `json:"embedded_field"`
}

func POSTHandlerWithId(c fiber.Ctx) error {
	id := c.Params("id", "no id provided!")
	queryParam := c.Query("queryParam", "no queryParam provided!")
	request_body := new(POSTHandlerWithIdRequestBody)
	if err := c.Bind().Body(request_body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "msg": "Invalid request body"})
	}

	// in real application, you'd want to validate the struct here,
	// however that would overcomplicate our basic example

	response := POSTHandlerWithIdResponse{
		Status:     "ok",
		Id:         id,
		QueryParam: queryParam,
		EmbeddedField: EmbeddedField{
			A: request_body.A,
			B: request_body.B,
			C: request_body.C,
		},
	}
	return c.Status(200).JSON(response)
}
