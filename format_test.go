package output

import (
	"testing"
)

func TestParseOutputFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    OutputFormat
		wantErr bool
	}{
		{"table", "table", OutputFormatTable, false},
		{"json", "json", OutputFormatJSON, false},
		{"csv", "csv", OutputFormatCSV, false},
		{"markdown", "markdown", OutputFormatMarkdown, false},
		{"d2", "d2", OutputFormatD2, false},
		{"yaml", "yaml", OutputFormatYAML, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOutputFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOutputFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormatString(t *testing.T) {
	tests := []struct {
		format OutputFormat
		want   string
	}{
		{OutputFormatTable, "table"},
		{OutputFormatJSON, "json"},
		{OutputFormatCSV, "csv"},
		{OutputFormatMarkdown, "markdown"},
		{OutputFormatD2, "d2"},
		{OutputFormatYAML, "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.format.String(); got != tt.want {
				t.Errorf("OutputFormat.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputFormatAllowedValues(t *testing.T) {
	got := OutputFormatTable.AllowedValues()
	want := []string{"table", "json", "csv", "markdown", "d2", "yaml"}

	if len(got) != len(want) {
		t.Errorf("AllowedValues() returned %d values, want %d", len(got), len(want))
	}

	for i, v := range got {
		if v != want[i] {
			t.Errorf("AllowedValues()[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestOutputFormatIsValid(t *testing.T) {
	tests := []struct {
		format OutputFormat
		want   bool
	}{
		{OutputFormatTable, true},
		{OutputFormatJSON, true},
		{OutputFormatCSV, true},
		{OutputFormatMarkdown, true},
		{OutputFormatD2, true},
		{OutputFormatYAML, true},
		{OutputFormat("invalid"), false},
		{OutputFormat(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.IsValid(); got != tt.want {
				t.Errorf("OutputFormat(%q).IsValid() = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}

func FuzzParseOutputFormat(f *testing.F) {
	f.Add("table")
	f.Add("json")
	f.Add("csv")
	f.Add("markdown")
	f.Add("d2")
	f.Add("yaml")
	f.Add("invalid")
	f.Add("")

	f.Fuzz(func(t *testing.T, s string) {
		format, err := ParseOutputFormat(s)
		if err != nil {
			if format != "" {
				t.Errorf("ParseOutputFormat(%q) returned error but non-empty format: %q", s, format)
			}
		}
		if format.IsValid() && err == nil {
			if string(format) != s {
				t.Errorf("ParseOutputFormat(%q) = %q, but IsValid() was true", s, format)
			}
		}
	})
}
