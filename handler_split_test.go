package shellshape

import "testing"

func TestSplit(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"default split", "split file.txt", "split <path>"},
		{"no args (stdin)", "split", "split"},
		{"stdin dash", "split -", "split -"},

		// Flags with arguments
		{"line count", "split -l 1000 largefile.log", "split -l N <path>"},
		{"byte count", "split -b 10M bigfile.tar", "split -b <size> <path>"},
		{"chunk count", "split -n 5 data.csv", "split -n N <path>"},
		{"pattern split", "split -p 'Chapter' book.txt", "split -p <pattern> <path>"},
		{"suffix length", "split -a 3 file.txt", "split -a N <path>"},

		// Boolean flags
		{"numeric suffix", "split -d -l 1000 file.txt", "split -d -l N <path>"},
		{"continue flag", "split -c -l 500 file.txt", "split -c -l N <path>"},

		// With output prefix
		{"with prefix", "split -l 100 file.txt chunk_", "split -l N <path> <prefix>"},
		{"byte with prefix", "split -b 512k data.csv part_", "split -b <size> <path> <prefix>"},

		// GNU long forms
		{"long lines", "split --lines=1000 file.txt", "split --lines=<val> <path>"},
		{"long bytes", "split --bytes=10M file.txt", "split --bytes=<val> <path>"},
		{"long number", "split --number=5 file.txt", "split --number=<val> <path>"},
		{"long suffix-length", "split --suffix-length=3 file.txt", "split --suffix-length=<val> <path>"},

		// Combined flags
		{"numeric suffix with line count", "split -d -l 500 -a 4 input.log out_", "split -d -l N -a N <path> <prefix>"},

		// Redirect
		{"with redirect", "split -l 100 < input.txt", "split -l N < <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("split -l 1000 server.log")
		b := Normalize("split -l 1000 access.log")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different line counts collide", func(t *testing.T) {
		a := Normalize("split -l 100 file.txt")
		b := Normalize("split -l 5000 file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different prefixes collide", func(t *testing.T) {
		a := Normalize("split -l 100 file.txt chunk_")
		b := Normalize("split -l 100 file.txt part_")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("split literal-arg")
		subshell := Normalize("split $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestCsplit(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"split at lines", "csplit file.txt 5 23", "csplit <path> <split-arg>+"},
		{"split at pattern", "csplit file.txt '/^Chapter/'", "csplit <path> <split-arg>"},
		{"pattern with repeat", "csplit file.txt '/^Chapter/' '{*}'", "csplit <path> <split-arg>+"},

		// Flags
		{"keep files on error", "csplit -k file.txt 100 '{19}'", "csplit -k <path> <split-arg>+"},
		{"silent", "csplit -s file.txt '/^---/' '{*}'", "csplit -s <path> <split-arg>+"},
		{"custom prefix", "csplit -f section_ file.txt '/^##/' '{*}'", "csplit -f <prefix> <path> <split-arg>+"},
		{"digit count", "csplit -n 4 file.txt 100", "csplit -n N <path> <split-arg>"},

		// Combined flags
		{"keep and silent", "csplit -k -s file.txt '/PATTERN/' '{5}'", "csplit -k -s <path> <split-arg>+"},
		{"prefix and digits", "csplit -f out_ -n 3 data.txt 50 '{9}'", "csplit -f <prefix> -n N <path> <split-arg>+"},

		// Stdin
		{"stdin dash", "csplit - 100 '{19}'", "csplit - <split-arg>+"},

		// Redirect
		{"with redirect", "csplit file.txt 100 2>/dev/null", "csplit <path> <split-arg> 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("csplit server.log '/^Error/' '{*}'")
		b := Normalize("csplit access.log '/^Error/' '{*}'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different patterns collide", func(t *testing.T) {
		a := Normalize("csplit file.txt '/^Chapter/'")
		b := Normalize("csplit file.txt '/^Section/'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("csplit literal-arg 100")
		subshell := Normalize("csplit $(dangerous-command) 100")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
