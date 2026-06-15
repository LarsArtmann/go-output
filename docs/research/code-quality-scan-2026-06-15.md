# Code Quality Scan — 2026-06-15

Scanned all 15 modules (`.` + 14 sub-modules + testhelpers/graphtest).

## Baseline Summary

| Check        | Result                                                                                       |
| ------------ | -------------------------------------------------------------------------------------------- |
| Build        | ✅ All 16 modules compile (`go build ./...` per module)                                       |
| `go vet`     | ✅ Clean across all modules                                                                  |
| Tests        | ✅ All pass (verified during later skills)                                                    |
| Lint         | ⚠️ 118 issues across 9 modules; 4 modules clean (d2, enum, escape, examples, graphtest)       |
| Duplication  | ✅ art-dupl: "Excellent code health" — 110 groups, all low/idiom (test patterns). 1 actionable |

## Issue Counts by Linter

| Linter                    | Count | Where                              | Severity |
| ------------------------- | ----- | ---------------------------------- | -------- |
| `errcheck`                | 37    | mostly tests (nom, tui, integration) + 6 in nom/inline_renderer.go (production) | High (prod) / Low (test) |
| `wsl_v5`                  | 26    | nom, tui, integration              | Low (style) |
| `depguard`                | 24    | all modules                        | **Config bug** (see #1) |
| `gochecknoglobals`        | 11    | tui/colors.go                      | Medium |
| `err113`                  | 7     | tests (integration, testhelpers, tui) | Low (test) |
| `staticcheck` (SA1019)    | 3     | tui/model_test.go (deprecated `EnsureBuild`) | Low (test) |
| `forcetypeassert`         | 3     | tui/event_sequence_test.go         | Low (test) |
| `embeddedstructfieldcheck`| 2     | nom/activity_display.go, nom/tree.go | Low |
| `dupl`                    | 2     | graph/registry_test.go             | Low (1 clone pair) |
| `predeclared`             | 1     | nom/activity_display_test.go (`copy`) | Low |
| `perfsprint`              | 1     | nom/tree_render.go:168             | Low |
| `nestif`                  | 1     | nom/inline_renderer.go:184         | Medium |

## Sorted Issue List (by impact)

### 🔴 P0 — Root Cause: depguard config not portable (24 false positives)

**File:** `.golangci.yml` (root, inherited by all sub-modules)

The `main` depguard rule (`.golangci.yml:213`) has **no `files:` filter**, so it matches
all files — including `_test.go`. It bans `bytes`, `testing`, `log`, and charmbracelet
sub-packages (`x/ansi`, `x/exp/golden`) that are legitimately needed in tests and in
`nom`/`tui` production code. This produces **24 false-positive** depguard errors across
8 modules.

Sub-modules have no own `.golangci.yml`, so they inherit the root config whose
module-specific allow-lists don't fit them.

**Fix options:**
1. Add `files: ["**/*.go", "!**/*_test.go"]` to the `main` rule so it only governs production code; let `default` (which allows `bytes`+`testing`) govern tests.
2. Add the needed stdlib/charmbracelet imports to `main`'s allow list.
3. (Cleanest long-term) per-module `.golangci.yml` — but high overhead for 15 modules.

### 🟠 P1 — Production code: ignored write errors (nom/inline_renderer.go)

6 `errcheck` violations on `fmt.Fprint`/`Fprintln`/`Fprintf` (lines 152, 163, 166, 180,
186, 188, 192). Writing to the renderer's writer ignores errors. For a renderer whose
entire job is writing output, silently dropping write errors can mask broken pipes.

### 🟠 P1 — Production code: banned imports in nom/tui

- `nom/timing_cache_persist.go:6` imports `log` (banned by depguard `default`)
- `nom/inline_renderer.go:12` imports `github.com/mattn/go-runewidth` (not in allow-list)
- `tui/lifecycle.go:4` imports `log`
- `tui/model.go:7` imports `github.com/charmbracelet/x/ansi`

These are blocked by the same depguard config (#1). Either the imports are legitimate
(and the config should allow them) or they violate intended import hygiene.

### 🟡 P2 — tui/colors.go: 11 package-level globals (gochecknoglobals)

`colorInfo`, `colorTitle`, ... `colorHelpBG` are mutable package-level variables.
Should be constants or encapsulated in a style registry/struct.

### 🟡 P2 — nom/inline_renderer.go:184 nestif (complexity 5)

The `if r.noColor { ... }` block nests too deeply. Extract sub-blocks into named helpers.

### 🟡 P2 — nom/tree_render.go:168 perfsprint

`fmt.Sprintf(" ⬅ depends on %s", ...)` can be string concatenation.

### 🟡 P2 — embeddedstructfieldcheck (nom)

`nom/activity_display.go:28` and `nom/tree.go:19` embed `DisplayState` after regular fields.
Move embedded fields first (Go convention).

### 🔵 P3 — Test code (low value, high churn)

- 31 `errcheck` in tests across nom/tui/integration (unchecked `AddActivity`, `Record`,
  `OnEvent`, `UpdateActivityStatus` return values).
- 7 `err113` dynamic `errors.New(...)` in tests (should use wrapped static errors).
- 3 `forcetypeassert` in tui/event_sequence_test.go (unchecked type assertions).
- 3 `staticcheck` SA1019 deprecated `EnsureBuild()` in tui/model_test.go.
- 26 `wsl_v5` whitespace nits in nom/tui/integration.
- 1 `predeclared`: `copy` shadows builtin in nom/activity_display_test.go:204.

### 🔵 P3 — Duplication (1 actionable clone)

`graph/registry_test.go`: `TestRenderDOTTableData` (lines 20-52) is a near-exact duplicate
of `TestRenderMermaidTableData` (lines 54-86). Extract a table-driven helper parameterized
on `output.Format`.

## Clean Modules ✅

d2, enum, escape, examples, testhelpers/graphtest — zero lint issues.
