package shellshape

import "regexp"

func init() {
	Register("ip", handleIP)
}

var (
	ipv4RE = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(/\d{1,2})?$`)
	ipv6RE = regexp.MustCompile(`^[0-9a-fA-F]*:[0-9a-fA-F:]+(/\d{1,3})?$`)
	macRE  = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)
)

func looksLikeAddr(tok string) bool {
	return ipv4RE.MatchString(tok) || macRE.MatchString(tok) ||
		(len(tok) > 2 && ipv6RE.MatchString(tok))
}

// handleIP handles the ip (iproute2) command.
// ip has subcommands (addr, link, route, neigh, tunnel, etc.) and uses
// keyword-value pairs rather than traditional flags for most arguments.
//
// Global flags: -4, -6, -s, -d, -h, -j, -p, -t, -brief, -o, -oneline are boolean.
// -f/-family and -n/-netns consume the next token as <val>.
// Keyword args: dev, via, to, src, scope, label, table, proto, type consume <val> or <addr>.
// Numeric keywords: mtu, txqueuelen, metric, preference, priority consume N.
// The keyword "address" followed by a MAC address collapses to <addr>.
// IP/CIDR addresses and MAC addresses collapse to <addr>.
// Remaining positionals use classifyToken.
func handleIP(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Boolean global flags that look like numbers but should be kept verbatim.
	boolGlobalFlags := map[string]bool{
		"-4": true, "-6": true,
	}

	// Global flags that consume the next token as a value.
	globalValFlags := map[string]bool{
		"-f": true, "-family": true,
		"-n": true, "-netns": true,
	}

	// Keyword arguments that take a device/name/value.
	valKeywords := map[string]bool{
		"dev":   true,
		"via":   true,
		"to":    true,
		"src":   true,
		"scope": true,
		"label": true,
		"table": true,
		"proto": true,
		"type":  true,
	}

	// Keyword arguments that take a numeric value.
	numKeywords := map[string]bool{
		"mtu":         true,
		"txqueuelen":  true,
		"metric":      true,
		"preference":  true,
		"priority":    true,
		"advmss":      true,
		"window":      true,
		"rtt":         true,
		"mpu":         true,
		"cwnd":        true,
		"initcwnd":    true,
		"initrwnd":    true,
		"ssthresh":    true,
		"realms":      true,
		"rttvar":      true,
		"hoplimit":    true,
		"qdisc":       true,
		"numtxqueues": true,
		"numrxqueues": true,
	}

	// The "address" keyword takes a MAC or IP address.
	addrKeywords := map[string]bool{
		"address": true,
		"peer":    true,
		"local":   true,
		"remote":  true,
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

		if boolGlobalFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		if globalValFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
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

		// Keyword-value pairs
		if valKeywords[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else if looksLikeAddr(args[i]) {
					result = append(result, "<addr>")
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		if numKeywords[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "N")
				}
				i++
			}
			continue
		}

		if addrKeywords[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<addr>")
				}
				i++
			}
			continue
		}

		// Bare IP/CIDR or MAC address
		if looksLikeAddr(tok) {
			result = append(result, "<addr>")
			i++
			continue
		}

		// Positional: classify generically
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
