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

// A line that happens to start with ": " but isn't actually a zsh
// EXTENDED_HISTORY entry (no semicolon at all, or a semicolon but no
// timestamp-shaped header before it) must pass through unchanged rather
// than being mistaken for one and mangled.
func TestParseLineStartingWithColonSpaceButNotExtendedHistory(t *testing.T) {
	input := ": no semicolon at all here\n: notatimestamp;echo hi\n"
	got := Parse(strings.NewReader(input))
	want := []string{": no semicolon at all here", ": notatimestamp;echo hi"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse() = %#v, want %#v", got, want)
	}
}

// A single pasted line far longer than bufio.Scanner's default token size
// must not silently drop every command after it — a history file with one
// absurd line (a base64 blob someone once piped through echo) should still
// yield every other command.
func TestParseSurvivesLineLongerThanScannerBuffer(t *testing.T) {
	huge := strings.Repeat("a", 2*1024*1024)
	input := "git status\n" + huge + "\ngit log\n"
	got := Parse(strings.NewReader(input))
	want := []string{"git status", huge, "git log"}

	if !reflect.DeepEqual(got, want) {
		gotLens := make([]int, len(got))
		for i, c := range got {
			gotLens[i] = len(c)
		}
		t.Fatalf("Parse() returned %d commands with lengths %v, want 3 commands (lengths [10 %d 7])", len(got), gotLens, len(huge))
	}
}
