package miner

import (
	"strconv"
	"testing"
)

// syntheticSharedPrefix builds history where 90% of lines are one-off
// commands sharing a single leading token ("echo unique-command-N"). This
// is Mine's worst case: every one of those commands lands in the same
// trailing-argument group, so the group's suffix set grows to roughly the
// full history size. Guards against the algorithm regressing back toward
// O(n) per-command map churn on the shared "covered" bookkeeping — see
// docs/BACKLOG.md story 3.2 (150k lines must mine in well under 3s even in
// wasm, which runs several times slower than native).
func syntheticSharedPrefix(n int) []string {
	repeated := []string{
		"git status --short --branch",
		"docker compose up -d",
		"git log --oneline --graph --all",
	}
	commands := make([]string, 0, n)
	for i := 0; i < n; i++ {
		if i%10 == 0 {
			commands = append(commands, repeated[i%len(repeated)])
		} else {
			commands = append(commands, "echo unique-command-"+strconv.Itoa(i))
		}
	}
	return commands
}

func BenchmarkMine150kSharedPrefix(b *testing.B) {
	commands := syntheticSharedPrefix(150_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Mine(commands)
	}
}
