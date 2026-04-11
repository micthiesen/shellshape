package shellshape

import "testing"

func TestClang(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"compile to executable", "clang main.c -o main", "clang <path> -o <path>"},
		{"compile only", "clang -c foo.c", "clang -c <path>"},
		{"compile multiple sources", "clang foo.c bar.c -o app", "clang <path>+ -o <path>"},
		{"bare clang", "clang main.c", "clang <path>"},

		// Common flag combinations
		{"warnings and optimization", "clang -Wall -Wextra -O2 -c file.c", "clang -Wall -Wextra -O2 -c <path>"},
		{"debug build", "clang -g -O0 -Wall main.c -o main", "clang -g -O0 -Wall <path> -o <path>"},
		{"emit llvm ir", "clang -S -emit-llvm file.c -o file.ll", "clang -S -emit-llvm <path> -o <path>"},

		// Flags with arguments (space-separated)
		{"include path", "clang -I /usr/include main.c", "clang -I <path>+"},
		{"lib path", "clang -L /usr/lib main.c -lm", "clang -L <path>+ -l <val>"},
		{"define macro", "clang -D DEBUG main.c", "clang -D <def> <path>"},
		{"undefine macro", "clang -U NDEBUG main.c", "clang -U <def> <path>"},
		{"output flag", "clang main.c -o /tmp/build/app", "clang <path> -o <path>"},
		{"target flag", "clang --target x86_64-linux-gnu main.c", "clang --target <val> <path>"},
		{"std flag", "clang -std c11 main.c", "clang -std <val> <path>"},
		{"isystem flag", "clang -isystem /usr/local/include main.c", "clang -isystem <path>+"},
		{"include file", "clang -include config.h main.c", "clang -include <path>+"},
		{"arch flag", "clang -arch arm64 main.c", "clang -arch <val> <path>"},
		{"linker flag", "clang -Xlinker --no-as-needed main.c", "clang -Xlinker <val> <path>"},
		{"x language", "clang -x c main.txt -o main", "clang -x <val> <path> -o <path>"},
		{"linker script", "clang -T linker.ld main.o -o firmware", "clang -T <val> <dotted-id> -o <path>"},

		// Fused flags (no space between flag and value)
		{"fused include", "clang -I/usr/include main.c", "clang -I <path>+"},
		{"fused lib path", "clang -L/usr/lib main.c", "clang -L <path>+"},
		{"fused define", "clang -DFOO=bar main.c", "clang -D <def> <path>"},
		{"fused define no value", "clang -DNDEBUG main.c", "clang -D <def> <path>"},
		{"fused undefine", "clang -UNDEBUG main.c", "clang -U <def> <path>"},
		{"fused link lib", "clang main.c -lm", "clang <path> -l <val>"},
		{"fused std", "clang -std=c++17 main.cpp", "clang -std=<val> <path>"},
		{"fused target", "clang --target=x86_64-linux-gnu main.c", "clang --target=<val> <path>"},
		{"fused sysroot", "clang --sysroot=/opt/sdk main.c", "clang --sysroot=<path> <path>"},
		{"fused march", "clang -march=native main.c", "clang -march=<val> <path>"},

		// clang++ usage
		{"clang++ basic", "clang++ main.cpp -o app", "clang++ <path> -o <path>"},
		{"clang++ with std", "clang++ -std=c++17 -O2 main.cpp -o app", "clang++ -std=<val> -O2 <path> -o <path>"},

		// gcc/g++ aliases
		{"gcc basic", "gcc main.c -o main", "gcc <path> -o <path>"},
		{"g++ basic", "g++ main.cpp -o app", "g++ <path> -o <path>"},
		{"cc basic", "cc main.c -o main", "cc <path> -o <path>"},

		// Optimization levels kept verbatim
		{"O0", "clang -O0 main.c", "clang -O0 <path>"},
		{"O1", "clang -O1 main.c", "clang -O1 <path>"},
		{"O2", "clang -O2 main.c", "clang -O2 <path>"},
		{"O3", "clang -O3 main.c", "clang -O3 <path>"},
		{"Os", "clang -Os main.c", "clang -Os <path>"},
		{"Oz", "clang -Oz main.c", "clang -Oz <path>"},
		{"Og", "clang -Og main.c", "clang -Og <path>"},
		{"Ofast", "clang -Ofast main.c", "clang -Ofast <path>"},

		// Warning flags kept verbatim
		{"Wall", "clang -Wall main.c", "clang -Wall <path>"},
		{"Werror", "clang -Werror main.c", "clang -Werror <path>"},

		// MF/MQ/MT dependency flags
		{"dependency output", "clang -MMD -MF deps/main.d -c main.c", "clang -MMD -MF <path> -c <path>"},
		{"dependency target", "clang -MT build/main.o -c main.c", "clang -MT <path> -c <path>"},

		// Redirects
		{"with redirect", "clang main.c -o main 2>&1", "clang <path> -o <path> 2>&1"},
		{"redirect to file", "clang main.c -o main 2> err.log", "clang <path> -o <path> 2> <path>"},

		// Complex real-world
		{"full build", "gcc -g -Wall -Werror -O2 -I./include -DVERSION=1 -L./lib -lssl -c src/main.c -o build/main.o", "gcc -g -Wall -Werror -O2 -I <path> -D <def> -L <path> -l <val> -c <path> -o <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different source files collide", func(t *testing.T) {
		a := Normalize("clang foo.c -o app")
		b := Normalize("clang bar.c -o app")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different define values collide", func(t *testing.T) {
		a := Normalize("clang -DFOO=1 main.c")
		b := Normalize("clang -DBAR=2 main.c")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different include paths collide", func(t *testing.T) {
		a := Normalize("clang -I/usr/include main.c")
		b := Normalize("clang -I/opt/local/include main.c")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different libraries collide", func(t *testing.T) {
		a := Normalize("clang main.c -lm")
		b := Normalize("clang main.c -lpthread")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("clang literal-arg")
		subshell := Normalize("clang $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
