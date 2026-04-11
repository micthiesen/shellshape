package shellshape

import "testing"

func TestMvn(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single phase", "mvn compile", "mvn compile"},
		{"clean package", "mvn clean package", "mvn clean package"},
		{"clean install", "mvn clean install", "mvn clean install"},
		{"plugin goal", "mvn exec:java", "mvn exec:java"},
		{"dependency tree", "mvn dependency:tree", "mvn dependency:tree"},

		// -D system properties
		{"skipTests boolean", "mvn package -DskipTests", "mvn package -DskipTests"},
		{"skip with value", "mvn clean package -Dmaven.test.skip=true", "mvn clean package -Dmaven.test.skip=<val>"},
		{"exec mainClass", `mvn exec:java -Dexec.mainClass="com.example.Main" -Dexec.args="arg1 arg2"`, "mvn exec:java -Dexec.mainClass=<val> -Dexec.args=<val>"},
		{"separate -D", "mvn -D skipTests package", "mvn -D skipTests package"},

		// Profile flag
		{"fused profile", "mvn clean -Pproduction package", "mvn clean -P <val> package"},
		{"separate profile", "mvn clean -P dev package", "mvn clean -P <val> package"},

		// Path flags
		{"pom file -f", "mvn -f /home/user/project/pom.xml install", "mvn -f <path> install"},
		{"settings -s", "mvn -s /path/to/settings.xml clean", "mvn -s <path> clean"},
		{"log file -l", "mvn clean install -l /tmp/build.log", "mvn clean install -l <path>"},

		// Threads
		{"threads -T", "mvn -T 4 clean install", "mvn -T N clean install"},
		{"threads long", "mvn --threads 8 package", "mvn --threads N package"},

		// Value flags
		{"resume-from", "mvn -rf :my-module install", "mvn -rf <val> install"},
		{"projects", "mvn -pl :core,:web install", "mvn -pl <val> install"},

		// Boolean flags preserved
		{"batch mode", "mvn -B clean package", "mvn -B clean package"},
		{"update snapshots", "mvn -U clean install", "mvn -U clean install"},
		{"debug", "mvn -X clean install", "mvn -X clean install"},
		{"offline", "mvn -o package", "mvn -o package"},

		// Long flags with =
		{"long fused define", "mvn --define maven.test.skip=true package", "mvn --define <val> package"},
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
	t.Run("different property values collide", func(t *testing.T) {
		a := Normalize("mvn clean package -Dapp.version=1.0.0")
		b := Normalize("mvn clean package -Dapp.version=2.5.3")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different profiles collide", func(t *testing.T) {
		a := Normalize("mvn clean package -Pdev")
		b := Normalize("mvn clean package -Pproduction")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different pom paths collide", func(t *testing.T) {
		a := Normalize("mvn -f /home/alice/project/pom.xml install")
		b := Normalize("mvn -f /home/bob/project/pom.xml install")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("mvn clean package")
		subshell := Normalize("mvn $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
