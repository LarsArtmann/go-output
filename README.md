# go-output

Reusable Go library for CLI applications providing consistent output formatting across multiple formats with type-safe enum-based configuration.

## Purpose

Standardizes output formatting across personal Go projects:

- `/Users/larsartmann/projects/project-meta/`
- `/Users/larsartmann/projects/projects-management-automation/`

## Quick Start

```go
import "github.com/larsartmann/go-output"

// JSON output
data, _ := output.MarshalJSONIndent(projects, "", "  ")
fmt.Println(string(data))

// Markdown table
md := output.NewMarkdownTable()
md.SetHeaders([]string{"Name", "Health", "Complexity"})
md.AddRow([]string{"Alpha", "90%", "7/10"})
out, _ := md.Render()
fmt.Println(out)

// CSV output
w := output.NewCSVWriter(os.Stdout)
w.WriteHeader([]string{"Name", "Value"})
w.WriteRow([]string{"Item", "123"})
w.Flush()
```

## Supported Formats

| Format     | Description                           | Package                                  |
| ---------- | ------------------------------------- | ---------------------------------------- |
| `table`    | Terminal tables with lipgloss styling | `github.com/larsartmann/go-output/table` |
| `json`     | JSON output with indentation          | `github.com/larsartmann/go-output`       |
| `csv`      | CSV export with headers               | `github.com/larsartmann/go-output`       |
| `markdown` | Markdown tables                       | `github.com/larsartmann/go-output`       |
| `d2`       | D2 diagram shapes                     | `github.com/larsartmann/go-output`       |
| `yaml`     | YAML serialization                    | `github.com/larsartmann/go-output`       |

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

## Type-Safe Enums

All configuration types provide validation and string conversion:

```go
// Parse with validation
format, err := output.ParseOutputFormat("json")
if err != nil {
    // handle error
}

// Check validity
if format.IsValid() {
    fmt.Println(format.String()) // "json"
}

// Get allowed values for CLI help
allowed := format.AllowedValues() // []string{"table", "json", "csv", ...}
```

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

# Test (includes benchmarks and fuzz tests)
just test

# Lint
just lint

# Full verification
just verify

# Run example
go run ./examples/basic/main.go markdown
```

## Examples

See [`examples/basic/main.go`](examples/basic/main.go) for a complete example demonstrating all formats.

## License

See [LICENSE](LICENSE) file for details.
