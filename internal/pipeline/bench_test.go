package pipeline

import (
	"strconv"
	"strings"
	"testing"
)

// syntheticHistory builds an n-line history mixing a handful of repeated
// commands with unique ones, roughly approximating a real user's history
// shape. Used to keep an eye on Run's performance as history size grows
// toward the 100k+ line files targeted by backlog story 3.2.
func syntheticHistory(lines int) string {
	var b strings.Builder
	repeated := []string{
		"git status --short --branch",
		"docker compose up -d",
		"git log --oneline --graph --all",
	}
	for i := 0; i < lines; i++ {
		if i%10 == 0 {
			b.WriteString(repeated[i%len(repeated)])
		} else {
			b.WriteString("echo unique-command-" + strconv.Itoa(i))
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func BenchmarkRun10k(b *testing.B) {
	history := syntheticHistory(10_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Run(history)
	}
}

func BenchmarkRun100k(b *testing.B) {
	history := syntheticHistory(100_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Run(history)
	}
}
