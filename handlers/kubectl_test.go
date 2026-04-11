package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

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

		// Subshell as resource type position (non-name-first subcommand)
		{"subshell as resource type", "kubectl get $(get_resource_type)", "kubectl get $(get_resource_type)"},

		// Subshell as value for verbatim value flag (-o)
		{"subshell as output value", "kubectl get pods -o $(pick_format)", "kubectl get pods -o $(pick_format)"},

		// Subshell as value for selector flag
		{"subshell as selector value", "kubectl get pods -l $(build_selector)", "kubectl get pods -l $(build_selector)"},

		// Subshell as value for val flag (-n)
		{"subshell as namespace value", "kubectl get pods -n $(get_ns)", "kubectl get pods -n $(get_ns)"},

		// Subshell as value for path flag (-f)
		{"subshell as filename value", "kubectl apply -f $(find_manifest)", "kubectl apply -f $(find_manifest)"},

		// Unknown/unrecognized flags kept verbatim
		{"unknown flag", "kubectl get pods --show-managed-fields", "kubectl get pods --show-managed-fields"},
		{"unknown short flag", "kubectl get pods -R", "kubectl get pods -R"},

		// Resource name that classifyToken recognizes as path or URL
		{"resource name is path", "kubectl delete pods /tmp/some-pod-file", "kubectl delete pods <path>"},
		{"resource name is url", "kubectl delete pods https://example.com/pod", "kubectl delete pods <https-uri>"},

		// Name-first subcommands
		{"port-forward", "kubectl port-forward my-pod 8080:80", "kubectl port-forward <word>+"},
		{"attach", "kubectl attach my-pod", "kubectl attach <word>"},
		{"cp", "kubectl cp my-pod:/tmp/file /local/path", "kubectl cp <path>+"},
		{"debug", "kubectl debug my-pod --image busybox", "kubectl debug <word> --image <val>"},
		{"run", "kubectl run my-pod --image nginx", "kubectl run <word> --image <val>"},

		// Additional boolean flags
		{"watch-only", "kubectl get pods --watch-only", "kubectl get pods --watch-only"},
		{"dry-run", "kubectl apply -f manifest.yaml --dry-run", "kubectl apply -f <path> --dry-run"},
		{"no-headers", "kubectl get pods --no-headers", "kubectl get pods --no-headers"},
		{"show-labels", "kubectl get pods --show-labels", "kubectl get pods --show-labels"},
		{"force delete", "kubectl delete pods my-pod --force", "kubectl delete pods <word> --force"},
		{"all flag", "kubectl delete pods --all", "kubectl delete pods --all"},
		{"privileged", "kubectl exec -it my-pod --privileged -- bash", "kubectl exec -it <word> --privileged -- bash"},
		{"cascade", "kubectl delete deployment my-deploy --cascade", "kubectl delete deployment <word> --cascade"},
		{"overwrite", "kubectl label pods my-pod env=prod --overwrite", "kubectl label pods <word>+ --overwrite"},
		{"prune", "kubectl apply -f dir/ --prune", "kubectl apply -f <path> --prune"},
		{"record", "kubectl apply -f manifest.yaml --record", "kubectl apply -f <path> --record"},
		{"verbose", "kubectl get pods --verbose", "kubectl get pods --verbose"},
		{"-v flag", "kubectl get pods -v", "kubectl get pods -v"},
		{"version", "kubectl version --version", "kubectl version --version"},
		{"help", "kubectl get --help", "kubectl get --help"},
		{"-d flag", "kubectl get pods -d", "kubectl get pods -d"},
		{"stdin flag", "kubectl run my-pod --stdin --image nginx", "kubectl run <word> --stdin --image <val>"},
		{"tty flag", "kubectl exec my-pod --tty -- bash", "kubectl exec <word> --tty -- bash"},
		{"-ti combined", "kubectl exec -ti my-pod -- bash", "kubectl exec -ti <word> -- bash"},

		// Redirect
		{"redirect output", "kubectl get pods > /tmp/pods.txt", "kubectl get pods > <path>"},

		// Edge: flag at end with no value
		{"flag at end no value", "kubectl get pods -n", "kubectl get pods -n"},

		// Additional path flags
		{"certificate-authority", "kubectl get pods --certificate-authority /etc/ca.crt", "kubectl get pods --certificate-authority <path>"},
		{"client-certificate", "kubectl get pods --client-certificate /etc/cert.pem", "kubectl get pods --client-certificate <path>"},
		{"client-key", "kubectl get pods --client-key /etc/key.pem", "kubectl get pods --client-key <path>"},
		{"cache-dir", "kubectl get pods --cache-dir /tmp/cache", "kubectl get pods --cache-dir <path>"},

		// Additional val flags
		{"type flag", "kubectl patch deployment my-deploy --type merge", "kubectl patch deployment <word> --type <val>"},
		{"replicas flag", "kubectl scale deployment my-deploy --replicas 3", "kubectl scale deployment <word> --replicas <val>"},
		{"port flag", "kubectl expose deployment my-deploy --port 80", "kubectl expose deployment <word> --port <val>"},
		{"timeout flag", "kubectl delete pods my-pod --timeout 30s", "kubectl delete pods <word> --timeout <val>"},
		{"grace-period flag", "kubectl delete pods my-pod --grace-period 0", "kubectl delete pods <word> --grace-period <val>"},
		{"service-account", "kubectl run my-pod --service-account admin --image nginx", "kubectl run <word> --service-account <val> --image <val>"},
		{"cluster flag", "kubectl config use-context --cluster my-cluster", "kubectl config use-context --cluster <val>"},
		{"user flag", "kubectl config set-credentials --user admin", "kubectl config set-credentials --user <val>"},
		{"server flag", "kubectl config set-cluster --server https://k8s.example.com", "kubectl config set-cluster --server <val>"},
		{"go-template flag", "kubectl get pods --go-template '{{.items}}'", "kubectl get pods --go-template <val>"},

		// Edit subcommand (resource type based)
		{"edit deployment", "kubectl edit deployment my-deploy", "kubectl edit deployment <word>"},
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
	t.Run("different pod names collide", func(t *testing.T) {
		a := shellshape.Normalize("kubectl get pods nginx-abc123")
		b := shellshape.Normalize("kubectl get pods redis-xyz789")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different namespaces collide", func(t *testing.T) {
		a := shellshape.Normalize("kubectl get pods -n production")
		b := shellshape.Normalize("kubectl get pods -n staging")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different selectors collide", func(t *testing.T) {
		a := shellshape.Normalize("kubectl get pods -l app=web")
		b := shellshape.Normalize("kubectl get pods -l app=api,tier=backend")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("kubectl get pods my-pod")
		subshell := shellshape.Normalize("kubectl get pods $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
