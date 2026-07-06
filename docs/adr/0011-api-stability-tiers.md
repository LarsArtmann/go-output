# ADR 011: API Stability Tiers

**Date:** 2026-07-07
**Status:** Accepted

## Context

With v0.30.0 introducing breaking changes, consumers need clear expectations about which APIs are stable and which may evolve. Without tiers, every consumer must treat the entire API as potentially volatile, which discourages adoption.

## Decision

Define three stability tiers for the public API:

| API Surface | Tier | Rationale |
|-------------|------|-----------|
| Root model types (`Table`, `Graph`, `GraphNode`, `GraphEdge`, `TreeNode`, `NodeStyle`, `EdgeStyle`) | **Frozen** | Core data model — breaking these breaks everything |
| Builder API (`NewGraphBuilder`, `AddNode`, `AddEdge`, `Build()`) | **Frozen** | Write-side contract — consumers depend on this to construct data |
| Render functions (`RenderDOT`, `WriteCSV`, `RenderJSON`, etc.) | **Frozen** | Read-side contract — consumers depend on this for output |
| Registry dispatch (`RegisterFormat`, `RegisterTableMarshaler`, `Render`) | **Frozen** | Extensibility contract — sub-modules and consumers depend on this |
| nom/ events + subscriber (`OnEvent`, `NOMSubscriber`, event types) | **Frozen** | Command-side contract — progress visualization depends on this |
| nom/ snapshot model (`ActivitySnapshot`, `SnapshotActivities`) | **Frozen** | Query-side contract — renderers depend on this |
| nom/ themes (`Theme`, `ThemeDracula`, etc.) | **Experimental** | New v0.23 feature, may evolve based on feedback |
| daghtml/ | **Experimental** | New module, still maturing |
| Functional option names per module (`WithDirection`, `WithColorMode`, etc.) | **Frozen** once v0.30.0 ships | Options are part of the public API |

## Consequences

- **Frozen** APIs will not change in a breaking way until v1.0.0. Bug fixes and additive changes (new methods, new options) are allowed.
- **Experimental** APIs may change between minor versions (v0.30.0 → v0.31.0). Consumers should pin versions.
- After v1.0.0, ALL APIs become Frozen. Breaking changes require a major version bump.
