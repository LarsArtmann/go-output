// Package serialization provides JSON, YAML, TOML, and JSONL output formatters
// for tabular, tree, and graph data.
//
// Table renderers implement output.TableRenderer via embedding output.TableDataStore.
// Tree and graph renderers implement output.TreeOutputRenderer and output.GraphRenderer
// respectively, marshaling data structures using the appropriate encoding library.
//
// Use MarshalJSONFromTableData, MarshalYAMLFromTableData, etc. for one-shot
// marshaling, or use the typed renderer constructors for incremental building.
package serialization
