# Contributing to Maya

Thanks for your interest in Maya! This is an exploratory prototype (see the
[Honest status](README.md#honest-status) section of the README), so the most valuable
contributions right now are the "make the existing thing solid" kind:

- **Fixing the currently failing tests** (`TestContainer_Build` in `internal/widgets`, and the
  broken benchmark).
- **Removing debug output** — there is diagnostic `println`/`fmt` logging in production code
  (`[SIGNAL]`, `[EFFECT]`, `[TEXT-SIGNAL]`, `[RENDER]`, …). It should go behind a flag or be
  deleted.
- **Finishing rough edges** — the `Memo`/`Effect` interface wiring, and making the app-level
  update path actually use the DOM renderer's selective update everywhere.
- **Cleaning up the `Makefile`**, which targets a non-existent `cmd/maya` package and tooling
  (TinyGo/WASI/WebGPU) the repo doesn't use.
- Documentation, examples, and small quality-of-life fixes.

By participating, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Prerequisites

- **Go 1.24 or newer** — required. The codebase uses `iter`, `weak`, `//go:wasmexport`, and
  other 1.24 features and will not build on older toolchains. Check with `go version`.
- A **browser with WebAssembly** to run the example.
- No Node, no bundler, no package manager. Everything is `go`.

## Getting set up

```bash
# 1. Fork the repo on GitHub, then clone your fork
git clone https://github.com/<your-username>/maya.git
cd maya

# 2. Add the upstream remote so you can keep your fork in sync
git remote add upstream https://github.com/JoaoHenriqueBarbosa/maya.git

# 3. There are no external dependencies to install — it's stdlib only.
#    Just confirm the host-side packages build and test:
go build ./internal/...
go test ./internal/...
```

To build and run the browser example, follow the steps in
[Getting started](README.md#build-and-run-the-example).

## Repository layout

```
maya/
├── maya.go          # public API (build tag `wasm`, depends on syscall/js)
├── exports.go       # //go:wasmexport JS bridge
├── internal/
│   ├── reactive/    # Signal, Effect, Memo/Computed, tracking, batching  ← the core
│   ├── widgets/     # Text, Button, Container, Column, Row, RenderObject
│   ├── render/      # Renderer interface + DOM + Canvas2D renderers
│   ├── core/        # widget tree / node model
│   ├── graph/       # dependency-graph helpers
│   └── workflow/    # scheduling scaffolding
└── examples/simple/ # runnable reactive counter
```

Most changes land in `internal/reactive/` (the interesting part) or `internal/widgets/`. Note
that `maya.go` / `exports.go` only compile under the `wasm` build tag — the host `go test`
exercises the `internal/*` packages that are plain Go.

## Development workflow

1. **Create a branch** off `main`:
   ```bash
   git switch -c fix/failing-container-test
   ```
2. **Make your change**, keeping it focused. If you touch the reactive core, run it under the
   race detector.
3. **Format, vet, and test** before you push:
   ```bash
   gofmt -s -w .
   go vet ./internal/...
   go test ./internal/...
   go test -race ./internal/reactive
   ```
   If your change is performance-sensitive, also run the benchmarks:
   ```bash
   go test -bench=. -benchmem ./...
   ```
   If your change affects the WASM build, confirm it still compiles:
   ```bash
   cd examples/simple && GOOS=js GOARCH=wasm go build -o /tmp/app.wasm main.go
   ```
4. **Commit** using Conventional Commits (see below).
5. **Push** to your fork and **open a pull request** against `JoaoHenriqueBarbosa/maya:main`.
   Fill in the PR template and check off the boxes.

## Commit messages — Conventional Commits

Please format commit messages as `type(scope): short description`. The scope is optional but
encouraged (e.g. `reactive`, `widgets`, `render`, `core`, `examples`, `docs`).

| Type       | When to use it                                                        |
|------------|-----------------------------------------------------------------------|
| `feat`     | A new feature or public API addition                                  |
| `fix`      | A bug fix                                                             |
| `docs`     | Documentation only                                                    |
| `test`     | Adding or fixing tests                                                |
| `refactor` | Code change that neither fixes a bug nor adds a feature              |
| `perf`     | A change that improves performance                                    |
| `style`    | Formatting / whitespace / `gofmt` (no behavior change)               |
| `build`    | Build system, `go.mod`, WASM build, Makefile                         |
| `ci`       | CI configuration                                                      |
| `chore`    | Housekeeping that doesn't fit the above                               |

Examples:

```
fix(widgets): make TestContainer_Build pass by initializing child bounds
refactor(reactive): remove debug println from Signal.Set
docs(readme): correct the wasm_exec.js path for Go 1.24
```

## Pull request expectations

- Keep PRs focused and reasonably small; unrelated changes should be separate PRs.
- Make sure `go test ./internal/...`, `go vet ./internal/...`, and `gofmt -s -l .` are clean.
- Describe **what** changed and **why**. If you fixed a failing test, say what was actually
  wrong.
- Don't add third-party dependencies. Maya is intentionally standard-library only — if you
  believe a dependency is genuinely necessary, open an issue to discuss it first.

## Reporting bugs and requesting features

Use the issue templates:
[bug report](.github/ISSUE_TEMPLATE/bug_report.yml) ·
[feature request](.github/ISSUE_TEMPLATE/feature_request.yml).

For security issues, please **do not** open a public issue — follow [SECURITY.md](SECURITY.md).

Thank you for helping improve Maya!
