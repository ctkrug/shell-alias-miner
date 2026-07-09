package alias

import (
	"testing"

	"github.com/ctkrug/shell-alias-miner/internal/miner"
)

func TestProposeGeneratesUniqueNames(t *testing.T) {
	candidates := []miner.Candidate{
		{Command: "git status --short", Count: 340},
		{Command: "git stash pop", Count: 12},
	}

	got := Propose(candidates)

	if len(got) != 2 {
		t.Fatalf("got %d proposals, want 2", len(got))
	}
	if got[0].Name == got[1].Name {
		t.Errorf("expected unique alias names, got %q twice", got[0].Name)
	}
}

func TestProposeDefinitionIsPasteable(t *testing.T) {
	candidates := []miner.Candidate{{Command: "git status --short", Count: 340}}
	got := Propose(candidates)

	want := `alias gs="git status --short"`
	if got[0].Definition != want {
		t.Errorf("Definition = %q, want %q", got[0].Definition, want)
	}
}

func TestProposeKeystrokesSavedScalesWithOccurrences(t *testing.T) {
	few := Propose([]miner.Candidate{{Command: "git status --short", Count: 1}})
	many := Propose([]miner.Candidate{{Command: "git status --short", Count: 340}})

	if many[0].KeystrokesSaved <= few[0].KeystrokesSaved {
		t.Errorf("expected KeystrokesSaved to grow with occurrences: few=%d many=%d",
			few[0].KeystrokesSaved, many[0].KeystrokesSaved)
	}
	if many[0].KeystrokesSaved != few[0].KeystrokesSaved*340 {
		t.Errorf("KeystrokesSaved = %d, want %d", many[0].KeystrokesSaved, few[0].KeystrokesSaved*340)
	}
}

func TestProposeManyCollisionsStayUnique(t *testing.T) {
	// All of these derive the same base name ("eu") from their first two
	// fields, forcing a long run of collisions in the numbered-suffix
	// resolver. Also guards against the O(n) probe-per-collision approach
	// that made Propose quadratic on inputs like this.
	candidates := make([]miner.Candidate, 2000)
	for i := range candidates {
		candidates[i] = miner.Candidate{Command: "echo unique-command", Count: 1}
	}

	got := Propose(candidates)

	seen := make(map[string]bool, len(got))
	for _, p := range got {
		if seen[p.Name] {
			t.Fatalf("duplicate alias name %q", p.Name)
		}
		seen[p.Name] = true
	}
}

func TestProposeNeverGoesNegative(t *testing.T) {
	// A short, already-terse command where the alias definition is longer
	// than the command itself must not report negative savings.
	got := Propose([]miner.Candidate{{Command: "ls", Count: 5}})
	if got[0].KeystrokesSaved < 0 {
		t.Errorf("KeystrokesSaved = %d, want >= 0", got[0].KeystrokesSaved)
	}
}
