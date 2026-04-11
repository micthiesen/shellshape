package shellshape

import "testing"

func TestLdd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single binary", "ldd /usr/bin/ls", "ldd <path>"},
		{"relative path", "ldd ./myapp", "ldd <path>"},
		{"bare name", "ldd mylib.so", "ldd <path>"},
		{"multiple binaries", "ldd /usr/lib/libfoo.so /usr/lib/libbar.so", "ldd <path>+"},

		// Boolean flags
		{"verbose short", "ldd -v /usr/bin/ls", "ldd -v <path>"},
		{"verbose long", "ldd --verbose /usr/bin/ls", "ldd --verbose <path>"},
		{"unused short", "ldd -u /usr/bin/ls", "ldd -u <path>"},
		{"unused long", "ldd --unused /usr/bin/ls", "ldd --unused <path>"},
		{"data relocs short", "ldd -d /usr/bin/ls", "ldd -d <path>"},
		{"data relocs long", "ldd --data-relocs /usr/bin/ls", "ldd --data-relocs <path>"},
		{"function relocs short", "ldd -r /usr/bin/ls", "ldd -r <path>"},
		{"function relocs long", "ldd --function-relocs /usr/bin/ls", "ldd --function-relocs <path>"},

		// Multiple flags
		{"verbose and unused", "ldd -v -u /usr/lib/libc.so", "ldd -v -u <path>"},

		// Redirects
		{"with redirect", "ldd /usr/bin/ls > /tmp/deps.txt", "ldd <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different binaries collide", func(t *testing.T) {
		a := Normalize("ldd /usr/bin/ls")
		b := Normalize("ldd /usr/bin/cat")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("ldd literal-arg")
		subshell := Normalize("ldd $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestOtool(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"show libs", "otool -L /usr/lib/libfoo.dylib", "otool -L <path>"},
		{"load commands", "otool -l myapp", "otool -l <path>"},
		{"mach header", "otool -h myapp", "otool -h <path>"},
		{"multiple files", "otool -L libfoo.dylib libbar.dylib", "otool -L <path>+"},

		// Boolean flags
		{"verbose", "otool -v -l myapp", "otool -v -l <path>"},
		{"symbolic disassembly", "otool -V -t myapp", "otool -V -t <path>"},
		{"no addresses", "otool -X -t myapp", "otool -X -t <path>"},
		{"combined flags", "otool -tVX myapp", "otool -tVX <path>"},
		{"archive header", "otool -a libfoo.a", "otool -a <path>"},

		// Flags with arguments
		{"arch flag", "otool -arch x86_64 myapp", "otool -arch <val> <path>"},
		{"arch arm64", "otool -arch arm64 libfoo.dylib", "otool -arch <val> <path>"},
		{"p flag symbol", "otool -p _main -t -v myapp", "otool -p <sym> -t -v <path>"},
		{"s flag two args", "otool -s __TEXT __text myapp", "otool -s <val>+ <path>"},

		// Mixed
		{"arch and libs", "otool -arch x86_64 -L myapp", "otool -arch <val> -L <path>"},
		{"verbose with arch", "otool -v -arch arm64 -l myapp", "otool -v -arch <val> -l <path>"},

		// Redirects
		{"with redirect", "otool -L myapp > /tmp/deps.txt", "otool -L <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different binaries collide", func(t *testing.T) {
		a := Normalize("otool -L /usr/lib/libfoo.dylib")
		b := Normalize("otool -L /Applications/MyApp.app/Contents/MacOS/MyApp")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different arch values collide", func(t *testing.T) {
		a := Normalize("otool -arch x86_64 -L myapp")
		b := Normalize("otool -arch arm64 -L myapp")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("otool -L literal-arg")
		subshell := Normalize("otool -L $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
