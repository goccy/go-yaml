package token

import (
	"fmt"
	"strings"
)

// Character type for character
type Character byte

const (
	// SequenceEntryCharacter character for sequence entry
	SequenceEntryCharacter Character = '-'
	// MappingKeyCharacter character for mapping key
	MappingKeyCharacter = '?'
	// MappingValueCharacter character for mapping value
	MappingValueCharacter = ':'
	// CollectEntryCharacter character for collect entry
	CollectEntryCharacter = ','
	// SequenceStartCharacter character for sequence start
	SequenceStartCharacter = '['
	// SequenceEndCharacter character for sequence end
	SequenceEndCharacter = ']'
	// MappingStartCharacter character for mapping start
	MappingStartCharacter = '{'
	// MappingEndCharacter character for mapping end
	MappingEndCharacter = '}'
	// CommentCharacter character for comment
	CommentCharacter = '#'
	// AnchorCharacter character for anchor
	AnchorCharacter = '&'
	// AliasCharacter character for alias
	AliasCharacter = '*'
	// TagCharacter character for tag
	TagCharacter = '!'
	// LiteralCharacter character for literal
	LiteralCharacter = '|'
	// FoldedCharacter character for folded
	FoldedCharacter = '>'
	// SingleQuoteCharacter character for single quote
	SingleQuoteCharacter = '\''
	// DoubleQuoteCharacter character for double quote
	DoubleQuoteCharacter = '"'
	// DirectiveCharacter character for directive
	DirectiveCharacter = '%'
	// SpaceCharacter character for space
	SpaceCharacter = ' '
	// LineBreakCharacter character for line break
	LineBreakCharacter = '\n'
)

// Type type identifier for token
type Type int

const (
	// UnknownType reserve for invalid type
	UnknownType Type = iota
	// DocumentHeaderType type for DocumentHeader token
	DocumentHeaderType
	// DocumentEndType type for DocumentEnd token
	DocumentEndType
	// SequenceEntryType type for SequenceEntry token
	SequenceEntryType
	// MappingKeyType type for MappingKey token
	MappingKeyType
	// MappingValueType type for MappingValue token
	MappingValueType
	// MergeKeyType type for MergeKey token
	MergeKeyType
	// CollectEntryType type for CollectEntry token
	CollectEntryType
	// SequenceStartType type for SequenceStart token
	SequenceStartType
	// SequenceEndType type for SequenceEnd token
	SequenceEndType
	// MappingStartType type for MappingStart token
	MappingStartType
	// MappingEndType type for MappingEnd token
	MappingEndType
	// CommentType type for Comment token
	CommentType
	// AnchorType type for Anchor token
	AnchorType
	// AliasType type for Alias token
	AliasType
	// TagType type for Tag token
	TagType
	// LiteralType type for Literal token
	LiteralType
	// FoldedType type for Folded token
	FoldedType
	// SingleQuoteType type for SingleQuote token
	SingleQuoteType
	// DoubleQuoteType type for DoubleQuote token
	DoubleQuoteType
	// DirectiveType type for Directive token
	DirectiveType
	// SpaceType type for Space token
	SpaceType
	// NullType type for Null token
	NullType
	// InfinityType type for Infinity token
	InfinityType
	// NanType type for Nan token
	NanType
	// IntegerType type for Integer token
	IntegerType
	// FloatType type for Float token
	FloatType
	// StringType type for String token
	StringType
	// BoolType type for Bool token
	BoolType
)

// String type identifier to text
func (t Type) String() string {
	switch t {
	case UnknownType:
		return "Unknown"
	case DocumentHeaderType:
		return "DocumentHeader"
	case DocumentEndType:
		return "DocumentEnd"
	case SequenceEntryType:
		return "SequenceEntry"
	case MappingKeyType:
		return "MappingKey"
	case MappingValueType:
		return "MappingValue"
	case MergeKeyType:
		return "MergeKey"
	case CollectEntryType:
		return "CollectEntry"
	case SequenceStartType:
		return "SequenceStart"
	case SequenceEndType:
		return "SequenceEnd"
	case MappingStartType:
		return "MappingStart"
	case MappingEndType:
		return "MappingEnd"
	case CommentType:
		return "Comment"
	case AnchorType:
		return "Anchor"
	case AliasType:
		return "Alias"
	case TagType:
		return "Tag"
	case LiteralType:
		return "Literal"
	case FoldedType:
		return "Folded"
	case SingleQuoteType:
		return "SingleQuote"
	case DoubleQuoteType:
		return "DoubleQuote"
	case DirectiveType:
		return "Directive"
	case SpaceType:
		return "Space"
	case StringType:
		return "String"
	case BoolType:
		return "Bool"
	case IntegerType:
		return "Integer"
	case FloatType:
		return "Float"
	case NullType:
		return "Null"
	case InfinityType:
		return "Infinity"
	case NanType:
		return "Nan"
	}
	return ""
}

// CharacterType type for character category
type CharacterType int

const (
	// CharacterTypeIndicator type of indicator character
	CharacterTypeIndicator CharacterType = iota
	// CharacterTypeWhiteSpace type of white space character
	CharacterTypeWhiteSpace
	// CharacterTypeMiscellaneous type of miscellaneous character
	CharacterTypeMiscellaneous
	// CharacterTypeEscaped type of escaped character
	CharacterTypeEscaped
)

// String character type identifier to text
func (c CharacterType) String() string {
	switch c {
	case CharacterTypeIndicator:
		return "Indicator"
	case CharacterTypeWhiteSpace:
		return "WhiteSpcae"
	case CharacterTypeMiscellaneous:
		return "Miscellaneous"
	case CharacterTypeEscaped:
		return "Escaped"
	}
	return ""
}

// Indicator type for indicator
type Indicator int

const (
	// NotIndicator not indicator
	NotIndicator Indicator = iota
	// BlockStructureIndicator indicator for block structure ( '-', '?', ':' )
	BlockStructureIndicator
	// FlowCollectionIndicator indicator for flow collection ( '[', ']', '{', '}', ',' )
	FlowCollectionIndicator
	// CommentIndicator indicator for comment ( '#' )
	CommentIndicator
	// NodePropertyIndicator indicator for node property ( '!', '&', '*' )
	NodePropertyIndicator
	// BlockScalarIndicator indicator for block scalar ( '|', '>' )
	BlockScalarIndicator
	// QuotedScalarIndicator indicator for quoted scalar ( ''', '"' )
	QuotedScalarIndicator
	// DirectiveIndicator indicator for directive ( '%' )
	DirectiveIndicator
	// InvalidUseOfReservedIndicator indicator for invalid use of reserved keyword ( '@', '`' )
	InvalidUseOfReservedIndicator
)

// String indicator to text
func (i Indicator) String() string {
	switch i {
	case NotIndicator:
		return "NotIndicator"
	case BlockStructureIndicator:
		return "BlockStructure"
	case FlowCollectionIndicator:
		return "FlowCollection"
	case CommentIndicator:
		return "Comment"
	case NodePropertyIndicator:
		return "NodeProperty"
	case BlockScalarIndicator:
		return "BlockScalar"
	case QuotedScalarIndicator:
		return "QuotedScalar"
	case DirectiveIndicator:
		return "Directive"
	case InvalidUseOfReservedIndicator:
		return "InvalidUseOfReserved"
	}
	return ""
}

// ReservedKeyword type of reserved keyword
type ReservedKeyword string

const (
	// Null `null` keyword
	Null ReservedKeyword = "null"
	// False `false` keyword
	False = "false"
	// True `true` keyword
	True = "true"
	// Infinity `.inf` keyword
	Infinity = ".inf"
	// NegativeInfinity `-.inf` keyword
	NegativeInfinity = "-.inf"
	// Nan `.nan` keyword
	Nan = ".nan"
)

var (
	// ReservedKeywordMap map for reserved keywords
	ReservedKeywordMap = map[ReservedKeyword]func(string, string, *Position) *Token{
		Null: func(value string, org string, pos *Position) *Token {
			return &Token{
				Type:          NullType,
				CharacterType: CharacterTypeMiscellaneous,
				Indicator:     NotIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		False: func(value string, org string, pos *Position) *Token {
			return &Token{
				Type:          BoolType,
				CharacterType: CharacterTypeMiscellaneous,
				Indicator:     NotIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		True: func(value string, org string, pos *Position) *Token {
			return &Token{
				Type:          BoolType,
				CharacterType: CharacterTypeMiscellaneous,
				Indicator:     NotIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		Infinity: func(value string, org string, pos *Position) *Token {
			return &Token{
				Type:          InfinityType,
				CharacterType: CharacterTypeMiscellaneous,
				Indicator:     NotIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		NegativeInfinity: func(value string, org string, pos *Position) *Token {
			return &Token{
				Type:          InfinityType,
				CharacterType: CharacterTypeMiscellaneous,
				Indicator:     NotIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		Nan: func(value string, org string, pos *Position) *Token {
			return &Token{
				Type:          NanType,
				CharacterType: CharacterTypeMiscellaneous,
				Indicator:     NotIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
	}
)

// ReservedTagKeyword type of reserved tag keyword
type ReservedTagKeyword string

const (
	// IntegerTag `!!int` tag
	IntegerTag ReservedTagKeyword = "!!int"
	// FloatTag `!!float` tag
	FloatTag = "!!float"
	// NullTag `!!null` tag
	NullTag = "!!null"
	// SequenceTag `!!seq` tag
	SequenceTag = "!!seq"
	// MappingTag `!!map` tag
	MappingTag = "!!map"
	// StringTag `!!str` tag
	StringTag = "!!str"
	// BinaryTag `!!binary` tag
	BinaryTag = "!!binary"
	// OrderedMapTag `!!omap` tag
	OrderedMapTag = "!!omap"
	// SetTag `!!set` tag
	SetTag = "!!set"
)

var (
	// ReservedTagKeywordMap map for reserved tag keywords
	ReservedTagKeywordMap = map[ReservedTagKeyword]func(string, string, *Position) *Token{
		IntegerTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		FloatTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		NullTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		SequenceTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		MappingTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		StringTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		BinaryTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		OrderedMapTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
		SetTag: func(value, org string, pos *Position) *Token {
			return &Token{
				Type:          TagType,
				CharacterType: CharacterTypeIndicator,
				Indicator:     NodePropertyIndicator,
				Value:         value,
				Origin:        org,
				Position:      pos,
			}
		},
	}
)

func isNumber(str string) (bool, bool) {
	if str == "-" || str == "." {
		return false, false
	}
	isFloat := false
	isMultipleDot := false
	for idx, c := range str {
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		case '.':
			if isFloat {
				isMultipleDot = true
			}
			isFloat = true
			continue
		case '-':
			if idx == 0 {
				continue
			}
		}
		return false, false
	}
	if isMultipleDot {
		return false, false
	}
	return true, isFloat
}

// IsNeedQuoted whether need quote for passed string or not
func IsNeedQuoted(value string) bool {
	if _, exists := ReservedKeywordMap[ReservedKeyword(value)]; exists {
		return true
	}
	if ok, _ := isNumber(value); ok {
		return true
	}
	if strings.IndexByte(value, ':') == 1 {
		return true
	}
	if strings.IndexByte(value, '#') > 0 {
		return true
	}
	for _, c := range value {
		if c == '\\' {
			return true
		}
	}
	return false
}

// New create reserved keyword token or number token and other string token
func New(value string, org string, pos *Position) *Token {
	fn := ReservedKeywordMap[ReservedKeyword(value)]
	if fn != nil {
		return fn(value, org, pos)
	}
	if ok, isFloat := isNumber(value); ok {
		tk := &Token{
			Type:          IntegerType,
			CharacterType: CharacterTypeMiscellaneous,
			Indicator:     NotIndicator,
			Value:         value,
			Origin:        org,
			Position:      pos,
		}
		if isFloat {
			tk.Type = FloatType
		}
		return tk
	}
	return &Token{
		Type:          StringType,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Position type for position in YAML document
type Position struct {
	Line        int
	Column      int
	Offset      int
	IndentNum   int
	IndentLevel int
}

// String position to text
func (p *Position) String() string {
	return fmt.Sprintf("[level:%d,line:%d,column:%d,offset:%d]", p.IndentLevel, p.Line, p.Column, p.Offset)
}

// Token type for token
type Token struct {
	Type          Type
	CharacterType CharacterType
	Indicator     Indicator
	Value         string
	Origin        string
	Position      *Position
	Next          *Token
	Prev          *Token
}

// PreviousType previous token type
func (t *Token) PreviousType() Type {
	if t.Prev != nil {
		return t.Prev.Type
	}
	return UnknownType
}

// NextType next token type
func (t *Token) NextType() Type {
	if t.Next != nil {
		return t.Next.Type
	}
	return UnknownType
}

// Tokens type of token collection
type Tokens []*Token

func (t *Tokens) add(tk *Token) {
	tokens := *t
	if len(tokens) == 0 {
		tokens = append(tokens, tk)
	} else {
		last := tokens[len(tokens)-1]
		last.Next = tk
		tk.Prev = last
		tokens = append(tokens, tk)
	}
	*t = tokens
}

// Add append new some tokens
func (t *Tokens) Add(tks ...*Token) {
	for _, tk := range tks {
		t.add(tk)
	}
}

// Dump dump all token structures for debugging
func (t Tokens) Dump() {
	for _, tk := range t {
		fmt.Printf("- %+v\n", tk)
	}
}

// SequenceEntry create token for SequenceEntry
func SequenceEntry(org string, pos *Position) *Token {
	return &Token{
		Type:          SequenceEntryType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockStructureIndicator,
		Value:         string(SequenceEntryCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// MappingKey create token for MappingKey
func MappingKey(pos *Position) *Token {
	return &Token{
		Type:          MappingKeyType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockStructureIndicator,
		Value:         string(MappingKeyCharacter),
		Origin:        string(MappingKeyCharacter),
		Position:      pos,
	}
}

// MappingValue create token for MappingValue
func MappingValue(pos *Position) *Token {
	return &Token{
		Type:          MappingValueType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockStructureIndicator,
		Value:         string(MappingValueCharacter),
		Origin:        string(MappingValueCharacter),
		Position:      pos,
	}
}

// CollectEntry create token for CollectEntry
func CollectEntry(org string, pos *Position) *Token {
	return &Token{
		Type:          CollectEntryType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(CollectEntryCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// SequenceStart create token for SequenceStart
func SequenceStart(org string, pos *Position) *Token {
	return &Token{
		Type:          SequenceStartType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(SequenceStartCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// SequenceEnd create token for SequenceEnd
func SequenceEnd(org string, pos *Position) *Token {
	return &Token{
		Type:          SequenceEndType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(SequenceEndCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// MappingStart create token for MappingStart
func MappingStart(org string, pos *Position) *Token {
	return &Token{
		Type:          MappingStartType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(MappingStartCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// MappingEnd create token for MappingEnd
func MappingEnd(org string, pos *Position) *Token {
	return &Token{
		Type:          MappingEndType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     FlowCollectionIndicator,
		Value:         string(MappingEndCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// Comment create token for Comment
func Comment(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          CommentType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     CommentIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Anchor create token for Anchor
func Anchor(org string, pos *Position) *Token {
	return &Token{
		Type:          AnchorType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     NodePropertyIndicator,
		Value:         string(AnchorCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// Alias create token for Alias
func Alias(org string, pos *Position) *Token {
	return &Token{
		Type:          AliasType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     NodePropertyIndicator,
		Value:         string(AliasCharacter),
		Origin:        org,
		Position:      pos,
	}
}

// Tag create token for Tag
func Tag(value string, org string, pos *Position) *Token {
	fn := ReservedTagKeywordMap[ReservedTagKeyword(value)]
	if fn != nil {
		return fn(value, org, pos)
	}
	return &Token{
		Type:          TagType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     NodePropertyIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Literal create token for Literal
func Literal(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          LiteralType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockScalarIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Folded create token for Folded
func Folded(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          FoldedType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     BlockScalarIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// SingleQuote create token for SingleQuote
func SingleQuote(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          SingleQuoteType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     QuotedScalarIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// DoubleQuote create token for DoubleQuote
func DoubleQuote(value string, org string, pos *Position) *Token {
	return &Token{
		Type:          DoubleQuoteType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     QuotedScalarIndicator,
		Value:         value,
		Origin:        org,
		Position:      pos,
	}
}

// Directive create token for Directive
func Directive(pos *Position) *Token {
	return &Token{
		Type:          DirectiveType,
		CharacterType: CharacterTypeIndicator,
		Indicator:     DirectiveIndicator,
		Value:         string(DirectiveCharacter),
		Origin:        string(DirectiveCharacter),
		Position:      pos,
	}
}

// Space create token for Space
func Space(pos *Position) *Token {
	return &Token{
		Type:          SpaceType,
		CharacterType: CharacterTypeWhiteSpace,
		Indicator:     NotIndicator,
		Value:         string(SpaceCharacter),
		Origin:        string(SpaceCharacter),
		Position:      pos,
	}
}

// MergeKey create token for MergeKey
func MergeKey(org string, pos *Position) *Token {
	return &Token{
		Type:          MergeKeyType,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         "<<",
		Origin:        org,
		Position:      pos,
	}
}

// DocumentHeader create token for DocumentHeader
func DocumentHeader(pos *Position) *Token {
	return &Token{
		Type:          DocumentHeaderType,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         "---",
		Origin:        "---",
		Position:      pos,
	}
}

// DocumentEnd create token for DocumentEnd
func DocumentEnd(pos *Position) *Token {
	return &Token{
		Type:          DocumentEndType,
		CharacterType: CharacterTypeMiscellaneous,
		Indicator:     NotIndicator,
		Value:         "...",
		Origin:        "...",
		Position:      pos,
	}
}
