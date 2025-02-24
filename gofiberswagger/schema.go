package gofiberswagger

import (
	"reflect"
	"strconv"
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

	ref := strings.ReplaceAll(strings.ReplaceAll(t.PkgPath(), "/", "_"), ".", "_") + t.Name()
	ref_path := "#/components/schemas/" + ref
	possibleSchema := getAcquiredSchemas(ref)
	if possibleSchema != nil {
		if t.Kind() == reflect.Struct {
			return &SchemaRef{
				Ref:        ref_path,
				Extensions: possibleSchema.Extensions,
				Origin:     possibleSchema.Origin,
				Value:      possibleSchema.Value,
			}
		}
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
			fieldKind := fieldType.Kind()
			if fieldKind == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			fieldName := field.Name

			// handle json tag
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
			jsonTagOptions := strings.Split(jsonTag, ",")
			if len(jsonTagOptions) > 0 && jsonTagOptions[0] != "" {
				fieldName = jsonTagOptions[0]
			}

			// create schema for the field
			var result *SchemaRef = nil
			switch fieldKind {
			case reflect.Struct:
				result = generateSchema(fieldType)

			case reflect.Slice, reflect.Array:
				result = &SchemaRef{
					Value: &Schema{
						Type:  &Types{"array"},
						Items: generateSchema(fieldType.Elem()),
					},
				}

			case reflect.Map:
				keyType := fieldType.Key().Kind()
				if keyType == reflect.String {
					valueSchema := generateSchema(fieldType.Elem())
					has := true
					result = &SchemaRef{
						Value: &Schema{
							Type: &Types{"object"},
							AdditionalProperties: AdditionalProperties{
								Has:    &has,
								Schema: valueSchema,
							},
						},
					}
				} else {
					result = &SchemaRef{
						Value: &Schema{
							Type: &Types{"object"},
						},
					}
				}

			default:
				result = &SchemaRef{
					Value: &Schema{
						Type:   &Types{getTypeString(fieldType)},
						Format: getFormatString(fieldType),
					},
				}
			}
			result.Value.Title = fieldName

			// handle json tag
			for i := 1; i < len(jsonTagOptions); i++ {
				option := jsonTagOptions[i]
				switch option {
				case "string":
					result.Value.Type = &Types{"string"}
				case "omitempty":
					result.Value.Description += "omitempty"
				case "omitzero":
					result.Value.Description += "omitzero"
				}
			}

			// handle validate tag
			validateTag := field.Tag.Get("validate")
			validateTagOptions := strings.Split(validateTag, ",")
			for _, validation := range validateTagOptions {
				switch {
				case validation == "required":
					schema.Required = append(schema.Required, fieldName)
					result.Value.AllowEmptyValue = false
				case strings.HasPrefix(validation, "min=") && (fieldKind == reflect.Slice || fieldKind == reflect.Array):
					if minValue, err := strconv.ParseUint(strings.TrimPrefix(validation, "min="), 10, 64); err == nil {
						result.Value.MinItems = minValue
					}
				case strings.HasPrefix(validation, "min="):
					if minValue, err := strconv.ParseFloat(strings.TrimPrefix(validation, "min="), 64); err == nil {
						result.Value.Min = &minValue
					}
				case strings.HasPrefix(validation, "max=") && (fieldKind == reflect.Slice || fieldKind == reflect.Array):
					if maxValue, err := strconv.ParseUint(strings.TrimPrefix(validation, "max="), 10, 64); err == nil {
						result.Value.MaxItems = &maxValue
					}
				case strings.HasPrefix(validation, "max="):
					if maxValue, err := strconv.ParseFloat(strings.TrimPrefix(validation, "max="), 64); err == nil {
						result.Value.Max = &maxValue
					}
				case strings.HasPrefix(validation, "minLength="):
					if minLen, err := strconv.ParseUint(strings.TrimPrefix(validation, "minLength="), 10, 64); err == nil {
						result.Value.MinLength = minLen
					}
				case strings.HasPrefix(validation, "maxLength="):
					if maxLen, err := strconv.ParseUint(strings.TrimPrefix(validation, "maxLength="), 10, 64); err == nil {
						result.Value.MaxLength = &maxLen
					}
				case strings.HasPrefix(validation, "uniqueItems"):
					result.Value.UniqueItems = true
				case strings.HasPrefix(validation, "oneof="):
					options := strings.Split(strings.TrimPrefix(validation, "oneof="), " ")
					if result.Value.OneOf == nil {
						result.Value.OneOf = []*SchemaRef{}
					}
					for _, option := range options {
						option_schema := NewStringSchema()
						option_schema.Default = option
						result.Value.OneOf = append(result.Value.OneOf, &SchemaRef{Value: option_schema})
					}
				}
			}

			schema.Properties[fieldName] = result
		}

		appendToAcquiredSchemas(ref, &SchemaRef{
			Value: schema,
		})
		return &SchemaRef{
			Ref:   ref_path,
			Value: schema,
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
