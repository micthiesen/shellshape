package shellshape

func init() {
	Register("nmap", handleNmap)
}

// handleNmap handles the nmap command.
// Port flags (-p) collapse the next arg to <ports>.
// Output flags (-oN, -oX, -oG, -oS, -oA) collapse the next arg to <path>.
// Input/exclude file flags (-iL, --excludefile, --script-args-file) collapse the next arg to <path>.
// Script flag (--script) collapses the next arg to <script>.
// Numeric flags (--top-ports, --min-rate, --max-rate, timing/parallelism, -g, etc.) collapse the next arg to N.
// Value flags (-D, -e, --spoof-mac, --dns-servers, --script-args) collapse the next arg to <val>.
// Host flags (-S, --exclude) collapse the next arg to <host>.
// All positionals (targets) become <host>.
func handleNmap(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	portFlags := map[string]bool{
		"-p": true,
	}

	pathFlags := map[string]bool{
		"-oN": true, "-oX": true, "-oG": true, "-oS": true, "-oA": true,
		"-iL": true, "--excludefile": true, "--script-args-file": true,
	}

	scriptFlags := map[string]bool{
		"--script": true,
	}

	numericFlags := map[string]bool{
		"--top-ports": true, "--port-ratio": true,
		"--min-rate": true, "--max-rate": true,
		"--min-parallelism": true, "--max-parallelism": true,
		"--host-timeout": true, "--max-retries": true,
		"--scan-delay": true, "--max-scan-delay": true,
		"--min-hostgroup": true, "--max-hostgroup": true,
		"--min-rtt-timeout": true, "--max-rtt-timeout": true,
		"--initial-rtt-timeout": true,
		"-g":                    true, "--source-port": true,
		"--data-length": true, "--mtu": true,
		"--ttl": true, "--max-os-tries": true,
		"--version-intensity": true,
	}

	valFlags := map[string]bool{
		"-D": true, "-e": true,
		"--spoof-mac": true, "--dns-servers": true,
		"--script-args": true, "--datadir": true,
		"--stylesheet": true, "--webxml": true,
		"--proxies": true,
	}

	hostFlags := map[string]bool{
		"-S": true, "--exclude": true,
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

		if portFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<ports>")
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if scriptFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<script>")
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if hostFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<host>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: target host/IP/range
		result = append(result, "<host>")
		i++
	}

	result = append(result, redirects...)
	return result
}
