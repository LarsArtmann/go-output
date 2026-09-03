# Full Code Review — In-Progress Status

**Date:** 2026-08-16 08:27
**Scope:** Deep full-code-review of all 19 modules (375 Go files, ~49.1k lines)
**Trigger:** User-requested deep review via full-code-review skill
**State:** INTERRUPTED mid-review (Tier 4% of the Pareto plan). Baseline verified green; root + nom + format-module reviews done; fixes applied; tui/diagrams/infra/report pending.

---

## a) What is FULLY done (and verified)

### Planning & baseline

- [x] Loaded full-code-review skill, architect checklist, pareto-planning skill, html-report-kit guide + template
- [x] Verified git clean before start
- [x] Pareto execution plan written: `docs/planning/2026-08-16_07-50_full-code-review-execution-plan.html` (+ `.d2` source + inlined SVG). Tiers: 1% (root + nom concurrency + green baseline), 4% (nom rest + tui), 20% (diagrams, formats, integration, infra), closeout (re-verify, HTML report, TODO_LIST harvest, commit)
- [x] **Baseline all green**: `nix run .#build` (19 modules), `nix run .#lint` (0 issues everywhere), `nix run .#test` (all ok), `nix run .#test-race` (nom + tui ok)

### Root module review (46 prod + 20 test files) — COMPLETE

Every root production file read and reviewed. Core Invariant verified by grep: **root has ZERO sub-module imports** (only examples/ and integration/ import sub-modules, as designed).

Reviewed: format.go, shape.go, enum.go, color.go, direction.go, envdetect.go, ids.go, renderer.go, registry.go, tabledata.go, table_helpers.go, graph.go, marshal.go, treenode.go, write_rendered.go, render_tabledata.go, streaming.go, marshal_helpers.go, table_builder.go, tree_builder.go, projections.go, tree_to_renderer.go, graph_tabledata.go, doc.go, go.mod + test sampling. Verdict: exceptionally clean; typed errors, branded IDs, generics used well, no split brains found.

### nom/ production review — COMPLETE (all 33 production files)

Personally read: inline_renderer.go (937 L), tree_render.go (926 L), nom_subscriber.go, subscriber_handlers.go, event.go, activity.go, state_accessors.go, tree.go, tree_building.go, tree_critical.go, tree_priority.go, tree_modification.go, status_registry.go, theme.go, activity_snapshot.go, terminal.go, timing_cache.go, timing_cache_persist.go, dag_summary.go, tree_render_layered.go (top half). Agent covered the remaining 11 files + test sampling. Two-mutex lock ordering (renderMu → tickMu) verified consistent; snapshot model honored; sealed event sum type intact.

### Format-module review (8 modules, ~100 files) — COMPLETE via agent + my spot verification

serialization/, markup/, markdown/, delimited/, table/, tree/, escape/, testhelpers/ (+graphtest). I personally verified and confirmed the two most critical agent findings before acting (table footer off-by-one, MermaidText escaping gap — both confirmed real).

### Fixes applied and verified (6 fixes, all module tests green after each)

1. **format.go** — `InvalidFormatError.Error()` now uses shared `EnumErrorMessage` (was hand-rolled 22-line builder; the exact duplication the helper was created to remove)
2. **shape.go** — `InvalidShapeError.Error()` same migration (also gains the nil-Allowed guard it was missing)
3. **nom/dag_summary.go** — removed dead/misleading leaf-count branch (both `if` and `else if` branches did the identical increment; the `!depSet[node.ID]` condition was dead) + removed the now-unused `depSet` map allocation
4. **table/table.go** — **footer styling off-by-one fixed**: `footerRowIndex` was `rowCount+1` but lipgloss 0-based row indices mean the footer lands at `rowCount`; bold/`WithFooterStyle` never applied in real output since v0.x. Now `footerRowIndex = rowCount`, sentinel `-1` (replacing `> 0` guard, also fixes footer-as-first-row edge), initialized in `New()`
5. **table/footer_test.go** — test had encoded the buggy expectation (want 3, actual correct 2); corrected + strengthened with real style assertions (bold at index 2, not at index 1)
6. Root: removed now-unused `strings` import

---

## b) What is PARTIALLY done

### Mermaid escaping hardening (interrupted mid-fix)

- Verified the vulnerability is real: `escape.MermaidText` (escape.go:100-112) escapes only `"`, `[`, `]`, `{`, `}`, `\n`. Mermaid renders labels as HTML (htmlLabels default), so `<img src=x onerror=...>` / `<script>` in node labels, edge labels, AND style values (fill/stroke/color via mermaid.go:149-157) passes the security boundary untouched. Mermaid also decodes `#nnn;` entity codes, so `#60;` can smuggle `<`.
- Confirmed goldens unaffected (labels contain no `<>` or `&`); usage sites mapped (graph/mermaid.go:74,81,149,153,157); test tables located (escape_test.go mermaidTextTestCases, fuzz assertions in escape + graph)
- **NOT yet applied**: the replacer fix (`&`,`<`,`>` at minimum; `#`/`;` decision open — see questions) + test cases

### nom/ tui cross-checks

Agent findings on nom test hygiene (blind sleeps, non-hermetic `NewNOMSubscriber()` reading real `~/.cache/nom-timing.csv`, tautological MarshalError tests in serialization) — read, reported, not fixed.

---

## c) NOT started

- **tui/ module review** (34 files) — untouched
- **d2/, graph/, plantuml/, daghtml/ diagram reviews** — untouched (graph only spot-checked via escape verification)
- **integration/, bdd/, examples/ reviews** — untouched
- **Infra review**: flake.nix, scripts/*.sh, .github/workflows (ci.yml, release.yml), go.mod Pattern B parity across 19 modules, .golangci.yml allow-lists, docs-vs-code drift — untouched
- **Final verification pass** (`nix run .#build && .#test && .#lint`) after all fixes — only per-module tests run so far
- **HTML review report** → `docs/reviews/2026-08-16_08-27_full-code-review.html` — not started
- **TODO_LIST.md harvest** of unfixed findings — not started
- **Commit** — working tree has 5 modified + 3 untracked files (all mine, all self-reviewed)

---

## d) What I totally fucked up / mistakes (full honesty)

1. **Broken CSS written into the plan HTML** — wrote `--text: var(--color-bone;` and `--problem: var(--color-red;` (unclosed parens) on first write. Caught + fixed immediately, but sloppy.
2. **Wrong test invocation** — ran `go test ./nom/` from repo root; Pattern B means root's module doesn't contain nom/. Should have used `nix run .#test` or cd'd into nom/. One wasted round trip; documented in AGENTS.md already, I failed to apply it.
3. **Footer fix initially broke a test** — the test encoded the buggy behavior (want 3). My first strengthened assertion then failed again because bold is only applied when `useColor` (ColorModeNever disables ALL styling including bold) — third attempt (ColorModeAlways) correct. Three round trips for a one-line bug; I should have read the style path fully before asserting.
4. **Interrupted with uncommitted work** — the 8 changed/new files are sitting uncommitted; if the session dies, the review trail is half-preserved (plan committed nowhere).
5. **Agent findings not yet independently verified** beyond the 2 criticals I confirmed (footer, mermaid). Agents can hallucinate line numbers; remaining findings must be spot-verified before fixing (e.g. renderMode race claim, XML footer drop claim).

---

## e) What to improve (process + codebase)

**Process (mine):**

- Commit artifacts earlier and more often (plan should have been committed the moment it was written)
- Verify-before-fix discipline for agent findings (2/2 verified so far — keep the ratio)
- Batch module test runs via flake instead of ad-hoc cd's

**Codebase (from findings so far, highest value first):**

- Security: finish Mermaid escaping; consider AsciiDoc newline escape + Markdown `|`/newline escaping (currently ZERO content escaping in markdown module)
- Data loss: `markup.WriteXML` silently drops `Table.Footer` while `MarshalXMLFromTable` keeps it (CQRS/registry drift)
- Correctness: nom `renderMode` read lock-free in render path (data race if `SetRenderMode` called during render); `ActivityCounts` has only 4 buckets while the status registry is open (custom statuses silently vanish from counts/percent); layered-mode `statusPriority` inverts the package-wide failed-first ordering
- Determinism: tree metadata map iteration (random output order), 2 missing `json.Deterministic(true)` in serialization/json.go, collapsed-category map order in layered mode
- API honesty: nom `Finish(err)` ignores its parameter; `Finish` panics on nil subscriber where `Draw` guards; yaml panic path (`MarshalYAML(chan)` panics instead of erroring)
- Test hygiene: 6 no-op MarshalError tests, 2 blind sleeps, non-hermetic subscriber constructions reading the real HOME cache

---

## f) Next tasks (Pareto order; up to 50)

**1% — the three verified/critical fixes (do first):**

1. Apply MermaidText escaping fix + test cases + fuzz seeds (decision needed on `#`/`;`)
2. Verify + fix nom renderMode data race (read under RLock in render dispatch or make construction-only)
3. Verify + fix markup WriteXML footer data loss (route through shared core with MarshalXMLFromTable)

**4% — nom correctness + determinism:**
4. Verify + fix ActivityCounts custom-status corruption (Other bucket or forbid)
5. Fix layered statusPriority to use Status.Interest() (agree with tree mode, custom statuses free)
6. Fix collapsed-category map iteration order (sort groups)
7. Add nil guard to SetActivityState
8. Add nil-subscriber guard to Finish (match Draw)
9. Resolve Finish(err) dead parameter (render or remove — API decision)
10. Rewrite nom/doc.go concurrency section (describes deleted shared-pointer design)
11. Delete snapshotRoots dead code (tree_accessors.go:78)
12. Fix layered separator misalignment ≥10 layers (drop the len(iota)-1 adjustment)
13. Replace separator rune-sniffing with a VisibleEntry kind field
14. Consolidate FormatDuration vs formatDuration near-duplicates in nom
15. Document/lock GetNode shared-pointer hazard

**tui review (not started):**
16. Review tui/model.go + state.go + messages.go + lifecycle.go
17. Review tui/view.go + render_nom.go + summary.go + colors.go + dot_export.go
18. Review tui/reporter.go
19. Review tui teatest harness (helpers, vt tests)
20. Check TUI consumes fixed nom APIs correctly (post-fix alignment)

**20% — serialization/markup/markdown fixes (verified findings):**
21. Add json.Deterministic(true) to serialization/json.go:32,:64
22. Unify JSONTableRenderer.Render with WriteJSON (byte-identical)
23. Unify the three empty-JSONL behaviors ("" everywhere)
24. Unify TOML empty-table behavior
25. Fix markdown/tree registry-vs-CQRS trailing-newline drift (double \n)
26. Add Markdown cell escaping (`|`, `\n`) — new escape function
27. Add AsciiDoc newline escaping; consider moving replacer into escape/
28. Wrap yaml.Marshal with recover-to-error; flip the panic test
29. Delete or make-real the 6 no-op MarshalError tests
30. Fix xml_test.go:162 boolean operator (`&&` → `||`)
31. Remove csv tableDataWriter alias indirection
32. Stale comment fixes (serialization/render.go:17 "json.NewEncoder")

**20% — remaining module reviews (not started):**
33. Review d2/ (22 files) — domain model + escaping
34. Review graph/ (20 files) — dot/mermaid + enums
35. Review plantuml/ (10 files)
36. Review daghtml/ (8 files) — JS template XSS surface
37. Review integration/ (20 files)
38. Review bdd/ (4 files) + examples/ (8 programs)
39. Review flake.nix + .golangci.yml + go.work.example
40. Review scripts/tag-release.sh + pre-tag-check.sh + setup-workspace.sh
41. Review .github/workflows/ci.yml + release.yml
42. Verify all 19 go.mod Pattern B invariants (sibling v0.0.0 sentinel)
43. Docs-vs-code drift check (AGENTS.md/README/CHANGELOG claims)

**Test hygiene:**
44. Replace blind sleeps in nom/table tests with renderNotify sync + t.Cleanup(Stop)
45. Make progress_events_test constructions hermetic (newTestSubscriber)
46. Strengthen tree_test.go max-height assertion (count lines)
47. Add golden with 2 metadata keys (tree determinism regression guard)
48. Add colored footer golden test (regression guard for fix #4)

**Closeout:**
49. Full `nix run .#build && .#test && .#lint` (+ race for nom/tui) verification pass
50. Write HTML review report to docs/reviews/ + harvest unfixed findings into TODO_LIST.md + commit everything

---

## g) Questions I cannot figure out myself (max 3)

1. **Mermaid escaping scope:** For `MermaidText`, escape `& < >` with HTML entities (matches htmlLabels default rendering, goldens unaffected) — and also neutralize `#`/`;` (mermaid entity-code injection like `#60;`) at the cost of turning literal `#` into `#35;`/`&#35;` in renderers that don't decode entities ("C# dev" → "C&#35; dev" in such renderers)? Full-hardening or minimal `&<>` only?
2. **Fix-on-the-spot boundary:** Keep fixing only _verified critical bugs_ (races, data loss, security) on the spot, and put the behavior-drift family (registry-vs-CQRS newline unification, empty-JSONL/TOML unification, markdown cell escaping) into TODO_LIST for a deliberate v0.38.0 breaking-change release — or unify everything now inside this review?
3. **nom `Finish(err)` API:** render the error (e.g. ✗ line with cause) as the parameter implies, or drop the parameter (breaking, next minor)? doc.go currently teaches `Finish(nil)`.

---

_Point-in-time snapshot. Review continues on instruction._
