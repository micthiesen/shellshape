package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("btrfs", handleBtrfs, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleBtrfs handles btrfs filesystem management commands.
// btrfs has two-level subcommands (e.g. btrfs subvolume list, btrfs filesystem df).
// The first level subcommand is consumed by the framework; the second level
// (sub-subcommand) is kept verbatim as a structural token.
// Path flags (-p, -f, -c for send; -r for receive) consume a path argument.
// Numeric positionals → N, path positionals → <path>.
// Sub-subcommand names are kept structural.
func handleBtrfs(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags that consume a path argument
	pathFlags := map[string]bool{
		"-p": true, "--parent": true,
		"-c": true, "--clone-src": true,
		"-f": true, "--file": true,
		"-o": true, "--output": true,
	}

	// Flags that consume a numeric argument
	numericFlags := map[string]bool{
		"-i": true, "--subvolid": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
	}

	// Known sub-subcommands per subcommand (kept structural)
	subSubcommands := map[string]map[string]bool{
		"subvolume": {
			"create": true, "delete": true, "list": true, "snapshot": true,
			"show": true, "get-default": true, "set-default": true, "find-new": true,
			"sync": true,
		},
		"filesystem": {
			"df": true, "du": true, "show": true, "sync": true, "defragment": true,
			"resize": true, "label": true, "usage": true, "mkswapfile": true,
		},
		"device": {
			"add": true, "delete": true, "remove": true, "scan": true,
			"ready": true, "stats": true, "usage": true,
		},
		"balance": {
			"start": true, "pause": true, "cancel": true, "resume": true, "status": true,
		},
		"scrub": {
			"start": true, "cancel": true, "resume": true, "status": true,
		},
		"property": {
			"get": true, "set": true, "list": true,
		},
		"quota": {
			"enable": true, "disable": true, "rescan": true,
		},
		"qgroup": {
			"assign": true, "remove": true, "create": true, "destroy": true,
			"show": true, "limit": true,
		},
		"replace": {
			"start": true, "status": true, "cancel": true,
		},
		"rescue": {
			"chunk-recover": true, "fix-device-size": true, "super-recover": true,
			"zero-log": true, "create-control-device": true,
		},
		"inspect-internal": {
			"inode-resolve": true, "logical-resolve": true, "subvolid-resolve": true,
			"rootid": true, "min-dev-size": true, "dump-tree": true, "dump-super": true,
			"tree-stats": true,
		},
	}

	var result []string
	i := 0

	// Check for sub-subcommand as first positional
	if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
		tok := args[i]
		if subs, ok := subSubcommands[subcommand]; ok && subs[tok] {
			result = append(result, tok)
			i++
		}
	}

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

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

		// Positional: numbers → N, paths → <path>, otherwise classify
		if shellshape.NumberRE.MatchString(tok) {
			result = append(result, "N")
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
