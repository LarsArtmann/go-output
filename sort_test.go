package output

import (
	"testing"
)

func TestParseSortBy(t *testing.T) {
	tests := []parseEnumTestCase[SortBy]{
		{"name", "name", SortByName, false},
		{"importance", "importance", SortByImportance, false},
		{"created_at", "created_at", SortByCreatedAt, false},
		{"updated_at", "updated_at", SortByUpdatedAt, false},
		{"health", "health", SortByHealth, false},
		{"complexity", "complexity", SortByComplexity, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}
	testParseEnum(t, "ParseSortBy", ParseSortBy, tests, func(a, b SortBy) bool { return a == b })
}

func TestSortByString(t *testing.T) {
	tests := []stringEnumTestCase[SortBy]{
		{SortByName, "name"},
		{SortByImportance, "importance"},
		{SortByCreatedAt, "created_at"},
		{SortByUpdatedAt, "updated_at"},
		{SortByHealth, "health"},
		{SortByComplexity, "complexity"},
	}
	testEnumString(t, "SortBy.String", tests, func(s SortBy) string { return s.String() })
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
