package lexer_test

import (
	"sort"
	"testing"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		YAML   string
		Tokens token.Tokens
	}{
		{
			YAML: `null
		`,
			Tokens: token.Tokens{
				{
					Type:          token.NullType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "null",
					Origin:        "null\n",
				},
			},
		},
		{
			YAML: `{}
		`,
			Tokens: token.Tokens{
				{
					Type:          token.MappingStartType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         "{",
					Origin:        "{",
				},
				{
					Type:          token.MappingEndType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         "}",
					Origin:        "}",
				},
			},
		},
		{
			YAML: `v: hi
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "hi",
					Origin:        " hi\n",
				},
			},
		},
		{
			YAML: `v: "true"
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "true",
					Origin:        " \"true\"",
				},
			},
		},
		{
			YAML: `v: "false"
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "false",
					Origin:        " \"false\"",
				},
			},
		},
		{
			YAML: `v: true
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.BoolType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "true",
					Origin:        " true\n",
				},
			},
		},
		{
			YAML: `v: false
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.BoolType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "false",
					Origin:        " false\n",
				},
			},
		},
		{
			YAML: `v: 10
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "10",
					Origin:        " 10\n",
				},
			},
		},
		{
			YAML: `v: -10
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "-10",
					Origin:        " -10\n",
				},
			},
		},
		{
			YAML: `v: 42
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "42",
					Origin:        " 42\n",
				},
			},
		},
		{
			YAML: `v: 4294967296
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "4294967296",
					Origin:        " 4294967296\n",
				},
			},
		},
		{
			YAML: `v: "10"
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "10",
					Origin:        " \"10\"",
				},
			},
		},
		{
			YAML: `v: 0.1
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.FloatType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "0.1",
					Origin:        " 0.1\n",
				},
			},
		},
		{
			YAML: `v: 0.99
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.FloatType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "0.99",
					Origin:        " 0.99\n",
				},
			},
		},
		{
			YAML: `v: -0.1
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.FloatType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "-0.1",
					Origin:        " -0.1\n",
				},
			},
		},
		{
			YAML: `v: .inf
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.InfinityType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         ".inf",
					Origin:        " .inf\n",
				},
			},
		},
		{
			YAML: `v: -.inf
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.InfinityType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "-.inf",
					Origin:        " -.inf\n",
				},
			},
		},
		{
			YAML: `v: .nan
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.NanType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         ".nan",
					Origin:        " .nan\n",
				},
			},
		},
		{
			YAML: `v: null
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.NullType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "null",
					Origin:        " null\n",
				},
			},
		},
		{
			YAML: `v: ""
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "",
					Origin:        " \"\"",
				},
			},
		},
		{
			YAML: `v:
		- A
		- B
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\n\t\t-",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "A",
					Origin:        " A\n",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\t\t-",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "B",
					Origin:        " B\n",
				},
			},
		},
		{
			YAML: `v:
		- A
		- |-
		 B
		 C
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\n\t\t-",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "A",
					Origin:        " A\n",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\t\t-",
				},
				{
					Type:          token.LiteralType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockScalarIndicator,
					Value:         "|-",
					Origin:        " |-\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\t B",
					Origin:        "\t\t B\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\t C",
					Origin:        "\t\t C\n",
				},
			},
		},
		{
			YAML: `v:
		- A
		- 1
		- B:
		 - 2
		 - 3
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "v",
					Origin:        "v",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\n\t\t-",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "A",
					Origin:        " A\n",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\t\t-",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "1",
					Origin:        " 1\n",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\t\t-",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "B",
					Origin:        " B",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\n\t\t -",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "2",
					Origin:        " 2\n",
				},
				{
					Type:          token.SequenceEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         "-",
					Origin:        "\t\t -",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "3",
					Origin:        " 3\n",
				},
			},
		},
		{
			YAML: `a:
		 b: c
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\t b",
					Origin:        "\n\t\t b",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "c",
					Origin:        " c\n",
				},
			},
		},
		{
			YAML: `a: '-'
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SingleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "-",
					Origin:        " '-'",
				},
			},
		},
		{
			YAML: `123
		`,
			Tokens: token.Tokens{
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "123",
					Origin:        "123\n",
				},
			},
		},
		{
			YAML: `hello: world
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "hello",
					Origin:        "hello",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "world",
					Origin:        " world\n",
				},
			},
		},
		{
			YAML: `a: null
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.NullType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "null",
					Origin:        " null\n",
				},
			},
		},
		{
			YAML: `a: {x: 1}
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.MappingStartType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         "{",
					Origin:        " {",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "x",
					Origin:        "x",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "1",
					Origin:        " 1",
				},
				{
					Type:          token.MappingEndType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         "}",
					Origin:        "}",
				},
			},
		},
		{
			YAML: `a: [1, 2]
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SequenceStartType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         "[",
					Origin:        " [",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "1",
					Origin:        "1",
				},
				{
					Type:          token.CollectEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         ",",
					Origin:        ",",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "2",
					Origin:        " 2",
				},
				{
					Type:          token.SequenceEndType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         "]",
					Origin:        "]",
				},
			},
		},
		{
			YAML: `t2: 2018-01-09T10:40:47Z
		t4: 2098-01-09T10:40:47Z
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "t2",
					Origin:        "t2",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "2018-01-09T10:40:47Z",
					Origin:        " 2018-01-09T10:40:47Z\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\tt4",
					Origin:        "\t\tt4",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "2098-01-09T10:40:47Z",
					Origin:        " 2098-01-09T10:40:47Z\n",
				},
			},
		},
		{
			YAML: `a: {b: c, d: e}
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.MappingStartType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         "{",
					Origin:        " {",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "b",
					Origin:        "b",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "c",
					Origin:        " c",
				},
				{
					Type:          token.CollectEntryType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         ",",
					Origin:        ",",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "d",
					Origin:        " d",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "e",
					Origin:        " e",
				},
				{
					Type:          token.MappingEndType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.FlowCollectionIndicator,
					Value:         "}",
					Origin:        "}",
				},
			},
		},
		{
			YAML: `a: 3s
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "3s",
					Origin:        " 3s\n",
				},
			},
		},
		{
			YAML: `a: <foo>
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "<foo>",
					Origin:        " <foo>\n",
				},
			},
		},
		{
			YAML: `a: "1:1"
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "1:1",
					Origin:        " \"1:1\"",
				},
			},
		},
		{
			YAML: `a: "\0"
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "\\0",
					Origin:        " \"\\0\"",
				},
			},
		},
		{
			YAML: `a: !!binary gIGC
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.TagType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.NodePropertyIndicator,
					Value:         "!!binary",
					Origin:        " !!binary ",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "gIGC",
					Origin:        "gIGC\n",
				},
			},
		},
		{
			YAML: `a: !!binary |
		 kJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJ
		 CQ
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.TagType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.NodePropertyIndicator,
					Value:         "!!binary",
					Origin:        " !!binary ",
				},
				{
					Type:          token.LiteralType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockScalarIndicator,
					Value:         "|",
					Origin:        "|\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\t kJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJ",
					Origin:        "\t\t kJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJCQkJ\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\t CQ",
					Origin:        "\t\t CQ\n",
				},
			},
		},
		{
			YAML: `b: 2
		a: 1
		d: 4
		c: 3
		sub:
		 e: 5
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "b",
					Origin:        "b",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "2",
					Origin:        " 2\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\ta",
					Origin:        "\t\ta",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "1",
					Origin:        " 1\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\td",
					Origin:        "\t\td",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "4",
					Origin:        " 4\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\tc",
					Origin:        "\t\tc",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "3",
					Origin:        " 3\n",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\tsub",
					Origin:        "\t\tsub",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\t e",
					Origin:        "\n\t\t e",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.IntegerType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "5",
					Origin:        " 5\n",
				},
			},
		},
		{
			YAML: `a: 1.2.3.4
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "1.2.3.4",
					Origin:        " 1.2.3.4\n",
				},
			},
		},
		{
			YAML: `a: "2015-02-24T18:19:39Z"
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "2015-02-24T18:19:39Z",
					Origin:        " \"2015-02-24T18:19:39Z\"",
				},
			},
		},
		{
			YAML: `a: 'b: c'
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SingleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "b: c",
					Origin:        " 'b: c'",
				},
			},
		},
		{
			YAML: `a: 'Hello #comment'
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SingleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "Hello #comment",
					Origin:        " 'Hello #comment'",
				},
			},
		},
		{
			YAML: `a: 100.5
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.FloatType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "100.5",
					Origin:        " 100.5\n",
				},
			},
		},
		{
			YAML: `a: bogus
		`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "bogus",
					Origin:        " bogus\n",
				},
			},
		},
		{
			YAML: `"a": double quoted map key`,
			Tokens: token.Tokens{
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "a",
					Origin:        "\"a\"",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "double quoted map key",
					Origin:        " double quoted map key",
				},
			},
		},
		{
			YAML: `'a': single quoted map key`,
			Tokens: token.Tokens{
				{
					Type:          token.SingleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "a",
					Origin:        "'a'",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "single quoted map key",
					Origin:        " single quoted map key",
				},
			},
		},
		{
			YAML: `a: "double quoted"
		b: "value map"`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "double quoted",
					Origin:        " \"double quoted\"",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\tb",
					Origin:        "\n\t\tb",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "value map",
					Origin:        " \"value map\"",
				},
			},
		},
		{
			YAML: `a: 'single quoted'
		b: 'value map'`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "a",
					Origin:        "a",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SingleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "single quoted",
					Origin:        " 'single quoted'",
				},
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "\t\tb",
					Origin:        "\n\t\tb",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SingleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "value map",
					Origin:        " 'value map'",
				},
			},
		},
		{
			YAML: `json: '\"expression\": \"thi:\"'`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "json",
					Origin:        "json",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.SingleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "\\\"expression\\\": \\\"thi:\\\"",
					Origin:        " '\\\"expression\\\": \\\"thi:\\\"'",
				},
			},
		},
		{
			YAML: `json: "\"expression\": \"thi:\""`,
			Tokens: token.Tokens{
				{
					Type:          token.StringType,
					CharacterType: token.CharacterTypeMiscellaneous,
					Indicator:     token.NotIndicator,
					Value:         "json",
					Origin:        "json",
				},
				{
					Type:          token.MappingValueType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.BlockStructureIndicator,
					Value:         ":",
					Origin:        ":",
				},
				{
					Type:          token.DoubleQuoteType,
					CharacterType: token.CharacterTypeIndicator,
					Indicator:     token.QuotedScalarIndicator,
					Value:         "\"expression\": \"thi:\"",
					Origin:        " \"\\\"expression\\\": \\\"thi:\\\"\"",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.YAML, func(t *testing.T) {
			tokens := lexer.Tokenize(test.YAML)
			if len(tokens) != len(test.Tokens) {
				t.Fatalf("Tokenize(%q) token count mismatch, expected: %d got: %d", test.YAML, len(test.Tokens), len(tokens))
			}
			for i := range test.Tokens {
				if tokens[i].Type != test.Tokens[i].Type {
					t.Errorf("Tokenize(%q)[%d] token.Type mismatch, expected: %s got: %s", test.YAML, i, test.Tokens[i].Type, tokens[i].Type)
				}
				if tokens[i].CharacterType != test.Tokens[i].CharacterType {
					t.Errorf("Tokenize(%q)[%d] token.CharacterType mismatch, expected: %s got: %s", test.YAML, i, test.Tokens[i].CharacterType, tokens[i].CharacterType)
				}
				if tokens[i].Indicator != test.Tokens[i].Indicator {
					t.Errorf("Tokenize(%q)[%d] token.Indicator mismatch, expected: %s got: %s", test.YAML, i, test.Tokens[i].Indicator, tokens[i].Indicator)
				}
				if tokens[i].Value != test.Tokens[i].Value {
					t.Errorf("Tokenize(%q)[%d] token.Value mismatch, expected: %q got: %q", test.YAML, i, test.Tokens[i].Value, tokens[i].Value)
				}
				if tokens[i].Origin != test.Tokens[i].Origin {
					t.Errorf("Tokenize(%q)[%d] token.Origin mismatch, expected: %q got: %q", test.YAML, i, test.Tokens[i].Origin, tokens[i].Origin)
				}
			}
		})
	}
}

type testToken struct {
	line   int
	column int
	value  string
}

func TestSingleLineToken_ValueLineColumnPosition(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect map[int]string // Column -> Value map.
	}{
		{
			name: "single quote, single value array",
			src:  "test: ['test']",
			expect: map[int]string{
				1:  "test",
				5:  ":",
				7:  "[",
				8:  "test",
				14: "]",
			},
		},
		{
			name: "double quote, single value array",
			src:  `test: ["test"]`,
			expect: map[int]string{
				1:  "test",
				5:  ":",
				7:  "[",
				8:  "test",
				14: "]",
			},
		},
		{
			name: "no quotes, single value array",
			src:  "test: [somevalue]",
			expect: map[int]string{
				1:  "test",
				5:  ":",
				7:  "[",
				8:  "somevalue",
				17: "]",
			},
		},
		{
			name: "single quote, multi value array",
			src:  "myarr: ['1','2','3', '444' , '55','66' ,  '77'  ]",
			expect: map[int]string{
				1:  "myarr",
				6:  ":",
				8:  "[",
				9:  "1",
				12: ",",
				13: "2",
				16: ",",
				17: "3",
				20: ",",
				22: "444",
				28: ",",
				30: "55",
				34: ",",
				35: "66",
				40: ",",
				43: "77",
				49: "]",
			},
		},
		{
			name: "double quote, multi value array",
			src:  `myarr: ["1","2","3", "444" , "55","66" ,  "77"  ]`,
			expect: map[int]string{
				1:  "myarr",
				6:  ":",
				8:  "[",
				9:  "1",
				12: ",",
				13: "2",
				16: ",",
				17: "3",
				20: ",",
				22: "444",
				28: ",",
				30: "55",
				34: ",",
				35: "66",
				40: ",",
				43: "77",
				49: "]",
			},
		},
		{
			name: "no quote, multi value array",
			src:  "numbers: [1, 5, 99,100, 3, 7 ]",
			expect: map[int]string{
				1:  "numbers",
				8:  ":",
				10: "[",
				11: "1",
				12: ",",
				14: "5",
				15: ",",
				17: "99",
				19: ",",
				20: "100",
				23: ",",
				25: "3",
				26: ",",
				28: "7",
				30: "]",
			},
		},
		{
			name: "double quotes, nested arrays",
			src:  `Strings: ["1",["2",["3"]]]`,
			expect: map[int]string{
				1:  "Strings",
				8:  ":",
				10: "[",
				11: "1",
				14: ",",
				15: "[",
				16: "2",
				19: ",",
				20: "[",
				21: "3",
				24: "]",
				25: "]",
				26: "]",
			},
		},
		{
			name: "mixed quotes, nested arrays",
			src:  `Values: [1,['2',"3",4,["5",6]]]`,
			expect: map[int]string{
				1:  "Values",
				7:  ":",
				9:  "[",
				10: "1",
				11: ",",
				12: "[",
				13: "2",
				16: ",",
				17: "3",
				20: ",",
				21: "4",
				22: ",",
				23: "[",
				24: "5",
				27: ",",
				28: "6",
				29: "]",
				30: "]",
				31: "]",
			},
		},
		{
			name: "double quote, empty array",
			src:  `Empty: ["", ""]`,
			expect: map[int]string{
				1:  "Empty",
				6:  ":",
				8:  "[",
				9:  "",
				11: ",",
				13: "",
				15: "]",
			},
		},
		{
			name: "double quote key",
			src:  `"a": b`,
			expect: map[int]string{
				1: "a",
				4: ":",
				6: "b",
			},
		},
		{
			name: "single quote key",
			src:  `'a': b`,
			expect: map[int]string{
				1: "a",
				4: ":",
				6: "b",
			},
		},
		{
			name: "double quote key and value",
			src:  `"a": "b"`,
			expect: map[int]string{
				1: "a",
				4: ":",
				6: "b",
			},
		},
		{
			name: "single quote key and value",
			src:  `'a': 'b'`,
			expect: map[int]string{
				1: "a",
				4: ":",
				6: "b",
			},
		},
		{
			name: "double quote key, single quote value",
			src:  `"a": 'b'`,
			expect: map[int]string{
				1: "a",
				4: ":",
				6: "b",
			},
		},
		{
			name: "single quote key, double quote value",
			src:  `'a': "b"`,
			expect: map[int]string{
				1: "a",
				4: ":",
				6: "b",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := lexer.Tokenize(tc.src)
			sort.Slice(got, func(i, j int) bool {
				return got[i].Position.Column < got[j].Position.Column
			})
			var expected []testToken
			for k, v := range tc.expect {
				tt := testToken{
					line:   1,
					column: k,
					value:  v,
				}
				expected = append(expected, tt)
			}
			sort.Slice(expected, func(i, j int) bool {
				return expected[i].column < expected[j].column
			})
			if len(got) != len(expected) {
				t.Errorf("Tokenize(%s) token count mismatch, expected:%d got:%d", tc.src, len(expected), len(got))
			}
			for i, tok := range got {
				if !tokenMatches(tok, expected[i]) {
					t.Errorf("Tokenize(%s) expected:%+v got line:%d column:%d value:%s", tc.src, expected[i], tok.Position.Line, tok.Position.Column, tok.Value)
				}
			}
		})
	}
}

func tokenMatches(t *token.Token, e testToken) bool {
	return t != nil && t.Position != nil &&
		t.Value == e.value &&
		t.Position.Line == e.line &&
		t.Position.Column == e.column
}

func TestMultiLineToken_ValueLineColumnPosition(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect []testToken
	}{
		{
			name: "double quote",
			src: `one: "1 2 3 4 5"
two: "1 2
3 4
5"
three: "1 2 3 4
5"`,
			expect: []testToken{
				{
					line:   1,
					column: 1,
					value:  "one",
				},
				{
					line:   1,
					column: 4,
					value:  ":",
				},
				{
					line:   1,
					column: 6,
					value:  "1 2 3 4 5",
				},
				{
					line:   2,
					column: 1,
					value:  "two",
				},
				{
					line:   2,
					column: 4,
					value:  ":",
				},
				{
					line:   2,
					column: 6,
					value:  "1 2 3 4 5",
				},
				{
					line:   5,
					column: 1,
					value:  "three",
				},
				{
					line:   5,
					column: 6,
					value:  ":",
				},
				{
					line:   5,
					column: 8,
					value:  "1 2 3 4 5",
				},
			},
		},
		{
			name: "single quote in an array",
			src: `arr: ['1', 'and
two']
last: 'hello'`,
			expect: []testToken{
				{
					line:   1,
					column: 1,
					value:  "arr",
				},
				{
					line:   1,
					column: 4,
					value:  ":",
				},
				{
					line:   1,
					column: 6,
					value:  "[",
				},
				{
					line:   1,
					column: 7,
					value:  "1",
				},
				{
					line:   1,
					column: 10,
					value:  ",",
				},
				{
					line:   1,
					column: 12,
					value:  "and two",
				},
				{
					line:   2,
					column: 5,
					value:  "]",
				},
				{
					line:   3,
					column: 1,
					value:  "last",
				},
				{
					line:   3,
					column: 5,
					value:  ":",
				},
				{
					line:   3,
					column: 7,
					value:  "hello",
				},
			},
		},
		{
			name: "single quote and double quote",
			src: `foo: "test




bar"
foo2: 'bar2'`,
			expect: []testToken{
				{
					line:   1,
					column: 1,
					value:  "foo",
				},
				{
					line:   1,
					column: 4,
					value:  ":",
				},
				{
					line:   1,
					column: 6,
					value:  "test     bar",
				},
				{
					line:   7,
					column: 1,
					value:  "foo2",
				},
				{
					line:   7,
					column: 5,
					value:  ":",
				},
				{
					line:   7,
					column: 7,
					value:  "bar2",
				},
			},
		},
		{
			name: "single and double quote map keys",
			src: `"a": test
'b': 1
c: true`,
			expect: []testToken{
				{
					line:   1,
					column: 1,
					value:  "a",
				},
				{
					line:   1,
					column: 4,
					value:  ":",
				},
				{
					line:   1,
					column: 6,
					value:  "test",
				},
				{
					line:   2,
					column: 1,
					value:  "b",
				},
				{
					line:   2,
					column: 4,
					value:  ":",
				},
				{
					line:   2,
					column: 6,
					value:  "1",
				},
				{
					line:   3,
					column: 1,
					value:  "c",
				},
				{
					line:   3,
					column: 2,
					value:  ":",
				},
				{
					line:   3,
					column: 4,
					value:  "true",
				},
			},
		},
		{
			name: "issue326",
			src: `a: |
  Text
b: 1`,
			expect: []testToken{
				{
					line:   1,
					column: 1,
					value:  "a",
				},
				{
					line:   1,
					column: 2,
					value:  ":",
				},
				{
					line:   1,
					column: 4,
					value:  "|",
				},
				{
					line:   2,
					column: 3,
					value:  "Text\n",
				},
				{
					line:   3,
					column: 1,
					value:  "b",
				},
				{
					line:   3,
					column: 2,
					value:  ":",
				},
				{
					line:   3,
					column: 4,
					value:  "1",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := lexer.Tokenize(tc.src)
			sort.Slice(got, func(i, j int) bool {
				// sort by line, then column
				if got[i].Position.Line < got[j].Position.Line {
					return true
				} else if got[i].Position.Line == got[j].Position.Line {
					return got[i].Position.Column < got[j].Position.Column
				}
				return false
			})
			sort.Slice(tc.expect, func(i, j int) bool {
				if tc.expect[i].line < tc.expect[j].line {
					return true
				} else if tc.expect[i].line == tc.expect[j].line {
					return tc.expect[i].column < tc.expect[j].column
				}
				return false
			})
			if len(got) != len(tc.expect) {
				t.Errorf("Tokenize() token count mismatch, expected:%d got:%d", len(tc.expect), len(got))
			}
			for i, tok := range got {
				if !tokenMatches(tok, tc.expect[i]) {
					t.Errorf("Tokenize() expected:%+v got line:%d column:%d value:%s", tc.expect[i], tok.Position.Line, tok.Position.Column, tok.Value)
				}
			}
		})
	}
}
