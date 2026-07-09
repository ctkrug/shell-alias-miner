// Package miner mines shell history for commands worth aliasing.
package miner

import "sort"

// Candidate is a command line seen one or more times in history, along
// with how many times it was seen.
type Candidate struct {
	Command string
	Count   int
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
		candidates = append(candidates, Candidate{Command: cmd, Count: n})
	}

	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].Count != candidates[j].Count {
			return candidates[i].Count > candidates[j].Count
		}
		return candidates[i].Command < candidates[j].Command
	})

	return candidates
}
