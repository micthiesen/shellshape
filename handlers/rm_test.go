package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestRm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare rm", "rm", "rm"},
		{"single file", "rm file.txt", "rm <path>"},
		{"absolute path", "rm /tmp/build", "rm <path>"},
		{"relative path", "rm ./old-output", "rm <path>"},

		// Flags (all boolean)
		{"force", "rm -f file.txt", "rm -f <path>"},
		{"recursive force", "rm -rf /tmp/build", "rm -rf <path>"},
		{"interactive", "rm -i file.txt", "rm -i <path>"},
		{"verbose", "rm -v file.txt", "rm -v <path>"},
		{"directory flag", "rm -d empty-dir/", "rm -d <path>"},
		{"recursive verbose force", "rm -rvf /opt/app", "rm -rvf <path>"},

		// Multiple positionals
		{"two files", "rm file1.txt file2.txt", "rm <path>+"},
		{"three dirs with rf", "rm -rf dir1/ dir2/ dir3/", "rm -rf <path>+"},

		// Double dash
		{"double dash", "rm -- -weird-file.txt", "rm -- <path>"},

		// Glob (classified as path due to extension)
		{"glob pattern", "rm *.log", "rm <path>"},

		// Tilde path
		{"tilde path", "rm ~/Downloads/old.zip", "rm <path>"},

		// Redirects
		{"with redirect", "rm file.txt 2>/dev/null", "rm <path> 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different file paths -> same shape
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("rm -rf /tmp/build")
		b := shellshape.Normalize("rm -rf /var/log/old")
		c := shellshape.Normalize("rm -rf ~/Downloads/junk")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("rm literal-arg")
		subshell := shellshape.Normalize("rm $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
