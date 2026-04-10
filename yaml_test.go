package output

import (
	"testing"
)

func TestMarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name:    "simple map",
			input:   map[string]int{"a": 1, "b": 2},
			wantErr: false,
		},
		{
			name:    "slice",
			input:   []int{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "string",
			input:   "hello",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := MarshalYAML(tt.input)
			assertMarshalError(t, "MarshalYAML", err, tt.wantErr)
			if tt.wantErr {
				return
			}

			if string(got) == "" {
				t.Error("MarshalYAML() produced empty output")
			}
		})
	}
}

func TestUnmarshalYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "map",
			data:    "a: 1\nb: 2",
			wantErr: false,
		},
		{
			name:    "slice",
			data:    "- 1\n- 2\n- 3",
			wantErr: false,
		},
		{
			name:    "invalid",
			data:    "invalid: yaml: [",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		testUnmarshalError(t, tt.name, tt.data, tt.wantErr, UnmarshalYAML, "UnmarshalYAML")
	}
}

func BenchmarkMarshalYAML(b *testing.B) {
	data := NewBenchmarkData()

	for b.Loop() {
		_, _ = MarshalYAML(data)
	}
}

func BenchmarkUnmarshalYAML(b *testing.B) {
	yamlData := []byte(`id: 12345
name: Test Project Alpha
items:
  - item1
  - item2
  - item3
  - item4
  - item5
count: 100
active: true
created_at: "2026-03-22T10:00:00Z"
updated_at: "2026-03-22T12:00:00Z"`)

	for b.Loop() {
		var result BenchmarkYAMLStruct

		_ = UnmarshalYAML(yamlData, &result)
	}
}
