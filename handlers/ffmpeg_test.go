package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFfmpeg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple transcode", "ffmpeg -i input.mp4 output.avi", "ffmpeg -i <path>+"},
		{"overwrite flag", "ffmpeg -y -i input.mp4 output.mp4", "ffmpeg -y -i <path>+"},

		// Codec flags kept verbatim
		{"video codec", "ffmpeg -i input.mp4 -c:v libx264 output.mp4", "ffmpeg -i <path> -c:v libx264 <path>"},
		{"audio codec", "ffmpeg -i input.mp4 -c:a aac output.mp4", "ffmpeg -i <path> -c:a aac <path>"},
		{"codec copy", "ffmpeg -i input.mp4 -c copy output.mkv", "ffmpeg -i <path> -c copy <path>"},
		{"vcodec alias", "ffmpeg -i input.mp4 -vcodec libx265 output.mp4", "ffmpeg -i <path> -vcodec libx265 <path>"},
		{"acodec alias", "ffmpeg -i input.mp4 -acodec libopus output.ogg", "ffmpeg -i <path> -acodec libopus <path>"},

		// Format flag kept verbatim
		{"format flag", "ffmpeg -f rawvideo -i input.raw output.mp4", "ffmpeg -f rawvideo -i <path>+"},

		// Bitrate collapsed
		{"video bitrate", "ffmpeg -i input.mp4 -b:v 5M output.mp4", "ffmpeg -i <path> -b:v <val> <path>"},
		{"audio bitrate", "ffmpeg -i input.mp4 -ab 128k output.mp4", "ffmpeg -i <path> -ab <val> <path>"},

		// Numeric flags
		{"framerate", "ffmpeg -i input.mp4 -r 30 output.mp4", "ffmpeg -i <path> -r N <path>"},
		{"audio rate", "ffmpeg -i input.mp4 -ar 44100 output.mp4", "ffmpeg -i <path> -ar N <path>"},
		{"audio channels", "ffmpeg -i input.mp4 -ac 2 output.mp4", "ffmpeg -i <path> -ac N <path>"},

		// Timestamps and duration
		{"seek start", "ffmpeg -i input.mp4 -ss 00:01:30 output.mp4", "ffmpeg -i <path> -ss <val> <path>"},
		{"seek to", "ffmpeg -i input.mp4 -to 00:05:00 output.mp4", "ffmpeg -i <path> -to <val> <path>"},
		{"duration", "ffmpeg -i input.mp4 -t 60 output.mp4", "ffmpeg -i <path> -t <val> <path>"},

		// Resolution
		{"resolution", "ffmpeg -i input.mp4 -s 1280x720 output.mp4", "ffmpeg -i <path> -s <val> <path>"},

		// Filters
		{"video filter", "ffmpeg -i input.mp4 -vf scale=1920:1080 output.mp4", "ffmpeg -i <path> -vf <filter> <path>"},
		{"audio filter", "ffmpeg -i input.mp4 -af volume=0.5 output.mp4", "ffmpeg -i <path> -af <filter> <path>"},

		// Map
		{"map stream", "ffmpeg -i input.mp4 -map 0:v -map 0:a output.mkv", "ffmpeg -i <path> -map <val> -map <val> <path>"},

		// Multiple inputs
		{"two inputs", "ffmpeg -i video.mp4 -i audio.aac -c copy output.mkv", "ffmpeg -i <path> -i <path> -c copy <path>"},

		// Boolean flags
		{"disable video", "ffmpeg -i input.mp4 -vn output.mp3", "ffmpeg -i <path> -vn <path>"},
		{"disable audio", "ffmpeg -i input.mp4 -an output.mp4", "ffmpeg -i <path> -an <path>"},

		// Loglevel
		{"loglevel", "ffmpeg -v quiet -i input.mp4 output.mp4", "ffmpeg -v quiet -i <path>+"},

		// Metadata
		{"metadata", "ffmpeg -i input.mp4 -metadata title=MyVideo output.mp4", "ffmpeg -i <path> -metadata <val> <path>"},

		// Complex real-world
		{"full transcode", "ffmpeg -y -i /tmp/src.mp4 -c:v libx264 -b:v 2M -c:a aac -b:a 128k -r 24 /tmp/out.mp4",
			"ffmpeg -y -i <path> -c:v libx264 -b:v <val> -c:a aac -b:a <val> -r N <path>"},

		// ffprobe alias
		{"ffprobe basic", "ffprobe input.mp4", "ffprobe <path>"},
		{"ffprobe with flags", "ffprobe -v quiet -show_format input.mp4", "ffprobe -v quiet -show_format <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different inputs collide", func(t *testing.T) {
		a := shellshape.Normalize("ffmpeg -i /path/to/video.mp4 output.avi")
		b := shellshape.Normalize("ffmpeg -i /other/file.mkv result.mp4")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different bitrates collide", func(t *testing.T) {
		a := shellshape.Normalize("ffmpeg -i in.mp4 -b:v 2M out.mp4")
		b := shellshape.Normalize("ffmpeg -i in.mp4 -b:v 5M out.mp4")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ffmpeg -i input.mp4 output.mp4")
		subshell := shellshape.Normalize("ffmpeg -i $(echo input.mp4) output.mp4")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
