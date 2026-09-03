# Appendix: Response to samber-do-auditlog Adoption Feedback

**Original report:** [`2026-06-17_adoption-feedback-from-samber-do-auditlog.html`](./2026-06-17_adoption-feedback-from-samber-do-auditlog.html)\
**Library version reviewed:** go-output v0.11.0\
**Response authored:** 2026-06-17\
**Status:** Addressed in `master` (`bbc8f89` and ancestors)

---

## Executive Summary

All four critical gaps and three of the four important gaps identified in the feedback report have been resolved. The one remaining important item (`io.Writer` support for diagram renderers) is already possible via the existing `StreamingRenderer` adapter, so no new API surface was added.

| Severity  | Resolved | Count |
| --------- | -------- | ----- |
| Critical  | Yes      | 4/4   |
| Important | Yes      | 3/4   |
| Nice      | Yes      | 1/1   |

---

## Issue-by-Issue Response

### Critical 1 — GraphStyle ignored in Mermaid and PlantUML

**Status:** ✅ Resolved

Mermaid now emits per-node `style <id> ...` directives for non-zero `GraphStyle` values, replacing the hardcoded pink `classDef default`. PlantUML now emits inline color specs (`#[#fill;line:stroke;text:font]`).

```go
// graph/mermaid.go
r.AddNode(output.GraphNode{
    ID:    output.NewBrandedID[output.GraphNodeIDBrand]("svc"),
    Label: output.NewBrandedID[output.GraphNodeLabelBrand]("Service"),
    Style: output.GraphStyle{FillColor: "#e8a838", StrokeColor: "#4a4030"},
})
```

**Files changed:** `graph/mermaid.go`, `plantuml/plantuml.go`, `graph/mermaid_test.go`, `plantuml/plantuml_test.go`

---

### Critical 2 — No edge deduplication

**Status:** ✅ Resolved

Added `GraphRendererState.DedupEdges()` in `output` package. It removes duplicate edges by `(from, to)` in-place, keeping the first occurrence. This is opt-in to preserve backwards compatibility.

```go
renderer.AddEdge(edgeA)
renderer.AddEdge(edgeA) // duplicate
renderer.DedupEdges()   // second edge removed
```

**Files changed:** `graph.go`, `graph_state_test.go`

---

### Critical 3 — Mermaid output locked behind markdown fence

**Status:** ✅ Resolved

`MermaidRenderer.SetCodeFence(bool)` toggles the fence. Default remains `true` for backwards compatibility.

```go
r := graph.NewMermaidRenderer()
r.SetCodeFence(false) // raw flowchart syntax
```

**Files changed:** `graph/mermaid.go`, `graph/mermaid_test.go`

---

### Critical 4 — Root go.mod transitively pulls YAML + TOML

**Status:** ✅ Resolved

Removed `github.com/larsartmann/go-output/serialization` from root `go.mod`. The root module now requires only `go-branded-id`, `delimited`, `enum`, `testhelpers`, and `x/term`. Importing `graph/`, `plantuml/`, or `d2/` no longer transitively pulls `go-faster/yaml`, `go-toml/v2`, `go-faster/jx`, or `segmentio/asm`.

```text
Before: 11 unwanted external modules for a diagram renderer
After:  0 unwanted external modules
```

**Files changed:** `go.mod`, `userjourney_test.go`

---

### Important 5 — SlugifyID gaps

**Status:** ✅ Resolved

`escape.SlugifyID` now replaces `. * [ ] { } ( )` in addition to spaces, hyphens, and slashes. This eliminates collisions like `foo.bar` ↔ `foo_bar`.

**Files changed:** `escape/escape.go`, `escape/escape_test.go`

---

### Important 6 — Hardcoded DOT layout attributes

**Status:** ✅ Resolved

DOT layout is now configurable via typed enums and setters:

```go
r := graph.NewDOTRenderer().
    SetRankDir(graph.RankDirLR).
    SetSplines(graph.SplineSpline).
    SetNodeSep("0.8").
    SetRankSep("1.0")
```

`RankDir` and `SplineStyle` are typed enums following the project's `Parse/String/IsValid/AllowedValues` pattern, making invalid values unrepresentable.

**Files changed:** `graph/dot.go`, `graph/dot_enum.go`, `graph/dot_enum_test.go`, `graph/dot_test.go`, `graph/go.mod`

---

### Important 7 — No io.Writer support for diagram renderers

**Status:** ✅ Already supported

The `output.StreamingRenderer` interface plus `StreamingRendererFromRenderer()` adapter already provides `Stream(io.Writer) error` for any renderer. For example:

```go
streaming := output.StreamingRendererFromRenderer(graph.NewMermaidRenderer())
err := streaming.Stream(w)
```

This is the intended pattern for diagrams. No new code was added.

---

### Nice 8 — Pre-v1 stability promise

**Status:** ✅ Resolved

ADR 006 now documents stable vs experimental diagram APIs. Stable APIs include `GraphRenderer`, core renderer constructors, `GraphRendererState`, and the new `DedupEdges`/`SetCodeFence`/typed-enum additions. D2's rich domain model and NOM/TUI remain experimental.

**Files changed:** `docs/adr/006-api-stability.md`

---

## Verification

- `nix run .#test` — all modules pass
- `nix run .#lint` — zero issues
- `go run ./examples/basic mermaid` — outputs raw flowchart with amber theme
- `go run ./examples/basic dot` — outputs left-to-right graph with configurable layout

---

## What Was Intentionally Not Implemented

- **`WithTheme(Theme)` global option:** The existing `GraphStyle` struct already covers per-node coloring and is now honored by all three renderers. A higher-level `Theme` abstraction would add a new type without additional expressiveness over `GraphStyle`. If a global theme becomes necessary later, it can be composed from `GraphStyle` defaults.

---

## Remaining Open Questions

None blocking. Potential follow-ups for future releases:

1. Should `DedupEdges()` become automatic in `Render()` rather than opt-in?
2. Should `MermaidRenderer` default to `codeFence = false` in a future v1?
3. Should edge labels participate in dedup keys when present?
