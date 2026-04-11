package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestAt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare at with simple time", "at noon", "at <time>"},
		{"at with multi-word time", "at now + 5 minutes", "at <time>"},
		{"at with clock time", "at 9:30 PM Tue", "at <time>"},
		{"at with midnight", "at midnight", "at <time>"},
		{"at with teatime", "at teatime", "at <time>"},

		// Flags with arguments
		{"file flag short", "at -f script.sh noon", "at -f <path> <time>"},
		{"file flag with path", "at -f /home/user/backup.sh 3:00 AM", "at -f <path> <time>"},
		{"queue flag", "at -q a now + 1 hour", "at -q <queue> <time>"},
		{"time flag", "at -t 202301011200", "at -t <time>"},
		{"cat job flag", "at -c 5", "at -c <job-id>"},

		// Boolean flags
		{"mail flag", "at -m now + 2 hours", "at -m <time>"},
		{"no-mail flag", "at -M noon", "at -M <time>"},
		{"verbose flag", "at -v teatime", "at -v <time>"},

		// Related commands
		{"atq bare", "atq", "atq"},
		{"atq with queue", "atq -q b", "atq -q <queue>"},
		{"atrm single job", "atrm 5", "atrm <job-id>"},
		{"atrm multiple jobs", "atrm 1 2 3", "atrm <job-id>"},
		{"batch bare", "batch", "batch"},

		// at -l and at -d (aliases for atq/atrm)
		{"at list", "at -l", "at -l"},
		{"at delete", "at -d 42", "at -d <job-id>"},

		// Redirect
		{"with redirect", "at noon < /tmp/commands.txt", "at <time> < <path>"},

		// Edge cases
		{"at with relative time days", "at now + 3 days", "at <time>"},
		{"at with date", "at 10:00 AM Jul 31", "at <time>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different times collide", func(t *testing.T) {
		a := shellshape.Normalize("at now + 5 minutes")
		b := shellshape.Normalize("at 9:30 PM Tue")
		c := shellshape.Normalize("at midnight")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("different job ids collide", func(t *testing.T) {
		a := shellshape.Normalize("atrm 1")
		b := shellshape.Normalize("atrm 99")
		if a != b {
			t.Errorf("expected same: %q, %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("at noon")
		subshell := shellshape.Normalize("at $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
