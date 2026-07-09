package history

import (
	"reflect"
	"strings"
	"testing"
)

func TestParsePlainBashHistory(t *testing.T) {
	input := "git status\n\nls -la\ngit commit -m \"wip\"\n"
	got := Parse(strings.NewReader(input))
	want := []string{"git status", "ls -la", `git commit -m "wip"`}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse() = %#v, want %#v", got, want)
	}
}

func TestParseExtendedHistory(t *testing.T) {
	input := ": 1700000000:0;git status\n: 1700000010:2;git commit -m \"wip\"\n"
	got := Parse(strings.NewReader(input))
	want := []string{"git status", `git commit -m "wip"`}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse() = %#v, want %#v", got, want)
	}
}

func TestParseLineContinuation(t *testing.T) {
	input := "docker run \\\n  --rm -it \\\n  ubuntu\ngit status\n"
	got := Parse(strings.NewReader(input))

	if len(got) != 2 {
		t.Fatalf("Parse() returned %d commands, want 2: %#v", len(got), got)
	}
	if !strings.Contains(got[0], "docker run") || !strings.Contains(got[0], "ubuntu") {
		t.Errorf("continuation not joined into one command: %q", got[0])
	}
	if got[1] != "git status" {
		t.Errorf("got[1] = %q, want %q", got[1], "git status")
	}
}

func TestParseSkipsBlankLines(t *testing.T) {
	input := "\n\nls\n\n\ncd /tmp\n\n"
	got := Parse(strings.NewReader(input))
	want := []string{"ls", "cd /tmp"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse() = %#v, want %#v", got, want)
	}
}
