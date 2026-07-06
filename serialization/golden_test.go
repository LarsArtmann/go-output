package serialization

import (
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/larsartmann/go-output"
)

func TestGolden_JSON_Table(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Status", "Duration"})
	data.AddRow([]string{"Build", "completed", "1.2s"})
	data.AddRow([]string{"Test", "running", "0.5s"})

	r := NewJSONTableRenderer()
	r.SetData(data)

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_YAML_Table(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Service", "Port", "Health"})
	data.AddRow([]string{"api", "8080", "healthy"})
	data.AddRow([]string{"worker", "9090", "degraded"})

	r := NewYAMLTableRenderer()
	r.SetData(data)

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}

func TestGolden_TOML_Table(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Key", "Value"})
	data.AddRow([]string{"name", "go-output"})
	data.AddRow([]string{"version", "1.0.0"})

	r := NewTOMLTableRenderer()
	r.SetData(data)

	got, err := r.Render()
	if err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, []byte(got))
}
