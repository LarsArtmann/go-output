package output

import (
	"bytes"
	"testing"
	"time"
)

func testUnmarshalError(
	t *testing.T,
	name, data string,
	wantErr bool,
	unmarshal func([]byte, any) error,
	funcName string,
) {
	t.Run(name, func(t *testing.T) {
		t.Parallel()

		var got any

		err := unmarshal([]byte(data), &got)
		if (err != nil) != wantErr {
			t.Errorf("%s() error = %v, wantErr %v", funcName, err, wantErr)
		}
	})
}

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
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestMarshalJSONIndent(t *testing.T) {
	t.Parallel()

	input := map[string]int{"a": 1}

	got, err := MarshalJSONIndent(input, "", "  ")
	if err != nil {
		t.Errorf("MarshalJSONIndent() error = %v", err)

		return
	}

	want := "{\n  \"a\": 1\n}"
	if string(got) != want {
		t.Errorf("MarshalJSONIndent() = %v, want %v", string(got), want)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		want    any
		wantErr bool
	}{
		{
			name:    "map",
			data:    `{"a":1}`,
			want:    &map[string]any{},
			wantErr: false,
		},
		{
			name:    "slice",
			data:    "[1,2,3]",
			want:    &[]any{},
			wantErr: false,
		},
		{
			name:    "invalid",
			data:    `{`,
			want:    nil,
			wantErr: true,
		},
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

func BenchmarkMarshalJSON(b *testing.B) {
	data := NewBenchmarkData()

	for b.Loop() {
		_, _ = MarshalJSON(data)
	}
}

func BenchmarkMarshalJSONIndent(b *testing.B) {
	data := NewBenchmarkData()

	for b.Loop() {
		_, _ = MarshalJSONIndent(data, "", "  ")
	}
}

type BenchmarkStruct struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Items     []string  `json:"items"`
	Count     int       `json:"count"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	jsonData := []byte(
		`{"id":12345,"name":"Test Project Alpha","items":["item1","item2","item3","item4","item5"],"count":100,"active":true,"created_at":"2026-03-22T10:00:00Z","updated_at":"2026-03-22T12:00:00Z"}`,
	)

	for b.Loop() {
		var result BenchmarkStruct

		_ = UnmarshalJSON(jsonData, &result)
	}
}
