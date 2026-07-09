package miner

import (
	"reflect"
	"testing"
)

func TestTokenizeSplitsOnWhitespace(t *testing.T) {
	got := tokenize("git commit -m")
	want := []string{"git", "commit", "-m"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize() = %#v, want %#v", got, want)
	}
}

func TestTokenizeKeepsQuotedSpanAsOneToken(t *testing.T) {
	got := tokenize(`git commit -m "fix the bug"`)
	want := []string{"git", "commit", "-m", `"fix the bug"`}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize() = %#v, want %#v", got, want)
	}
}

func TestTokenizeKeepsSingleQuotedSpanAsOneToken(t *testing.T) {
	got := tokenize(`echo 'multi word arg'`)
	want := []string{"echo", "'multi word arg'"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize() = %#v, want %#v", got, want)
	}
}

func TestTokenizeEmptyInput(t *testing.T) {
	got := tokenize("")
	if len(got) != 0 {
		t.Errorf("tokenize(\"\") = %#v, want empty", got)
	}
}

func TestTokenizeCollapsesRepeatedWhitespace(t *testing.T) {
	got := tokenize("git   status")
	want := []string{"git", "status"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize() = %#v, want %#v", got, want)
	}
}

// An unterminated quote must not hang or panic — the rest of the line
// (including any embedded spaces) becomes one trailing token instead.
func TestTokenizeUnbalancedQuoteNeverPanics(t *testing.T) {
	got := tokenize(`git commit -m "oops never closed`)
	want := []string{"git", "commit", "-m", `"oops never closed`}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize() = %#v, want %#v", got, want)
	}
}

func TestTokenizeUnicodeAndEmoji(t *testing.T) {
	got := tokenize(`git commit -m "🎉 done"`)
	want := []string{"git", "commit", "-m", `"🎉 done"`}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tokenize() = %#v, want %#v", got, want)
	}
}
