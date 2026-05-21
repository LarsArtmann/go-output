package output

import (
	"testing"

	"github.com/larsartmann/go-output/testhelpers"
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
	testAllowedValues(
		t,
		"AllowedValues",
		SortByName.AllowedValues(),
		[]string{"name", "importance", "created_at", "updated_at", "health", "complexity"},
	)
}

func TestSortByIsValid(t *testing.T) {
	t.Parallel()

	testhelpers.TestEnumIsValid(t, []SortBy{
		SortByName,
		SortByImportance,
		SortByCreatedAt,
		SortByUpdatedAt,
		SortByHealth,
		SortByComplexity,
		SortBy("invalid"),
		SortBy(""),
	}, []bool{
		true,
		true,
		true,
		true,
		true,
		true,
		false,
		false,
	})
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
		fuzzEnumTest(t, s, ParseSortBy, "ParseSortBy")
	})
}
