package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestBase64(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"encode file", "base64 file.txt", "base64 <path>"},
		{"encode path", "base64 /tmp/data.bin", "base64 <path>"},
		{"decode flag short", "base64 -d /tmp/encoded", "base64 -d <path>"},
		{"decode flag long", "base64 --decode ./data.txt", "base64 --decode <path>"},
		{"macOS decode flag", "base64 -D ./encoded.txt", "base64 -D <path>"},

		// Flags with arguments
		{"wrap columns", "base64 -w 76 file.txt", "base64 -w N <path>"},
		{"wrap long", "base64 --wrap 0 file.txt", "base64 --wrap N <path>"},
		{"break macOS", "base64 -b 76 file.txt", "base64 -b N <path>"},
		{"break long", "base64 --break 64 file.txt", "base64 --break N <path>"},
		{"output file", "base64 -o /tmp/output.txt /tmp/input.txt", "base64 -o <path>+"},
		{"output long", "base64 --output /tmp/result.txt /tmp/input.txt", "base64 --output <path>+"},

		// Combined flags
		{"decode with output", "base64 -D -o /tmp/decoded.txt ./input.txt", "base64 -D -o <path>+"},
		{"decode wrap zero", "base64 -d -w 0", "base64 -d -w N"},
		{"ignore garbage", "base64 -d -i ./encoded.txt", "base64 -d -i <path>"},

		// No positionals (stdin)
		{"no args", "base64", "base64"},
		{"decode only", "base64 -d", "base64 -d"},

		// Redirect
		{"stdin redirect", "base64 < /tmp/input.txt", "base64 < <path>"},
		{"stdout redirect", "base64 file.txt > /tmp/out.txt", "base64 <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("base64 /tmp/secret.bin")
		b := shellshape.Normalize("base64 /home/user/photo.jpg")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different output paths collide", func(t *testing.T) {
		a := shellshape.Normalize("base64 -o /tmp/a.txt /tmp/input1.txt")
		b := shellshape.Normalize("base64 -o /home/user/b.txt /var/data.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("base64 literal-arg")
		subshell := shellshape.Normalize("base64 $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
