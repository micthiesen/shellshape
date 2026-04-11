package shellshape

import "testing"

func TestGo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// go test — the most common subcommand with rich flags
		{"test bare", "go test", "go test"},
		{"test current pkg", "go test .", "go test ."},
		{"test all pkgs", "go test ./...", "go test <path>"},
		{"test verbose", "go test -v ./...", "go test -v <path>"},
		{"test run pattern", "go test -run TestFoo ./...", "go test -run <pattern> <path>"},
		{"test run complex regex", "go test -run ^TestIP$ -v ./...", "go test -run <pattern> -v <path>"},
		// Note: TestSudo|TestSu with unquoted | is parsed as a pipe by the shell
		// tokenizer, so -run only gets "TestSudo". Quoted version works properly.
		{"test run pipe regex unquoted", "go test -run TestSudo|TestSu -v ./...", "go test -run <pattern> | TestSu -v <path>"},
		{"test bench", "go test -bench BenchmarkSort -benchmem ./...", "go test -bench <pattern> -benchmem <path>"},
		{"test count", "go test -count 5 ./...", "go test -count N <path>"},
		{"test count=N fused", "go test -count=1 ./...", "go test -count=N <path>"},
		{"test timeout", "go test -timeout 30s ./...", "go test -timeout <duration> <path>"},
		{"test timeout=val fused", "go test -timeout=5m ./...", "go test -timeout=<duration> <path>"},
		{"test coverprofile", "go test -coverprofile coverage.out ./...", "go test -coverprofile <path>+"},
		{"test cpuprofile", "go test -cpuprofile cpu.prof ./...", "go test -cpuprofile <path>+"},
		{"test memprofile", "go test -memprofile mem.prof ./...", "go test -memprofile <path>+"},
		{"test mutexprofile", "go test -mutexprofile mutex.prof", "go test -mutexprofile <path>"},
		{"test blockprofile", "go test -blockprofile block.prof", "go test -blockprofile <path>"},
		{"test trace", "go test -trace trace.out ./...", "go test -trace <path>+"},
		{"test outputdir", "go test -outputdir /tmp/results ./...", "go test -outputdir <path>+"},
		{"test covermode", "go test -covermode atomic ./...", "go test -covermode <val> <path>"},
		{"test coverpkg", "go test -coverpkg ./pkg/... ./...", "go test -coverpkg <val> <path>"},
		{"test benchtime", "go test -benchtime 10s ./...", "go test -benchtime <val> <path>"},
		{"test race", "go test -race ./...", "go test -race <path>"},
		{"test short", "go test -short ./...", "go test -short <path>"},
		{"test json", "go test -json ./...", "go test -json <path>"},
		{"test failfast", "go test -failfast ./...", "go test -failfast <path>"},
		{"test parallel", "go test -parallel 4 ./...", "go test -parallel N <path>"},
		{"test cpu", "go test -cpu 1,2,4 ./...", "go test -cpu <val> <path>"},
		{"test tags", "go test -tags integration ./...", "go test -tags <val> <path>"},
		{"test run specific file", "go test -run TestSomething -v ./pkg/foo", "go test -run <pattern> -v <path>"},
		{"test multiple pkgs", "go test ./pkg/a ./pkg/b", "go test <path>+"},
		{"test with ldflags", "go test -ldflags \"-X main.v=1\" ./...", "go test -ldflags <val> <path>"},

		// go build
		{"build bare", "go build", "go build"},
		{"build current", "go build .", "go build ."},
		{"build all", "go build ./...", "go build <path>"},
		{"build output", "go build -o /tmp/foo .", "go build -o <path> ."},
		{"build output path", "go build -o /tmp/foo ./cmd/app", "go build -o <path>+"},
		{"build race", "go build -race ./cmd/server", "go build -race <path>"},
		{"build ldflags", "go build -ldflags \"-s -w\" .", "go build -ldflags <val> ."},
		{"build gcflags", "go build -gcflags \"all=-N -l\" .", "go build -gcflags <val> ."},
		{"build tags", "go build -tags production ./cmd/app", "go build -tags <val> <path>"},
		{"build verbose", "go build -v ./...", "go build -v <path>"},
		{"build trimpath", "go build -trimpath -o app .", "go build -trimpath -o <path> ."},

		// go run
		{"run file", "go run main.go", "go run <path>"},
		{"run dot", "go run .", "go run ."},
		{"run pkg", "go run ./cmd/app", "go run <path>"},
		{"run with args", "go run main.go -flag value", "go run <path> -flag value"},
		{"run race", "go run -race main.go", "go run -race <path>"},

		// go install
		{"install pkg", "go install golang.org/x/tools/gopls@latest", "go install golang.org/x/tools/gopls@latest"},
		{"install local", "go install ./cmd/app", "go install <path>"},

		// go get
		{"get pkg", "go get github.com/foo/bar", "go get github.com/foo/bar"},
		{"get versioned", "go get github.com/foo/bar@v1.2.3", "go get github.com/foo/bar@v1.2.3"},
		{"get update", "go get -u ./...", "go get -u <path>"},
		{"get multiple", "go get github.com/foo/bar github.com/baz/qux", "go get github.com/foo/bar github.com/baz/qux"},

		// go mod
		{"mod tidy", "go mod tidy", "go mod tidy"},
		{"mod download", "go mod download", "go mod download"},
		{"mod vendor", "go mod vendor", "go mod vendor"},
		{"mod init", "go mod init github.com/user/repo", "go mod init github.com/user/repo"},
		{"mod why", "go mod why golang.org/x/sys", "go mod why golang.org/x/sys"},
		{"mod edit replace", "go mod edit -replace github.com/foo=../foo", "go mod edit -replace <val>"},
		{"mod edit require", "go mod edit -require github.com/foo@v1.0.0", "go mod edit -require <val>"},
		{"mod graph", "go mod graph", "go mod graph"},
		{"mod verify", "go mod verify", "go mod verify"},

		// go vet
		{"vet all", "go vet ./...", "go vet <path>"},
		{"vet current", "go vet .", "go vet ."},
		{"vet verbose", "go vet -v ./...", "go vet -v <path>"},
		{"vet with tags", "go vet -tags integration ./...", "go vet -tags <val> <path>"},

		// go fmt
		{"fmt all", "go fmt ./...", "go fmt <path>"},

		// go doc
		{"doc package", "go doc strings", "go doc strings"},
		{"doc function", "go doc strings.Replace", "go doc strings.Replace"},
		{"doc method", "go doc sql.DB.Query", "go doc sql.DB.Query"},
		{"doc all", "go doc -all strings", "go doc -all strings"},
		{"doc src", "go doc -src strings.Replace", "go doc -src strings.Replace"},

		// go env
		{"env bare", "go env", "go env"},
		{"env var", "go env GOPATH", "go env GOPATH"},
		{"env multiple vars", "go env GOPATH GOROOT", "go env GOPATH GOROOT"},
		{"env write", "go env -w GOBIN=/usr/local/bin", "go env -w GOBIN=<val>"},
		{"env json", "go env -json", "go env -json"},

		// go list
		{"list all", "go list ./...", "go list <path>"},
		{"list json", "go list -json ./...", "go list -json <path>"},
		{"list modules", "go list -m all", "go list -m all"},
		{"list format", "go list -f '{{.Dir}}' ./...", "go list -f <val> <path>"},

		// go generate
		{"generate all", "go generate ./...", "go generate <path>"},

		// go clean
		{"clean bare", "go clean", "go clean"},
		{"clean cache", "go clean -cache", "go clean -cache"},
		{"clean modcache", "go clean -modcache", "go clean -modcache"},
		{"clean testcache", "go clean -testcache", "go clean -testcache"},

		// go tool
		{"tool pprof", "go tool pprof cpu.prof", "go tool pprof <dotted-id>"},
		{"tool trace", "go tool trace trace.out", "go tool trace <dotted-id>"},
		{"tool pprof path", "go tool pprof ./cpu.prof", "go tool pprof <path>"},
		{"tool dist list", "go tool dist list", "go tool dist list"},

		// Shared build flags across subcommands
		{"build p flag", "go build -p 4 .", "go build -p N ."},
		{"test asmflags", "go test -asmflags all=-trimpath ./...", "go test -asmflags <val> <path>"},
		{"build gccgoflags", "go build -gccgoflags -O2 .", "go build -gccgoflags <val> ."},

		// Redirects
		{"test redirect", "go test ./... > output.txt", "go test <path> > <path>"},
		{"test stderr redirect", "go test ./... 2>&1", "go test <path> 2>&1"},
		{"vet stderr", "go vet ./... 2>&1", "go vet <path> 2>&1"},

		// Edge cases
		{"bare go", "go", "go"},
		{"version", "go version", "go version"},
		{"help", "go help build", "go help build"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different test names collide", func(t *testing.T) {
		a := Normalize("go test -run TestFoo -v ./...")
		b := Normalize("go test -run TestBar -v ./...")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	t.Run("different build outputs collide", func(t *testing.T) {
		a := Normalize("go build -o /tmp/app1 .")
		b := Normalize("go build -o /tmp/app2 .")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	t.Run("different timeout values collide", func(t *testing.T) {
		a := Normalize("go test -timeout 30s ./...")
		b := Normalize("go test -timeout 5m ./...")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("go test ./pkg/foo")
		subshell := Normalize("go test $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
