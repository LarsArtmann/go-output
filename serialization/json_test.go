package serialization

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{
			name:    "simple map",
			input:   map[string]int{"a": 1, "b": 2},
			want:    `{"a":1,"b":2}`,
			wantErr: false,
		},
		{
			name:    "slice",
			input:   []int{1, 2, 3},
			want:    "[1,2,3]",
			wantErr: false,
		},
		{
			name:    "string",
			input:   "hello",
			want:    `"hello"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := MarshalJSON(tt.input)
			assertMarshalError(t, "MarshalJSON", err, tt.wantErr)

			if tt.wantErr {
				return
			}

			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{name: "map", data: `{"a":1}`, wantErr: false},
		{name: "slice", data: "[1,2,3]", wantErr: false},
		{name: "invalid", data: `{`, wantErr: true},
	}

	for _, tt := range tests {
		testUnmarshalError(t, tt.name, tt.data, tt.wantErr, UnmarshalJSON, "UnmarshalJSON")
	}
}

func TestJSONWriter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	w := NewJSONWriter(&buf)

	data := map[string]int{"test": 42}

	err := w.Encode(data)
	if err != nil {
		t.Errorf("JSONWriter.Encode() error = %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("JSONWriter.Encode() produced empty output")
	}
}

func TestNewJSONWriter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	w := NewJSONWriter(&buf)

	if w.Writer != &buf {
		t.Error("NewJSONWriter() did not set Writer correctly")
	}
}

func TestJSONTableRenderer(t *testing.T) {
	t.Parallel()

	t.Run("renders table as JSON array of objects", func(t *testing.T) {
		t.Parallel()

		r := NewJSONTableRenderer()
		r.SetHeaders([]string{"Name", "Age"})
		r.AddRow([]string{"Alice", "30"})
		r.AddRow([]string{"Bob", "25"})

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertOutputContains(t, got, `"Name": "Alice"`)
		assertOutputContains(t, got, `"Age": "25"`)
	})

	t.Run("nil data returns empty array", func(t *testing.T) {
		t.Parallel()

		r := NewJSONTableRenderer()

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if got != "[]" {
			t.Errorf("Render() = %q, want []", got)
		}
	})

	t.Run("short row omits missing cells", func(t *testing.T) {
		t.Parallel()

		r := NewJSONTableRenderer()
		r.SetHeaders([]string{"A", "B", "C"})
		r.AddRow([]string{"1"})

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertOutputContains(t, got, `"A": "1"`)

		if strings.Contains(got, `"B"`) {
			t.Errorf("Render() = %q, want B absent for short row", got)
		}
	})
}

func TestJSONTableRendererNoHeaders(t *testing.T) {
	t.Parallel()

	r := NewJSONTableRenderer()
	r.AddRow([]string{"a"})

	got, err := r.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got != "[]" {
		t.Errorf("Render() = %q, want []", got)
	}
}

func TestMarshalJSONError(t *testing.T) {
	t.Parallel()

	_, err := MarshalJSON(make(chan int))
	if err == nil {
		t.Fatal("expected error for unmarshalable type")
	}
}

func TestJSONWriterEncodeError(t *testing.T) {
	t.Parallel()

	w := NewJSONWriter(&errorWriter{})

	err := w.Encode(map[string]string{"k": "v"})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "encode json", "error should mention encode json")
}

type benchmarkStruct struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Items     []string  `json:"items"`
	Count     int       `json:"count"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func BenchmarkMarshalJSON(b *testing.B) {
	data := benchmarkStruct{
		ID: 12345, Name: "Test Project Alpha",
		Items: []string{"item1", "item2", "item3", "item4", "item5"},
		Count: 100, Active: true,
	}

	for b.Loop() {
		_, _ = MarshalJSON(data)
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	jsonData := []byte(
		`{"id":12345,"name":"Test Project Alpha","items":["item1","item2","item3","item4","item5"],"count":100,"active":true,"created_at":"2026-03-22T10:00:00Z","updated_at":"2026-03-22T12:00:00Z"}`,
	)

	for b.Loop() {
		var result benchmarkStruct

		_ = UnmarshalJSON(jsonData, &result)
	}
}
