package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDbusSend(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{
			"session with string args",
			`dbus-send --session --dest=org.kde.KWin --print-reply /Scripting org.kde.kwin.Scripting.loadScript string:"/tmp/close-lmstudio.js" string:"close-lmstudio"`,
			"dbus-send --session --dest=<val> --print-reply <path> <dotted-id> string:<val> string:<val>",
		},
		{
			"system bus signal",
			"dbus-send --system --type=signal / com.example.TestSignal",
			"dbus-send --system --type=<val> <path> <dotted-id>",
		},
		{
			"int32 argument",
			"dbus-send --session --dest=org.freedesktop.ExampleService --print-reply /org/freedesktop/sample/object org.freedesktop.ExampleInterface.ExampleMethod int32:47",
			"dbus-send --session --dest=<val> --print-reply <path> <dotted-id> int32:<val>",
		},
		{
			"boolean argument",
			"dbus-send --session --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.PlayPause boolean:true",
			"dbus-send --session --dest=<val> <path> <dotted-id> boolean:<val>",
		},
		{
			"multiple typed args",
			"dbus-send --session --dest=org.example.Service /Path org.example.Method string:hello int32:42 double:3.14",
			"dbus-send --session --dest=<val> <path> <dotted-id> string:<val> int32:<val> double:<val>",
		},
		{
			"reply timeout",
			"dbus-send --session --dest=org.example.Service --reply-timeout=5000 --print-reply /Path org.example.Method",
			"dbus-send --session --dest=<val> --reply-timeout=N --print-reply <path> <dotted-id>",
		},
		{
			"print-reply with literal",
			"dbus-send --session --dest=org.example.Service --print-reply=literal /Path org.example.Method",
			"dbus-send --session --dest=<val> --print-reply=literal <path> <dotted-id>",
		},
		{
			"objpath typed arg",
			"dbus-send --session --dest=org.example.Service /Path org.example.Method objpath:/org/example/object",
			"dbus-send --session --dest=<val> <path> <dotted-id> objpath:<val>",
		},
		{
			"byte and uint32 args",
			"dbus-send --session --dest=org.example.Service /Path org.example.Method byte:255 uint32:1000",
			"dbus-send --session --dest=<val> <path> <dotted-id> byte:<val> uint32:<val>",
		},
		{
			"no typed args",
			"dbus-send --session --dest=org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.ListNames",
			"dbus-send --session --dest=<val> <path> <dotted-id>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different string values collide", func(t *testing.T) {
		a := shellshape.Normalize(`dbus-send --session --dest=org.kde.KWin /Scripting org.kde.kwin.Scripting.loadScript string:"close-lmstudio"`)
		b := shellshape.Normalize(`dbus-send --session --dest=org.kde.KWin /Scripting org.kde.kwin.Scripting.loadScript string:"close-lmstudio2"`)
		if a != b {
			t.Errorf("expected same shape:\n  a = %q\n  b = %q", a, b)
		}
	})

	t.Run("different int32 values collide", func(t *testing.T) {
		a := shellshape.Normalize("dbus-send --session --dest=org.example.Svc /Path org.example.Method int32:42")
		b := shellshape.Normalize("dbus-send --session --dest=org.example.Svc /Path org.example.Method int32:99")
		if a != b {
			t.Errorf("expected same shape:\n  a = %q\n  b = %q", a, b)
		}
	})

	t.Run("different dest values collide", func(t *testing.T) {
		a := shellshape.Normalize("dbus-send --session --dest=org.kde.KWin /Path org.example.Method")
		b := shellshape.Normalize("dbus-send --session --dest=org.freedesktop.DBus /Path org.example.Method")
		if a != b {
			t.Errorf("expected same shape:\n  a = %q\n  b = %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("dbus-send --session --dest=org.kde.KWin /Path org.example.Method string:hello")
		subshell := shellshape.Normalize("dbus-send --session --dest=org.kde.KWin /Path org.example.Method $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
