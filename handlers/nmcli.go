package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("nmcli", handleNmcli, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleNmcli handles nmcli (NetworkManager CLI).
// Since HasSubcommands is set, the normalizer extracts the first token as
// the subcommand. nmcli supports global flags before the subcommand (-t, -p,
// -f, etc.), so the handler detects when a flag was consumed as the subcommand
// and emits the real subcommand from the remaining tokens.
//
// Sub-subcommands (show, up, down, add, modify, delete, wifi, connect, etc.)
// are kept structural. Keyword-value pairs (con-name, ifname, password, type,
// file, ssid, iface) keep the keyword and collapse the value.
// Connection names, UUIDs, interface names, SSIDs, passwords -> <val>.
func handleNmcli(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Global boolean flags (no argument consumed).
	globalBoolFlags := map[string]bool{
		"-t": true, "--terse": true,
		"-p": true, "--pretty": true,
		"-e": true,
	}

	// Global flags that consume the next token as <val>.
	globalValFlags := map[string]bool{
		"-f": true, "--fields": true,
		"-m": true, "--mode": true,
		"--colors": true,
		"--wait":   true,
		"-w":       true,
	}

	// nmcli sub-subcommands that are structural (kept verbatim).
	structuralTokens := map[string]bool{
		"show": true, "up": true, "down": true,
		"add": true, "modify": true, "delete": true,
		"status": true, "list": true, "rescan": true,
		"wifi": true, "connect": true, "disconnect": true,
		"hotspot": true, "import": true, "export": true,
		"reload": true, "load": true, "clone": true,
		"edit": true, "monitor": true,
		"on": true, "off": true,
		"secret": true, "polkit": true, "all": true,
		"connectivity": true, "hostname": true, "permissions": true, "logging": true,
		"wwan": true,
	}

	// Keyword arguments that consume the next token as <val>.
	valKeywords := map[string]bool{
		"con-name": true, "ifname": true, "password": true,
		"ssid": true, "iface": true,
	}

	// Keyword arguments that consume the next token verbatim (structural).
	verbatimKeywords := map[string]bool{
		"type": true,
	}

	// Keyword arguments that consume the next token as <path>.
	pathKeywords := map[string]bool{
		"file": true,
	}

	leadingCategories := []shellshape.FlagCategory{
		{Flags: globalValFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0

	// Handle a leading flag that the normalizer extracted as the subcommand.
	// Boolean flags (globalBoolFlags) don't need a category entry; the
	// helper keeps them verbatim as "unknown flags".
	result, subcommand, i = shellshape.RepairLeadingFlagSubcommand(
		subcommand, args, i, result, leadingCategories, nil,
	)

	// For connection modify: after the sub-subcommand, first positional is the
	// connection name (<val>), then everything else is property-value pairs (all <val>).
	isModify := false
	if subcommand == "connection" || subcommand == "con" || subcommand == "c" {
		// Check if next structural token is "modify"
		if i < len(args) && args[i] == "modify" {
			isModify = true
		}
	}

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Boolean flags
		if globalBoolFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Global val flags that appear after subcommand
		if globalValFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Other flags: keep verbatim (e.g. --active, --show-secrets)
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Structural tokens (sub-subcommands, on/off, wifi, etc.)
		if structuralTokens[tok] {
			if tok == "modify" {
				isModify = true
			}
			result = append(result, tok)
			i++
			continue
		}

		// Keyword-value pairs
		if valKeywords[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if verbatimKeywords[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, args[i])
				}
				i++
			}
			continue
		}

		if pathKeywords[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		// For modify: all remaining positionals are <val> (conn name + property-value pairs)
		if isModify {
			result = append(result, "<val>")
			i++
			continue
		}

		// Default: positional -> <val> (connection names, interface names, UUIDs, SSIDs)
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
