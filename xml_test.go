package output

import (
	"strings"
	"testing"
)

func TestXMLWriterWriteHeader(t *testing.T) {
	t.Parallel()
	x := NewXMLWriter()
	err := x.WriteHeader([]string{"Name", "Value"})
	if err != nil {
		t.Fatalf("WriteHeader() error = %v", err)
	}
	result := x.String()
	if !strings.Contains(result, "<headers>") {
		t.Error("XML should contain <headers>")
	}
	if !strings.Contains(result, "<column>Name</column>") {
		t.Error("XML should contain <column>Name</column>")
	}
}

func TestXMLWriterWriteRow(t *testing.T) {
	t.Parallel()
	x := NewXMLWriter()
	_ = x.WriteHeader([]string{"Name", "Value"})
	err := x.WriteRow([]string{"test", "123"})
	if err != nil {
		t.Fatalf("WriteRow() error = %v", err)
	}
	result := x.String()
	if !strings.Contains(result, "<row>") {
		t.Error("XML should contain <row>")
	}
	if !strings.Contains(result, "<cell>test</cell>") {
		t.Error("XML should contain <cell>test</cell>")
	}
}

func TestXMLWriterWriteRows(t *testing.T) {
	t.Parallel()
	x := NewXMLWriter()
	_ = x.WriteHeader([]string{"Name", "Value"})
	err := x.WriteRows([][]string{
		{"a", "1"},
		{"b", "2"},
	})
	if err != nil {
		t.Fatalf("WriteRows() error = %v", err)
	}
	result := x.String()
	if strings.Count(result, "<row>") != 2 {
		t.Errorf("XML should contain 2 <row> elements")
	}
}

func TestXMLWriterEscape(t *testing.T) {
	t.Parallel()
	x := NewXMLWriter()
	_ = x.WriteHeader([]string{"Name"})
	_ = x.WriteRow([]string{"<script>alert('xss')</script>"})
	result := x.String()
	if strings.Contains(result, "<script>") {
		t.Error("XML should escape <script> tags")
	}
	if !strings.Contains(result, "&lt;script&gt;") {
		t.Error("XML should contain escaped &lt;script&gt;")
	}
}

func TestMarshalXMLFromTableDataNil(t *testing.T) {
	t.Parallel()
	data, err := MarshalXMLFromTableData(nil)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
	}
	result := string(data)
	if !strings.Contains(result, "<?xml") {
		t.Error("XML should contain XML declaration")
	}
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
	if !strings.Contains(output, "<table>") {
		t.Error("XML should contain <table>")
	}
	if !strings.Contains(output, "<headers>") {
		t.Error("XML should contain <headers>")
	}
	if !strings.Contains(output, "<row>") {
		t.Error("XML should contain <row>")
	}
}

func TestMarshalXMLFromTableDataEmpty(t *testing.T) {
	t.Parallel()
	data := NewTableData([]string{})
	result, err := MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
	}
	output := string(result)
	if !strings.Contains(output, "</table>") {
		t.Error("XML should contain closing </table>")
	}
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
	if !strings.Contains(result, "  <name>") {
		t.Error("Indented XML should contain indentation")
	}
}
