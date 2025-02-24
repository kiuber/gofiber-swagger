# Golang Fiber Swagger generation

This library is basically a wrapper around [github.com/getkin/kin-openapi](https://github.com/getkin/kin-openapi) (with re-exported types) for the [fiber web framework](https://gofiber.io/).

It uses context from `fiber.App` to generate routes (with parameters) and then uses additional context provided by the user of the library to generate `openapi` schemas from types.

You can either:

- a) Use the `gofiberswagger.NewRouter` to create a router which acts like the `fiber.Router`, but takes `RouteInfo` for swagger docs as the second argument.
- b)

```go
// examples/basic/main.go
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

	// wildcard path parameters get automatically recognized as well
	wildcard_group := router.Group("/wildcards/")
	wildcard_group.Get("/1/*", nil, HelloHandler)
	wildcard_group.Get("/2/*/*", nil, HelloHandler)

	// You can easily specify Request body and response body type.
	// From that, a schema will get generated. This schema respects the type given + the `json` and `validate` tags.
	req_body_group := router.Group("/request_body/")
	req_body_group.Post("/", &gofiberswagger.RouteInfo{
		RequestBody: gofiberswagger.NewRequestBodyJSON[HandlerWithRequestBodyRequest](),
		Responses: gofiberswagger.NewResponses(
			gofiberswagger.NewResponseInfo[HandlerWithIdResponse]("200", "example response 👀"),
		),
	}, HandlerWithRequestBody)

	// Register swagger. Without this line, nothing will get generated.
	// You probably want to skip this line for production builds.
	gofiberswagger.Register(app, gofiberswagger.DefaultConfig)

	app.Listen(":3000")
}

// ----- Hello Handler and it's types ----- //
func HelloHandler(c fiber.Ctx) error {
	return c.SendStatus(200)
}

// ----- Handler with Id response in body and it's types ----- //
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

// ----- Handler with custom body & embedded struct and it's types ----- //
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

// ----- Handler with custom request body it's types ----- //
type HandlerWithRequestBodyRequest struct {
	A int32    `json:"a" validate:"required,min=1,max=10"`
	B string   `json:"b" validate:"required"`
	C []string `json:"c" validate:"require,min=1"`
}

func HandlerWithRequestBody(c fiber.Ctx) error {
	request_body := new(HandlerWithRequestBodyRequest)
	if err := c.Bind().Body(request_body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "msg": "Invalid request body"})
	}
	// in real application, you'd want to validate the struct here,
	// however that would overcomplicate our basic example

	response := HandlerWithBodyResponse{
		Status: "ok",
		EmbeddedField: EmbeddedField{
			A: request_body.A,
			B: request_body.B,
			C: request_body.C,
		},
	}
	return c.Status(200).JSON(response)
}

```

```go
// examples/basic/custom-config.go

package main

import (
	"html/template"

	"github.com/TDiblik/gofiber-swagger/gofiberswagger"
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	// Create wrapper around the fiber router
	router := gofiberswagger.NewRouter(app)
	router.Get("/", nil, HelloHandler)

	// Provide the custom options you wanna.
	// You can go all out and basically use ANY openapi property!
	gofiberswagger.Register(app, gofiberswagger.Config{
		Swagger: gofiberswagger.SwaggerConfig{
			OpenAPI: "3.0.0",
			Info: &gofiberswagger.Info{
				Title:   "Title inside the generated files.",
				Version: "0.0.1",
			},
		},
		SwaggerUI: gofiberswagger.SwaggerUIConfig{
			URL:    "/swagger/swagger.yaml",
			Title:  "Swagger UI - title of the swagger UI page",
			Layout: "StandaloneLayout",
			Plugins: []template.JS{
				template.JS("SwaggerUIBundle.plugins.DownloadUrl"),
			},
			Presets: []template.JS{
				template.JS("SwaggerUIBundle.presets.apis"),
				template.JS("SwaggerUIStandalonePreset"),
			},
			DeepLinking:              true,
			DefaultModelsExpandDepth: 1,
			DefaultModelExpandDepth:  1,
			DefaultModelRendering:    "example",
			DocExpaosion:             "list",
			SyntaxHighlight: &gofiberswagger.SyntaxHighlightConfig{
				Activate: true,
				Theme:    "agate",
			},
			ShowMutatedRequest:     true,
			DisplayRequestDuration: true,
			PersistAuthorization:   true,
		},
		CreateSwaggerFiles:       true,
		SwaggerFilesPath:         "./generated/swagger",
		AppendMethodToTags:       false,
		FilterOutAppUse:          true,
		RequiredAuth:             nil,
		AutomaticallyRequireAuth: false,
	})

	app.Listen(":3000")
}

// ----- Hello Handler and it's types ----- //
func HelloHandler(c fiber.Ctx) error {
	return c.SendStatus(200)
}

```
