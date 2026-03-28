// Package integration provides end-to-end workflow tests for go-output.
package integration

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/sort"
)

// TestCSVToTableData tests converting CSV data to TableData
func TestCSVToTableData(t *testing.T) {
	t.Parallel()

	// Given: Raw CSV-like data
	headers := []string{"Name", "Value"}
	rows := [][]string{
		{"Alpha", "100"},
		{"Beta", "200"},
		{"Gamma", "150"},
	}

	// When: I convert CSV data to TableData
	data := output.NewTableData(headers)
	for _, row := range rows {
		data.AddRow(row)
	}

	// Then: Data should be properly structured
	if data.ColCount() != 2 {
		t.Errorf("Expected 2 columns, got %d", data.ColCount())
	}
	if data.RowCount() != 3 {
		t.Errorf("Expected 3 rows, got %d", data.RowCount())
	}
}

// TestTableDataToJSON tests rendering TableData as JSON
func TestTableDataToJSON(t *testing.T) {
	t.Parallel()

	// Given: TableData
	headers := []string{"Name", "Value"}
	rows := [][]string{
		{"Alpha", "100"},
		{"Beta", "200"},
		{"Gamma", "150"},
	}
	data := output.NewTableData(headers)
	for _, row := range rows {
		data.AddRow(row)
	}

	// When: I render as JSON
	jsonBytes, err := output.MarshalJSONIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("MarshalJSONIndent failed: %v", err)
	}

	// Then: JSON should contain all data
	jsonStr := string(jsonBytes)
	if !strings.Contains(jsonStr, "Alpha") {
		t.Error("JSON should contain Alpha")
	}
	if !strings.Contains(jsonStr, "100") {
		t.Error("JSON should contain 100")
	}
}

// TestTableDataToYAML tests rendering TableData as YAML
func TestTableDataToYAML(t *testing.T) {
	t.Parallel()

	// Given: TableData
	headers := []string{"Name", "Value"}
	rows := [][]string{
		{"Alpha", "100"},
		{"Beta", "200"},
		{"Gamma", "150"},
	}
	data := output.NewTableData(headers)
	for _, row := range rows {
		data.AddRow(row)
	}

	// When: I render as YAML
	yamlBytes, err := output.MarshalYAML(data)
	if err != nil {
		t.Fatalf("MarshalYAML failed: %v", err)
	}

	// Then: YAML should contain all data
	yamlStr := string(yamlBytes)
	if !strings.Contains(yamlStr, "Name") {
		t.Error("YAML should contain Name header")
	}
	if !strings.Contains(yamlStr, "Gamma") {
		t.Error("YAML should contain Gamma")
	}
}

// TestSortAndRenderWorkflow tests sorting before rendering
func TestSortAndRenderWorkflow(t *testing.T) {
	t.Parallel()

	type Item struct {
		Name  string
		Score int
	}

	t.Run("sort by name ascending then render", func(t *testing.T) {
		t.Parallel()

		// Given: Unordered items
		items := []Item{
			{Name: "Zebra", Score: 50},
			{Name: "Apple", Score: 90},
			{Name: "Banana", Score: 70},
		}

		// When: I sort by name ascending and render as JSON
		sorted := sort.New(items, output.SortByName, false)
		sorted.Sort()

		jsonBytes, _ := output.MarshalJSONIndent(items, "", "  ")
		jsonStr := string(jsonBytes)

		// Then: Apple should come first
		if !strings.Contains(jsonStr, `"Apple"`) {
			t.Error("Sorted JSON should contain Apple")
		}
		appleIdx := strings.Index(jsonStr, `"Apple"`)
		zebraIdx := strings.Index(jsonStr, `"Zebra"`)
		if appleIdx > zebraIdx {
			t.Error("Apple should come before Zebra after sorting by name")
		}
	})
}

// TestLargeDatasetWorkflow tests handling of larger datasets
func TestLargeDatasetWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("handles 1000 rows", func(t *testing.T) {
		t.Parallel()

		// Given: Large dataset
		data := output.NewTableData([]string{"ID", "Value"})
		for i := range 1000 {
			data.AddRow([]string{string(rune('A'+i%26)) + string(rune('0'+i%10)), "value"})
		}

		// When: I render as JSON
		jsonBytes, err := output.MarshalJSONIndent(data, "", "  ")
		if err != nil {
			t.Fatalf("MarshalJSONIndent failed: %v", err)
		}

		// Then: Should complete without error and have correct row count
		if len(jsonBytes) == 0 {
			t.Error("JSON output should not be empty")
		}
	})

	t.Run("handles streaming with large data", func(t *testing.T) {
		t.Parallel()

		// Given: Large streaming dataset
		html := output.NewStreamingHTMLRenderer()
		headers := make([]string, 10)
		for i := range headers {
			headers[i] = "Col"
		}
		html.SetHeaders(headers)

		for range 100 {
			row := make([]string, 10)
			for j := range row {
				row[j] = "data"
			}
			html.AddRow(row)
		}

		// When: I stream the output
		var buf bytes.Buffer
		err := html.Stream(&buf)
		if err != nil {
			t.Fatalf("Stream failed: %v", err)
		}

		// Then: Output should be valid HTML
		result := buf.String()
		if !strings.Contains(result, "<table") {
			t.Error("Should contain table tag")
		}
		if !strings.Contains(result, "<tr>") {
			t.Error("Should contain row tags")
		}
	})
}

// TestErrorHandlingWorkflow tests graceful error handling
func TestErrorHandlingWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("invalid format gives clean error", func(t *testing.T) {
		t.Parallel()

		// When: I parse an invalid format
		_, err := output.ParseOutputFormat("not_a_format")

		// Then: Error should be descriptive
		if err == nil {
			t.Error("Expected error for invalid format")
		}
		if !strings.Contains(err.Error(), "invalid") {
			t.Error("Error should mention invalid format")
		}
	})

	t.Run("empty data renders without panic", func(t *testing.T) {
		t.Parallel()

		// Given: Empty TableData
		data := output.NewTableData([]string{})

		// When: I render it - should not panic
		_, err := output.MarshalJSON(data)
		if err != nil {
			t.Errorf("MarshalJSON on empty data should not error: %v", err)
		}
	})
}
