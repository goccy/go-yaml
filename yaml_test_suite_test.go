package yaml_test

import (
	"bytes"
	"encoding/json"
	"io"
	"runtime/debug"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/testdata/yaml-test-suite"
)

const (
	wip                   = true
	successCountThreshold = 199
)

func fatal(t *testing.T, msg string, args ...any) {
	t.Helper()
	if wip {
		t.Logf(msg, args...)
		return
	}
	t.Fatalf(msg, args...)
}

func TestYAMLTestSuite(t *testing.T) {
	tests, err := yamltestsuite.TestSuites()
	if err != nil {
		t.Fatal(err)
	}

	var (
		success int
		failure int
	)
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			defer func() {
				if e := recover(); e != nil {
					failure++
					fatal(t, "panic occurred.\n[input]\n%s\nstack[%s]", string(test.InYAML), debug.Stack())
					return
				}
			}()

			if test.Error {
				var v any
				if err := yaml.Unmarshal(test.InYAML, &v); err == nil {
					failure++
					fatal(t, "expected error.\n[input]\n%s\n", string(test.InYAML))
					return
				}
				success++
				return
			}

			dec := yaml.NewDecoder(bytes.NewReader(test.InYAML))
			var idx int
			for {
				var v any
				if err := dec.Decode(&v); err != nil {
					if err == io.EOF {
						break
					}
					failure++
					fatal(t, err.Error())
					return
				}
				if len(test.InJSON) <= idx {
					failure++
					fatal(t, "expected json value is nothing but got %v", v)
					return
				}
				expected, err := json.Marshal(test.InJSON[idx])
				if err != nil {
					fatal(t, "failed to encode json value: %v", err)
					return
				}
				got, err := json.Marshal(v)
				if err != nil {
					fatal(t, "failed to encode json value: %v", err)
					return
				}
				if !bytes.Equal(expected, got) {
					failure++
					fatal(t, "json mismatch [%s]:\n[expected]\n%s\n[got]\n%s\n", test.Name, string(expected), string(got))
					return
				}
				idx++
			}
			success++
		})
	}
	total := len(tests)
	if success+failure == total {
		t.Logf("yaml-test-suite result: success/total = %d/%d (%f %%)\n", success, total, float32(success)/float32(total)*100)
	}
	if success < successCountThreshold {
		// degrade occurred.
		t.Fatalf("expected success count is over %d but got %d", successCountThreshold, success)
	}
}
