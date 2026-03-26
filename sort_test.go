package output

import (
	"testing"
)

func TestParseSortBy(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    SortBy
		wantErr bool
	}{
		{"name", "name", SortByName, false},
		{"importance", "importance", SortByImportance, false},
		{"created_at", "created_at", SortByCreatedAt, false},
		{"updated_at", "updated_at", SortByUpdatedAt, false},
		{"health", "health", SortByHealth, false},
		{"complexity", "complexity", SortByComplexity, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseSortBy(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSortBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseSortBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortByString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		sortBy SortBy
		want   string
	}{
		{SortByName, "name"},
		{SortByImportance, "importance"},
		{SortByCreatedAt, "created_at"},
		{SortByUpdatedAt, "updated_at"},
		{SortByHealth, "health"},
		{SortByComplexity, "complexity"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.sortBy.String(); got != tt.want {
				t.Errorf("SortBy.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortByAllowedValues(t *testing.T) {
	t.Parallel()
	got := SortByName.AllowedValues()
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

func TestSortByIsValid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		sortBy SortBy
		want   bool
	}{
		{SortByName, true},
		{SortByImportance, true},
		{SortByCreatedAt, true},
		{SortByUpdatedAt, true},
		{SortByHealth, true},
		{SortByComplexity, true},
		{SortBy("invalid"), false},
		{SortBy(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.sortBy), func(t *testing.T) {
			t.Parallel()
			if got := tt.sortBy.IsValid(); got != tt.want {
				t.Errorf("SortBy(%q).IsValid() = %v, want %v", tt.sortBy, got, tt.want)
			}
		})
	}
}

func FuzzParseSortBy(f *testing.F) {
	f.Add("name")
	f.Add("importance")
	f.Add("created_at")
	f.Add("updated_at")
	f.Add("health")
	f.Add("complexity")
	f.Add("invalid")
	f.Add("")

	f.Fuzz(func(t *testing.T, s string) {
		sortBy, err := ParseSortBy(s)
		if err != nil {
			if sortBy != "" {
				t.Errorf("ParseSortBy(%q) returned error but non-empty sortBy: %q", s, sortBy)
			}
		}
		if sortBy.IsValid() && err == nil {
			if string(sortBy) != s {
				t.Errorf("ParseSortBy(%q) = %q, but IsValid() was true", s, sortBy)
			}
		}
	})
}
