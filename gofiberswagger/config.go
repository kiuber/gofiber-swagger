package gofiberswagger

import "github.com/getkin/kin-openapi/openapi3"

// https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.1.1.md#openapi-object
type SwaggerConfig = openapi3.T

type Config struct {
	Swagger                  SwaggerConfig
	SwaggerUI                SwaggerUIConfig
	CreateSwaggerFiles       bool
	SwaggerFilesPath         string
	AppendMethodToTags       bool
	FilterOutAppUse          bool
	RequiredAuth             *openapi3.SecurityRequirements
	AutomaticallyRequireAuth bool
	CallbackBeforeGenerate   func(config *Config) error
}

var DefaultSwaggerConfig = SwaggerConfig{
	OpenAPI: "3.1.1",
	Info: &Info{
		Title:   DefaultUIConfig.Title,
		Version: "0.0.1",
	},
	Components: &Components{},
	Paths:      &Paths{},
}
var DefaultConfig = Config{
	Swagger:                  DefaultSwaggerConfig,
	SwaggerUI:                DefaultUIConfig,
	CreateSwaggerFiles:       true,
	SwaggerFilesPath:         "./generated/swagger",
	AppendMethodToTags:       false,
	FilterOutAppUse:          true,
	RequiredAuth:             nil,
	AutomaticallyRequireAuth: false,
	CallbackBeforeGenerate:   nil,
}

func swaggerConfigDefault(config SwaggerConfig) SwaggerConfig {
	cfg := config

	if cfg.Info == nil {
		cfg.Info = DefaultSwaggerConfig.Info
	}
	if cfg.Info.Title == "" {
		cfg.Info.Title = DefaultSwaggerConfig.Info.Title
	}
	if cfg.Info.Version == "" {
		cfg.Info.Version = DefaultSwaggerConfig.Info.Version
	}

	if cfg.Paths == nil {
		cfg.Paths = DefaultSwaggerConfig.Paths
	}

	if cfg.Components == nil {
		cfg.Components = DefaultSwaggerConfig.Components
	}
	if cfg.Components.Schemas == nil {
		cfg.Components.Schemas = make(map[string]*SchemaRef)
	}

	return cfg
}
