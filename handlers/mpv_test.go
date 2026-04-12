package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMpv(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"play file", "mpv video.mp4", "mpv <path>"},
		{"play url", "mpv https://example.com/stream.m3u8", "mpv <https-uri>"},

		// Boolean flags
		{"fullscreen", "mpv --fullscreen video.mp4", "mpv --fullscreen <path>"},
		{"no video", "mpv --no-video music.flac", "mpv --no-video <path>"},
		{"no audio", "mpv --no-audio video.mp4", "mpv --no-audio <path>"},

		// Fused value flags
		{"volume", "mpv --volume=50 video.mp4", "mpv --volume=N <path>"},
		{"start time", "mpv --start=01:30 video.mp4", "mpv --start=<val> <path>"},
		{"length", "mpv --length=00:30 video.mp4", "mpv --length=<val> <path>"},
		{"sub file", "mpv --sub-file=/tmp/subs.srt video.mp4", "mpv --sub-file=<path> <path>"},

		// Profile (structural, kept verbatim)
		{"profile", "mpv --profile=gpu-hq video.mp4", "mpv --profile=gpu-hq <path>"},

		// Generic long flags with values get <val>
		{"screen", "mpv --screen=1 video.mp4", "mpv --screen=<val> <path>"},

		// Combined
		{"fullscreen volume", "mpv --fullscreen --volume=80 video.mp4", "mpv --fullscreen --volume=N <path>"},
		{"no video start", "mpv --no-video --start=00:05:00 music.mp3", "mpv --no-video --start=<val> <path>"},

		// Multiple positionals
		{"playlist files", "mpv video1.mp4 video2.mkv", "mpv <path>+"},
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
		a := shellshape.Normalize("mpv --fullscreen /tmp/movie1.mp4")
		b := shellshape.Normalize("mpv --fullscreen /tmp/movie2.mkv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different volumes collide", func(t *testing.T) {
		a := shellshape.Normalize("mpv --volume=30 video.mp4")
		b := shellshape.Normalize("mpv --volume=80 video.mp4")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("mpv literal-arg")
		subshell := shellshape.Normalize("mpv $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
