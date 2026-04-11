package shellshape

import "testing"

func TestHelm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// install subcommand
		{"install basic", "helm install my-release bitnami/nginx", "helm install my-release <path>"},
		{"install with values file", "helm install my-release bitnami/nginx -f values.yaml", "helm install my-release <path> -f <path>"},
		{"install with long values", "helm install my-release bitnami/nginx --values values.yaml", "helm install my-release <path> --values <path>"},
		{"install with set", "helm install my-release bitnami/nginx --set image.tag=latest", "helm install my-release <path> --set <val>"},
		{"install with set-string", "helm install my-release bitnami/nginx --set-string nodeSelector.role=worker", "helm install my-release <path> --set-string <val>"},
		{"install with set-file", "helm install my-release bitnami/nginx --set-file config=conf.yaml", "helm install my-release <path> --set-file <val>"},
		{"install with namespace", "helm install my-release bitnami/nginx -n production", "helm install my-release <path> -n <val>"},
		{"install with long namespace", "helm install my-release bitnami/nginx --namespace production", "helm install my-release <path> --namespace <val>"},
		{"install with version", "helm install my-release bitnami/nginx --version 13.2.1", "helm install my-release <path> --version <val>"},
		{"install with wait", "helm install my-release bitnami/nginx --wait", "helm install my-release <path> --wait"},
		{"install with atomic", "helm install my-release bitnami/nginx --atomic", "helm install my-release <path> --atomic"},
		{"install with create-namespace", "helm install my-release bitnami/nginx --create-namespace -n staging", "helm install my-release <path> --create-namespace -n <val>"},
		{"install with timeout", "helm install my-release bitnami/nginx --timeout 10m", "helm install my-release <path> --timeout <val>"},
		{"install with kubeconfig", "helm install my-release bitnami/nginx --kubeconfig /home/user/.kube/config", "helm install my-release <path> --kubeconfig <path>"},
		{"install dry-run debug", "helm install my-release bitnami/nginx --dry-run --debug", "helm install my-release <path> --dry-run --debug"},
		{"install full flags", "helm install my-release bitnami/nginx -f values.yaml --set foo=bar -n prod --version 1.0.0 --wait --atomic", "helm install my-release <path> -f <path> --set <val> -n <val> --version <val> --wait --atomic"},
		{"install local chart", "helm install my-release ./my-chart", "helm install my-release <path>"},

		// upgrade subcommand
		{"upgrade basic", "helm upgrade my-release bitnami/nginx", "helm upgrade my-release <path>"},
		{"upgrade with install", "helm upgrade --install my-release bitnami/redis --version 17.0.0 --wait", "helm upgrade --install my-release <path> --version <val> --wait"},
		{"upgrade with set and namespace", "helm upgrade my-release ./chart --set foo=bar --set baz=qux -n prod", "helm upgrade my-release <path> --set <val> --set <val> -n <val>"},
		{"upgrade with force", "helm upgrade my-release bitnami/nginx --force", "helm upgrade my-release <path> --force"},

		// uninstall subcommand
		{"uninstall basic", "helm uninstall my-release", "helm uninstall my-release"},
		{"uninstall with namespace", "helm uninstall my-release -n staging", "helm uninstall my-release -n <val>"},

		// list subcommand
		{"list basic", "helm list", "helm list"},
		{"list with namespace", "helm list -n kube-system", "helm list -n <val>"},
		{"list all namespaces", "helm list -A", "helm list -A"},

		// repo subcommand
		{"repo add", "helm repo add bitnami https://charts.bitnami.com/bitnami", "helm repo add bitnami <https-uri>"},
		{"repo remove", "helm repo remove bitnami", "helm repo remove bitnami"},
		{"repo update", "helm repo update", "helm repo update"},

		// other subcommands
		{"rollback", "helm rollback my-release 2", "helm rollback my-release N"},
		{"status", "helm status my-release", "helm status my-release"},
		{"get values", "helm get values my-release", "helm get values my-release"},

		// subshell as value-flag argument
		{"set with subshell", "helm install my-release bitnami/nginx --set $(get-value)", "helm install my-release <path> --set $(get-value)"},
		// subshell as path-flag argument
		{"values file subshell", "helm install my-release bitnami/nginx -f $(find-values)", "helm install my-release <path> -f $(find-values)"},
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
	t.Run("different set values collide", func(t *testing.T) {
		a := Normalize("helm install my-release bitnami/nginx --set image.tag=v1.0")
		b := Normalize("helm install my-release bitnami/nginx --set image.tag=v2.0")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different namespace values collide", func(t *testing.T) {
		a := Normalize("helm install my-release bitnami/nginx -n staging")
		b := Normalize("helm install my-release bitnami/nginx -n production")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different values files collide", func(t *testing.T) {
		a := Normalize("helm install my-release bitnami/nginx -f values-dev.yaml")
		b := Normalize("helm install my-release bitnami/nginx -f values-prod.yaml")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("helm install my-release bitnami/nginx")
		subshell := Normalize("helm install $(dangerous-command) bitnami/nginx")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
