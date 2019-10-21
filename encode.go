package yaml

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/printer"
	"github.com/goccy/go-yaml/token"
	"golang.org/x/xerrors"
)

const (
	// DefaultIndentSpaces default number of space for indent
	DefaultIndentSpaces = 2
)

// Encoder writes YAML values to an output stream.
type Encoder struct {
	writer             io.Writer
	opts               []EncodeOption
	indent             int
	anchorPtrToNameMap map[uintptr]string

	line        int
	column      int
	offset      int
	indentNum   int
	indentLevel int
}

// NewEncoder returns a new encoder that writes to w.
// The Encoder should be closed after use to flush all data to w.
func NewEncoder(w io.Writer, opts ...EncodeOption) *Encoder {
	return &Encoder{
		writer:             w,
		opts:               opts,
		indent:             DefaultIndentSpaces,
		anchorPtrToNameMap: map[uintptr]string{},
		line:               1,
		column:             1,
		offset:             0,
	}
}

// Close closes the encoder by writing any remaining data.
// It does not write a stream terminating string "...".
func (e *Encoder) Close() error {
	return nil
}

// Encode writes the YAML encoding of v to the stream.
// If multiple items are encoded to the stream,
// the second and subsequent document will be preceded with a "---" document separator,
// but the first will not.
//
// See the documentation for Marshal for details about the conversion of Go values to YAML.
func (e *Encoder) Encode(v interface{}) error {
	node, err := e.encodeValue(reflect.ValueOf(v), 1)
	if err != nil {
		return xerrors.Errorf("failed to encode value: %w", err)
	}
	var p printer.Printer
	e.writer.Write(p.PrintNode(node))
	return nil
}

func (e *Encoder) encodeValue(v reflect.Value, column int) (ast.Node, error) {
	switch v.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.encodeInt(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeUint(v.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return e.encodeFloat(v.Float()), nil
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return e.encodeNil(), nil
		}
		return e.encodeValue(v.Elem(), column)
	case reflect.String:
		return e.encodeString(v.String(), column), nil
	case reflect.Bool:
		return e.encodeBool(v.Bool()), nil
	case reflect.Slice:
		return e.encodeSlice(v), nil
	case reflect.Struct:
		return e.encodeStruct(v, column)
	case reflect.Map:
		return e.encodeMap(v, column), nil
	default:
		fmt.Printf("valueType = [%s]\n", v.Type().String())
	}
	return nil, nil
}

func (e *Encoder) pos(column int) *token.Position {
	return &token.Position{
		Line:        e.line,
		Column:      column,
		Offset:      e.offset,
		IndentNum:   e.indentNum,
		IndentLevel: e.indentLevel,
	}
}

func (e *Encoder) encodeNil() ast.Node {
	value := "null"
	return ast.Null(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeInt(v int64) ast.Node {
	value := fmt.Sprint(v)
	return ast.Integer(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeUint(v uint64) ast.Node {
	value := fmt.Sprint(v)
	return ast.Integer(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeFloat(v float64) ast.Node {
	if v == math.Inf(0) {
		value := ".inf"
		return ast.Infinity(token.New(value, value, e.pos(e.column)))
	} else if v == math.Inf(-1) {
		value := "-.inf"
		return ast.Infinity(token.New(value, value, e.pos(e.column)))
	} else if math.IsNaN(v) {
		value := ".nan"
		return ast.Nan(token.New(value, value, e.pos(e.column)))
	}
	value := fmt.Sprintf("%f", v)
	fvalue := strings.Split(value, ".")
	if len(fvalue) > 1 {
		precision := fvalue[1]
		precisionNum := 1
		for i := len(precision) - 1; i >= 0; i-- {
			if precision[i] != '0' {
				precisionNum = i + 1
				break
			}
		}
		value = strconv.FormatFloat(v, 'f', precisionNum, 64)
	}
	return ast.Float(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeString(v string, column int) ast.Node {
	if token.IsNeedQuoted(v) {
		v = strconv.Quote(v)
	}
	return ast.String(token.New(v, v, e.pos(column)))
}

func (e *Encoder) encodeBool(v bool) ast.Node {
	value := fmt.Sprint(v)
	return ast.Bool(token.New(value, value, e.pos(e.column)))
}

func (e *Encoder) encodeSlice(value reflect.Value) ast.Node {
	sequence := &ast.SequenceNode{
		Start:  token.New("-", "-", e.pos(e.column)),
		Values: []ast.Node{},
	}
	for i := 0; i < value.Len(); i++ {
		node, err := e.encodeValue(value.Index(i), e.column)
		if err != nil {
			panic(err)
			return nil
		}
		sequence.Values = append(sequence.Values, node)
	}
	return sequence
}

func (e *Encoder) encodeMap(value reflect.Value, column int) ast.Node {
	node := &ast.MappingCollectionNode{
		Start:  token.New("", "", e.pos(column)),
		Values: []ast.Node{},
	}
	keys := []string{}
	for _, k := range value.MapKeys() {
		keys = append(keys, k.Interface().(string))
	}
	sort.Strings(keys)
	for _, key := range keys {
		k := reflect.ValueOf(key)
		v := value.MapIndex(k)
		value, err := e.encodeValue(v, column)
		if err != nil {
			return nil
		}
		if c, ok := value.(*ast.MappingCollectionNode); ok {
			for _, value := range c.Values {
				if mvnode, ok := value.(*ast.MappingValueNode); ok {
					mvnode.Key.GetToken().Position.Column += e.indent
				}
			}
		}
		node.Values = append(node.Values, &ast.MappingValueNode{
			Key:   e.encodeString(k.Interface().(string), column),
			Value: value,
		})
	}
	return node
}

// IsZeroer is used to check whether an object is zero to determine
// whether it should be omitted when marshaling with the omitempty flag.
// One notable implementation is time.Time.
type IsZeroer interface {
	IsZero() bool
}

func (e *Encoder) isZeroValue(v reflect.Value) bool {
	kind := v.Kind()
	if z, ok := v.Interface().(IsZeroer); ok {
		if (kind == reflect.Ptr || kind == reflect.Interface) && v.IsNil() {
			return true
		}
		return z.IsZero()
	}
	switch kind {
	case reflect.String:
		return len(v.String()) == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Struct:
		vt := v.Type()
		for i := v.NumField() - 1; i >= 0; i-- {
			if vt.Field(i).PkgPath != "" {
				continue // private field
			}
			if !e.isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	}
	return false
}

func (e *Encoder) encodeStruct(value reflect.Value, column int) (ast.Node, error) {
	node := &ast.MappingCollectionNode{
		Start:  token.New("", "", e.pos(column)),
		Values: []ast.Node{},
	}
	structType := value.Type()
	structFieldMap, err := structFieldMap(structType)
	if err != nil {
		return nil, xerrors.Errorf("failed to get struct field map: %w", err)
	}
	for i := 0; i < value.NumField(); i++ {
		field := structType.Field(i)
		if isIgnoredStructField(field) {
			continue
		}
		fieldValue := value.FieldByName(field.Name)
		structField := structFieldMap[field.Name]
		if structField.IsOmitEmpty && e.isZeroValue(fieldValue) {
			// omit encoding
			continue
		}
		value, err := e.encodeValue(fieldValue, column)
		if err != nil {
			return nil, xerrors.Errorf("failed to encode value: %w", err)
		}
		if c, ok := value.(*ast.MappingCollectionNode); ok {
			for _, value := range c.Values {
				if mvnode, ok := value.(*ast.MappingValueNode); ok {
					mvnode.Key.GetToken().Position.Column += e.indent
				}
			}
		}
		key := e.encodeString(structField.RenderName, column)
		if structField.AnchorName != "" {
			anchorName := structField.AnchorName
			if fieldValue.Kind() == reflect.Ptr {
				e.anchorPtrToNameMap[fieldValue.Pointer()] = anchorName
			}
			value = &ast.AnchorNode{
				Start: token.New("&", "&", e.pos(column)),
				Name:  ast.String(token.New(anchorName, anchorName, e.pos(column))),
				Value: value,
			}
		} else if structField.IsAutoAnchor {
			anchorName := structField.RenderName
			if fieldValue.Kind() == reflect.Ptr {
				e.anchorPtrToNameMap[fieldValue.Pointer()] = anchorName
			}
			value = &ast.AnchorNode{
				Start: token.New("&", "&", e.pos(column)),
				Name:  ast.String(token.New(anchorName, anchorName, e.pos(column))),
				Value: value,
			}
		} else if structField.IsAutoAlias {
			if fieldValue.Kind() != reflect.Ptr {
				return nil, xerrors.Errorf(
					"%s in struct is not pointer type. but required automatically alias detection",
					structField.FieldName,
				)
			}
			anchorName := e.anchorPtrToNameMap[fieldValue.Pointer()]
			if anchorName == "" {
				return nil, xerrors.Errorf(
					"cannot find anchor name from pointer address for automatically alias detection",
				)
			}
			aliasName := anchorName
			value = &ast.AliasNode{
				Start: token.New("*", "*", e.pos(column)),
				Value: ast.String(token.New(aliasName, aliasName, e.pos(column))),
			}
			if structField.IsInline {
				// if both used alias and inline, output `<<: *alias`
				key = ast.MergeKey(token.New("<<", "<<", e.pos(column)))
			}
		} else if structField.AliasName != "" {
			aliasName := structField.AliasName
			value = &ast.AliasNode{
				Start: token.New("*", "*", e.pos(column)),
				Value: ast.String(token.New(aliasName, aliasName, e.pos(column))),
			}
			if structField.IsInline {
				// if both used alias and inline, output `<<: *alias`
				key = ast.MergeKey(token.New("<<", "<<", e.pos(column)))
			}
		}
		node.Values = append(node.Values, &ast.MappingValueNode{
			Key:   key,
			Value: value,
		})
	}
	return node, nil
}
