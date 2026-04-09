package yaml_test

import (
	"bytes"
	"encoding/json"
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
		"00A: 0000A",
		"{\"000\":0000A,",
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
		if err := yaml.Unmarshal(src, &v); err != nil {
			t.Log(err.Error())
		}
	})
}

// FuzzEncodeFromJSON checks that any JSON encoded value can also be encoded as YAML... and decoded.
func FuzzEncodeFromJSON(f *testing.F) {
	f.Add(`null`)
	f.Add(`""`)
	f.Add(`0`)
	f.Add(`true`)
	f.Add(`false`)
	f.Add(`{}`)
	f.Add(`[]`)
	f.Add(`[[]]`)
	f.Add(`{"a":[]}`)

	f.Fuzz(func(t *testing.T, s string) {

		var v interface{}
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			t.Skip("not valid JSON")
		}

		t.Logf("JSON %q", s)
		t.Logf("Go   %T %#[1]v <%[1]x>", v)

		// Encode as YAML
		b, err := yaml.Marshal(v)
		if err != nil {
			t.Error(err)
		}
		t.Logf("YAML %q <%[1]x>", b)

		// Decode as YAML
		var v2 interface{}
		if err := yaml.Unmarshal(b, &v2); err != nil {
			t.Error(err)
		}

		t.Logf("Go   %T %#[1]v <%[1]x>", v2)

		/*
			// Handling of number is different, so we can't have universal exact matching
			if !reflect.DeepEqual(v2, v) {
				t.Errorf("mismatch:\n-      got: %#v\n- expected: %#v", v2, v)
			}
		*/

		b2, err := yaml.Marshal(v2)
		if err != nil {
			t.Error(err)
		}
		t.Logf("YAML %q <%[1]x>", b2)

		if !bytes.Equal(b, b2) {
			t.Errorf("Marshal->Unmarshal->Marshal mismatch:\n- expected: %q\n- got:      %q", b, b2)
		}

	})
}

func TestEncodeString(t *testing.T) {
	b, _ := yaml.Marshal(`\n`)
	t.Logf("%q <%[1]x>", string(b))
}
