package nom

import (
	"testing"
)

func TestSymbolConstants(t *testing.T) {
	t.Parallel()

	symbols := map[string]Symbol{
		"SymbolRunning":   SymbolRunning,
		"SymbolCompleted": SymbolCompleted,
		"SymbolPending":   SymbolPending,
		"SymbolFailed":    SymbolFailed,
		"SymbolDownload":  SymbolDownload,
		"SymbolUpload":    SymbolUpload,
		"SymbolAverage":   SymbolAverage,
		"SymbolTotal":     SymbolTotal,
	}

	for name, value := range symbols {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if value == "" {
				t.Errorf("%s should not be empty", name)
			}
		})
	}
}
