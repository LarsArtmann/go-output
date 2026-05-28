// Package output provides a reusable Go library for CLI applications offering
// consistent output formatting across 16 formats (Table, JSON, CSV, TSV,
// Markdown, XML, YAML, HTML, Tree, D2, Mermaid, DOT, JSONL, AsciiDoc, TOML,
// PlantUML) with type-safe enum-based configuration and a Shape capability matrix.
//
// # Quick Start
//
// Use TableData as the single source of truth for tabular data:
//
//	data := output.NewTableData([]string{"Name", "Status"})
//	data.AddRow([]string{"Project A", "Active"})
//	data.Footer = []string{"Total", "1"}
//
//	Render to any supported format:
//	var buf bytes.Buffer
//	_ = output.RenderTableData(data, output.FormatMarkdown, output.RenderOptions{Writer: &buf})
//
// # Architecture
//
// The library uses a multi-module workspace where each format family is an
// independent Go module. Users import only the formats they need:
//
//	import "github.com/larsartmann/go-output"                     // core + Markdown + Tree
//	import "github.com/larsartmann/go-output/table"               // terminal tables (optional)
//	import "github.com/larsartmann/go-output/serialization"       // JSON + YAML + TOML + JSONL (optional)
//	import "github.com/larsartmann/go-output/delimited"           // CSV + TSV (optional)
//	import "github.com/larsartmann/go-output/markup"              // XML + HTML + AsciiDoc (optional)
//	import "github.com/larsartmann/go-output/d2"                  // D2 diagrams (optional)
//	import "github.com/larsartmann/go-output/graph"               // DOT + Mermaid (optional)
//	import "github.com/larsartmann/go-output/plantuml"            // PlantUML diagrams (optional)
//
// The root module has ZERO imports from sub-modules, ensuring users get
// zero transitive dependencies from formats they don't use.
package output
