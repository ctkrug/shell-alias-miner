package alias

import (
	"testing"
	"unicode"
)

// nameFor's result is pasted verbatim as a shell identifier, so its only
// real contract — beyond "never panic" — is "every character is a letter
// or digit" (bash/zsh both accept Unicode letters in alias/function names
// fine; it's operators, quotes, and whitespace that break). Fuzz it
// directly against that invariant rather than relying on the hand-picked
// operator/quote examples in name_test.go to be exhaustive.
func FuzzNameFor(f *testing.F) {
	seeds := []string{
		"",
		"git status --short",
		"true && false",
		"make test || echo fail",
		`curl -H "Content-Type: application/json" https://example.com`,
		"!!!",
		"--- --- ---",
		"🎉 party --now",
		"'''",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, command string) {
		name := nameFor(command)
		for _, r := range name {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				t.Fatalf("nameFor(%q) = %q, contains disallowed character %q", command, name, r)
			}
		}
	})
}
