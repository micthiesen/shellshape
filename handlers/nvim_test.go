package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNvim(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "nvim file.txt", "nvim <path>"},
		{"multiple files", "nvim file1.go file2.go", "nvim <path>+"},
		{"absolute path", "nvim /etc/hosts", "nvim <path>"},

		// Line number with +N
		{"plus line number", "nvim +42 file.txt", "nvim +N <path>"},
		{"plus line number large", "nvim +100 main.go", "nvim +N <path>"},

		// Plus command (+/pattern, +cmd)
		{"plus search pattern", "nvim +/TODO file.txt", "nvim +<val> <path>"},
		{"plus command", "nvim +PlugInstall", "nvim +<val>"},

		// Flags with arguments
		{"ex command", "nvim -c 'set number' file.txt", "nvim -c <val> <path>"},
		{"multiple ex commands", "nvim -c 'set nu' -c 'syntax on' file.txt", "nvim -c <val> -c <val> <path>"},
		{"custom vimrc", "nvim -u ~/.vimrc file.txt", "nvim -u <path>+"},
		{"session file", "nvim -S session.vim", "nvim -S <path>"},

		// Boolean flags
		{"diff mode", "nvim -d file1.txt file2.txt", "nvim -d <path>+"},
		{"readonly", "nvim -R file.txt", "nvim -R <path>"},
		{"tabs", "nvim -p file1.go file2.go", "nvim -p <path>+"},
		{"hsplit", "nvim -o file1.go file2.go", "nvim -o <path>+"},
		{"vsplit", "nvim -O file1.go file2.go", "nvim -O <path>+"},
		{"headless", "nvim --headless +PlugInstall +qa", "nvim --headless +<val> +<val>"},
		{"clean", "nvim --clean file.txt", "nvim --clean <path>"},

		// Aliases
		{"vim alias", "vim file.txt", "vim <path>"},
		{"vi alias", "vi file.txt", "vi <path>"},

		// Redirect
		{"with redirect", "nvim file.txt > /dev/null 2>&1", "nvim <path> > <path> 2>&1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("nvim /home/user/code.go")
		b := shellshape.Normalize("nvim /tmp/scratch.py")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different line numbers collide", func(t *testing.T) {
		a := shellshape.Normalize("nvim +10 file.go")
		b := shellshape.Normalize("nvim +999 other.py")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("nvim file.txt")
		subshell := shellshape.Normalize("nvim $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
