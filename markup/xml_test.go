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

	data := output.NewTableData([]string{"Name", "Value"})
	data.AddRow([]string{"test", "123"})

	result, err := MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
	}

	outputStr := string(result)
	assertContains(t, outputStr, "<table>", "XML should contain <table>")
	assertContains(t, outputStr, "<headers>", "XML should contain <headers>")
	assertContains(t, outputStr, "<row>", "XML should contain <row>")
}

func TestMarshalXMLFromTableDataEmpty(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{})

	result, err := MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
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

func TestXMLWriterWriteHeaderError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&errorWriter{})

	err := x.WriteHeader([]string{"Name"})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "xml header", "error should mention xml header")
}

func TestXMLWriterWriteHeaderTableOpenError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&writeNThenFailWriter{remaining: 1})

	err := x.WriteHeader([]string{"Name"})
	if err == nil {
		t.Fatal("expected error on table open")
	}

	assertContains(t, err.Error(), "table open", "error should mention table open")
}

func TestXMLWriterWriteHeaderHeadersOpenError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&writeNThenFailWriter{remaining: 2})

	err := x.WriteHeader([]string{"Name"})
	if err == nil {
		t.Fatal("expected error on headers open")
	}

	assertContains(t, err.Error(), "headers open", "error should mention headers open")
}

func TestXMLWriterWriteHeaderColumnsError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&writeNThenFailWriter{remaining: 3})

	err := x.WriteHeader([]string{"Name"})
	if err == nil {
		t.Fatal("expected error on columns")
	}

	assertContains(t, err.Error(), "columns", "error should mention columns")
}

func TestXMLWriterWriteHeaderHeadersCloseError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&writeNThenFailWriter{remaining: 6})

	err := x.WriteHeader([]string{"Name"})
	if err == nil {
		t.Fatal("expected error on headers close")
	}

	assertContains(t, err.Error(), "headers close", "error should mention headers close")
}

func TestXMLWriterWriteHeaderRowsOpenError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&writeNThenFailWriter{remaining: 7})

	err := x.WriteHeader([]string{"Name"})
	if err == nil {
		t.Fatal("expected error on rows open")
	}

	assertContains(t, err.Error(), "rows open", "error should mention rows open")
}

func TestXMLWriterWriteFooterTableCloseError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&writeNThenFailWriter{remaining: 1})

	err := x.WriteFooter()
	if err == nil {
		t.Fatal("expected error on table close")
	}

	assertContains(t, err.Error(), "table close", "error should mention table close")
}

func TestXMLWriterWriteRowError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&errorWriter{})

	err := x.WriteRow([]string{"test"})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "row", "error should mention row")
}

func TestXMLWriterWriteRowsError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&errorWriter{})

	err := x.WriteRows([][]string{{"a", "b"}})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "row", "error should mention row")
}

func TestXMLWriterWriteFooterError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&errorWriter{})

	err := x.WriteFooter()
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "rows close", "error should mention rows close")
}

func TestMarshalXMLFromTableDataNoHeaders(t *testing.T) {
	t.Parallel()

	data := output.NewTableData(nil)
	data.AddRow([]string{"a", "b"})

	result, err := MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
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
