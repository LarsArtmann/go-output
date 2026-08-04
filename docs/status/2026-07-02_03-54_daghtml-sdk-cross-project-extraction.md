# Status: DAG HTML SDK — Cross-Project Extraction & Integration

**Date:** 2026-07-02 03:54
**Scope:** `go-output/daghtml` SDK + `go-workflow-auditlog` + `samber-do-auditlog`

---


> **✅ Resolved (2026-08-04):**
>
> All go-output-side items done. daghtml included in v0.23.0. AGENTS.md module map updated. CHANGELOG written. daghtml lint clean (0 issues). The samber-do-auditlog completion is external. daghtml is independently publishable (zero-dep).

---

## Executive Summary

Extracted the shared Sugiyama DAG visualization (500-line JS algorithm, CSS theme, SVG rendering) from two independent auditlog projects into a reusable zero-dependency Go SDK module (`go-output/daghtml`). One consumer (`go-workflow-auditlog`) is fully refactored, tested, and committed. The second consumer (`samber-do-auditlog`) is ~70% done — adapter and templ modifications are written but uncommitted, `go.mod` not yet updated, build broken.

---

## a) FULLY DONE

### go-output/daghtml module (COMMITTED — 3 commits)

- **Module created** with zero external dependencies (stdlib only: `html/template`, `encoding/json`, `embed`)
- **Data model**: `DAG{Nodes, Edges}`, `Node{ID, Label, Color, Tooltip, Error}`, `Edge{From, To}`
- **Public API**: `Render()`, `Write()`, `GraphHTML()`, `StyleSheet()`, `Script()` + functional options (`WithTitle`, `WithSubtitle`, `WithContainerID`, `WithDataID`, `WithHeight`, `WithFooter`)
- **Embedded assets**: `graph.css` (dark theme with CSS custom properties, `--transient` variable, `.dag-container` class selectors) and `graph.js` (generic Sugiyama algorithm: Kahn's rank assignment, 4-pass barycenter crossing reduction, median alignment, cubic bezier edges, pan/zoom/touch/click-highlight)
- **17 tests passing** — covers HTML structure, JSON data injection, XSS safety, edge dedup, nil-slice handling, all options, round-trip serialization
- **Wired into build system**: `flake.nix` modules list, `go.work.example`, `.golangci.yml` depguard allow-lists
- **CSS fix committed**: `#graph-container` → `.dag-container` class selector (enables multi-graph support via `WithContainerID`); added missing `--transient` CSS variable
- **JS fix committed**: removed literal `<script` string from JS comment (was causing false-positive XSS fuzz test failures in consumers)

### go-workflow-auditlog (COMMITTED — 1 commit)

- **daghtml_adapter.go**: converts `WorkflowReport.Steps` → `daghtml.DAG` (nodes with status colors, tooltips, error flags; edges from dependency refs)
- **html_render.go**: injects DAG JSON + daghtml JS into the HTML template (added 2 new `%s` verbs)
- **dashboard.js**: replaced 500-line `renderGraph()` with 1-line `initDAGGraph("graph-container", "dag-data")`
- **dashboard.css**: graph controls changed from ID-based (`#graph-zoom-in`) to class-based (`.graph-zoom-in`) selectors
- **Tests fixed and passing** (244 tests): updated `assertHTMLStructure` (3→5 script tags), `TestWriteHTML_GraphFailedNodeDot` (now checks `error:true` in DAG JSON), `assertNoRawScriptInjection` (now checks JSON data blocks specifically, not whole HTML), golden file regenerated

---

## b) PARTIALLY DONE

### samber-do-auditlog (~70% — UNCOMMITTED, BUILD BROKEN)

- **daghtml_adapter.go** — WRITTEN but untracked. Converts `Report.Services` → `daghtml.DAG` using `diagramNodeID()` for node IDs, service type/status color mapping, tooltips with invocation count + build duration + errors
- **html.templ** — MODIFIED but uncommitted. Replaced 306-line `renderGraph()` with placeholder `// DAGHTML_JS_INJECTION_POINT`. Graph controls changed from ID to class selectors. Added `@templ.JSONScript("dag-data", buildDAGHTML(report))` for DAG data injection
- **html.go** — MODIFIED but uncommitted. Uses `bytes.Buffer` + `strings.Replace` to inject daghtml JS at the marker point (workaround for templ v0.3.1020 not evaluating expressions inside `<script>` blocks)
- **go.mod** — NOT YET UPDATED. Missing `daghtml` dependency + `replace` directive. Build fails with `no required module provides package github.com/larsartmann/go-output/daghtml`
- **templ generate** — NOT YET RUN. `html_templ.go` is stale (still has the old renderGraph code)
- **Tests** — NOT YET RUN. Expected to have similar script-tag-count and golden file failures as go-workflow-auditlog had

---

## c) NOT STARTED

- **samber-do-auditlog go.mod**: Add `daghtml` require + replace directive, `go mod tidy`
- **samber-do-auditlog templ generate**: Regenerate `html_templ.go` from modified `html.templ`
- **samber-do-auditlog tests**: Run, fix failures (script tag count, golden file, any XSS checks)
- **samber-do-auditlog commit**: Commit all changes
- **AGENTS.md update**: Module map in go-output AGENTS.md missing `daghtml/` entry
- **Release tagging**: No version tags created for any of the three repos
- **Push**: No pushes to remote for any repo

---

## d) TOTALLY FUCKED UP

Nothing is irrecoverably broken. However:

1. **Stale `dag.templ` LSP errors** — The LSP keeps reporting errors for `daghtml/dag.templ` which was deleted (we switched to `html/template`). The file doesn't exist on disk. This is a phantom diagnostic from the templ LSP server caching a deleted file. Harmless but noisy.

2. **templ limitation workaround** — The `strings.Replace` injection point in samber-do-auditlog's `html.go` is a hack. templ v0.3.1020 doesn't evaluate Go expressions inside `<script>` text nodes. The cleaner approach would be upgrading templ or using `@templ.JSONScript` for the JS too (but that wraps in a script tag, which doesn't work for inline JS that needs to be in the same scope as other dashboard JS).

3. **go-workflow-auditlog go.mod replace hack** — Uses `replace github.com/larsartmann/go-output/daghtml => ../go-output/daghtml` (a local filesystem replace). This is correct for development but won't work for external consumers until daghtml is published with a real version tag.

---

## e) WHAT WE SHOULD IMPROVE

### Architecture

1. **daghtml should support custom node metadata** — Currently `Node.Tooltip` is a single string. For richer dashboards, consider a `Metadata map[string]any` field that the JS could render as a structured tooltip.
2. **Theme customization** — The CSS variables are hardcoded in `graph.css`. Consider a `Theme` option that lets consumers override the palette without post-processing the CSS string.
3. **daghtml could live outside go-output** — It has zero go-output dependencies. It could be a standalone repo (`go-daghtml`) for broader adoption. Keeping it in go-output is fine for now but limits discoverability.

### Testing

4. **No golden file test in daghtml** — Every other renderer module (table, tree, graph, d2, etc.) has `golden_test.go`. daghtml should too.
5. **No browser/VT testing** — The JS algorithm is tested only structurally (string contains). A headless browser test (or at least a Node.js execution test) would verify the layout algorithm actually works.

### Process

6. **No CI integration** — daghtml is in the `flake.nix` modules list, but Go checks aren't in `nix flake check` (sandbox blocks `go mod download`). CI configuration hasn't been updated.
7. **Semver strategy unclear** — go-output is at v0.22.0 but daghtml has no version. Need to decide: does daghtml inherit go-output's version, or version independently like testhelpers?

---

## f) TOP 25 THINGS TO DO NEXT

| #   | Task                                                                     | Repo                 | Impact   | Effort     |
| --- | ------------------------------------------------------------------------ | -------------------- | -------- | ---------- |
| 1   | Add `daghtml` to samber-do-auditlog `go.mod` + `go mod tidy`             | samber-do-auditlog   | Critical | 5min       |
| 2   | Run `templ generate` in samber-do-auditlog                               | samber-do-auditlog   | Critical | 1min       |
| 3   | Fix samber-do-auditlog tests (script count, golden, XSS)                 | samber-do-auditlog   | Critical | 15min      |
| 4   | Commit samber-do-auditlog daghtml refactor                               | samber-do-auditlog   | Critical | 5min       |
| 5   | Add `daghtml/` to go-output AGENTS.md module map                         | go-output            | High     | 5min       |
| 6   | Commit AGENTS.md update                                                  | go-output            | High     | 2min       |
| 7   | Add golden file test to daghtml (`golden_test.go`)                       | go-output            | Medium   | 15min      |
| 8   | Tag go-output v0.23.0 (includes daghtml module)                          | go-output            | High     | 5min       |
| 9   | Push go-output to remote                                                 | go-output            | High     | 1min       |
| 10  | Update go-workflow-auditlog go.mod to `daghtml v0.23.0` (remove replace) | go-workflow-auditlog | High     | 5min       |
| 11  | Tag go-workflow-auditlog release                                         | go-workflow-auditlog | Medium   | 5min       |
| 12  | Push go-workflow-auditlog to remote                                      | go-workflow-auditlog | High     | 1min       |
| 13  | Update samber-do-auditlog go.mod to `daghtml v0.23.0` (remove replace)   | samber-do-auditlog   | High     | 5min       |
| 14  | Tag samber-do-auditlog release                                           | samber-do-auditlog   | Medium   | 5min       |
| 15  | Push samber-do-auditlog to remote                                        | samber-do-auditlog   | High     | 1min       |
| 16  | Add `daghtml` to go-output README.md module list                         | go-output            | Low      | 5min       |
| 17  | Add daghtml usage example to go-output examples/                         | go-output            | Low      | 15min      |
| 18  | Consider `Theme` option for daghtml CSS customization                    | go-output            | Low      | 30min      |
| 19  | Consider `Metadata map[string]any` on daghtml Node                       | go-output            | Low      | 20min      |
| 20  | Add `--lazy`, `--eager`, `--alias` CSS vars to daghtml graph.css         | go-output            | Low      | 5min       |
| 21  | Verify `nix run .#build` builds daghtml                                  | go-output            | Medium   | 5min       |
| 22  | Verify `nix run .#test` tests daghtml                                    | go-output            | Medium   | 5min       |
| 23  | Verify `nix run .#lint` lints daghtml                                    | go-output            | Medium   | 5min       |
| 24  | Update go-output CHANGELOG.md with daghtml addition                      | go-output            | Low      | 5min       |
| 25  | Evaluate whether daghtml should be a standalone repo                     | go-output            | Low      | Discussion |

---

## g) TOP QUESTION I CANNOT FIGURE OUT MYSELF

**#1: Versioning strategy for daghtml.**

daghtml is zero-dependency (like testhelpers), but it's not independently versioned in the go.mod — it uses Pattern B (`v0.0.0-00010101000000` + replace). However, both consuming projects (go-workflow-auditlog and samber-do-auditlog) are EXTERNAL to the go-output workspace and currently reference go-output via real published versions (v0.21.0, v0.22.0).

**The question:** When we tag go-output v0.23.0, should daghtml be:

- **(A) Tagged as `go-output/daghtml v0.23.0`** (inheriting the parent version, like d2/graph/table/etc.) — This means consumers do `go get github.com/larsartmann/go-output/daghtml@v0.23.0`. But wait — daghtml has NO go-output dependencies, so it doesn't follow the same lockstep model. Also, the Pattern B replace in go-workflow-auditlog currently points to a local path.
- **(B) Treated like testhelpers** (independently versioned, own tag sequence like `daghtml/v0.1.0`) — More correct architecturally but adds another version coordinate to track.
- **(C) Just part of the go-output monorepo tag** — Consumers use `replace` locally and `go get` the parent for releases.

The consuming projects currently use real versions for all go-output sub-modules. daghtml is the first sub-module that's both zero-dep AND consumed by external repos AND not testhelpers. I cannot determine which pattern is correct without guidance.

---

## Build & Test Summary

| Repo                 | Build       | Tests           | Committed          | Pushed |
| -------------------- | ----------- | --------------- | ------------------ | ------ |
| go-output (daghtml)  | ✅          | ✅ 17/17 pass   | ✅ 3 commits       | ❌     |
| go-workflow-auditlog | ✅          | ✅ 244/244 pass | ✅ 1 commit        | ❌     |
| samber-do-auditlog   | ❌ (go.mod) | ❌ (not run)    | ❌ (3 files dirty) | ❌     |

## File Inventory

### go-output/daghtml/ (11 files, committed)

| File              | Purpose                                                          |
| ----------------- | ---------------------------------------------------------------- |
| `go.mod`          | Zero-dependency module (stdlib only)                             |
| `doc.go`          | Package documentation                                            |
| `types.go`        | `DAG`, `Node`, `Edge` data model                                 |
| `options.go`      | Functional options (`WithTitle`, `WithHeight`, etc.)             |
| `render.go`       | `Render()`, `Write()`, `GraphHTML()`, `StyleSheet()`, `Script()` |
| `assets.go`       | `go:embed` for CSS/JS + JSON serialization                       |
| `graph.css`       | Dark theme CSS with CSS custom properties                        |
| `graph.js`        | Generic Sugiyama DAG layout algorithm (vanilla JS)               |
| `daghtml_test.go` | 16 unit tests                                                    |
| `example_test.go` | Godoc example                                                    |
| `go.sum`          | (empty — no deps)                                                |

### go-workflow-auditlog (4 files changed, committed)

| File                            | Change                                             |
| ------------------------------- | -------------------------------------------------- |
| `daghtml_adapter.go`            | NEW — Report → DAG converter                       |
| `html_render.go`                | Added DAG JSON + daghtml JS injection              |
| `dashboard.js`                  | Replaced 500-line renderGraph with 1-line SDK call |
| `html_test.go` + `fuzz_test.go` | Fixed script-tag count + XSS check                 |

### samber-do-auditlog (3 files changed, UNCOMMITTED)

| File                 | Change                                                  |
| -------------------- | ------------------------------------------------------- |
| `daghtml_adapter.go` | NEW — Report → DAG converter (untracked)                |
| `html.templ`         | Replaced 306-line renderGraph with SDK call placeholder |
| `html.go`            | Added daghtml JS injection via strings.Replace          |

