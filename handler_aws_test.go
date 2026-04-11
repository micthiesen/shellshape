package shellshape

import "testing"

func TestAws(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage - second subcommand kept verbatim
		{"s3 cp local to remote", "aws s3 cp file.txt s3://my-bucket/key", "aws s3 cp <path> <s3-uri>"},
		{"s3 ls", "aws s3 ls s3://my-bucket/prefix/", "aws s3 ls <s3-uri>"},
		{"s3 sync", "aws s3 sync ./build s3://deploy-bucket --delete", "aws s3 sync <path> <s3-uri> --delete"},
		{"s3 rm recursive", "aws s3 rm s3://my-bucket/old/ --recursive", "aws s3 rm <s3-uri> --recursive"},
		{"s3 mb", "aws s3 mb s3://new-bucket", "aws s3 mb <s3-uri>"},

		// Lambda
		{"lambda invoke", "aws lambda invoke --function-name myFunc output.json", "aws lambda invoke --function-name <val> <path>"},
		{"lambda list-functions", "aws lambda list-functions --region eu-west-1", "aws lambda list-functions --region eu-west-1"},

		// ECS
		{"ecs describe-tasks", "aws ecs describe-tasks --cluster prod --tasks abc123", "aws ecs describe-tasks --cluster <val> --tasks <val>"},

		// EC2 with query
		{"ec2 with query", "aws ec2 describe-instances --query Reservations[].Instances[].InstanceId --output table", "aws ec2 describe-instances --query <query> --output table"},

		// STS no extra args
		{"sts get-caller-identity", "aws sts get-caller-identity", "aws sts get-caller-identity"},
		{"sts with output", "aws sts get-caller-identity --output json", "aws sts get-caller-identity --output json"},

		// Global flags
		{"profile flag", "aws s3 ls --profile staging", "aws s3 ls --profile <val>"},
		{"region flag", "aws s3 ls --region us-west-2", "aws s3 ls --region us-west-2"},
		{"debug flag", "aws s3 ls --debug", "aws s3 ls --debug"},
		{"endpoint-url", "aws s3 ls --endpoint-url http://localhost:4566", "aws s3 ls --endpoint-url <val>"},

		// Numeric timeout flags
		{"cli-read-timeout", "aws s3 cp bigfile.tar s3://bucket/key --cli-read-timeout 300", "aws s3 cp <path> <s3-uri> --cli-read-timeout N"},

		// No subcommand after service (e.g. just service help)
		{"no second subcommand", "aws configure", "aws configure"},

		// Multiple positionals
		{"s3 cp with flags interspersed", "aws s3 cp s3://src-bucket/key s3://dst-bucket/key --region us-east-1", "aws s3 cp <s3-uri> <s3-uri> --region us-east-1"},

		// --bucket flag
		{"s3api with bucket", "aws s3api list-objects --bucket my-bucket --prefix logs/", "aws s3api list-objects --bucket <val> --prefix <val>"},

		// --ca-bundle
		{"ca-bundle", "aws s3 ls --ca-bundle /path/to/cert.pem", "aws s3 ls --ca-bundle <path>"},
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
	t.Run("different function names collide", func(t *testing.T) {
		a := Normalize("aws lambda invoke --function-name myFunc out.json")
		b := Normalize("aws lambda invoke --function-name otherFunc out.json")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different queries collide", func(t *testing.T) {
		a := Normalize("aws ec2 describe-instances --query Reservations[].Instances[].InstanceId")
		b := Normalize("aws ec2 describe-instances --query SecurityGroups[].GroupId")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different s3 paths collide", func(t *testing.T) {
		a := Normalize("aws s3 cp data.csv s3://bucket-a/uploads/data.csv")
		b := Normalize("aws s3 cp report.pdf s3://bucket-b/reports/report.pdf")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("aws lambda invoke --function-name myFunc out.json")
		subshell := Normalize("aws lambda invoke --function-name $(get-func-name) out.json")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
