# Full Code Review — In-Progress Status (Session 3)

**Date:** 2026-08-16 08:55
**Scope:** Deep full-code-review of all 19 modules (375 Go files, ~49.1k lines) via full-code-review skill
**Trigger:** User-requested: "READ, UNDERSTAND, RESEARCH, REFLECT... Execute and Verify one step at a time... Keep going until done"
**State:** INTERRUPTED mid-tui-review at user request for this status report. All 5 verified-critical fixes from the prior session's question list are DONE and COMMITTED. A 6th fix (TUI mouse-click mapping) is APPLIED BUT UNVERIFIED — **tui currently does not compile** (see b).

---

## a) What is FULLY done (and verified green) — COMMITTED as `c4e6069`

### 1. MermaidText HTML-injection fix (security, verified real)

- `escape/escape.go`: `MermaidText` now escapes `&`→`&amp;`, `<`→`&lt;`, `>`→`&gt;` (Mermaid renders labels as HTML by default — raw `<script>`/`<img onerror>` in node labels, edge labels, AND style values previously crossed the security boundary), PLUS a `mermaidEntityGuard` regex `#([0-9A-Za-z]+;)` → `&#35;$1` that neutralizes Mermaid's `#60;`-style entity-decode smuggling.
- Design constraint honored: hex colors in style directives (`fill:#ff0000` — never semicolon-terminated) and plain hashes ("C# dev") pass through untouched; goldens render byte-identical.
- Tests: 9 new table cases (script tag, img onerror, numeric/named entity smuggle, typed-entity guard, hex color, "C# dev"), fuzz seeds + `<`/`>` assertions in escape AND graph fuzz corpora. Fuzz ran 5M execs clean.
- Verification: escape + graph + integration tests green, examples build.

### 2. nom renderMode data race (verified real)

- Agent claim confirmed in production: `tui/model.go` toggles `SetRenderMode` from the Update loop while the tick loop renders. The 4 render dispatch sites in `nom/tree_render.go` read `dt.renderMode` without `dt.mu`.
- Fix: all 4 dispatch sites now read via the locked accessor `dt.RenderMode()`; `WithRenderMode` option now uses the setter instead of raw field write.
- New `-race` regression test `nom/tree_render_mode_race_test.go` (500 toggles × renders across all 4 entry points) — passes under `-race`.

### 3. markup WriteXML footer loss + CQRS/registry drift (verified real, worse than reported)

- `WriteXML` AND the registry path (`renderXMLTable` → WriteXML) silently dropped `Table.Footer`, while `MarshalXMLFromTable` kept it. Also `WriteHeader` emitted an empty `<headers>` block where Marshal omitted it.
- Fix: `XMLWriter.WriteFooter(footer []string)` gained the footer-row parameter (documented breaking signature change); `WriteHeader` guards empty cols; `MarshalXMLFromTable` now delegates to `WriteXML` for non-nil tables so both paths are byte-identical BY CONSTRUCTION.
- New regression tests: footer-not-dropped, empty-headers omission, and WriteXML≡MarshalXMLFromTable parity. Goldens unaffected (footer-less fixture). All markup tests green.

### 4. nom ActivityCounts loses custom statuses (verified real)

- The open status registry (`RegisterStatus` — documented for "skipped"/"cached") made activities in custom statuses INVISIBLE to counts/Total/CompletionPercent. Even the brute-force `recount` test helper dropped them, so the cache-consistency invariant test could never catch it.
- Fix: new `Other` bucket (+ `SymbolOther "⋯"`) aggregated by `adjustStatusCount`/`recount`; `Total()`/`Summary()`/`SummaryColored()` include it. Custom→core transitions decrement correctly.
- New tests: custom-status-not-lost (registry-backed), plus nil guard for `SetActivityState`. All nom tests green.

### 5. nom `SetActivityState` + `Finish` nil guards (verified real)

- `SetActivityState(id, nil)` panicked on `activity.Status`; `Finish` panicked on nil subscriber where `Draw` guards. Both now no-op defensively.
- `Finish(workflowErr)` dead parameter resolved per the documented contract (twice documented: caller owns the post-run summary; integration test asserts no completion line): signature kept, doc comment now states the parameter contract explicitly. Signature cleanup queued for v0.38.0 TODO harvest.
- New test: `TestInlineRenderer_Finish_NilSubscriber`.

### Decisions made autonomously (prior session's 3 open questions)

1. **Mermaid escaping scope**: full hardening — `&<>` + targeted `#entity;` guard (NOT blanket `#` escaping, which would break hex colors in styles).
2. **Fix-on-the-spot boundary**: verified bugs (security/data-loss/races/wrong-output) fixed now; intentional behavior-drift unification (registry-vs-CQRS empty payloads, Finish signature removal) → TODO_LIST for v0.38.0.
3. **Finish(err)**: keep signature + honest docs now; removal queued for v0.38.0.

---

## b) What is PARTIALLY done (UNCOMMITTED, on top of `c4e6069`)

### 6. TUI mouse-click mapping bug — FIX APPLIED, **BUILD CURRENTLY BROKEN**

**Verified real (third "test-encodes-bug" case this review):** the mouse click handler mapped `mouse.Y` via `treeStartLine(2) + chromeLinesAboveTree(5)` = first tree line at screen line 7. Empirical probe of the ACTUAL rendered layout: tree starts at line **3** (no message) / line **5** (with current message). Every NOM-mode click therefore selected the node 2–4 lines BELOW the one clicked (or nothing). All existing mouse tests passed only because they computed click positions with the same wrong formula.

Applied (tui/model.go, tui/render_nom.go, tui/model_mouse_test.go — modified, uncommitted):
- `renderNOMStyle` now computes `treeStartLine` from the real layout (`nomTreeStartLine()`: title block 2 + section gap 1 + optional message 2).
- `chromeLinesAboveTree` const deleted; click mapping is `mouse.Y - m.treeStartLine`.
- `model_mouse_test.go` restructured to render first (populating visibleEntries AND the true mapping), then click; right-click test de-hardcoded.

**⚠ NOT DONE — this is why tui does not compile:**
- `tui/layered_test.go:45,62` still reference the deleted `chromeLinesAboveTree` → `go vet` fails.
- Remediation (next session, ~5 min): drop the `+ chromeLinesAboveTree` from both clickY lines (render first, as in model_mouse_test.go), then `GOEXPERIMENT=jsonv2 go test ./...` in tui/, then add a layout-pinning regression test asserting `lines[treeStartLine]` is actually the first tree row (render → split → assert), so the mapping can never silently drift again.
- `git status`: 3 modified files, uncommitted. Note: the auto-git daemon may commit this BROKEN state — if a commit lands before the fix, follow-up commit must repair `layered_test.go`.

---

## c) What is NOT started (remaining Pareto tiers, unchanged from prior plan)

- **Diagram module reviews**: d2/ (22 files), graph/ remainder (mermaid path spot-checked via escape fix), plantuml/ (10), daghtml/ (8, JS template XSS surface)
- **Test-module reviews**: integration/ (partially exercised via fixes), bdd/, examples/
- **Infra review**: flake.nix, scripts/*.sh, .github/workflows (ci.yml, release.yml), 19× go.mod Pattern B parity, .golangci.yml allow-lists, docs-vs-code drift (AGENTS.md/README/CHANGELOG)
- **Remaining verified-findings fixes** (from prior session's list, all confirmed by earlier verification):
  - serialization/json.go:32,64 — add `json.Deterministic(true)` (2 sites)
  - markdown — zero content escaping (`|`, `\n` break tables) — new escape path
  - tree/tree.go:115-119 — metadata map iteration → sort keys (nondeterministic output)
  - nom/tree_render_layered.go — statusPriority inversion (failed should sort first), collapsed-category map order, ≥10-layer separator misalignment, separator rune-sniffing → VisibleEntry kind
  - serialization/yaml.go — panic path (chan input) → recover-to-error
  - nom/doc.go concurrency section describes the deleted shared-pointer design
  - nom/tree_accessors.go:78 `snapshotRoots` dead code
  - nom FormatDuration vs formatDuration near-duplicates
- **Test hygiene**: 6 no-op MarshalError tests (serialization), xml_test.go:162 `&&`→`||`, blind sleeps (nom/inline_renderer_test.go:226, timing_cache_test.go:216), non-hermetic `NewNOMSubscriber` in progress_events_test.go, csv tableDataWriter alias
- **Closeout**: final `nix run .#build && .#test && .#lint` (+ race nom/tui), HTML report → `docs/reviews/2026-08-16_<HH-MM>_full-code-review.html` (Bauhaus Dark kit), TODO_LIST.md harvest, final commit

---

## d) What I totally fucked up / mistakes

1. **Left the tree broken mid-edit** — deleted `chromeLinesAboveTree` from production but hadn't updated `layered_test.go` when interrupted. The commit boundary should have been after ALL references compile. One `go vet` would have caught it before the interrupt; I ran the probe test instead of vetting first.
2. **Probe test left in tree?** — No: `mousemap_probe_test.go` was created, run, and `rm`'d in the same command. Verified gone.
3. **Three round trips on the escape regex** — first attempt used RE2-unsupported lookahead `(?=...)` (panic at init), second edit failed on a whitespace mismatch, and I introduced a "the the" comment duplication while fixing grammar. All caught and fixed, but sloppy for a security fix; the fuzz+table suite is what actually proved correctness.
4. **Repeated pattern worth internalizing**: this review has now found THREE cases where tests encoded the bug they should have caught (table footer index, mouse click formula, counts recount helper). When a test derives its expectation from the same constant/formula as the implementation, it verifies nothing.

---

## e) What to improve (highest value first)

1. **Finish the tui mouse fix** (b) — compile, test, layout-pinning regression test, commit.
2. **Keep the verify-empirically discipline** — every remaining agent claim gets a probe/render/test BEFORE fixing (6/6 verified-real so far this session: mermaid, renderMode race, XML footer, counts, nil guards, mouse mapping).
3. **Consider a "layout contract" test for the TUI**: render → assert structural line positions (title/message/tree/summary) so treeStartLine-style drift is structurally impossible.
4. **Remaining module reviews** are mostly volume — delegate d2/plantuml/daghtml to an agent with verify-before-fix, personally review daghtml's JS template injection surface.

---

## f) Next tasks (Pareto order)

1. Fix `layered_test.go` chromeLinesAboveTree refs → tui tests green → layout-pinning regression test → commit mouse-mapping fix
2. Review remainder of tui/ tests (teatest harness already read; regression_test.go, event_*_test.go skim for encoding-bugs)
3. Agent-review d2/ + plantuml/ + daghtml/ (verify before fix; personally audit daghtml JS/CSS injection + template typing)
4. Review integration/, bdd/, examples/
5. Infra review (flake, scripts, CI, go.mod parity, docs drift)
6. Remaining verified-findings batch (see c for the full list)
7. Test-hygiene batch
8. Final verification: `nix run .#build && .#test && .#lint && .#test-race`
9. HTML review report → docs/reviews/ + TODO_LIST.md harvest + final commit

---

## g) Questions I cannot figure out myself

None blocking. The three prior questions were resolved autonomously (see a). One judgment call for v0.38.0 harvest (not blocking): whether `Finish(workflowErr)` should drop the parameter or finally render the error — current documented contract says caller owns it; I lean drop-the-parameter.

---

*Point-in-time snapshot. Review resumes at task f.1 on instruction.*
