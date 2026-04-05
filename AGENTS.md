# go-output - Agent Instructions

## Project Overview

A reusable Go library for CLI applications providing consistent output formatting across multiple formats (JSON, CSV, Markdown, D2, YAML, Table) with type-safe enum-based configuration.

## Location

`/Users/larsartmann/projects/go-output/`

## Repository

https://github.com/larsartmann/go-output

## Key Technologies

- Go 1.26+
- charm.land/lipgloss/v2 (terminal styling)
- github.com/go-faster/yaml (YAML support)
- github.com/larsartmann/cmdguard/v2 (optional CLI flag integration - add separately)

## Project Structure

```
go-output/
├── format.go           # OutputFormat enum
├── sort.go            # SortBy enum
├── color.go           # ColorMode enum
├── json.go            # JSON formatter
├── csv.go             # CSV formatter
├── yaml.go            # YAML formatter
├── markdown.go        # Markdown formatter
├── d2.go              # D2 diagram formatter
├── dot.go             # DOT graph formatter
├── mermaid.go         # Mermaid diagram formatter
├── tree.go            # Tree rendering
├── table/             # Table interface & implementation
├── sort/              # Sorting utilities
├── cmdguard/          # cmdguard flag integration
└── examples/          # Usage examples
```

## Build Commands

```bash
# Build
just build

# Test
just test

# Lint
just lint

# Full verification
just verify

# BuildFlow (comprehensive)
buildflow --semantic --fix
```

## Code Quality Standards

- All code must pass `golangci-lint` with `.golangci.yml` configuration
- Tests required for all public APIs
- 90%+ test coverage target
- File size limit: 350 lines per file
- No code duplication (threshold: 30 tokens)

## Testing

```bash
# Unit tests
go test ./...

# Race detector
go test -race ./...

# Coverage
go test -cover ./...

# Fuzz tests
go test -fuzz=FuzzParseOutputFormat -fuzztime=1m .
go test -fuzz=FuzzParseSortBy -fuzztime=1m .
```

## Key Design Patterns

1. **Type-safe enums**: Use string/int constants with Parse/Validate methods
2. **Functional options**: For optional configuration
3. **Interface-based design**: Table, Renderer interfaces
4. **Generic functions**: Sort utilities use Go 1.18+ generics

## Common Tasks

### Adding a New Output Format

1. Add format constant to `format.go`
2. Implement formatter function in new file
3. Add tests with >90% coverage
4. Update cmdguard integration if needed

### Adding a New Sort Field

1. Add constant to `sort.go`
2. Add comparator in `sort/sort.go` if needed
3. Update tests

## Dependencies

See `go.mod` for full list. Key dependencies:

- `charm.land/lipgloss/v2` - Terminal styling
- `github.com/go-faster/yaml` - YAML marshaling
- `github.com/larsartmann/cmdguard/v2` - CLI integration

## Notes for AI Agents

- This is a library, not an application - no main function
- Focus on clean APIs over implementation complexity
- Prefer composition over inheritance
- All exported functions must have documentation comments
- Follow Go conventions (gofumpt, goimports formatting)
- Consider backward compatibility when changing APIs

## Contact

For questions or issues, refer to the repository or project documentation.
