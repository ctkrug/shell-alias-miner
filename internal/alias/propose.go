// Package alias turns mined command candidates into ready-to-paste shell
// alias definitions with a keystrokes-saved estimate.
package alias

import (
	"fmt"
	"strings"

	"github.com/ctkrug/shell-alias-miner/internal/miner"
)

// Proposal is a ready-to-paste alias for a repeated command.
type Proposal struct {
	// Name is the suggested short alias name.
	Name string
	// Command is the original, full command being replaced.
	Command string
	// Definition is the shell snippet to paste into .zshrc/.bashrc.
	Definition string
	// Occurrences is how many times Command was seen in history.
	Occurrences int
	// KeystrokesSaved is the total character count saved across every
	// occurrence, assuming the alias is used in place of the full command.
	KeystrokesSaved int
}

// Propose builds one Proposal per candidate, using name as the alias name
// for the single highest-value candidate and a numbered suffix for the
// rest (cmd2, cmd3, ...) to keep names unique and short.
//
// candidates should already be sorted by relevance (see miner.CountFrequencies);
// Propose preserves that order in its output.
func Propose(candidates []miner.Candidate) []Proposal {
	proposals := make([]Proposal, 0, len(candidates))
	used := map[string]bool{}

	for _, c := range candidates {
		name := uniqueName(nameFor(c.Command), used)
		used[name] = true

		def := fmt.Sprintf("alias %s=%q", name, c.Command)

		// Savings are per future invocation: typing the short name instead
		// of the full command, summed over every time it was already used.
		// The one-time cost of pasting the definition itself doesn't count
		// against it — it's paid once, the savings compound forever after.
		saved := (len(c.Command) - len(name)) * c.Count
		if saved < 0 {
			saved = 0
		}

		proposals = append(proposals, Proposal{
			Name:            name,
			Command:         c.Command,
			Definition:      def,
			Occurrences:     c.Count,
			KeystrokesSaved: saved,
		})
	}

	return proposals
}

// nameFor derives a short alias name from a command's leading tokens,
// e.g. "git status --short" -> "gs".
func nameFor(command string) string {
	fields := strings.Fields(command)
	if len(fields) == 0 {
		return "cmd"
	}

	var b strings.Builder
	for i, f := range fields {
		if i >= 3 {
			break
		}
		if strings.HasPrefix(f, "-") {
			continue
		}
		r := []rune(f)
		b.WriteRune(r[0])
	}

	if b.Len() == 0 {
		return "cmd"
	}
	return strings.ToLower(b.String())
}

func uniqueName(base string, used map[string]bool) string {
	if !used[base] {
		return base
	}
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s%d", base, i)
		if !used[candidate] {
			return candidate
		}
	}
}
