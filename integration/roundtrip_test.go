package integration

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-faster/yaml"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/testhelpers"
)

// roundTripData returns a consistent Table for round-trip testing.
func roundTripData() *output.Table {
	data := output.NewTable([]string{"Name", "Score", "Active"})
	data.AddRow([]string{"Alice", "95", "true"})
	data.AddRow([]string{"Bob", "87", "false"})

	return data
}

// renderTable is a small helper that renders Table to the given format
// and returns the output string. It fails the test on error.
func renderTable(t *testing.T, data *output.Table, format output.Format) string {
	t.Helper()

	var buf bytes.Buffer

	if err := output.RenderTable(data, format, output.RenderOptions{Writer: &buf}); err != nil {
		t.Fatalf("RenderTable(%s): %v", format, err)
	}

	return buf.String()
}

// parseDelimited is a small helper that parses CSV/TSV-like output into records.
func parseDelimited(t *testing.T, input string, comma rune) [][]string {
	t.Helper()

	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = comma

	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("parseDelimited (comma=%q): %v", comma, err)
	}

	return records
}

// assertMapRow checks that parsed[row] contains the given key/value pairs.
// Fails the test on the first mismatch. Use for round-trip verification of
// parsed structured output. kv must contain an even number of arguments
// (key, value pairs).
func assertMapRow(t *testing.T, parsed []map[string]string, row int, kv ...string) {
	t.Helper()

	for i := 0; i+1 < len(kv); i += 2 {
		key, want := kv[i], kv[i+1]
		if got := parsed[row][key]; got != want {
			t.Errorf("row %d = %v, want %s=%q", row, parsed[row], key, want)
		}
	}
}

func assertCell(t *testing.T, name string, records [][]string, row, col int, want string) {
	t.Helper()

	if row >= len(records) {
		t.Errorf("%s: row %d out of range (have %d)", name, row, len(records))

		return
	}

	if col >= len(records[row]) {
		t.Errorf("%s: col %d out of range for row %d (have %d)", name, col, row, len(records[row]))

		return
	}

	if records[row][col] != want {
		t.Errorf("%s: records[%d][%d] = %q, want %q", name, row, col, records[row][col], want)
	}
}

// TestRoundTripJSON verifies Table → JSON renderer → parse → verify.
func TestRoundTripJSON(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := serialization.NewJSONTableRenderer()
	renderer.SetData(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("JSON Render: %v", err)
	}

	var parsed []map[string]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("Parse JSON: %v", err)
	}

	if len(parsed) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(parsed))
	}

	assertMapRow(t, parsed, 0, "Name", "Alice", "Score", "95")
	assertMapRow(t, parsed, 1, "Name", "Bob", "Active", "false")
}

// TestRoundTripCSV verifies Table → CSV → parse → verify.
func TestRoundTripCSV(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	result := renderTable(t, data, output.FormatCSV)
	records := parseDelimited(t, result, ',')

	if len(records) != 3 {
		t.Fatalf("expected 3 records (header + 2 rows), got %d", len(records))
	}

	assertCell(t, "CSV header[0]", records, 0, 0, "Name")
	assertCell(t, "CSV header[1]", records, 0, 1, "Score")
	assertCell(t, "CSV header[2]", records, 0, 2, "Active")
	assertCell(t, "CSV row[1][0]", records, 1, 0, "Alice")
	assertCell(t, "CSV row[1][1]", records, 1, 1, "95")
	assertCell(t, "CSV row[2][0]", records, 2, 0, "Bob")
	assertCell(t, "CSV row[2][2]", records, 2, 2, "false")
}

// TestRoundTripTSV verifies Table → TSV → parse → verify.
func TestRoundTripTSV(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	result := renderTable(t, data, output.FormatTSV)
	records := parseDelimited(t, result, '\t')

	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}

	assertCell(t, "TSV header[0]", records, 0, 0, "Name")
	assertCell(t, "TSV row[1][0]", records, 1, 0, "Alice")
}

// TestRoundTripYAML verifies Table → YAML → parse → verify.
func TestRoundTripYAML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	var buf bytes.Buffer

	err := output.RenderTable(data, output.FormatYAML, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable YAML: %v", err)
	}

	var parsed []map[string]string
	if err := yaml.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Fatalf("Parse YAML: %v", err)
	}

	if len(parsed) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(parsed))
	}

	if parsed[0]["Name"] != "Alice" {
		t.Errorf("row 0 Name = %q, want %q", parsed[0]["Name"], "Alice")
	}
}

// TestRoundTripTOML verifies Table → TOML → verify content.
// TOML does not support top-level arrays, so rows are wrapped under a [[row]]
// key using array-of-tables syntax.
func TestRoundTripTOML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	b, err := serialization.MarshalTOMLFromTable(data)
	if err != nil {
		t.Fatalf("MarshalTOMLFromTable: %v", err)
	}

	result := string(b)

	if !strings.HasPrefix(result, "[[row]]") {
		t.Errorf("TOML should start with [[row]] array-of-tables syntax, got: %s", result[:min(len(result), 50)])
	}

	testhelpers.AssertContains(t, result, "Name = 'Alice'", "TOML should contain Alice row")
	testhelpers.AssertContains(t, result, "Score = '95'", "TOML should contain Score 95")
	testhelpers.AssertContains(t, result, "Active = 'false'", "TOML should contain Active false")
	testhelpers.AssertContains(t, result, "Name = 'Bob'", "TOML should contain Bob row")
}

// TestRoundTripJSONL verifies Table → JSONL → parse → verify.
func TestRoundTripJSONL(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	var buf bytes.Buffer

	err := output.RenderTable(data, output.FormatJSONL, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable JSONL: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 JSONL lines, got %d", len(lines))
	}

	var row0 map[string]string
	if err := json.Unmarshal([]byte(lines[0]), &row0); err != nil {
		t.Fatalf("Parse JSONL line 0: %v", err)
	}

	if row0["Name"] != "Alice" {
		t.Errorf("line 0 Name = %q, want %q", row0["Name"], "Alice")
	}

	var row1 map[string]string
	if err := json.Unmarshal([]byte(lines[1]), &row1); err != nil {
		t.Fatalf("Parse JSONL line 1: %v", err)
	}

	if row1["Name"] != "Bob" {
		t.Errorf("line 1 Name = %q, want %q", row1["Name"], "Bob")
	}
}
