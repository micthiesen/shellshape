package shellshape

import "testing"

func TestGradle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage - tasks are structural, kept verbatim
		{"single task", "gradle build", "gradle build"},
		{"multiple tasks", "gradle clean build", "gradle clean build"},
		{"task list", "gradle tasks --all", "gradle tasks --all"},
		{"init with type", "gradle init --type java-library", "gradle init --type <val>"},

		// Exclude task
		{"exclude task short", "gradle build -x test", "gradle build -x <task>"},
		{"exclude task long", "gradle build --exclude-task test", "gradle build --exclude-task <task>"},

		// Path flags
		{"build file short", "gradle -b custom.gradle build", "gradle -b <path> build"},
		{"build file long", "gradle --build-file /home/user/build.gradle build", "gradle --build-file <path> build"},
		{"project dir short", "gradle -p /path/to/project build", "gradle -p <path> build"},
		{"settings file", "gradle --settings-file settings.gradle test", "gradle --settings-file <path> test"},
		{"init script", "gradle -I /path/init.gradle build", "gradle -I <path> build"},
		{"include build", "gradle --include-build ../other build", "gradle --include-build <path> build"},
		{"gradle user home", "gradle -g /custom/gradle build", "gradle -g <path> build"},
		{"project cache dir", "gradle --project-cache-dir /tmp/cache build", "gradle --project-cache-dir <path> build"},

		// Value flags
		{"max workers", "gradle build --max-workers 4", "gradle build --max-workers N"},
		{"priority", "gradle build --priority low", "gradle build --priority <val>"},
		{"console", "gradle build --console plain", "gradle build --console <val>"},
		{"warning mode", "gradle build --warning-mode all", "gradle build --warning-mode <val>"},

		// Boolean flags kept verbatim
		{"offline", "gradle build --offline", "gradle build --offline"},
		{"parallel", "gradle build --parallel", "gradle build --parallel"},
		{"daemon", "gradle --daemon build", "gradle --daemon build"},
		{"no daemon", "gradle --no-daemon build", "gradle --no-daemon build"},
		{"refresh deps", "gradle clean build --refresh-dependencies", "gradle clean build --refresh-dependencies"},
		{"stacktrace", "gradle build --stacktrace", "gradle build --stacktrace"},
		{"debug", "gradle build -d", "gradle build -d"},
		{"quiet", "gradle build -q", "gradle build -q"},
		{"scan", "gradle build --scan", "gradle build --scan"},
		{"continuous", "gradle build --continuous", "gradle build --continuous"},
		{"dry run", "gradle build -m", "gradle build -m"},

		// System properties and project properties (fused with value)
		{"system prop", "gradle -Dorg.gradle.debug=true build", "gradle <sysprop> build"},
		{"project prop", "gradle -Pmyprop=myvalue build", "gradle <projprop> build"},

		// gradlew alias
		{"gradlew", "gradlew assembleRelease", "gradlew assembleRelease"},
		{"gradlew with flags", "gradlew build --offline -x test", "gradlew build --offline -x <task>"},

		// Fused --flag=value
		{"fused type", "gradle init --type=java-library", "gradle init --type=<val>"},
		{"fused max workers", "gradle build --max-workers=4", "gradle build --max-workers=N"},

		// Complex real-world
		{"full build", "gradle clean build -x test --parallel --max-workers 4", "gradle clean build -x <task> --parallel --max-workers N"},
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
	t.Run("different excluded tasks collide", func(t *testing.T) {
		a := Normalize("gradle build -x test")
		b := Normalize("gradle build -x lint")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different project dirs collide", func(t *testing.T) {
		a := Normalize("gradle -p /home/alice/project build")
		b := Normalize("gradle -p /home/bob/project build")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different system props collide", func(t *testing.T) {
		a := Normalize("gradle -Dfoo=bar build")
		b := Normalize("gradle -Dbaz=qux build")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("gradle build")
		subshell := Normalize("gradle $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
