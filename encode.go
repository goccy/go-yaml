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
	DefaultIndentSpaces = 2
)

// An Encoder writes YAML values to an output stream.
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
	node, err := e.encodeValue(v, 1)
	if err != nil {
		return xerrors.Errorf("failed to encode value: %w", err)
	}
	var p printer.Printer
	e.writer.Write(p.PrintNode(node))
	return nil
}

func (e *Encoder) encodeValue(v interface{}, column int) (ast.Node, error) {
	if v == nil {
		return e.encodeNil(), nil
	}
	value := reflect.ValueOf(v)
	valueType := reflect.TypeOf(v)
	switch valueType.Kind() {
	case reflect.Int:
		return e.encodeInt(int64(v.(int))), nil
	case reflect.Int8:
		return e.encodeInt(int64(v.(int8))), nil
	case reflect.Int16:
		return e.encodeInt(int64(v.(int16))), nil
	case reflect.Int32:
		return e.encodeInt(int64(v.(int32))), nil
	case reflect.Int64:
		return e.encodeInt(v.(int64)), nil
	case reflect.Uint:
		return e.encodeUint(uint64(v.(uint))), nil
	case reflect.Uint8:
		return e.encodeUint(uint64(v.(uint8))), nil
	case reflect.Uint16:
		return e.encodeUint(uint64(v.(uint16))), nil
	case reflect.Uint32:
		return e.encodeUint(uint64(v.(uint32))), nil
	case reflect.Uint64:
		return e.encodeUint(v.(uint64)), nil
	case reflect.Float32:
		return e.encodeFloat(float64(v.(float32))), nil
	case reflect.Float64:
		return e.encodeFloat(v.(float64)), nil
	case reflect.Ptr:
		if value.IsNil() {
			return e.encodeNil(), nil
		}
		return e.encodeValue(value.Elem().Interface(), column)
	case reflect.String:
		return e.encodeString(v.(string), column), nil
	case reflect.Bool:
		return e.encodeBool(v.(bool)), nil
	case reflect.Slice:
		return e.encodeSlice(value), nil
	case reflect.Struct:
		return e.encodeStruct(value, column), nil
	case reflect.Map:
		return e.encodeMap(value, column), nil
	default:
		fmt.Printf("valueType = [%s]\n", valueType.String())
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
		node, err := e.encodeValue(value.Index(i).Interface(), e.column)
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
		value, err := e.encodeValue(v.Interface(), column)
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

func (e *Encoder) encodeStruct(value reflect.Value, column int) ast.Node {
	node := &ast.MappingCollectionNode{
		Start:  token.New("", "", e.pos(column)),
		Values: []ast.Node{},
	}
	structType := value.Type()
	structFieldMap, err := structFieldMap(structType)
	if err != nil {
		return nil
	}
	for i := 0; i < value.NumField(); i++ {
		field := structType.Field(i)
		if isIgnoredStructField(field) {
			continue
		}
		fieldValue := value.FieldByName(field.Name)
		structField := structFieldMap[field.Name]
		value, err := e.encodeValue(fieldValue.Interface(), column)
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
				// TODO: error handling
			}
			anchorName := e.anchorPtrToNameMap[fieldValue.Pointer()]
			if anchorName == "" {
				// TODO: error handling
			}
			aliasName := anchorName
			value = &ast.AliasNode{
				Start: token.New("*", "*", e.pos(column)),
				Value: ast.String(token.New(aliasName, aliasName, e.pos(column))),
			}
		} else if structField.AliasName != "" {
			aliasName := structField.AliasName
			value = &ast.AliasNode{
				Start: token.New("*", "*", e.pos(column)),
				Value: ast.String(token.New(aliasName, aliasName, e.pos(column))),
			}
		}
		node.Values = append(node.Values, &ast.MappingValueNode{
			Key:   e.encodeString(structField.RenderName, column),
			Value: value,
		})
	}
	return node
}
