package shellshape

func init() {
	for _, name := range []string{"groupadd", "groupmod", "groupdel"} {
		Register(name, handleGroupadd)
	}
}

// handleGroupadd handles groupadd, groupmod, and groupdel commands.
// All three follow the same pattern: flags followed by a group name positional.
// The group name is collapsed to <group>. Flags that consume arguments get
// appropriate placeholders (N for GIDs, <path> for directories, <str> for
// passwords/users/keys, <group> for new-name).
func handleGroupadd(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a GID (numeric).
	gidFlags := map[string]bool{
		"-g": true, "--gid": true,
	}
	// Flags whose next token is a directory path.
	pathFlags := map[string]bool{
		"-R": true, "--root": true,
		"-P": true, "--prefix": true,
	}
	// Flags whose next token is a string value (password, key=value, user list).
	strFlags := map[string]bool{
		"-p": true, "--password": true,
		"-K": true, "--key": true,
		"-U": true, "--users": true,
	}
	// Flags whose next token is a group name.
	groupFlags := map[string]bool{
		"-n": true, "--new-name": true,
	}

	categories := []flagCategory{
		{gidFlags, "N"},
		{pathFlags, "<path>"},
		{strFlags, "<str>"},
		{groupFlags, "<group>"},
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

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: group name
		result = append(result, "<group>")
		i++
	}

	result = append(result, redirects...)
	return result
}
