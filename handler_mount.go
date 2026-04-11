package shellshape

func init() {
	for _, name := range []string{"mount", "umount"} {
		Register(name, handleMount)
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
	args, redirects := splitRedirects(tokens)

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

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if typeFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<type>")
				i++
			}
			continue
		}

		if optsFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<opts>")
				i++
			}
			continue
		}

		if labelFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<label>")
				i++
			}
			continue
		}

		if uuidFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<uuid>")
				i++
			}
			continue
		}

		if namespaceFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device or mount point path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
