package gentest

import (
	"errors"
	"html"
	"strings"
	"testing"
)

var errTest = errors.New("test error")

type mockHTMLRenderer struct {
	headers []string
	rows    [][]string
}

func (m *mockHTMLRenderer) SetHeaders(h []string) { m.headers = h }

func (m *mockHTMLRenderer) AddRow(r []string) { m.rows = append(m.rows, r) }

func (m *mockHTMLRenderer) Render() (string, error) {
	var b strings.Builder

	b.WriteString("<table>")

	for _, h := range m.headers {
		b.WriteString("<th>")
		b.WriteString(html.EscapeString(h))
		b.WriteString("</th>")
	}

	for _, row := range m.rows {
		for _, cell := range row {
			b.WriteString("<td>")
			b.WriteString(html.EscapeString(cell))
			b.WriteString("</td>")
		}
	}

	b.WriteString("</table>")

	return b.String(), nil
}

func TestAssertOutputContains(t *testing.T) {
	t.Parallel()

	t.Run("contains substring", func(t *testing.T) {
		t.Parallel()

		AssertOutputContains(t, "hello world", "hello")
	})

	t.Run("missing substring", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertOutputContains(mock, "hello world", "missing")

		if !mock.Failed() {
			t.Error("expected test to fail for missing substring")
		}
	})
}

func TestAssertValidJSON(t *testing.T) {
	t.Parallel()

	t.Run("valid JSON", func(t *testing.T) {
		t.Parallel()

		AssertValidJSON(t, `{"key": "value"}`)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertValidJSON(mock, "not json")

		if !mock.Failed() {
			t.Error("expected test to fail for invalid JSON")
		}
	})
}

func TestAssertValidYAML(t *testing.T) {
	t.Parallel()

	t.Run("valid YAML", func(t *testing.T) {
		t.Parallel()

		AssertValidYAML(t, "key: value\n")
	})

	t.Run("invalid YAML", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertValidYAML(mock, ":\n  :\n  :")

		if !mock.Failed() {
			t.Error("expected test to fail for invalid YAML")
		}
	})
}

func TestAssertMarshalError(t *testing.T) {
	t.Parallel()

	t.Run("error expected and occurred", func(t *testing.T) {
		t.Parallel()

		AssertMarshalError(t, "test", errTest, true)
	})

	t.Run("no error expected and none occurred", func(t *testing.T) {
		t.Parallel()

		AssertMarshalError(t, "test", nil, false)
	})

	t.Run("error expected but none occurred", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertMarshalError(mock, "test", nil, true)

		if !mock.Failed() {
			t.Error("expected test to fail when error expected but none occurred")
		}
	})

	t.Run("no error expected but error occurred", func(t *testing.T) {
		t.Parallel()

		mock := &testing.T{}

		AssertMarshalError(mock, "test", errTest, false)

		if !mock.Failed() {
			t.Error("expected test to fail when error not expected but occurred")
		}
	})
}

func TestAssertHTMLEscape(t *testing.T) {
	t.Parallel()

	t.Run("properly escapes", func(t *testing.T) {
		t.Parallel()

		AssertHTMLEscape(t, func() HTMLEscapeTestRenderer {
			return &mockHTMLRenderer{}
		}, "mock")
	})
}
