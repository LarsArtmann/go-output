# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- `enum` package with generic enum utilities
- `Format.Category()` method for format classification
- `FormatTSV` constant for TSV output format
- TSV formatter implementation
- Map-based format classification (replaces switch statements)

### Changed

- Refactored `ParseFormat`, `ParseSortBy`, `ParseColorMode` to use `enum` helpers
- Refactored `AllowedValues()` methods to use `enum` helpers

### Deprecated

- `OutputFormat` type alias - will be removed in v2.0
- `OutputFormat*` constants - will be removed in v2.0
- `ParseOutputFormat()` function - will be removed in v2.0

### Removed

### Fixed

### Security

## [0.1.0] - 2026-01-01

### Added

- Initial release
