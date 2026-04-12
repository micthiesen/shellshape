package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("iperf3", handleIperf3)
}

// handleIperf3 handles iperf3 (network bandwidth testing).
// Boolean flags: -s/--server, -R/--reverse, -u/--udp, -D/--daemon, -J/--json,
//
//	-V/--verbose, --version, --help, -1/--one-off, -4/-6.
//
// Val flags: -c/--client (host), -b/--bandwidth, -f/--format, --bind.
// Numeric flags: -p/--port, -t/--time, -n/--bytes, -P/--parallel, -i/--interval,
//
//	-l/--length, -w/--window, -M/--set-mss, -k/--blockcount, --connect-timeout.
//
// Path flags: --logfile, --pidfile.
func handleIperf3(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-c": true, "--client": true,
		"-b": true, "--bandwidth": true,
		"-f": true, "--format": true,
		"--bind":     true,
		"-B":         true,
		"--cport":    true,
		"--affinity": true,
	}

	numericFlags := map[string]bool{
		"-p": true, "--port": true,
		"-t": true, "--time": true,
		"-n": true, "--bytes": true,
		"-P": true, "--parallel": true,
		"-i": true, "--interval": true,
		"-l": true, "--length": true,
		"-w": true, "--window": true,
		"-M": true, "--set-mss": true,
		"-k": true, "--blockcount": true,
		"--connect-timeout": true,
		"-O":                true, "--omit": true,
	}

	pathFlags := map[string]bool{
		"--logfile": true,
		"--pidfile": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: pathFlags, Placeholder: "<path>"},
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

		// Fused flags like --port=5201
		if f, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, f)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// All other flags: boolean, keep verbatim
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positionals (rare for iperf3, classify generically)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
