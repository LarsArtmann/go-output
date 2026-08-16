package integration

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/markdown"
	"github.com/larsartmann/go-output/markup"
	"github.com/larsartmann/go-output/plantuml"
	"github.com/larsartmann/go-output/table"
	"github.com/larsartmann/go-output/testhelpers"
)

// TestRoundTripXML verifies Table → XML → parse → verify.
func TestRoundTripXML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	xmlBytes, err := markup.MarshalXMLFromTable(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTable: %v", err)
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

	if len(parsed.Rows[0].Cells) == 0 {
		t.Fatal("row 0 has no cells — XML round-trip dropped cell data")
	}

	if parsed.Rows[0].Cells[0].Value != "Alice" {
		t.Errorf("row 0 cell 0 = %q, want %q", parsed.Rows[0].Cells[0].Value, "Alice")
	}
}

// TestRoundTripHTML verifies Table → HTML → verify structure.
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

// TestRoundTripMarkdown verifies Table → Markdown → verify structure.
func TestRoundTripMarkdown(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	md := markdown.NewMarkdownTableFromTable(data)

	out, err := md.Render()
	if err != nil {
		t.Fatalf("Markdown Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "Name", "Markdown should contain Name header")
	testhelpers.AssertContains(t, out, "---", "Markdown should contain separator")
	testhelpers.AssertContains(t, out, "Alice", "Markdown should contain Alice row")
	testhelpers.AssertContains(t, out, "Bob", "Markdown should contain Bob row")
}

// TestRoundTripTable verifies Table → lipgloss Table → verify content.
func TestRoundTripTable(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	tbl := table.FromTable(data, table.WithColorMode(output.ColorModeNever))

	out, err := tbl.Render()
	if err != nil {
		t.Fatalf("Table Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "Name", "Table should contain header")
	testhelpers.AssertContains(t, out, "Alice", "Table should contain Alice")
	testhelpers.AssertContains(t, out, "Bob", "Table should contain Bob")
}

// TestRoundTripTree verifies Table → Tree → verify non-empty output.
func TestRoundTripTree(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	var buf bytes.Buffer

	err := output.RenderTable(data, output.FormatTree, output.RenderOptions{Writer: &buf})
	if err != nil {
		t.Fatalf("RenderTable Tree: %v", err)
	}

	if buf.String() == "" {
		t.Error("Tree output should not be empty")
	}
}

// TestRoundTripAsciiDoc verifies Table → AsciiDoc → verify structure.
func TestRoundTripAsciiDoc(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	b, err := markup.MarshalAsciiDocFromTable(data)
	if err != nil {
		t.Fatalf("MarshalAsciiDocFromTable: %v", err)
	}

	result := string(b)
	testhelpers.AssertContains(t, result, "|===", "AsciiDoc should contain table border")
	testhelpers.AssertContains(t, result, "Name", "AsciiDoc should contain header")
	testhelpers.AssertContains(t, result, "Alice", "AsciiDoc should contain Alice")
}

// TestRoundTripD2 verifies Table → D2 → verify structure.
func TestRoundTripD2(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := d2.NewD2FromTable(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("D2 Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "Alice", "D2 should contain Alice node")
	testhelpers.AssertContains(t, out, "Bob", "D2 should contain Bob node")
}

// TestRoundTripMermaid verifies Table → Mermaid → verify structure.
func TestRoundTripMermaid(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := graph.NewMermaidFromTable(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("Mermaid Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "flowchart", "Mermaid should contain flowchart keyword")
	testhelpers.AssertContains(t, out, "Alice", "Mermaid should contain Alice node")
}

// TestRoundTripDOT verifies Table → DOT → verify structure.
func TestRoundTripDOT(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := graph.NewDOTFromTable(data)

	out, err := renderer.Render()
	if err != nil {
		t.Fatalf("DOT Render: %v", err)
	}

	testhelpers.AssertContains(t, out, "digraph", "DOT should contain digraph keyword")
	testhelpers.AssertContains(t, out, "Alice", "DOT should contain Alice node")
}

// TestRoundTripPlantUML verifies Table → PlantUML → verify structure.
func TestRoundTripPlantUML(t *testing.T) {
	t.Parallel()

	data := roundTripData()

	renderer := plantuml.NewPlantUMLFromTable(data)

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

	data := output.NewTable([]string{"Item", "Qty"})
	data.AddRow([]string{"Apple", "5"})
	data.AddRow([]string{"Banana", "3"})
	data.Footer = []string{"Total", "8"}

	t.Run("csv footer round-trip", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		err := output.RenderTable(data, output.FormatCSV, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTable CSV: %v", err)
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

		err := output.RenderTable(data, output.FormatTSV, output.RenderOptions{Writer: &buf})
		if err != nil {
			t.Fatalf("RenderTable TSV: %v", err)
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
