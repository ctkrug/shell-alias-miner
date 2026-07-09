package miner

import "unicode"

// tokenize splits a command line into shell words, treating a single- or
// double-quoted span as one token (quotes included verbatim) so an argument
// like a commit message doesn't get split on its internal spaces.
func tokenize(command string) []string {
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
