package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("resolvectl", handleResolvectl, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleResolvectl handles resolvectl (systemd-resolved management).
// Subcommands: query, status, statistics, flush-caches, dns, domain, llmnr,
// mdns, dnssec, nta, revert, openpgp, tlsa, service.
// Domain names, IPs, interface names -> <val>.
// Global flags like --type, -t, -p consume values; -i consumes a numeric interface index.
func handleResolvectl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Global flags that consume the next token as <val>.
	valFlags := map[string]bool{
		"-t": true, "--type": true,
		"-c": true, "--class": true,
		"-p": true, "--protocol": true,
		"--legend": true,
	}

	// Global flags that consume the next token as N (interface index).
	numericFlags := map[string]bool{
		"-i": true, "--interface": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: numericFlags, Placeholder: "N"},
	}

	var result []string
	i := 0

	result, _, i = shellshape.RepairLeadingFlagSubcommand(
		subcommand, args, i, result, categories, nil,
	)

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Fused flags like --type=MX
		if f, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, f)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Other flags: keep verbatim
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// All positionals: domains, IPs, interface names -> <val>
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
