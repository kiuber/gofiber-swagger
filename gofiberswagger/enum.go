package gofiberswagger

import "reflect"

type ISwaggerEnum interface {
	EnumValues() []any
}

func implementsSwaggerEnum(t reflect.Type) bool {
	if t == nil {
		return false
	}

	tKind := t.Kind()
	enumInterface := reflect.TypeOf((*ISwaggerEnum)(nil)).Elem()
	return t.Implements(enumInterface) || (tKind != reflect.Pointer && tKind != reflect.UnsafePointer && tKind != reflect.Invalid && reflect.PointerTo(t).Implements(enumInterface))
}

// call implementsSwaggerEnum beforehand!
func getSwaggerEnumValues(t reflect.Type) []any {
	if !implementsSwaggerEnum(t) {
		return []any{}
	}

	var instance ISwaggerEnum
	enumInterface := reflect.TypeOf((*ISwaggerEnum)(nil)).Elem()
	if t.Implements(enumInterface) {
		instance = reflect.New(t).Elem().Interface().(ISwaggerEnum)
	}
	if reflect.PointerTo(t).Implements(enumInterface) {
		instance = reflect.New(t).Interface().(ISwaggerEnum)
	}
	return instance.EnumValues()
}
