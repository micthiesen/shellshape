package shellshape

import (
	"regexp"
	"strings"
)

func init() {
	for _, name := range []string{"clang", "clang++", "gcc", "g++", "cc", "c++"} {
		Register(name, handleClang)
	}
}

// Fused short-flag prefixes: -I<path>, -L<path>, -D<def>, -U<def>, -l<lib>.
var clangFusedShortRE = regexp.MustCompile(`^-(I|L)(.+)$`)
var clangFusedDefRE = regexp.MustCompile(`^-(D|U)(.+)$`)
var clangFusedLibRE = regexp.MustCompile(`^-l(.+)$`)

// Fused long-flag=value patterns.
var clangFusedLongRE = regexp.MustCompile(`^(-std|-march|-mcpu|-mtune|-mabi|-mfpu|-mfloat-abi|--target)=(.+)$`)
var clangFusedLongPathRE = regexp.MustCompile(`^(--sysroot)=(.+)$`)

// handleClang handles clang, clang++, gcc, g++, cc, c++.
//
// Path flags (-o, -I, -L, -isystem, etc.) collapse the next arg to <path>.
// Define flags (-D, -U) collapse the next arg to <def>.
// Value flags (-l, -std, --target, -arch, -x, etc.) collapse the next arg to <val>.
// Fused forms like -I/path, -DFOO, -lm, -std=c++17, --target=... are split and collapsed.
// Optimization flags (-O0, -O2, -Os, etc.) and warning flags (-W...) are kept verbatim.
// Positionals use classifyToken.
func handleClang(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true, "--output": true,
		"-I": true, "-L": true,
		"-isystem": true, "-isysroot": true,
		"-iquote": true, "-idirafter": true,
		"-iprefix": true, "-iwithprefix": true, "-iwithprefixbefore": true,
		"-F": true, "-B": true,
		"-MF": true, "-MQ": true, "-MT": true,
		"-include": true, "-imacros": true,
	}

	defFlags := map[string]bool{
		"-D": true, "-U": true,
	}

	valFlags := map[string]bool{
		"-l": true, "-Xlinker": true,
		"-target": true, "--target": true,
		"-arch": true, "-march": true, "-mcpu": true, "-mtune": true,
		"-mabi": true, "-mfpu": true, "-mfloat-abi": true,
		"-std": true, "--std": true,
		"-x": true, "-T": true, "-e": true,
		"--serialize-diagnostics": true, "-MJ": true,
		"-rpath": true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{defFlags, "<def>"},
		{valFlags, "<val>"},
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

		// Fused short flags: -I/path, -L/path
		if m := clangFusedShortRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-"+m[1], "<path>")
			i++
			continue
		}

		// Fused define flags: -DFOO, -UFOO
		if m := clangFusedDefRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-"+m[1], "<def>")
			i++
			continue
		}

		// Fused library: -lfoo
		if m := clangFusedLibRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-l", "<val>")
			i++
			continue
		}

		// Fused long flag=value (val): -std=c++17, --target=..., -march=...
		if m := clangFusedLongRE.FindStringSubmatch(tok); m != nil {
			result = append(result, m[1]+"=<val>")
			i++
			continue
		}

		// Fused long flag=value (path): --sysroot=...
		if m := clangFusedLongPathRE.FindStringSubmatch(tok); m != nil {
			result = append(result, m[1]+"=<path>")
			i++
			continue
		}

		// Space-separated flag categories
		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Other flags (including -O*, -W*, -f*, -m (boolean), -g, -c, -S, -E, etc.)
		if isFlagToken(tok) {
			// Check for --flag=value on remaining long flags
			if strings.Contains(tok, "=") {
				eqIdx := strings.IndexByte(tok, '=')
				key := tok[:eqIdx]
				if pathFlags[key] {
					result = append(result, key+"=<path>")
				} else if defFlags[key] {
					result = append(result, key+"=<def>")
				} else if valFlags[key] {
					result = append(result, key+"=<val>")
				} else {
					result = append(result, tok)
				}
				i++
				continue
			}
			result = append(result, tok)
			i++
			continue
		}

		// Positional: source files, object files, etc.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
