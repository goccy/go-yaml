package yamltestsuite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type TestSuite struct {
	Name    string
	InYAML  []byte
	InJSON  []any
	OutYAML []byte
	Error   bool
}

func curDir() string {
	_, file, _, _ := runtime.Caller(0) //nolint:dogsled
	return filepath.Dir(file)
}

func TestSuites() ([]*TestSuite, error) {
	dir := curDir()
	testMap := make(map[string]*TestSuite)
	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			// this file.
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		name := strings.TrimPrefix(path, dir+"/")
		name = strings.TrimSuffix(name, "/"+filepath.Base(name))
		if _, exists := testMap[name]; !exists {
			testMap[name] = &TestSuite{}
		}
		f, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fileName := filepath.Base(path)
		switch fileName {
		case "in.yaml":
			testMap[name].InYAML = f
		case "in.json":
			dec := json.NewDecoder(bytes.NewReader(f))
			var inJSON []any
			for {
				var v any
				if err := dec.Decode(&v); err != nil {
					if err == io.EOF {
						break
					}
					return fmt.Errorf("failed to decode json: %s: %s: %w", name, string(f), err)
				}
				inJSON = append(inJSON, v)
			}
			testMap[name].InJSON = inJSON
		case "out.yaml":
			testMap[name].OutYAML = f
		case "error":
			testMap[name].Error = true
		}
		testMap[name].Name = name
		return nil
	}); err != nil {
		return nil, err
	}

	tests := make([]*TestSuite, 0, len(testMap))
	for _, test := range testMap {
		if test.InYAML == nil {
			continue
		}
		tests = append(tests, test)
	}
	sort.Slice(tests, func(i, j int) bool {
		return tests[i].Name < tests[j].Name
	})
	return tests, nil
}
