package output

import (
	"strings"
	"testing"
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

func TestMarshalXMLFromTableDataNil(t *testing.T) {
	t.Parallel()

	data, err := MarshalXMLFromTableData(nil)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
	}

	result := string(data)
	assertContains(t, result, "<?xml", "XML should contain XML declaration")
}

func TestMarshalXMLFromTableDataWithData(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{"Name", "Value"})
	data.AddRow([]string{"test", "123"})

	result, err := MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
	}

	output := string(result)
	assertContains(t, output, "<table>", "XML should contain <table>")
	assertContains(t, output, "<headers>", "XML should contain <headers>")
	assertContains(t, output, "<row>", "XML should contain <row>")
}

func TestMarshalXMLFromTableDataEmpty(t *testing.T) {
	t.Parallel()

	data := NewTableData([]string{})

	result, err := MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
	}

	output := string(result)
	assertContains(t, output, "</table>", "XML should contain closing </table>")
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
