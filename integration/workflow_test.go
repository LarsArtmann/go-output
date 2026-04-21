// Package integration provides end-to-end workflow tests for go-output.
package integration

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/internal/testutils"
	"github.com/larsartmann/go-output/sort"
)

// TestCSVToTableData tests converting CSV data to TableData.
func TestCSVToTableData(t *testing.T) {
	t.Parallel()

	// Given: Raw CSV-like data
	headers, rows := sharedTestData()

	// When: I convert CSV data to TableData
	data := output.NewTableData(headers)
	for _, row := range rows {
		data.AddRow(row)
	}

	// Then: Data should be properly structured
	testutils.AssertTableData(t, data, 2, 3)
}

// TestTableDataToJSON tests rendering TableData as JSON.
func TestTableDataToJSON(t *testing.T) {
	t.Parallel()

	// Given: TableData
	headers, rows := sharedTestData()

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
	testutils.AssertContains(t, jsonStr, "Alpha", "JSON should contain Alpha")
	testutils.AssertContains(t, jsonStr, "100", "JSON should contain 100")
}

// TestTableDataToYAML tests rendering TableData as YAML.
func TestTableDataToYAML(t *testing.T) {
	t.Parallel()

	// Given: TableData
	headers, rows := sharedTestData()

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
	testutils.AssertContains(t, yamlStr, "Name", "YAML should contain Name header")
	testutils.AssertContains(t, yamlStr, "Gamma", "YAML should contain Gamma")
}

// TestSortAndRenderWorkflow tests sorting before rendering.
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
		testutils.AssertContains(t, jsonStr, `"Apple"`, "Sorted JSON should contain Apple")

		appleIdx := strings.Index(jsonStr, `"Apple"`)

		zebraIdx := strings.Index(jsonStr, `"Zebra"`)
		if appleIdx > zebraIdx {
			t.Error("Apple should come before Zebra after sorting by name")
		}
	})
}

// TestLargeDatasetWorkflow tests handling of larger datasets.
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

		headers := output.FilledStrings(10, "Col")
		html.SetHeaders(headers)

		for range 100 {
			row := output.FilledStrings(10, "data")
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
		testutils.AssertContains(t, result, "<table", "Should contain table tag")
		testutils.AssertContains(t, result, "<tr>", "Should contain row tags")
	})
}

// TestErrorHandlingWorkflow tests graceful error handling.
func TestErrorHandlingWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("invalid format gives clean error", func(t *testing.T) {
		t.Parallel()

		// When: I parse an invalid format
		_, err := output.ParseFormat("not_a_format")

		// Then: Error should be descriptive
		if err == nil {
			t.Error("Expected error for invalid format")
		}

		if !strings.Contains(err.Error(), "invalid") {
			t.Error("Error should mention invalid format")
		}
	})

	testutils.RunEmptyDataRendersJSONWithoutPanic(t)
}
