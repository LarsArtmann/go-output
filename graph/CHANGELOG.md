# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

## [0.9.0] - 2026-06-12

### Added

- DOT and Mermaid now register as `TableDataMarshaler` via `init()` — enables `RenderTableData()` dispatch for `FormatDOT` and `FormatMermaid`.
- Registry dispatch tests for DOT and Mermaid table data rendering.

### Deprecated

- `NewGraphNodeID` / `NewGraphNodeLabel` — use `output.NewBrandedID` directly.

### Changed

- `dt.Build()` errors now propagated instead of silently discarded.

## [0.7.0] - 2026-06-09

### Added

- Initial changelog entry.

## [0.1.0] - 2026-01-01

### Added

- Initial release
