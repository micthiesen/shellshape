package shellshape

import "testing"

func TestFly(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// deploy (no second subcommand)
		{"deploy simple", "fly deploy", "fly deploy"},
		{"deploy with app", "fly deploy --app my-app", "fly deploy --app <val>"},
		{"deploy with dockerfile", "fly deploy --dockerfile Dockerfile.prod", "fly deploy --dockerfile <path>"},
		{"deploy with strategy", "fly deploy -a my-app --strategy canary", "fly deploy -a <val> --strategy <val>"},
		{"deploy with image", "fly deploy -i docker.io/org/app:latest", "fly deploy -i <val>"},
		{"deploy detach", "fly deploy --detach --no-cache", "fly deploy --detach --no-cache"},
		{"deploy with env", "fly deploy -e FOO=bar -e BAZ=qux", "fly deploy -e <val> -e <val>"},
		{"deploy with build-arg", "fly deploy --build-arg VERSION=1.0 --build-arg ENV=prod", "fly deploy --build-arg <val> --build-arg <val>"},
		{"deploy remote only", "fly deploy --remote-only -a my-app", "fly deploy --remote-only -a <val>"},
		{"deploy with config", "fly deploy -c ./fly.toml", "fly deploy -c <path>"},
		{"deploy with region", "fly deploy --primary-region lax", "fly deploy --primary-region <val>"},

		// launch (no second subcommand)
		{"launch simple", "fly launch", "fly launch"},
		{"launch with name", "fly launch --name my-app --region lax", "fly launch --name <val> --region <val>"},
		{"launch with org", "fly launch --org personal --now", "fly launch --org <val> --now"},
		{"launch with dockerfile", "fly launch --dockerfile Dockerfile", "fly launch --dockerfile <path>"},
		{"launch yes", "fly launch -y", "fly launch -y"},
		{"launch with image", "fly launch --image nginx:latest --name web", "fly launch --image <val> --name <val>"},

		// apps (has second subcommand)
		{"apps list", "fly apps list", "fly apps list"},
		{"apps create", "fly apps create my-app --org personal", "fly apps create my-app --org <val>"},
		{"apps destroy", "fly apps destroy my-app", "fly apps destroy my-app"},

		// machine (has second subcommand)
		{"machine list", "fly machine list -a my-app", "fly machine list -a <val>"},
		{"machine stop", "fly machine stop 123abc -a my-app", "fly machine stop 123abc -a <val>"},
		{"machine destroy", "fly machine destroy 123abc --force", "fly machine destroy 123abc --force"},

		// secrets (has second subcommand)
		{"secrets set", "fly secrets set DATABASE_URL=postgres://host/db", "fly secrets set <path>"},
		{"secrets set simple", "fly secrets set MY_SECRET=abc123", "fly secrets set MY_SECRET=abc123"},
		{"secrets list", "fly secrets list -a my-app", "fly secrets list -a <val>"},
		{"secrets unset", "fly secrets unset SECRET_KEY -a my-app", "fly secrets unset SECRET_KEY -a <val>"},

		// scale (has second subcommand)
		{"scale count", "fly scale count 3 -a my-app", "fly scale count N -a <val>"},
		{"scale vm", "fly scale vm shared-cpu-1x -a my-app", "fly scale vm shared-cpu-1x -a <val>"},
		{"scale memory", "fly scale memory 512 -a my-app", "fly scale memory N -a <val>"},
		{"scale show", "fly scale show -a my-app", "fly scale show -a <val>"},

		// status (no second subcommand)
		{"status", "fly status --app my-app", "fly status --app <val>"},

		// logs (no second subcommand)
		{"logs", "fly logs --app my-app --region ord", "fly logs --app <val> --region <val>"},

		// ssh (has second subcommand)
		{"ssh console", "fly ssh console -a my-app", "fly ssh console -a <val>"},
		{"ssh sftp", "fly ssh sftp get /app/data.db ./data.db -a my-app", "fly ssh sftp get <path>+ -a <val>"},

		// volumes (has second subcommand)
		{"volumes create", "fly volumes create data --region lax --size 10", "fly volumes create data --region <val> --size <val>"},
		{"volumes list", "fly volumes list -a my-app", "fly volumes list -a <val>"},

		// ips (has second subcommand)
		{"ips list", "fly ips list -a my-app", "fly ips list -a <val>"},
		{"ips allocate-v4", "fly ips allocate-v4 --shared -a my-app", "fly ips allocate-v4 --shared -a <val>"},

		// certs (has second subcommand)
		{"certs add", "fly certs add example.com -a my-app", "fly certs add <dotted-id> -a <val>"},
		{"certs list", "fly certs list -a my-app", "fly certs list -a <val>"},

		// dashboard, open, version, doctor (no second subcommand)
		{"dashboard", "fly dashboard", "fly dashboard"},
		{"open", "fly open /path/to/page", "fly open <path>"},
		{"version", "fly version", "fly version"},

		// proxy (no second subcommand)
		{"proxy", "fly proxy 5432 -a my-db", "fly proxy N -a <val>"},
		{"proxy with bind", "fly proxy 5432:5432 -a my-db", "fly proxy 5432:5432 -a <val>"},

		// auth (has second subcommand)
		{"auth login", "fly auth login", "fly auth login"},
		{"auth token", "fly auth token", "fly auth token"},

		// config (has second subcommand)
		{"config show", "fly config show -a my-app", "fly config show -a <val>"},

		// vm-size flags
		{"deploy vm flags", "fly deploy --vm-size shared-cpu-1x --vm-memory 512 --vm-cpus 2", "fly deploy --vm-size <val> --vm-memory <val> --vm-cpus <val>"},

		// flyctl alias
		{"flyctl alias", "flyctl deploy --app my-app", "flyctl deploy --app <val>"},
		{"flyctl apps list", "flyctl apps list", "flyctl apps list"},

		// fused flags
		{"fused app flag", "fly deploy --app=my-app", "fly deploy --app=<val>"},
		{"fused region flag", "fly deploy --region=lax", "fly deploy --region=<val>"},
		{"fused dockerfile", "fly deploy --dockerfile=Dockerfile.prod", "fly deploy --dockerfile=<path>"},

		// verbose/debug (boolean)
		{"verbose flag", "fly status --verbose --app my-app", "fly status --verbose --app <val>"},
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
	t.Run("different app names collide", func(t *testing.T) {
		a := Normalize("fly deploy --app my-app-1")
		b := Normalize("fly deploy --app my-app-2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different env values collide", func(t *testing.T) {
		a := Normalize("fly deploy -e FOO=bar")
		b := Normalize("fly deploy -e BAZ=qux")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different regions collide", func(t *testing.T) {
		a := Normalize("fly deploy --region lax")
		b := Normalize("fly deploy --region ord")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("fly deploy --app my-app")
		subshell := Normalize("fly deploy --app $(get-app-name)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
