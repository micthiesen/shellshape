package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("rsync", handleRsync)
}

// handleRsync handles rsync commands.
// -e/--rsh consumes the next arg as <rsh> (remote shell).
// -f/--filter, --exclude, --include consume the next arg as <filter>.
// --exclude-from, --include-from, --files-from consume the next arg as <path>.
// --backup-dir, --compare-dest, --copy-dest, --link-dest, --partial-dir,
// --log-file, --password-file, --read-batch, --write-batch, --only-write-batch,
// --rsync-path, -T/--temp-dir consume the next arg as <path>.
// All other positionals are classified as paths (including remote host:path specs).
func handleRsync(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	rshFlags := map[string]bool{"-e": true, "--rsh": true}
	filterFlags := map[string]bool{
		"-f": true, "--filter": true,
		"--exclude": true, "--include": true,
	}
	pathFlags := map[string]bool{
		"--exclude-from": true, "--include-from": true, "--files-from": true,
		"--backup-dir": true, "--compare-dest": true, "--copy-dest": true,
		"--link-dest": true, "--partial-dir": true, "--log-file": true,
		"--password-file": true, "--read-batch": true, "--write-batch": true,
		"--only-write-batch": true, "--rsync-path": true,
		"-T": true, "--temp-dir": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: rshFlags, Placeholder: "<rsh>"},
		{Flags: filterFlags, Placeholder: "<filter>"},
		{Flags: pathFlags, Placeholder: "<path>"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional: classify, but treat rsync remote specs (host:path) as <path>.
		result = append(result, classifyRsyncPath(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// classifyRsyncPath classifies a token as a path, handling rsync remote
// specs like "user@host:/path" or "host:/path" which would otherwise be
// classified as <git-uri> or left verbatim.
func classifyRsyncPath(tok string) string {
	// Let URLs (rsync://, ssh://, etc.) pass through to classifyToken.
	if strings.Contains(tok, "://") {
		return shellshape.ClassifyToken(tok)
	}
	// Remote spec: contains ":" with path component (host:path or user@host:path)
	if idx := strings.Index(tok, ":"); idx > 0 {
		return "<path>"
	}
	return shellshape.ClassifyToken(tok)
}
