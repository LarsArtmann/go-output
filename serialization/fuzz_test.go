package serialization

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func FuzzMarshalJSONFromTableData(f *testing.F) {
	f.Add("Name", "Age", "Alice", "30")
	f.Add("", "", "", "")
	f.Add("a,b", "c,d", "e,f", "g,h")

	f.Fuzz(func(_ *testing.T, h1, h2, v1, v2 string) {
		data := output.NewTableData([]string{h1, h2})
		data.AddRow([]string{v1, v2})

		renderer := NewJSONTableRenderer()
		renderer.SetData(data)
		_, _ = renderer.Render()
	})
}

func FuzzMarshalYAMLFromTableData(f *testing.F) {
	f.Add("Name", "Age", "Alice", "30")
	f.Add("", "", "", "")

	f.Fuzz(func(_ *testing.T, h1, h2, v1, v2 string) {
		data := output.NewTableData([]string{h1, h2})
		data.AddRow([]string{v1, v2})

		renderer := NewYAMLTableRenderer()
		renderer.SetData(data)
		_, _ = renderer.Render()
	})
}

func FuzzMarshalTOMLFromTableData(f *testing.F) {
	f.Add("Name", "Age", "Alice", "30")
	f.Add("", "", "", "")

	f.Fuzz(func(_ *testing.T, h1, h2, v1, v2 string) {
		data := output.NewTableData([]string{h1, h2})
		data.AddRow([]string{v1, v2})

		_, _ = MarshalTOMLFromTableData(data)
	})
}

func FuzzMarshalJSONLFromTableData(f *testing.F) {
	f.Add("Name", "Age", "Alice", "30")
	f.Add("", "", "", "")

	f.Fuzz(func(_ *testing.T, h1, h2, v1, v2 string) {
		data := output.NewTableData([]string{h1, h2})
		data.AddRow([]string{v1, v2})

		_, _ = MarshalJSONLFromTableData(data)
	})
}
