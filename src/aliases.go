package swagger

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/uuid"
)

type Operation = openapi3.Operation
type Origin = openapi3.Origin
type Location = openapi3.Location
type Parameters = openapi3.Parameters
type ParameterRef = openapi3.ParameterRef
type Parameter = openapi3.Parameter
type Examples = openapi3.Examples
type ExampleRef = openapi3.ExampleRef
type Example = openapi3.Example
type Content = openapi3.Content
type MediaType = openapi3.MediaType
type Encoding = openapi3.Encoding
type Headers = openapi3.Headers
type HeaderRef = openapi3.HeaderRef
type Header = openapi3.Header
type SchemaRef = openapi3.SchemaRef
type Schema = openapi3.Schema
type SchemaRefs = openapi3.SchemaRefs
type Types = openapi3.Types
type ExternalDocs = openapi3.ExternalDocs
type RequestBodyRef = openapi3.RequestBodyRef
type RequestBody = openapi3.RequestBody
type Responses = openapi3.Responses
type ResponseRef = openapi3.ResponseRef
type Callbacks = openapi3.Callbacks
type Links = openapi3.Links
type ParametersMap = openapi3.ParametersMap
type RequestBodies = openapi3.RequestBodies
type ResponseBodies = openapi3.ResponseBodies
type Schemas = openapi3.Schemas
type Components = openapi3.Components
type Response = openapi3.Response
type SecuritySchemes = openapi3.SecuritySchemes
type SecuritySchemeRef = openapi3.SecuritySchemeRef
type SecurityScheme = openapi3.SecurityScheme
type OAuthFlows = openapi3.OAuthFlows
type OAuthFlow = openapi3.OAuthFlow
type StringMap = openapi3.StringMap
type Info = openapi3.Info
type Contact = openapi3.Contact
type License = openapi3.License
type Paths = openapi3.Paths
type PathItem = openapi3.PathItem
type Servers = openapi3.Servers
type Server = openapi3.Server
type ServerVariable = openapi3.ServerVariable
type SecurityRequirements = openapi3.SecurityRequirements
type SecurityRequirement = openapi3.SecurityRequirement
type Tags = openapi3.Tags
type Tag = openapi3.Tag
type NewResponsesOption = openapi3.NewResponsesOption

func NewRequestBodyJSON[T any]() *RequestBodyRef {
	return NewRequestBodyJSONExtended[T]("", false)
}
func NewRequestBodyJSONExtended[T any](description string, required bool) *RequestBodyRef {
	request_body := openapi3.NewRequestBody()
	request_body.WithDescription(description)
	request_body.WithRequired(required)
	schema := CreateSchema[T]()
	request_body.WithJSONSchemaRef(schema)
	return &RequestBodyRef{Ref: schema.Ref, Value: request_body}
}

func NewRequestBodyFormData[T any]() *RequestBodyRef {
	return NewRequestBodyFormDataExtended[T]("", false)
}
func NewRequestBodyFormDataExtended[T any](description string, required bool) *RequestBodyRef {
	request_body := openapi3.NewRequestBody()
	request_body.WithDescription(description)
	request_body.WithRequired(required)
	schema := CreateSchema[T]()
	request_body.WithFormDataSchemaRef(schema)
	return &RequestBodyRef{Ref: schema.Ref, Value: request_body}
}

func NewResponses(responses ...*ResponseRef) *Responses {
	output := &Responses{}
	for _, r := range responses {
		output.Set(uuid.NewString(), r)
	}
	return output
}
func NewResponseJSON[T any]() *ResponseRef {
	response := openapi3.NewResponse()
	schema := CreateSchema[T]()
	response.WithJSONSchemaRef(schema)
	return &ResponseRef{Value: response}
}

func NewPathParameter(name string) *ParameterRef {
	return &ParameterRef{Value: openapi3.NewPathParameter(name)}
}

func NewQueryParameter(name string) *ParameterRef {
	return &ParameterRef{Value: openapi3.NewQueryParameter(name)}
}

func NewHeaderParameter(name string) *ParameterRef {
	return &ParameterRef{Value: openapi3.NewHeaderParameter(name)}
}

func NewCookieParameter(name string) *ParameterRef {
	return &ParameterRef{Value: openapi3.NewCookieParameter(name)}
}

func NewPathParameterRaw(name string) *Parameter {
	return openapi3.NewPathParameter(name)
}

func NewQueryParameterRaw(name string) *Parameter {
	return openapi3.NewQueryParameter(name)
}

func NewHeaderParameterRaw(name string) *Parameter {
	return openapi3.NewHeaderParameter(name)
}

func NewCookieParameterRaw(name string) *Parameter {
	return openapi3.NewCookieParameter(name)
}

var AcquiredSchemas map[string]*SchemaRef

func appendToAcquiredSchemas(ref string, schema *SchemaRef) {
	if AcquiredSchemas == nil {
		AcquiredSchemas = make(map[string]*SchemaRef)
	}
	if schema != nil {
		AcquiredSchemas[ref] = schema
	}
}
func getAcquiredSchemas(ref string) *SchemaRef {
	if AcquiredSchemas == nil {
		return nil
	}

	schema := AcquiredSchemas[ref]
	if schema == nil {
		return nil
	}

	return schema
}

func CreateSchema[T any]() *SchemaRef {
	var t T
	return generateSchema(reflect.TypeOf(t))
}

func generateSchema(t reflect.Type) *SchemaRef {
	ref := t.PkgPath() + t.Name()
	possible_schema := getAcquiredSchemas(ref)
	if possible_schema != nil {
		return possible_schema
	}

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	schema := &Schema{
		Type: &Types{
			getTypeString(t),
		},
		Format:     getFormatString(t),
		Properties: make(Schemas),
		// todo: check whether it's required by the validation library
	}

	if t.Kind() == reflect.Struct {
		schema.Title = t.Name()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if t.Kind() == reflect.Struct {
				fieldSchema := generateSchema(field.Type)
				schema.Properties[field.Name] = fieldSchema
				continue
			}
			schema.Properties[field.Name] = &SchemaRef{
				Value: &Schema{
					Type: &Types{
						getTypeString(field.Type),
					},
					Format: getFormatString(field.Type),
					// todo: check whether it's required by the validation library
				},
			}
		}
		appendToAcquiredSchemas(ref, &SchemaRef{
			Value: schema,
		})
		return &SchemaRef{
			Value: schema,
			Ref:   "#/components/schemas/" + ref,
		}
	}

	return &SchemaRef{
		Value: schema,
	}
}

func getTypeString(t reflect.Type) string {
	var typeString string
	switch t.Kind() {
	case reflect.String:
		typeString = "string"
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		typeString = "integer"
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		typeString = "integer"
		break
	case reflect.Float32, reflect.Float64:
		typeString = "number"
		break
	case reflect.Bool:
		typeString = "boolean"
		break
	case reflect.Struct:
		typeString = "object"
		break
	case reflect.Slice, reflect.Array:
		typeString = "array"
		break
	default:
		typeString = "unknown"
		break
	}
	return typeString
}

func getFormatString(t reflect.Type) string {
	var formatString string
	switch t.Kind() {
	case reflect.Int32:
		formatString = "i32"
		break
	case reflect.Int64:
		formatString = "i64"
		break
	default:
		formatString = ""
		break
	}
	return formatString
}
