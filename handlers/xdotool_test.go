package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestXdotool(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// search subcommand
		{"search name", `xdotool search --name "LM Studio"`, "xdotool search --name <val>"},
		{"search class", `xdotool search --class "firefox"`, "xdotool search --class <val>"},
		{"search classname", `xdotool search --classname "Navigator"`, "xdotool search --classname <val>"},
		{"search onlyvisible name", `xdotool search --onlyvisible --name "Chrome"`, "xdotool search --onlyvisible --name <val>"},
		{"search pattern positional", `xdotool search "LM Studio"`, "xdotool search <val>"},
		{"search pid", `xdotool search --pid 1234`, "xdotool search --pid <val>"},

		// key subcommand
		{"key literal", "xdotool key ctrl+a", "xdotool key ctrl+a"},
		{"key enter", "xdotool key KP_Enter", "xdotool key KP_Enter"},
		{"key repeat delay", "xdotool key --repeat 3600 --delay 1000 b", "xdotool key --repeat N --delay N b"},

		// type subcommand
		{"type text", `xdotool type "Hello world"`, "xdotool type <str>"},
		{"type delay", `xdotool type --delay 500 "Hello world"`, "xdotool type --delay N <str>"},

		// window subcommands
		{"windowfocus id", "xdotool windowfocus 12345", "xdotool windowfocus N"},
		{"windowfocus sync", "xdotool windowfocus --sync 12345", "xdotool windowfocus --sync N"},
		{"windowactivate id", "xdotool windowactivate 12345", "xdotool windowactivate N"},
		{"getactivewindow", "xdotool getactivewindow", "xdotool getactivewindow"},

		// click/mouse
		{"click button", "xdotool click 3", "xdotool click N"},
		{"click repeat", "xdotool click --repeat 5 --delay 100 1", "xdotool click --repeat N --delay N N"},
		{"mousemove", "xdotool mousemove 100 200", "xdotool mousemove N N"},

		// windowsize/windowmove
		{"windowsize", "xdotool windowsize 12345 800 600", "xdotool windowsize N N N"},
		{"windowmove", "xdotool windowmove 12345 100 200", "xdotool windowmove N N N"},

		// aliases
		{"kdotool search", `kdotool search "LM Studio"`, "kdotool search <val>"},
		{"ydotool type", `ydotool type "Hello world"`, "ydotool type <str>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different search names collide", func(t *testing.T) {
		a := shellshape.Normalize(`xdotool search --name "LM Studio"`)
		b := shellshape.Normalize(`xdotool search --name "Firefox"`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different typed text collide", func(t *testing.T) {
		a := shellshape.Normalize(`xdotool type "Hello world"`)
		b := shellshape.Normalize(`xdotool type "Goodbye world"`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different window IDs collide", func(t *testing.T) {
		a := shellshape.Normalize("xdotool windowfocus 12345")
		b := shellshape.Normalize("xdotool windowfocus 67890")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize(`xdotool search --name "Firefox"`)
		subshell := shellshape.Normalize("xdotool search --name $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
