package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestJava(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"class only", "java HelloWorld", "java HelloWorld"},
		{"class with package", "java com.example.Main", "java com.example.Main"},
		{"class with args", "java HelloWorld arg1 arg2", "java HelloWorld <arg>+"},
		{"version", "java -version", "java -version"},
		{"no args", "java", "java"},

		// Jar mode
		{"jar", "java -jar app.jar", "java -jar <path>"},
		{"jar with args", "java -jar app.jar --port 3000 foo", "java -jar <path> <arg>+"},
		{"jar with classpath", "java -cp lib/* -jar myapp.jar", "java -cp <path> -jar <path>"},

		// Classpath flags
		{"classpath short", "java -cp /usr/lib/java Main", "java -cp <path> Main"},
		{"classpath long", "java -classpath /opt/lib Main", "java -classpath <path> Main"},
		{"class-path double dash", "java --class-path /opt/lib Main", "java --class-path <path> Main"},

		// Module path
		{"module-path", "java --module-path mods -m com.app/Main", "java --module-path <path> -m <val>"},
		{"module-path short", "java -p mods -m com.app/Main", "java -p <path> -m <val>"},

		// System properties (fused -D)
		{"sysprop fused", "java -Dserver.port=8080 Main", "java -D<val> Main"},
		{"sysprop bare", "java -Dfoo Main", "java -D<val> Main"},
		{"multiple sysprops", "java -Dfoo=bar -Dbaz=qux Main", "java -D<val> -D<val> Main"},

		// Memory flags (fused -X)
		{"xms", "java -Xms256m Main", "java -X<val> Main"},
		{"xmx", "java -Xmx2g Main", "java -X<val> Main"},
		{"xss", "java -Xss1m -jar app.jar", "java -X<val> -jar <path>"},

		// Agent flags
		{"javaagent", "java -javaagent:/path/to/agent.jar Main", "java -javaagent:<val> Main"},
		{"agentlib", "java -agentlib:jdwp=transport=dt_socket,server=y Main", "java -agentlib:<val> Main"},
		{"agentpath", "java -agentpath:/path/to/agent.so Main", "java -agentpath:<val> Main"},

		// Add-* module flags
		{"add-modules", "java --add-modules java.sql Main", "java --add-modules <val> Main"},
		{"add-opens", "java --add-opens java.base/java.lang=ALL-UNNAMED Main", "java --add-opens <val> Main"},
		{"add-exports", "java --add-exports java.base/sun.nio.ch=ALL-UNNAMED Main", "java --add-exports <val> Main"},
		{"add-reads", "java --add-reads java.base=ALL-UNNAMED Main", "java --add-reads <val> Main"},

		// Boolean flags
		{"verbose", "java -verbose Main", "java -verbose Main"},
		{"enable-assertions", "java -ea Main", "java -ea Main"},
		{"server", "java -server Main", "java -server Main"},
		{"showversion", "java -showversion -jar app.jar", "java -showversion -jar <path>"},

		// Complex invocations
		{"complex", "java -Xmx4g -cp lib/a.jar:lib/b.jar -Dspring.profiles.active=prod com.example.App --server.port=8080",
			"java -X<val> -cp <path> -D<val> com.example.App <arg>"},

		// Redirect
		{"redirect", "java -jar app.jar > output.log 2>&1", "java -jar <path> > <path> 2>&1"},
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
	t.Run("different jars collide", func(t *testing.T) {
		a := shellshape.Normalize("java -jar server.jar")
		b := shellshape.Normalize("java -jar worker.jar")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different args collide", func(t *testing.T) {
		a := shellshape.Normalize("java -jar app.jar --port 3000")
		b := shellshape.Normalize("java -jar app.jar --host localhost")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("java -jar app.jar")
		subshell := shellshape.Normalize("java -jar $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestJavac(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "javac HelloWorld.java", "javac <path>"},
		{"multiple files", "javac Main.java Utils.java", "javac <path>+"},
		{"version", "javac -version", "javac -version"},
		{"no args", "javac", "javac"},

		// Destination flags
		{"output dir", "javac -d build HelloWorld.java", "javac -d <path>+"},
		{"source dir", "javac -sourcepath src Main.java", "javac -sourcepath <path>+"},
		{"source-path long", "javac --source-path src Main.java", "javac --source-path <path>+"},

		// Classpath flags
		{"classpath short", "javac -cp lib/a.jar Main.java", "javac -cp <path>+"},
		{"classpath long", "javac -classpath lib/a.jar:lib/b.jar Main.java", "javac -classpath <path>+"},
		{"class-path double dash", "javac --class-path lib Main.java", "javac --class-path <path>+"},

		// Source/target version
		{"source version", "javac -source 11 Main.java", "javac -source <val> <path>"},
		{"target version", "javac -target 11 Main.java", "javac -target <val> <path>"},
		{"release", "javac --release 17 Main.java", "javac --release <val> <path>"},
		{"source long", "javac --source 11 Main.java", "javac --source <val> <path>"},
		{"target long", "javac --target 11 Main.java", "javac --target <val> <path>"},

		// Encoding
		{"encoding", "javac -encoding UTF-8 Main.java", "javac -encoding <val> <path>"},

		// Processor flags
		{"processor", "javac -processor com.example.MyProcessor Main.java", "javac -processor <val> <path>"},
		{"processorpath", "javac -processorpath /opt/processors Main.java", "javac -processorpath <path>+"},
		{"processor-path long", "javac --processor-path /opt/proc Main.java", "javac --processor-path <path>+"},
		{"processor-module-path", "javac --processor-module-path mods Main.java", "javac --processor-module-path <path>+"},

		// Module flags
		{"module-path", "javac --module-path mods Main.java", "javac --module-path <path>+"},
		{"module-path short", "javac -p mods Main.java", "javac -p <path>+"},
		{"module", "javac --module com.app Main.java", "javac --module <val> <path>"},
		{"module short", "javac -m com.app Main.java", "javac -m <val> <path>"},
		{"add-modules", "javac --add-modules java.sql Main.java", "javac --add-modules <val> <path>"},

		// Boolean flags
		{"debug info", "javac -g HelloWorld.java", "javac -g <path>"},
		{"verbose", "javac -verbose Main.java", "javac -verbose <path>"},
		{"deprecation", "javac -deprecation Main.java", "javac -deprecation <path>"},
		{"werror", "javac -Werror Main.java", "javac -Werror <path>"},
		{"nowarn", "javac -nowarn Main.java", "javac -nowarn <path>"},

		// Header output dir
		{"header dir", "javac -h gen Main.java", "javac -h <path>+"},
		{"source output dir", "javac -s gen Main.java", "javac -s <path>+"},

		// Complex
		{"complex", "javac -cp lib/a.jar -d build -source 11 -target 11 src/Main.java src/Utils.java",
			"javac -cp <path> -d <path> -source <val> -target <val> <path>+"},

		// Redirect
		{"redirect", "javac Main.java 2> errors.log", "javac <path> 2> <path>"},
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
		a := shellshape.Normalize("javac Main.java")
		b := shellshape.Normalize("javac App.java")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different classpaths collide", func(t *testing.T) {
		a := shellshape.Normalize("javac -cp lib/a.jar Main.java")
		b := shellshape.Normalize("javac -cp lib/b.jar App.java")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("javac Main.java")
		subshell := shellshape.Normalize("javac $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
