package integration

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/go-faster/yaml"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/table"
	"github.com/larsartmann/go-output/testhelpers"
)

// roundTripData returns a consistent TableData for round-trip testing.
func roundTripData() *output.TableData {
	data := output.NewTableData([]string{"Name", "Score", "Active"})
	data.AddRow([]string{"Alice", "95", "true"})
	data.AddRow([]string{"Bob", "87", "false"})

	return data
}

// renderTableData is a small helper that renders TableData to the given format
// and returns the output string. It fails the test on error.
func renderTableData(t *testing.T, data *output.TableData, format output.Format) string {
	t.Helper()

	var buf bytes.Buffer

	if err := output.RenderTableData(data, format, output.RenderOptions{Writer: &buf}); err != nil {
		t.Fatalf("RenderTableData(%s): %v", format, err)
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

// TestRoundTripJSON verifies TableData → JSON renderer → parse → verify.
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

// TestRoundTripCSV verifies TableData → CSV → parse → verify.
func TestRoundTripCSV(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	result := renderTableData(t, data, output.FormatCSV)
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

// TestRoundTripTSV verifies TableData → TSV → parse → verify.
func TestRoundTripTSV(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	result := renderTableData(t, data, output.FormatTSV)
	records := parseDelimited(t, result, '\t')

	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}

	assertCell(t, "TSV header[0]", records, 0, 0, "Name")
	assertCell(t, "TSV row[1][0]", records, 1, 0, "Alice")
}

// TestRoundTripYAML verifies TableData → YAML → parse → verify.
func TestRoundTripYAML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	var buf bytes.Buffer

	err := output.RenderTableData(data, output.FormatYAML, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData YAML: %v", err)
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

// TestRoundTripTOML verifies TableData → TOML → verify content.
// TOML does not support top-level arrays, so rows are wrapped under a [[row]]
// key using array-of-tables syntax.
func TestRoundTripTOML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	b, err := serialization.MarshalTOMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalTOMLFromTableData: %v", err)
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

// TestRoundTripJSONL verifies TableData → JSONL → parse → verify.
func TestRoundTripJSONL(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	var buf bytes.Buffer

	err := output.RenderTableData(data, output.FormatJSONL, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData JSONL: %v", err)
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

// TestRoundTripXML verifies TableData → XML → parse → verify.
func TestRoundTripXML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	xmlBytes, err := markup.MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData: %v", err)
	}

	type xmlColumn struct {
		XMLName xml.Name `xml:"column"`
		Value   string   `xml:",chardata"`
	}

	type xmlCell struct {
		XMLName xml.Name `xml:"cell"`
		Value   string   `xml:",chardata"`
	}

	type xmlRow struct {
		XMLName xml.Name  `xml:"row"`
		Cells   []xmlCell `xml:"cell"`
	}

	type xmlTable struct {
		XMLName xml.Name    `xml:"table"`
		Headers []xmlColumn `xml:"headers>column"`
		Rows    []xmlRow    `xml:"rows>row"`
	}

	var parsed xmlTable
	if err := xml.Unmarshal(xmlBytes, &parsed); err != nil {
		t.Fatalf("Parse XML: %v\nXML: %s", err, string(xmlBytes))
	}

	if len(parsed.Headers) != 3 {
		t.Errorf("header count = %d, want 3: %v", len(parsed.Headers), parsed.Headers)
	}

	if parsed.Headers[0].Value != "Name" {
		t.Errorf("header[0] = %q, want %q", parsed.Headers[0].Value, "Name")
	}

	if len(parsed.Rows) != 2 {
		t.Fatalf("row count = %d, want 2", len(parsed.Rows))
	}

	if len(parsed.Rows[0].Cells) > 0 && parsed.Rows[0].Cells[0].Value != "Alice" {
		t.Errorf("row 0 cell 0 = %q, want %q", parsed.Rows[0].Cells[0].Value, "Alice")
	}
}

// TestRoundTripHTML verifies TableData → HTML → verify structure.
func TestRoundTripHTML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := markup.NewHTMLRenderer()
	renderer.SetData(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("HTML Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "<table", "HTML should contain table tag")
	testhelpers.AssertContains(t, out, "<th>Name</th>", "HTML should contain Name header")
	testhelpers.AssertContains(t, out, "<td>Alice</td>", "HTML should contain Alice cell")
	testhelpers.AssertContains(t, out, "<td>Bob</td>", "HTML should contain Bob cell")
	testhelpers.AssertContains(t, out, "</table>", "HTML should close table tag")
}

// TestRoundTripMarkdown verifies TableData → Markdown → verify structure.
func TestRoundTripMarkdown(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	md := output.NewMarkdownTableFromData(data)

	out, err := md.Render()
	if err != nil {
		t.Fatalf("Markdown Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "Name", "Markdown should contain Name header")
	testhelpers.AssertContains(t, out, "---", "Markdown should contain separator")
	testhelpers.AssertContains(t, out, "Alice", "Markdown should contain Alice row")
	testhelpers.AssertContains(t, out, "Bob", "Markdown should contain Bob row")
}

// TestRoundTripTable verifies TableData → lipgloss Table → verify content.
func TestRoundTripTable(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	tbl := table.FromTableData(data, table.WithColorMode(output.ColorModeNever))

	out, err := tbl.Render()
	if err != nil {
		t.Fatalf("Table Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "Name", "Table should contain header")
	testhelpers.AssertContains(t, out, "Alice", "Table should contain Alice")
	testhelpers.AssertContains(t, out, "Bob", "Table should contain Bob")
}

// TestRoundTripTree verifies TableData → Tree → verify non-empty output.
func TestRoundTripTree(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	var buf bytes.Buffer

	err := output.RenderTableData(data, output.FormatTree, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTableData Tree: %v", err)
	}

	if buf.String() == "" {
		t.Error("Tree output should not be empty")
	}
}

// TestRoundTripAsciiDoc verifies TableData → AsciiDoc → verify structure.
func TestRoundTripAsciiDoc(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	b, err := markup.MarshalAsciiDocFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalAsciiDocFromTableData: %v", err)
	}

	result := string(b)
	testhelpers.AssertContains(t, result, "|===", "AsciiDoc should contain table border")
	testhelpers.AssertContains(t, result, "Name", "AsciiDoc should contain header")
	testhelpers.AssertContains(t, result, "Alice", "AsciiDoc should contain Alice")
}

// TestRoundTripD2 verifies TableData → D2 → verify structure.
func TestRoundTripD2(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := d2.D2FromTableData(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("D2 Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "Alice", "D2 should contain Alice node")
	testhelpers.AssertContains(t, out, "Bob", "D2 should contain Bob node")
}

// TestRoundTripMermaid verifies TableData → Mermaid → verify structure.
func TestRoundTripMermaid(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := graph.MermaidFromTableData(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Mermaid Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "flowchart", "Mermaid should contain flowchart keyword")
	testhelpers.AssertContains(t, out, "Alice", "Mermaid should contain Alice node")
}

// TestRoundTripDOT verifies TableData → DOT → verify structure.
func TestRoundTripDOT(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := graph.DOTFromTableData(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("DOT Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "digraph", "DOT should contain digraph keyword")
	testhelpers.AssertContains(t, out, "Alice", "DOT should contain Alice node")
}

// TestRoundTripPlantUML verifies TableData → PlantUML → verify structure.
func TestRoundTripPlantUML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := plantuml.PlantUMLFromTableData(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("PlantUML Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "@startuml", "PlantUML should contain start tag")
	testhelpers.AssertContains(t, out, "Alice", "PlantUML should contain Alice node")
	testhelpers.AssertContains(t, out, "@enduml", "PlantUML should contain end tag")
}

// TestRoundTripFooter verifies footer data survives round-trip in parseable formats.
func TestRoundTripFooter(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Item", "Qty"})
	data.AddRow([]string{"Apple", "5"})
	data.AddRow([]string{"Banana", "3"})
	data.Footer = []string{"Total", "8"}

	t.Run("csv footer round-trip", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		err := output.RenderTableData(data, output.FormatCSV, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTableData CSV: %v", err)
		}

		reader := csv.NewReader(&buf)

		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Parse CSV: %v", err)
		}

		if len(records) != 4 {
			t.Fatalf("expected 4 records (header + 2 rows + footer), got %d", len(records))
		}

		if records[3][0] != "Total" || records[3][1] != "8" {
			t.Errorf("footer = %v, want [Total 8]", records[3])
		}
	})

	t.Run("tsv footer round-trip", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		err := output.RenderTableData(data, output.FormatTSV, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTableData TSV: %v", err)
		}

		reader := csv.NewReader(&buf)
		reader.Comma = '\t'

		records, err := reader.ReadAll()
		if err != nil {
			t.Fatalf("Parse TSV: %v", err)
		}

		if len(records) != 4 {
			t.Fatalf("expected 4 records, got %d", len(records))
		}

		if records[3][0] != "Total" {
			t.Errorf("footer[0] = %q, want %q", records[3][0], "Total")
		}
	})
}
