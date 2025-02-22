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
			if fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}

			fieldName := field.Name
			field_schema_additonal_info := &Schema{}

			// handle json tag
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
			jsonOptions := strings.Split(jsonTag, ",")
			if len(jsonOptions) > 0 {
				if jsonOptions[0] != "" {
					field_schema_additonal_info.Title = fieldName
					fieldName = jsonOptions[0]
				}

				for i := 1; i < len(jsonOptions); i++ {
					option := jsonOptions[i]
					switch option {
					case "string":
						field_schema_additonal_info.Type = &Types{"string"}
					case "omitempty":
						field_schema_additonal_info.Description += "omitempty"
					case "omitzero":
						field_schema_additonal_info.Description += "omitzero"
					}
				}
			}

			// handle validate tag
			validateTag := field.Tag.Get("validate")
			validationOptions := strings.Split(validateTag, ",")
			for _, validation := range validationOptions {
				switch {
				case validation == "required":
					schema.Required = append(schema.Required, fieldName)
				// case strings.HasPrefix(validation, "min="):
				// 	minValue := strings.TrimPrefix(validation, "min=")
				// 	field_schema_additonal_info.Min = minValue
				// case strings.HasPrefix(validation, "max="):
				// 	maxValue := strings.TrimPrefix(validation, "max=")
				// 	field_schema_additonal_info.Max = maxValue
				// case strings.HasPrefix(validation, "minLength="):
				// 	minLen := strings.TrimPrefix(validation, "minLength=")
				// 	field_schema_additonal_info.MinLength = minLen
				// case strings.HasPrefix(validation, "maxLength="):
				// 	maxLen := strings.TrimPrefix(validation, "maxLength=")
				// 	field_schema_additonal_info.MaxLength = maxLen
				// case strings.HasPrefix(validation, "pattern="):
				// 	pattern := strings.TrimPrefix(validation, "pattern=")
				// 	field_schema_additonal_info.Pattern = pattern
				// case strings.HasPrefix(validation, "enum="):
				// 	enumValues := strings.TrimPrefix(validation, "enum=")
				// 	field_schema_additonal_info.Enum = strings.Split(enumValues, "|")
				// case strings.HasPrefix(validation, "minItems="):
				// 	minItems := strings.TrimPrefix(validation, "minItems=")
				// 	field_schema_additonal_info.MinItems = minItems
				// case strings.HasPrefix(validation, "maxItems="):
				// 	maxItems := strings.TrimPrefix(validation, "maxItems=")
				// 	field_schema_additonal_info.MaxItems = maxItems
				// case strings.HasPrefix(validation, "uniqueItems"):
				// 	field_schema_additonal_info.UniqueItems = true
				// case strings.HasPrefix(validation, "multipleOf="):
				// 	multipleOf := strings.TrimPrefix(validation, "multipleOf=")
				// 	field_schema_additonal_info.MultipleOf = multipleOf
				default:
					continue
				}
			}

			// create schema for the field
			var result *SchemaRef = nil
			switch fieldType.Kind() {
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

			// overwrite result with aditional preset info taken from tags
			if field_schema_additonal_info.Title != "" {
				result.Value.Title = field_schema_additonal_info.Title
			}
			if field_schema_additonal_info.Type != nil {
				result.Value.Type = field_schema_additonal_info.Type
			}
			if field_schema_additonal_info.Description != "" {
				result.Value.Description += field_schema_additonal_info.Description
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
