package structinfo

import "reflect"

// IsFieldKindSupported checks if a specific reflect kind of the field is supported by the Builder.
func IsFieldKindSupported(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.String, reflect.Bool:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func IsFieldModification(name string, typeKind reflect.Kind) bool {
	return (name == "CreatedAt" || name == "CreatedBy" || name == "ModifiedAt" || name == "ModifiedBy") && typeKind == reflect.Int64
}
