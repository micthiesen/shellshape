package handlers

import (
	"strings"

	shellshape "github.com/micthiesen/shellshape"
)

func init() {
	for _, name := range []string{".confirm-sudo.sh", "confirm-sudo.sh"} {
		shellshape.Register(name, handleConfirmSudo)
	}
}

// handleConfirmSudo handles confirm-sudo.sh, a privilege-escalation wrapper
// that passes all arguments through to the inner command (no flags of its own).
// The remaining tokens are re-normalized as a nested command so that inner
// handlers (pacman, systemctl, etc.) apply their specialized rules.
func handleConfirmSudo(subcommand string, tokens []string) []string {
	if len(tokens) == 0 {
		return nil
	}

	inner := strings.Join(tokens, " ")
	normalized := shellshape.Normalize(inner)
	if normalized == "" {
		return nil
	}

	return strings.Fields(normalized)
}
