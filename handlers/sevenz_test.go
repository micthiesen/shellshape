package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSevenZ(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"add archive", "7z a archive.7z file1 file2", "7z a <path>+"},
		{"extract", "7z x archive.7z", "7z x <path>"},
		{"list", "7z l archive.7z", "7z l <path>"},
		{"test", "7z t archive.7z", "7z t <path>"},
		{"extract flat", "7z e archive.7z", "7z e <path>"},
		{"delete from archive", "7z d archive.7z file.txt", "7z d <path>+"},
		{"update archive", "7z u archive.7z newfile", "7z u <path>+"},

		// Fused flags
		{"output dir", "7z x archive.7z -oOutDir", "7z x <path> -o<path>"},
		{"password", "7z x archive.7z -psecret", "7z x <path> -p<val>"},
		{"compression level", "7z a -mx=9 archive.7z files/", "7z a -mx=N <path>+"},
		{"archive type", "7z a -t7z archive.7z dir/", "7z a -t7z <path>+"},
		{"archive type zip", "7z a -tzip archive.zip dir/", "7z a -tzip <path>+"},

		// Boolean flags
		{"yes flag", "7z x -y archive.7z", "7z x -y <path>"},

		// Aliases
		{"7za", "7za a archive.7z file", "7za a <path>+"},
		{"7zr", "7zr a archive.7z file", "7zr a <path>+"},

		// Combined
		{"complex extract", "7z x -y -oOutDir -psecret archive.7z", "7z x -y -o<path> -p<val> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different archives collide", func(t *testing.T) {
		a := shellshape.Normalize("7z x backup.7z")
		b := shellshape.Normalize("7z x photos.tar")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("7z x archive.7z")
		subshell := shellshape.Normalize("7z x $(echo archive.7z)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
