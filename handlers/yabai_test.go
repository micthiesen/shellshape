package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestYabai(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Boolean management flags stay verbatim.
		{"restart service", "yabai --restart-service", "yabai --restart-service"},
		{"start service", "yabai --start-service", "yabai --start-service"},
		{"stop service", "yabai --stop-service", "yabai --stop-service"},
		{"load sa", "yabai --load-sa", "yabai --load-sa"},

		// -m <domain> <subcommand> queries.
		{"query windows", "yabai -m query --windows", "yabai -m query --windows"},
		{"query spaces", "yabai -m query --spaces", "yabai -m query --spaces"},
		{"query displays", "yabai -m query --displays", "yabai -m query --displays"},

		// Window management with selectors.
		{"window focus next", "yabai -m window --focus next", "yabai -m window --focus <val>"},
		{"window focus prev", "yabai -m window --focus prev", "yabai -m window --focus <val>"},
		{"window swap recent", "yabai -m window --swap recent", "yabai -m window --swap <val>"},
		{"window focus by id", "yabai -m window --focus 12345", "yabai -m window --focus N"},

		// Space management.
		{"space rotate", "yabai -m space --rotate 270", "yabai -m space --rotate N"},
		{"space focus", "yabai -m space --focus 2", "yabai -m space --focus N"},

		// Config get/set.
		{"config set", "yabai -m config window_border on", "yabai -m config window_border <val>"},
		{"config set numeric", "yabai -m config top_padding 10", "yabai -m config top_padding N"},
		{"config get", "yabai -m config window_border", "yabai -m config window_border"},

		// Display selectors.
		{"display focus", "yabai -m display --focus next", "yabai -m display --focus <val>"},

		// Pipelines and redirects.
		{"with redirect", "yabai --restart-service 2>&1", "yabai --restart-service 2>&1"},
		{"piped query", "yabai -m query --displays | jq .", "yabai -m query --displays | jq <filter>"},

		// Compound from the approvals review.
		{
			"compound restart and query",
			"yabai --restart-service 2>&1 && sleep 3 && yabai -m query --displays | jq '.[].frame' && echo done",
			"yabai --restart-service 2>&1 && sleep N && yabai -m query --displays | jq <filter> && echo <str>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different selectors and values should produce the same shape.
	t.Run("different selectors collide", func(t *testing.T) {
		a := shellshape.Normalize("yabai -m window --focus next")
		b := shellshape.Normalize("yabai -m window --focus prev")
		c := shellshape.Normalize("yabai -m window --focus recent")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("different config values collide", func(t *testing.T) {
		a := shellshape.Normalize("yabai -m config window_border on")
		b := shellshape.Normalize("yabai -m config window_border off")
		if a != b {
			t.Errorf("expected same: %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse (mandatory).
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("yabai -m window --focus next")
		subshell := shellshape.Normalize("yabai -m window --focus $(rm -rf /tmp/foo)")
		if benign == subshell {
			t.Error("yabai with subshell must not collapse to same shape as literal")
		}
		if !contains(subshell, "$(") {
			t.Errorf("expected subshell marker in %q", subshell)
		}
	})
}
