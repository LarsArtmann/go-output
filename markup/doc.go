// Package markup provides XML, HTML, AsciiDoc, and streaming HTML output
// formatters for tabular data.
//
// All renderers implement the output.TableRenderer interface via embedding
// output.TableDataStore. HTML output uses semantic <thead>/<tbody>/<tfoot>
// elements with CSS classes for styling. Streaming HTML writes incrementally
// to an io.Writer for memory-efficient rendering of large datasets.
package markup
