//go:build !tinygo

package yaml

import (
	"reflect"
)

func convertibleTo(src reflect.Value, typ reflect.Type) bool {
	return src.Type().ConvertibleTo(typ)
}
