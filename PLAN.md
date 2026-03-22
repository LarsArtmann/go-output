# go-output Library Plan

## Overview

A reusable Go library for CLI applications that provides consistent output formatting across multiple formats (JSON, CSV, Markdown, D2, YAML, Table) with type-safe enum-based configuration and optional cmdguard integration.

**Location:** `/Users/larsartmann/projects/go-output/`

---

## Core Enums

| Enum | Values | Purpose |
|------|--------|---------|
| `OutputFormat` | `table`, `json`, `csv`, `markdown`, `d2`, `yaml` | Output format selection |
| `SortBy` | `name`, `importance`, `created_at`, `updated_at`, `health`, `complexity` | Sort field selection |
| `ColorMode` | `auto`, `always`, `never` | Color output control |

---

## Package Structure

```
go-output/
├── go.mod
├── justfile
├── README.md
│
├── format.go                    # OutputFormat enum
├── sort.go                     # SortBy enum + sorting helpers
├── color.go                    # ColorMode enum + ANSI helpers
│
├── json.go                     # JSON formatter
├── csv.go                      # CSV formatter
├── yaml.go                     # YAML formatter
├── markdown.go                 # Markdown formatter
├── d2.go                       # D2 diagram formatter
│
├── table/
│   ├── table.go                # Table interface & generic implementation
│   ├── config.go               # TableConfig interface
│   ├── lipgloss.go             # Lipgloss-based table builder
│   └── styles.go               # Pre-defined styles
│
├── sort/
│   ├── sorter.go               # Generic sorting interface
│   ├── comparators.go          # Type-specific comparators
│   └── adapter.go              # Sort adapter for cmdguard integration
│
├── cmdguard/
│   ├── format.go               # OutputFormat flag with validation
│   ├── sort.go                 # SortBy flag with validation
│   └── color.go                # ColorMode flag with validation
│
└── examples/
    ├── basic/main.go
    ├── table/main.go
    └── cmdguard/main.go
```

---

## Phase 1: Foundation

| Step | Task | Details |
|------|------|---------|
| 1.1 | Create repository | `mkdir -p /Users/larsartmann/projects/go-output` |
| 1.2 | Initialize go.mod | Module: `github.com/larsartmann/go-output` |
| 1.3 | Create `format.go` | `OutputFormat` enum with Parse, String, AllowedValues |
| 1.4 | Create `sort.go` | `SortBy` enum with Parse, String, AllowedValues |
| 1.5 | Create `color.go` | `ColorMode` enum with Parse, ShouldColor, ToANSI |

### format.go Specification

```go
type OutputFormat int

const (
    OutputFormatTable OutputFormat = iota
    OutputFormatJSON
    OutputFormatCSV
    OutputFormatMarkdown
    OutputFormatD2
    OutputFormatYAML
)

func ParseOutputFormat(s string) (OutputFormat, error)
func (f OutputFormat) String() string
func (f OutputFormat) AllowedValues() []string
func (f OutputFormat) IsValid() bool
```

### sort.go Specification

```go
type SortBy string

const (
    SortByName SortBy = "name"
    SortByImportance SortBy = "importance"
    SortByCreatedAt SortBy = "created_at"
    SortByUpdatedAt SortBy = "updated_at"
    SortByHealth SortBy = "health"
    SortByComplexity SortBy = "complexity"
)

func ParseSortBy(s string) (SortBy, error)
func (s SortBy) String() string
func (s SortBy) AllowedValues() []string
```

### color.go Specification

```go
type ColorMode string

const (
    ColorModeAuto ColorMode = "auto"
    ColorModeAlways ColorMode = "always"
    ColorModeNever ColorMode = "never"
)

func ParseColorMode(s string) (ColorMode, error)
func (c ColorMode) ShouldColor() bool  // Respects NO_COLOR, CI env vars, TTY
func (c ColorMode) ToANSI() string     // Returns escape code or ""
```

---

## Phase 2: Core Formatters

| Step | Task | Details |
|------|------|---------|
| 2.1 | Create `json.go` | Marshal/Indent, encoder options |
| 2.2 | Create `csv.go` | NewWriter, WriteHeader, WriteRow |
| 2.3 | Create `yaml.go` | Marshal/Unmarshal using go.yaml |
| 2.4 | Create `markdown.go` | Table generator with alignment |
| 2.5 | Create `d2.go` | SQL table shape generator |

### json.go Specification

```go
func MarshalJSON(v any) ([]byte, error)
func MarshalJSONIndent(v any, prefix, indent string) ([]byte, error)
type JSONWriter struct { Writer io.Writer }
func (j *JSONWriter) Encode(v any) error
```

### csv.go Specification

```go
type CSVWriter struct { Writer io.Writer }
func NewCSVWriter(w io.Writer) *CSVWriter
func (c *CSVWriter) WriteHeader(cols []string) error
func (c *CSVWriter) WriteRow(values []string) error
func (c *CSVWriter) Flush()
```

---

## Phase 3: Table System

| Step | Task | Details |
|------|------|---------|
| 3.1 | Create `table/config.go` | TableConfig interface |
| 3.2 | Create `table/styles.go` | BorderStyle, Color, Alignment |
| 3.3 | Create `table/lipgloss.go` | Lipgloss table builder |
| 3.4 | Create `table/table.go` | Generic table interface |

### Table Interface

```go
type TableConfig interface {
    GetBorderStyle() BorderStyle
    GetHeaderStyle() TableCellStyle
    GetRowStyle() TableCellStyle
    GetAlternateRowStyle() TableCellStyle
    ShouldUseAlternateRows() bool
    GetColumnConfigs() []TableColumnConfig
    GetBorderColor() Color
    GetPadding() (int, int)
    ShouldShowBorders() bool
}

type Table interface {
    SetHeaders([]string) error
    AddRow([]string) error
    Render() (string, error)
}

func NewTable(config TableConfig) Table
func (t *Table) WithStyleFunc(fn func(row, col int) Style) Table
```

---

## Phase 4: Sorting System

| Step | Task | Details |
|------|------|---------|
| 4.1 | Create `sort/comparators.go` | CompareString, CompareInt, CompareTime |
| 4.2 | Create `sort/sorter.go` | SortFunc adapter |
| 4.3 | Create `sort/adapter.go` | SortBy to Comparator conversion |

### Sorting Interface

```go
type Comparator func(a, b any) int

type Sorter struct {
    Items any       // slice to sort
    By    SortBy    // sort field
    Desc  bool      // descending order
}

func (s *Sorter) Sort() error
func (s *Sorter) SortFunc(cmp Comparator) *Sorter
```

---

## Phase 5: cmdguard Integration

| Step | Task | Details |
|------|------|---------|
| 5.1 | Create `cmdguard/format.go` | OutputFormat flag tag support |
| 5.2 | Create `cmdguard/sort.go` | SortBy flag tag support |
| 5.3 | Create `cmdguard/color.go` | ColorMode flag with auto-detection |

### cmdguard Usage

```go
import "github.com/larsartmann/go-output/cmdguard"

type MyFlags struct {
    Format output.OutputFormat `flag:"format" default:"table"`
    SortBy output.SortBy       `flag:"sort-by" default:"name" help:"Sort by (name, importance)"`
    Color  output.ColorMode    `flag:"color" default:"auto"`
}

// Flag registry automatically:
// - Validates against allowed values
// - Provides bash/zsh completion
// - Shows help with all options
```

---

## Phase 6: Polish

| Step | Task | Details |
|------|------|---------|
| 6.1 | Create `justfile` | build, test, lint, format |
| 6.2 | Add tests | Unit tests for all formatters and enums |
| 6.3 | Create README | Usage documentation with examples |
| 6.4 | Add examples | Basic, table, and cmdguard usage |

---

## Dependencies

```go
require (
    github.com/charmbracelet/lipgloss/v2 v2
    go.yaml.in/yaml/v4 v4
    github.com/larsartmann/cmdguard/v2 v2
)
```

---

## Usage Examples

### Basic Usage

```go
import "github.com/larsartmann/go-output"

func main() {
    data := []User{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}

    // JSON
    jsonBytes, _ := output.MarshalJSONIndent(data, "", "  ")
    fmt.Println(string(jsonBytes))

    // CSV
    writer := output.NewCSVWriter(os.Stdout)
    writer.WriteHeader([]string{"Name", "Age"})
    writer.WriteRow([]string{"Alice", "30"})
    writer.Flush()

    // Markdown
    md := output.MarkdownTable([]string{"Name", "Age"}, data)
    fmt.Println(md)
}
```

### Table with Sorting

```go
table := output.NewTable(config).
    SetHeaders([]string{"Name", "Importance", "Health"})

for _, p := range projects {
    table.AddRow([]string{p.Name, p.Importance.String(), p.Health.String()})
}

sorted, _ := table.Render()
fmt.Println(sorted)
```

### cmdguard Integration

```go
import "github.com/larsartmann/go-output/cmdguard"

type ListFlags struct {
    Format output.OutputFormat `flag:"format" default:"table"`
    SortBy output.SortBy       `flag:"sort-by" default:"name"`
    Color  output.ColorMode    `flag:"color" default:"auto"`
}

cmd := v2.Command[AppConfig, *ListFlags]{
    Use:   "list",
    Short: "List projects",
    Flags: &ListFlags{},
    RunE: func(ctx context.Context, cfg *AppConfig, flags *ListFlags) error {
        projects := getProjects()

        sorter := output.Sort(projects, flags.SortBy, true)
        sorter.Sort()

        table := output.NewTable(tableConfig).
            SetHeaders([]string{"Name", "Importance"})

        for _, p := range projects {
            table.AddRow([]string{p.Name, p.Importance.String()})
        }

        output.Print(table, flags.Format, flags.Color.ShouldColor())
        return nil
    },
}
```

---

## Motivation

This library solves the problem of scattered, inconsistent output formatting across CLI applications:

| Before | After |
|--------|-------|
| Inline JSON/CSV/Markdown in each command | Shared, tested formatters |
| String-based format validation | Type-safe enums |
| Manual color detection | Auto ColorMode with env var support |
| Custom sorting logic per command | Generic SortBy with Comparator adapter |
| No cmdguard integration | Flags with validation and completion |

---

## Migration Path

1. **Phase 1-2:** Extract formatters from `projects-management-automation` and `project-meta`
2. **Phase 3-4:** Add table system with sorting
3. **Phase 5:** Integrate with cmdguard for type-safe flags
4. **Phase 6:** Polish and document

---

## Status

- [ ] Phase 1: Foundation
- [ ] Phase 2: Core Formatters
- [ ] Phase 3: Table System
- [ ] Phase 4: Sorting System
- [ ] Phase 5: cmdguard Integration
- [ ] Phase 6: Polish
