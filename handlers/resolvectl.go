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

	// Handle global flag consumed as subcommand by the normalizer.
	if shellshape.IsFlagToken(subcommand) {
		// Check for --flag=value that was consumed as subcommand.
		if fused, ok := shellshape.ConsumeFusedFlag(subcommand, categories); ok {
			// The normalizer already emitted the subcommand token; we need to
			// replace it. But we can't - we can only append. The fused flag
			// is already in the output. We just need to handle the remaining.
			// Actually the normalizer emitted the raw subcommand token.
			// We can't change it. Let's just handle the remaining tokens.
			_ = fused
		} else if placeholder, ok := shellshape.MatchFlagCategory(subcommand, categories); ok {
			if i < len(args) && !shellshape.IsFlagToken(args[i]) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, placeholder)
				}
				i++
			}
		}
		// Consume more global flags before the real subcommand.
		for i < len(args) && shellshape.IsFlagToken(args[i]) {
			flag := args[i]
			if placeholder, ok := shellshape.MatchFlagCategory(flag, categories); ok {
				result, i = shellshape.ConsumeFlagArg(flag, args, i, result, placeholder)
			} else if f, ok := shellshape.ConsumeFusedFlag(flag, categories); ok {
				result = append(result, f)
				i++
			} else {
				result = append(result, flag)
				i++
			}
		}
		// Emit the real subcommand.
		if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
			result = append(result, args[i])
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
