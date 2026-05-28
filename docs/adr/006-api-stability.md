# ADR: Pre-v1 API Stability Guarantees

**Date:** 2026-05-28
**Status:** ACCEPTED & IMPLEMENTED

## Context

go-output has 13 modules, 16 output formats, and ~228 exported symbols. The library is at v0.6.0. Before tagging v1.0.0, the public API must be audited and locked in.

## Decision

All exported symbols across all modules are now **frozen**. The following guarantees apply:

### Frozen Interfaces

- `Renderer` — `Render() (string, error)`
- `TableRenderer` — `SetHeaders([]string)`, `AddRow([]string)`, plus `Renderer`
- `TreeOutputRenderer` — `SetRoot(*TreeNode)`, plus `Renderer`
- `GraphRenderer` — `SetNodes([]GraphNode)`, `SetEdges([]GraphEdge)`, plus `Renderer`
- `StreamingRenderer` — `Stream(io.Writer) error`, plus `Renderer`

### Frozen Core Types

- `Format` (16 constants), `Shape` (3 constants), `ColorMode` (3 constants)
- `TableData`, `TreeNode`, `GraphNode`, `GraphEdge`, `GraphRendererMixin`
- All branded ID types (`D2NodeID`, `TreeNodeID`, `GraphNodeID`, etc.)

### Non-Breaking Changes Allowed

- New `Format` constants
- New `Shape` constants
- New methods on existing types
- New `RenderOptions` fields
- New sub-modules

### API Issues Found and Fixed

1. **Capability matrix incomplete**: `FormatD2`, `FormatMermaid`, `FormatDOT`, `FormatPlantUML` were missing `ShapeTree` despite having `*FromTree()` conversion functions. `FormatTOML` was missing `ShapeGraph` despite having `TOMLGraphRenderer`. Fixed in `shape.go`.

2. **`RenderOptions.GraphID` is dead code**: No registered `TableDataMarshaler` reads this field. DOT returns `UnsupportedFormatError` from `RenderTableData`. Kept as a no-op field for future use — not breaking to leave, would be breaking to remove.

3. **`MarshalTSV(data any)`**: Uses type switch over `any` instead of concrete `[][]string`. This is intentional — the function also handles `[]string` (single row). A generic overload would add complexity for minimal gain.

## Consequences

- Users can rely on all current APIs remaining stable through v1.x
- The 5 interfaces above form the extension points for new formats
- `internal/` packages remain exempt from stability guarantees
