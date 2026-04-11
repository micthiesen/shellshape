package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"iptables", "ip6tables"} {
		shellshape.Register(name, handleIptables)
	}
}

// handleIptables handles iptables and ip6tables.
//
// Command flags (-A, -D, -I, -N, -X, -P, -L, -F, -Z) take a chain name that
// stays verbatim. After -D or -I, a bare number positional becomes N (rule
// number). Table (-t), protocol (-p), jump (-j), goto (-g), match (-m), and
// --reject-with keep their argument verbatim since those are structural.
// Address flags (-s, -d) collapse to <addr>. Interface flags (-i, -o) collapse
// to <val>. Port and match-module option flags collapse to <val> or <str>.
func handleIptables(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Command flags that consume a chain name (kept verbatim).
	chainFlags := map[string]bool{
		"-A": true, "--append": true,
		"-D": true, "--delete": true,
		"-I": true, "--insert": true,
		"-N": true, "--new-chain": true,
		"-X": true, "--delete-chain": true,
		"-P": true, "--policy": true,
		"-L": true, "--list": true,
		"-F": true, "--flush": true,
		"-Z": true, "--zero": true,
		"-C": true, "--check": true,
		"-R": true, "--replace": true,
		"-E": true, "--rename-chain": true,
	}

	// Flags whose argument is structural and stays verbatim.
	verbatimFlags := map[string]bool{
		"-t": true, "--table": true,
		"-j": true, "--jump": true,
		"-g": true, "--goto": true,
		"-p": true, "--protocol": true,
		"-m": true, "--match": true,
		"--reject-with": true,
	}

	// Flags whose argument is an address, collapsed to <addr>.
	addrFlags := map[string]bool{
		"-s": true, "--source": true, "--src": true,
		"-d": true, "--destination": true, "--dst": true,
	}

	// Flags whose argument is data, collapsed to <val>.
	valFlags := map[string]bool{
		"-i": true, "--in-interface": true,
		"-o": true, "--out-interface": true,
		"--dport": true, "--destination-port": true,
		"--sport": true, "--source-port": true,
		"--to-source": true, "--to-destination": true, "--to-ports": true,
		"--state": true, "--ctstate": true, "--ctstatus": true,
		"--limit": true, "--limit-burst": true,
		"--uid-owner": true, "--gid-owner": true,
		"--log-level": true,
		"--icmp-type": true, "--icmpv6-type": true,
		"--mac-source": true,
		"--set":        true,
		"--tcp-flags":  true,
		"--mark":       true, "--set-mark": true,
		"--tos": true, "--set-tos": true,
		"--ttl-set": true,
	}

	// Flags whose argument is a string, collapsed to <str>.
	strFlags := map[string]bool{
		"--comment":    true,
		"--log-prefix": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: addrFlags, Placeholder: "<addr>"},
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: strFlags, Placeholder: "<str>"},
	}

	// Track whether we just saw -D or -I (next positional number = rule number).
	expectRuleNum := false

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			expectRuleNum = false
			i++
			continue
		}

		if chainFlags[tok] {
			result = append(result, tok)
			// -D, -I, -R can have an optional rule number after the chain.
			expectRuleNum = tok == "-D" || tok == "--delete" ||
				tok == "-I" || tok == "--insert" ||
				tok == "-R" || tok == "--replace"
			i++
			// Consume the chain name if present (next non-flag token).
			if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
				result = append(result, args[i]) // chain name verbatim
				i++
			}
			continue
		}

		if verbatimFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			expectRuleNum = false
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			expectRuleNum = false
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			expectRuleNum = false
			i++
			continue
		}

		// Positional: if we expect a rule number after -D/-I chain, collapse to N.
		if expectRuleNum && shellshape.NumberRE.MatchString(tok) {
			result = append(result, "N")
			expectRuleNum = false
			i++
			continue
		}

		expectRuleNum = false
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
