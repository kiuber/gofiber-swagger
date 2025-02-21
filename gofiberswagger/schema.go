package gofiberswagger

import (
	"reflect"
	"strings"
)

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
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	ref := t.PkgPath() + t.Name()
	possibleSchema := getAcquiredSchemas(ref)
	if possibleSchema != nil {
		return possibleSchema
	}

	schema := &Schema{
		Type:       &Types{getTypeString(t)},
		Format:     getFormatString(t),
		Properties: make(Schemas),
		Required:   []string{},
	}

	if t.Kind() == reflect.Struct {
		schema.Title = t.Name()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldType := field.Type
			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			// Check if field is required using the validate tag
			validateTag := field.Tag.Get("validate")
			if isFieldRequired(validateTag) {
				schema.Required = append(schema.Required, field.Name)
			}

			// Recursively resolve struct fields
			if fieldType.Kind() == reflect.Struct {
				fieldSchema := generateSchema(fieldType)
				schema.Properties[field.Name] = fieldSchema
				continue
			}

			// Handle slices and arrays
			if fieldType.Kind() == reflect.Slice || fieldType.Kind() == reflect.Array {
				itemSchema := generateSchema(fieldType.Elem())
				schema.Properties[field.Name] = &SchemaRef{
					Value: &Schema{
						Type:  &Types{"array"},
						Items: itemSchema,
					},
				}
				continue
			}

			// Handle maps as objects with additional properties
			if fieldType.Kind() == reflect.Map {
				keyType := fieldType.Key().Kind()
				if keyType == reflect.String {
					valueSchema := generateSchema(fieldType.Elem())
					has := true
					schema.Properties[field.Name] = &SchemaRef{
						Value: &Schema{
							Type: &Types{"object"},
							AdditionalProperties: AdditionalProperties{
								Has:    &has,
								Schema: valueSchema,
							},
						},
					}
				} else {
					// Unsupported key type
					schema.Properties[field.Name] = &SchemaRef{
						Value: &Schema{
							Type: &Types{"object"},
						},
					}
				}
				continue
			}

			// Inline other field types
			schema.Properties[field.Name] = &SchemaRef{
				Value: &Schema{
					Type:   &Types{getTypeString(fieldType)},
					Format: getFormatString(fieldType),
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
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Struct:
		return "object"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map:
		return "object"
	default:
		return "unknown"
	}
}

func getFormatString(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Int32:
		return "int32"
	case reflect.Int64:
		return "int64"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.String:
		if t.Name() == "Time" {
			return "date-time"
		}
	}
	return ""
}

func isFieldRequired(tag string) bool {
	if tag == "" {
		return false
	}
	validations := strings.Split(tag, ",")
	for _, validation := range validations {
		if validation == "required" {
			return true
		}
	}
	return false
}
