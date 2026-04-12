package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFastfetch(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no args", "fastfetch", "fastfetch"},
		{"logo", "fastfetch --logo arch", "fastfetch --logo arch"},
		{"config path", "fastfetch --config /home/user/.config/fastfetch/config.jsonc", "fastfetch --config <path>"},
		{"short config", "fastfetch -c /home/user/.config/fastfetch/config.jsonc", "fastfetch -c <path>"},
		{"load-config", "fastfetch --load-config /tmp/my-config.jsonc", "fastfetch --load-config <path>"},
		{"color", "fastfetch --color blue", "fastfetch --color <val>"},
		{"separator", "fastfetch --separator ' - '", "fastfetch --separator <val>"},
		{"structure", "fastfetch --structure Title:OS:Kernel:Uptime", "fastfetch --structure <val>"},
		{"logo and color", "fastfetch --logo ubuntu --color yellow", "fastfetch --logo ubuntu --color <val>"},
		{"pipe mode", "fastfetch --pipe", "fastfetch --pipe"},
		{"multiple flags", "fastfetch --logo arch --pipe --stat", "fastfetch --logo arch --pipe --stat"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different configs collide", func(t *testing.T) {
		a := shellshape.Normalize("fastfetch --config /home/alice/.config/fastfetch/config.jsonc")
		b := shellshape.Normalize("fastfetch --config /home/bob/.config/fastfetch/config.jsonc")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("fastfetch --logo arch")
		subshell := shellshape.Normalize("fastfetch $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
