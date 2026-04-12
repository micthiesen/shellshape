package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestZed(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "zed file.txt", "zed <path>"},
		{"directory", "zed ~/Code/project", "zed <path>"},
		{"absolute path", "zed /etc/hosts", "zed <path>"},
		{"multiple files", "zed file1.go file2.go", "zed <path>+"},
		{"path with line col", "zed src/main.rs:42:10", "zed <path>"},

		// Boolean flags
		{"new window", "zed --new file.txt", "zed --new <path>"},
		{"add to workspace", "zed --add file.txt", "zed --add <path>"},
		{"wait short", "zed -w file.txt", "zed -w <path>"},
		{"wait long", "zed --wait file.txt", "zed --wait <path>"},
		{"foreground", "zed --foreground file.txt", "zed --foreground <path>"},

		// Combinations
		{"new and wait", "zed --new -w src/main.go", "zed --new -w <path>"},
		{"current dir", "zed .", "zed ."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("zed /home/user/project/main.go")
		b := shellshape.Normalize("zed /tmp/scratch.py")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("zed file.txt")
		subshell := shellshape.Normalize("zed $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
