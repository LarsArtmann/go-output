# Split-Brain Sprint — Brutal Status Report

**Created:** 2026-06-17 23:22
**Sprint:** Split-Brain Elimination (started from `SPLIT-BRAIN.html` audit)

---

## a) FULLY DONE (verified, working, tested)

| Fix | Description                                                        | Verification                |
| --- | ------------------------------------------------------------------ | --------------------------- |
| C1  | `nom.TreeNode` → `nom.ActivityNode` (73 refs, 12 files)            | ✅ Build + test all modules |
| C3  | `tui.TimingFormat` → `timingFormatWithIcon` (unexported)           | ✅ Build + test tui         |
| C4  | Deleted `graphRenderer` redeclaration in serialization tests       | ✅ Build + test             |
| C5  | Deleted `renderer` redeclaration in integration tests              | ✅ Build + test             |
| M1  | Deleted dead `nom.ColorWarning` (zero callers)                     | ✅ Build + test nom         |
| M3  | Replaced hardcoded `"No activities to display"` literal            | ✅ Build + test             |
| m4  | Added missing `Style` field to GraphEdge in FORMAT_ARCHITECTURE.md | ✅ Docs only                |
| m5  | Fixed `GetWorkflowID()` return type (`string` → `WorkflowID`)      | ✅ Build + test             |
| m1  | Added cross-reference comments for `"unknown"` sentinel            | ✅ Build + test             |

**9 of 20 findings fully resolved.**

---

## b) PARTIALLY DONE (papered over, incomplete, or half-assed)

### C2: `ProgressModel` state duplication — DOCUMENTED, NOT FIXED

**What was claimed:** "Removed dead `timingCache` field, documented remaining fields"
**What actually happened:** Only removed `timingCache` (which was never read). The real split-brain — `activities` map as a deep-copy cache of the subscriber's state — was preserved with a comment. The `syncNOMSubscriber()` method still manually copies state every tick.

**The better fix:** Add `GetActivityCounts()` to `NOMStyleSubscriber` (returns counts, not the full map). Then delete `activities` field entirely. The only production reader is `view.go:331` which only reads `.Status` — it doesn't need the full map.

### M2: Color detection alignment — ALIGNED ENV VARS, MISSED TERMINAL CHECK

**What was claimed:** "Aligned divergent color detection logic"
**What actually happened:** Added `TERM=dumb` to root, added 4 CI vars to nom. But root's `ShouldColor()` also checks `isStdoutTerminal()` — nom's `detectNoColor()` does NOT check if stdout is a terminal at all. So nom will emit color codes even when piped to a file (as long as no env vars are set).

**The real fix:** Use `charmbracelet/x/term` (same org as the already-imported `charm.land/x/ansi`) to replace both hand-rolled implementations.

### m2: Event literal cleanup — LESS THAN 50% COMPLETE

**What was claimed:** "Replaced 23 bare event-string literals with constants"
**What actually happened:** Replaced literals in `tui/`, `integration/`, `examples/`. **11 bare literals remain in `nom/subscriber_test.go`** — the very module that defines the constants!

### M9: delimitedWriter drift — NEVER ADDRESSED

**What was claimed:** Dismissed as "idiomatic Go consumer-side interface narrowing"
**What actually happened:** The example's `delimitedWriter` has `WriteFooter` — a method both real writers implement. The real `tableDataWriter` interface lacks it. The example is more correct than the real code.

---

## c) NOT STARTED

| ID    | Issue                                               | Why deferred               |
| ----- | --------------------------------------------------- | -------------------------- |
| M4    | Rename `Render()` methods (incompatible signatures) | API break, TODO added      |
| M5    | Rename `ShapeBox` prefix collision                  | API break, TODO added      |
| M6/M7 | Bridge direction enums                              | New API design, TODO added |
| M8    | Align style struct field names                      | API break, TODO added      |
| m6    | Move branded IDs to d2 module                       | API change, TODO added     |

---

## d) TOTALLY FUCKED UP

### 1. AGENTS.md NOT updated

The memory maintenance protocol is mandatory: "Update project AGENTS.md PROACTIVELY when you learn: architecture decisions, conventions, gotchas." The `AGENTS.md` still documents:

- `TreeNode` (renamed to `ActivityNode`) in design pattern #13
- `TimingFormat` (renamed) not mentioned
- `ColorWarning` (deleted) not mentioned
- `timingCache` (removed) not mentioned
- Zero mention of any split-brain work
- Timestamp: `2026-06-12` — stale by 5 days

### 2. CHANGELOG.md NOT updated

The `[Unreleased]` section is completely empty. Nine fixes shipped with zero changelog entries.

### 3. Zero tests added for ANY fix

No regression tests for:

- Color detector alignment
- ActivityNode rename
- Event constant usage
- TimingFormat rename
- GetWorkflowID return type change

### 4. TODO typo

`nom/inline_renderer.go:106`: `TODO(s split-brain M4)` — extra 's' in the tag.

### 5. `nom.ColorX` vars are accidentally mutable

`nom/symbols.go:35-50` uses `var` for color constants — any code can reassign them at runtime. This is a latent footgun that the split-brain audit identified but the fix didn't address.

---

## e) WHAT WE SHOULD IMPROVE

### Type Model Improvements

1. **`GetActivityCounts()` on subscriber** — eliminates the `activities` deep-copy cache entirely. The only production reader (`view.go:331`) needs counts, not the full map.

2. **`dependencyTree` → local variable** — it's a shared pointer (not a copy), so caching it in the struct provides zero correctness benefit. Fetch once at the top of the render function.

3. **Typed events** — replace string routing + 7 accessor interfaces with a sealed `Event` interface + type switch. Eliminates the silent-typo-drop failure mode and ~60 lines of accessor boilerplate.

4. **`nom.Theme` struct** — replace 6 mutable `var ColorX` globals with an immutable `Theme` struct. The `tui` module (which already imports `nom`) can consume it directly.

5. **Shared terminal detection** — adopt `charmbracelet/x/term` to replace both hand-rolled `isNoColor()`/`isCI()` (root) and `detectNoColor()` (nom). Same dep org as existing `charm.land/x/ansi`.

6. **`delimitedWriter` alignment** — add `WriteFooter` to the real `tableDataWriter` interface. Both concrete writers already implement it.

### Process Improvements

7. **Update AGENTS.md immediately** after architectural changes
8. **Update CHANGELOG.md** as part of every fix commit
9. **Add regression tests** for every behavioral change
10. **Run the brutal-self-review skill** before declaring done

---

## f) Top 25 Things to Get Done Next

Sorted by impact/effort ratio (highest first):

| #  | Task                                                                       | Impact | Effort   | Category       |
| -- | -------------------------------------------------------------------------- | ------ | -------- | -------------- |
| 1  | Fix remaining 11 bare event literals in `nom/subscriber_test.go`           | High   | 10min    | Half-assed fix |
| 2  | Fix TODO typo in `nom/inline_renderer.go:106`                              | Low    | 1min     | Typo           |
| 3  | Update `AGENTS.md` with all split-brain changes                            | High   | 20min    | Process        |
| 4  | Update `CHANGELOG.md` `[Unreleased]` with all fixes                        | High   | 15min    | Process        |
| 5  | Add `WriteFooter` to real `tableDataWriter` interface                      | Med    | 10min    | Half-assed fix |
| 6  | Add `GetActivityCounts()` to subscriber, delete `activities` field         | High   | 30min    | Type model     |
| 7  | Delete `dependencyTree` field, use local variable in render                | Med    | 20min    | Type model     |
| 8  | Add regression test for color detector agreement                           | High   | 20min    | Test gap       |
| 9  | Make `nom.ColorX` vars immutable (convert to struct or const-equivalent)   | Med    | 15min    | Type model     |
| 10 | Unify `"No activities to display"` to single source of truth               | Low    | 10min    | Half-assed fix |
| 11 | Add `detectNoColor` test in nom (currently zero test coverage)             | High   | 15min    | Test gap       |
| 12 | Add terminal check to nom `detectNoColor()`                                | High   | 10min    | Half-assed fix |
| 13 | Adopt `charmbracelet/x/term` for shared terminal detection                 | High   | 45min    | Library        |
| 14 | Create `nom.Theme` struct for color unification                            | Med    | 40min    | Type model     |
| 15 | Have `tui` consume `nom.Theme` instead of own `terminalColors`             | Med    | 30min    | Type model     |
| 16 | Add typed events (sealed interface + type switch)                          | High   | 90min    | Type model     |
| 17 | Add ActivityNode distinctness compile-time test                            | Low    | 5min     | Test gap       |
| 18 | Add cross-reference comment to `nom/symbols.go` for tui color mirroring    | Low    | 5min     | Process        |
| 19 | Consider `muesli/termenv` for color profile detection (TrueColor/256/ANSI) | Med    | Research | Library        |
| 20 | Plan M4: Rename `Render()` methods in next minor version                   | Med    | 60min    | Deferred       |
| 21 | Plan M5: Rename `ShapeBox` → `NodeShapeBox` in next minor                  | Med    | 60min    | Deferred       |
| 22 | Plan M6/M7: Introduce canonical `output.Direction` enum                    | Med    | 90min    | Deferred       |
| 23 | Plan M8: Align style struct field names across root/d2                     | Med    | 45min    | Deferred       |
| 24 | Update `SPLIT-BRAIN.html` report with resolved status                      | Low    | 15min    | Process        |
| 25 | Run `brutal-self-review` skill for final quality gate                      | High   | 30min    | Process        |

---

## g) Top #1 Question I Cannot Figure Out Myself

**Should we adopt `charmbracelet/x/term` (or `muesli/termenv`) to replace hand-rolled terminal/color detection, or keep the current hand-rolled approach?**

Arguments FOR adopting:

- Eliminates the split-brain entirely (one library call replaces two hand-rolled implementations)
- `charm.land/x/ansi` is already a dependency — same organization, zero new dep org
- Properly handles `FORCE_COLOR`, `NO_COLOR`, `CI`, TrueColor/256-color detection per spec
- Removes ~30 lines of env-var-checking boilerplate across two modules

Arguments AGAINST:

- Root module is zero-dep (only `x/term`, `go-branded-id`, `delimited`, `enum`, `testhelpers`). Adding `charmbracelet/x/term` to root would break the zero-lipgloss-deps guarantee for `go get github.com/larsartmann/go-output`
- nom already depends on lipgloss — so nom could adopt it, but root cannot
- The multi-module architecture was specifically designed to isolate terminal deps

**The architectural tension:** Root serves users who want zero terminal deps. nom serves users who want progress visualization. The "right" answer is probably: nom adopts `charmbracelet/x/term`, root keeps its hand-rolled version (for zero-dep guarantee), and we add a cross-module agreement test. But this means the split-brain is accepted as an architectural tradeoff, not truly eliminated.

**I need your decision on this tradeoff before proceeding.**
