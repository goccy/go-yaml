package yaml

import (
	"bytes"
	"strings"
	"testing"
)

// TestInlineAfterDepth tests the depth-based inline formatting functionality
func TestInlineAfterDepth(t *testing.T) {
	tests := []struct {
		name     string
		depth    int
		input    interface{}
		expected string
	}{
		{
			name:  "simple nested structure with depth 1",
			depth: 1,
			input: struct {
				Level1 struct {
					Level2 struct {
						Data string `yaml:"data"`
					} `yaml:"level2"`
				} `yaml:"level1"`
			}{
				Level1: struct {
					Level2 struct {
						Data string `yaml:"data"`
					} `yaml:"level2"`
				}{
					Level2: struct {
						Data string `yaml:"data"`
					}{
						Data: "test",
					},
				},
			},
			expected: "level1:\n  level2: {data: test}\n",
		},
		{
			name:  "nested structure with depth 2",
			depth: 2,
			input: struct {
				Level1 struct {
					Level2 struct {
						Level3 struct {
							Data string `yaml:"data"`
						} `yaml:"level3"`
					} `yaml:"level2"`
				} `yaml:"level1"`
			}{
				Level1: struct {
					Level2 struct {
						Level3 struct {
							Data string `yaml:"data"`
						} `yaml:"level3"`
					} `yaml:"level2"`
				}{
					Level2: struct {
						Level3 struct {
							Data string `yaml:"data"`
						} `yaml:"level3"`
					}{
						Level3: struct {
							Data string `yaml:"data"`
						}{
							Data: "deep",
						},
					},
				},
			},
			expected: "level1:\n  level2:\n    level3: {data: deep}\n",
		},
		{
			name:  "map with nested structures",
			depth: 1,
			input: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": map[string]interface{}{
						"value": "test",
					},
				},
			},
			expected: "outer:\n  inner: {value: test}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			enc := NewEncoder(&buf, InlineAfterDepth(tt.depth))

			err := enc.Encode(tt.input)
			if err != nil {
				t.Fatalf("Failed to encode: %v", err)
			}

			result := buf.String()
			if strings.TrimSpace(result) != strings.TrimSpace(tt.expected) {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}

// TestInlineAfterDepthWithInvalidDepth tests error handling for invalid depth values
func TestInlineAfterDepthWithInvalidDepth(t *testing.T) {
	var buf bytes.Buffer

	// Test negative depth
	enc := NewEncoder(&buf)
	err := InlineAfterDepth(-1)(enc)
	if err == nil {
		t.Error("Expected error for negative depth, got nil")
	}

	expectedError := "inline depth must be non-negative, got -1"
	if err.Error() != expectedError {
		t.Errorf("Expected error message: %s, got: %s", expectedError, err.Error())
	}
}

// TestInlineAfterDepthDisabled tests that the feature is disabled by default
func TestInlineAfterDepthDisabled(t *testing.T) {
	input := struct {
		Level1 struct {
			Level2 struct {
				Data string `yaml:"data"`
			} `yaml:"level2"`
		} `yaml:"level1"`
	}{
		Level1: struct {
			Level2 struct {
				Data string `yaml:"data"`
			} `yaml:"level2"`
		}{
			Level2: struct {
				Data string `yaml:"data"`
			}{
				Data: "test",
			},
		},
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf) // No InlineAfterDepth option

	err := enc.Encode(input)
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	result := buf.String()
	expected := "level1:\n  level2:\n    data: test\n"

	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Expected normal formatting without inline option:\n%s\nGot:\n%s", expected, result)
	}
}

// TestInlineAfterDepthWithComplexStructure tests with more complex nested structures
func TestInlineAfterDepthWithComplexStructure(t *testing.T) {
	type DeepStruct struct {
		Name  string            `yaml:"name"`
		Value int               `yaml:"value"`
		Meta  map[string]string `yaml:"meta"`
	}

	input := struct {
		Config struct {
			Database struct {
				Primary   DeepStruct `yaml:"primary"`
				Secondary DeepStruct `yaml:"secondary"`
			} `yaml:"database"`
			Cache struct {
				Redis DeepStruct `yaml:"redis"`
			} `yaml:"cache"`
		} `yaml:"config"`
	}{
		Config: struct {
			Database struct {
				Primary   DeepStruct `yaml:"primary"`
				Secondary DeepStruct `yaml:"secondary"`
			} `yaml:"database"`
			Cache struct {
				Redis DeepStruct `yaml:"redis"`
			} `yaml:"cache"`
		}{
			Database: struct {
				Primary   DeepStruct `yaml:"primary"`
				Secondary DeepStruct `yaml:"secondary"`
			}{
				Primary: DeepStruct{
					Name:  "main",
					Value: 100,
					Meta:  map[string]string{"type": "postgres"},
				},
				Secondary: DeepStruct{
					Name:  "backup",
					Value: 50,
					Meta:  map[string]string{"type": "mysql"},
				},
			},
			Cache: struct {
				Redis DeepStruct `yaml:"redis"`
			}{
				Redis: DeepStruct{
					Name:  "cache",
					Value: 200,
					Meta:  map[string]string{"type": "redis"},
				},
			},
		},
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf, InlineAfterDepth(2))

	err := enc.Encode(input)
	if err != nil {
		t.Fatalf("Failed to encode: %v", err)
	}

	result := buf.String()

	// Check that deep structures (level 3+) are inlined
	if !strings.Contains(result, "{") {
		t.Error("Expected inline formatting (with braces) for deep structures")
	}

	// Check that top-level structures are not inlined
	lines := strings.Split(result, "\n")
	if len(lines) < 5 { // Should have multiple lines for non-inlined top levels
		t.Error("Expected multi-line formatting for top-level structures")
	}
}
