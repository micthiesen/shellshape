package handlers

import (
	"strings"

	shellshape "github.com/micthiesen/shellshape"
)

func init() {
	shellshape.Register("kopia", handleKopia, shellshape.HandlerOptions{HasSubcommands: true})
}

// Known kopia sub-subcommands that should stay verbatim.
var kopiaSubSubcommands = map[string]bool{
	"create": true, "list": true, "restore": true, "delete": true,
	"connect": true, "disconnect": true, "status": true, "set": true,
	"start": true, "stop": true, "refresh": true, "show": true,
	"estimate": true, "expire": true, "fix": true, "migrate": true,
	"verify": true, "get": true, "remove": true, "add": true,
	"info": true, "ls": true, "cat": true,
}

// Known structural tokens (backend names, etc.) that stay verbatim.
var kopiaStructuralTokens = map[string]bool{
	"filesystem": true, "s3": true, "b2": true, "gcs": true,
	"azure": true, "sftp": true, "webdav": true, "rclone": true,
	"all": false, // "all" for mount is a positional, not structural
}

// handleKopia handles the kopia backup tool.
func handleKopia(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"--path":            true,
		"--config-file":     true,
		"--cache-directory": true,
		"--log-dir":         true,
		"--source":          true,
	}

	numericFlags := map[string]bool{
		"--keep-latest":        true,
		"--keep-hourly":        true,
		"--keep-daily":         true,
		"--keep-weekly":        true,
		"--keep-monthly":       true,
		"--keep-annual":        true,
		"--parallel":           true,
		"--max-upload-speed":   true,
		"--max-download-speed": true,
	}

	valFlags := map[string]bool{
		"--bucket":            true,
		"--access-key":        true,
		"--secret-access-key": true,
		"--region":            true,
		"--endpoint":          true,
		"--prefix":            true,
		"--password":          true,
		"--address":           true,
		"--add-ignore":        true,
		"--description":       true,
		"--override-hostname": true,
		"--override-username": true,
		"--token":             true,
		"--server-url":        true,
		"--container":         true,
		"--storage-account":   true,
		"--compression":       true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	subSubSeen := false
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Fused --flag=value
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Sub-subcommand: keep verbatim (create, list, restore, etc.)
		if !subSubSeen && kopiaSubSubcommands[tok] {
			result = append(result, tok)
			subSubSeen = true
			i++
			continue
		}

		// Backend type names after sub-subcommand: keep verbatim
		if kopiaStructuralTokens[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: paths stay <path>, everything else → <val>
		result = append(result, kopiaClassifyPositional(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func kopiaClassifyPositional(tok string) string {
	classified := shellshape.ClassifyToken(tok)
	if classified == "<path>" || strings.HasSuffix(classified, "-uri>") {
		return classified
	}
	return "<val>"
}
