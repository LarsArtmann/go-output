# Naming Review Report — 2026-06-15

## Executive Summary

- **81 exported types** reviewed across 15 modules
- **0** Manager/Handler/Processor/Helper/Util classes
- **0** Impl/Concrete/Default suffixes
- **0** Abstract/Base prefixes
- **1** "vague" type name: `TableData` (legitimate domain concept — a table of data)
- **1** vague verb: `HandleError` (example code only)
- **1** terminology split-brain: `Marshaler` (registry) vs `Renderer` (everywhere else)
- Naming health: **Excellent**. This is one of the cleanest-named Go codebases the review has seen.

## Automated Detection Results

| Smell                                   | Hits | Verdict                              |
| --------------------------------------- | ---- | ------------------------------------ |
| Vague nouns (Data/Info/Record/...)      | 1    | `TableData` — domain-accurate, keep  |
| Manager/Handler/Processor/Helper/Util   | 0    | Clean                                |
| Impl/Concrete/Default suffix            | 0    | Clean                                |
| Abstract/Base prefix                    | 0    | Clean                                |
| Vague verbs (do/handle/process/manage)  | 1    | `HandleError` (examples only)        |

## Strengths (Good Naming)

- `Renderer` / `TableRenderer` / `GraphRenderer` / `TreeOutputRenderer` — honest, role-based interface names
- `Format` / `Shape` / `ColorMode` / `Alignment` — single-word domain enums
- `BrandedID[T]` / `D2NodeID` / `TreeNodeID` — phantom-type pattern, type-encodes intent
- `GraphRendererState` — composition root, name says "shared state"
- `NodeEdgeAppender` — interface name describes the capability precisely
- `TableDataStore` — honest: it's the store/backing for table data
- `delimited` / `markup` / `serialization` / `escape` — module names are single-purpose verbs/nouns

## Issues

### 🟡 Medium — Terminology split-brain: `Marshaler` vs `Renderer`

The registry layer uses `Marshaler` while every concrete type and the root interface use `Renderer`:

| Layer             | Name                  | Concept                          |
| ----------------- | --------------------- | -------------------------------- |
| Root interface    | `Renderer`            | `Render() (string, error)`       |
| Concrete types    | `JSONTableRenderer` … | implement `Renderer`             |
| Registry (table)  | `TableDataMarshaler`  | `func(w io.Writer, ...) error`   |
| Registry (any)    | `AnyDataMarshaler`    | `func(w io.Writer, ...) error`   |

The distinction is *defensible* (marshal → writer; render → string), but two terms for "produce output" is a mild split-brain. Consider unifying to `TableDataRenderer`/`AnyDataRenderer` for consistency — **low priority**, and changing exported registry types touches ADR 006 (API stability).

### 🔵 Low — `Get`-prefix accessors (27 functions)

- **nom accessor interfaces** (`GetEventType`, `GetDuration`, `GetError`, `GetWorkflowID`, ...) — the Get-prefix is baked into the type-assertion accessor pattern. Go style discourages `Get`, but this is consistent within `nom/` and changing it is a broad, breaking refactor. **Accept — note for future API revision.**
- **`TableData.GetHeaders/GetRows/GetFooter`** — these duplicate the already-exported fields `Headers`/`Rows`/`Footer`. Callers could use `data.Headers` directly. However they ARE used in production (`table/table.go:149`) and are API-frozen (ADR 006). **Accept.**
- **`nom.GetOperationSymbol(operationType)`** — top-level fn; could be `OperationSymbol`. Trivial. **Note.**

### 🔵 Low — `HandleError` (examples/shared/shared.go:12)

Vague verb "Handle". It panics on error. Better: `Must(err)` or `PanicOnError`. Example code only — not in the public library API.

## Fix Recommendations

| # | Issue                      | Priority | Action          | Breaking? |
| - | -------------------------- | -------- | --------------- | --------- |
| 1 | Marshaler vs Renderer      | Low      | Rename in next major version (post-v1) | Yes |
| 2 | Get-prefix in nom          | Low      | Accept (interface-driven) | — |
| 3 | TableData Get* redundancy  | Low      | Accept (API-frozen) | — |
| 4 | `HandleError` in examples  | Trivial  | Rename to `Must` | No (examples) |

No fixes executed: all substantive findings are either API-frozen (ADR 006) or low-value style preferences. The naming is already at a high standard.
