# go-output Comprehensive Completion Plan

**Generated:** 2026-03-28  
**Date:** Saturday, March 28, 2026  
**Branch:** master  
**Author:** Crush AI  

---

## Executive Summary

This document contains **ALL TODOs** extracted from the March 27 comprehensive status review, broken down into actionable tasks of **max 12 minutes each**, sorted by **priority score** (Importance × Impact × Customer Value ÷ Effort).

| Category | Count | High Priority |
|----------|-------|---------------|
| Critical Bugs | 1 | 1 |
| Technical Debt | 3 | 2 |
| Code Quality | 5 | 2 |
| Testing Improvements | 4 | 2 |
| Documentation | 3 | 1 |
| New Features | 8 | 2 |
| Future Enhancements | 12 | 0 |
| **Total** | **36** | **10** |

---

## Priority Calculation

**Priority Score = (Importance × Impact × Customer Value) ÷ Effort**

| Score Range | Priority Level | Action |
|-------------|----------------|--------|
| 15-25 | 🔴 CRITICAL | Fix immediately |
| 10-14 | 🟠 HIGH | Do next sprint |
| 5-9 | 🟡 MEDIUM | Backlog |
| 1-4 | 🟢 LOW | Future |

---

## 🔴 CRITICAL (Score: 15-25)

### C1: Fix GraphNode.GetStyle() Bug
**File:** `graph.go:38-43`  
**Score:** 20 (Importance: 5, Impact: 4, Value: 5, Effort: 2)  

| # | Task | Time | Status |
|---|------|------|--------|
| C1.1 | Read graph.go and understand GetStyle() context | 3 min | ⬜ |
| C1.2 | Fix GetStyle() to return actual style when set | 4 min | ⬜ |
| C1.3 | Add test case for GetStyle() with non-zero style | 3 min | ⬜ |
| C1.4 | Run tests to verify fix | 2 min | ⬜ |

---

## 🟠 HIGH PRIORITY (Score: 10-14)

### H1: Add cmdguard to go.mod
**Issue:** CLI integration code exists but dependency not declared  
**Score:** 14 (Importance: 4, Impact: 4, Value: 4, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| H1.1 | Check cmdguard v2 availability on GitHub | 2 min | ⬜ |
| H1.2 | Add cmdguard as test dependency to go.mod | 3 min | ⬜ |
| H1.3 | Verify cmdguard tests still pass | 3 min | ⬜ |
| H1.4 | Update documentation if needed | 4 min | ⬜ |

### H2: Standardize Render() Error Returns
**Issue:** MarkdownTable returns error, HTMLRenderer doesn't  
**Score:** 12 (Importance: 3, Impact: 4, Value: 4, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| H2.1 | Audit all Render() methods for error returns | 4 min | ⬜ |
| H2.2 | Add error returns to methods missing them | 5 min | ⬜ |
| H2.3 | Update tests for error handling | 3 min | ⬜ |

### H3: Split format_test.go (465 lines → target 350)
**Issue:** 32.9% over file size limit  
**Score:** 12 (Importance: 3, Impact: 3, Value: 3, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| H3.1 | Analyze format_test.go structure | 3 min | ⬜ |
| H3.2 | Extract Format enum tests to format_enum_test.go | 4 min | ⬜ |
| H3.3 | Extract TableData tests to table_data_test.go | 4 min | ⬜ |
| H3.4 | Run tests to verify split | 1 min | ⬜ |

### H4: Split format.go (365 lines → target 350)
**Issue:** 4.3% over file size limit  
**Score:** 12 (Importance: 3, Impact: 2, Value: 3, Effort: 2)  

| # | Task | Time | Status |
|---|------|------|--------|
| H4.1 | Analyze format.go structure | 3 min | ⬜ |
| H4.2 | Move GraphNode/GetStyle to graph.go | 4 min | ⬜ |
| H4.3 | Move TreeNode to tree.go if appropriate | 4 min | ⬜ |
| H4.4 | Verify no circular dependencies | 1 min | ⬜ |

### H5: Add User Journey Tests
**Issue:** Tests are implementation-focused, not user-focused  
**Score:** 12 (Importance: 4, Impact: 4, Value: 5, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| H5.1 | Create userjourney_test.go with CLI dev stories | 5 min | ⬜ |
| H5.2 | Add "As a CLI developer" test naming | 4 min | ⬜ |
| H5.3 | Add Given/When/Then comments to existing tests | 3 min | ⬜ |

### H6: Refactor custom containsString in Tests
**File:** `graph_test.go:212-223`  
**Score:** 10 (Importance: 2, Impact: 2, Value: 2, Effort: 2)  

| # | Task | Time | Status |
|---|------|------|--------|
| H6.1 | Find all containsString custom functions | 2 min | ⬜ |
| H6.2 | Replace with strings.Contains | 3 min | ⬜ |
| H6.3 | Run tests to verify | 2 min | ⬜ |

### H7: Document Streaming Semantics
**Issue:** Streaming adapter collects all before writing  
**Score:** 10 (Importance: 3, Impact: 3, Value: 3, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| H7.1 | Read streaming.go and understand adapter pattern | 4 min | ⬜ |
| H7.2 | Add code comments explaining semantic | 4 min | ⬜ |
| H7.3 | Update README if streaming section exists | 4 min | ⬜ |

### H8: Fix Cyclomatic Complexity - TestTableData
**File:** `format_test.go:251` (complexity 12, max 10)  
**Score:** 10 (Importance: 3, Impact: 3, Value: 3, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| H8.1 | Analyze TestTableData function | 3 min | ⬜ |
| H8.2 | Extract sub-cases to subtests | 5 min | ⬜ |
| H8.3 | Run golangci-lint to verify | 4 min | ⬜ |

### H9: Fix Cyclomatic Complexity - Stream
**File:** `streaming.go:56` (complexity 11, max 10)  
**Score:** 10 (Importance: 3, Impact: 3, Value: 3, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| H9.1 | Analyze Stream function | 3 min | ⬜ |
| H9.2 | Extract error handling logic | 5 min | ⬜ |
| H9.3 | Run golangci-lint to verify | 4 min | ⬜ |

### H10: Add Integration Tests for Workflows
**Score:** 11 (Importance: 4, Impact: 4, Value: 4, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| H10.1 | Create workflow_test.go | 4 min | ⬜ |
| H10.2 | Add format conversion workflow tests | 4 min | ⬜ |
| H10.3 | Add CLI flag integration tests | 4 min | ⬜ |

---

## 🟡 MEDIUM PRIORITY (Score: 5-9)

### M1: Add Property-Based Tests (go-fuzz)
**Score:** 8 (Importance: 3, Impact: 3, Value: 4, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| M1.1 | Review existing fuzz tests | 3 min | ⬜ |
| M1.2 | Add fuzz test for CSV formatter | 4 min | ⬜ |
| M1.3 | Add fuzz test for Markdown renderer | 4 min | ⬜ |
| M1.4 | Run fuzz tests for 1 minute | 1 min | ⬜ |

### M2: Add Configuration Options per Format
**Score:** 8 (Importance: 3, Impact: 3, Value: 4, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| M2.1 | Design configuration interface | 5 min | ⬜ |
| M2.2 | Add options to JSON renderer | 3 min | ⬜ |
| M2.3 | Add options to CSV renderer | 3 min | ⬜ |
| M2.4 | Add tests for options | 1 min | ⬜ |

### M3: Extract Hardcoded CSS for HTML
**Score:** 6 (Importance: 2, Impact: 2, Value: 3, Effort: 2)  

| # | Task | Time | Status |
|---|------|------|--------|
| M3.1 | Find hardcoded CSS in html.go | 3 min | ⬜ |
| M3.2 | Extract to constants or template file | 5 min | ⬜ |
| M3.3 | Add tests for CSS customization | 4 min | ⬜ |

### M4: Add TSV Format
**Score:** 6 (Importance: 2, Impact: 3, Value: 3, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| M4.1 | Create tsv.go based on csv.go | 6 min | ⬜ |
| M4.2 | Add tests with >90% coverage | 4 min | ⬜ |
| M4.3 | Add to format registry | 2 min | ⬜ |

### M5: Improve D2FromTableData for Multiple Tables
**Score:** 6 (Importance: 2, Impact: 3, Value: 3, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| M5.1 | Review current D2FromTableData implementation | 4 min | ⬜ |
| M5.2 | Add multi-table support | 5 min | ⬜ |
| M5.3 | Add tests | 3 min | ⬜ |

### M6: Add Release Workflow
**Score:** 6 (Importance: 3, Impact: 3, Value: 3, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| M6.1 | Design release workflow | 4 min | ⬜ |
| M6.2 | Create .github/workflows/release.yml | 5 min | ⬜ |
| M6.3 | Test workflow in PR | 3 min | ⬜ |

### M7: Add XML Format
**Score:** 6 (Importance: 2, Impact: 3, Value: 3, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| M7.1 | Create xml.go | 6 min | ⬜ |
| M7.2 | Add tests with >90% coverage | 4 min | ⬜ |
| M7.3 | Add to format registry | 2 min | ⬜ |

### M8: Add Performance Benchmarks for All Formats
**Score:** 8 (Importance: 3, Impact: 4, Value: 4, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| M8.1 | Review existing benchmarks_test.go | 3 min | ⬜ |
| M8.2 | Add CSV benchmark | 2 min | ⬜ |
| M8.3 | Add Markdown benchmark | 2 min | ⬜ |
| M8.4 | Add YAML benchmark | 2 min | ⬜ |
| M8.5 | Add HTML benchmark | 2 min | ⬜ |
| M8.6 | Run benchmarks and save baseline | 1 min | ⬜ |

### M9: Create API Documentation (Go.dev)
**Score:** 6 (Importance: 3, Impact: 3, Value: 3, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| M9.1 | Review godoc requirements | 3 min | ⬜ |
| M9.2 | Ensure all exported funcs have doc comments | 6 min | ⬜ |
| M9.3 | Publish to pkg.go.dev | 3 min | ⬜ |

---

## 🟢 LOW PRIORITY (Score: 1-4)

### L1: Add CSV Dialect Support (Excel-compatible)
**Score:** 4 (Importance: 1, Impact: 2, Value: 3, Effort: 2)  

| # | Task | Time | Status |
|---|------|------|--------|
| L1.1 | Research CSV dialect options | 4 min | ⬜ |
| L1.2 | Add configuration to CSVWriter | 4 min | ⬜ |
| L1.3 | Add tests | 4 min | ⬜ |

### L2: Add Markdown Extensions (GFM tables)
**Score:** 4 (Importance: 1, Impact: 2, Value: 3, Effort: 2)  

| # | Task | Time | Status |
|---|------|------|--------|
| L2.1 | Review GFM table spec | 4 min | ⬜ |
| L2.2 | Add GFM options to MarkdownTable | 5 min | ⬜ |
| L2.3 | Add tests | 3 min | ⬜ |

### L3: Add ANSI Color Support
**Score:** 4 (Importance: 1, Impact: 2, Value: 3, Effort: 2)  

| # | Task | Time | Status |
|---|------|------|--------|
| L3.1 | Review lipgloss color usage | 3 min | ⬜ |
| L3.2 | Add color presets to color.go | 5 min | ⬜ |
| L3.3 | Add tests | 4 min | ⬜ |

### L4: Add JSON Schema Generation
**Score:** 4 (Importance: 1, Impact: 2, Value: 3, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| L4.1 | Research JSON Schema generation | 4 min | ⬜ |
| L4.2 | Create schema.go | 6 min | ⬜ |
| L4.3 | Add tests | 2 min | ⬜ |

### L5: Add Plugin System for Dynamic Formats
**Score:** 3 (Importance: 1, Impact: 2, Value: 2, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| L5.1 | Design plugin interface | 6 min | ⬜ |
| L5.2 | Implement plugin registry | 5 min | ⬜ |
| L5.3 | Add example plugin | 3 min | ⬜ |
| L5.4 | Add tests | 2 min | ⬜ |

### L6: WebAssembly Support
**Score:** 2 (Importance: 1, Impact: 1, Value: 2, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| L6.1 | Research WASM compatibility issues | 4 min | ⬜ |
| L6.2 | Add WASM build tags if needed | 6 min | ⬜ |
| L6.3 | Create example WASM page | 4 min | ⬜ |

### L7: Internationalization Support
**Score:** 2 (Importance: 1, Impact: 1, Value: 2, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| L7.1 | Design i18n interface | 6 min | ⬜ |
| L7.2 | Add locale support to formatters | 5 min | ⬜ |
| L7.3 | Add tests | 3 min | ⬜ |

### L8: Streaming Benchmarks (Memory Usage)
**Score:** 3 (Importance: 1, Impact: 2, Value: 2, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| L8.1 | Design memory benchmark | 4 min | ⬜ |
| L8.2 | Implement memory profiler | 5 min | ⬜ |
| L8.3 | Document findings | 3 min | ⬜ |

### L9: Fuzzing Corpus for Dedicated Fuzzing
**Score:** 3 (Importance: 1, Impact: 2, Value: 2, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| L9.1 | Create testdata/fuzz directory | 2 min | ⬜ |
| L9.2 | Add corpus files for each formatter | 6 min | ⬜ |
| L9.3 | Document corpus usage | 4 min | ⬜ |

### L10: Performance Regression CI
**Score:** 4 (Importance: 2, Impact: 3, Value: 2, Effort: 4)  

| # | Task | Time | Status |
|---|------|------|--------|
| L10.1 | Design benchmark comparison approach | 4 min | ⬜ |
| L10.2 | Create benchmark comparison workflow | 6 min | ⬜ |
| L10.3 | Test in PR | 4 min | ⬜ |

### L11: Add Phantom Type Violations (8 violations)
**Score:** 3 (Importance: 1, Impact: 1, Value: 2, Effort: 3)  

| # | Task | Time | Status |
|---|------|------|--------|
| L11.1 | List all 8 phantom type violations | 3 min | ⬜ |
| L11.2 | Fix violations in graph.go | 4 min | ⬜ |
| L11.3 | Fix violations in other files | 4 min | ⬜ |
| L11.4 | Run strong-id linter to verify | 1 min | ⬜ |

### L12: DOT/Mermaid Mixin Opportunity
**Score:** 2 (Importance: 1, Impact: 1, Value: 1, Effort: 2)  

| # | Task | Time | Status |
|---|------|------|--------|
| L12.1 | Analyze shared fields between DOT and Mermaid | 4 min | ⬜ |
| L12.2 | Extract shared interface if beneficial | 6 min | ⬜ |
| L12.3 | Update tests | 2 min | ⬜ |

---

## Summary Table (All Tasks Sorted by Priority)

|| # | Task | Score | Time | Priority |
||---|-----|-------|------|----------|
|| 1 | **C1: Fix GraphNode.GetStyle() Bug** | 20 | 12 min | 🔴 CRITICAL |
|| 2 | **H1: Add cmdguard to go.mod** | 14 | 12 min | 🟠 HIGH |
|| 3 | **H2: Standardize Render() Error Returns** | 12 | 12 min | 🟠 HIGH |
|| 4 | **H3: Split format_test.go** | 12 | 12 min | 🟠 HIGH |
|| 5 | **H4: Split format.go** | 12 | 12 min | 🟠 HIGH |
|| 6 | **H5: Add User Journey Tests** | 12 | 12 min | 🟠 HIGH |
|| 7 | **H8: Fix TestTableData Complexity** | 10 | 12 min | 🟠 HIGH |
|| 8 | **H9: Fix Stream Complexity** | 10 | 12 min | 🟠 HIGH |
|| 9 | **H6: Refactor containsString** | 10 | 7 min | 🟠 HIGH |
|| 10 | **H7: Document Streaming Semantics** | 10 | 12 min | 🟠 HIGH |
|| 11 | **H10: Add Workflow Integration Tests** | 11 | 12 min | 🟠 HIGH |
|| 12 | **M1: Add Property-Based Tests** | 8 | 12 min | 🟡 MEDIUM |
|| 13 | **M2: Add Configuration Options** | 8 | 12 min | 🟡 MEDIUM |
|| 14 | **M8: Add Benchmarks for All Formats** | 8 | 12 min | 🟡 MEDIUM |
|| 15 | **M3: Extract Hardcoded CSS** | 6 | 12 min | 🟡 MEDIUM |
|| 16 | **M4: Add TSV Format** | 6 | 12 min | 🟡 MEDIUM |
|| 17 | **M5: Improve D2 Multi-Table** | 6 | 12 min | 🟡 MEDIUM |
|| 18 | **M6: Add Release Workflow** | 6 | 12 min | 🟡 MEDIUM |
|| 19 | **M7: Add XML Format** | 6 | 12 min | 🟡 MEDIUM |
|| 20 | **M9: Create API Documentation** | 6 | 12 min | 🟡 MEDIUM |
|| 21 | **L1: Add CSV Dialect Support** | 4 | 12 min | 🟢 LOW |
|| 22 | **L2: Add Markdown Extensions (GFM)** | 4 | 12 min | 🟢 LOW |
|| 23 | **L3: Add ANSI Color Support** | 4 | 12 min | 🟢 LOW |
|| 24 | **L4: Add JSON Schema Generation** | 4 | 12 min | 🟢 LOW |
|| 25 | **L10: Performance Regression CI** | 4 | 12 min | 🟢 LOW |
|| 26 | **L5: Add Plugin System** | 3 | 12 min | 🟢 LOW |
|| 27 | **L8: Streaming Benchmarks** | 3 | 12 min | 🟢 LOW |
|| 28 | **L9: Fuzzing Corpus** | 3 | 12 min | 🟢 LOW |
|| 29 | **L11: Fix Phantom Type Violations** | 3 | 12 min | 🟢 LOW |
|| 30 | **L6: WebAssembly Support** | 2 | 12 min | 🟢 LOW |
|| 31 | **L7: Internationalization** | 2 | 12 min | 🟢 LOW |
|| 32 | **L12: DOT/Mermaid Mixin** | 2 | 12 min | 🟢 LOW |

---

## Execution Order

### Phase 1: Critical Fixes (Week 1)
| Order | Task | Est. Time |
|-------|------|-----------|
| 1 | C1: Fix GraphNode.GetStyle() | 12 min |
| 2 | H6: Refactor containsString | 7 min |
| 3 | H2: Standardize Render() | 12 min |

### Phase 2: High Priority (Week 2)
| Order | Task | Est. Time |
|-------|------|-----------|
| 4 | H1: Add cmdguard to go.mod | 12 min |
| 5 | H7: Document Streaming | 12 min |
| 6 | H3: Split format_test.go | 12 min |
| 7 | H4: Split format.go | 12 min |
| 8 | H8: Fix TestTableData | 12 min |
| 9 | H9: Fix Stream complexity | 12 min |

### Phase 3: Testing Improvements (Week 3)
| Order | Task | Est. Time |
|-------|------|-----------|
| 10 | H5: Add User Journey Tests | 12 min |
| 11 | H10: Add Workflow Tests | 12 min |
| 12 | M1: Add Property-Based Tests | 12 min |

### Phase 4: Medium Priority (Week 4)
| Order | Task | Est. Time |
|-------|------|-----------|
| 13 | M8: Add Format Benchmarks | 12 min |
| 14 | M2: Add Configuration Options | 12 min |
| 15 | M3: Extract Hardcoded CSS | 12 min |
| 16 | M9: API Documentation | 12 min |
| 17 | M6: Release Workflow | 12 min |
| 18 | M4: Add TSV Format | 12 min |
| 19 | M7: Add XML Format | 12 min |
| 20 | M5: D2 Multi-Table | 12 min |

### Phase 5: Low Priority (Future)
All remaining L-series tasks can be completed in future sprints as needed.

---

## Total Estimated Time

| Phase | Tasks | Time |
|-------|-------|------|
| Phase 1: Critical | 3 | 31 min |
| Phase 2: High | 6 | 72 min |
| Phase 3: Testing | 3 | 36 min |
| Phase 4: Medium | 8 | 96 min |
| Phase 5: Low | 12 | 144 min |
| **Total** | **32** | **379 min** |

**Estimated completion: ~6 hours of focused work**

---

## Sign-off

**Planner:** Crush AI  
**Date:** 2026-03-28  
**Status:** ✅ PLAN COMPLETE - Ready for Execution

---

*Generated from comprehensive analysis of March 27, 2026 status review*
