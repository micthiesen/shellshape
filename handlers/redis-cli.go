package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("redis-cli", handleRedisCli)
}

// handleRedisCli handles the redis-cli command.
// Connection flags (-h, -s, -u, -a, --pass, --user, etc.) consume their argument.
// Numeric flags (-p, -n, -r, -i, --count, etc.) collapse next arg to N.
// Path flags (--eval, --rdb, --functions-rdb) collapse next arg to <path>.
// Pattern flags (--pattern, --quoted-pattern) collapse next arg to <pattern>.
// The first positional is the Redis command (<redis-cmd>).
// Remaining positionals become <arg>.
func handleRedisCli(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-p": true, "-n": true, "-r": true, "-i": true,
		"--port": true, "--db": true, "--repeat": true, "--interval": true,
		"--count": true, "--pipe-timeout": true,
		"--intrinsic-latency": true, "--memkeys-samples": true,
		"--lru-test": true,
	}

	hostFlags := map[string]bool{
		"-h": true, "--host": true,
	}

	pathFlags := map[string]bool{
		"-s": true, "--eval": true, "--rdb": true, "--functions-rdb": true,
	}

	strFlags := map[string]bool{
		"-a": true, "--pass": true, "--user": true,
		"-d": true, "-D": true, "-X": true, "--output": true,
	}

	uriFlags := map[string]bool{
		"-u": true,
	}

	patternFlags := map[string]bool{
		"--pattern": true, "--quoted-pattern": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: hostFlags, Placeholder: "<host>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: strFlags, Placeholder: "<str>"},
		{Flags: uriFlags, Placeholder: "<uri>"},
		{Flags: patternFlags, Placeholder: "<pattern>"},
	}

	var result []string
	redisCmdAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if !redisCmdAssigned {
				redisCmdAssigned = true
			}
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

		// Positional
		if !redisCmdAssigned {
			result = append(result, "<redis-cmd>")
			redisCmdAssigned = true
		} else {
			result = append(result, "<arg>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
