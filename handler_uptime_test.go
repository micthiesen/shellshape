package shellshape

import "testing"

func TestUptime(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare uptime", "uptime", "uptime"},
		{"pretty flag short", "uptime -p", "uptime -p"},
		{"pretty flag long", "uptime --pretty", "uptime --pretty"},
		{"since flag short", "uptime -s", "uptime -s"},
		{"since flag long", "uptime --since", "uptime --since"},
		{"version flag", "uptime -V", "uptime -V"},
		{"help flag", "uptime -h", "uptime -h"},

		// With redirect
		{"with redirect", "uptime > /tmp/out.txt", "uptime > <path>"},

		// In pipeline
		{"in pipeline", "uptime && echo done", "uptime && echo <str>"},

		// Unexpected positional gets classified
		{"unexpected positional", "uptime somefile.txt", "uptime <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different redirects collide", func(t *testing.T) {
		a := Normalize("uptime > /tmp/a.txt")
		b := Normalize("uptime > /var/log/b.log")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("uptime somefile")
		subshell := Normalize("uptime $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
