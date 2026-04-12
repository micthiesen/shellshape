package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFfprobe(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"basic file", "ffprobe input.mp4", "ffprobe <path>"},
		{"basic url", "ffprobe https://example.com/video.mp4", "ffprobe <https-uri>"},

		// Verbosity
		{"quiet mode", "ffprobe -v quiet input.mp4", "ffprobe -v quiet <path>"},
		{"error level", "ffprobe -v error input.mp4", "ffprobe -v error <path>"},

		// Show flags (boolean/structural)
		{"show format", "ffprobe -show_format input.mp4", "ffprobe -show_format <path>"},
		{"show streams", "ffprobe -show_streams input.mp4", "ffprobe -show_streams <path>"},
		{"show entries", "ffprobe -show_entries format=duration input.mp4", "ffprobe -show_entries format=duration <path>"},
		{"show format and streams", "ffprobe -v quiet -show_format -show_streams input.mp4", "ffprobe -v quiet -show_format -show_streams <path>"},

		// Output format
		{"json output", "ffprobe -of json input.mp4", "ffprobe -of json <path>"},
		{"csv output", "ffprobe -of csv input.mp4", "ffprobe -of csv <path>"},
		{"print format", "ffprobe -print_format flat input.mp4", "ffprobe -print_format flat <path>"},

		// Select streams
		{"select video", "ffprobe -select_streams v input.mp4", "ffprobe -select_streams <val> <path>"},
		{"select audio", "ffprobe -select_streams a:0 input.mp4", "ffprobe -select_streams <val> <path>"},

		// Input flag
		{"input flag", "ffprobe -i input.mp4", "ffprobe -i <path>"},

		// Combined typical usage
		{"full probe", "ffprobe -v quiet -show_format -show_streams -of json input.mp4", "ffprobe -v quiet -show_format -show_streams -of json <path>"},

		// With redirect
		{"redirect output", "ffprobe -v quiet input.mp4 > info.json", "ffprobe -v quiet <path> > <path>"},
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
		a := shellshape.Normalize("ffprobe -v quiet /tmp/video1.mp4")
		b := shellshape.Normalize("ffprobe -v quiet /tmp/video2.mkv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ffprobe literal-arg")
		subshell := shellshape.Normalize("ffprobe $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
