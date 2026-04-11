package shellshape

func init() {
	Register("ssh", handleSsh)
}

// handleSsh handles the ssh command.
// The first positional becomes <host> (hostnames, user@host, IPs).
// Remaining positionals after the host are remote command tokens, collapsed to <cmd>.
// Flag arguments are classified by type: paths, ports, strings.
func handleSsh(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a file path.
	pathFlags := map[string]bool{
		"-i": true, "-F": true, "-E": true, "-S": true, "-I": true,
	}
	// Flags whose next token is a port number.
	portFlags := map[string]bool{
		"-p": true,
	}
	// Flags whose next token is a host/destination.
	hostFlags := map[string]bool{
		"-J": true,
	}
	// Flags whose next token is an opaque string argument.
	strFlags := map[string]bool{
		"-l": true, "-o": true, "-c": true, "-m": true, "-e": true,
		"-B": true, "-b": true, "-O": true, "-Q": true, "-P": true,
		"-w": true, "-W": true, "-D": true, "-L": true, "-R": true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{portFlags, "N"},
		{hostFlags, "<jump-host>"},
		{strFlags, "<str>"},
	}

	var result []string
	hostAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			if !hostAssigned {
				hostAssigned = true
			}
			i++
			continue
		}

		// Once host is assigned, everything remaining is the remote command.
		if hostAssigned {
			result = append(result, "<cmd>")
			i++
			continue
		}

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional is the destination host.
		result = append(result, "<host>")
		hostAssigned = true
		i++
	}

	result = append(result, redirects...)
	return result
}
