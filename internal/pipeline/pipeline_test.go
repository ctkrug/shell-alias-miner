package pipeline

import "testing"

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
