package alias

import (
	"path/filepath"
	"regexp"
	"strings"
)

// secretFlagPattern matches common long-form flags that take a credential
// as their value, however it's spelled ("--password secret",
// "--password=secret", "--api-key=...").
var secretFlagPattern = regexp.MustCompile(`(?i)--(password|token|api[-_]?key|secret)\b`)

// inlineCredentialTools are commands with a well-known "-p<value>" (no
// space) convention for passing a password directly on the command line,
// as opposed to "-p" alone, which prompts interactively.
var inlineCredentialTools = map[string]bool{
	"mysql": true,
	"mongo": true,
}

// containsSecret reports whether command looks like it bakes in a
// credential, so proposing an alias or function for it would hand the user
// a snippet that leaks a password/token every time it's pasted somewhere.
func containsSecret(command string) bool {
	if secretFlagPattern.MatchString(command) {
		return true
	}

	fields := strings.Fields(command)
	if len(fields) == 0 {
		return false
	}
	// The history line may invoke the tool by an absolute/relative path
	// (PATH resolution order, a versioned wrapper) or in unusual case, so
	// match on the lowercased base name rather than the raw first field.
	tool := strings.ToLower(filepath.Base(fields[0]))
	if !inlineCredentialTools[tool] {
		return false
	}
	for _, f := range fields[1:] {
		if strings.HasPrefix(f, "-p") && f != "-p" {
			return true
		}
	}
	return false
}
