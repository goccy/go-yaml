package token

import "testing"

// Regression test for goccy/go-yaml#916: scientific-notation numbers without a
// mantissa dot (e.g. "1e+20") must be recognized as floats so that a float64
// survives an Unmarshal -> Marshal -> Unmarshal roundtrip instead of decoding
// back to a string.
func TestToNumberScientificNotation(t *testing.T) {
	floatCases := []string{
		"1e+20", "1e20", "1E+20", "-1e+20", "1.5e+10", "1.0e+20",
		"1E-5", "0e0", "9e9",
	}
	for _, s := range floatCases {
		num := ToNumber(s)
		if num == nil {
			t.Errorf("ToNumber(%q) = nil, want a float number", s)
			continue
		}
		if num.Type != NumberTypeFloat {
			t.Errorf("ToNumber(%q).Type = %q, want %q", s, num.Type, NumberTypeFloat)
		}
	}

	// Values that merely contain 'e'/'E' but are not scientific notation must
	// not be treated as numbers.
	notNumber := []string{
		"e10", "1e", "1e+", "1ee5", "abc", "0x1e", "1.2.3e4", "1e1.5",
	}
	for _, s := range notNumber {
		if hasDecimalExponent(stripSignAndUnderscore(s)) {
			t.Errorf("hasDecimalExponent(%q) = true, want false", s)
		}
	}
}

// stripSignAndUnderscore mirrors the normalization toNumber applies before the
// type switch, so the helper is exercised on the same input shape.
func stripSignAndUnderscore(s string) string {
	out := make([]byte, 0, len(s))
	if len(s) > 0 && (s[0] == '+' || s[0] == '-') {
		s = s[1:]
	}
	for i := 0; i < len(s); i++ {
		if s[i] != '_' {
			out = append(out, s[i])
		}
	}
	return string(out)
}
