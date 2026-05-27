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

func TestXMLWriterWriteHeaderPartialErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		remaining   int
		wantErrPart string
	}{
		{"table open", 1, "table open"},
		{"headers open", 2, "headers open"},
		{"columns", 3, "columns"},
		{"headers close", 6, "headers close"},
		{"rows open", 7, "rows open"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			x := NewXMLWriter(&writeNThenFailWriter{Remaining: tt.remaining})

			err := x.WriteHeader([]string{"Name"})
			if err == nil {
				t.Fatal("expected error")
			}

			assertContains(t, err.Error(), tt.wantErrPart, "error should mention "+tt.wantErrPart)
		})
	}
}

func TestXMLWriterWriteFooterTableCloseError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&writeNThenFailWriter{Remaining: 1})

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

func TestMarshalXMLFromTableDataWithFooter(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name", "Count"})
	data.AddRow([]string{"Alice", "10"})
	data.Footer = []string{"Total", "10"}

	result, err := MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
	}

	outputStr := string(result)
	assertContains(t, outputStr, "<footer>", "XML should contain <footer>")
	assertContains(t, outputStr, "Total", "XML footer should contain footer text")
	assertContains(t, outputStr, "</footer>", "XML should contain </footer>")
}

func TestMarshalXMLFromTableDataNoFooter(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Alice"})

	result, err := MarshalXMLFromTableData(data)
	if err != nil {
		t.Fatalf("MarshalXMLFromTableData() error = %v", err)
	}

	if strings.Contains(string(result), "<footer>") {
		t.Error("XML should not contain <footer> when no footer is set")
	}
}
