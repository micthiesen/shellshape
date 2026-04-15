package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("networkctl", handleNetworkctl, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleNetworkctl handles networkctl (systemd-networkd management).
// Subcommands: list, status, up, down, lldp, label, renew, forcerenew,
// reconfigure, reload, delete.
// Interface names (positionals) -> <val>.
// Boolean flags: --no-pager, --no-legend, -a/--all, -s/--stats.
// Fused flags like --json=short are kept with value verbatim (small finite set).
func handleNetworkctl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	i := 0

	result, _, i = shellshape.RepairLeadingFlagSubcommand(
		subcommand, args, i, result, nil, nil,
	)

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Flags: keep verbatim
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// All positionals: interface names -> <val>
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
