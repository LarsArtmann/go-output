# ADR 010: True DAG Topology — Dependency-Aware Rendering

**Date:** 2026-07-02
**Status:** ACCEPTED & IMPLEMENTED
**Deciders:** Lars Artmann

## Context

The NOM dependency tree originally rendered activities as a simple tree — each node appeared once under a single parent, with no awareness of cross-cutting dependencies. This produced several problems:

1. **Misleading structure**: A node with 3 dependencies appeared as a child of just one, hiding the real blocking relationships. Users couldn't see why a task was still pending when its other (non-display) deps hadn't completed.

2. **No DAG-level intelligence**: Critical-path analysis, blockage detection, and parallelism metering all require the full dependency graph. A tree-only structure made these features impossible.

3. **nix-output-monitor mismatch**: NOM (the reference implementation) tracks true DAG topology and uses it for layered rendering. Our tree-only model diverged from this design.

## Decision

Store the full DAG topology (all dependency edges) on `ActivityNode.Deps` and make it the structural source of truth. The display tree is computed from `Deps` at `Build()` time.

### Key design choices:

- **`Deps []ActivityID` is the source of truth**: Each entry means "this node depends on depID". `AddActivity(id, deps)` records them directly. The display parent is selected as the deepest dep (or first root for nodes with no deps).
- **`AllNodes() []DAGNode`**: Returns read-only snapshots of all nodes with their full dep sets, for external consumers (DOT export, analysis).
- **`computeCriticalPath()`**: Fixpoint longest-path algorithm over the DAG using `Deps`. Back-tracks from the maximum-total node to find all nodes on the critical path.
- **Layered mode uses depth from DAG**: `collectLayeredEntries` groups by `node.Depth` which is computed from the DAG topology, not the display tree depth.
- **Placeholder skipping**: Nodes that only exist as dependency targets (never registered as real activities) are skipped in rendering — they have no snapshot label.

## Consequences

**Positive:**

- Critical-path analysis, parallelism metering, and blockage detection all "just work" because they have access to the full graph.
- Layered mode renders true parallelism structure (siblings by depth, not by display parent).
- DOT export via `AllNodes()` produces the real DAG, not the display approximation.

**Negative:**

- `ActivityNode` is heavier: carries both `Deps` (DAG edges) and `Parent`/`Children` (display tree). The display tree is derived at Build() time.
- Potential confusion: `Children` is NOT the inverse of `Deps`. A child in the display tree may depend on other nodes that are not its display parent.

**Neutral:**

- `ExtraDeps()` exposes non-parent deps for the optional "↳ dep" sub-line, making cross-edges visible without cluttering the default view.
