package history

import (
	"os"
	"reflect"
	"testing"
)

// Both fixtures encode the same 11-command sequence, one in plain
// bash_history form and one in zsh's EXTENDED_HISTORY form. Parsing either
// must produce the identical command list, with no timestamp artifacts
// leaking through in the zsh case.
func TestParseFixturesAgree(t *testing.T) {
	want := []string{
		"ls -la",
		"cd ~/projects",
		"git status",
		"git status",
		"git add -A",
		`git commit -m "wip"`,
		"docker ps -a",
		"git status",
		"npm install",
		"npm run build",
		"git status",
	}

	for _, name := range []string{"bash_history", "zsh_extended_history"} {
		f, err := os.Open("../../testdata/" + name)
		if err != nil {
			t.Fatalf("open %s: %v", name, err)
		}
		defer f.Close()

		got := Parse(f)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Parse(%s) = %#v, want %#v", name, got, want)
		}
	}
}
