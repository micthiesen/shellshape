package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSwift(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// swift build
		{"build bare", "swift build", "swift build"},
		{"build verbose", "swift build -v", "swift build -v"},
		{"build release", "swift build -c release", "swift build -c <val>"},
		{"build configuration long", "swift build --configuration debug", "swift build --configuration <val>"},
		{"build product", "swift build --product MyApp", "swift build --product <val>"},
		{"build target", "swift build --target MyLib", "swift build --target <val>"},
		{"build package-path", "swift build --package-path /tmp/myproject", "swift build --package-path <path>"},
		{"build scratch-path", "swift build --scratch-path .build-debug", "swift build --scratch-path <path>"},
		{"build jobs", "swift build -j 8", "swift build -j N"},
		{"build jobs long", "swift build --jobs 4", "swift build --jobs N"},
		{"build swift-sdk", "swift build --swift-sdk wasm32-unknown-wasi", "swift build --swift-sdk <val>"},
		{"build traits", "swift build --traits Trait1,Trait2", "swift build --traits <val>"},
		{"build static stdlib", "swift build --static-swift-stdlib", "swift build --static-swift-stdlib"},

		// swift test
		{"test bare", "swift test", "swift test"},
		{"test verbose", "swift test -v", "swift test -v"},
		{"test filter", "swift test --filter MyTests", "swift test --filter <pattern>"},
		{"test skip", "swift test --skip PerformanceTests", "swift test --skip <pattern>"},
		{"test parallel", "swift test --parallel", "swift test --parallel"},
		{"test num-workers", "swift test --num-workers 4", "swift test --num-workers N"},
		{"test xunit-output", "swift test --xunit-output results.xml", "swift test --xunit-output <path>"},
		{"test code-coverage", "swift test --enable-code-coverage", "swift test --enable-code-coverage"},
		{"test skip-build", "swift test --skip-build", "swift test --skip-build"},
		{"test list", "swift test -l", "swift test -l"},
		{"test specifier short", "swift test -s MyTarget.MyTestCase", "swift test -s <val>"},
		{"test attachments-path", "swift test --attachments-path /tmp/att", "swift test --attachments-path <path>"},

		// swift run
		{"run bare", "swift run", "swift run"},
		{"run executable", "swift run MyApp", "swift run MyApp"},
		{"run executable with args", "swift run MyApp --port 8080 --host localhost", "swift run MyApp --port 8080 --host localhost"},
		{"run skip-build", "swift run --skip-build MyApp", "swift run --skip-build MyApp"},
		{"run with config", "swift run -c release MyApp arg1 arg2", "swift run -c <val> MyApp arg1 arg2"},
		{"run package-path", "swift run --package-path /tmp/proj MyApp", "swift run --package-path <path> MyApp"},

		// swift package
		{"package init", "swift package init", "swift package init"},
		{"package init type", "swift package init --type library", "swift package init --type <val>"},
		{"package init name", "swift package init --name MyLib", "swift package init --name <val>"},
		{"package update", "swift package update", "swift package update"},
		{"package resolve", "swift package resolve", "swift package resolve"},
		{"package clean", "swift package clean", "swift package clean"},
		{"package reset", "swift package reset", "swift package reset"},
		{"package show-dependencies", "swift package show-dependencies", "swift package show-dependencies"},
		{"package dump-package", "swift package dump-package", "swift package dump-package"},

		// swift repl
		{"repl bare", "swift repl", "swift repl"},

		// swift with no subcommand (compiler mode, like swiftc)
		{"swift file", "swift main.swift", "swift <path>"},
		{"swift file with -O", "swift -O main.swift", "swift -O <path>"},
		{"swift with -e", "swift -e print(42)", "swift -e <val>"},

		// Combined flags
		{"build combined", "swift build -c release --product MyTool -j 4", "swift build -c <val> --product <val> -j N"},
		{"test combined", "swift test --filter UnitTests --parallel --num-workers 8", "swift test --filter <pattern> --parallel --num-workers N"},

		// Fused flags with =
		{"build config fused", "swift build --configuration=release", "swift build --configuration=<val>"},
		{"test filter fused", "swift test --filter=MyTests", "swift test --filter=<pattern>"},
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
	t.Run("different configurations collide", func(t *testing.T) {
		a := shellshape.Normalize("swift build -c release")
		b := shellshape.Normalize("swift build -c debug")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different test filters collide", func(t *testing.T) {
		a := shellshape.Normalize("swift test --filter UnitTests")
		b := shellshape.Normalize("swift test --filter IntegrationTests")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different package paths collide", func(t *testing.T) {
		a := shellshape.Normalize("swift build --package-path /Users/alice/proj")
		b := shellshape.Normalize("swift build --package-path /Users/bob/proj")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("swift build --product MyApp")
		subshell := shellshape.Normalize("swift build --product $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestSwiftc(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "swiftc", "swiftc"},
		{"single file", "swiftc main.swift", "swiftc <path>"},
		{"multiple files", "swiftc main.swift utils.swift", "swiftc <path>+"},

		// Output flag
		{"output", "swiftc -o myapp main.swift", "swiftc -o <path>+"},
		{"output long path", "swiftc -o /usr/local/bin/myapp main.swift", "swiftc -o <path>+"},

		// Optimization levels (boolean)
		{"optimize", "swiftc -O main.swift", "swiftc -O <path>"},
		{"optimize size", "swiftc -Osize main.swift", "swiftc -Osize <path>"},
		{"optimize unchecked", "swiftc -Ounchecked main.swift", "swiftc -Ounchecked <path>"},

		// Target and module flags
		{"target", "swiftc -target arm64-apple-macosx14.0 main.swift", "swiftc -target <val> <path>"},
		{"module-name", "swiftc -module-name MyModule main.swift", "swiftc -module-name <val> <path>"},
		{"module-cache-path", "swiftc -module-cache-path /tmp/cache main.swift", "swiftc -module-cache-path <path>+"},

		// Include/link paths
		{"include path", "swiftc -I /usr/local/include main.swift", "swiftc -I <path>+"},
		{"lib path", "swiftc -L /usr/local/lib main.swift", "swiftc -L <path>+"},
		{"link lib", "swiftc -l sqlite3 main.swift", "swiftc -l <val> <path>"},
		{"framework", "swiftc -framework Foundation main.swift", "swiftc -framework <val> <path>"},
		{"framework search", "swiftc -F /Library/Frameworks main.swift", "swiftc -F <path>+"},

		// SDK
		{"sdk", "swiftc -sdk /Applications/Xcode.app/Contents/Developer/Platforms/MacOSX.platform/Developer/SDKs/MacOSX.sdk main.swift",
			"swiftc -sdk <path>+"},

		// Emit modes (boolean)
		{"emit-object", "swiftc -emit-object main.swift", "swiftc -emit-object <path>"},
		{"emit-library", "swiftc -emit-library main.swift", "swiftc -emit-library <path>"},
		{"emit-executable", "swiftc -emit-executable main.swift", "swiftc -emit-executable <path>"},

		// Pass-through flags
		{"Xcc", "swiftc -Xcc -DDEBUG main.swift", "swiftc -Xcc <val> <path>"},
		{"Xlinker", "swiftc -Xlinker -rpath -Xlinker /usr/lib main.swift", "swiftc -Xlinker <val> -Xlinker <val> <path>"},

		// Conditional compilation
		{"define", "swiftc -D DEBUG main.swift", "swiftc -D <val> <path>"},

		// Swift version
		{"swift-version", "swiftc -swift-version 5 main.swift", "swiftc -swift-version <val> <path>"},

		// Whole module
		{"whole-module", "swiftc -whole-module-optimization main.swift utils.swift", "swiftc -whole-module-optimization <path>+"},

		// Num threads
		{"num-threads", "swiftc -num-threads 4 main.swift", "swiftc -num-threads N <path>"},

		// Working directory
		{"working-directory", "swiftc -working-directory /tmp/build main.swift", "swiftc -working-directory <path>+"},

		// Fused flag with =
		{"cxx-interop fused", "swiftc -cxx-interoperability-mode=default main.swift", "swiftc -cxx-interoperability-mode=<val> <path>"},
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
	t.Run("different source files collide", func(t *testing.T) {
		a := shellshape.Normalize("swiftc foo.swift")
		b := shellshape.Normalize("swiftc bar.swift")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different output paths collide", func(t *testing.T) {
		a := shellshape.Normalize("swiftc -o /tmp/app1 main.swift")
		b := shellshape.Normalize("swiftc -o /tmp/app2 main.swift")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("swiftc main.swift")
		subshell := shellshape.Normalize("swiftc $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
