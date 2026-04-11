package shellshape

func init() {
	Register("dpkg", handleDpkg)
}

// handleDpkg handles dpkg (Debian package manager) commands.
// dpkg uses action flags rather than subcommands to determine behavior.
//
// Actions where positionals are package names (structural, kept verbatim):
//
//	-r/--remove, -P/--purge, -L/--listfiles, -s/--status,
//	-p/--print-avail, --configure
//
// Actions where positionals are file paths (collapsed to <path>):
//
//	-i/--install, -c/--contents, --unpack, -b/--build
//
// Actions where positionals are search/glob patterns (collapsed to <pattern>):
//
//	-l/--list, -S/--search
//
// Flags like --admindir, --root, --instdir, --log consume a path argument.
func handleDpkg(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Action flags that determine how positionals are interpreted.
	pathActions := map[string]bool{
		"-i": true, "--install": true,
		"-c": true, "--contents": true,
		"--unpack": true,
		"-b":       true, "--build": true,
	}
	patternActions := map[string]bool{
		"-l": true, "--list": true,
		"-S": true, "--search": true,
	}
	structuralActions := map[string]bool{
		"-r": true, "--remove": true,
		"-P": true, "--purge": true,
		"-L": true, "--listfiles": true,
		"-s": true, "--status": true,
		"-p": true, "--print-avail": true,
		"--configure": true,
	}

	// Flags that consume the next token as a path.
	pathFlags := map[string]bool{
		"--admindir": true, "--instdir": true,
		"--root": true, "--log": true,
	}

	// First pass: determine the action mode from the flags present.
	mode := "" // "path", "pattern", "structural", or "" (unknown)
	for _, tok := range args {
		if pathActions[tok] {
			mode = "path"
			break
		}
		if patternActions[tok] {
			mode = "pattern"
			break
		}
		if structuralActions[tok] {
			mode = "structural"
			break
		}
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: classify based on the action mode.
		switch mode {
		case "path":
			result = append(result, classifyToken(tok))
			i++
		case "pattern":
			result = append(result, "<pattern>")
			i++
		case "structural":
			result = append(result, tok)
			i++
		default:
			result = append(result, classifyToken(tok))
			i++
		}
	}

	result = append(result, redirects...)
	return result
}
