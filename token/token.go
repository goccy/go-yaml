package token

import "fmt"

type Character byte

const (
	SequenceEntryCharacter Character = '-'
	MappingKeyCharacter              = '?'
	MappingValueCharacter            = ':'
	CollectEntryCharacter            = ','
	SequenceStartCharacter           = '['
	SequenceEndCharacter             = ']'
	MappingStartCharacter            = '{'
	MappingEndCharacter              = '}'
	CommentCharacter                 = '#'
	AnchorCharacter                  = '&'
	AliasCharacter                   = '*'
	TagCharacter                     = '!'
	LiteralCharacter                 = '|'
	FoldedCharacter                  = '>'
	SingleQuoteCharacter             = '\''
	DoubleQuoteCharacter             = '"'
	DirectiveCharacter               = '%'
	SpaceCharacter                   = ' '
	TabCharacter                     = '\t'
	LineBreakCharacter               = '\n'
)

type Type int

const (
	UnknownType Type = iota
	DocumentHeaderType
	DocumentEndType
	SequenceEntryType
	MappingKeyType
	MappingValueType
	MergeKeyType
	CollectEntryType
	SequenceStartType
	SequenceEndType
	MappingStartType
	MappingEndType
	CommentType
	AnchorType
	AliasType
	TagType
	LiteralType
	FoldedType
	SingleQuoteType
	DoubleQuoteType
	DirectiveType
	SpaceType
	TabType
	NullType
	InfinityType
	NanType
	IntegerType
	FloatType
	StringType
	BoolType
)

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
	case TabType:
		return "Tab"
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

type CharacterType int

const (
	CharacterTypeIndicator CharacterType = iota
	CharacterTypeWhiteSpace
	CharacterTypeMiscellaneous
	CharacterTypeEscaped
)

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

type Indicator int

const (
	NotIndicator                  Indicator = iota
	BlockStructureIndicator                 // '-', '?', ':'
	FlowCollectionIndicator                 // '[', ']', '{', '}', ','
	CommentIndicator                        // '#'
	NodePropertyIndicator                   // '!', '&', '*'
	BlockScalarIndicator                    // '|', '>'
	QuotedScalarIndicator                   // ''', '"'
	DirectiveIndicator                      // '%'
	InvalidUseOfReservedIndicator           // '@', '`'
)

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

type ReservedKeyword string

const (
	Null             ReservedKeyword = "null"
	False                            = "false"
	True                             = "true"
	Infinity                         = ".inf"
	NegativeInfinity                 = "-.inf"
	Nan                              = ".nan"
)

var (
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

type ReservedTagKeyword string

const (
	IntegerTag    ReservedTagKeyword = "!!int"
	FloatTag                         = "!!float"
	NullTag                          = "!!null"
	SequenceTag                      = "!!seq"
	MappingTag                       = "!!map"
	StringTag                        = "!!str"
	BinaryTag                        = "!!binary"
	OrderedMapTag                    = "!!omap"
	SetTag                           = "!!set"
)

var (
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

type Position struct {
	Line        int
	Column      int
	Offset      int
	IndentNum   int
	IndentLevel int
}

func (p *Position) String() string {
	return fmt.Sprintf("[level:%d,line:%d,column:%d,offset:%d]", p.IndentLevel, p.Line, p.Column, p.Offset)
}

func (p *Position) IndentString() string {
	spaceNum := p.IndentLevel * 2
	space := ""
	for i := 0; i < spaceNum; i++ {
		space += " "
	}
	return space
}

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

func (t *Token) NextType() Type {
	if t.Next != nil {
		return t.Next.Type
	}
	return UnknownType
}

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

func (t *Tokens) Add(tks ...*Token) {
	for _, tk := range tks {
		t.add(tk)
	}
}

func (t Tokens) Dump() {
	for _, tk := range t {
		fmt.Printf("- %+v\n", tk)
	}
}

var (
	SequenceEntry = func(org string, pos *Position) *Token {
		return &Token{
			Type:          SequenceEntryType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     BlockStructureIndicator,
			Value:         string(SequenceEntryCharacter),
			Origin:        org,
			Position:      pos,
		}
	}
	MappingKey = func(pos *Position) *Token {
		return &Token{
			Type:          MappingKeyType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     BlockStructureIndicator,
			Value:         string(MappingKeyCharacter),
			Origin:        string(MappingKeyCharacter),
			Position:      pos,
		}
	}
	MappingValue = func(pos *Position) *Token {
		return &Token{
			Type:          MappingValueType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     BlockStructureIndicator,
			Value:         string(MappingValueCharacter),
			Origin:        string(MappingValueCharacter),
			Position:      pos,
		}
	}
	CollectEntry = func(org string, pos *Position) *Token {
		return &Token{
			Type:          CollectEntryType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     FlowCollectionIndicator,
			Value:         string(CollectEntryCharacter),
			Origin:        org,
			Position:      pos,
		}
	}
	SequenceStart = func(org string, pos *Position) *Token {
		return &Token{
			Type:          SequenceStartType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     FlowCollectionIndicator,
			Value:         string(SequenceStartCharacter),
			Origin:        org,
			Position:      pos,
		}
	}
	SequenceEnd = func(org string, pos *Position) *Token {
		return &Token{
			Type:          SequenceEndType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     FlowCollectionIndicator,
			Value:         string(SequenceEndCharacter),
			Origin:        org,
			Position:      pos,
		}
	}
	MappingStart = func(org string, pos *Position) *Token {
		return &Token{
			Type:          MappingStartType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     FlowCollectionIndicator,
			Value:         string(MappingStartCharacter),
			Origin:        org,
			Position:      pos,
		}
	}
	MappingEnd = func(org string, pos *Position) *Token {
		return &Token{
			Type:          MappingEndType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     FlowCollectionIndicator,
			Value:         string(MappingEndCharacter),
			Origin:        org,
			Position:      pos,
		}
	}
	Comment = func(value string, org string, pos *Position) *Token {
		return &Token{
			Type:          CommentType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     CommentIndicator,
			Value:         value,
			Origin:        org,
			Position:      pos,
		}
	}
	Anchor = func(org string, pos *Position) *Token {
		return &Token{
			Type:          AnchorType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     NodePropertyIndicator,
			Value:         string(AnchorCharacter),
			Origin:        org,
			Position:      pos,
		}
	}
	Alias = func(org string, pos *Position) *Token {
		return &Token{
			Type:          AliasType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     NodePropertyIndicator,
			Value:         string(AliasCharacter),
			Origin:        org,
			Position:      pos,
		}
	}
	Tag = func(value string, org string, pos *Position) *Token {
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
	Literal = func(value string, org string, pos *Position) *Token {
		return &Token{
			Type:          LiteralType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     BlockScalarIndicator,
			Value:         value,
			Origin:        org,
			Position:      pos,
		}
	}
	Folded = func(value string, org string, pos *Position) *Token {
		return &Token{
			Type:          FoldedType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     BlockScalarIndicator,
			Value:         value,
			Origin:        org,
			Position:      pos,
		}
	}
	SingleQuote = func(value string, org string, pos *Position) *Token {
		return &Token{
			Type:          SingleQuoteType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     QuotedScalarIndicator,
			Value:         value,
			Origin:        org,
			Position:      pos,
		}
	}
	DoubleQuote = func(value string, org string, pos *Position) *Token {
		return &Token{
			Type:          DoubleQuoteType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     QuotedScalarIndicator,
			Value:         value,
			Origin:        org,
			Position:      pos,
		}
	}
	Directive = func(pos *Position) *Token {
		return &Token{
			Type:          DirectiveType,
			CharacterType: CharacterTypeIndicator,
			Indicator:     DirectiveIndicator,
			Value:         string(DirectiveCharacter),
			Origin:        string(DirectiveCharacter),
			Position:      pos,
		}
	}
	Space = func(pos *Position) *Token {
		return &Token{
			Type:          SpaceType,
			CharacterType: CharacterTypeWhiteSpace,
			Indicator:     NotIndicator,
			Value:         string(SpaceCharacter),
			Origin:        string(SpaceCharacter),
			Position:      pos,
		}
	}
	Tab = func(pos *Position) *Token {
		return &Token{
			Type:          TabType,
			CharacterType: CharacterTypeWhiteSpace,
			Indicator:     NotIndicator,
			Value:         string(TabCharacter),
			Origin:        string(TabCharacter),
			Position:      pos,
		}
	}
	MergeKey = func(pos *Position) *Token {
		return &Token{
			Type:          MergeKeyType,
			CharacterType: CharacterTypeMiscellaneous,
			Indicator:     NotIndicator,
			Value:         "<<",
			Origin:        "<<",
			Position:      pos,
		}
	}
	DocumentHeader = func(pos *Position) *Token {
		return &Token{
			Type:          DocumentHeaderType,
			CharacterType: CharacterTypeMiscellaneous,
			Indicator:     NotIndicator,
			Value:         "---",
			Origin:        "---",
			Position:      pos,
		}
	}
	DocumentEnd = func(pos *Position) *Token {
		return &Token{
			Type:          DocumentEndType,
			CharacterType: CharacterTypeMiscellaneous,
			Indicator:     NotIndicator,
			Value:         "...",
			Origin:        "...",
			Position:      pos,
		}
	}
)
