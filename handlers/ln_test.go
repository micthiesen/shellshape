package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLn(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare ln", "ln", "ln"},
		{"hard link two paths", "ln /usr/local/bin/fooprog-1.0 /usr/local/bin/fooprog", "ln <path>+"},
		{"symlink", "ln -s /usr/src /home/src", "ln -s <path>+"},
		{"force symlink", "ln -sf /usr/bin/python3 /usr/local/bin/python", "ln -sf <path>+"},

		// Flags (all boolean)
		{"verbose", "ln -v source.txt link.txt", "ln -v <path>+"},
		{"interactive", "ln -i old.txt new.txt", "ln -i <path>+"},
		{"combined flags", "ln -sfv /opt/app/bin/run /usr/local/bin/run", "ln -sfv <path>+"},
		{"no-deref", "ln -shf /new/target /existing/link", "ln -shf <path>+"},

		// Multiple source files → target dir (collapsed)
		{"multiple sources", "ln -s /a/foo /b/bar /c/baz /target/dir/", "ln -s <path>+"},

		// Relative paths
		{"relative source", "ln -s ../lib/libfoo.so .", "ln -s <path> ."},
		{"dot-slash target", "ln -s ../conf/app.yaml ./config.yaml", "ln -s <path>+"},

		// Single file
		{"single source", "ln -s /usr/src", "ln -s <path>"},

		// Redirects
		{"with redirect", "ln -s /a /b 2>/dev/null", "ln -s <path>+ 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different paths → same shape
	t.Run("different targets collide", func(t *testing.T) {
		a := shellshape.Normalize("ln -s /usr/src /home/src")
		b := shellshape.Normalize("ln -s /opt/app /var/link")
		c := shellshape.Normalize("ln -s ~/dotfiles/.bashrc ~/.bashrc")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ln literal-arg")
		subshell := shellshape.Normalize("ln $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
