package shellshape

import "testing"

func TestStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple file", "strings /usr/bin/ls", "strings <path>"},
		{"relative file", "strings binary.o", "strings <path>"},

		// Flags with arguments
		{"-n minimum length", "strings -n 8 binary.o", "strings -n N <path>"},
		{"--bytes minimum length", "strings --bytes 10 binary.o", "strings --bytes N <path>"},
		{"-t radix", "strings -t x myfile", "strings -t x <path>"},
		{"-e encoding", "strings -e S file.bin", "strings -e S <path>"},

		// Boolean flags
		{"-a flag", "strings -a foo.bin", "strings -a <path>"},
		{"-o flag", "strings -a -o foo.bin", "strings -a -o <path>"},
		{"-f flag", "strings -f binary.o", "strings -f <path>"},

		// Multiple positionals
		{"multiple files", "strings file1.o file2.o file3.o", "strings <path>+"},
		{"flags and multiple files", "strings -a -n 6 file1.o file2.o", "strings -a -n N <path>+"},

		// Double dash
		{"double dash", "strings -- file1.o file2.o", "strings -- <path>+"},

		// Combined flags and arguments
		{"combined", "strings -a -t o -n 10 binary.o", "strings -a -t o -n N <path>"},

		// Redirect
		{"redirect", "strings binary.o > output.txt", "strings <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("strings -n 8 /usr/bin/ls")
		b := Normalize("strings -n 4 /usr/bin/cat")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("strings literal-arg")
		subshell := Normalize("strings $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
