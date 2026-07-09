// Package pipeline composes the parse -> mine -> propose stages into the
// single call the wasm entrypoint (and its tests) use.
package pipeline

import (
	"strings"

	"github.com/ctkrug/shell-alias-miner/internal/alias"
	"github.com/ctkrug/shell-alias-miner/internal/history"
	"github.com/ctkrug/shell-alias-miner/internal/miner"
)

// Run parses raw shell history text and returns alias proposals ranked by
// keystrokes saved. It contains no syscall/js dependency so it can be
// exercised by ordinary Go tests; cmd/wasm is a thin wrapper around it.
func Run(historyText string) []alias.Proposal {
	commands := history.Parse(strings.NewReader(historyText))
	candidates := miner.Mine(commands)
	return alias.Propose(candidates)
}
