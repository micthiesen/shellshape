package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage (target is consumed as subcommand by framework)
		{"bare make", "make", "make"},
		{"single target", "make build", "make build"},
		{"multiple targets", "make clean all", "make clean all"},

		// Flags after target (handler processes these)
		{"jobs after target", "make build -j 8", "make build -j N"},
		{"makefile after target", "make build -f custom.mk", "make build -f <path>"},
		{"directory after target", "make build -C src/dir", "make build -C <path>"},
		{"include dir after target", "make build -I /usr/include", "make build -I <path>"},
		{"keep going after target", "make build -k", "make build -k"},
		{"dry run after target", "make build -n", "make build -n"},
		{"always make after target", "make build -B", "make build -B"},
		{"silent after target", "make build -s", "make build -s"},

		// Long flags with =value
		{"long file flag", "make build --file=Makefile.dev", "make build --file=<path>"},
		{"long directory flag", "make build --directory=/tmp/build", "make build --directory=<path>"},
		{"long jobs flag", "make build --jobs=4", "make build --jobs=N"},

		// Variable assignments
		{"variable assignment", "make build CC=gcc", "make build CC=<val>"},
		{"multiple vars", "make build CC=gcc CFLAGS=-O2", "make build CC=<val> CFLAGS=<val>"},
		{"var with path value", "make install DESTDIR=/usr/local", "make install DESTDIR=<val>"},

		// Combined flags and vars
		{"flags and vars", "make build -k -j 4 CC=clang", "make build -k -j N CC=<val>"},
		{"multiple flags after target", "make test -j 8 -C src/dir", "make test -j N -C <path>"},

		// Old-file and what-if flags
		{"old-file flag", "make build -o main.o", "make build -o <path>"},
		{"what-if flag", "make build -W header.h", "make build -W <path>"},

		// Flags before target (flag becomes subcommand, handler sees dangling arg)
		{"jobs before target", "make -j 8 build", "make -j N build"},
		{"makefile before target", "make -f custom.mk build", "make -f <path> build"},
		{"directory before target", "make -C /some/dir build", "make -C <path> build"},
		{"boolean flag before target", "make -k test", "make -k test"},

		// Redirects
		{"with redirect", "make build > log.txt", "make build > <path>"},
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
	t.Run("different variable values collide", func(t *testing.T) {
		a := shellshape.Normalize("make build CC=gcc")
		b := shellshape.Normalize("make build CC=clang")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different job counts collide", func(t *testing.T) {
		a := shellshape.Normalize("make build -j 4")
		b := shellshape.Normalize("make build -j 16")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different makefiles collide", func(t *testing.T) {
		a := shellshape.Normalize("make build -f Makefile.dev")
		b := shellshape.Normalize("make build -f Makefile.prod")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("make literal-arg")
		subshell := shellshape.Normalize("make $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
