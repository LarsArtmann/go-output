# ADR 007: nom/ Composition via Root Graph Types

**Date:** 2026-06-18
**Status:** ACCEPTED (partially implemented)

## Context

The `nom/` module reinvented root's graph concepts with inferior types:

| nom/ (before)                                                          | root/ (already existed)                          | Problem                                                                |
| ---------------------------------------------------------------------- | ------------------------------------------------ | ---------------------------------------------------------------------- |
| `ActivityNode` + `ActivityDisplayState` (two types, duplicated fields) | `output.GraphNode` (one type)                    | Split-brain: `SyncActivityTimingToTree()` required before every render |
| `ActivityID string` (plain)                                            | `output.GraphNodeID` (branded phantom type)      | No type safety                                                         |
| `InlineRenderer.Render()` (void)                                       | `output.Renderer.Render() (string, error)`       | Contract violation (split-brain M4)                                    |
| No diagram export possible                                             | `output.GraphRenderer` (DOT/Mermaid/D2/PlantUML) | Entire feature class blocked                                           |

## Decision

**Make nom/ import and embed root's graph types.**

1. **`nom.Activity`** embeds `output.GraphNode` — unified source of truth carrying structural identity (ID, Label, Shape, Style) + temporal domain fields (Status, StartTime, EndTime, EstimatedTime, Err).

2. **`ActivityStatus.NodeShape()` / `.GraphStyle()`** — mappers from domain status to root's visual types, enabling the SAME status to drive both terminal lipgloss styling AND diagram export.

3. **`ActivityStore`** — map-backed store with `Nodes() []output.GraphNode` / `Edges() []output.GraphEdge` projections. Any `output.GraphRenderer` can consume live progress state.

4. **`ActivityNode`** (dependency tree) embeds `Activity` instead of `DisplayState` — tree nodes now carry GraphNode data.

5. **`MultiSubscriber`** — `io.MultiWriter`-style fanout for `EventSubscriber`, enabling logging + TUI + metrics from one event stream.

## Killer Feature Unlocked

```go
// Export current build progress as a DOT diagram — zero new rendering code
dot := graph.NewDOTRenderer()
dot.SetNodes(subscriber.Store().Nodes())
dot.SetEdges(subscriber.Store().Edges())
diagram, _ := dot.Render()
```

## Trade-offs

| Concern                   | Decision    | Rationale                                                                                               |
| ------------------------- | ----------- | ------------------------------------------------------------------------------------------------------- |
| nom/ gains root dep       | Accepted    | Root deps are tiny (branded-id, enum, x/term). No lipgloss/bubbletea/yaml/toml pulled.                  |
| Color model stays split   | Intentional | `GraphStyle` (hex strings) for diagram export; lipgloss `color.Color` for terminal. Different concerns. |
| Bridge sync (syncToStore) | Temporary   | Will be eliminated when subscriber migrates from `ActivityDisplayState` to `Activity` directly.         |

## Implementation Status

| Component                                      | Status                              |
| ---------------------------------------------- | ----------------------------------- |
| `Activity` type (embeds GraphNode)             | ✅ Done                             |
| `ActivityStatus.NodeShape()` / `.GraphStyle()` | ✅ Done                             |
| `ActivityStore` with `Nodes()`/`Edges()`       | ✅ Done                             |
| `ActivityNode` embeds `Activity`               | ✅ Done                             |
| `MultiSubscriber` fanout                       | ✅ Done                             |
| `ActivityStore` wired into subscriber          | ✅ Done (bridge sync)               |
| Diagram export tests                           | ✅ Done (3 tests, race-clean)       |
| `ActivityDisplayState` → `Activity` migration  | 🔲 Future sprint                    |
| `SyncActivityTimingToTree` elimination         | 🔲 Future sprint                    |
| `InlineRenderer.Render()` contract fix (M4)    | ✅ Done — `Draw()` / `RenderString` |
| `tui/` migration to new types                  | 🔲 Future sprint                    |

## Consequences

- nom/ now depends on `github.com/larsartmann/go-output` (root)
- Diagram export of progress state is possible with 3 lines of code
- Multiple subscribers can receive one event stream via `MultiSubscriber`
- The `DisplayState` struct and `SyncActivityTimingToTree` are tech debt to be eliminated
