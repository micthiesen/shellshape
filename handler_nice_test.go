package shellshape

import "testing"

func TestNice(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare nice", "nice", "nice"},
		{"nice with command", "nice make", "nice make"},
		{"nice with command and args", "nice gcc -o main main.c", "nice gcc -o main <path>"},
		{"nice date", "nice date", "nice date"},

		// Flags with arguments
		{"nice -n with command", "nice -n 10 make", "nice -n N make"},
		{"nice -n with command and args", "nice -n 5 gcc -o output main.c", "nice -n N gcc -o output <path>"},
		{"nice -n negative", "nice -n -20 date", "nice -n N date"},

		// Multiple positionals after utility
		{"command with adjacent paths", "nice tar -czf archive.tar.gz /home/user", "nice tar -czf <path>+"},
		{"command with flags and args", "nice find /tmp -name '*.log'", "nice find <path> -name <path>"},

		// Redirect
		{"with redirect", "nice sort file.txt > out.txt", "nice sort <path> > <path>"},

		// Pipeline
		{"in pipeline", "nice make && nice date", "nice make && nice date"},
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
	t.Run("different priorities collide", func(t *testing.T) {
		a := Normalize("nice -n 10 make")
		b := Normalize("nice -n 15 make")
		c := Normalize("nice -n 5 make")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("different command args collide", func(t *testing.T) {
		a := Normalize("nice gcc -o main foo.c")
		b := Normalize("nice gcc -o main bar.c")
		if a != b {
			t.Errorf("expected same: %q, %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("nice make target")
		subshell := Normalize("nice $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
