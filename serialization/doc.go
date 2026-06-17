// Package serialization provides JSON, YAML, TOML, and JSONL output formatters
// for tabular, tree, and graph data.
//
// Table renderers implement output.TableRenderer via embedding output.TableDataStore.
// Tree and graph renderers implement output.TreeOutputRenderer and output.GraphRenderer
// respectively, marshaling data structures using the appropriate encoding library.
//
// Use MarshalJSON, MarshalYAML, or MarshalTOML for one-shot marshaling of any
// value, or MarshalTOMLFromTableData / MarshalJSONLFromTableData for TableData.
// JSON and YAML TableData use the renderer-based dispatch (JSONTableRenderer,
// YAMLTableRenderer) via the output.RenderTableData registry.
package serialization
