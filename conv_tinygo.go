//go:build tinygo

package yaml

import (
	"reflect"
)

func convertibleTo(src reflect.Value, typ reflect.Type) bool {
	srck, typk := src.Kind(), typ.Kind()
	if srck == typk {
		return true
	}
	switch srck {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch typk {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return true
		case reflect.Float32, reflect.Float64:
			return true
		case reflect.String:
			return true
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch typk {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return true
		case reflect.Float32, reflect.Float64:
			return true
		case reflect.String:
			return true
		}

	case reflect.Float32, reflect.Float64:
		switch typk {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return true
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return true
		case reflect.Float32, reflect.Float64:
			return true
		}

	case reflect.Slice:
		if typk == reflect.String /*&& !src.Type().Elem().isNamed()*/ {
			switch src.Type().Elem().Kind() {
			case reflect.Uint8, reflect.Int32:
				return true
			}
		}

	case reflect.String:
		switch typk {
		case reflect.Slice:
			switch typ.Elem().Kind() {
			case reflect.Uint8, reflect.Int32:
				return true
			}
			return false
		}
	}

	return false
}
