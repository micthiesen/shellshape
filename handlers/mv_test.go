package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMv(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare mv", "mv", "mv"},
		{"rename file", "mv old.txt new.txt", "mv <path>+"},
		{"move to dir", "mv file.txt /tmp/", "mv <path>+"},
		{"absolute paths", "mv /home/user/a.txt /home/user/b.txt", "mv <path>+"},

		// Boolean flags
		{"force", "mv -f foo.txt bar.txt", "mv -f <path>+"},
		{"interactive", "mv -i foo.txt bar.txt", "mv -i <path>+"},
		{"no-clobber", "mv -n foo.txt bar.txt", "mv -n <path>+"},
		{"verbose", "mv -v foo.txt bar.txt", "mv -v <path>+"},
		{"combined flags", "mv -fv old.txt new.txt", "mv -fv <path>+"},

		// Multiple source files
		{"multiple sources", "mv a.txt b.txt c.txt dest/", "mv <path>+"},

		// GNU -t flag (target directory, consumes next arg)
		// Note: -t's <path> arg merges with positional <path> via collapseRepeatedPlaceholders
		{"target dir flag", "mv -t /dest a.txt b.txt", "mv -t <path>+"},
		{"target dir long", "mv --target-directory /dest a.txt", "mv --target-directory <path>+"},

		// GNU -S flag (suffix, consumes next arg)
		{"suffix flag", "mv -S .bak foo.txt bar.txt", "mv -S <str> <path>+"},
		{"suffix long", "mv --suffix .bak foo.txt bar.txt", "mv --suffix <str> <path>+"},

		// Redirects
		{"with redirect", "mv foo bar > /tmp/log", "mv <path>+ > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different file paths → same shape
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("mv /tmp/a.txt /tmp/b.txt")
		b := shellshape.Normalize("mv config.yaml backup.yaml")
		c := shellshape.Normalize("mv ~/doc.pdf ~/archive/doc.pdf")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("mv literal-arg dest/")
		subshell := shellshape.Normalize("mv $(dangerous-command) dest/")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
