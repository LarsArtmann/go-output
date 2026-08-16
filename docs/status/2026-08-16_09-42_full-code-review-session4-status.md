# Full Code Review — Session 4 Status

**Date:** 2026-08-16 09:42
**Scope:** Continuation of the deep full-code review (session 4 of the interrupted review; plan: `docs/planning/2026-08-16_09-00_SUPERB-code-review-completion-plan.md`)
**Trigger:** User work order: READ/UNDERSTAND/RESEARCH/REFLECT → Pareto breakdown (1%/4%/20%) → comprehensive 2-level plan with tables → execute + verify one step at a time → no VERSCHLIMMBESSER → plan .md with mermaid graph → commit (+push at closeout)
**State:** Tier 1% + Tier 4% essentially COMPLETE. 4 commits this session, all green at commit time. Tier 20% reviews partially done (d2/plantuml/graph via agent + fixes applied; daghtml/integration/bdd/examples/infra NOT started). Closeout NOT started.

---

## a) What is FULLY done (all committed, module tests green at commit time)

### Session commits (this session, `b28aa31` → `ff22ec3` on top of `e58c8bf`)

1. **`b28aa31` — tui build repair + layout contract test.** Fixed `layered_test.go` (still referencing the deleted `chromeLinesAboveTree` — the daemon had committed the broken state during the last interrupt); restructured both layered click tests to render first (real `treeStartLine`); restored the toggle-off preset my earlier restructure dropped; added `TestNOMStyle_TreeStartLineMatchesRenderedLayout` which renders with/without a message and asserts the line at `treeStartLine` IS the first tree row — the mapping can never silently drift again.

2. **Plan artifact committed** — `docs/planning/2026-08-16_09-00_SUPERB-code-review-completion-plan.md` (part of `4ab5e1a`): Pareto tiers, Level-1 task table (17 tasks, 10-30 min), Level-2 subtask table (40 subtasks ≤12 min), mermaid.js execution graph, rules of engagement.

3. **`4ab5e1a` — escaping/correctness batch** (biggest commit this session; every fix verified at the exact lines first):
   - **T1 Markdown cell escaping**: new `escape.MarkdownCell` (`\`→`\\`, `|`→`\|`, `\n`→`<br>`, `\r` dropped); `MarkdownTable.Render` escapes headers/rows/footer up front so column widths are computed on escaped text (padding stays aligned); 2 new tests lock row integrity (no extra rows, 3 unescaped pipes/line) and width consistency. Module has no goldens → zero drift.
   - **T2 JSON determinism**: `json.Deterministic(true)` added to `MarshalJSON` and `JSONWriter.Encode` (serialization/json.go) — map keys now sort on every JSON path.
   - **Agent-verified findings fixed** (agent reviewed d2/plantuml/graph file-by-file; I verified each claim at the cited lines — all 13 real; 9 fixed, 4 harvested):
     - **P1 (High)**: PlantUML emitted `start --> : flows to end` (label before target — target swallowed, edge dangling). Now `from --> to : label`. Golden that had pinned the bug regenerated; diff verified to be exactly this fix.
     - **P2 (High)**: PlantUML IDs could inject raw newlines/directives (`a\n@enduml`) — last renderer with that vector. New `escape.PlantUMLID` (slug → allowlist [A-Za-z0-9_] → "node" fallback); tree-node IDs now sanitized too (they bypassed sanitization entirely via `plantUMLTreeNodeID`).
     - **D1/D2 (Medium)**: d2 SQL `Constraint` now quoted (was raw — block-close + style injection); `Direction`/`NodeShape`/`ArrowType` gated by `IsValid()` at render.
     - **D3 (Medium)**: `d2.Write` renders from a shallow copy — options no longer permanently mutate the caller's diagram.
     - **D4/D5 (Low)**: `\r` added to the shared D2/DOT replacer; empty d2 tree-node ID falls back to "node".
     - **G1/G2 (Medium)**: DOT graph IDs quoted via `%q` when not plain alphanumeric (hostile `x\nlabel="injected"` stays one line, inert); invalid `RankDir`/`SplineStyle`/`LineStyle` omitted at render. 2 injection regression tests added (my first assertion for G1 was naive — "contains 'injected'" also matches the safe quoted form; corrected to assert the one-line-header invariant).
   - P3 (`WithDiagramType` dead option — API removal decision), D6/G3 (registry-vs-CQRS trailing-newline family) → harvested for v0.38.0 (behavior-drift family).

4. **`ff22ec3` — T3/T4/T5 verified findings**:
   - **T3 (verified)**: layered `statusPriority` ranked running > pending > failed > completed while tree mode orders by registry `Interest()` (failed first) — failed sank to 3rd in layered view and custom statuses sank below everything. Layered now sorts by `Status.Interest()`; hand-rolled ladder deleted; collapsed-category groups sorted (map-order nondeterminism).
   - **T4**: tree metadata annotations sort keys; regression test renders identical 3-key metadata 20×, asserts byte-stable sorted output.
   - **T5**: `MarshalYAML` recovers encoder panics (chan/func inputs) into errors; the test asserting the PANIC flipped to assert the error.

### Cumulative review state

All 6 session-3 criticals + 13 session-4 verified fixes committed. Modules now personally reviewed: root, nom (all 33 prod files), tui (all 12 prod files), escape, markdown (render path), serialization (json/yaml), markup (xml), table (footer path), plantuml, d2, graph (dot+mermaid via verification work).

---

## b) What is PARTIALLY done

- **Agent review of d2/plantuml/graph**: complete and verified/fixed/harvested (see above). But my own file-by-file pass over d2/ (22 files) and graph/ (20 files) beyond the fixed sites has NOT happened — the agent covered them; spot-verification was targeted at findings only.
- **Test hygiene (plan T6/T7/T8)**: NOT started — 6 no-op MarshalError tests, xml_test.go:162 boolean operator, 2 blind sleeps, non-hermetic progress_events tests, nom/doc.go stale concurrency section, `snapshotRoots` dead code, FormatDuration near-dupes, colored-footer golden.

---

## c) What is NOT started

- **daghtml/ XSS/template audit (T10)** — highest residual security surface (JS/HTML template injection) — personal audit pending.
- **integration/ + bdd/ + examples/ review (T11)**.
- **Infra review (T12/T13)**: 19× go.mod Pattern B parity, .golangci allow-lists, scripts, CI workflows, docs drift.
- **Closeout (T14-T17)**: full `nix run .#build && .#test && .#lint && .#test-race-all` gate, HTML review report → `docs/reviews/`, TODO_LIST.md harvest (accumulated: P3, D6/G3 newline family, Finish(err) signature, separator ≥10-layer alignment, separator rune-sniffing → VisibleEntry kind, v0.38.0 batch), AGENTS.md pattern updates, final commit + **push** (user explicitly ordered push this session).

---

## d) What I totally fucked up / mistakes

1. **Daemon committed broken tui mid-edit (last interrupt)** — repaired first thing this session; lesson applied: **finished the in-flight nom edit + tests + commit BEFORE starting this status report**, so the tree is green and committed, not a landmine.
2. **Python-heredoc edit failures ×5** — heredoc mangled `\n`/`\\` escapes in Go test literals (PlantUML test got a literal newline; markdown backslash assertion collapsed; two replaces silently didn't match). Each burned a round trip and one left a broken build caught by `go test`. Correct move (used from G1 onward): read the file, edit with the Edit tool, never embed Go escape-heavy literals through shell heredocs.
3. **Naive G1 test assertion** — `Contains(got, "injected")` also matches the SAFE quoted form. Rewrote to the real invariant (header stays one line, quoted form). Also double-escaped DOT IDs (`%q` + `escape.DOT`) before catching it via the test.
4. **Placeholder-edit slip in T3** — one multiedit block landed a no-op loop (`_ = cat; _ = count`) instead of the sorted version; caught immediately and replaced, but it's exactly the VERSCHLIMMBESSER risk the user threatened my balls over. Tests + review-of-own-diff before commit is what caught it.
5. **Earlier claim vs reality**: my "aligned widths" test invariant for markdown was never true of the format (separator rows legitimately differ in raw length); replaced with the actual invariants (header≡row length, padding to escaped width, separator = width+1 dashes).

---

## e) What to improve

1. **Never heredoc Go string literals** — always Edit tool after View. (Systemic this session; 5 incidents.)
2. **Agent + verify split works** — 13/13 agent findings verified real at cited lines; keep the pattern for daghtml/integration/bdd/examples (T10/T11).
3. **Commit-per-green-batch held** — 4 clean commits, each with module tests green; continue.
4. **Remaining high-value**: daghtml audit before hygiene work — security surface beats cosmetics.

---

## f) Next tasks (Pareto order, up to 50)

1. T10 daghtml/ audit (8 files): template.CSS/JS/HTML typing, dagToJSON data flow, XSS surface — personal, ~20 min
2. T11 integration/ (20 files) skim: coverage gaps + tests-encoding-bugs
3. T11 bdd/ (4) + examples/ (8 programs) skim
4. T12a verify 19× go.mod Pattern B sentinels (scriptable in one pass)
5. T12b scripts/*.sh + .github/workflows/{ci,release}.yml review
6. T12c .golangci depguard allow-list audit vs actual imports
7. T13 AGENTS.md/README/CHANGELOG claim spot-check
8. T6a xml_test.go:162 `&&`→`||` (5 min)
9. T6b rewrite 6 no-op MarshalError tests to assert real errors
10. T7a replace blind sleep nom/inline_renderer_test.go:226 with renderNotify + t.Cleanup(Stop)
11. T7b replace blind sleep nom/timing_cache_test.go:216
12. T7c hermetic newTestSubscriber in progress_events_test.go
13. T7d colored footer golden (regression guard for the v0.36-era fix)
14. T8a rewrite nom/doc.go concurrency section (snapshot model)
15. T8b delete `snapshotRoots` dead code (tree_accessors.go:78)
16. T8c consolidate FormatDuration/formatDuration near-duplicates
17. T14 full verification: `nix run .#build && .#test && .#lint && .#test-race-all`
18. T15 HTML review report → `docs/reviews/2026-08-16_<HH-MM>_full-code-review.html` (Bauhaus Dark: stat cards, issue cards incl. all fixed+harvested findings, badge table)
19. T16 TODO_LIST.md harvest: v0.38.0 behavior-drift family (P3 WithDiagramType, D6/G3 newline unification, Finish(err) signature), layered separator ≥10 alignment, separator rune-sniffing→kind field, plus any T10-T13 findings
20. T16b AGENTS.md updates: new invariants (escape.MarkdownCell/PlantUMLID, DOT graphID quoting, layered Interest ordering, ActivityCounts Other bucket, tui layout pin test)
21. T17 final commit + PUSH (user-ordered)

---

## g) Questions I cannot figure out myself

None blocking. Standing judgment call (unchanged, non-blocking): v0.38.0 should decide `Finish(err)` drop-parameter vs render, and whether the registry-vs-CQRS newline unification lands as documented breaking change — both already queued for the TODO harvest with my recommendation noted.

---

*Point-in-time snapshot. Review resumes at task f.1 on instruction. Tree is green and fully committed.*
