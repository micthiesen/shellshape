package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestGdbus(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"call basic", "gdbus call --session --dest org.freedesktop.DBus --object-path /org/freedesktop/DBus --method org.freedesktop.DBus.ListNames", "gdbus call --session --dest <val> --object-path <path> --method <dotted-id>"},
		{"introspect", "gdbus introspect --system --dest org.bluez --object-path /", "gdbus introspect --system --dest <val> --object-path <path>"},
		{"monitor", "gdbus monitor --session --dest org.freedesktop.DBus", "gdbus monitor --session --dest <val>"},
		{"wait", "gdbus wait --session org.freedesktop.portal.Desktop", "gdbus wait --session <val>"},
		{"emit", "gdbus emit --object-path /com/example --signal com.example.Signal", "gdbus emit --object-path <path> --signal <dotted-id>"},

		// Method call with positional args (method arguments)
		{"call with args", "gdbus call --session --dest org.freedesktop.DBus --object-path /org/freedesktop/DBus --method org.freedesktop.DBus.GetNameOwner com.example.App", "gdbus call --session --dest <val> --object-path <path> --method <dotted-id> <val>"},

		// Timeout flag
		{"timeout", "gdbus call --session --dest org.freedesktop.DBus --object-path /org/freedesktop/DBus --method org.freedesktop.DBus.ListNames --timeout 5", "gdbus call --session --dest <val> --object-path <path> --method <dotted-id> --timeout N"},

		// Multiple positional args
		{"multiple args", "gdbus call --session --dest org.example --object-path /org/example --method org.example.Set key1 value1", "gdbus call --session --dest <val> --object-path <path> --method <dotted-id> <val>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test
	t.Run("different dest values collide", func(t *testing.T) {
		a := shellshape.Normalize("gdbus call --session --dest org.freedesktop.DBus --object-path /org/freedesktop/DBus --method org.freedesktop.DBus.ListNames")
		b := shellshape.Normalize("gdbus call --session --dest org.bluez --object-path /org/bluez --method org.bluez.Adapter1.StartDiscovery")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("gdbus call --session --dest org.freedesktop.DBus --object-path /org/freedesktop/DBus --method org.freedesktop.DBus.ListNames")
		subshell := shellshape.Normalize("gdbus call --session --dest $(evil-cmd) --object-path /org/freedesktop/DBus --method org.freedesktop.DBus.ListNames")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
