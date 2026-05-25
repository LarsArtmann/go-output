package serialization

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
		{name: "simple map", input: map[string]int{"a": 1, "b": 2}, wantErr: false},
		{name: "slice", input: []int{1, 2, 3}, wantErr: false},
		{name: "string", input: "hello", wantErr: false},
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
		{name: "map", data: "a: 1\nb: 2", wantErr: false},
		{name: "slice", data: "- 1\n- 2\n- 3", wantErr: false},
		{name: "invalid", data: "invalid: yaml: [", wantErr: true},
	}

	for _, tt := range tests {
		testUnmarshalError(t, tt.name, tt.data, tt.wantErr, UnmarshalYAML, "UnmarshalYAML")
	}
}

func TestYAMLTableRenderer(t *testing.T) {
	t.Parallel()

	t.Run("renders table as YAML sequence", func(t *testing.T) {
		t.Parallel()

		r := NewYAMLTableRenderer()
		r.SetHeaders([]string{"Name", "Age"})
		r.AddRow([]string{"Alice", "30"})
		r.AddRow([]string{"Bob", "25"})

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertOutputContains(t, got, "Name: Alice")
		assertOutputContains(t, got, `Age: "25"`)
	})

	t.Run("nil data returns empty array", func(t *testing.T) {
		t.Parallel()

		r := NewYAMLTableRenderer()

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if got != "[]\n" {
			t.Errorf("Render() = %q, want []\\n", got)
		}
	})

	t.Run("no headers returns empty array", func(t *testing.T) {
		t.Parallel()

		r := NewYAMLTableRenderer()
		r.AddRow([]string{"a"})

		got, err := r.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if got != "[]\n" {
			t.Errorf("Render() = %q, want []\\n", got)
		}
	})
}

func TestMarshalYAMLError(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unmarshalable type")
		}
	}()

	_, _ = MarshalYAML(make(chan int))
}

type benchmarkYAMLStruct struct {
	ID        int      `yaml:"id"`
	Name      string   `yaml:"name"`
	Items     []string `yaml:"items"`
	Count     int      `yaml:"count"`
	Active    bool     `yaml:"active"`
	CreatedAt string   `yaml:"created_at"`
	UpdatedAt string   `yaml:"updated_at"`
}

func BenchmarkMarshalYAML(b *testing.B) {
	data := benchmarkYAMLStruct{
		ID: 12345, Name: "Test Project Alpha",
		Items: []string{"item1", "item2", "item3", "item4", "item5"},
		Count: 100, Active: true,
	}

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
		var result benchmarkYAMLStruct

		_ = UnmarshalYAML(yamlData, &result)
	}
}
