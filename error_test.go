package yaml

import (
	"testing"

	"golang.org/x/xerrors"

	"github.com/goccy/go-yaml/internal/errors"
	"github.com/goccy/go-yaml/token"
)

func TestAsSyntaxError(t *testing.T) {
	tests := []struct {
		input    error
		expected *SyntaxError
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
			expected: &SyntaxError{
				Msg:   "some error",
				Token: &token.Token{Value: "123"},
			},
		},
		{
			input: xerrors.Errorf(
				"something went wrong: %w",
				errors.ErrSyntax("some error", &token.Token{Value: "123"})),
			expected: &SyntaxError{
				Msg:   "some error",
				Token: &token.Token{Value: "123"},
			},
		},
		{
			input: errUnknownField("unknown field", &token.Token{Value: "123"}),
			expected: &SyntaxError{
				Msg:   "unknown field",
				Token: &token.Token{Value: "123"},
			},
		},
	}
	for _, test := range tests {
		syntaxErr := AsSyntaxError(test.input)
		if test.expected == nil {
			if syntaxErr != nil {
				t.Fatalf("wanted nil, but go %v", syntaxErr)
			}
			return
		}
		if *test.expected.Token != *syntaxErr.Token && test.expected.Msg != syntaxErr.Msg {
			t.Fatalf("unexpected output.\nexpect:\n[%v]\nactual:\n[%v]", test.expected, syntaxErr)
		}
	}
}
