package shellshape

func init() {
	Register("lsof", handleLsof)
}

// handleLsof handles the lsof command.
// -i takes an optional network spec (collapse to <net-spec> if present).
// -p (PID), -u (user), -c (command name), -d (FD numbers) consume the next token.
// +D and -D consume a path argument.
// +c consumes a numeric argument.
// Positionals are file paths (classifyToken).
func handleLsof(_ string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags that consume the next token and their placeholder.
	valueFlags := map[string]string{
		"-p": "<pid>",
		"-u": "<user>",
		"-c": "<name>",
		"-d": "<fd>",
		"-D": "<path>",
		"+D": "<path>",
		"+d": "<path>",
		"-A": "<path>",
		"-k": "<path>",
		"-K": "<path>",
		"-s": "<val>",
		"-g": "<val>",
	}

	numericFlags := map[string]bool{
		"+c": true,
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

		// Handle -i: optional network spec argument.
		if tok == "-i" {
			result = append(result, tok)
			i++
			// Peek at next token: if it looks like a network spec, consume it.
			if i < len(args) && isNetworkSpec(args[i]) {
				result = append(result, "<net-spec>")
				i++
			}
			continue
		}

		// Flags that consume a value.
		if ph, ok := valueFlags[tok]; ok {
			result, i = consumeFlagArg(tok, args, i, result, ph)
			continue
		}

		// Numeric value flags.
		if numericFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		// Regular flags (boolean).
		if isFlagToken(tok) || tok == "+w" || tok == "+E" || tok == "+L" {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// isNetworkSpec returns true if the token looks like a lsof -i network spec.
// Network specs typically contain colons (e.g., :8080, TCP:80, @host:port)
// or start with protocol indicators (TCP, UDP, 4, 6).
func isNetworkSpec(tok string) bool {
	if isFlagToken(tok) || isSubshellToken(tok) {
		return false
	}
	// Starts with + (lsof flag like +D) - not a network spec.
	if len(tok) > 0 && tok[0] == '+' {
		return false
	}
	// Contains colon: :8080, TCP:*, @host:port
	if len(tok) > 0 && containsByte(tok, ':') {
		return true
	}
	// Bare protocol: TCP, UDP, tcp, udp
	upper := toUpper(tok)
	if upper == "TCP" || upper == "UDP" {
		return true
	}
	// Bare address family: 4, 6
	if tok == "4" || tok == "6" {
		return true
	}
	// Starts with @ (host spec)
	if len(tok) > 0 && tok[0] == '@' {
		return true
	}
	return false
}

func containsByte(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return true
		}
	}
	return false
}

func toUpper(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
