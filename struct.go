package yaml

import (
	"reflect"
	"strings"

	"golang.org/x/xerrors"
)

const (
	// StructTagName tag keyword for Marshal/Unmarshal
	StructTagName = "yaml"
)

// StructField information for each the field in structure
type StructField struct {
	FieldName    string
	RenderName   string
	AnchorName   string
	AliasName    string
	IsAutoAnchor bool
	IsAutoAlias  bool
	IsOmitEmpty  bool
	IsFlow       bool
	IsInline     bool
}

func structField(field reflect.StructField) *StructField {
	tag := field.Tag.Get(StructTagName)
	fieldName := strings.ToLower(field.Name)
	options := strings.Split(tag, ",")
	if len(options) > 0 {
		if options[0] != "" {
			fieldName = options[0]
		}
	}
	structField := &StructField{
		FieldName:  field.Name,
		RenderName: fieldName,
	}
	if len(options) > 1 {
		for _, opt := range options[1:] {
			switch {
			case opt == "omitempty":
				structField.IsOmitEmpty = true
			case opt == "flow":
				structField.IsFlow = true
			case opt == "inline":
				structField.IsInline = true
			case strings.HasPrefix(opt, "anchor"):
				anchor := strings.Split(opt, "=")
				if len(anchor) > 1 {
					structField.AnchorName = anchor[1]
				} else {
					structField.IsAutoAnchor = true
				}
			case strings.HasPrefix(opt, "alias"):
				alias := strings.Split(opt, "=")
				if len(alias) > 1 {
					structField.AliasName = alias[1]
				} else {
					structField.IsAutoAlias = true
				}
			default:
			}
		}
	}
	return structField
}

func isIgnoredStructField(field reflect.StructField) bool {
	if field.PkgPath != "" && !field.Anonymous {
		// private field
		return true
	}
	tag := field.Tag.Get(StructTagName)
	if tag == "-" {
		return true
	}
	return false
}

type StructFieldMap map[string]*StructField

func (m StructFieldMap) isIncludedRenderName(name string) bool {
	for _, v := range m {
		if v.RenderName == name {
			return true
		}
	}
	return false
}

func structFieldMap(structType reflect.Type) (StructFieldMap, error) {
	structFieldMap := StructFieldMap{}
	renderNameMap := map[string]struct{}{}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if isIgnoredStructField(field) {
			continue
		}
		structField := structField(field)
		if _, exists := renderNameMap[structField.RenderName]; exists {
			return nil, xerrors.Errorf("duplicated struct field name %s", structField.RenderName)
		}
		structFieldMap[structField.FieldName] = structField
		renderNameMap[structField.RenderName] = struct{}{}
	}
	return structFieldMap, nil
}
