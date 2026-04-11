package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTerraform(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"init simple", "terraform init", "terraform init"},
		{"plan simple", "terraform plan", "terraform plan"},
		{"apply simple", "terraform apply", "terraform apply"},
		{"destroy simple", "terraform destroy", "terraform destroy"},
		{"output simple", "terraform output", "terraform output"},
		{"validate", "terraform validate", "terraform validate"},
		{"fmt", "terraform fmt", "terraform fmt"},

		// Boolean flags
		{"apply auto-approve", "terraform apply -auto-approve", "terraform apply -auto-approve"},
		{"plan no-color", "terraform plan -no-color", "terraform plan -no-color"},
		{"output json", "terraform output -json", "terraform output -json"},
		{"plan compact-warnings", "terraform plan -compact-warnings", "terraform plan -compact-warnings"},
		{"plan detailed-exitcode", "terraform plan -detailed-exitcode", "terraform plan -detailed-exitcode"},
		{"init reconfigure", "terraform init -reconfigure", "terraform init -reconfigure"},
		{"init upgrade", "terraform init -upgrade", "terraform init -upgrade"},
		{"init migrate-state", "terraform init -migrate-state", "terraform init -migrate-state"},

		// Flags with =value kept verbatim (isFlagToken passthrough)
		{"apply lock equals", "terraform apply -lock=false", "terraform apply -lock=false"},
		{"apply input equals", "terraform apply -input=false", "terraform apply -input=false"},
		{"plan refresh equals", "terraform plan -refresh=false", "terraform plan -refresh=false"},

		// Value flags (separate arg)
		{"plan var", "terraform plan -var 'region=us-east-1'", "terraform plan -var <val>"},
		{"plan multiple vars", "terraform plan -var 'a=1' -var 'b=2'", "terraform plan -var <val> -var <val>"},
		{"destroy target", "terraform destroy -target aws_instance.web", "terraform destroy -target <val>"},
		{"plan target multiple", "terraform plan -target aws_instance.web -target aws_s3_bucket.data", "terraform plan -target <val> -target <val>"},
		{"init backend-config", "terraform init -backend-config 'bucket=my-bucket'", "terraform init -backend-config <val>"},
		{"init backend-config kv", "terraform init -backend-config key=value", "terraform init -backend-config <val>"},

		// Value flags (=value syntax)
		{"var with equals", "terraform plan -var='region=us-east-1'", "terraform plan -var=<val>"},
		{"target with equals", "terraform destroy -target=aws_instance.web", "terraform destroy -target=<val>"},

		// Path flags (separate arg)
		{"plan var-file", "terraform plan -var-file prod.tfvars", "terraform plan -var-file <path>"},
		{"apply state", "terraform apply -state terraform.tfstate", "terraform apply -state <path>"},
		{"plan out", "terraform plan -out plan.out", "terraform plan -out <path>"},

		// Path flags (=value syntax)
		{"var-file with equals", "terraform plan -var-file=prod.tfvars", "terraform plan -var-file=<path>"},

		// Numeric flags
		{"apply parallelism", "terraform apply -parallelism 10", "terraform apply -parallelism N"},

		// Positionals (classifyToken)
		{"apply plan file", "terraform apply ./plan.out", "terraform apply <path>"},
		{"apply plan file auto-approve", "terraform apply -auto-approve ./plan.out", "terraform apply -auto-approve <path>"},
		{"import resource", "terraform import aws_instance.web i-1234567890abcdef", "terraform import <val>+"},
		{"state show", "terraform state show aws_instance.web", "terraform state show <dotted-id>"},
		{"state rm", "terraform state rm aws_instance.web", "terraform state rm <dotted-id>"},

		// Complex combinations
		{"plan full", "terraform plan -var-file prod.tfvars -out plan.out -parallelism 5 -no-color", "terraform plan -var-file <path> -out <path> -parallelism N -no-color"},
		{"apply full", "terraform apply -auto-approve -var 'env=prod' -var-file vars.tfvars -parallelism 10", "terraform apply -auto-approve -var <val> -var-file <path> -parallelism N"},
		{"destroy full", "terraform destroy -auto-approve -target aws_instance.web -var 'env=staging'", "terraform destroy -auto-approve -target <val> -var <val>"},
		{"init full", "terraform init -backend-config 'bucket=tf-state' -reconfigure -upgrade", "terraform init -backend-config <val> -reconfigure -upgrade"},
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
	t.Run("different var values collide", func(t *testing.T) {
		a := shellshape.Normalize("terraform plan -var 'region=us-east-1'")
		b := shellshape.Normalize("terraform plan -var 'region=eu-west-1'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different var-file paths collide", func(t *testing.T) {
		a := shellshape.Normalize("terraform plan -var-file prod.tfvars")
		b := shellshape.Normalize("terraform plan -var-file staging.tfvars")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different target resources collide", func(t *testing.T) {
		a := shellshape.Normalize("terraform destroy -target aws_instance.web")
		b := shellshape.Normalize("terraform destroy -target aws_s3_bucket.data")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different import args collide", func(t *testing.T) {
		a := shellshape.Normalize("terraform import aws_instance.web i-abc123")
		b := shellshape.Normalize("terraform import aws_s3_bucket.data my-bucket")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("terraform apply ./plan.out")
		subshell := shellshape.Normalize("terraform apply $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
