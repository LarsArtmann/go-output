package cmdguard

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func testParseHelper[T EnumValue](
	t *testing.T,
	flagName string,
	newFlag func(*T) *EnumFlag[T],
	tests []struct {
		name    string
		input   string
		want    T
		wantErr bool
	},
) {
	t.Helper()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var val T

			flag := newFlag(&val)

			err := flag.Parse(testCase.input)
			if (err != nil) != testCase.wantErr {
				t.Errorf("%s.Parse() error = %v, wantErr %v", flagName, err, testCase.wantErr)

				return
			}

			if !testCase.wantErr && val != testCase.want {
				t.Errorf("%s.Parse() = %v, want %v", flagName, val, testCase.want)
			}
		})
	}
}

func testNewFlagHelper[T EnumValue](
	t *testing.T,
	newFlag func(*T) *EnumFlag[T],
	val *T,
) {
	t.Helper()

	flag := newFlag(val)
	if flag == nil {
		t.Fatal("NewFlag() returned nil")
	}

	if flag.value != val {
		t.Error("NewFlag() did not set value correctly")
	}
}

func TestNewFlag(t *testing.T) {
	t.Parallel()

	t.Run("ColorModeFlag", func(t *testing.T) {
		t.Parallel()

		var val output.ColorMode
		testNewFlagHelper(t, NewColorModeFlag, &val)
	})

	t.Run("OutputFormatFlag", func(t *testing.T) {
		t.Parallel()

		var val output.Format
		testNewFlagHelper(t, NewOutputFormatFlag, &val)
	})

	t.Run("SortByFlag", func(t *testing.T) {
		t.Parallel()

		var val output.SortBy
		testNewFlagHelper(t, NewSortByFlag, &val)
	})
}

func TestColorModeFlag_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    output.ColorMode
		wantErr bool
	}{
		{"auto", "auto", output.ColorModeAuto, false},
		{"always", "always", output.ColorModeAlways, false},
		{"never", "never", output.ColorModeNever, false},
		{"invalid", "invalid", output.ColorModeAuto, true},
		{"empty", "", output.ColorModeAuto, true},
	}

	testParseHelper(t, "ColorModeFlag", NewColorModeFlag, tests)
}

func TestColorModeFlag_AllowedValues(t *testing.T) {
	t.Parallel()

	val := output.ColorModeAuto
	flag := NewColorModeFlag(&val)
	got := flag.AllowedValues()
	want := []string{"auto", "always", "never"}

	if len(got) != len(want) {
		t.Errorf("AllowedValues() returned %d values, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("AllowedValues()[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func testDefaultHelper[T EnumValue](
	t *testing.T,
	newFlag func(*T) *EnumFlag[T],
	defaultVal T,
) {
	t.Helper()

	flag := newFlag(&defaultVal)
	if got := flag.Default(); got != defaultVal.String() {
		t.Errorf("Default() = %v, want %v", got, defaultVal.String())
	}
}

func TestColorModeFlag_Default(t *testing.T) {
	t.Parallel()
	testDefaultHelper(t, NewColorModeFlag, output.ColorModeAuto)
}

func TestOutputFormatFlag_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    output.Format
		wantErr bool
	}{
		{"table", "table", output.FormatTable, false},
		{"json", "json", output.FormatJSON, false},
		{"csv", "csv", output.FormatCSV, false},
		{"tsv", "tsv", output.FormatTSV, false},
		{"xml", "xml", output.FormatXML, false},
		{"markdown", "markdown", output.FormatMarkdown, false},
		{"d2", "d2", output.FormatD2, false},
		{"yaml", "yaml", output.FormatYAML, false},
		{"html", "html", output.FormatHTML, false},
		{"tree", "tree", output.FormatTree, false},
		{"mermaid", "mermaid", output.FormatMermaid, false},
		{"dot", "dot", output.FormatDOT, false},
		{"invalid", "invalid", output.FormatTable, true},
		{"empty", "", output.FormatTable, true},
	}

	testParseHelper(t, "OutputFormatFlag", NewOutputFormatFlag, tests)
}

func TestOutputFormatFlag_AllowedValues(t *testing.T) {
	t.Parallel()

	val := output.FormatTable
	flag := NewOutputFormatFlag(&val)
	got := flag.AllowedValues()
	want := []string{
		"table",
		"json",
		"csv",
		"tsv",
		"markdown",
		"xml",
		"d2",
		"yaml",
		"html",
		"tree",
		"mermaid",
		"dot",
	}

	if len(got) != len(want) {
		t.Errorf("AllowedValues() returned %d values, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("AllowedValues()[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestOutputFormatFlag_Default(t *testing.T) {
	t.Parallel()
	testDefaultHelper(t, NewOutputFormatFlag, output.OutputFormatJSON)
}

func TestSortByFlag_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    output.SortBy
		wantErr bool
	}{
		{"name", "name", output.SortByName, false},
		{"importance", "importance", output.SortByImportance, false},
		{"created_at", "created_at", output.SortByCreatedAt, false},
		{"updated_at", "updated_at", output.SortByUpdatedAt, false},
		{"health", "health", output.SortByHealth, false},
		{"complexity", "complexity", output.SortByComplexity, false},
		{"invalid", "invalid", output.SortByName, true},
		{"empty", "", output.SortByName, true},
	}

	testParseHelper(t, "SortByFlag", NewSortByFlag, tests)
}

func TestSortByFlag_AllowedValues(t *testing.T) {
	t.Parallel()

	val := output.SortByName
	flag := NewSortByFlag(&val)
	got := flag.AllowedValues()
	want := []string{"name", "importance", "created_at", "updated_at", "health", "complexity"}

	if len(got) != len(want) {
		t.Errorf("AllowedValues() returned %d values, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("AllowedValues()[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestSortByFlag_Default(t *testing.T) {
	t.Parallel()
	testDefaultHelper(t, NewSortByFlag, output.SortByName)
}

func TestNewColorModeFlag(t *testing.T) {
	t.Parallel()

	val := output.ColorModeAuto

	flag := NewColorModeFlag(&val)
	if flag == nil {
		t.Fatal("NewColorModeFlag() returned nil")
	}

	if flag.value != &val {
		t.Error("NewColorModeFlag() did not set value correctly")
	}
}

func TestNewOutputFormatFlag(t *testing.T) {
	t.Parallel()

	val := output.OutputFormatTable

	flag := NewOutputFormatFlag(&val)
	if flag == nil {
		t.Fatal("NewOutputFormatFlag() returned nil")
	}

	if flag.value != &val {
		t.Error("NewOutputFormatFlag() did not set value correctly")
	}
}

func TestNewSortByFlag(t *testing.T) {
	t.Parallel()

	val := output.SortByName

	flag := NewSortByFlag(&val)
	if flag == nil {
		t.Fatal("NewSortByFlag() returned nil")
	}

	if flag.value != &val {
		t.Error("NewSortByFlag() did not set value correctly")
	}
}
