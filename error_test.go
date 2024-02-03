package yaml

import (
	"reflect"
	"testing"

	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/token"
	"golang.org/x/xerrors"
)

func TestAsSyntaxError(t *testing.T) {
	tests := []struct {
		input    error
		expected *TokenScopedError
	}{
		{
			input:    nil,
			expected: nil,
		},
		{
			input:    xerrors.New("dummy test"),
			expected: nil,
		},
		{
			input: errors.ErrSyntax("some error", &token.Token{Value: "123"}),
			expected: &TokenScopedError{
				Msg:   "some error",
				Token: &token.Token{Value: "123"},
			},
		},
		{
			input: xerrors.Errorf(
				"something went wrong: %w",
				errors.ErrSyntax("some error", &token.Token{Value: "123"})),
			expected: &TokenScopedError{
				Msg:   "some error",
				Token: &token.Token{Value: "123"},
			},
		},
		{
			input: errUnknownField("unknown field", &token.Token{Value: "123"}),
			expected: &TokenScopedError{
				Msg:   "unknown field",
				Token: &token.Token{Value: "123"},
			},
		},
		{
			input: errDuplicateKey("duplicate key", &token.Token{Value: "123"}),
			expected: &TokenScopedError{
				Msg:   "duplicate key",
				Token: &token.Token{Value: "123"},
			},
		},
		{
			input: errTypeMismatch(reflect.TypeOf("string"), reflect.TypeOf(0), &token.Token{Value: "123"}),
			expected: &TokenScopedError{
				Msg:   "cannot unmarshal int into Go value of type string",
				Token: &token.Token{Value: "123"},
			},
		},
	}
	for _, test := range tests {
		syntaxErr := AsTokenScopedError(test.input)
		if test.expected == nil {
			if syntaxErr != nil {
				t.Fatalf("wanted nil, but go %v", syntaxErr)
			}
			continue
		}
		if syntaxErr == nil {
			t.Fatalf("must not be nil")
		}
		if *test.expected.Token != *syntaxErr.Token || test.expected.Msg != syntaxErr.Msg {
			t.Fatalf("unexpected output.\nexpect:\n[%v]\nactual:\n[%v]", test.expected, syntaxErr)
		}
	}
}
