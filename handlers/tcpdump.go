package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("tcpdump", handleTcpdump)
}

// handleTcpdump handles the tcpdump command.
// Interface flag (-i) collapses the next arg to <val>.
// File flags (-w, -r, -F, -V) collapse the next arg to <path>.
// Numeric flags (-c, -s, -C, -W, -G, -B) collapse the next arg to N.
// Val flags (-E, -T, -z, -Z, -y, -Q, -j, --time-stamp-precision) collapse the next arg to <val>.
// All positionals are BPF filter tokens and collapse to <filter>.
func handleTcpdump(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-w": true, "-r": true, "-F": true, "-V": true,
	}

	numericFlags := map[string]bool{
		"-c": true, "-s": true, "--snapshot-length": true,
		"-C": true, "-W": true, "-G": true,
		"-B": true, "--buffer-size": true,
	}

	valFlags := map[string]bool{
		"-i": true, "--interface": true,
		"-E": true, "-T": true,
		"-z": true, "-Z": true, "--relinquish-privileges": true,
		"-y": true, "--linktype": true,
		"-Q": true, "--direction": true,
		"-j": true, "--time-stamp-type": true,
		"--time-stamp-precision": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: BPF filter token
		result = append(result, "<filter>")
		i++
	}

	result = append(result, redirects...)
	return result
}
