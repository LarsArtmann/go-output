package output

import (
	"testing"
)

func TestMarshalYAML(t *testing.T) {
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
			got, err := MarshalYAML(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) == "" {
				t.Error("MarshalYAML() produced empty output")
			}
		})
	}
}

func TestUnmarshalYAML(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			var got any
			err := UnmarshalYAML([]byte(tt.data), &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type yamlBenchmarkStruct struct {
	ID        int      `yaml:"id"`
	Name      string   `yaml:"name"`
	Items     []string `yaml:"items"`
	Count     int      `yaml:"count"`
	Active    bool     `yaml:"active"`
	CreatedAt string   `yaml:"created_at"`
	UpdatedAt string   `yaml:"updated_at"`
}

func BenchmarkMarshalYAML(b *testing.B) {
	data := yamlBenchmarkStruct{
		ID:        12345,
		Name:      "Test Project Alpha",
		Items:     []string{"item1", "item2", "item3", "item4", "item5"},
		Count:     100,
		Active:    true,
		CreatedAt: "2026-03-22T10:00:00Z",
		UpdatedAt: "2026-03-22T12:00:00Z",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result yamlBenchmarkStruct
		_ = UnmarshalYAML(yamlData, &result)
	}
}
