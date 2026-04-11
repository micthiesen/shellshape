package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	for _, name := range []string{"useradd", "usermod", "userdel"} {
		shellshape.Register(name, handleUseradd)
	}
}

// handleUseradd handles useradd, usermod, and userdel.
// The username (last positional) is preserved verbatim because it
// identifies the target user. Flags that take data arguments are
// collapsed to appropriate placeholders.
func handleUseradd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-d": true, "--home": true,
		"-b": true, "--base-dir": true,
		"-k": true, "--skel": true,
		"-R": true, "--root": true,
		"-s": true, "--shell": true,
	}

	groupFlag := map[string]bool{
		"-g": true, "--gid": true,
	}

	groupsFlag := map[string]bool{
		"-G": true, "--groups": true,
	}

	numericFlags := map[string]bool{
		"-u": true, "--uid": true,
		"-f": true, "--inactive": true,
	}

	commentFlags := map[string]bool{
		"-c": true, "--comment": true,
	}

	valueFlags := map[string]bool{
		"-e": true, "--expiredate": true,
		"-p": true, "--password": true,
		"-l": true, "--login": true,
		"-K": true, "--key": true,
		"-Z": true, "--selinux-user": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: groupFlag, Placeholder: "<group>"},
		{Flags: groupsFlag, Placeholder: "<groups>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: commentFlags, Placeholder: "<comment>"},
		{Flags: valueFlags, Placeholder: "<value>"},
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

		// Handle fused short flags like -aG where -a is boolean
		// and -G takes an argument. Only check unambiguous flags
		// (groups) -- numeric flags like -f are overloaded
		// (--inactive vs --force) and cannot be safely detected.
		if shellshape.IsFlagToken(tok) && !strings.HasPrefix(tok, "--") && len(tok) > 2 {
			lastChar := tok[len(tok)-1:]
			if groupsFlag["-"+lastChar] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<groups>")
				continue
			}
			if groupFlag["-"+lastChar] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<group>")
				continue
			}
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: username -- keep verbatim.
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
