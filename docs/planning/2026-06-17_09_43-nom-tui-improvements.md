# Comprehensive Plan: nom/ + tui/ Improvements

**Date:** 2026-06-17 09:43\
**Goal:** Make the generic concurrent-activity progress UI reliable and useful — NOM rendering inspiration only, no nix-specific features.\
**Pareto framing:**

- **1% → 51%:** Priority-ordered visibility (running/failed first).
- **4% → 64%:** Width truncation + logical line accounting (no screen corruption).
- **20% → 80%:** Timing median, redraw quality, terminal-state summary, cleanups.

---

## Dependency Graph

```d2
direction: down

Ranking: 1.1 Define status ranking
SortKey: 1.2 Add SortKey type
Collect: 1.3 Build candidate list with SortKey
Context: 1.4 Preserve tree context for top-N
RenderUpdate: 1.5 Update Render/VisibleNodes
UnitTests: 1.6 Priority-order unit tests
GoldenUpdate: 1.7 Update golden tests
ExampleUpdate: 1.8 Update nom_progress example

WidthHelper: 2.1 Terminal width helper
AvailWidth: 2.2 Compute available width
Truncate: 2.3 Truncate long lines
LineCount: 2.4 Logical/physical line count
WidthTests: 2.5 Width/truncation tests
TUICheck: 2.6 Verify TUI uses same logic

MedianHelper: 3.1 Median helper
ReplaceMedian: 3.2 Replace GetAverage with GetMedian
UsageMedian: 3.3 Update subscriber usage

RefreshAPI: 4.1 Refresh signal API
Loop: 4.2 Select-based Start loop
RedrawTests: 4.3 Redraw timing tests

StatusInline: 5.1 Final status line in InlineRenderer
StatusTUI: 5.2 Final status line in ProgressModel
StatusTests: 5.3 Update status tests/goldens

Cleanup1: 6.1 Remove unused errTestFailure
Cleanup2: 6.2 Fix WriteString warning
Verify1: 6.3 go test ./... nom
Verify2: 6.4 go test ./... tui
Verify3: 6.5 golangci-lint nom + tui

Ranking -> SortKey
SortKey -> Collect
Collect -> Context
Context -> RenderUpdate
RenderUpdate -> UnitTests
RenderUpdate -> GoldenUpdate
ExampleUpdate -> RenderUpdate

WidthHelper -> AvailWidth
AvailWidth -> Truncate
Truncate -> LineCount
LineCount -> WidthTests
Truncate -> TUICheck

MedianHelper -> ReplaceMedian
ReplaceMedian -> UsageMedian

RefreshAPI -> Loop
Loop -> RedrawTests

StatusInline -> StatusTests
StatusTUI -> StatusTests
```

---

## Task Table

| Order | ID  | Task                                                           | Area       | Impact | Effort | Customer Value                       | Est. min |
| ----: | --- | -------------------------------------------------------------- | ---------- | ------ | ------ | ------------------------------------ | -------- |
|     1 | 1.1 | Define activity-interest ranking type and status→rank mapping  | Visibility | High   | Low    | Core value: see hot activities first | 10       |
|     2 | 1.2 | Add `SortKey` struct + `Less` method (status, elapsed, depth)  | Visibility | High   | Low    | Enables deterministic ordering       | 10       |
|     3 | 1.3 | Build candidate visible-node list scored by `SortKey`          | Visibility | High   | Med    | Hot nodes compete for screen space   | 12       |
|     4 | 1.4 | Preserve parent/child context when selecting top-N nodes       | Visibility | High   | Med    | Tree prefixes still render correctly | 12       |
|     5 | 1.5 | Update `Render` and `VisibleNodes` to use priority ordering    | Visibility | High   | Low    | Wires ordering into output           | 10       |
|     6 | 1.6 | Add unit tests for priority ordering                           | Visibility | High   | Med    | Prevents regression                  | 12       |
|     7 | 1.7 | Update golden tests/expected outputs                           | Visibility | High   | Low    | Keeps CI green                       | 10       |
|     8 | 1.8 | Update `examples/nom_progress` to demo concurrent ordering     | Visibility | Low    | Low    | Shows the feature                    | 10       |
|     9 | 2.1 | Add terminal-width helper (`GetTerminalWidth`)                 | Width      | High   | Low    | Prerequisite for truncation          | 10       |
|    10 | 2.2 | Compute available render width (account for prefix + timing)   | Width      | High   | Low    | Knows how much room the text has     | 10       |
|    11 | 2.3 | Truncate activity names/annotations with `…`                   | Width      | High   | Med    | No overflow/corruption               | 12       |
|    12 | 2.4 | Track logical vs physical line count in `InlineRenderer`       | Width      | High   | Med    | Correct cursor-up redraw             | 12       |
|    13 | 2.5 | Add tests for truncation and narrow-terminal behavior          | Width      | High   | Med    | Verifies correctness                 | 12       |
|    14 | 2.6 | Verify TUI `renderNOMStyle` uses same truncation logic         | Width      | Low    | Low    | Consistency across surfaces          | 8        |
|    15 | 3.1 | Add median helper for `[]time.Duration`                        | Timing     | Med    | Low    | Robust statistics                    | 8        |
|    16 | 3.2 | Replace `GetAverage` with `GetMedian`; update tests            | Timing     | Med    | Low    | Better estimates                     | 10       |
|    17 | 3.3 | Update subscriber usage from average to median                 | Timing     | Low    | Low    | Wiring                               | 5        |
|    18 | 4.1 | Add refresh-signal API to `InlineRenderer`                     | Redraw     | Med    | Low    | On-demand updates                    | 10       |
|    19 | 4.2 | Rewrite `Start` loop: ctx + refresh + ticker + 1s max-frame    | Redraw     | Med    | Med    | Smooth, efficient timers             | 12       |
|    20 | 4.3 | Add tests for refresh timing behavior                          | Redraw     | Med    | Med    | Confidence                           | 12       |
|    21 | 5.1 | Add terminal-state line to `InlineRenderer.Finish`             | Summary    | Med    | Low    | Clear completion state               | 10       |
|    22 | 5.2 | Add terminal-state line to `ProgressModel` summary             | Summary    | Med    | Low    | Consistent TUI output                | 10       |
|    23 | 5.3 | Update tests/goldens for final status line                     | Summary    | Low    | Med    | CI green                             | 10       |
|    24 | 6.1 | Remove unused `errTestFailure` in `tui/test_sentinels_test.go` | Cleanup    | Low    | Low    | Lint cleanliness                     | 5        |
|    25 | 6.2 | Fix inefficient `WriteString` warning in `inline_renderer.go`  | Cleanup    | Low    | Low    | Lint cleanliness                     | 5        |
|    26 | 6.3 | Run `go test ./...` in `nom` module                            | Verify     | High   | Low    | Must pass                            | 10       |
|    27 | 6.4 | Run `go test ./...` in `tui` module                            | Verify     | High   | Low    | Must pass                            | 10       |
|    28 | 6.5 | Run `golangci-lint` on `nom` and `tui`                         | Verify     | High   | Low    | Must pass                            | 10       |

**Total tasks:** 28\
**Total estimated time:** ~270 minutes (4.5 h)\
**Max task size:** 12 minutes

---

## Risk Notes

- **Priority ordering must not orphan children.** If a child is high-interest but its parent low-interest, include the parent as a collapsed/hidden context line or ensure prefixes are still correct.
- **Width detection** should work for non-terminal writers (e.g. tests use `bytes.Buffer`). Fall back gracefully to a default width.
- **Line-count accounting** interacts with priority ordering (fewer visible lines → fewer wrap issues), but both must land together for correctness.
- **API surface is frozen** per ADR 006; only add new exported symbols, don't change existing signatures.
