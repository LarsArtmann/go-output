# Status Report: go-output — Zero-Lint Stable State

**Date:** 2026-05-25 21:06
**Scope:** Full codebase audit after zero-lint achievement, test helper dedup, and BuildFlow alignment
**State:** All 12 modules BUILD ✅ TEST ✅ RACE ✅ LINT ✅ — zero issues, working tree clean

---

## A) FULLY DONE ✅

### Multi-Module Architecture — 12 Independent Modules

| Module           | Coverage | Lint | Files            | Role                                            |
| ---------------- | -------- | ---- | ---------------- | ----------------------------------------------- |
| root (`output`)  | 88.8%    | 0    | 26 src + 18 test | Core types, Markdown, Tree, Streaming, Registry |
| `delimited/`     | 84.8%    | 0    | 3 src + 4 test   | CSV + TSV writers and formatters                |
| `serialization/` | 84.4%    | 0    | 4 src + 5 test   | JSON + YAML marshaling and renderers            |
| `markup/`        | 87.9%    | 0    | 4 src + 5 test   | XML + HTML + Streaming HTML renderers           |
| `d2/`            | 100%     | 0    | 5 src + 6 test   | D2 diagram renderer (rich domain model)         |
| `graph/`         | 96.0%    | 0    | 4 src + 5 test   | DOT + Mermaid renderers                         |
| `enum/`          | 100%     | 0    | 1 src + 2 test   | Generic enum utilities                          |
| `escape/`        | 100%     | 0    | 1 src + 1 test   | Format-specific escaping                        |
| `testhelpers/`   | 81.1%    | 0    | 2 src + 1 test   | Shared test assertions + writer types           |
| `table/`         | 100%     | 0    | 1 src + 1 test   | Lipgloss terminal tables                        |
| `integration/`   | 82.8%    | 0    | 1 src + 5 test   | Cross-module integration tests                  |
| `examples/`      | n/a      | 0    | 3 src            | Usage examples (no tests)                       |

**Totals:** 105 Go files, 14,781 LOC, 45 source files, 60 test files.

### Lint Compliance — 80+ Linters, Zero Issues

Every module passes `golangci-lint run` cleanly with 80+ enabled linters including:
gochecknoglobals, gochecknoinits, wsl_v5, goimports, gci, depguard, wrapcheck, dupl,
staticcheck, govet, errcheck, cyclop, gocritic, revive, gosec, and 65+ more.

### Root Module Invariant

**Root has ZERO sub-module imports in production code.** Verified via `go mod graph`.
Users who `go get github.com/larsartmann/go-output` get zero transitive deps from
sub-modules they don't import (no lipgloss, no yaml, no d2/graph deps).

### Dependency Isolation

| Dependency                             | Isolated To           | In Root? |
| -------------------------------------- | --------------------- | -------- |
| `charm.land/lipgloss/v2`               | `table/` only         | ❌       |
| `github.com/go-faster/yaml`            | `serialization/` only | ❌       |
| `github.com/larsartmann/go-branded-id` | root only             | ✅       |
| `golang.org/x/term`                    | root only             | ✅       |

### CI Pipeline

4-step CI (build, test, tidy, vulncheck) covers all 12 modules via loop in `.github/workflows/ci.yml`.

### BuildFlow Pre-Commit

- BuildFlow pre-commit passes **without `--no-verify`**
- Auto-formatters (goimports, gofumpt, golines) aligned with golangci-lint config
- `gci` sections configured with `Prefix(github.com/larsartmann/go-output)` to match goimports `local-prefixes`
- Auto-deduplicated test writer types into `testhelpers/writers.go`

### Race Detector

All 12 modules pass with `-race` flag — zero data races detected.

### Session Commits (this conversation)

| Commit    | Description                                                                         |
| --------- | ----------------------------------------------------------------------------------- |
| `e502c79` | `chore: lint compliance, test helpers refactor, and documentation formatting`       |
| `7d94d14` | `chore(deps): update flake.lock with latest nixpkgs`                                |
| `5f9a849` | `docs(status): add zero-lint achievement comprehensive status report`               |
| `ae06615` | `refactor(testhelpers): extract shared test writer types to testhelpers/writers.go` |

### What Was Fixed This Session

1. **42 → 0 lint issues** across all 12 modules
2. **goimports/gci import grouping** — proper stdlib → third-party → local sections (6 files)
3. **wsl_v5 whitespace** — 19 blank-line fixes across 9 files
4. **gochecknoinits nolint** — 5 registry init() functions with justification comments
5. **gochecknoglobals nolint** — emptyYAML constant + test helper re-exports
6. **wrapcheck config** — ignore intra-project packages in `.golangci.yml`
7. **dupl dedup** — extracted `testRenderTableData` shared helper in delimited/
8. **SA1019 nolint** — deprecated registry API test in integration/
9. **gocritic config** — removed 3 already-disabled checks
10. **gci/goimports conflict** — configured gci sections to match goimports local-prefixes
11. **BuildFlow alignment** — resolved pre-commit hook formatting conflicts
12. **Test helper extraction** — `errorWriter` + `writeNThenFailWriter` → `testhelpers/writers.go`

---

## B) PARTIALLY DONE 🟡

### Coverage Gaps (modules below 90% target)

| Module               | Coverage | Gap  | What's Missing                                                                                |
| -------------------- | -------- | ---- | --------------------------------------------------------------------------------------------- |
| `integration/`       | 82.8%    | 7.2% | Error paths in cross-module tests, edge cases                                                 |
| `serialization/`     | 84.4%    | 5.6% | JSON/YAML graph renderer error paths                                                          |
| `delimited/`         | 84.8%    | 5.2% | DelimitedWriter error paths, TSV type-switch                                                  |
| `root (output)`      | 88.8%    | 1.2% | Minor error paths in render_tabledata, streaming                                              |
| `gentest (internal)` | 87.5%    | 2.5% | One untested assertion helper                                                                 |
| `testhelpers/`       | 81.1%    | 8.9% | New `writers.go` has uncovered paths (ErrorWriter.Write, WriteNThenFailWriter.Write branches) |

### Documentation

- ✅ README.md, CHANGELOG.md, AGENTS.md up to date
- ✅ ADR-001, ADR-002, ADR-003 written
- 🟡 ADR-004 (delimited/serialization/markup extraction rationale) not yet written

---

## C) NOT STARTED ⬜

### From TODO_LIST.md (5 open items)

| #  | Item                                                   | Priority | Status                 |
| -- | ------------------------------------------------------ | -------- | ---------------------- |
| 20 | Move `internal/gentest` to `testhelpers/gentest`?      | P3       | Needs Decision         |
| 21 | Duplicated test helpers in graph/ (depends on #20)     | P3       | Blocked                |
| 24 | Pre-commit: `go-structure-linter` false positives      | P4       | External tool issue    |
| 26 | flake.nix Go build/test/lint                           | P4       | Blocked by Nix sandbox |
| 34 | Shape-specific renderer constructors (ADR-002 Phase 2) | P6       | Future                 |

### Future Formats (TODO P6)

| #  | Format                         | Estimated Effort |
| -- | ------------------------------ | ---------------- |
| 35 | TOML format (new module)       | 2 hr             |
| 36 | JSONL format (new renderer)    | 1 hr             |
| 37 | PlantUML format (new module)   | 4 hr             |
| 38 | AsciiDoc format (new renderer) | 1 hr             |

### Other Not Started

- Fuzz tests for delimited/, serialization/, markup/
- JSON/YAML streaming renderers
- Pre-v1 API stability audit (#39)
- Community posting to r/golang, Awesome Go (#40)
- Version tagging (v0.5.0 decision)
- ADR-004 for module extraction rationale

---

## D) TOTALLY FUCKED UP 💥

**Nothing is totally fucked up.** The codebase is in its best-ever state:

- Zero build failures across all 12 modules
- Zero test failures (with race detector)
- Zero lint issues (80+ linters)
- Zero circular dependencies
- Zero leaked dependencies across module boundaries
- BuildFlow pre-commit passes cleanly
- Working tree is clean, all changes pushed

**Minor annoyances (not bugs):**

- `go.work` is gitignored — local dev requires manual workspace setup (documented in AGENTS.md)
- `go mod tidy` produces expected remote resolution errors for local replace directives
- `testhelpers/` coverage dropped from 93.8% → 81.1% because BuildFlow extracted `writers.go` without adding tests for the new exported types
- `serialization/testhelpers_test.go` still has its own local `errorWriter` (BuildFlow didn't dedup it since serialization is a separate module)

---

## E) WHAT WE SHOULD IMPROVE

### Architecture & Type Model

1. **`TableData` is `[][]string`-only** — no type-safe column access. Consider `TableData[T]` or a `Row` type with named fields.
2. **Error types are inconsistent** — `InvalidFormatError`, `UnsupportedFormatError`, `InvalidColorModeError` are separate. Consider unified `output.Error` with kind enum.
3. **Registry API is deprecated but still exists** — `Register()`/`Create()`/`Unregister()`/`RegisteredFormats()`/`IsRegistered()` are deprecated but not removed. Dead code or keep for v1?
4. **No `io.WriterTo` integration** — Renderers return `(string, error)`. Adding `WriteTo(io.Writer)` would be more idiomatic.
5. **Streaming API inconsistency** — CSV/TSV use `Writer`, HTML uses `StreamingHTMLRenderer`, JSON uses `JSONWriter.Encode()`. Should share a common interface.

### Test Infrastructure

6. **testhelpers coverage dropped** — 93.8% → 81.1% after `writers.go` extraction. Need tests for `ErrorWriter.Write`, `WriteNThenFailWriter.Write` success/failure branches.
7. **Test helper duplication remains in serialization/** — `serialization/testhelpers_test.go` still has local `errorWriter` (separate module, not deduped).
8. **No coverage enforcement in CI** — Could drop silently. Add minimum threshold.
9. **No snapshot testing** — `assertContains` is fragile. Snapshot tests would catch unintended output changes.

### Developer Experience

10. **No companion CLI** — A `go-output` CLI tool (like `jq`) would dramatically increase adoption.
11. **No benchmarks in CI** — Performance could regress silently.

### Build & Release

12. **No version tag** — 5+ BREAKING changes bundled in `[Unreleased]`. Need v0.5.0 decision.
13. **No goreleaser** — Release process is manual.

---

## F) TOP 25 THINGS TO DO NEXT

### HIGH IMPACT, LOW EFFORT (do first)

| #  | Task                                                                                                                    | Impact | Effort |
| -- | ----------------------------------------------------------------------------------------------------------------------- | ------ | ------ |
| 1  | **Fix testhelpers/ coverage** — Add tests for `ErrorWriter`, `WriteNThenFailWriter` in writers.go (dropped 93.8%→81.1%) | High   | 15 min |
| 2  | **Tag v0.5.0** — Bundle BREAKING changes, update CHANGELOG date                                                         | High   | 5 min  |
| 3  | **Add fuzz tests to delimited/, serialization/, markup/** — Copy pattern from root/d2/graph                             | High   | 1 hr   |
| 4  | **Write ADR-004** — Document delimited/serialization/markup extraction rationale                                        | Med    | 20 min |
| 5  | **Add JSONL format** — Trivial extension of JSON writer, high demand                                                    | High   | 1 hr   |
| 6  | **Fix integration/ coverage to 90%+** — Add error path tests                                                            | Med    | 30 min |
| 7  | **Fix serialization/ coverage to 90%+** — Add graph renderer error tests                                                | Med    | 20 min |
| 8  | **Fix delimited/ coverage to 90%+** — Add DelimitedWriter error tests                                                   | Med    | 20 min |
| 9  | **Add coverage threshold to CI** — Fail if any module drops below 80%                                                   | Med    | 15 min |
| 10 | **Deduplicate serialization/ test helpers** — Use testhelpers.ErrorWriter via type alias                                | Low    | 10 min |

### HIGH IMPACT, MEDIUM EFFORT

| #  | Task                                                                                            | Impact | Effort |
| -- | ----------------------------------------------------------------------------------------------- | ------ | ------ |
| 11 | **Move internal/gentest to testhelpers/gentest** — Eliminate test helper duplication (TODO #20) | Med    | 30 min |
| 12 | **Deduplicate graph/ test helpers** — Depends on #11 (TODO #21)                                 | Med    | 15 min |
| 13 | **Unified error types** — `output.Error` with `Kind` enum                                       | High   | 2 hr   |
| 14 | **Consistent streaming API** — Common `StreamingWriter` interface                               | High   | 3 hr   |
| 15 | **Add TOML format** — New module, popular in Go ecosystem                                       | Med    | 2 hr   |
| 16 | **Pre-v1 API audit** — Review all public APIs for stability                                     | High   | 3 hr   |
| 17 | **goreleaser config** — Automate releases                                                       | Med    | 1 hr   |
| 18 | **Move userjourney_test.go to integration/** — Cleaner separation                               | Low    | 15 min |

### MEDIUM IMPACT, HIGHER EFFORT

| #  | Task                                                           | Impact    | Effort |
| -- | -------------------------------------------------------------- | --------- | ------ |
| 19 | **Add snapshot testing** — Catch unintended output changes     | Med       | 2 hr   |
| 20 | **Generic TableData[T]** — Type-safe column access             | Very High | 1 day  |
| 21 | **go-output CLI tool** — Companion CLI for format conversion   | Very High | 2 days |
| 22 | **Add PlantUML format** — New module with UML diagram types    | Med       | 4 hr   |
| 23 | **Benchmark regression CI** — Track performance across commits | Med       | 2 hr   |
| 24 | **Shape-specific renderer constructors** — ADR-002 Phase 2     | Med       | 3 hr   |
| 25 | **Community posting** — r/golang, Awesome Go                   | High      | 1 hr   |

---

## G) TOP QUESTION FOR LARS

**Should we remove the deprecated Registry API before v0.5.0?**

Current state:

- 5 deprecated functions: `Register`, `Create`, `Unregister`, `RegisteredFormats`, `IsRegistered`
- Only tested in `integration/format_test.go` with `//nolint:staticcheck`
- Replaced by `RegisterTableDataMarshaler` for TableData dispatch
- But nothing replaces the generic `Renderer` factory use case (runtime dispatch for custom formats)

Options:

1. **Remove now** — Clean break, bundle with existing BREAKING changes in v0.5.0
2. **Keep for v0.x** — Remove in v1.0 after providing alternative
3. **Replace with better API** — `RegisterRenderer(format, factory)` that doesn't overlap with TableDataMarshaler
