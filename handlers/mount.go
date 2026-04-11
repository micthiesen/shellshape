package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"mount", "umount"} {
		shellshape.Register(name, handleMount)
	}
}

// handleMount handles the mount and umount commands.
// -t/--types consumes a filesystem type (<type>).
// -o/--options and -O/--test-opts consume mount options (<opts>).
// -L/--label consumes a label (<label>).
// -U/--uuid consumes a UUID (<uuid>).
// -N/--namespace consumes a numeric namespace (N).
// --source, --target, -T/--fstab consume a path (<path>).
// All other flags are boolean. Positionals are paths via classifyToken.
func handleMount(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	typeFlags := map[string]bool{
		"-t": true, "--types": true,
	}
	optsFlags := map[string]bool{
		"-o": true, "--options": true,
		"-O": true, "--test-opts": true,
	}
	labelFlags := map[string]bool{
		"-L": true, "--label": true,
	}
	uuidFlags := map[string]bool{
		"-U": true, "--uuid": true,
	}
	namespaceFlags := map[string]bool{
		"-N": true, "--namespace": true,
	}
	pathFlags := map[string]bool{
		"--source": true, "--target": true,
		"-T": true, "--fstab": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: typeFlags, Placeholder: "<type>"},
		{Flags: optsFlags, Placeholder: "<opts>"},
		{Flags: labelFlags, Placeholder: "<label>"},
		{Flags: uuidFlags, Placeholder: "<uuid>"},
		{Flags: namespaceFlags, Placeholder: "N"},
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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device or mount point path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
