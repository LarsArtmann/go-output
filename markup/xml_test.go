package markup

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestXMLWriterWriteHeader(t *testing.T) {
	t.Parallel()

	var buf strings.Builder

	x := NewXMLWriter(&buf)

	err := x.WriteHeader([]string{"Name", "Value"})
	if err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}

	err = x.WriteFooter()
	if err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	result := buf.String()
	assertContains(t, result, "<headers>", "XML should contain <headers>")
	assertContains(t, result, "<column>Name</column>", "XML should contain <column>Name</column>")
}

func TestXMLWriterWriteRow(t *testing.T) {
	t.Parallel()

	var buf strings.Builder

	x := NewXMLWriter(&buf)

	_ = x.WriteHeader([]string{"Name", "Value"})

	err := x.WriteRow([]string{"test", "123"})
	if err != nil {
		t.Fatalf("WriteRow() error = %v", err)
	}

	err = x.WriteFooter()
	if err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	result := buf.String()
	assertContains(t, result, "<row>", "XML should contain <row>")
	assertContains(t, result, "<cell>test</cell>", "XML should contain <cell>test</cell>")
}

func TestXMLWriterWriteRows(t *testing.T) {
	t.Parallel()

	var buf strings.Builder

	x := NewXMLWriter(&buf)

	_ = x.WriteHeader([]string{"Name", "Value"})

	err := x.WriteRows([][]string{
		{"a", "1"},
		{"b", "2"},
	})
	if err != nil {
		t.Fatalf("WriteRows() error = %v", err)
	}

	err = x.WriteFooter()
	if err != nil {
		t.Fatalf("WriteFooter() error = %v", err)
	}

	result := buf.String()
	if strings.Count(result, "<row>") != 2 {
		t.Errorf("XML should contain 2 <row> elements")
	}
}

func TestXMLWriterEscape(t *testing.T) {
	t.Parallel()

	var buf strings.Builder

	x := NewXMLWriter(&buf)

	_ = x.WriteHeader([]string{"Name"})
	_ = x.WriteRow([]string{"<script>alert('xss')</script>"})
	_ = x.WriteFooter()

	result := buf.String()
	if strings.Contains(result, "<script>") {
		t.Error("XML should escape <script> tags")
	}

	assertContains(t, result, "&lt;script&gt;", "XML should contain escaped &lt;script&gt;")
}

func TestMarshalXMLFromTableNil(t *testing.T) {
	t.Parallel()

	data, err := MarshalXMLFromTable(nil)
	if err != nil {
		t.Fatalf("MarshalXMLFromTable() error = %v", err)
	}

	result := string(data)
	assertContains(t, result, "<?xml", "XML should contain XML declaration")
}

func TestMarshalXMLFromTableWithData(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Value"})
	data.AddRow([]string{"test", "123"})

	result, err := MarshalXMLFromTable(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTable() error = %v", err)
	}

	outputStr := string(result)
	assertContains(t, outputStr, "<table>", "XML should contain <table>")
	assertContains(t, outputStr, "<headers>", "XML should contain <headers>")
	assertContains(t, outputStr, "<row>", "XML should contain <row>")
}

func TestMarshalXMLFromTableEmpty(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{})

	result, err := MarshalXMLFromTable(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTable() error = %v", err)
	}

	outputStr := string(result)
	assertContains(t, outputStr, "</table>", "XML should contain closing </table>")
}

func TestMarshalXML(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}

	person := Person{Name: "John", Age: 30}

	data, err := MarshalXML(person)
	if err != nil {
		t.Fatalf("MarshalXML() error = %v", err)
	}

	result := string(data)
	if !strings.Contains(result, "<Person>") && !strings.Contains(result, "<name>John</name>") {
		t.Error("XML should contain person data")
	}
}

func TestMarshalXMLIndent(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string `xml:"name"`
		Age  int    `xml:"age"`
	}

	person := Person{Name: "John", Age: 30}

	data, err := MarshalXMLIndent(person, "", "  ")
	if err != nil {
		t.Fatalf("MarshalXMLIndent() error = %v", err)
	}

	result := string(data)
	assertContains(t, result, "  <name>", "Indented XML should contain indentation")
}

func TestMarshalXMLFromTableNoHeaders(t *testing.T) {
	t.Parallel()

	data := output.NewTable(nil)
	data.AddRow([]string{"a", "b"})

	result, err := MarshalXMLFromTable(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTable() error = %v", err)
	}

	outputStr := string(result)
	assertContains(t, outputStr, "<rows>", "should contain rows")
	assertContains(t, outputStr, "<row>", "should contain row")
	assertContains(t, outputStr, "</table>", "should contain table close")
}

func TestMarshalXMLIndentError(t *testing.T) {
	t.Parallel()

	_, err := MarshalXMLIndent(make(chan int), "", "  ")
	if err == nil {
		t.Fatal("expected error for unmarshalable type")
	}
}

func TestMarshalXMLFromTableWithFooter(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Count"})
	data.AddRow([]string{"Alice", "10"})
	data.Footer = []string{"Total", "10"}

	result, err := MarshalXMLFromTable(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTable() error = %v", err)
	}

	outputStr := string(result)
	assertContains(t, outputStr, "<footer>", "XML should contain <footer>")
	assertContains(t, outputStr, "Total", "XML footer should contain footer text")
	assertContains(t, outputStr, "</footer>", "XML should contain </footer>")
}

func TestMarshalXMLFromTableNoFooter(t *testing.T) {
	t.Parallel()

	data := output.NewTableWithRow([]string{"Name"}, "Alice")

	result, err := MarshalXMLFromTable(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTable() error = %v", err)
	}

	if strings.Contains(string(result), "<footer>") {
		t.Error("XML should not contain <footer> when no footer is set")
	}
}
