package gofiberswagger

import (
	"math/rand/v2"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var acquiredSchemas map[string]*SchemaRef

func appendToAcquiredSchemas(ref string, schema *SchemaRef) {
	if acquiredSchemas == nil {
		acquiredSchemas = make(map[string]*SchemaRef)
	}
	if schema != nil {
		acquiredSchemas[ref] = schema
	}
}
func getAcquiredSchemas(ref string) *SchemaRef {
	if acquiredSchemas == nil {
		return nil
	}

	schema := acquiredSchemas[ref]
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

	tName := t.Name()
	if tName == "" {
		var genPartOfName string

		if genPart, err := uuid.NewUUID(); err == nil {
			genPartOfName = genPart.String()
		} else {
			genPartOfName = strconv.Itoa(rand.Int())
		}

		tName = "generated-" + genPartOfName
	}

	ref := strings.ReplaceAll(strings.ReplaceAll(t.PkgPath(), "/", "_"), ".", "_") + tName
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

	schema := getDefaultSchema(t)
	if t.Kind() == reflect.Struct {
		schema.Title = tName
		schema.Type = &Types{"object"}
		for i := range t.NumField() {
			field := t.Field(i)

			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}

			fieldType := field.Type
			fieldTypeName := fieldType.Name()
			fieldTypePkgPath := fieldType.PkgPath()
			fieldKind := fieldType.Kind()
			isNullable := false
			if fieldKind == reflect.Pointer {
				fieldKind = fieldType.Elem().Kind()
				fieldType = fieldType.Elem()
				isNullable = true
			}

			// for debugging purposes:
			// log.Println(fieldType, fieldTypeName, fieldKind)

			// create schema for the field. First handle special cases!
			var result *SchemaRef = nil
			switch {
			// skip channels and functions
			case fieldKind == reflect.Func, fieldKind == reflect.Chan:
				continue

			// handle time type
			case fieldKind == reflect.Struct && fieldType == timeType:
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"string"},
					Format: "date-time",
				}}

			// handle file uploads
			case fieldKind == reflect.Struct && fieldTypeName == "FileHeader" && fieldTypePkgPath == "mime/multipart":
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"string"},
					Format: "binary",
				}}

			// handle uuid.UUID
			case fieldKind == reflect.Array && fieldTypeName == "UUID" && fieldType.Elem().Kind() == reflect.Uint8:
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"string"},
					Format: "uuid",
				}}

			// handle uuid.NullUUID
			case fieldKind == reflect.Struct && fieldTypeName == "NullUUID" && isNullType(fieldType, "UUID"):
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"string"},
					Format: "uuid",
				}}

			// handle sql.NullBool
			case fieldKind == reflect.Struct && fieldTypeName == "NullBool" && isNullType(fieldType, "Bool"):
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type: &Types{"boolean"},
				}}

			// handle sql.NullByte
			case fieldKind == reflect.Struct && fieldTypeName == "NullByte" && isNullType(fieldType, "Byte"):
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"string"},
					Format: "byte",
				}}

			// handle sql.NullInt16
			case fieldKind == reflect.Struct && fieldTypeName == "NullInt16" && isNullType(fieldType, "Int16"):
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type: &Types{"integer"},
					Min:  &minInt16,
					Max:  &maxInt16,
				}}

			// handle sql.NullInt32
			case fieldKind == reflect.Struct && fieldTypeName == "NullInt32" && isNullType(fieldType, "Int32"):
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"integer"},
					Format: "int32",
					Min:    &minInt32,
					Max:    &maxInt32,
				}}

			// handle sql.NullInt64
			case fieldKind == reflect.Struct && fieldTypeName == "NullInt64" && isNullType(fieldType, "Int64"):
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"integer"},
					Format: "int64",
					Min:    &minInt64,
					Max:    &maxInt64,
				}}

			// handle sql.NullFloat64
			case fieldKind == reflect.Struct && fieldTypeName == "NullFloat64" && isNullType(fieldType, "Float64"):
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"number"},
					Format: "double",
				}}

			// handle sql.NullTime
			case fieldKind == reflect.Struct && fieldTypeName == "NullTime" && isNullType(fieldType, "Time"): // todo: we could also check whether the Time field is of time.Time type
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type:   &Types{"string"},
					Format: "date-time",
				}}

			// handle sql.NullString
			case fieldKind == reflect.Struct && fieldTypeName == "NullString" && isNullType(fieldType, "String"):
				isNullable = true
				result = &SchemaRef{Value: &Schema{
					Type: &Types{"string"},
				}}

			// handle bytes
			case fieldKind == reflect.Slice && fieldType.Elem().Kind() == reflect.Uint8:
				if fieldType == rawMessageType {
					result = &SchemaRef{Value: &Schema{}}
				} else {
					result = &SchemaRef{Value: &Schema{
						Type:   &Types{"string"},
						Format: "byte",
					}}
				}

			// handle map[string]object
			case fieldKind == reflect.Map && fieldType.Key().Kind() == reflect.String:
				valueSchema := generateSchema(fieldType.Elem())
				has := true
				result = &SchemaRef{Value: &Schema{
					Type: &Types{"object"},
					AdditionalProperties: AdditionalProperties{
						Has:    &has,
						Schema: valueSchema,
					},
				}}

			// handle general structs
			case fieldKind == reflect.Struct:
				result = generateSchema(fieldType)

			// handle general slices / arrays
			case fieldKind == reflect.Slice, fieldKind == reflect.Array:
				result = &SchemaRef{Value: &Schema{
					Type:  &Types{"array"},
					Items: generateSchema(fieldType.Elem()),
				}}

			// handle general maps
			case fieldKind == reflect.Map:
				result = &SchemaRef{Value: &Schema{
					Type: &Types{"object"},
				}}

			// generated default schema for non-special types (string/int/etc)
			default:
				result = &SchemaRef{
					Value: getDefaultSchema(fieldType),
				}
			}
			result.Value.Nullable = isNullable

			// handle json tag
			fieldName := field.Name
			jsonTagOptions := strings.Split(jsonTag, ",")
			if len(jsonTagOptions) > 0 && jsonTagOptions[0] != "" {
				fieldName = jsonTagOptions[0]
			}
			for i := 1; i < len(jsonTagOptions); i++ {
				option := jsonTagOptions[i]
				switch option {
				case "string":
					result.Value.Type = &Types{"string"}
				case "omitempty":
					result.Value.Description += " omitempty "
				case "omitzero":
					result.Value.Description += " omitzero "
				}
			}

			// handle validate tag
			validateTag := field.Tag.Get("validate")
			validateTagOptions := strings.Split(validateTag, ",")
			for _, validation := range validateTagOptions {
				switch {
				case validation == "required":
					schema.Required = append(schema.Required, fieldName)
					result.Value.Nullable = false
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
				case strings.HasPrefix(validation, "omitnil"):
					result.Value.Description += " omitnil "
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
			result.Value.Title = fieldName
			result.Value.Description = strings.ReplaceAll(result.Value.Description, "  ", "")

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

func getDefaultSchema(t reflect.Type) *Schema {
	schema := Schema{
		Properties: make(Schemas),
		Required:   []string{},
	}
	switch t.Kind() {
	case reflect.Bool:
		schema.Type = &Types{"boolean"}

	case reflect.Int:
		schema.Type = &Types{"integer"}
		schema.Min = &minInt
		schema.Max = &maxInt
	case reflect.Int8:
		schema.Type = &Types{"integer"}
		schema.Min = &minInt8
		schema.Max = &maxInt8
	case reflect.Int16:
		schema.Type = &Types{"integer"}
		schema.Min = &minInt16
		schema.Max = &maxInt16
	case reflect.Int32:
		schema.Type = &Types{"integer"}
		schema.Format = "int32"
		schema.Min = &minInt32
		schema.Max = &maxInt32
	case reflect.Int64:
		schema.Type = &Types{"integer"}
		schema.Format = "int64"
		schema.Min = &minInt64
		schema.Max = &maxInt64
	case reflect.Uint:
		schema.Type = &Types{"integer"}
		schema.Min = &zeroInt
		schema.Max = &maxUint
	case reflect.Uint8:
		schema.Type = &Types{"integer"}
		schema.Min = &zeroInt
		schema.Max = &maxUint8
	case reflect.Uint16:
		schema.Type = &Types{"integer"}
		schema.Min = &zeroInt
		schema.Max = &maxUint16
	case reflect.Uint32:
		schema.Type = &Types{"integer"}
		schema.Min = &zeroInt
		schema.Max = &maxUint32
	case reflect.Uint64:
		schema.Type = &Types{"integer"}
		schema.Min = &zeroInt
		schema.Max = &maxUint64

	case reflect.Float32:
		schema.Type = &Types{"number"}
		schema.Format = "float"
	case reflect.Float64:
		schema.Type = &Types{"number"}
		schema.Format = "double"

	case reflect.String:
		schema.Type = &Types{"string"}

	case reflect.Array:
		if t.Name() == "UUID" && t.Elem().Kind() == reflect.Uint8 {
			schema.Type = &Types{"string"}
			schema.Format = "uuid"
		}
	}

	return &schema
}

func isNullType(fieldType reflect.Type, uniqueFieldName string) bool {
	_, has_valid := fieldType.FieldByName("Valid")
	if !has_valid {
		return false
	}
	_, ok_unique := fieldType.FieldByName(uniqueFieldName)
	return ok_unique
}
