package miner

import "testing"

func TestCountFrequenciesOrdersByCountThenText(t *testing.T) {
	commands := []string{
		"ls -la",
		"git status",
		"git status",
		"ls -la",
		"git status",
		"cd /tmp",
	}

	got := CountFrequencies(commands)

	if len(got) != 3 {
		t.Fatalf("got %d candidates, want 3: %#v", len(got), got)
	}
	if got[0].Command != "git status" || got[0].Count != 3 {
		t.Errorf("got[0] = %#v, want {git status 3}", got[0])
	}
	if got[1].Command != "ls -la" || got[1].Count != 2 {
		t.Errorf("got[1] = %#v, want {ls -la 2}", got[1])
	}
	if got[2].Command != "cd /tmp" || got[2].Count != 1 {
		t.Errorf("got[2] = %#v, want {cd /tmp 1}", got[2])
	}
}

func TestCountFrequenciesIgnoresBlankLines(t *testing.T) {
	got := CountFrequencies([]string{"", "ls", "", "ls"})
	if len(got) != 1 {
		t.Fatalf("got %d candidates, want 1: %#v", len(got), got)
	}
	if got[0].Count != 2 {
		t.Errorf("got[0].Count = %d, want 2", got[0].Count)
	}
}

func TestCountFrequenciesEmptyInput(t *testing.T) {
	got := CountFrequencies(nil)
	if len(got) != 0 {
		t.Fatalf("got %d candidates, want 0", len(got))
	}
}
