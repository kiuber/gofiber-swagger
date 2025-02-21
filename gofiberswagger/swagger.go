package gofiberswagger

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v3"
	"gopkg.in/yaml.v3"
)

func Register(app *fiber.App, config Config) error {
	config.Swagger = swaggerConfigDefault(config.Swagger)
	config.SwaggerUI = swaggerUIConfigDefault(config.SwaggerUI)

	for k, v := range AcquiredSchemas {
		if config.Swagger.Components.Schemas[k] == nil {
			config.Swagger.Components.Schemas[k] = v
		}
	}

	routes := app.GetRoutes(config.FilterOutAppUse)
	for _, route := range routes {
		path_item := config.Swagger.Paths.Find(route.Path)
		if path_item == nil {
			path_item = &openapi3.PathItem{}
		}

		operation := getAcquiredRoutesInfo(route.Method, route.Path)
		if operation == nil {
			operation = &RouteInfo{}
		}

		for _, param := range route.Params {
			operation.AddParameter(openapi3.NewPathParameter(param))
		}
		if config.AppendMethodToTags {
			operation.Tags = append(operation.Tags, route.Method)
		}

		switch route.Method {
		case "POST":
			path_item.Post = operation
		case "CONNECT":
			path_item.Connect = operation
		case "DELETE":
			path_item.Delete = operation
		case "GET":
			path_item.Get = operation
		case "HEAD":
			path_item.Head = operation
		case "OPTIONS":
			path_item.Options = operation
		case "PATCH":
			path_item.Patch = operation
		case "PUT":
			path_item.Put = operation
		case "TRACE":
			path_item.Trace = operation
		default:
			log.Println("gofiber-swagger: unable to translate operation \"", route.Method, "\", skipping...")
		}

		config.Swagger.Paths.Set(route.Path, path_item)
	}

	index_page, err := generateIndexPage(swaggerUIConfigDefault(config.SwaggerUI))
	if err != nil {
		return err
	}
	schema_as_json, schema_as_yaml, err := generateOpenApiSchema(config.Swagger)
	if err != nil {
		return err
	}

	if config.CreateSwaggerFiles && !fiber.IsChild() {
		if config.SwaggerFilesPath == "" {
			return errors.New("CreateSwaggerFiles was set to true, however SwaggerFilesPaths was left empty.")
		}
		createSwaggerFiles(config.SwaggerFilesPath, index_page, schema_as_json, schema_as_yaml)
	}

	swagger_routes := app.Group("/swagger")
	swagger_routes.Get("/", func(c fiber.Ctx) error {
		return c.Type("html").Send(index_page)
	})
	swagger_routes.Get("/swagger.json", func(c fiber.Ctx) error {
		return c.Type("json").Send(schema_as_json)
	})
	swagger_routes.Get("/swagger.yaml", func(c fiber.Ctx) error {
		return c.Type("yaml").Send(schema_as_yaml)
	})

	return nil
}

func generateIndexPage(ui_config SwaggerUIConfig) (index_page []byte, err error) {
	index_tpl, err := template.New("swagger_index.html").Parse(indexPageTmpl)
	if err != nil {
		log.Println("gofiber-swagger: error while parsing the swagger index template -> ", err)
		return nil, err
	}
	index_tpl_buf := bytes.NewBufferString("")
	err = index_tpl.Execute(index_tpl_buf, ui_config)
	if err != nil {
		log.Println("gofiber-swagger: error while executing the swagger index template -> ", err)
		return nil, err
	}
	return index_tpl_buf.Bytes(), nil
}

func generateOpenApiSchema(schema openapi3.T) (as_json, as_yaml []byte, err error) {
	schema_as_yaml_raw, err := schema.MarshalYAML()
	if err != nil {
		log.Println("gofiber-swagger: error while creating the yaml schema -> ", err)
		return nil, nil, err
	}
	schema_as_yaml, err := yaml.Marshal(schema_as_yaml_raw)
	if err != nil {
		log.Println("gofiber-swagger: error while converting the yaml schema to yaml -> ", err)
		return nil, nil, err
	}

	schema_as_json, err := schema.MarshalJSON()
	if err != nil {
		log.Println("gofiber-swagger: error while creating the json schema -> ", err)
		return nil, nil, err
	}

	return schema_as_json, schema_as_yaml, nil
}

func createSwaggerFiles(target_folder_path string, index_page []byte, schema_as_json []byte, schema_as_yaml []byte) error {
	var creation_perms os.FileMode = 0o766

	swagger_dir := filepath.Dir(target_folder_path)
	if err := os.MkdirAll(swagger_dir, creation_perms); err != nil {
		return errors.Join(errors.New("unable to create file directory for swagger files"), err)
	}

	if err := os.WriteFile(filepath.Join(swagger_dir, "index.html"), index_page, creation_perms); err != nil {
		return errors.Join(errors.New("unable to create index.html for swagger files"), err)
	}

	if err := os.WriteFile(filepath.Join(swagger_dir, "swagger.json"), schema_as_json, creation_perms); err != nil {
		return errors.Join(errors.New("unable to create swagger.json for swagger files"), err)
	}

	if err := os.WriteFile(filepath.Join(swagger_dir, "swagger.yaml"), schema_as_yaml, creation_perms); err != nil {
		return errors.Join(errors.New("unable to create swagger.yaml for swagger files"), err)
	}

	return nil
}
