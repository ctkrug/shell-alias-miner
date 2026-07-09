package miner

import "testing"

func TestMineCollapsesVaryingTrailingArgumentIntoOneTemplate(t *testing.T) {
	var commands []string
	for i := 0; i < 50; i++ {
		commands = append(commands, `git commit -m "message `+string(rune('a'+i%26))+`"`)
	}

	got := Mine(commands)

	if len(got) != 1 {
		t.Fatalf("got %d candidates, want 1: %#v", len(got), got)
	}
	if got[0].Kind != KindTemplate {
		t.Errorf("Kind = %q, want %q", got[0].Kind, KindTemplate)
	}
	if got[0].Prefix != "git commit -m" {
		t.Errorf("Prefix = %q, want %q", got[0].Prefix, "git commit -m")
	}
	if got[0].Count != 50 {
		t.Errorf("Count = %d, want 50", got[0].Count)
	}
}

func TestMineNoVariationMatchesCountFrequencies(t *testing.T) {
	commands := []string{
		"git status --short",
		"git status --short",
		"ls -la",
		"git status --short",
		"cd /tmp",
	}

	mined := Mine(commands)
	baseline := CountFrequencies(commands)

	if len(mined) != len(baseline) {
		t.Fatalf("got %d candidates from Mine, want %d (baseline)", len(mined), len(baseline))
	}
	for i := range baseline {
		if mined[i].Command != baseline[i].Command || mined[i].Count != baseline[i].Count {
			t.Errorf("mined[%d] = {%q %d}, want {%q %d}",
				i, mined[i].Command, mined[i].Count, baseline[i].Command, baseline[i].Count)
		}
		if mined[i].Kind != KindExact {
			t.Errorf("mined[%d].Kind = %q, want %q", i, mined[i].Kind, KindExact)
		}
	}
}

func TestMineSingleTokenCommandsNeverBecomeTemplates(t *testing.T) {
	got := Mine([]string{"ls", "ls", "pwd"})

	for _, c := range got {
		if c.Kind == KindTemplate {
			t.Errorf("single-token command produced a template candidate: %#v", c)
		}
	}
}

func TestMineRequiresAtLeastTwoDistinctSuffixes(t *testing.T) {
	// Same prefix, same single suffix, repeated: no variation observed,
	// so this must stay an exact-match candidate, not a template.
	got := Mine([]string{`git commit -m "same"`, `git commit -m "same"`})

	if len(got) != 1 {
		t.Fatalf("got %d candidates, want 1: %#v", len(got), got)
	}
	if got[0].Kind != KindExact {
		t.Errorf("Kind = %q, want %q", got[0].Kind, KindExact)
	}
}

func TestMineEmptyInput(t *testing.T) {
	got := Mine(nil)
	if len(got) != 0 {
		t.Fatalf("got %d candidates, want 0", len(got))
	}
}

func TestMineMixesExactAndTemplateCandidates(t *testing.T) {
	var commands []string
	for i := 0; i < 10; i++ {
		commands = append(commands, `docker run -it image-`+string(rune('a'+i)))
	}
	for i := 0; i < 5; i++ {
		commands = append(commands, "git status --short")
	}

	got := Mine(commands)

	var sawTemplate, sawExact bool
	for _, c := range got {
		switch c.Kind {
		case KindTemplate:
			sawTemplate = true
			if c.Count != 10 {
				t.Errorf("template Count = %d, want 10", c.Count)
			}
		case KindExact:
			sawExact = true
			if c.Command != "git status --short" || c.Count != 5 {
				t.Errorf("exact candidate = %#v, want {git status --short 5}", c)
			}
		}
	}
	if !sawTemplate || !sawExact {
		t.Fatalf("expected both a template and an exact candidate, got %#v", got)
	}
}
