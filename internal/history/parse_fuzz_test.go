package history

import (
	"strings"
	"testing"
)

// Parse's contract is "never panic, best-effort on anything" — fuzz it
// directly against arbitrary bytes (not just well-formed history files)
// rather than hand-enumerating malformed shapes.
func FuzzParse(f *testing.F) {
	seeds := []string{
		"",
		"\n\n\n",
		": 1700000000:0;git status\n",
		": not-a-timestamp;git status\n",
		"git status\\\n",
		"git status\\\ngit log\\\n",
		strings.Repeat("a", 5000),
		"git commit -m \"🎉\"\n",
		"\x00\x01\x02binary\xff\xfe",
		": 1:2:3;weird;semicolons;here\n",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, input string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Parse(%q) panicked: %v", input, r)
			}
		}()
		Parse(strings.NewReader(input))
	})
}
