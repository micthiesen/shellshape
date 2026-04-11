package shellshape

import "regexp"

func init() {
	for _, name := range []string{
		"mkfs",
		"mkfs.ext2", "mkfs.ext3", "mkfs.ext4",
		"mkfs.xfs", "mkfs.btrfs",
		"mkfs.vfat", "mkfs.fat",
		"mkfs.ntfs", "mkfs.exfat",
		"mkfs.f2fs", "mkfs.minix",
	} {
		Register(name, handleMkfs)
	}
}

var mkfsKeyValRE = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)=(.+)$`)

// handleMkfs handles mkfs and its filesystem-specific variants (mkfs.ext4, mkfs.xfs, etc.).
//
// -t/--type: next arg is filesystem type (preserved verbatim, structural).
// -L, -n: next arg is a label (collapsed to <label>).
// -b, -i, -I, -m, -N, -s, -F (vfat FAT size): next arg is numeric (collapsed to N).
//
//	For XFS-style key=value args (e.g. size=4096), the key is preserved and the value collapsed.
//
// -U: next arg is a UUID (collapsed via classifyToken which handles UUIDs).
// -O: next arg is a feature list (preserved verbatim, structural).
// All other flags are boolean. Positionals are device paths via classifyToken.
func handleMkfs(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	typeFlags := map[string]bool{"-t": true, "--type": true}
	labelFlags := map[string]bool{"-L": true, "-n": true}
	numericFlags := map[string]bool{
		"-b": true, "-i": true, "-I": true,
		"-m": true, "-N": true, "-s": true,
	}
	featureFlags := map[string]bool{"-O": true}
	uuidFlags := map[string]bool{"-U": true}

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
				result = append(result, args[i])
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

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, collapseKeyValNumeric(args[i]))
				i++
			}
			continue
		}

		if featureFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if uuidFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, classifyToken(args[i]))
				i++
			}
			continue
		}

		// -F is special: in ext* it's boolean (force), in vfat it takes a numeric arg (FAT size).
		// We can't reliably distinguish at parse time, so we check if the next
		// token looks numeric. If so, consume it; otherwise treat -F as boolean.
		if tok == "-F" {
			result = append(result, tok)
			i++
			if i < len(args) && numberRE.MatchString(args[i]) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// collapseKeyValNumeric handles XFS-style key=value args (e.g. "size=4096").
// If the value is numeric, it collapses to "key=N". Otherwise falls back to classifyToken.
func collapseKeyValNumeric(tok string) string {
	if m := mkfsKeyValRE.FindStringSubmatch(tok); m != nil {
		key, val := m[1], m[2]
		if numberRE.MatchString(val) {
			return key + "=N"
		}
		return tok
	}
	if numberRE.MatchString(tok) {
		return "N"
	}
	return classifyToken(tok)
}
