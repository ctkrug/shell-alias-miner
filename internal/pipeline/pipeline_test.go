package pipeline

import (
	"strings"
	"testing"
)

func TestRunEndToEnd(t *testing.T) {
	history := ""
	for i := 0; i < 340; i++ {
		history += "git status --short --branch\n"
	}
	history += "ls -la\n"

	got := Run(history)

	if len(got) == 0 {
		t.Fatal("Run() returned no proposals")
	}
	top := got[0]
	if top.Command != "git status --short --branch" {
		t.Errorf("top proposal Command = %q, want the 340x repeated command", top.Command)
	}
	if top.Occurrences != 340 {
		t.Errorf("top proposal Occurrences = %d, want 340", top.Occurrences)
	}
	if top.KeystrokesSaved <= 0 {
		t.Errorf("top proposal KeystrokesSaved = %d, want > 0", top.KeystrokesSaved)
	}
}

func TestRunEmptyHistory(t *testing.T) {
	got := Run("")
	if len(got) != 0 {
		t.Errorf("Run(\"\") returned %d proposals, want 0", len(got))
	}
}

func TestRunCollapsesVaryingCommitMessagesIntoOneFunction(t *testing.T) {
	var history strings.Builder
	for i := 0; i < 50; i++ {
		history.WriteString(`git commit -m "message ` + string(rune('a'+i%26)) + "\"\n")
	}

	got := Run(history.String())

	if len(got) != 1 {
		t.Fatalf("got %d proposals, want 1: %#v", len(got), got)
	}
	top := got[0]
	if top.Kind != "function" {
		t.Errorf("top proposal Kind = %q, want %q", top.Kind, "function")
	}
	if top.Occurrences != 50 {
		t.Errorf("top proposal Occurrences = %d, want 50", top.Occurrences)
	}
	if !strings.Contains(top.Definition, `"$1"`) {
		t.Errorf("top proposal Definition = %q, want it to reference \"$1\"", top.Definition)
	}
}
