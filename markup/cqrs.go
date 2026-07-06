package markup

import (
	"fmt"
	"io"
	"strings"

	"github.com/larsartmann/go-output"
)

// WriteXML writes a Table as XML to the provided writer.
func WriteXML(w io.Writer, data *output.Table) error {
	b, err := MarshalXMLFromTable(data)
	if err != nil {
		return fmt.Errorf("marshal xml: %w", err)
	}

	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write xml output: %w", err)
	}

	return nil
}

// RenderXML renders a Table as an XML string.
func RenderXML(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteXML(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteAsciiDoc writes a Table as AsciiDoc to the provided writer.
func WriteAsciiDoc(w io.Writer, data *output.Table) error {
	b, err := MarshalAsciiDocFromTable(data)
	if err != nil {
		return fmt.Errorf("marshal asciidoc: %w", err)
	}

	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write asciidoc output: %w", err)
	}

	return nil
}

// RenderAsciiDoc renders a Table as an AsciiDoc string.
func RenderAsciiDoc(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteAsciiDoc(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteHTML writes a Table as HTML to the provided writer.
func WriteHTML(w io.Writer, data *output.Table) error {
	return renderHTMLTable(w, data, output.RenderOptions{})
}

// RenderHTML renders a Table as an HTML string.
func RenderHTML(data *output.Table) (string, error) {
	var buf strings.Builder
	if err := WriteHTML(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
