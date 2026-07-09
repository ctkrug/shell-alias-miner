package miner

import (
	"strings"
	"unicode"
)

// tokenize splits a command line into shell words, treating a single- or
// double-quoted span as one token (quotes included verbatim) so an argument
// like a commit message doesn't get split on its internal spaces.
func tokenize(command string) []string {
	// The overwhelming majority of commands have no quotes at all; skip the
	// rune-by-rune scan below and let the stdlib's tuned byte scanner do it
	// for the common case. This matters at 100k+-line scale.
	if !strings.ContainsAny(command, `'"`) {
		return strings.Fields(command)
	}

	var (
		tokens []string
		b      []rune
		quote  rune
	)

	flush := func() {
		if len(b) > 0 {
			tokens = append(tokens, string(b))
			b = b[:0]
		}
	}

	for _, r := range command {
		switch {
		case quote != 0:
			b = append(b, r)
			if r == quote {
				quote = 0
			}
		case r == '\'' || r == '"':
			quote = r
			b = append(b, r)
		case unicode.IsSpace(r):
			flush()
		default:
			b = append(b, r)
		}
	}
	flush()

	return tokens
}
