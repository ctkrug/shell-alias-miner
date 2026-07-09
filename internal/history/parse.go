// Package history parses shell history files into a flat list of command
// lines, hiding the differences between the zsh and bash history formats.
package history

import (
	"bufio"
	"io"
	"strings"
)

// Parse reads a shell history file from r and returns each command line in
// the order it was run. It transparently handles:
//
//   - plain bash_history (one command per line)
//   - zsh's EXTENDED_HISTORY format (": <timestamp>:<duration>;<command>")
//   - zsh line continuations, where a trailing "\" joins the next line
//     into the same logical command
//
// Blank lines are skipped. Parse never fails on malformed input; it does
// its best with whatever it can read.
func Parse(r io.Reader) []string {
	var (
		commands []string
		pending  strings.Builder
		inCont   bool
	)

	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	flush := func() {
		cmd := strings.TrimSpace(pending.String())
		if cmd != "" {
			commands = append(commands, cmd)
		}
		pending.Reset()
	}

	for scanner.Scan() {
		line := scanner.Text()

		if !inCont {
			if cmd, ok := stripExtendedHistoryPrefix(line); ok {
				line = cmd
			}
		}

		if strings.HasSuffix(line, "\\") {
			pending.WriteString(strings.TrimSuffix(line, "\\"))
			pending.WriteByte('\n')
			inCont = true
			continue
		}

		pending.WriteString(line)
		flush()
		inCont = false
	}

	// A trailing continuation with no terminating line still counts.
	if inCont {
		flush()
	}

	return commands
}

// stripExtendedHistoryPrefix strips the ": <ts>:<duration>;" prefix that
// zsh writes when EXTENDED_HISTORY is enabled, returning the bare command
// and true if the prefix was present.
func stripExtendedHistoryPrefix(line string) (string, bool) {
	if !strings.HasPrefix(line, ": ") {
		return line, false
	}
	semi := strings.IndexByte(line, ';')
	if semi < 0 {
		return line, false
	}
	header := line[2:semi]
	if !strings.Contains(header, ":") {
		return line, false
	}
	return line[semi+1:], true
}
