package yaml_test

import (
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func FuzzUnmarshalToMap(f *testing.F) {
	const validYAML = `
id: 1
message: Hello, World
verified: true
`

	invalidYAML := []string{
		"0::",
		"{0",
		"*-0",
		">\n>",
		"&{0",
		"0_",
		"0\n:",
		"0\n-",
		"0\n0",
		"0\n0\n",
		"0\n0\n0",
		"0\n0\n0\n",
		"0\n0\n0\n0",
		"0\n0\n0\n0\n",
		"0\n0\n0\n0\n0",
		"0\n0\n0\n0\n0\n",
		"0\n0\n0\n0\n0\n0",
		"0\n0\n0\n0\n0\n0\n",
		"",
	}

	f.Add([]byte(validYAML))
	for _, s := range invalidYAML {
		f.Add([]byte(s))
		f.Add([]byte(validYAML + s))
		f.Add([]byte(s + validYAML))
		f.Add([]byte(s + validYAML + s))
		f.Add([]byte(strings.Repeat(s, 3)))
	}

	f.Fuzz(func(t *testing.T, src []byte) {
		v := map[string]any{}
		_ = yaml.Unmarshal(src, &v)
	})
}
