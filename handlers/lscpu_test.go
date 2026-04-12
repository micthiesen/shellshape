package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLscpu(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bare", "lscpu", "lscpu"},
		{"json output", "lscpu --json", "lscpu --json"},
		{"extended", "lscpu -e", "lscpu -e"},
		{"extended long", "lscpu --extended", "lscpu --extended"},
		{"all extended", "lscpu -a -e", "lscpu -a -e"},
		{"offline extended", "lscpu --offline --extended", "lscpu --offline --extended"},
		{"extended with columns fused", "lscpu --extended=CPU,CORE", "lscpu --extended=<val>"},
		{"parse with columns fused", "lscpu -p=CPU,SOCKET", "lscpu -p=<val>"},
		{"parse long fused", "lscpu --parse=CPU", "lscpu --parse=<val>"},
		{"hex output", "lscpu -x", "lscpu -x"},
		{"with redirect", "lscpu > /tmp/cpu.txt", "lscpu > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different fused values collide", func(t *testing.T) {
		a := shellshape.Normalize("lscpu --extended=CPU,CORE")
		b := shellshape.Normalize("lscpu --extended=SOCKET,NODE")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("lscpu --json")
		subshell := shellshape.Normalize("lscpu $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
