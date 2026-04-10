# ShellShape

## Quick Reference

```bash
go build ./...    # Build all packages
go test ./...     # Run all tests
go vet ./...      # Static analysis
go fmt ./...      # Format code
go run .          # Run the main package
```

**Always run `go vet ./... && go test ./...` after making changes.**

## Code Style

- Standard Go formatting via `gofmt` (tabs, no config needed)
- Follow Go idioms: short variable names in tight scope, explicit error handling
- No over-engineering: simple solutions, no unnecessary abstractions
