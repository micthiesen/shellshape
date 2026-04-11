package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDotnet(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// build subcommand
		{"build bare", "dotnet build", "dotnet build"},
		{"build with configuration", "dotnet build -c Release", "dotnet build -c <val>"},
		{"build with long configuration", "dotnet build --configuration Release", "dotnet build --configuration <val>"},
		{"build with output", "dotnet build -o ./bin/output", "dotnet build -o <path>"},
		{"build with framework", "dotnet build -f net8.0", "dotnet build -f <val>"},
		{"build with runtime", "dotnet build -r linux-x64", "dotnet build -r <val>"},
		{"build with property", "dotnet build -p:Version=1.0.0", "dotnet build -p:<val>"},
		{"build with verbosity", "dotnet build -v quiet", "dotnet build -v <val>"},
		{"build with no-restore", "dotnet build --no-restore", "dotnet build --no-restore"},
		{"build full flags", "dotnet build -c Release -o ./out -f net8.0 --no-restore", "dotnet build -c <val> -o <path> -f <val> --no-restore"},
		{"build project path", "dotnet build ./src/MyApp/MyApp.csproj", "dotnet build <path>"},

		// run subcommand
		{"run bare", "dotnet run", "dotnet run"},
		{"run with project", "dotnet run --project ./src/MyApp", "dotnet run --project <path>"},
		{"run with configuration", "dotnet run -c Release", "dotnet run -c <val>"},
		{"run with launch-profile", "dotnet run --launch-profile dev", "dotnet run --launch-profile <val>"},
		{"run with args after double dash", "dotnet run -- --port 8080", "dotnet run -- --port 8080"},
		{"run with project and args", "dotnet run --project ./src/MyApp -- arg1 arg2", "dotnet run --project <path> -- arg1 arg2"},

		// test subcommand
		{"test bare", "dotnet test", "dotnet test"},
		{"test with filter", "dotnet test --filter FullyQualifiedName~MyTest", "dotnet test --filter <pattern>"},
		{"test with settings", "dotnet test --settings ./test.runsettings", "dotnet test --settings <path>"},
		{"test with results dir", "dotnet test --results-directory ./results", "dotnet test --results-directory <path>"},
		{"test with logger", "dotnet test --logger trx", "dotnet test --logger <val>"},
		{"test with collect", "dotnet test --collect \"Code Coverage\"", "dotnet test --collect <val>"},
		{"test with configuration", "dotnet test -c Release --no-build", "dotnet test -c <val> --no-build"},
		{"test project path", "dotnet test ./tests/MyTests.csproj", "dotnet test <path>"},

		// publish subcommand
		{"publish bare", "dotnet publish", "dotnet publish"},
		{"publish with configuration and runtime", "dotnet publish -c Release -r linux-x64 --self-contained", "dotnet publish -c <val> -r <val> --self-contained"},
		{"publish with output", "dotnet publish -o ./publish", "dotnet publish -o <path>"},

		// new subcommand — template name is structural
		{"new console", "dotnet new console", "dotnet new console"},
		{"new webapi", "dotnet new webapi", "dotnet new webapi"},
		{"new with name", "dotnet new console -n MyApp", "dotnet new console -n <val>"},
		{"new with output", "dotnet new classlib -o ./libs/MyLib", "dotnet new classlib -o <path>"},

		// add/remove package — package name is structural
		{"add package", "dotnet add package Newtonsoft.Json", "dotnet add package Newtonsoft.Json"},
		{"add package with version", "dotnet add package Serilog --version 3.1.0", "dotnet add package Serilog --version <val>"},
		{"add reference", "dotnet add reference ../MyLib/MyLib.csproj", "dotnet add reference <path>"},
		{"remove package", "dotnet remove package Newtonsoft.Json", "dotnet remove package Newtonsoft.Json"},

		// tool subcommand — tool name is structural
		{"tool install global", "dotnet tool install dotnet-ef -g", "dotnet tool install dotnet-ef -g"},
		{"tool install with version", "dotnet tool install dotnet-ef --version 8.0.0", "dotnet tool install dotnet-ef --version <val>"},
		{"tool uninstall", "dotnet tool uninstall dotnet-ef -g", "dotnet tool uninstall dotnet-ef -g"},
		{"tool list", "dotnet tool list -g", "dotnet tool list -g"},

		// restore / clean
		{"restore bare", "dotnet restore", "dotnet restore"},
		{"restore with source", "dotnet restore --source https://api.nuget.org/v3/index.json", "dotnet restore --source <val>"},
		{"clean bare", "dotnet clean", "dotnet clean"},

		// nuget subcommand
		{"nuget push", "dotnet nuget push ./pkg.nupkg --source https://api.nuget.org/v3/index.json --api-key abc123", "dotnet nuget push <path> --source <val> --api-key <val>"},

		// watch subcommand
		{"watch run", "dotnet watch run", "dotnet watch run"},
		{"watch test", "dotnet watch test", "dotnet watch test"},

		// sln subcommand
		{"sln add", "dotnet sln add ./src/MyApp/MyApp.csproj", "dotnet sln add <path>"},
		{"sln list", "dotnet sln list", "dotnet sln list"},

		// format subcommand
		{"format bare", "dotnet format", "dotnet format"},
		{"format with verbosity", "dotnet format -v diag", "dotnet format -v <val>"},

		// run a dll directly — framework treats the dll as the subcommand (verbatim)
		{"run dll", "dotnet myapp.dll", "dotnet myapp.dll"},
		{"run dll with path", "dotnet ./bin/Release/myapp.dll", "dotnet ./bin/Release/myapp.dll"},

		// fused --flag=value
		{"fused configuration", "dotnet build --configuration=Release", "dotnet build --configuration=<val>"},
		{"fused output", "dotnet build --output=./out", "dotnet build --output=<path>"},

		// MSBuild property with /p: and -p:
		{"msbuild property long", "dotnet build --property:Configuration=Release", "dotnet build --property:<val>"},
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
	t.Run("different configurations collide", func(t *testing.T) {
		a := shellshape.Normalize("dotnet build -c Release")
		b := shellshape.Normalize("dotnet build -c Debug")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different output paths collide", func(t *testing.T) {
		a := shellshape.Normalize("dotnet publish -o ./publish/linux")
		b := shellshape.Normalize("dotnet publish -o ./publish/windows")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different test filters collide", func(t *testing.T) {
		a := shellshape.Normalize("dotnet test --filter FullyQualifiedName~TestA")
		b := shellshape.Normalize("dotnet test --filter FullyQualifiedName~TestB")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("dotnet literal-arg")
		subshell := shellshape.Normalize("dotnet $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
