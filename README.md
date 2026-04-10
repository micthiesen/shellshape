# shellshape

Standalone open-source tool for normalizing shell commands into stable
"shapes" for use as cache keys, fingerprints, or classification inputs.

```
$ echo "git log abc1234..HEAD --oneline -10" | shellshape
git log <range> --oneline N
```

This document captures the design decisions, architecture plan, and
implementation notes for extracting the normalizer from claude-approvals
into a standalone Go binary called `shellshape`.

## Origin

The normalizer was built as part of claude-approvals, a local permission
gate for Claude Code's Bash tool. It turns concrete bash commands into
stable "shapes" so that `git log --oneline -10` and `git log --oneline -50`
hit the same cache entry. The normalizer is the most generally useful
piece of that system and has no dependency on the rest of the approvals
pipeline (classifier, cache, feedback loop).

The Python implementation lives in `normalize.py` in this directory, with
150+ tests in `tests/test_normalize.py`. That code is the reference
implementation and test oracle for the Go rewrite.

## Why this is useful outside approvals

- **Command audit logging**: collapse thousands of unique commands into
  a manageable set of shapes for dashboards and alerting
- **Permission/policy systems**: write rules against shapes instead of
  exact command strings
- **Shell history deduplication**: cluster semantically identical
  commands that differ only in file paths or arguments
- **Security classification**: feed shapes to a classifier instead of
  raw commands (smaller input space, better accuracy)
- **Rate limiting / abuse detection**: count executions by shape, not
  by exact string

Nothing else does this. Shell parsers (bashlex, mvdan/sh) produce ASTs
but don't reduce commands to equivalence classes. Shell history tools
(atuin, hstr) deduplicate exact matches. SIEM tools log commands
verbatim. shellshape fills the gap between "parse" and "classify."

## Name

**shellshape**. The GNOME tiling window manager project with this name
is archived and inactive. The name is not registered on any package
registry (pypi, crates.io, homebrew, npm). We take it.

## Language: Go

Decision factors:

- **Ship speed**: Go compiles fast, cross-compiles trivially
  (`GOOS=linux GOARCH=arm64 go build`), and produces static binaries
  with no runtime dependencies.
- **Contributor accessibility**: The tool is string manipulation. Go's
  `strings`, `regexp`, and `unicode` packages are all you need.
  Contributors who know bash but not systems programming can read and
  write Go. Rust's ownership model would be friction for zero benefit.
- **Toolchain**: `go test` is built in, fast, and has no setup.
  `goreleaser` handles homebrew taps and GitHub releases in one config.
- **Performance**: Irrelevant. The Python implementation runs in <1ms.
  Any compiled language is fine. Go's regexp is slower than Rust's but
  we're matching short strings, not scanning gigabytes.
- **Distribution**: `go install`, homebrew tap via goreleaser, single
  binary download from GitHub releases. No runtime, no dependencies.

What we give up vs Rust: better type system, cargo (genuinely great),
WebAssembly story if we ever wanted a browser version. None of these
matter for this tool right now.

## CLI interface

Minimal. One job, one way to use it.

```
# Pipe a command in, get a shape out
echo "git log --oneline -10" | shellshape

# Or pass as argument
shellshape "git log --oneline -10"

# Multiple commands (one per line)
shellshape < commands.txt
```

Output is one shape per input line, to stdout. No flags needed for
normal use. Possible future flags:

- `--json` for structured output (input, shape, exe, handler used)
- `--explain` to show which rule classified each token (for debugging)

No `--mode` or `--style` or `--config` for normalization behavior.
One opinionated default. See "not configurable" below.

## Core normalization rules

These are not debatable. Every user wants these:

| Input pattern              | Output            | Why                          |
|---------------------------|-------------------|------------------------------|
| `/tmp/foo/bar.txt`        | `<path>`          | Paths are data               |
| `./relative/path`         | `<path>`          | Paths are data               |
| `~/home/thing`            | `<path>`          | Paths are data               |
| `foo.py`, `settings.json` | `<path>`          | Known file extensions        |
| `-10`, `42`, `3.14`       | `N`               | Numbers are data             |
| `https://api.foo.com/x`   | `<https-uri>`     | URLs are data                |
| `s3://bucket/prefix/`     | `<s3-uri>`        | URIs are data                |
| `git@github.com:o/r.git`  | `<git-uri>`       | Git remotes are data         |
| `abc1234def` (7-40 hex)   | `<hash>`          | Git SHAs are data            |
| `550e8400-e29b-...`       | `<uuid>`          | UUIDs are data               |
| `abc1234..HEAD`           | `<range>`         | Git rev-ranges are data      |
| `HEAD~5`, `HEAD^^`        | `<rev>`           | HEAD variants are data       |
| `FOO=bar cmd`             | `FOO=<val> cmd`   | Env var values are data      |
| `--flag=value`            | `--flag=<val>`    | Long flag values are data    |
| `--flag`, `-rf`           | kept verbatim     | Flags are structure          |
| `<path> <path> <path>`    | `<path>+`         | Repeated placeholders merge  |
| `<<'EOF'\n...\nEOF`       | `<heredoc>`       | Heredoc bodies are data      |
| `tests.test_foo.Bar`      | `<dotted-id>`     | Dotted module paths are data |
| `# comment text`          | stripped           | Comments are noise           |
| `\<newline>`              | joined            | Line continuations are noise |

## Per-executable handlers

This is where the real value is. The generic tokenizer gets you 80%.
Handlers for specific executables get you to 95%+.

Current handlers in the Python implementation:

- **echo/printf**: All positionals collapse to single `<str>`. Flags
  preserved. Subshells preserved verbatim (safety invariant).
- **sed**: Script expressions (`-e` arg or first positional) collapse
  to `<sed-expr>`. Address ranges like `100,230p` also collapse.
- **grep/egrep/fgrep/rg/ag/ack**: First positional (or `-e` arg) is
  `<pattern>`. `-A/-B/-C/-m` consume a numeric arg as `N`. Remaining
  positionals are files.
- **find**: `-name/-iname/-path/-ipath` args collapse to `<pattern>`.
  Other positionals are paths. Primaries like `-type f` kept verbatim.

Handlers that would be worth adding:

- **curl/wget**: URL is already handled generically, but `-H` header
  values, `-d` data bodies, `-o` output paths could be handler-aware.
- **awk**: First positional is a program, collapse to `<awk-expr>`.
- **docker**: Nested subcommand structure (`docker compose up` vs
  `docker run`).
- **jq**: First positional is a filter expression.
- **xargs**: Everything after xargs flags is a sub-command.
- **sort/cut/tr/wc/head/tail**: Mostly handled by generic rules but
  could be tightened.

## Handler architecture

Each handler is a function with a standard signature:

```go
type HandlerFunc func(tokens []string, h Helpers) []string
```

`Helpers` provides the shared utilities every handler needs:

- `IsFlag(tok) bool`
- `IsSubshell(tok) bool`
- `ClassifyToken(tok) string`
- `SplitRedirects(tokens) (args, redirects []string)`

Handlers receive the token list starting after the executable (and
subcommand, if applicable). They return a list of classified tokens.

Registration is a simple map:

```go
var handlers = map[string]HandlerFunc{
    "grep":  handleGrep,
    "egrep": handleGrep,
    "rg":    handleGrep,
    "sed":   handleSed,
    "echo":  handleEcho,
    // ...
}
```

Each handler lives in its own file (`handler_grep.go`,
`handler_sed.go`, etc.) for contributor ergonomics. A new contributor
copies an existing handler, adapts it, adds tests, submits a PR.

## The subshell-preservation invariant

This is the one safety-critical rule. Per-executable handlers MUST
preserve `$(...)` tokens verbatim instead of collapsing them to a
placeholder.

**Why**: If `echo hello` is classified SAFE and cached under shape
`echo <str>`, then `echo $(rm -rf /)` hitting the same shape would
auto-allow a destructive subshell from a benign cache entry. Silent
cache-poisoning.

**How**: Each handler checks `IsSubshell(tok)` before its positional
collapse branch and keeps the token verbatim if true.

This invariant must be enforced in tests for every handler. The test
template: assert that `cmd $(dangerous)` produces a different shape
than `cmd benign-literal`.

## Not configurable

There is one normalization behavior. No modes, no profiles, no config
files. Reasons:

- The core rules (paths, numbers, URLs, hashes) are not debatable.
  Nobody wants different SHAs to produce different shapes.
- Configuration creates compatibility problems. If two systems use
  different shellshape configs, their shapes don't match.
- The long tail is handled by adding handlers, not by tuning knobs.
  A handler for `curl` is better than a config option for "collapse
  HTTP headers."
- The value proposition is "it just works." Adding `--aggressive` vs
  `--conservative` modes doubles the test surface for no real benefit.

The only future exception might be `--json` for output format. That's
presentation, not normalization behavior.

## Testing strategy

Tests are the product. The normalizer is only as good as its test
corpus.

### Structure

```
normalize_test.go          # core rules, token classification
handler_grep_test.go       # per-handler tests
handler_sed_test.go
handler_echo_test.go
handler_find_test.go
testdata/
  corpus.txt               # real-world commands, one per line
  corpus_expected.txt      # expected shapes, aligned by line number
```

### Test patterns

Every test case should follow one of these shapes:

1. **Direct assertion**: Input command produces expected shape.
2. **Collision assertion**: Two commands that differ only in data
   produce the same shape. This is the cache-hit property.
3. **Non-collision assertion**: Two semantically different commands
   produce different shapes. Prevents over-normalization.
4. **Subshell safety**: `cmd $(dangerous)` produces a different shape
   than `cmd literal`. One per handler, never delete these.

### Growing the corpus

The review workflow from claude-approvals transfers directly:

1. Collect real commands from users (opt-in logging, submitted examples,
   CI command logs).
2. Query for single-hit shapes (commands that didn't collide with
   anything). These are normalization gaps.
3. Sort alphabetically and scan for near-duplicates that differ in one
   token. That token is the leak.
4. Write failing test, implement fix, run full suite.

Contribution template for new handlers: "Here are 10 real commands for
<executable> that should collapse to 3-4 shapes. Here are the tests.
Here is the handler."

## Distribution

- **GitHub releases**: goreleaser builds binaries for
  darwin/linux x amd64/arm64. Attached to GitHub releases.
- **Homebrew**: goreleaser generates a tap formula automatically.
  `brew install <user>/tap/shellshape`.
- **go install**: `go install github.com/<user>/shellshape@latest`.
- **Docker**: Maybe. Single binary means docker is rarely needed but
  some CI environments prefer it.

## Project structure

```
shellshape/
  main.go                  # CLI entry point, stdin/argv handling
  normalize.go             # core pipeline: split, tokenize, classify
  normalize_test.go        # core tests
  token.go                 # token classification (_classify_token equiv)
  token_test.go
  handlers.go              # registry + Helpers type
  handler_echo.go          # per-executable handlers
  handler_echo_test.go
  handler_grep.go
  handler_grep_test.go
  handler_sed.go
  handler_sed_test.go
  handler_find.go
  handler_find_test.go
  testdata/
    corpus.txt
    corpus_expected.txt
  .goreleaser.yml
  go.mod
  README.md
  LICENSE                  # MIT
```

## Migration from Python

The Python `normalize.py` is the reference implementation. The Go
rewrite should produce identical output for every test case in the
Python suite. The migration path:

1. Port `tests/test_normalize.py` test cases to Go table-driven tests.
   This is step one, before writing any implementation.
2. Implement core pipeline (comment strip, continuation join, heredoc
   collapse, top-level split, subshell normalization, shlex-equivalent
   tokenizer).
3. Implement token classification.
4. Implement handlers one at a time, each with its ported tests.
5. Run the corpus test (all Python test inputs through Go, diff outputs).
6. Ship v0.1.0.

The Python implementation stays in claude-approvals and continues to
be used there. shellshape is a separate project. If shellshape
eventually becomes a binary dependency of claude-approvals, the
Python normalizer can be replaced with a subprocess call or removed
entirely.

## What shellshape is NOT

- **Not a shell parser.** It doesn't produce an AST. It doesn't
  understand bash grammar. It handles the 95% case with shlex-style
  tokenization and degrades gracefully on the rest.
- **Not a security tool.** It's a normalizer. It can be used as an
  input to a security classifier (like claude-approvals does) but it
  doesn't make safety decisions itself.
- **Not configurable.** One behavior, one output. No modes, no
  profiles, no plugins at runtime. Extensibility is through code
  contributions (new handlers), not configuration.
- **Not a linter or analyzer.** It doesn't tell you if your command
  is correct, idiomatic, or dangerous. It tells you what shape it is.
