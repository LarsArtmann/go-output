package serialization

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func FuzzMarshalJSONFromTable(f *testing.F) {
	f.Add("Name", "Age", "Alice", "30")
	f.Add("", "", "", "")
	f.Add("a,b", "c,d", "e,f", "g,h")

	f.Fuzz(func(_ *testing.T, h1, h2, v1, v2 string) {
		data := output.NewTable([]string{h1, h2})
		data.AddRow([]string{v1, v2})

		renderer := NewJSONTableRenderer()
		renderer.SetData(data)
		_, _ = renderer.Render()
	})
}

func FuzzMarshalYAMLFromTable(f *testing.F) {
	f.Add("Name", "Age", "Alice", "30")
	f.Add("", "", "", "")

	f.Fuzz(func(_ *testing.T, h1, h2, v1, v2 string) {
		data := output.NewTable([]string{h1, h2})
		data.AddRow([]string{v1, v2})

		renderer := NewYAMLTableRenderer()
		renderer.SetData(data)
		_, _ = renderer.Render()
	})
}

func FuzzMarshalTOMLFromTable(f *testing.F) {
	f.Add("Name", "Age", "Alice", "30")
	f.Add("", "", "", "")

	f.Fuzz(func(_ *testing.T, h1, h2, v1, v2 string) {
		data := output.NewTable([]string{h1, h2})
		data.AddRow([]string{v1, v2})

		_, _ = MarshalTOMLFromTable(data)
	})
}

func FuzzMarshalJSONLFromTable(f *testing.F) {
	f.Add("Name", "Age", "Alice", "30")
	f.Add("", "", "", "")

	f.Fuzz(func(_ *testing.T, h1, h2, v1, v2 string) {
		data := output.NewTable([]string{h1, h2})
		data.AddRow([]string{v1, v2})

		_, _ = MarshalJSONLFromTable(data)
	})
}
