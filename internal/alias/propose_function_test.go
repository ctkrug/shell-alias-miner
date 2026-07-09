package alias

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/ctkrug/shell-alias-miner/internal/miner"
)

func TestProposeFunctionDefinitionIsFunctionShaped(t *testing.T) {
	candidates := []miner.Candidate{
		{Command: `git commit -m "fix bug"`, Prefix: "git commit -m", Count: 50, Kind: miner.KindTemplate},
	}

	got := Propose(candidates)

	if got[0].Kind != KindFunction {
		t.Errorf("Kind = %q, want %q", got[0].Kind, KindFunction)
	}
	if !strings.HasPrefix(got[0].Definition, "function ") {
		t.Errorf("Definition = %q, want it to start with \"function \"", got[0].Definition)
	}
	if !strings.Contains(got[0].Definition, `"$1"`) {
		t.Errorf("Definition = %q, want it to reference \"$1\"", got[0].Definition)
	}
}

func TestProposeTemplateKeystrokesSavedScalesWithOccurrences(t *testing.T) {
	few := Propose([]miner.Candidate{{Prefix: "git commit -m", Command: `git commit -m "a"`, Count: 1, Kind: miner.KindTemplate}})
	many := Propose([]miner.Candidate{{Prefix: "git commit -m", Command: `git commit -m "a"`, Count: 50, Kind: miner.KindTemplate}})

	if many[0].KeystrokesSaved != few[0].KeystrokesSaved*50 {
		t.Errorf("KeystrokesSaved = %d, want %d", many[0].KeystrokesSaved, few[0].KeystrokesSaved*50)
	}
}

func TestProposeTemplateNeverGoesNegative(t *testing.T) {
	// A one-token prefix where the function name ends up no shorter than
	// the prefix itself must not report negative savings.
	got := Propose([]miner.Candidate{{Prefix: "cd", Command: "cd /tmp", Count: 5, Kind: miner.KindTemplate}})
	if got[0].KeystrokesSaved < 0 {
		t.Errorf("KeystrokesSaved = %d, want >= 0", got[0].KeystrokesSaved)
	}
}

// TestProposeFunctionDefinitionReproducesCommandInBash verifies acceptance
// criterion 1.3: sourcing the generated function in a real shell and
// calling it with an argument must reproduce the original command with
// that argument substituted in the correct position. It uses an "echo"
// stand-in prefix rather than a real git/docker command so the test has no
// side effects.
func TestProposeFunctionDefinitionReproducesCommandInBash(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}

	candidates := []miner.Candidate{
		{Command: `echo tagged "first release"`, Prefix: "echo tagged", Count: 2, Kind: miner.KindTemplate},
	}
	proposal := Propose(candidates)[0]

	script := proposal.Definition + "\n" + proposal.Name + ` "second release"`
	out, err := exec.Command("bash", "-c", script).CombinedOutput()
	if err != nil {
		t.Fatalf("bash -c %q failed: %v\noutput: %s", script, err, out)
	}

	got := strings.TrimSpace(string(out))
	want := "tagged second release"
	if got != want {
		t.Errorf("function call output = %q, want %q", got, want)
	}
}
