# shellshape [![CI](https://github.com/micthiesen/shellshape/actions/workflows/ci.yml/badge.svg)](https://github.com/micthiesen/shellshape/actions/workflows/ci.yml)

Normalize shell commands into stable "shapes" for use as cache keys,
fingerprints, or classification inputs.

```
$ shellshape "curl -s -H 'Auth: Bearer tok123' https://api.example.com/users | jq '.data[]'"
curl -s -H <header> <https-uri> | jq <filter>

$ shellshape "git log abc1234..HEAD --oneline -10"
git log <range> --oneline N

$ shellshape "docker run --rm -v /home/app:/app -p 8080:80 nginx:latest"
docker run --rm -v <val> -p <val> nginx:latest
```

Same logical command, same shape, regardless of specific paths, numbers,
URLs, hashes, or quoted data. Over 240 commands have specialized handlers
that understand their flag grammar (what takes an argument, what's boolean,
what's structural vs data).

## Install

```bash
go install github.com/micthiesen/shellshape/cmd/shellshape@latest
```

Or build from source:

```bash
git clone https://github.com/micthiesen/shellshape.git
cd shellshape
go build -o /usr/local/bin/shellshape ./cmd/shellshape
```

## Usage

```bash
# Pass a command as an argument
shellshape "rm -rf /tmp/foo/bar"

# Or pipe commands in (one per line)
echo "git log --oneline -10" | shellshape

# Batch normalize
shellshape < commands.txt
```

## What it normalizes

| Input | Shape | Rule |
|---|---|---|
| `/tmp/foo/bar.txt` | `<path>` | Absolute, relative, home, or extension-based paths |
| `-10`, `42`, `3.14` | `N` | Numbers |
| `https://api.foo.com/x` | `<https-uri>` | URLs (scheme preserved) |
| `s3://bucket/prefix/` | `<s3-uri>` | Any URI scheme |
| `git@github.com:o/r.git` | `<git-uri>` | Git SSH remotes |
| `abc1234def` (7-40 hex) | `<hash>` | Git SHAs |
| `550e8400-e29b-...` | `<uuid>` | UUIDs |
| `abc1234..HEAD` | `<range>` | Git rev ranges |
| `HEAD~5`, `HEAD^^` | `<rev>` | HEAD variants |
| `FOO=bar cmd` | `FOO=<val> cmd` | Env var values |
| `--flag=value` | `--flag=<val>` | Long flag values |
| `--flag`, `-rf` | kept | Flags are structure |
| `<path> <path> <path>` | `<path>+` | Repeated placeholders merge |
| `<<'EOF'...EOF` | `<heredoc>` | Heredoc bodies |
| `tests.test_foo.Bar` | `<dotted-id>` | Dotted module paths |
| `# comment` | stripped | Comments are noise |
| `\`+newline | joined | Line continuations |

Per-command handlers add domain-specific placeholders beyond the generic
rules above. For example, `curl -H` collapses to `<header>`, `grep`'s
first positional becomes `<pattern>`, and `jq`'s filter becomes `<filter>`.

### Compound commands

Shell compound commands (`for`, `while`, `until`, `if`) are recognized as
single structural units. Internal separators (`;`, `&&`, `|`) inside the
compound don't split it apart, and the body is recursively normalized:

```
$ shellshape "for cmd in doctor test show stats; do python3 app \$cmd --help; done"
for cmd in <val>+ ; do python3 app <arg>+ ; done

$ shellshape "if test -f /tmp/config.yml; then cat /tmp/config.yml; else echo missing; fi"
if test -f <path> ; then cat <path> ; else echo <str> ; fi
```

For-in iteration values are collapsed to `<val>+`. C-style `for ((...))` loop
expressions are collapsed to `((<expr>))`. Nested compounds are handled
correctly.

### Subshell safety

Subshell expressions like `$(...)` are recursively normalized but never
collapsed to a data placeholder. `echo hello` and `echo $(rm -rf /)`
always produce different shapes. This is enforced by tests on every handler.

## Use as a library

```go
import (
    "github.com/micthiesen/shellshape"
    _ "github.com/micthiesen/shellshape/handlers" // register all handlers
)

shape := shellshape.Normalize("curl -s https://api.example.com | jq '.data'")
// "curl -s <https-uri> | jq <filter>"

exe := shellshape.ExecutableOf(shape)
// "curl"
```

The blank import of `handlers` is required to register all per-command
handlers via their `init()` functions. Without it, only the generic token
classifier runs.

## Adding a handler

Each handler is a self-contained file in the `handlers/` package that registers
itself via `init()`:

```go
// handlers/curl.go
package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
    shellshape.Register("curl", handleCurl)
}

func handleCurl(subcommand string, tokens []string) []string {
    // ...
}
```

No other files need to be modified. For commands with subcommands (like `docker run`):

```go
func init() {
    shellshape.Register("docker", handleDocker, shellshape.HandlerOptions{HasSubcommands: true})
}
```

Every handler must:
1. Call `shellshape.SplitRedirects(tokens)` first
2. Check `shellshape.IsSubshellToken(tok)` before collapsing any positional
3. Have a subshell safety test

PRs for new handlers are welcome. Handlers target the latest version of each
command on Unix-like systems.

## Why

- **Permission/policy systems**: write rules against shapes instead of
  exact command strings
- **Command audit logging**: collapse thousands of unique commands into
  a manageable set of shapes
- **Shell history deduplication**: cluster semantically identical commands
- **Security classification**: smaller input space, better accuracy
- **Rate limiting**: count executions by shape, not by exact string

## Alternatives

Using an LLM. But it's slower and not deterministic. So why not just get an LLM
to write the rules carefully once instead?

## License

MIT
