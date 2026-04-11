package shellshape

func init() {
	Register("nft", handleNft)
}

// handleNft handles nft (nftables) commands.
// nft's grammar is: nft [flags] <command> <object-type> [<family>] <table> [<chain>] [<rule-spec>...]
//
// Commands (add, delete, list, flush, insert, replace, create, export, etc.)
// are kept verbatim. Object types (table, chain, rule, set, etc.) are structural.
// Address families (ip, ip6, inet, arp, bridge, netdev) are structural.
// Table and chain names are data, collapsed to <word>.
// Rule specifications (everything after chain name for rule commands) collapse to <val>.
// The "handle" keyword in delete operations is structural, followed by a numeric N.
func handleNft(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-f": true, "--file": true,
		"-I": true, "--includepath": true,
	}

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"-D": true, "--define": true,
	}

	// nft commands: kept verbatim.
	commands := map[string]bool{
		"add": true, "create": true, "delete": true, "list": true,
		"flush": true, "rename": true, "insert": true, "replace": true,
		"export": true, "monitor": true, "describe": true, "import": true,
		"reset": true,
	}

	// Object types: kept verbatim.
	objectTypes := map[string]bool{
		"table": true, "chain": true, "rule": true,
		"set": true, "map": true, "element": true,
		"flowtable": true, "counter": true, "quota": true,
		"ct": true, "limit": true, "meter": true,
		"secmark": true, "synproxy": true, "ruleset": true,
		"tables": true, "chains": true, "sets": true, "maps": true,
		"counters": true, "quotas": true, "meters": true,
	}

	// Address families are structural.
	families := map[string]bool{
		"ip": true, "ip6": true, "inet": true,
		"arp": true, "bridge": true, "netdev": true,
	}

	// Objects that have trailing rule body or data after names.
	bodyObjects := map[string]bool{
		"rule": true, "element": true,
	}

	// How many name positionals to expect after family for each object type.
	nameCount := map[string]int{
		"table": 1, "chain": 2, "rule": 2, "element": 2,
		"set": 2, "map": 2, "flowtable": 2,
		"counter": 2, "quota": 2, "ct": 2,
		"limit": 2, "meter": 2, "secmark": 2, "synproxy": 2,
		"ruleset": 0,
	}

	// Commands where all positionals after command are structural.
	verbatimCommands := map[string]bool{
		"export": true, "monitor": true, "describe": true,
	}

	var result []string
	command := ""
	objectType := ""
	familySeen := false
	namesCollected := 0
	inBody := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Path flags
		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		// Value flags
		if valFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Boolean flags
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Command word (first non-flag positional)
		if command == "" {
			if commands[tok] {
				command = tok
				result = append(result, tok)
				i++
				continue
			}
			// Unknown command: keep verbatim and fall through
			result = append(result, tok)
			i++
			continue
		}

		// For verbatim commands (export, monitor, describe), keep all tokens as-is.
		if verbatimCommands[command] {
			result = append(result, tok)
			i++
			continue
		}

		// If we're in rule/element body, collapse remaining positionals.
		if inBody {
			// Special case: "handle" keyword followed by a number.
			if tok == "handle" {
				result = append(result, "handle")
				i++
				if i < len(args) {
					result = append(result, "N")
					i++
				}
				continue
			}
			// Consume all remaining non-flag, non-subshell positionals into <val>.
			hasBody := false
			for i < len(args) {
				t := args[i]
				if isSubshellToken(t) {
					if hasBody {
						result = append(result, "<val>")
						hasBody = false
					}
					result = append(result, t)
					i++
					continue
				}
				if isFlagToken(t) {
					break
				}
				hasBody = true
				i++
			}
			if hasBody {
				result = append(result, "<val>")
			}
			continue
		}

		// Object type
		if objectType == "" && objectTypes[tok] {
			objectType = tok
			result = append(result, tok)
			i++
			continue
		}

		// After object type: check for address family
		if objectType != "" && !familySeen {
			if families[tok] {
				familySeen = true
				result = append(result, tok)
				i++
				continue
			}
			// No family specified, proceed to name collection.
			familySeen = true
		}

		// Collect name positionals (table, chain, set names)
		maxNames := nameCount[objectType]
		if namesCollected < maxNames {
			result = append(result, "<word>")
			namesCollected++
			i++
			if namesCollected == maxNames && bodyObjects[objectType] {
				inBody = true
			}
			continue
		}

		// Fallback: classify remaining positionals
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
