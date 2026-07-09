package miner

import (
	"sort"
	"strings"
)

// Mine finds both exact-repeat candidates and n-gram template candidates:
// commands that share every token except a varying trailing argument (a
// commit message, a filename, a container tag). A command only becomes a
// template candidate if at least two distinct trailing arguments are
// observed for the same leading tokens — a group with only one observed
// suffix has no variation to catch and is left as a plain exact-match
// candidate, so history with no argument variation mines identically to
// CountFrequencies.
//
// Results are sorted by Count descending, then Command ascending, matching
// CountFrequencies's ordering.
func Mine(commands []string) []Candidate {
	exactCounts := make(map[string]int, len(commands))
	for _, c := range commands {
		if c == "" {
			continue
		}
		exactCounts[c]++
	}

	// prefixOf remembers, per unique command, which group it landed in (if
	// any), as a flat slice rather than a second map keyed by full command
	// text — cheap to build here and cheap to re-scan below, avoiding a
	// large "covered" set keyed by (potentially long) command strings.
	type prefixOf struct {
		cmd    string
		prefix string // "" if cmd has fewer than 2 tokens
	}

	groups := map[string]*templateGroup{}
	commandPrefixes := make([]prefixOf, 0, len(exactCounts))

	for cmd := range exactCounts {
		tokens := tokenize(cmd)
		if len(tokens) < 2 {
			commandPrefixes = append(commandPrefixes, prefixOf{cmd: cmd})
			continue // no trailing argument to vary
		}

		prefix := strings.Join(tokens[:len(tokens)-1], " ")
		suffix := tokens[len(tokens)-1]

		g, ok := groups[prefix]
		if !ok {
			g = &templateGroup{suffixes: map[string]bool{}}
			groups[prefix] = g
		}
		g.count += exactCounts[cmd]
		g.suffixes[suffix] = true
		if g.example == "" || cmd < g.example {
			g.example = cmd
		}

		commandPrefixes = append(commandPrefixes, prefixOf{cmd: cmd, prefix: prefix})
	}

	candidates := make([]Candidate, 0, len(exactCounts))

	for prefix, g := range groups {
		if len(g.suffixes) < 2 {
			continue // every occurrence used the same argument: no variation
		}
		candidates = append(candidates, Candidate{
			Command: g.example,
			Prefix:  prefix,
			Count:   g.count,
			Kind:    KindTemplate,
		})
	}

	for _, e := range commandPrefixes {
		if e.prefix != "" && len(groups[e.prefix].suffixes) >= 2 {
			continue // covered by a template candidate emitted above
		}
		candidates = append(candidates, Candidate{Command: e.cmd, Count: exactCounts[e.cmd], Kind: KindExact})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Count != candidates[j].Count {
			return candidates[i].Count > candidates[j].Count
		}
		return candidates[i].Command < candidates[j].Command
	})

	return candidates
}

// templateGroup accumulates every command sharing one trailing-argument
// prefix while Mine scans the deduplicated command set.
type templateGroup struct {
	count    int
	suffixes map[string]bool
	example  string
}
