# ADR 012: CQRS Streaming via Standard Encoders + Registry Rewire

**Date:** 2026-07-06
**Status:** ACCEPTED
**Supersedes:** The old registry-dispatch path (`renderViaRenderer` / `renderDelimitedTable`)

---

## Context

The CQRS architecture (v0.30.0) introduced `WriteXxx(w io.Writer, data)` functions for every format. The initial implementation built full strings in memory, then wrote them — defeating the streaming promise. Session 2 fixed this by rewiring `WriteXxx` to use standard encoders directly (`json.NewEncoder(w).Encode()`, `yaml.NewEncoder(w).Encode()`, etc.).

However, the registry dispatch path (`output.RenderTable(data, FormatJSON, opts)`) still called the OLD renderer structs internally (`renderViaRenderer` → `JSONTableRenderer.Render()` → `MarshalIndent` → `Fprint`). This created a split brain:

- CQRS path: `serialization.WriteJSON(w, data)` → `json.NewEncoder` → streams, adds trailing `\n`
- Registry path: `output.RenderTable(data, FormatJSON, opts)` → `MarshalIndent` → buffers, no trailing `\n`

Two code paths for the same format, producing slightly different output.

## Decision

**Rewire the registry dispatch to call the CQRS streaming functions.**

### What Changed

1. **serialization/json.go**: `renderJSONTable` now calls `WriteJSON(w, data)` instead of `renderViaRenderer(..., NewJSONTableRenderer(), ...)`.
2. **serialization/yaml.go**: Same — `renderYAMLTable` calls `WriteYAML`.
3. **serialization/toml.go**: Same — `renderTOMLTable` calls `WriteTOML`.
4. **serialization/jsonl.go**: Same — `renderJSONLTable` calls `WriteJSONL`.
5. **delimited/csv.go**: Registry closure calls `WriteCSV` instead of `renderDelimitedTable(..., MarshalCSVFromTable, ...)`.
6. **delimited/tsv.go**: Same — calls `WriteTSV`.
7. **serialization/render.go**: Removed `renderViaRenderer` and `dataSetter` (dead code after rewire). Kept `renderTable` (still used by old renderer structs' `Render()` methods).
8. **delimited/csv.go**: Removed `renderDelimitedTable` (dead code).

### What Stayed

- **Old renderer structs** (`JSONTableRenderer`, `YAMLTableRenderer`, etc.) remain for direct callers. They produce output WITHOUT trailing `\n` (via `MarshalIndent`). This is the legacy path. Full deletion deferred to v0.31.0.
- **HTML and AsciiDoc** still buffer — no streaming writer exists for those formats.
- **XML** kept on the old path (`MarshalXMLFromTable`) — no golden test yet to verify output equivalence.

## Consequences

### Positive

- **One code path per format.** Registry and CQRS produce byte-for-byte identical output (proven by `TestCQRS_StreamVsRegistry_JSON/CSV` in `integration/cqrs_test.go`).
- **Streaming everywhere.** The registry path no longer builds full strings in memory before writing.
- **Less code.** Dead helper functions removed. The CQRS `WriteXxx` functions are the single implementation.

### Negative

- **Trailing newline behavior change.** `output.RenderTable(data, FormatJSON, opts)` now outputs a trailing `\n`. Consumers doing exact-output comparison may break. This is documented in CHANGELOG and locked in by CQRS golden files.
- **Two output formats coexist.** The old renderer structs (`NewJSONTableRenderer().Render()`) produce no-`\n` output; the registry path produces with-`\n`. This is temporary — resolved when renderer structs are deleted in v0.31.0.

## Verification

- **Byte-for-byte equivalence**: `TestCQRS_StreamVsRegistry_JSON` and `TestCQRS_StreamVsRegistry_CSV` assert `cqrsBuf.String() == registryBuf.String()`.
- **Golden files**: `testdata/TestGolden_CQRS_JSON.golden` (and YAML/TOML/JSONL/CSV/TSV) lock in the exact streaming output.
- **All 19 modules**: build, test, lint, race, govulncheck — all green.
