package alias

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/ctkrug/shell-alias-miner/internal/miner"
)

// nameFor derives a name meant to be pasted directly into a shell as an
// alias/function name. It must never contain a character that would make
// the generated definition invalid — or, worse, break out of the
// definition's own quoting.
func TestNameForNeverEmitsShellMetacharacters(t *testing.T) {
	cases := []string{
		"true && false",
		"make test || echo fail",
		"ls; pwd",
		`curl -H "Content-Type: application/json" https://example.com`,
		`echo "hi there" foo`,
	}

	for _, c := range cases {
		name := nameFor(c)
		for _, r := range name {
			if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyz0123456789", r) {
				t.Errorf("nameFor(%q) = %q, contains disallowed character %q", c, name, r)
			}
		}
	}
}

// End-to-end: a candidate built from an operator-containing command must
// still produce a Definition that is itself syntactically valid shell
// (checked via bash -n, a parse-only pass — a plain "alias name=..." isn't
// usable later in the very same -c script regardless of name validity, so
// this checks the thing that actually matters here: the generated line
// parses cleanly).
func TestProposeDefinitionValidInBashForOperatorCommand(t *testing.T) {
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}

	got := Propose([]miner.Candidate{{Command: "true && false", Count: 5, Kind: miner.KindExact}})
	if len(got) != 1 {
		t.Fatalf("got %d proposals, want 1", len(got))
	}

	def := got[0].Definition
	out, err := exec.Command("bash", "-n", "-c", def).CombinedOutput()
	if err != nil {
		t.Fatalf("bash rejected generated definition %q: %v\n%s", def, err, out)
	}
}
