// Package miner mines shell history for commands worth aliasing.
package miner

import "sort"

// PatternKind distinguishes a verbatim-repeat candidate from one whose
// trailing argument varies across occurrences.
type PatternKind string

const (
	// KindExact is a command seen with byte-identical text every time.
	KindExact PatternKind = "exact"
	// KindTemplate is a command whose leading tokens repeat but whose
	// trailing token (an argument) varies between occurrences.
	KindTemplate PatternKind = "template"
)

// Candidate is a command pattern seen one or more times in history, along
// with how many times it was seen.
type Candidate struct {
	// Command is a representative full command line: for Kind == KindExact
	// it's the literal repeated command; for Kind == KindTemplate it's one
	// example instance of the pattern.
	Command string
	// Prefix is the fixed leading portion of a KindTemplate candidate (all
	// tokens except the varying trailing argument). Empty for KindExact.
	Prefix string
	Count  int
	// Kind defaults to the zero value "" for callers that only care about
	// exact-match counting; treat "" the same as KindExact.
	Kind PatternKind
}

// CountFrequencies counts how many times each exact command line occurs in
// commands and returns the results sorted by count descending, then by
// command text ascending for a stable order among ties.
//
// This is the baseline miner: it only catches commands repeated verbatim.
// Catching commands that repeat with varying arguments (e.g. a different
// commit message each time) requires n-gram template mining, tracked as
// follow-up work in docs/BACKLOG.md.
func CountFrequencies(commands []string) []Candidate {
	counts := make(map[string]int, len(commands))
	for _, c := range commands {
		if c == "" {
			continue
		}
		counts[c]++
	}

	candidates := make([]Candidate, 0, len(counts))
	for cmd, n := range counts {
		candidates = append(candidates, Candidate{Command: cmd, Count: n, Kind: KindExact})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Count != candidates[j].Count {
			return candidates[i].Count > candidates[j].Count
		}
		return candidates[i].Command < candidates[j].Command
	})

	return candidates
}
