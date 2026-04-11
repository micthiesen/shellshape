package shellshape

func init() {
	Register("ab", handleAb)
	Register("wrk", handleWrk)
	Register("hey", handleHey)
}

// handleAb handles the Apache Bench (ab) command.
// URL positionals are classified via classifyToken (producing <http-uri>, etc.).
// Numeric flags (-n, -c, -s, -t, -v, -b) collapse the next arg to N.
// Header flags (-H) collapse the next arg to <header>.
// Path flags (-e, -E, -g, -p, -u) collapse the next arg to <path>.
// Data flags (-A, -C, -P) collapse the next arg to <data>.
// Other argument-consuming flags (-T, -m, -X, -f, -B, -x, -y, -z, -Z) collapse the next arg to <val>.
func handleAb(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	categories := []flagCategory{
		{map[string]bool{
			"-n": true, "-c": true, "-s": true, "-t": true, "-v": true, "-b": true,
		}, "N"},
		{map[string]bool{"-H": true}, "<header>"},
		{map[string]bool{
			"-e": true, "-E": true, "-g": true, "-p": true, "-u": true,
		}, "<path>"},
		{map[string]bool{"-A": true, "-C": true, "-P": true}, "<data>"},
		{map[string]bool{
			"-T": true, "-m": true, "-X": true, "-f": true, "-B": true,
			"-x": true, "-y": true, "-z": true, "-Z": true,
		}, "<val>"},
	}

	return httpBenchLoop(args, categories, redirects)
}

// handleWrk handles the wrk HTTP benchmarking tool.
// Numeric flags (-c/--connections, -t/--threads) collapse the next arg to N.
// Header flags (-H/--header) collapse the next arg to <header>.
// Path flags (-s/--script) collapse the next arg to <path>.
// Duration/other flags (-d/--duration, --timeout) collapse the next arg to <val>.
func handleWrk(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	categories := []flagCategory{
		{map[string]bool{
			"-c": true, "--connections": true,
			"-t": true, "--threads": true,
		}, "N"},
		{map[string]bool{"-H": true, "--header": true}, "<header>"},
		{map[string]bool{"-s": true, "--script": true}, "<path>"},
		{map[string]bool{
			"-d": true, "--duration": true, "--timeout": true,
		}, "<val>"},
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
	args, redirects := splitRedirects(tokens)

	categories := []flagCategory{
		{map[string]bool{
			"-n": true, "-c": true, "-q": true, "-t": true, "-cpus": true,
		}, "N"},
		{map[string]bool{"-H": true}, "<header>"},
		{map[string]bool{"-D": true}, "<path>"},
		{map[string]bool{"-d": true, "-a": true}, "<data>"},
		{map[string]bool{
			"-m": true, "-T": true, "-A": true, "-o": true,
			"-x": true, "-z": true, "-host": true,
		}, "<val>"},
	}

	return httpBenchLoop(args, categories, redirects)
}

// httpBenchLoop is the shared token-processing loop for HTTP benchmarking tools.
func httpBenchLoop(args []string, categories []flagCategory, redirects []string) []string {
	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
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

		// Positional: classify (URLs become <http-uri>, etc.)
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
