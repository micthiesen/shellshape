package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestRateMirrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"arch basic", "rate-mirrors arch", "rate-mirrors arch"},
		{"cachyos basic", "rate-mirrors cachyos", "rate-mirrors cachyos"},
		{"chaotic-aur", "rate-mirrors chaotic-aur", "rate-mirrors chaotic-aur"},
		{"endeavouros", "rate-mirrors endeavouros", "rate-mirrors endeavouros"},
		{"manjaro", "rate-mirrors manjaro", "rate-mirrors manjaro"},

		// Flags after subcommand
		{"country", "rate-mirrors arch --country US", "rate-mirrors arch --country <val>"},
		{"protocol structural", "rate-mirrors arch --protocol https", "rate-mirrors arch --protocol https"},
		{"max-delay", "rate-mirrors arch --max-delay 7200", "rate-mirrors arch --max-delay N"},
		{"save path", "rate-mirrors arch --save /etc/pacman.d/mirrorlist", "rate-mirrors arch --save <path>"},
		{"fetch-mirrors-timeout", "rate-mirrors cachyos --fetch-mirrors-timeout 30", "rate-mirrors cachyos --fetch-mirrors-timeout N"},
		{"allow-root", "rate-mirrors arch --allow-root", "rate-mirrors arch --allow-root"},
		{"concurrency", "rate-mirrors arch --concurrency 16", "rate-mirrors arch --concurrency N"},

		// Combined flags
		{"combined", "rate-mirrors arch --country DE --protocol https --max-delay 3600 --save /tmp/mirrors",
			"rate-mirrors arch --country <val> --protocol https --max-delay N --save <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test: different country codes produce same shape
	t.Run("different countries collide", func(t *testing.T) {
		a := shellshape.Normalize("rate-mirrors arch --country US")
		b := shellshape.Normalize("rate-mirrors arch --country DE")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("rate-mirrors arch --country US")
		subshell := shellshape.Normalize("rate-mirrors arch --country $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
