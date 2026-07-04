# Maya

**A fine-grained reactive UI framework for Go, compiled to WebAssembly — with zero external dependencies.**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go 1.24](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![WebAssembly](https://img.shields.io/badge/target-WASM-654FF0?logo=webassembly&logoColor=white)](https://webassembly.org/)
[![Dependencies](https://img.shields.io/badge/dependencies-none%20(stdlib%20only)-brightgreen.svg)](go.mod)
[![Lines of code](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/JoaoHenriqueBarbosa/maya/main/.github/badges/loc.json)](internal)

> **What this is: a completed experiment with a conclusion.** Maya asked one concrete
> question — *can you build SolidJS-style fine-grained reactivity in idiomatic Go, in the
> browser via WebAssembly, using nothing but the standard library?* — and answered it by
> **building the whole thing and measuring it**: signals, effects, memos, batching, two
> renderers, a WASM bundle. The measured answer is **no — Go's runtime is the wrong substrate
> for this model** (the GC and goroutine scheduler impose overhead that a fine-grained,
> allocation-light invalidation model fights rather than uses). That answer *is the
> deliverable* — see [The result](#the-result-an-engineering-conclusion). The framework works
> and is instructive to read; it was deliberately stopped once the question was answered, so
> treat it as a reference experiment, not a dependency. [Honest status](#honest-status) is
> candid about the rough edges that came with stopping on purpose.

---

## What is this?

Maya is an experiment in bringing **fine-grained reactivity** — the model popularized by
[SolidJS](https://www.solidjs.com/) — to Go on the web. Instead of a virtual DOM that diffs a
whole tree on every change, Maya tracks dependencies at the level of individual reactive values
(`Signal`s) and updates only the exact pieces of the UI that actually depend on what changed.

The whole thing compiles to a single `.wasm` file with the **standard Go compiler** (`GOOS=js
GOARCH=wasm`) and runs in the browser. There are **no third-party dependencies** — only the Go
standard library, leaning on newer features (`iter`, `weak`, `sync/atomic`,
`runtime.AddCleanup`, `syscall/js`, `//go:wasmexport`).

```go
counter := maya.Signal(0)

doubled := maya.Memo(func() int {
    return counter.Get() * 2 // re-runs only when `counter` changes
})

// Only the text node bound to `counter` is touched in the DOM when you click.
maya.TextSignal(counter, func(v int) string { return fmt.Sprintf("%d", v) })
```

## Highlights

- **Fine-grained reactive core** — `Signal`, `Memo`/`Computed`, and `Effect` with **automatic
  dependency tracking**. Reading a signal inside an effect registers the dependency; you never
  wire subscriptions by hand.
- **Versioned change detection** — each `Signal` carries an atomic version counter; effects
  remember the version they last saw, so re-runs are skipped when nothing they read actually
  moved.
- **Lazy, cached derivations** — `Memo`/`Computed` recompute lazily and cache their result
  behind a `weak.Pointer`, so unused memos can be collected.
- **Batched updates** — mutations are coalesced and flushed on a ~16ms (≈60fps) tick, so a
  burst of `Set`s produces a single render pass.
- **Two interchangeable renderers behind one interface** — a **DOM** renderer that does *real
  selective updates* (it mutates only the affected `textContent`, not the subtree) and a
  **Canvas2D** renderer (full redraw). Pick at runtime via `window.MAYA_RENDERER`.
- **Declarative widget API** — `Container`, `Column`, `Row`, `Text`, `Title`, `Button`, plus the
  reactive bindings `TextSignal` / `TextMemo`.
- **Zero external dependencies** — standard library only. Nothing to audit in `go.sum`
  (there is no `go.sum`).
- **Measured, not guessed** — the reactive core and widgets ship with a unit-test suite and
  microbenchmarks (tree traversal, signals, memos, effects, batching). The benchmarks are the
  point: they are what turned "Go feels wrong for this" into a measured engineering conclusion.
  (Note: the suite does not currently pass on `main` — see [Honest status](#honest-status).)

## Requirements

- **Go 1.24 or newer** (the code uses `iter`, `weak`, `//go:wasmexport`, and other 1.24
  features — it will not build on older toolchains).
- A **browser with WebAssembly** support to run the example.
- No package manager, no Node, no bundler. The build is `go build`.

## Getting started

Maya is a library plus a runnable example. There is no installable CLI yet.

```bash
git clone https://github.com/JoaoHenriqueBarbosa/maya.git
cd maya
go test ./internal/...   # run the reactive-core + widgets test suite
```

### Build and run the example

The example lives in [`examples/simple/`](examples/simple) — a reactive counter with `Memo`,
`Computed`, and both renderers wired up.

```bash
cd examples/simple

# 1. Compile the app to WebAssembly with the standard Go compiler
GOOS=js GOARCH=wasm go build -o app.wasm main.go

# 2. Make sure the Go WASM JS shim is present (ships with your Go install)
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" .   # path may be misc/wasm on older Go

# 3. Serve the directory (the tiny bundled server works too)
go run server.go 8080
#   → http://localhost:8080/index_dom.html      (DOM renderer)
#   → http://localhost:8080/index_canvas.html   (Canvas2D renderer)
#   → http://localhost:8080/compare.html        (both, side by side)
```

The standard-compiler build produces a **~3.2 MB** `.wasm` (the Go runtime is baked in; there
is no TinyGo path — see [Honest status](#honest-status)).

### Writing a component

A Maya app is a function that returns a widget tree. Reactive values drive the parts of the
tree that are bound to them:

```go
//go:build wasm

package main

import (
    "fmt"

    "github.com/maya-framework/maya"
    "github.com/maya-framework/maya/internal/widgets"
)

func main() {
    count := maya.Signal(0)

    app := maya.New(func() widgets.WidgetImpl {
        return maya.Column(
            maya.Title("Counter"),
            maya.Row(
                maya.Text("Count: "),
                // Only this text node updates when `count` changes.
                maya.TextSignal(count, func(v int) string { return fmt.Sprintf("%d", v) }),
            ),
            maya.Button("Increment", func() { count.Set(count.Get() + 1) }),
        )
    })

    app.Run()
}
```

Choose the renderer from the HTML host **before** loading the module:

```html
<script>
  window.MAYA_RENDERER = 'dom'; // or 'canvas'; defaults to 'dom'
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
    .then((r) => go.run(r.instance));
</script>
```

## How the reactivity works

The interesting part of Maya is `internal/reactive/`. In short:

- **`Signal[T]`** holds a value, an atomic version counter, and a set of observing effects.
  `Get()` records the current effect as an observer (dependency tracking); `Set()` bumps the
  version and notifies observers. `Peek()` reads without tracking.
- **`Effect`** runs a function, and while it runs, any signal it reads registers itself in the
  effect's dependency map — keyed by the *version* seen at that read. On re-run, the effect
  clears stale dependencies and re-tracks, so the dependency graph is always current.
- **`Memo` / `Computed`** are effects that produce a value. They compute lazily and cache the
  result behind a `weak.Pointer`, recomputing only when an upstream signal changes.
- **Batching** collects dirty work and flushes it on a timer (~16ms), so N synchronous `Set`s
  collapse into one render.
- **Rendering** is decoupled through the `Renderer` interface. The DOM renderer performs
  genuinely selective updates (it locates the bound node and rewrites only its text); the
  Canvas2D renderer repaints the frame.

This is the design that lets a click on "Increment" touch a single DOM text node instead of
re-rendering the tree.

## Development

```bash
go test ./internal/...              # unit tests (core, reactive, widgets, graph)
go test -race ./internal/reactive   # the reactive core under the race detector
go test -bench=. -benchmem ./...    # microbenchmarks
go vet ./internal/...               # static checks
gofmt -s -l .                       # formatting

# Build the example to WASM (see "Getting started" for the full run steps)
cd examples/simple && GOOS=js GOARCH=wasm go build -o app.wasm main.go
```

> **Note:** the top-level `Makefile` predates the current layout and targets a `cmd/maya`
> package that does not exist, plus optional tooling (TinyGo, WASI, WebGPU codegen) that this
> repository does not use. Prefer the `go` commands above. Cleaning up or removing the Makefile
> is a good first contribution — see [CONTRIBUTING.md](CONTRIBUTING.md).

## Architecture

```
maya/
├── maya.go                 # public API (build tag `wasm`): New, Signal, Memo, Container, Text, Button, TextSignal, ...
├── exports.go              # //go:wasmexport bridge to JS (DOM-ready, event dispatch)
├── internal/
│   ├── reactive/           # the reactive core: Signal, Effect, Memo/Computed, tracking, batching
│   ├── widgets/            # BaseWidget, Text, Button, Container, Column, Row, RenderObject
│   ├── render/             # Renderer interface + DOM renderer (selective) + Canvas2D renderer (redraw)
│   ├── core/               # widget tree + node model (DFS traversal, dirty flags)
│   ├── graph/              # dependency-graph helpers
│   └── workflow/           # scheduling scaffolding
└── examples/
    └── simple/             # runnable reactive counter (DOM + Canvas, side-by-side compare page)
```

The public surface (`maya.go`) only compiles under the `wasm` build tag, because it depends on
`syscall/js`. The `internal/reactive`, `internal/core`, and `internal/graph` packages are plain
Go and are what the host-side test suite exercises.

## Honest status

This repository is a **spike**, and it's worth being upfront about what that means:

- **It was built in a single day** and is **not under active development**. Treat it as a
  reference/experiment, not a dependency.
- **The test suite does not currently pass on `main`.** At least `TestContainer_Build` (widgets)
  and one benchmark are broken. Fixing the suite is a great entry point for contributors.
- **There is debug output in production code.** The reactive core and the public bindings print
  diagnostic lines (`[SIGNAL] ...`, `[EFFECT] ...`, `[TEXT-SIGNAL] ...`) via `println`/`fmt`.
  This is noise that belongs behind a logging flag or removed entirely.
- **Some pieces are stubbed or rough** — e.g. parts of the `Memo` integration with the effect
  interface, and the app-level "selective" update currently falls back to a full pipeline pass
  in places despite the DOM renderer supporting true selective updates.
- **The bundle is ~3.2 MB.** That's the standard Go WASM runtime; there is no size-optimized
  (TinyGo) build path in this repo.

### Not implemented (despite what older design docs in this repo say)

The `.md` design documents in this repository (`OVERVIEW.md`, `BREAKDOWN.md`, `ROADMAP.md`,
etc.) were written as an ambitious vision and describe a number of things that **were never
built**. To be clear, the following are **not** part of the code and should not be assumed:

- WebGPU / WGSL compute-shader rendering
- A TinyGo build path or an ~80 KB bundle
- Spring-physics animation
- R-Tree spatial indexing

If it isn't in `internal/` or `maya.go`, it doesn't exist. Please read the code, not the vision
docs, when deciding what Maya can do.

## The result: an engineering conclusion

The point of Maya was never to ship a framework. It was to answer a question by building and
measuring, and then to **trust the measurement enough to stop**.

The conclusion: **Go is the wrong substrate for fine-grained reactivity.** The model — the one
SolidJS uses — wants deterministic, synchronous, allocation-light signal propagation with
precise invalidation. Go's runtime pulls the other way: the garbage collector and the
goroutine scheduler add overhead exactly on the hot path where this model needs to be lean.
You *can* build it — this repository is the proof, down to a working reactive counter running
in the browser under two renderers — but you spend your effort fighting the runtime instead of
leaning on it, and the WASM payload (~3.2 MB of baked-in Go runtime) is heavy next to
JS-native solutions. The benchmarks in `internal/` are what moved this from a hunch to a
finding.

That is why the project stopped. Not because it failed — because it **succeeded at being an
experiment**: it produced a clear, measured answer, and the answer was "not this substrate."
Knowing when to invalidate an approach on the evidence, and to bank the lesson instead of
sinking more time into it, is the whole exercise. The conclusion is the deliverable.

*(If you want the theory this validates: it's why Solid is JavaScript, and why reactive
frameworks that care about the hot path avoid heavily-GC'd, scheduler-mediated runtimes there.
Maya is one data point that arrives at that conclusion the hard way — by building it.)*

## Contributing

Contributions are welcome — especially "make the existing thing solid" work (fixing the failing
tests, removing debug output, finishing the `Memo`/`Effect` wiring, cleaning up the Makefile).
See [CONTRIBUTING.md](CONTRIBUTING.md) for setup and the PR workflow, and
[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) for community expectations.

## Security

For vulnerability reports, see [SECURITY.md](SECURITY.md).

## License

Released under the [MIT License](LICENSE).
