package gofiberswagger

import (
	"github.com/getkin/kin-openapi/openapi3"
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
type AdditionalProperties = openapi3.AdditionalProperties

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

func NewResponses(responses map[string]*ResponseRef) *Responses {
	output := &Responses{}
	for k, v := range responses {
		output.Set(k, v)
	}
	return output
}
func NewResponseJSON[T any](description string) *ResponseRef {
	response := openapi3.NewResponse()
	schema := CreateSchema[T]()
	response.WithJSONSchemaRef(schema)
	response.WithDescription(description)
	return &ResponseRef{Value: response}
}

func NewPathParameter(name string) *ParameterRef {
	return NewPathParameterExtended(name, &Schema{
		Type: &Types{"string"},
	})
}

func NewPathParameterWithType(name string, Type string) *ParameterRef {
	return NewPathParameterExtended(name, &Schema{
		Type: &Types{Type},
	})
}

func NewPathParameterExtended(name string, schema *Schema) *ParameterRef {
	return &ParameterRef{Value: openapi3.NewPathParameter(name).WithSchema(schema)}
}

func NewQueryParameter(name string) *ParameterRef {
	return NewQueryParameterExtended(name, &Schema{
		Type: &Types{"string"},
	})
}

func NewQueryParameterWithType(name string, Type string) *ParameterRef {
	return NewQueryParameterExtended(name, &Schema{
		Type: &Types{Type},
	})
}

func NewQueryParameterExtended(name string, schema *Schema) *ParameterRef {
	return &ParameterRef{Value: openapi3.NewQueryParameter(name).WithSchema(schema)}
}

func NewHeaderParameter(name string) *ParameterRef {
	return NewHeaderParameterExtended(name, &Schema{
		Type: &Types{"string"},
	})
}

func NewHeaderParameterWithType(name string, Type string) *ParameterRef {
	return NewHeaderParameterExtended(name, &Schema{
		Type: &Types{Type},
	})
}

func NewHeaderParameterExtended(name string, schema *Schema) *ParameterRef {
	return &ParameterRef{Value: openapi3.NewHeaderParameter(name).WithSchema(schema)}
}

func NewCookieParameter(name string) *ParameterRef {
	return NewCookieParameterExtended(name, &Schema{
		Type: &Types{"string"},
	})
}

func NewCookieParameterWithType(name string, Type string) *ParameterRef {
	return NewCookieParameterExtended(name, &Schema{
		Type: &Types{Type},
	})
}

func NewCookieParameterExtended(name string, schema *Schema) *ParameterRef {
	return &ParameterRef{Value: openapi3.NewCookieParameter(name).WithSchema(schema)}
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
