# Root Package Split Proposal

**Date:** 2026-05-25

## Problem

Root has **54 `.go` files** (~8K lines) in a single `package output`. Three cross-cutting unexported helpers create coupling that blocks extraction:

| Helper | Used by | Block extraction? |
|---|---|---|
| `marshal.go` (`marshal`, `unmarshal`, `brandedValue`) | JSON, YAML, XML, json_renderers, yaml_renderers | Yes — must be exported or moved |
| `markup.go` (`writeMarkupRow`, `writeMarkupColumns`) | XML, HTML | Yes — must be exported or moved |
| `delimited.go` (`DelimitedWriter`) | CSV, TSV | Yes — must be exported or moved |

## Dependency Map (Current)

```
                    renderer.go ──────────────────────────────────┐
                    format.go ────────────────────────────────┐   │
                    shape.go                                  │   │
                    ids.go ─────────────────────────────────┐ │   │
                    marshal.go ─────────────────────────┐   │ │   │
                    tabledata.go ───────────────────┐   │   │ │   │
                    delimited.go ──────────────┐    │   │   │ │   │
                    markup.go ────────────┐    │   │   │   │ │   │
                                                 │   │   │   │ │   │
    csv.go ──────────────────────────────────►   ✓   ✓   │   │   │
    tsv.go ──────────────────────────────────►   ✓   ✓   │   │   │
    xml.go ──────────────────────────────►       ✓   ✓   ✓   │   │
    html.go ─────────────────────────────►   ✓       ✓   ✓   │   │
    markdown.go ─────────────────────────►               ✓   │   │
    tree.go ─────────────────────────────►   ✓       ✓       │   │
    json.go ─────────────────────────────►   ✓   ✓   ✓   ✓   ✓   │
    json_renderers.go ───────────────────►   ✓           ✓   ✓   │
    yaml.go ─────────────────────────────►   ✓   ✓   ✓   ✓   ✓   │
    yaml_renderers.go ───────────────────►   ✓           ✓   ✓   │
    graph.go ────────────────────────────►               ✓       │
    graph_tabledata.go ─────────────────►               ✓   ✓       (also graph.go)
    streaming.go ───────────────────────►                   ✓   ✓
    render_tabledata.go ────────────────► calls csv,tsv,md,xml,yaml,html,tree
```

Legend: ✓ = file depends on that column's file.

### Key Structural Patterns

1. **`tableDataBase`** (in `tabledata.go`) is embedded by: `json.go`, `yaml.go`, `html.go`, `streaming.go`
2. **`GraphRendererMixin`** (in `graph.go`) is embedded by: `json_renderers.go`, `yaml_renderers.go`
3. **`markup.go`** helpers are used by: `xml.go`, `html.go`
4. **`delimited.go`** is used by: `csv.go`, `tsv.go`
5. **`marshal.go`** unexported helpers used by: `json.go`, `yaml.go`, `xml.go`, `json_renderers.go`, `yaml_renderers.go`
6. **`render_tabledata.go`** is the top-level dispatcher depending on nearly everything

---

## Strategy A: Extract Per-Format Modules (Recommended)

Move each format family into its own sub-module. Export the shared helpers.

```
go-output/                        # Root: core types only (~20 files)
├── format.go, shape.go, renderer.go, ids.go, registry.go
├── sort.go, color.go, slices.go, tabledata.go, tree.go, graph.go
├── graph_tabledata.go            # stays — uses root types only
│
├── delimited/                    # NEW MODULE: shared CSV/TSV plumbing
│   ├── delimited.go              # DelimitedWriter (exported)
│   └── go.mod
│
├── csv/                          # NEW MODULE
│   ├── csv.go
│   └── go.mod → root, delimited
│
├── tsv/                          # NEW MODULE
│   ├── tsv.go
│   └── go.mod → root, delimited
│
├── json/                         # NEW MODULE
│   ├── json.go, json_renderers.go
│   └── go.mod → root
│
├── yaml/                         # NEW MODULE
│   ├── yaml.go, yaml_renderers.go
│   └── go.mod → root, go-faster/yaml
│
├── xml/                          # NEW MODULE
│   ├── xml.go
│   └── go.mod → root, escape
│
├── html/                         # NEW MODULE
│   ├── html.go, streaming.go
│   └── go.mod → root
│
├── markdown/                     # NEW MODULE
│   ├── markdown.go
│   └── go.mod → root
```

Root drops from ~54 to ~20 files (core types + graph/tree + `render_tabledata.go` which becomes a thin dispatcher importing format modules via registry).

### Key Changes Required

1. **`marshal()` / `unmarshal()` / `brandedValue()`** → export as `Marshal()` / `Unmarshal()` / `BrandedValue()` in root (they're generically useful)
2. **`markup.go` helpers** → export in root or move into a shared `markup/` package
3. **`DelimitedWriter`** → move to `delimited/` sub-module (it's only used by CSV/TSV)
4. **`render_tabledata.go`** stays in root as a dispatcher, but calls registered renderers instead of direct function calls
5. Each format module registers itself via `init()` or users call `Register()` explicitly

### Result

Users who only want JSON get zero CSV/TSV/XML/HTML/Markdown/YAML deps.

---

## Strategy B: Group by Domain (Simpler)

Fewer modules, less granular isolation:

```
go-output/                        # Root: core types + lightweight formatters
├── format.go, shape.go, renderer.go, ids.go, registry.go
├── sort.go, color.go, slices.go, tabledata.go
├── tree.go, graph.go, graph_tabledata.go
├── render_tabledata.go
├── markdown.go                   # stays — zero deps, tiny
│
├── serialization/                # NEW MODULE: JSON + YAML (shared marshal logic)
│   ├── json.go, json_renderers.go, yaml.go, yaml_renderers.go, marshal.go
│   └── go.mod → root
│
├── delimited/                    # NEW MODULE: CSV + TSV
│   ├── csv.go, tsv.go, delimited.go
│   └── go.mod → root
│
├── markup/                       # NEW MODULE: XML + HTML (shared markup logic)
│   ├── xml.go, html.go, streaming.go, markup.go
│   └── go.mod → root, escape
```

Root drops from ~54 to ~25 files. Less isolation but simpler.

---

## Comparison

| | Strategy A (per-format) | Strategy B (grouped) |
|---|---|---|
| Modules added | 7 new | 3 new |
| Root file reduction | ~54 → ~20 | ~54 → ~25 |
| Dep isolation | Best — per-format | Good — per-family |
| Maintenance | More go.mod files | Fewer go.mod files |
| User ergonomics | `import "go-output/json"` | `import "go-output/serialization"` |

## Decision

**Strategy B** — chosen 2026-05-25. Group by domain keeps the module count manageable (3 new) while still isolating heavy dependencies (go-faster/yaml, encoding/csv). Root stays coherent with core types + lightweight formatters.

## Risks and Considerations

- **Breaking change:** Users currently importing `output.NewJSONRenderer()` would need to change to `json.NewRenderer()` or use the registry
- **Registry migration:** `render_tabledata.go` currently calls format functions directly — must switch to registry-based dispatch
- **Test files:** ~20 test files move with their respective format modules
- **`tableDataBase` embedding:** Must remain accessible from format modules — either exported in root or moved to a shared package
- **`GraphRendererMixin` embedding:** Same concern for JSON/YAML graph renderers
