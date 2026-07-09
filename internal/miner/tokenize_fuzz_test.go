package miner

import "testing"

// tokenize has no error return — its only contract is "never panic, and
// never lose or invent characters that change the joined-back-together
// length by more than the quote/space bytes it's allowed to consume." Fuzz
// it directly rather than special-casing every adversarial input by hand.
func FuzzTokenize(f *testing.F) {
	seeds := []string{
		"",
		"   ",
		`git commit -m "fix the bug"`,
		`echo 'multi word arg'`,
		`git commit -m "unterminated`,
		`git commit -m 'unterminated`,
		"🎉🎉🎉",
		`"""""`,
		"'''''",
		"\t\n\r",
		`git commit -m "🎉 done"`,
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, command string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("tokenize(%q) panicked: %v", command, r)
			}
		}()
		tokenize(command)
	})
}
