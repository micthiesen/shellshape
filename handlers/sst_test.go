package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSst(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands (subcommand already extracted by normalizer)
		{"deploy with stage", "sst deploy --stage prod", "sst deploy --stage prod"},
		{"dev with stage", "sst dev --stage dev", "sst dev --stage dev"},
		{"remove with stage", "sst remove --stage michael", "sst remove --stage michael"},
		{"build bare", "sst build", "sst build"},

		// Profile and region flags collapse
		{"deploy with profile", "sst deploy --stage prod --profile my-aws", "sst deploy --stage prod --profile <val>"},
		{"deploy with region", "sst deploy --stage prod --region us-east-1", "sst deploy --stage prod --region <val>"},

		// Verbose flag preserved
		{"dev verbose", "sst dev --stage dev --verbose", "sst dev --stage dev --verbose"},

		// Secret subcommand: name stays verbatim, value collapses
		{"secret set", "sst secret set DatabaseUrl my-secret-value", "sst secret set DatabaseUrl <str>"},
		{"secret set url value", "sst secret set DatabaseUrl postgres://user:pass@host/db", "sst secret set DatabaseUrl <str>"},
		{"secret list", "sst secret list --stage prod", "sst secret list --stage prod"},
		{"secret remove", "sst secret remove DatabaseUrl", "sst secret remove DatabaseUrl"},

		// Shell and tunnel subcommands
		{"shell bare", "sst shell", "sst shell"},
		{"tunnel bare", "sst tunnel", "sst tunnel"},

		// Positionals use classifyToken
		{"deploy with path config", "sst deploy --stage prod ./custom/path", "sst deploy --stage prod <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different stages do not collide", func(t *testing.T) {
		a := shellshape.Normalize("sst deploy --stage prod")
		b := shellshape.Normalize("sst deploy --stage dev")
		if a == b {
			t.Errorf("stages should be preserved: %q vs %q", a, b)
		}
	})

	t.Run("different profiles collide", func(t *testing.T) {
		a := shellshape.Normalize("sst deploy --stage prod --profile default")
		b := shellshape.Normalize("sst deploy --stage prod --profile company-sso")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different secret values collide", func(t *testing.T) {
		a := shellshape.Normalize("sst secret set MyKey value-one")
		b := shellshape.Normalize("sst secret set MyKey value-two")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("sst deploy --stage prod")
		subshell := shellshape.Normalize("sst deploy $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
