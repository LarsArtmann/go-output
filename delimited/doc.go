// Package delimited provides CSV and TSV output formatters for tabular data.
//
// Use NewCSVWriter or NewTSVWriter for streaming (row-by-row) output, or
// MarshalCSVFromTable / MarshalTSVFromTable for one-shot marshaling
// from output.Table.
//
// Both writers support WriteHeader, WriteRow, WriteFooter, WriteRows, Flush,
// and Error methods for full control over delimited output generation.
package delimited
