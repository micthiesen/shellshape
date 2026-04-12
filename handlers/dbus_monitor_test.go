package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDbusMonitor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "dbus-monitor", "dbus-monitor"},
		{"system bus", "dbus-monitor --system", "dbus-monitor --system"},
		{"session bus", "dbus-monitor --session", "dbus-monitor --session"},
		{"profile mode", "dbus-monitor --profile", "dbus-monitor --profile"},
		{"monitor mode", "dbus-monitor --monitor", "dbus-monitor --monitor"},

		// Address flag
		{"address", "dbus-monitor --address unix:path=/run/user/1000/bus", "dbus-monitor --address <val>"},

		// Filter expressions → <val>
		{"single filter", `dbus-monitor "type='signal'"`, "dbus-monitor <val>"},
		{"filter with interface", `dbus-monitor "type='signal',interface='org.kde.KWin'"`, "dbus-monitor <val>"},
		{"system with filter", `dbus-monitor --system "type='signal',sender='org.freedesktop.DBus'"`, "dbus-monitor --system <val>"},
		{"multiple filters", `dbus-monitor "type='signal'" "type='method_call'"`, "dbus-monitor <val>+"},

		// Combined
		{"session with filter", `dbus-monitor --session "type='signal',interface='org.freedesktop.Notifications'"`,
			"dbus-monitor --session <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test: different filter expressions produce same shape
	t.Run("different filters collide", func(t *testing.T) {
		a := shellshape.Normalize(`dbus-monitor --system "type='signal',interface='org.kde.KWin'"`)
		b := shellshape.Normalize(`dbus-monitor --system "type='method_call',sender='org.freedesktop.DBus'"`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize(`dbus-monitor "type='signal'"`)
		subshell := shellshape.Normalize("dbus-monitor $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
