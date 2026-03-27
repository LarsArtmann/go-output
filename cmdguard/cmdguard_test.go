package cmdguard

import (
	"testing"

	"github.com/larsartmann/go-output"
)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var val output.ColorMode
			flag := NewColorModeFlag(&val)
			err := flag.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ColorModeFlag.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && val != tt.want {
				t.Errorf("ColorModeFlag.Parse() = %v, want %v", val, tt.want)
			}
		})
	}
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

func TestColorModeFlag_Default(t *testing.T) {
	t.Parallel()
	val := output.ColorModeAuto
	flag := NewColorModeFlag(&val)
	if got := flag.Default(); got != "auto" {
		t.Errorf("ColorModeFlag.Default() = %v, want auto", got)
	}
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var val output.Format
			flag := NewOutputFormatFlag(&val)
			err := flag.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("OutputFormatFlag.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && val != tt.want {
				t.Errorf("OutputFormatFlag.Parse() = %v, want %v", val, tt.want)
			}
		})
	}
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
		"markdown",
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
	val := output.OutputFormatJSON
	flag := NewOutputFormatFlag(&val)
	if got := flag.Default(); got != "json" {
		t.Errorf("OutputFormatFlag.Default() = %v, want json", got)
	}
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var val output.SortBy
			flag := NewSortByFlag(&val)
			err := flag.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SortByFlag.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && val != tt.want {
				t.Errorf("SortByFlag.Parse() = %v, want %v", val, tt.want)
			}
		})
	}
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
	val := output.SortByName
	flag := NewSortByFlag(&val)
	if got := flag.Default(); got != "name" {
		t.Errorf("SortByFlag.Default() = %v, want name", got)
	}
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
