package shellshape

import "testing"

func TestKubectl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage - resource type kept verbatim, resource names collapsed
		{"get pods", "kubectl get pods", "kubectl get pods"},
		{"get specific pod", "kubectl get pods my-pod", "kubectl get pods <word>"},
		{"describe pod", "kubectl describe pod my-pod", "kubectl describe pod <word>"},
		{"delete deployment", "kubectl delete deployment my-deploy", "kubectl delete deployment <word>"},
		{"logs", "kubectl logs my-pod", "kubectl logs <word>"},
		{"apply file", "kubectl apply -f manifest.yaml", "kubectl apply -f <path>"},
		{"apply recursive", "kubectl apply -f ./k8s/ --recursive", "kubectl apply -f <path> --recursive"},

		// Namespace flag
		{"get pods namespace", "kubectl get pods -n production", "kubectl get pods -n <val>"},
		{"get pods long namespace", "kubectl get pods --namespace kube-system", "kubectl get pods --namespace <val>"},

		// Output flag kept verbatim
		{"output json", "kubectl get pods -o json", "kubectl get pods -o json"},
		{"output yaml", "kubectl get pods --output yaml", "kubectl get pods --output yaml"},
		{"output wide", "kubectl get pods -o wide", "kubectl get pods -o wide"},

		// Label selector collapsed
		{"selector", "kubectl get pods -l app=web", "kubectl get pods -l <selector>"},
		{"long selector", "kubectl get pods --selector app=web,tier=frontend", "kubectl get pods --selector <selector>"},
		{"field selector", "kubectl get pods --field-selector status.phase=Running", "kubectl get pods --field-selector <selector>"},

		// Container flag
		{"container", "kubectl logs my-pod -c sidecar", "kubectl logs <word> -c <val>"},
		{"long container", "kubectl logs my-pod --container sidecar", "kubectl logs <word> --container <val>"},

		// Context and kubeconfig
		{"context", "kubectl get pods --context staging", "kubectl get pods --context <val>"},
		{"kubeconfig", "kubectl get pods --kubeconfig /home/user/.kube/config", "kubectl get pods --kubeconfig <path>"},

		// Exec with double dash
		{"exec", "kubectl exec -it my-pod -- bash", "kubectl exec -it <word> -- bash"},
		{"exec with container", "kubectl exec -it my-pod -c app -- /bin/sh", "kubectl exec -it <word> -c <val> -- /bin/sh"},

		// Boolean flags
		{"all namespaces", "kubectl get pods -A", "kubectl get pods -A"},
		{"watch", "kubectl get pods -w", "kubectl get pods -w"},
		{"all namespaces long", "kubectl get pods --all-namespaces", "kubectl get pods --all-namespaces"},

		// Filename as path
		{"filename path", "kubectl apply --filename /home/user/deploy.yaml", "kubectl apply --filename <path>"},

		// Multiple resource names
		{"delete multiple", "kubectl delete pods pod1 pod2 pod3", "kubectl delete pods <word>+"},

		// Sort-by and template
		{"sort by", "kubectl get pods --sort-by .metadata.name", "kubectl get pods --sort-by <val>"},
		{"template", "kubectl get pods -o go-template --template '{{.items}}'", "kubectl get pods -o go-template --template <val>"},
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
	t.Run("different pod names collide", func(t *testing.T) {
		a := Normalize("kubectl get pods nginx-abc123")
		b := Normalize("kubectl get pods redis-xyz789")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different namespaces collide", func(t *testing.T) {
		a := Normalize("kubectl get pods -n production")
		b := Normalize("kubectl get pods -n staging")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different selectors collide", func(t *testing.T) {
		a := Normalize("kubectl get pods -l app=web")
		b := Normalize("kubectl get pods -l app=api,tier=backend")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("kubectl get pods my-pod")
		subshell := Normalize("kubectl get pods $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
