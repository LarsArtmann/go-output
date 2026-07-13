# TODO_LIST.md — go-output

**Last updated:** 2026-07-13
**Open items:** 2
**Blocked:** 0

---

## Open Items

| #   | Task                                                                                                               | Effort | Status                        |
| --- | ------------------------------------------------------------------------------------------------------------------ | ------ | ----------------------------- |
| 14  | **Community: Post to r/golang, submit to Awesome Go**                                                              | Low    | Open (needs owner account)    |
| 16  | **Cut `v1.0.0` tag** — API frozen (ADR 006); CHANGELOG + full checklist done; all v0.30.x breaking changes shipped | Low    | Prepared — awaiting owner tag |

---

## Recently Resolved (2026-07-06 — v0.30.0 arc)

Details in CHANGELOG.md and git history.

- **v0.30.0 breaking changes shipped**: 7 breaking-change commits — deletions (B1-B10), renames (C1-C8), D2 prefix drop, GraphBuilder split, full CQRS architecture
- **CQRS architecture complete**: 3 builders, immutable Graph, pure-function renderers for all 16 formats, cross-shape projections
- **Registry dispatch rewired**: all table formats stream via CQRS streaming functions (byte-for-byte identical output proven)
- **Golden-file tests**: JSON, YAML, TOML, JSONL, CSV, TSV, XML, HTML, AsciiDoc — all locked in
- **v0.30.1-v0.30.4 patch releases**: version ref fixes, Pattern B sentinel migration, documentation website

## Recently Resolved (2026-07-05 — P0-P7 brutal review)

- **6 concurrency bugs fixed** (appName race, write-lock-held-across-I/O, unsynchronized renderNotify, unbounded saveAsync goroutines, swallowed save errors, unsynchronized showParallelism)
- **10 dead exported symbols deprecated then deleted** in v0.30.0
- **5 split brains resolved** (MsgNoActivities unified, Colors global documented, Direction bridge documented)
- **Draw() complexity** reduced (cyclop 20 → under 10 via decomposition)
- **Dead-writer detection** added to InlineRenderer
- **Build() cycle detection** added

## Recently Resolved (earlier sessions)

Details in CHANGELOG.md and git history.

- Pattern B versioning migration (all sibling deps use v0.0.0 sentinel + replace)
- enum/ + envdetect/ merged into root
- NOM BuildFlow integration (ActivityProgress, ActivityRetrying, EstimatedTotalRemaining)
- DAG topology overhaul (true DAG, layered display, critical-path analysis)
- Theme system, activity categories, parallelism meter, status registry
- daghtml module (zero-dep SVG DAG visualization)
- Split-brain elimination (20/20 findings resolved)
- FormatDuration bug fix, escaping vulnerability fixes, branded IDs, sealed Event sum type
