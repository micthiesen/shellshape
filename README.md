# shellshape

Normalize shell commands into stable "shapes" for use as cache keys,
fingerprints, or classification inputs.

```
$ shellshape "git log abc1234..HEAD --oneline -10"
git log <range> --oneline N

$ shellshape "grep -rn 'TODO' src/"
grep -rn <pattern> <path>

$ shellshape "FOO=bar python3 -c 'print(1)'"
FOO=<val> python3 -c <code>
```

Same logical command, same shape, regardless of specific paths, numbers,
URLs, hashes, or quoted data.

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

## Per-command handlers

The generic tokenizer handles most commands. Handlers add knowledge
about specific commands where positional arguments have known semantics:

| Handler | Commands | What it knows |
|---|---|---|
| echo | `echo`, `printf` | All positionals collapse to `<str>` |
| grep | `grep`, `egrep`, `rg`, `ag`, `ack` | First positional is `<pattern>`, `-A`/`-B`/`-C` take numeric args |
| sed | `sed` | Script expressions collapse to `<sed-expr>` |
| find | `find` | `-name`/`-path` args are `<pattern>`, positionals are paths |
| sqlite3 | `sqlite3`, `sqlite` | First positional is `<path>` (database), second is `<sql>` |

### Subshell safety

Handlers preserve `$(...)` verbatim. `echo hello` and `echo $(rm -rf /)`
always produce different shapes. This is enforced by tests on every handler.

## Use as a library

```go
import shellshape "github.com/micthiesen/shellshape"

shape := shellshape.Normalize("git log --oneline -10")
// "git log --oneline N"

exe := shellshape.ExecutableOf(shape)
// "git"
```

## Adding a handler

Each handler is a self-contained file that registers itself via `init()`:

```go
// handler_curl.go
package shellshape

func init() {
    Register("curl", handleCurl)
}

func handleCurl(subcommand string, tokens []string) []string {
    // ...
}
```

No other files need to be modified. For commands with subcommands (like `docker run`):

```go
func init() {
    Register("docker", handleDocker, HandlerOptions{HasSubcommands: true})
}
```

Every handler must:
1. Call `splitRedirects(tokens)` first
2. Check `isSubshellToken(tok)` before collapsing any positional
3. Have a subshell safety test

## Why

- **Permission/policy systems**: write rules against shapes instead of
  exact command strings
- **Command audit logging**: collapse thousands of unique commands into
  a manageable set of shapes
- **Shell history deduplication**: cluster semantically identical commands
- **Security classification**: smaller input space, better accuracy
- **Rate limiting**: count executions by shape, not by exact string

## License

MIT
