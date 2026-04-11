package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestCmake(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic configure usage
		{"configure with path", "cmake /home/user/myproject", "cmake <path>"},
		{"configure with dot", "cmake .", "cmake ."},
		{"configure relative path", "cmake ../src", "cmake <path>"},

		// -D variable definitions
		{"fused -D flag", "cmake -DCMAKE_BUILD_TYPE=Release .", "cmake -D <val> ."},
		{"separate -D flag", "cmake -D CMAKE_BUILD_TYPE=Release .", "cmake -D <val> ."},
		{"multiple -D flags", "cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/usr/local .", "cmake -D <val> -D <val> ."},
		{"typed -D flag", "cmake -DCMAKE_CXX_COMPILER:FILEPATH=/usr/bin/g++ .", "cmake -D <val> ."},

		// Source and build directories
		{"source and build dirs", "cmake -S src -B build", "cmake -S <path> -B <path>"},
		{"source and build full paths", "cmake -S /home/user/project -B /tmp/build", "cmake -S <path> -B <path>"},

		// Generator selection
		{"generator flag", `cmake -G "Unix Makefiles" .`, "cmake -G <val> ."},
		{"generator with toolset", `cmake -G "Visual Studio 16" -T v142 -A x64 .`, "cmake -G <val> -T <val> -A <val> ."},

		// Cache and toolchain
		{"cache preload", "cmake -C /path/to/cache.cmake .", "cmake -C <path> ."},
		{"toolchain file", "cmake --toolchain /path/to/toolchain.cmake .", "cmake --toolchain <path> ."},
		{"install prefix", "cmake --install-prefix /usr/local .", "cmake --install-prefix <path> ."},

		// Preset
		{"preset", "cmake --preset=release", "cmake --preset=<val>"},
		{"preset separate", "cmake --preset release", "cmake --preset <val>"},

		// Undefine
		{"undefine variable", "cmake -U CMAKE_BUILD_TYPE .", "cmake -U <val> ."},

		// Warning flags (boolean, kept verbatim)
		{"warning flags", "cmake -Wdev -Werror=dev .", "cmake -Wdev -Werror=dev ."},

		// Boolean flags
		{"fresh flag", "cmake --fresh -S src -B build", "cmake --fresh -S <path> -B <path>"},
		{"trace flag", "cmake --trace .", "cmake --trace ."},
		{"debug output", "cmake --debug-output .", "cmake --debug-output ."},

		// Log level
		{"log level fused", "cmake --log-level=DEBUG .", "cmake --log-level=<val> ."},
		{"trace format", "cmake --trace-format=json-v1 .", "cmake --trace-format=<val> ."},
		{"trace source", "cmake --trace-source /path/to/file.cmake .", "cmake --trace-source <path> ."},
		{"graphviz", "cmake --graphviz=deps.dot .", "cmake --graphviz=<path> ."},

		// --build mode
		{"build simple", "cmake --build .", "cmake --build ."},
		{"build with path", "cmake --build /tmp/build", "cmake --build <path>"},
		{"build with target", "cmake --build . --target install", "cmake --build . --target <val>"},
		{"build short target", "cmake --build . -t myapp", "cmake --build . -t <val>"},
		{"build with config", "cmake --build build --config Release", "cmake --build build --config <val>"},
		{"build parallel", "cmake --build . -j 4", "cmake --build . -j N"},
		{"build parallel long", "cmake --build . --parallel 8", "cmake --build . --parallel N"},
		{"build clean first", "cmake --build . --clean-first", "cmake --build . --clean-first"},
		{"build verbose", "cmake --build . -v", "cmake --build . -v"},
		{"build full", "cmake --build build --target all --config Debug -j 8 --clean-first -v", "cmake --build build --target <val> --config <val> -j N --clean-first -v"},

		// --install mode
		{"install simple", "cmake --install ./build", "cmake --install <path>"},
		{"install with prefix", "cmake --install ./build --prefix /usr/local", "cmake --install <path> --prefix <path>"},
		{"install with strip", "cmake --install ./build --strip", "cmake --install <path> --strip"},
		{"install with component", "cmake --install ./build --component Runtime", "cmake --install <path> --component <val>"},
		{"install full", "cmake --install ./build --prefix /opt/app --config Release --component Runtime --strip", "cmake --install <path> --prefix <path> --config <val> --component <val> --strip"},

		// -P script mode
		{"script mode", "cmake -P build_script.cmake", "cmake -P <path>"},
		{"script with vars", "cmake -DVAR=value -P script.cmake", "cmake -D <val> -P <path>"},

		// -E tool mode
		{"E copy", "cmake -E copy file1.txt file2.txt", "cmake -E copy <path>+"},
		{"E make_directory", "cmake -E make_directory /tmp/output", "cmake -E make_directory <path>"},
		{"E echo", "cmake -E echo hello", "cmake -E echo hello"},
		{"E remove", "cmake -E rm -rf /tmp/build", "cmake -E rm -rf <path>"},
		{"E chdir", "cmake -E chdir /tmp/build make", "cmake -E chdir <path> make"},
		{"E tar", "cmake -E tar czf archive.tar.gz src/", "cmake -E tar czf <path>+"},
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
	t.Run("different -D values collide", func(t *testing.T) {
		a := shellshape.Normalize("cmake -DCMAKE_BUILD_TYPE=Release .")
		b := shellshape.Normalize("cmake -DCMAKE_BUILD_TYPE=Debug .")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different generators collide", func(t *testing.T) {
		a := shellshape.Normalize(`cmake -G "Unix Makefiles" .`)
		b := shellshape.Normalize(`cmake -G Ninja .`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different build directories collide", func(t *testing.T) {
		a := shellshape.Normalize("cmake --build /home/user/project/build")
		b := shellshape.Normalize("cmake --build /tmp/other-build")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("cmake literal-arg")
		subshell := shellshape.Normalize("cmake $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
