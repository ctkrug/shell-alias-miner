// Package alias turns mined command candidates into ready-to-paste shell
// alias definitions with a keystrokes-saved estimate.
package alias

import (
	"fmt"
	"strings"

	"github.com/ctkrug/shell-alias-miner/internal/miner"
)

// Proposal is a ready-to-paste alias or function for a repeated command.
type Proposal struct {
	// Name is the suggested short alias/function name.
	Name string
	// Command is the original, full command being replaced (for a
	// KindFunction proposal, one representative instance of it).
	Command string
	// Definition is the shell snippet to paste into .zshrc/.bashrc.
	Definition string
	// Occurrences is how many times Command was seen in history.
	Occurrences int
	// KeystrokesSaved is the total character count saved across every
	// occurrence, assuming the proposal is used in place of the full
	// command.
	KeystrokesSaved int
	// Kind is "alias" for a fixed command or "function" for one with a
	// varying trailing argument.
	Kind string
}

const (
	// KindAlias marks a Proposal for a fixed command with no argument.
	KindAlias = "alias"
	// KindFunction marks a Proposal for a command with a varying trailing
	// argument, proposed as a shell function rather than a plain alias
	// (a plain alias can't take an argument).
	KindFunction = "function"
)

// Propose builds one Proposal per candidate, using name as the alias name
// for the single highest-value candidate and a numbered suffix for the
// rest (cmd2, cmd3, ...) to keep names unique and short.
//
// candidates should already be sorted by relevance (see miner.CountFrequencies);
// Propose preserves that order in its output.
func Propose(candidates []miner.Candidate) []Proposal {
	proposals := make([]Proposal, 0, len(candidates))
	// next[base] is the next numbered suffix to try for that base name, so
	// resolving a collision is O(1) instead of re-probing from scratch on
	// every call — that probe-from-2 approach goes quadratic when many
	// candidates collide on the same base (e.g. a long run of "echo ..."
	// commands, which all derive the same one- or two-letter base).
	next := map[string]int{}

	for _, c := range candidates {
		if containsSecret(c.Command) {
			// Never hand back a pasteable snippet that bakes in a password
			// or token — skip the candidate entirely rather than propose
			// something unsafe.
			continue
		}
		if c.Kind == miner.KindTemplate {
			proposals = append(proposals, proposeFunction(c, next))
			continue
		}
		proposals = append(proposals, proposeAlias(c, next))
	}

	return proposals
}

// proposeAlias builds a fixed-command alias proposal.
func proposeAlias(c miner.Candidate, next map[string]int) Proposal {
	name := uniqueName(nameFor(c.Command), next)

	def := fmt.Sprintf("alias %s=%q", name, c.Command)

	// Savings are per future invocation: typing the short name instead
	// of the full command, summed over every time it was already used.
	// The one-time cost of pasting the definition itself doesn't count
	// against it — it's paid once, the savings compound forever after.
	saved := (len(c.Command) - len(name)) * c.Count
	if saved < 0 {
		saved = 0
	}

	return Proposal{
		Name:            name,
		Command:         c.Command,
		Definition:      def,
		Occurrences:     c.Count,
		KeystrokesSaved: saved,
		Kind:            KindAlias,
	}
}

// proposeFunction builds a varying-argument function proposal. A plain
// alias can't take an argument, so a candidate whose trailing token varies
// gets a one-parameter shell function instead, with the fixed prefix baked
// in and "$1" substituted where the varying argument goes.
func proposeFunction(c miner.Candidate, next map[string]int) Proposal {
	name := uniqueName(nameFor(c.Prefix), next)

	def := fmt.Sprintf(`function %s() { %s "$1"; }`, name, c.Prefix)

	// Same accounting as the alias case, but relative to Prefix: the
	// argument itself is typed either way (as part of the original command
	// or as the function's "$1"), so only the fixed prefix-vs-name delta
	// is a real saving, and it's the same delta on every invocation
	// regardless of how long that occurrence's argument was.
	saved := (len(c.Prefix) - len(name)) * c.Count
	if saved < 0 {
		saved = 0
	}

	return Proposal{
		Name:            name,
		Command:         c.Command,
		Definition:      def,
		Occurrences:     c.Count,
		KeystrokesSaved: saved,
		Kind:            KindFunction,
	}
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

// uniqueName returns base the first time it's requested, then base2, base3,
// ... on each subsequent request for the same base. next tracks, per base,
// the next suffix to hand out, so repeated collisions resolve in O(1)
// instead of re-scanning from the start each time.
func uniqueName(base string, next map[string]int) string {
	n, seen := next[base]
	if !seen {
		next[base] = 2
		return base
	}
	next[base] = n + 1
	return fmt.Sprintf("%s%d", base, n)
}
