# go-output

Reusable Go library for CLI applications providing consistent output formatting across multiple formats with type-safe enum-based configuration.

## Purpose

Standardizes output formatting across personal Go projects:

- `/Users/larsartmann/projects/project-meta/`
- `/Users/larsartmann/projects/projects-management-automation/`

## Supported Formats

| Format     | Description                           |
| ---------- | ------------------------------------- |
| `table`    | Terminal tables with lipgloss styling |
| `json`     | JSON output with indentation          |
| `csv`      | CSV export with headers               |
| `markdown` | Markdown tables                       |
| `d2`       | D2 diagram shapes                     |
| `yaml`     | YAML serialization                    |

## Supported Sort Options

| Option       | Description               |
| ------------ | ------------------------- |
| `name`       | Sort by name              |
| `importance` | Sort by importance level  |
| `created_at` | Sort by creation date     |
| `updated_at` | Sort by last update       |
| `health`     | Sort by health score      |
| `complexity` | Sort by complexity metric |

## Color Modes

| Mode     | Description                                    |
| -------- | ---------------------------------------------- |
| `auto`   | Respect `NO_COLOR`, CI env vars, TTY detection |
| `always` | Force ANSI colors                              |
| `never`  | Disable colors                                 |

## CLI Flag Integration

Integrates with [cmdguard](https://github.com/larsartmann/cmdguard) for type-safe flags:

```go
type ListFlags struct {
    Format output.OutputFormat `flag:"format" default:"table" help:"Output format (table, json, csv, markdown, d2, yaml)"`
    SortBy output.SortBy       `flag:"sort-by" default:"name" help:"Sort by (name, importance, created_at, updated_at, health, complexity)"`
    Color  output.ColorMode    `flag:"color" default:"auto" help:"Color mode (auto, always, never)"`
}
```

Flags validate against allowed values and provide bash/zsh completion.

## Installation

```bash
go get github.com/larsartmann/go-output
```

## Dependencies

```go
require (
    github.com/charmbracelet/lipgloss/v2 v2
    go.yaml.in/yaml/v4 v4
    github.com/larsartmann/cmdguard/v2 v2
)
```

## Development

```bash
# Build
just build

# Test
just test

# Lint
just lint
```

## License

See [LICENSE](LICENSE) file for details.
