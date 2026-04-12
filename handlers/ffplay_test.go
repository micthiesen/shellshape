package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFfplay(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"play file", "ffplay video.mp4", "ffplay <path>"},
		{"play url", "ffplay https://example.com/stream.m3u8", "ffplay <https-uri>"},

		// Boolean flags
		{"autoexit", "ffplay -autoexit video.mp4", "ffplay -autoexit <path>"},
		{"fullscreen", "ffplay -fs video.mp4", "ffplay -fs <path>"},
		{"no audio", "ffplay -an video.mp4", "ffplay -an <path>"},
		{"no video", "ffplay -vn music.mp3", "ffplay -vn <path>"},

		// Value flags
		{"seek start", "ffplay -ss 01:30 video.mp4", "ffplay -ss <val> <path>"},
		{"duration", "ffplay -t 60 video.mp4", "ffplay -t <val> <path>"},

		// Numeric flags
		{"volume", "ffplay -volume 50 video.mp4", "ffplay -volume N <path>"},
		{"width", "ffplay -x 1280 video.mp4", "ffplay -x N <path>"},
		{"height", "ffplay -y 720 video.mp4", "ffplay -y N <path>"},
		{"loop", "ffplay -loop 0 video.mp4", "ffplay -loop N <path>"},

		// Combined flags
		{"fullscreen autoexit", "ffplay -fs -autoexit video.mp4", "ffplay -fs -autoexit <path>"},
		{"sized window", "ffplay -x 1280 -y 720 -autoexit video.mp4", "ffplay -x N -y N -autoexit <path>"},

		// With redirect
		{"redirect stderr", "ffplay -v quiet video.mp4 2>/dev/null", "ffplay -v quiet <path> 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("ffplay -autoexit /tmp/video1.mp4")
		b := shellshape.Normalize("ffplay -autoexit /tmp/video2.mkv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ffplay literal-arg")
		subshell := shellshape.Normalize("ffplay $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
