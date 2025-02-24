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
			DocExpansion:             "list",
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

	// You can now see your:
	// - UI at /swagger/
	// - json at /swagger/swagger.json
	// - yaml at /swagger/swagger.yaml

	app.Listen(":3000")
}

// ----- Hello Handler and it's types ----- //
func HelloHandler(c fiber.Ctx) error {
	return c.SendStatus(200)
}
