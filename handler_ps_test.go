package shellshape

import "testing"

func TestPs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "ps", "ps"},
		{"bsd aux", "ps aux", "ps aux"},
		{"bsd axjf", "ps axjf", "ps axjf"},
		{"dash ef", "ps -ef", "ps -ef"},
		{"long all", "ps -A", "ps -A"},

		// Flags with arguments
		{"output format -o", "ps -eo pid,comm,%cpu", "ps -eo <fmt>"},
		{"output format -O", "ps -O %mem", "ps -O <fmt>"},
		{"user -u", "ps -u root", "ps -u <user>"},
		{"user -U", "ps -U michael", "ps -U <user>"},
		{"pid -p", "ps -p 1234", "ps -p <pid>"},
		{"pid list -p", "ps -p 1234,5678,9999", "ps -p <pid>"},
		{"group -G", "ps -G 100", "ps -G <gid>"},
		{"process group -g", "ps -g 42", "ps -g <grp>"},
		{"terminal -t", "ps -t tty1", "ps -t <tty>"},

		// Consuming flag at end of args (no next token)
		{"trailing consuming flag", "ps -o", "ps -o"},
		// Subshell as standalone token
		{"subshell standalone", "ps $(echo aux)", "ps $(echo <str>)"},
		// Subshell as next token after consuming flag
		{"subshell after consuming flag", "ps -u $(whoami)", "ps -u $(whoami)"},

		// Combinations
		{"format and pid", "ps -o pid,args -p 42", "ps -o <fmt> -p <pid>"},
		{"user and full", "ps -fu root", "ps -fu <user>"},
		{"ef with redirect", "ps -ef > procs.txt", "ps -ef > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different users collide", func(t *testing.T) {
		a := Normalize("ps -u root")
		b := Normalize("ps -u nobody")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different pids collide", func(t *testing.T) {
		a := Normalize("ps -p 1234")
		b := Normalize("ps -p 9876")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("ps -p 1234")
		subshell := Normalize("ps -p $(pgrep nginx)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
