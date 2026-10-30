package yaml_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

func TestExplicitBlockIndentTrailingEmptyLines(t *testing.T) {
	endings := []string{"\n\n", "\n\n\n", "\n \n", "\n "}
	for _, style := range []string{"|", ">"} {
		for _, suffix := range []string{"2", "2-", "-2", "2+", "+2"} {
			header := style + suffix
			for _, lineBreak := range []string{"\n", "\r\n", "\r"} {
				for _, ending := range endings {
					name := fmt.Sprintf("%s/break=%q/end=%q", header, lineBreak, ending)
					t.Run(name, func(t *testing.T) {
						source := strings.ReplaceAll("key: "+header+"\n  first\n  second"+ending, "\n", lineBreak)
						var got map[string]string
						if err := yaml.Unmarshal([]byte(source), &got); err != nil {
							t.Fatalf("valid block scalar rejected: %v", err)
						}
						want := "first\nsecond"
						if style == ">" {
							want = "first second"
						}
						switch {
						case strings.Contains(suffix, "-"):
						case strings.Contains(suffix, "+"):
							want += strings.Repeat("\n", strings.Count(ending, "\n"))
						default:
							want += "\n"
						}
						if got["key"] != want {
							t.Errorf("decoded value = %q, want %q", got["key"], want)
						}
					})
				}
			}
		}
	}
}

func TestExplicitBlockIndentRejectsUnderIndentedContent(t *testing.T) {
	for _, header := range []string{"|2", "|2-", ">+2"} {
		for _, content := range []string{" x", "  first\n x", "\tx"} {
			for _, ending := range []string{"", "\n"} {
				t.Run(fmt.Sprintf("%s/content=%q/end=%q", header, content, ending), func(t *testing.T) {
					var got map[string]string
					if err := yaml.Unmarshal([]byte("key: "+header+"\n"+content+ending), &got); err == nil {
						t.Errorf("accepted under-indented content: %#v", got)
					}
				})
			}
		}
	}
}
