package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("ab", handleAb)
	shellshape.Register("wrk", handleWrk)
	shellshape.Register("hey", handleHey)
}

// handleAb handles the Apache Bench (ab) command.
// URL positionals are classified via classifyToken (producing <http-uri>, etc.).
// Numeric flags (-n, -c, -s, -t, -v, -b) collapse the next arg to N.
// Header flags (-H) collapse the next arg to <header>.
// Path flags (-e, -E, -g, -p, -u) collapse the next arg to <path>.
// Data flags (-A, -C, -P) collapse the next arg to <data>.
// Other argument-consuming flags (-T, -m, -X, -f, -B, -x, -y, -z, -Z) collapse the next arg to <val>.
func handleAb(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"-n": true, "-c": true, "-s": true, "-t": true, "-v": true, "-b": true,
		}, Placeholder: "N"},
		{Flags: map[string]bool{"-H": true}, Placeholder: "<header>"},
		{Flags: map[string]bool{
			"-e": true, "-E": true, "-g": true, "-p": true, "-u": true,
		}, Placeholder: "<path>"},
		{Flags: map[string]bool{"-A": true, "-C": true, "-P": true}, Placeholder: "<data>"},
		{Flags: map[string]bool{
			"-T": true, "-m": true, "-X": true, "-f": true, "-B": true,
			"-x": true, "-y": true, "-z": true, "-Z": true,
		}, Placeholder: "<val>"},
	}

	return httpBenchLoop(args, categories, redirects)
}

// handleWrk handles the wrk HTTP benchmarking tool.
// Numeric flags (-c/--connections, -t/--threads) collapse the next arg to N.
// Header flags (-H/--header) collapse the next arg to <header>.
// Path flags (-s/--script) collapse the next arg to <path>.
// Duration/other flags (-d/--duration, --timeout) collapse the next arg to <val>.
func handleWrk(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"-c": true, "--connections": true,
			"-t": true, "--threads": true,
		}, Placeholder: "N"},
		{Flags: map[string]bool{"-H": true, "--header": true}, Placeholder: "<header>"},
		{Flags: map[string]bool{"-s": true, "--script": true}, Placeholder: "<path>"},
		{Flags: map[string]bool{
			"-d": true, "--duration": true, "--timeout": true,
		}, Placeholder: "<val>"},
	}

	return httpBenchLoop(args, categories, redirects)
}

// handleHey handles the hey HTTP load generator.
// Numeric flags (-n, -c, -q, -t, -cpus) collapse the next arg to N.
// Header flags (-H) collapse the next arg to <header>.
// Path flags (-D) collapse the next arg to <path>.
// Data flags (-d, -a) collapse the next arg to <data>.
// Other argument-consuming flags (-m, -T, -A, -o, -x, -z, -host) collapse the next arg to <val>.
func handleHey(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"-n": true, "-c": true, "-q": true, "-t": true, "-cpus": true,
		}, Placeholder: "N"},
		{Flags: map[string]bool{"-H": true}, Placeholder: "<header>"},
		{Flags: map[string]bool{"-D": true}, Placeholder: "<path>"},
		{Flags: map[string]bool{"-d": true, "-a": true}, Placeholder: "<data>"},
		{Flags: map[string]bool{
			"-m": true, "-T": true, "-A": true, "-o": true,
			"-x": true, "-z": true, "-host": true,
		}, Placeholder: "<val>"},
	}

	return httpBenchLoop(args, categories, redirects)
}

// httpBenchLoop is the shared token-processing loop for HTTP benchmarking tools.
func httpBenchLoop(args []string, categories []shellshape.FlagCategory, redirects []string) []string {
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

		// Positional: classify (URLs become <http-uri>, etc.)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
