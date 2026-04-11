package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDocker(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// run subcommand
		{"run simple", "docker run nginx", "docker run nginx"},
		{"run interactive rm", "docker run -it --rm ubuntu bash", "docker run -it --rm ubuntu bash"},
		{"run detached with name", "docker run -d --name myapp nginx", "docker run -d --name <val> nginx"},
		{"run with port", "docker run -p 8080:80 nginx", "docker run -p <val> nginx"},
		{"run with volume", "docker run -v /host/path:/container/path nginx", "docker run -v <val> nginx"},
		{"run with env", "docker run -e MYSQL_ROOT_PASSWORD=secret mysql", "docker run -e <val> mysql"},
		{"run with workdir", "docker run -w /app node npm start", "docker run -w <path> node npm start"},
		{"run with network", "docker run --network host nginx", "docker run --network <val> nginx"},
		{"run full flags", "docker run -d --rm --name web -p 8080:80 -v /data:/data -e FOO=bar nginx", "docker run -d --rm --name <val> -p <val> -v <val> -e <val> nginx"},

		// build subcommand
		{"build simple", "docker build .", "docker build ."},
		{"build with tag", "docker build -t myapp:latest .", "docker build -t <val> ."},
		{"build with file", "docker build -f Dockerfile.prod .", "docker build -f <path> ."},
		{"build with build-arg", "docker build --build-arg VERSION=1.0 -t app .", "docker build --build-arg <val> -t <val> ."},
		{"build with path context", "docker build -t app ./src", "docker build -t <val> <path>"},

		// exec subcommand
		{"exec simple", "docker exec -it mycontainer bash", "docker exec -it mycontainer bash"},
		{"exec with env", "docker exec -e FOO=bar mycontainer sh", "docker exec -e <val> mycontainer sh"},
		{"exec with workdir", "docker exec -w /app mycontainer ls", "docker exec -w <path> mycontainer ls"},

		// compose subcommand
		{"compose up", "docker compose up -d", "docker compose up -d"},
		{"compose down", "docker compose down", "docker compose down"},
		{"compose logs", "docker compose logs -f", "docker compose logs -f"},
		{"compose with file", "docker compose -f docker-compose.prod.yml up -d", "docker compose -f <path> up -d"},

		// other flags
		{"run with platform", "docker run --platform linux/amd64 node", "docker run --platform <val> node"},
		{"run with entrypoint", "docker run --entrypoint /bin/sh nginx", "docker run --entrypoint <val> nginx"},
		{"run with user", "docker run -u 1000:1000 nginx", "docker run -u <val> nginx"},
		{"run with mount", "docker run --mount type=bind,source=/src,target=/app nginx", "docker run --mount <val> nginx"},
		{"run with restart", "docker run --restart always nginx", "docker run --restart <val> nginx"},
		{"run with label", "docker run -l app=web nginx", "docker run -l <val> nginx"},
		{"run with memory", "docker run -m 512m nginx", "docker run -m <val> nginx"},
		{"run long volume", "docker run --volume /a:/b nginx", "docker run --volume <val> nginx"},
		{"run long publish", "docker run --publish 3000:3000 nginx", "docker run --publish <val> nginx"},
		{"run long env", "docker run --env DB_HOST=localhost postgres", "docker run --env <val> postgres"},
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
	t.Run("different port mappings collide", func(t *testing.T) {
		a := shellshape.Normalize("docker run -p 8080:80 nginx")
		b := shellshape.Normalize("docker run -p 3000:3000 nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different env values collide", func(t *testing.T) {
		a := shellshape.Normalize("docker run -e FOO=bar nginx")
		b := shellshape.Normalize("docker run -e BAZ=qux nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different volume paths collide", func(t *testing.T) {
		a := shellshape.Normalize("docker run -v /home/user/data:/data nginx")
		b := shellshape.Normalize("docker run -v /tmp/test:/app nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("docker run literal-arg")
		subshell := shellshape.Normalize("docker run $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})

	t.Run("subshell as standalone flag value preserved", func(t *testing.T) {
		benign := shellshape.Normalize("docker run --name myapp nginx")
		subshell := shellshape.Normalize("docker run --name $(generate-name) nginx")
		if benign == subshell {
			t.Error("subshell as standalone flag value must produce different shape")
		}
	})
}

func TestPodman(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"run simple", "podman run nginx", "podman run nginx"},
		{"run detached with name", "podman run -d --name myapp nginx", "podman run -d --name <val> nginx"},
		{"run with port", "podman run -p 8080:80 nginx", "podman run -p <val> nginx"},
		{"run with volume", "podman run -v /host:/container nginx", "podman run -v <val> nginx"},
		{"run with env", "podman run -e SECRET=abc postgres", "podman run -e <val> postgres"},
		{"build with tag", "podman build -t myapp:latest .", "podman build -t <val> ."},
		{"exec interactive", "podman exec -it mycontainer sh", "podman exec -it mycontainer sh"},
		{"compose up", "podman compose up -d", "podman compose up -d"},
		{"run with workdir", "podman run -w /app node npm start", "podman run -w <path> node npm start"},
		{"run full flags", "podman run -d --rm --name web -p 3000:80 -v /data:/data -e FOO=bar nginx", "podman run -d --rm --name <val> -p <val> -v <val> -e <val> nginx"},
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
	t.Run("different port mappings collide", func(t *testing.T) {
		a := shellshape.Normalize("podman run -p 8080:80 nginx")
		b := shellshape.Normalize("podman run -p 3000:3000 nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("podman run literal-arg")
		subshell := shellshape.Normalize("podman run $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
